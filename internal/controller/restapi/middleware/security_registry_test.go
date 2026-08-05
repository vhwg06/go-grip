package middleware_test

import (
	"sync"
	"testing"

	"github.com/evrone/go-clean-template/internal/controller/restapi/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityRegistryImmutabilityAndDefensiveCopy(t *testing.T) {
	registry := middleware.NewSecurityRegistry()

	roles := []string{"admin", "editor"}
	scopes := []string{"users:read", "users:write"}

	policy := middleware.SecurityPolicy{
		OperationID: "getUserByID",
		IsPublic:    false,
		RequireJWT:  true,
		Roles:       roles,
		Scopes:      scopes,
		RateLimit:   10,
	}

	registry.RegisterPolicy(policy)

	// Mutate original slice to test defensive copy during registration
	roles[0] = "hacked"
	scopes[0] = "hacked"

	retrievedPolicy, found := registry.GetPolicy("getUserByID")
	require.True(t, found)
	assert.Equal(t, "admin", retrievedPolicy.Roles[0])
	assert.Equal(t, "users:read", retrievedPolicy.Scopes[0])

	// Freeze registry
	specYAML := []byte(`
openapi: 3.0.3
info:
  title: Test API
  version: 1.0.0
paths:
  /users/{id}:
    get:
      operationId: getUserByID
      responses:
        '200':
          description: OK
`)
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(specYAML)
	require.NoError(t, err)

	err = registry.BuildAndFreeze(doc)
	require.NoError(t, err)
	assert.True(t, registry.IsFrozen())

	// Test panic on RegisterPolicy after freeze
	assert.Panics(t, func() {
		registry.RegisterPolicy(middleware.SecurityPolicy{
			OperationID: "newOp",
		})
	}, "Must panic when registering policy after freeze")
}

func TestSecurityRegistryExhaustiveStartupGuard(t *testing.T) {
	registry := middleware.NewSecurityRegistry()
	registry.RegisterPolicy(middleware.SecurityPolicy{
		OperationID: "getUserByID",
		RequireJWT:  true,
	})

	specWithMissingOpYAML := []byte(`
openapi: 3.0.3
info:
  title: Test API
  version: 1.0.0
paths:
  /users/{id}:
    get:
      operationId: getUserByID
      responses:
        '200':
          description: OK
    delete:
      operationId: deleteUserByID
      responses:
        '204':
          description: No Content
`)

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(specWithMissingOpYAML)
	require.NoError(t, err)

	err = registry.BuildAndFreeze(doc)
	require.Error(t, err, "Must return error when an operationId in spec is missing from registry")
	assert.Contains(t, err.Error(), "exhaustive startup guard failed: generated operationId 'deleteUserByID'")
}

func TestSecurityRegistryConcurrentAccess(t *testing.T) {
	registry := middleware.NewSecurityRegistry()
	registry.RegisterPolicy(middleware.SecurityPolicy{
		OperationID: "getUserByID",
		RequireJWT:  true,
		Roles:       []string{"admin"},
		Scopes:      []string{"users:read"},
	})

	err := registry.BuildAndFreeze(nil)
	require.NoError(t, err)

	var wg sync.WaitGroup
	const goroutines = 50
	const iterations = 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				p, found := registry.GetPolicy("getUserByID")
				assert.True(t, found)
				assert.Equal(t, "getUserByID", p.OperationID)
				assert.Equal(t, "admin", p.Roles[0])
				assert.Equal(t, "users:read", p.Scopes[0])

				// Ensure caller mutating returned policy doesn't corrupt registry
				p.Roles[0] = "mutated"
				p.Scopes[0] = "mutated"
			}
		}()
	}

	wg.Wait()

	// Verify registry internal state was untouched by concurrent callers
	finalPolicy, found := registry.GetPolicy("getUserByID")
	require.True(t, found)
	assert.Equal(t, "admin", finalPolicy.Roles[0])
	assert.Equal(t, "users:read", finalPolicy.Scopes[0])
}
