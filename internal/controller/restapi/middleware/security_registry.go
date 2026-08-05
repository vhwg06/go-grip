package middleware

import (
	"fmt"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
)

// SecurityPolicy defines authorization and security rules for a specific operationId.
type SecurityPolicy struct {
	OperationID string
	IsPublic    bool
	RequireJWT  bool
	Roles       []string
	Scopes      []string
	RateLimit   int
}

// Clone returns a deep copy of the SecurityPolicy to prevent external mutation.
func (p SecurityPolicy) Clone() SecurityPolicy {
	rolesCopy := make([]string, len(p.Roles))
	copy(rolesCopy, p.Roles)

	scopesCopy := make([]string, len(p.Scopes))
	copy(scopesCopy, p.Scopes)

	return SecurityPolicy{
		OperationID: p.OperationID,
		IsPublic:    p.IsPublic,
		RequireJWT:  p.RequireJWT,
		Roles:       rolesCopy,
		Scopes:      scopesCopy,
		RateLimit:   p.RateLimit,
	}
}

// SecurityRegistry manages the immutable security policy registry for Generated API operations.
type SecurityRegistry struct {
	mu       sync.RWMutex
	policies map[string]SecurityPolicy
	isFrozen bool
}

// NewSecurityRegistry creates a new, unfrozen SecurityRegistry instance.
func NewSecurityRegistry() *SecurityRegistry {
	return &SecurityRegistry{
		policies: make(map[string]SecurityPolicy),
	}
}

// RegisterPolicy registers a security policy for a specific operationId.
// Panics if called after the registry has been frozen.
func (r *SecurityRegistry) RegisterPolicy(policy SecurityPolicy) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isFrozen {
		panic("cannot register policy: SecurityRegistry is frozen and immutable")
	}

	if policy.OperationID == "" {
		panic("cannot register policy with empty operationId")
	}

	r.policies[policy.OperationID] = policy.Clone()
}

// BuildAndFreeze verifies that all generated operations in the OpenAPI spec have registered policies,
// then freezes the registry to make it immutable. Panics or returns error if exhaustive check fails.
func (r *SecurityRegistry) BuildAndFreeze(doc *openapi3.T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isFrozen {
		return nil
	}

	// Exhaustive Startup Guard Check: Verify every generated operationId in spec is registered
	if doc != nil && doc.Paths != nil {
		for _, pathItem := range doc.Paths.Map() {
			if pathItem == nil {
				continue
			}

			ops := []*openapi3.Operation{
				pathItem.Get,
				pathItem.Post,
				pathItem.Put,
				pathItem.Patch,
				pathItem.Delete,
			}

			for _, op := range ops {
				if op == nil || op.OperationID == "" {
					continue
				}

				if _, exists := r.policies[op.OperationID]; !exists {
					return fmt.Errorf("exhaustive startup guard failed: generated operationId '%s' in spec is missing from SecurityRegistry", op.OperationID)
				}
			}
		}
	}

	r.isFrozen = true
	return nil
}

// IsFrozen returns true if the registry is frozen.
func (r *SecurityRegistry) IsFrozen() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.isFrozen
}

// GetPolicy retrieves a deep copy of the SecurityPolicy for the given operationId.
// Safe for concurrent access across multiple goroutines.
func (r *SecurityRegistry) GetPolicy(operationID string) (SecurityPolicy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	policy, exists := r.policies[operationID]
	if !exists {
		return SecurityPolicy{}, false
	}
	return policy.Clone(), true
}
