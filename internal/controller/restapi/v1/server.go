package v1

import (
	"context"

	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	"github.com/evrone/go-clean-template/config"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/admin"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/auth"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/cart"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/catalog"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/checkout"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/content"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/importer"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/lead"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/media"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/notification"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/orders"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/profile"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/reviews"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/user"
	"github.com/evrone/go-clean-template/internal/controller/restapi/v1/wishlist"
	cartmodule "github.com/evrone/go-clean-template/internal/module/cart"
	catalogmodule "github.com/evrone/go-clean-template/internal/module/catalog"
	"github.com/evrone/go-clean-template/internal/module/catalog/catalogbase"
	importermodule "github.com/evrone/go-clean-template/internal/module/importer"
	leadmodule "github.com/evrone/go-clean-template/internal/module/lead"
	mediamodule "github.com/evrone/go-clean-template/internal/module/media"
	ordermodule "github.com/evrone/go-clean-template/internal/module/order"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
	wishlistmodule "github.com/evrone/go-clean-template/internal/module/wishlist"
	notificationmodule "github.com/evrone/go-clean-template/internal/module/notification"
	contentmodule "github.com/evrone/go-clean-template/internal/module/content"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
)

var _ openapi.StrictServerInterface = (*Server)(nil)

// Server is the delivery layer composition root for OpenAPI generated handlers.
type Server struct {
	cfg         *config.Config
	userUC      usermodule.UserUseCase
	catalogUC   catalogmodule.CatalogUseCase
	catalogBase catalogbase.UseCase
	authUC      usermodule.AuthUseCase
	checkoutUC  ordermodule.CheckoutUseCase
	ordersUC    ordermodule.OrdersUseCase
	profileUC   usermodule.ProfileUseCase
	adminUC     usermodule.AdminUseCase
	wishlistUC  wishlistmodule.WishlistUseCase
	notifyUC    notificationmodule.NotificationCenterUseCase
	mediaUC     mediamodule.MediaUseCase
	homepageUC  contentmodule.HomepageUseCase
	cartUC      cartmodule.CartUseCase
	leadUC      leadmodule.LeadUseCase
	contentUC   contentmodule.ContentUseCase
	importerUC  importermodule.ImporterUseCase
	jwtManager  *jwt.Manager
	logger      logger.Interface
}

// NewServer constructs the OpenAPI Server composition root.
func NewServer(
	cfg *config.Config,
	userUC usermodule.UserUseCase,
	catalogUC catalogmodule.CatalogUseCase,
	catalogBase catalogbase.UseCase,
	authUC usermodule.AuthUseCase,
	checkoutUC ordermodule.CheckoutUseCase,
	ordersUC ordermodule.OrdersUseCase,
	profileUC usermodule.ProfileUseCase,
	adminUC usermodule.AdminUseCase,
	wishlistUC wishlistmodule.WishlistUseCase,
	notifyUC notificationmodule.NotificationCenterUseCase,
	mediaUC mediamodule.MediaUseCase,
	homepageUC contentmodule.HomepageUseCase,
	cartUC cartmodule.CartUseCase,
	leadUC leadmodule.LeadUseCase,
	contentUC contentmodule.ContentUseCase,
	importerUC importermodule.ImporterUseCase,
	jwtManager *jwt.Manager,
	l logger.Interface,
) *Server {
	return &Server{
		cfg:         cfg,
		userUC:      userUC,
		catalogUC:   catalogUC,
		catalogBase: catalogBase,
		authUC:      authUC,
		checkoutUC:  checkoutUC,
		ordersUC:    ordersUC,
		profileUC:   profileUC,
		adminUC:     adminUC,
		wishlistUC:  wishlistUC,
		notifyUC:    notifyUC,
		mediaUC:     mediaUC,
		homepageUC:  homepageUC,
		cartUC:      cartUC,
		leadUC:      leadUC,
		contentUC:   contentUC,
		importerUC:  importerUC,
		jwtManager:  jwtManager,
		logger:      l,
	}
}

// Handler instantiations for vertical capabilities

func (s *Server) adminHandler() *admin.Handler {
	return admin.NewHandler(s.adminUC, s.logger)
}

func (s *Server) authHandler() *auth.Handler {
	return auth.NewHandler(s.authUC, s.userUC, s.logger)
}

func (s *Server) cartHandler() *cart.Handler {
	return cart.NewHandler(s.cartUC, s.logger)
}

func (s *Server) catalogHandler() *catalog.Handler {
	return catalog.NewHandler(s.catalogUC, s.logger)
}

func (s *Server) checkoutHandler() *checkout.Handler {
	return checkout.NewHandler(s.checkoutUC, s.logger)
}

func (s *Server) contentHandler() *content.Handler {
	return content.NewHandler(s.contentUC, s.homepageUC, s.logger)
}

func (s *Server) importerHandler() *importer.Handler {
	return importer.NewHandler(s.importerUC, s.logger)
}

func (s *Server) leadHandler() *lead.Handler {
	return lead.NewHandler(s.leadUC, s.logger)
}

func (s *Server) mediaHandler() *media.Handler {
	return media.NewHandler(s.mediaUC, s.logger)
}

func (s *Server) notificationHandler() *notification.Handler {
	return notification.NewHandler(s.notifyUC, s.logger)
}

func (s *Server) ordersHandler() *orders.Handler {
	return orders.NewHandler(s.ordersUC, s.logger)
}

func (s *Server) profileHandler() *profile.Handler {
	return profile.NewHandler(s.profileUC, s.logger)
}

func (s *Server) userHandler() *user.Handler {
	return user.NewHandler(s.userUC, s.logger)
}

func (s *Server) wishlistHandler() *wishlist.Handler {
	return wishlist.NewHandler(s.wishlistUC, s.logger)
}

func (s *Server) reviewsHandler() *reviews.Handler {
	return reviews.NewHandler(s.logger)
}

// -----------------------------------------------------------------------------
// OpenAPI StrictServerInterface Delegation Methods
// -----------------------------------------------------------------------------

// Admin Capability
func (s *Server) GetAdminDashboardStats(ctx context.Context, request openapi.GetAdminDashboardStatsRequestObject) (openapi.GetAdminDashboardStatsResponseObject, error) {
	return s.adminHandler().GetAdminDashboardStats(ctx, request)
}

func (s *Server) ListAdminAuditLogs(ctx context.Context, request openapi.ListAdminAuditLogsRequestObject) (openapi.ListAdminAuditLogsResponseObject, error) {
	return s.adminHandler().ListAdminAuditLogs(ctx, request)
}

// Auth Capability
func (s *Server) RegisterUser(ctx context.Context, request openapi.RegisterUserRequestObject) (openapi.RegisterUserResponseObject, error) {
	return s.authHandler().RegisterUser(ctx, request)
}

func (s *Server) LoginUser(ctx context.Context, request openapi.LoginUserRequestObject) (openapi.LoginUserResponseObject, error) {
	return s.authHandler().LoginUser(ctx, request)
}

func (s *Server) RefreshToken(ctx context.Context, request openapi.RefreshTokenRequestObject) (openapi.RefreshTokenResponseObject, error) {
	return s.authHandler().RefreshToken(ctx, request)
}

func (s *Server) LogoutUser(ctx context.Context, request openapi.LogoutUserRequestObject) (openapi.LogoutUserResponseObject, error) {
	return s.authHandler().LogoutUser(ctx, request)
}

func (s *Server) GetCurrentUser(ctx context.Context, request openapi.GetCurrentUserRequestObject) (openapi.GetCurrentUserResponseObject, error) {
	return s.authHandler().GetCurrentUser(ctx, request)
}

// Cart Capability
func (s *Server) GetMyCart(ctx context.Context, request openapi.GetMyCartRequestObject) (openapi.GetMyCartResponseObject, error) {
	return s.cartHandler().GetMyCart(ctx, request)
}

func (s *Server) AddToCart(ctx context.Context, request openapi.AddToCartRequestObject) (openapi.AddToCartResponseObject, error) {
	return s.cartHandler().AddToCart(ctx, request)
}

func (s *Server) UpdateCartItem(ctx context.Context, request openapi.UpdateCartItemRequestObject) (openapi.UpdateCartItemResponseObject, error) {
	return s.cartHandler().UpdateCartItem(ctx, request)
}

func (s *Server) RemoveCartItem(ctx context.Context, request openapi.RemoveCartItemRequestObject) (openapi.RemoveCartItemResponseObject, error) {
	return s.cartHandler().RemoveCartItem(ctx, request)
}

func (s *Server) ClearCart(ctx context.Context, request openapi.ClearCartRequestObject) (openapi.ClearCartResponseObject, error) {
	return s.cartHandler().ClearCart(ctx, request)
}

// Catalog Capability
func (s *Server) ListProducts(ctx context.Context, request openapi.ListProductsRequestObject) (openapi.ListProductsResponseObject, error) {
	return s.catalogHandler().ListProducts(ctx, request)
}

func (s *Server) CreateProduct(ctx context.Context, request openapi.CreateProductRequestObject) (openapi.CreateProductResponseObject, error) {
	return s.catalogHandler().CreateProduct(ctx, request)
}

func (s *Server) GetProductByID(ctx context.Context, request openapi.GetProductByIDRequestObject) (openapi.GetProductByIDResponseObject, error) {
	return s.catalogHandler().GetProductByID(ctx, request)
}

func (s *Server) UpdateProduct(ctx context.Context, request openapi.UpdateProductRequestObject) (openapi.UpdateProductResponseObject, error) {
	return s.catalogHandler().UpdateProduct(ctx, request)
}

func (s *Server) DeleteProduct(ctx context.Context, request openapi.DeleteProductRequestObject) (openapi.DeleteProductResponseObject, error) {
	return s.catalogHandler().DeleteProduct(ctx, request)
}

func (s *Server) ListCategories(ctx context.Context, request openapi.ListCategoriesRequestObject) (openapi.ListCategoriesResponseObject, error) {
	return s.catalogHandler().ListCategories(ctx, request)
}

func (s *Server) CreateCategory(ctx context.Context, request openapi.CreateCategoryRequestObject) (openapi.CreateCategoryResponseObject, error) {
	return s.catalogHandler().CreateCategory(ctx, request)
}

func (s *Server) ListTags(ctx context.Context, request openapi.ListTagsRequestObject) (openapi.ListTagsResponseObject, error) {
	return s.catalogHandler().ListTags(ctx, request)
}

func (s *Server) CreateTag(ctx context.Context, request openapi.CreateTagRequestObject) (openapi.CreateTagResponseObject, error) {
	return s.catalogHandler().CreateTag(ctx, request)
}

// Checkout Capability
func (s *Server) GetCheckoutPreview(ctx context.Context, request openapi.GetCheckoutPreviewRequestObject) (openapi.GetCheckoutPreviewResponseObject, error) {
	return s.checkoutHandler().GetCheckoutPreview(ctx, request)
}

func (s *Server) PreviewCheckout(ctx context.Context, request openapi.PreviewCheckoutRequestObject) (openapi.PreviewCheckoutResponseObject, error) {
	return s.checkoutHandler().PreviewCheckout(ctx, request)
}

func (s *Server) CreateCheckoutOrder(ctx context.Context, request openapi.CreateCheckoutOrderRequestObject) (openapi.CreateCheckoutOrderResponseObject, error) {
	return s.checkoutHandler().CreateCheckoutOrder(ctx, request)
}

func (s *Server) GetPaymentParams(ctx context.Context, request openapi.GetPaymentParamsRequestObject) (openapi.GetPaymentParamsResponseObject, error) {
	return s.checkoutHandler().GetPaymentParams(ctx, request)
}

func (s *Server) PaymentNotify(ctx context.Context, request openapi.PaymentNotifyRequestObject) (openapi.PaymentNotifyResponseObject, error) {
	return s.checkoutHandler().PaymentNotify(ctx, request)
}

func (s *Server) GetPaymentStatus(ctx context.Context, request openapi.GetPaymentStatusRequestObject) (openapi.GetPaymentStatusResponseObject, error) {
	return s.checkoutHandler().GetPaymentStatus(ctx, request)
}

// Content & Homepage Capability
func (s *Server) GetStaticPage(ctx context.Context, request openapi.GetStaticPageRequestObject) (openapi.GetStaticPageResponseObject, error) {
	return s.contentHandler().GetStaticPage(ctx, request)
}

func (s *Server) GetHomepageConfig(ctx context.Context, request openapi.GetHomepageConfigRequestObject) (openapi.GetHomepageConfigResponseObject, error) {
	return s.contentHandler().GetHomepageConfig(ctx, request)
}

// Importer Capability
func (s *Server) ExecuteImport(ctx context.Context, request openapi.ExecuteImportRequestObject) (openapi.ExecuteImportResponseObject, error) {
	return s.importerHandler().ExecuteImport(ctx, request)
}

// Lead Capability
func (s *Server) SubmitLead(ctx context.Context, request openapi.SubmitLeadRequestObject) (openapi.SubmitLeadResponseObject, error) {
	return s.leadHandler().SubmitLead(ctx, request)
}

func (s *Server) ListLeads(ctx context.Context, request openapi.ListLeadsRequestObject) (openapi.ListLeadsResponseObject, error) {
	return s.leadHandler().ListLeads(ctx, request)
}

// Media Capability
func (s *Server) UploadMedia(ctx context.Context, request openapi.UploadMediaRequestObject) (openapi.UploadMediaResponseObject, error) {
	return s.mediaHandler().UploadMedia(ctx, request)
}

func (s *Server) DeleteMedia(ctx context.Context, request openapi.DeleteMediaRequestObject) (openapi.DeleteMediaResponseObject, error) {
	return s.mediaHandler().DeleteMedia(ctx, request)
}

// Notification Capability
func (s *Server) ListNotifications(ctx context.Context, request openapi.ListNotificationsRequestObject) (openapi.ListNotificationsResponseObject, error) {
	return s.notificationHandler().ListNotifications(ctx, request)
}

func (s *Server) MarkAllNotificationsRead(ctx context.Context, request openapi.MarkAllNotificationsReadRequestObject) (openapi.MarkAllNotificationsReadResponseObject, error) {
	return s.notificationHandler().MarkAllNotificationsRead(ctx, request)
}

func (s *Server) GetUnreadNotificationCount(ctx context.Context, request openapi.GetUnreadNotificationCountRequestObject) (openapi.GetUnreadNotificationCountResponseObject, error) {
	return s.notificationHandler().GetUnreadNotificationCount(ctx, request)
}

func (s *Server) MarkNotificationRead(ctx context.Context, request openapi.MarkNotificationReadRequestObject) (openapi.MarkNotificationReadResponseObject, error) {
	return s.notificationHandler().MarkNotificationRead(ctx, request)
}

// Reviews Capability
func (s *Server) GetProductReviews(ctx context.Context, request openapi.GetProductReviewsRequestObject) (openapi.GetProductReviewsResponseObject, error) {
	return s.reviewsHandler().GetProductReviews(ctx, request)
}

func (s *Server) CreateReview(ctx context.Context, request openapi.CreateReviewRequestObject) (openapi.CreateReviewResponseObject, error) {
	return s.reviewsHandler().CreateReview(ctx, request)
}

func (s *Server) DeleteReview(ctx context.Context, request openapi.DeleteReviewRequestObject) (openapi.DeleteReviewResponseObject, error) {
	return s.reviewsHandler().DeleteReview(ctx, request)
}

// Orders Capability
func (s *Server) ListOrders(ctx context.Context, request openapi.ListOrdersRequestObject) (openapi.ListOrdersResponseObject, error) {
	return s.ordersHandler().ListOrders(ctx, request)
}

func (s *Server) GetOrderByID(ctx context.Context, request openapi.GetOrderByIDRequestObject) (openapi.GetOrderByIDResponseObject, error) {
	return s.ordersHandler().GetOrderByID(ctx, request)
}

func (s *Server) RequestOrderRefund(ctx context.Context, request openapi.RequestOrderRefundRequestObject) (openapi.RequestOrderRefundResponseObject, error) {
	return s.ordersHandler().RequestOrderRefund(ctx, request)
}

// Profile Capability
func (s *Server) GetAccountProfile(ctx context.Context, request openapi.GetAccountProfileRequestObject) (openapi.GetAccountProfileResponseObject, error) {
	return s.profileHandler().GetAccountProfile(ctx, request)
}

func (s *Server) UpdateAccountProfile(ctx context.Context, request openapi.UpdateAccountProfileRequestObject) (openapi.UpdateAccountProfileResponseObject, error) {
	return s.profileHandler().UpdateAccountProfile(ctx, request)
}

// User Capability
func (s *Server) GetMyProfile(ctx context.Context, request openapi.GetMyProfileRequestObject) (openapi.GetMyProfileResponseObject, error) {
	return s.userHandler().GetMyProfile(ctx, request)
}

func (s *Server) UpdateMyProfile(ctx context.Context, request openapi.UpdateMyProfileRequestObject) (openapi.UpdateMyProfileResponseObject, error) {
	return s.userHandler().UpdateMyProfile(ctx, request)
}

func (s *Server) ListUsers(ctx context.Context, request openapi.ListUsersRequestObject) (openapi.ListUsersResponseObject, error) {
	return s.userHandler().ListUsers(ctx, request)
}

func (s *Server) CreateAdminUser(ctx context.Context, request openapi.CreateAdminUserRequestObject) (openapi.CreateAdminUserResponseObject, error) {
	return s.userHandler().CreateAdminUser(ctx, request)
}

func (s *Server) GetUserByID(ctx context.Context, request openapi.GetUserByIDRequestObject) (openapi.GetUserByIDResponseObject, error) {
	return s.userHandler().GetUserByID(ctx, request)
}

func (s *Server) LockUser(ctx context.Context, request openapi.LockUserRequestObject) (openapi.LockUserResponseObject, error) {
	return s.userHandler().LockUser(ctx, request)
}

func (s *Server) UnlockUser(ctx context.Context, request openapi.UnlockUserRequestObject) (openapi.UnlockUserResponseObject, error) {
	return s.userHandler().UnlockUser(ctx, request)
}

// Wishlist Capability
func (s *Server) GetMyWishlist(ctx context.Context, request openapi.GetMyWishlistRequestObject) (openapi.GetMyWishlistResponseObject, error) {
	return s.wishlistHandler().GetMyWishlist(ctx, request)
}

func (s *Server) AddToWishlistDirect(ctx context.Context, request openapi.AddToWishlistDirectRequestObject) (openapi.AddToWishlistDirectResponseObject, error) {
	return s.wishlistHandler().AddToWishlistDirect(ctx, request)
}

func (s *Server) RemoveFromWishlistDirect(ctx context.Context, request openapi.RemoveFromWishlistDirectRequestObject) (openapi.RemoveFromWishlistDirectResponseObject, error) {
	return s.wishlistHandler().RemoveFromWishlistDirect(ctx, request)
}

func (s *Server) VoteWishlistItem(ctx context.Context, request openapi.VoteWishlistItemRequestObject) (openapi.VoteWishlistItemResponseObject, error) {
	return s.wishlistHandler().VoteWishlistItem(ctx, request)
}

func (s *Server) AddToWishlist(ctx context.Context, request openapi.AddToWishlistRequestObject) (openapi.AddToWishlistResponseObject, error) {
	return s.wishlistHandler().AddToWishlist(ctx, request)
}

func (s *Server) RemoveFromWishlist(ctx context.Context, request openapi.RemoveFromWishlistRequestObject) (openapi.RemoveFromWishlistResponseObject, error) {
	return s.wishlistHandler().RemoveFromWishlist(ctx, request)
}
