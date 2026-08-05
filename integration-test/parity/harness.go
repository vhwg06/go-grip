package parity

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/controller/restapi"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/internal/usecase/catalogbase"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ParityTestCase represents a table-driven HTTP parity test case.
type ParityTestCase struct {
	Name                 string
	Method               string
	Path                 string
	Headers              map[string]string
	AuthToken            string
	Body                 interface{}
	ExpectedStatus       int
	ExpectedHeaders      map[string]string
	ExpectedBodyContains []string
	ExpectedJSONFields   map[string]interface{}
	CustomAssertions     func(t *testing.T, resp *http.Response, body []byte)
}

// TestHarness encapsulates test dependencies and fiber application.
type TestHarness struct {
	App         *fiber.App
	Config      *config.Config
	JWTManager  *jwt.Manager
	AuthUC      *StubAuth
	UserUC      *StubUser
	TaskUC      *StubTask
	CatalogUC   *StubCatalog
	CatalogBase *StubCatalogBase
	Logger      logger.Interface
}

// HarnessOption allows configuring the test harness.
type HarnessOption func(*TestHarness)

// WithAuthUC sets custom Auth usecase.
func WithAuthUC(uc *StubAuth) HarnessOption {
	return func(h *TestHarness) {
		h.AuthUC = uc
	}
}

// WithUserUC sets custom User usecase.
func WithUserUC(uc *StubUser) HarnessOption {
	return func(h *TestHarness) {
		h.UserUC = uc
	}
}

// NewTestHarness initializes a test harness with manual router wired.
func NewTestHarness(t *testing.T, opts ...HarnessOption) *TestHarness {
	t.Helper()

	cfg := &config.Config{}

	jwtManager := jwt.New("parity-test-secret", time.Hour)
	l := logger.New("error")

	h := &TestHarness{
		Config:      cfg,
		JWTManager:  jwtManager,
		AuthUC:      &StubAuth{},
		UserUC:      &StubUser{},
		TaskUC:      &StubTask{},
		CatalogUC:   &StubCatalog{},
		CatalogBase: &StubCatalogBase{},
		Logger:      l,
	}

	for _, opt := range opts {
		opt(h)
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	var noopUcase usecase.Translation
	var noopCheckout usecase.Checkout
	var noopOrders usecase.Orders
	var noopProfile usecase.Profile
	var noopAdmin usecase.Admin
	var noopWishlist usecase.Wishlist
	var noopNotify usecase.NotificationCenter
	var noopMedia usecase.Media
	var noopHomepage usecase.Homepage
	var noopCart usecase.Cart
	var noopLead usecase.Lead
	var noopContent usecase.Content
	var noopImporter usecase.Importer

	restapi.NewRouter(
		app,
		cfg,
		noopUcase,
		h.UserUC,
		h.TaskUC,
		h.CatalogUC,
		h.CatalogBase,
		h.AuthUC,
		noopCheckout,
		noopOrders,
		noopProfile,
		noopAdmin,
		noopWishlist,
		noopNotify,
		noopMedia,
		noopHomepage,
		noopCart,
		noopLead,
		noopContent,
		noopImporter,
		jwtManager,
		l,
	)

	h.App = app
	return h
}

// CreateAuthToken generates a signed JWT token for test requests.
func (h *TestHarness) CreateAuthToken(userID, username string, isAdmin bool) (string, error) {
	return h.JWTManager.GenerateTokenWithProfile(userID, username, isAdmin)
}

// RunParityTest executes a single table-driven parity test case against the Fiber app.
func (h *TestHarness) RunParityTest(t *testing.T, tc ParityTestCase) {
	t.Helper()

	t.Run(tc.Name, func(t *testing.T) {
		var bodyReader io.Reader
		if tc.Body != nil {
			switch v := tc.Body.(type) {
			case string:
				bodyReader = bytes.NewBufferString(v)
			case []byte:
				bodyReader = bytes.NewBuffer(v)
			default:
				jsonBytes, err := json.Marshal(tc.Body)
				require.NoError(t, err, "Failed to marshal request body")
				bodyReader = bytes.NewBuffer(jsonBytes)
			}
		}

		req := httptest.NewRequest(tc.Method, tc.Path, bodyReader)
		req.Header.Set("X-Playwright-Test", "true")

		if tc.Body != nil && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		for k, v := range tc.Headers {
			req.Header.Set(k, v)
		}

		if tc.AuthToken != "" {
			req.Header.Set("Authorization", "Bearer "+tc.AuthToken)
		}

		resp, err := h.App.Test(req, 5000)
		require.NoError(t, err, "app.Test execution failed")
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err, "Failed to read response body")

		require.Equal(t, tc.ExpectedStatus, resp.StatusCode,
			"Status code mismatch for %s %s. Response body: %s", tc.Method, tc.Path, string(respBody))

		for k, v := range tc.ExpectedHeaders {
			assert.Equal(t, v, resp.Header.Get(k), "Header mismatch for key %s", k)
		}

		for _, fragment := range tc.ExpectedBodyContains {
			assert.Contains(t, string(respBody), fragment, "Response body missing expected fragment")
		}

		if len(tc.ExpectedJSONFields) > 0 {
			var jsonResp map[string]interface{}
			err := json.Unmarshal(respBody, &jsonResp)
			require.NoError(t, err, "Response body is not valid JSON: %s", string(respBody))

			for field, expectedVal := range tc.ExpectedJSONFields {
				assert.Equal(t, expectedVal, jsonResp[field], "JSON field mismatch for key %s", field)
			}
		}

		if tc.CustomAssertions != nil {
			tc.CustomAssertions(t, resp, respBody)
		}
	})
}

// Lightweight Usecase Stubs for Parity Testing Harness

type StubAuth struct {
	LoginFn   func(ctx context.Context, email, password string) (entity.User, string, string, error)
	RefreshFn func(ctx context.Context, refreshToken string) (string, string, error)
	LogoutFn  func(ctx context.Context, actor entity.Actor, refreshToken string) error
	MeFn      func(ctx context.Context, actor entity.Actor) (entity.User, error)
}

func (s *StubAuth) Login(ctx context.Context, email, password string) (entity.User, string, string, error) {
	if s.LoginFn != nil {
		return s.LoginFn(ctx, email, password)
	}
	return entity.User{}, "", "", entity.ErrInvalidCredentials
}

func (s *StubAuth) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	if s.RefreshFn != nil {
		return s.RefreshFn(ctx, refreshToken)
	}
	return "", "", entity.ErrUnauthorized
}

func (s *StubAuth) Logout(ctx context.Context, actor entity.Actor, refreshToken string) error {
	if s.LogoutFn != nil {
		return s.LogoutFn(ctx, actor, refreshToken)
	}
	return nil
}

func (s *StubAuth) Me(ctx context.Context, actor entity.Actor) (entity.User, error) {
	if s.MeFn != nil {
		return s.MeFn(ctx, actor)
	}
	return entity.User{ID: actor.UserID, Email: actor.Email}, nil
}

type StubUser struct {
	RegisterFn        func(ctx context.Context, username, email, password string) (entity.User, error)
	LoginFn           func(ctx context.Context, email, password string) (string, error)
	GetUserFn         func(ctx context.Context, userID string) (entity.User, error)
	ListFn            func(ctx context.Context, limit, offset int) ([]entity.User, int, error)
	CreateAdminUserFn func(ctx context.Context, actorID, username, email, password string, role entity.RoleName) (entity.User, error)
	UpdateProfileFn   func(ctx context.Context, userID, displayName string) (entity.User, error)
	LockFn            func(ctx context.Context, actorID, userID string) error
	UnlockFn          func(ctx context.Context, actorID, userID string) error
}

func (s *StubUser) Register(ctx context.Context, username, email, password string) (entity.User, error) {
	if s.RegisterFn != nil {
		return s.RegisterFn(ctx, username, email, password)
	}
	return entity.User{ID: "user-123", Username: username, Email: email, Role: entity.RoleSubscriber}, nil
}

func (s *StubUser) Login(ctx context.Context, email, password string) (string, error) {
	if s.LoginFn != nil {
		return s.LoginFn(ctx, email, password)
	}
	return "mock-jwt-token", nil
}

func (s *StubUser) GetUser(ctx context.Context, userID string) (entity.User, error) {
	if s.GetUserFn != nil {
		return s.GetUserFn(ctx, userID)
	}
	return entity.User{ID: userID, Email: "user@example.com"}, nil
}

func (s *StubUser) List(ctx context.Context, limit, offset int) ([]entity.User, int, error) {
	if s.ListFn != nil {
		return s.ListFn(ctx, limit, offset)
	}
	return []entity.User{{ID: "user-1", Email: "user1@example.com"}}, 1, nil
}

func (s *StubUser) CreateAdminUser(ctx context.Context, actorID, username, email, password string, role entity.RoleName) (entity.User, error) {
	if s.CreateAdminUserFn != nil {
		return s.CreateAdminUserFn(ctx, actorID, username, email, password, role)
	}
	return entity.User{ID: "admin-1", Username: username, Email: email, Role: role}, nil
}

func (s *StubUser) UpdateProfile(ctx context.Context, userID, displayName string) (entity.User, error) {
	if s.UpdateProfileFn != nil {
		return s.UpdateProfileFn(ctx, userID, displayName)
	}
	return entity.User{ID: userID, DisplayName: displayName}, nil
}

func (s *StubUser) Lock(ctx context.Context, actorID, userID string) error {
	if s.LockFn != nil {
		return s.LockFn(ctx, actorID, userID)
	}
	return nil
}

func (s *StubUser) Unlock(ctx context.Context, actorID, userID string) error {
	if s.UnlockFn != nil {
		return s.UnlockFn(ctx, actorID, userID)
	}
	return nil
}

type StubTask struct{}

func (s *StubTask) Create(ctx context.Context, userID, title, description string) (entity.Task, error) {
	return entity.Task{ID: "task-1", UserID: userID, Title: title, Description: description}, nil
}
func (s *StubTask) Get(ctx context.Context, userID, taskID string) (entity.Task, error) {
	return entity.Task{ID: taskID, UserID: userID, Title: "Task 1"}, nil
}
func (s *StubTask) List(ctx context.Context, userID string, status *entity.TaskStatus, limit, offset int) ([]entity.Task, int, error) {
	return []entity.Task{{ID: "task-1", UserID: userID, Title: "Task 1"}}, 1, nil
}
func (s *StubTask) Update(ctx context.Context, userID, taskID, title, description string) (entity.Task, error) {
	return entity.Task{ID: taskID, UserID: userID, Title: title, Description: description}, nil
}
func (s *StubTask) Transition(ctx context.Context, userID, taskID string, newStatus entity.TaskStatus) (entity.Task, error) {
	return entity.Task{ID: taskID, UserID: userID, Status: newStatus}, nil
}
func (s *StubTask) Delete(ctx context.Context, userID, taskID string) error {
	return nil
}

type StubCatalog struct{}

func (s *StubCatalog) CreateProduct(ctx context.Context, product entity.Product) (entity.Product, error) {
	return product, nil
}
func (s *StubCatalog) ListProducts(ctx context.Context, filter entity.ProductFilter) ([]entity.Product, int, error) {
	return []entity.Product{}, 0, nil
}
func (s *StubCatalog) GetProduct(ctx context.Context, id string) (entity.Product, error) {
	return entity.Product{ID: id}, nil
}
func (s *StubCatalog) UpdateProduct(ctx context.Context, product entity.Product) (entity.Product, error) {
	return product, nil
}
func (s *StubCatalog) DeleteProduct(ctx context.Context, id string) error {
	return nil
}
func (s *StubCatalog) CreateCategory(ctx context.Context, category entity.Category) (entity.Category, error) {
	return category, nil
}
func (s *StubCatalog) ListCategories(ctx context.Context) ([]entity.Category, error) {
	return []entity.Category{}, nil
}
func (s *StubCatalog) CreateTag(ctx context.Context, tag entity.Tag) (entity.Tag, error) {
	return tag, nil
}
func (s *StubCatalog) ListTags(ctx context.Context) ([]entity.Tag, error) {
	return []entity.Tag{}, nil
}

type StubCatalogBase struct{ catalogbase.UseCase }
