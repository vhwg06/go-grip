package middleware

import (
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
)

// RoleCompatibilityAdapter translates legacy User roles/claims into OpenAPI capability scopes.
type RoleCompatibilityAdapter struct {
	adminUsers map[string]struct{}
}

// NewRoleCompatibilityAdapter creates an adapter instance with configured admin emails/usernames.
func NewRoleCompatibilityAdapter(adminUsersCSV string) *RoleCompatibilityAdapter {
	adminMap := make(map[string]struct{})
	for _, user := range strings.Split(adminUsersCSV, ",") {
		trimmed := strings.TrimSpace(user)
		if trimmed != "" {
			adminMap[trimmed] = struct{}{}
		}
	}

	return &RoleCompatibilityAdapter{
		adminUsers: adminMap,
	}
}

// IsAdminUser checks whether the given actor is an administrator.
func (a *RoleCompatibilityAdapter) IsAdminUser(actor entity.Actor) bool {
	if actor.IsAdmin {
		return true
	}
	if actor.Email != "" {
		if _, ok := a.adminUsers[actor.Email]; ok {
			return true
		}
	}
	if actor.Username != "" {
		if _, ok := a.adminUsers[actor.Username]; ok {
			return true
		}
	}
	return false
}

// DeriveActorScopes computes the active capability scopes granted to an authenticated or anonymous actor.
func (a *RoleCompatibilityAdapter) DeriveActorScopes(actor entity.Actor) map[string]struct{} {
	scopes := make(map[string]struct{})

	// Blocked actors or anonymous actors get no scopes
	if actor.IsBlocked || (actor.UserID == "" && actor.Email == "" && actor.Username == "") {
		return scopes
	}

	// Standard User Capabilities granted to any authenticated user
	userScopes := []string{
		"auth:refresh",
		"auth:logout",
		"profile:read",
		"profile:write",
		"users:read",
		"tasks:read",
		"tasks:write",
		"translation:read",
		"translation:write",
		"orders:read",
		"orders:write",
		"wishlist:read",
		"wishlist:write",
		"reviews:read",
		"reviews:write",
		"notifications:read",
		"notifications:write",
		"cart:read",
		"cart:write",
	}

	for _, s := range userScopes {
		scopes[s] = struct{}{}
	}

	// Admin Capabilities granted if actor is Admin
	if a.IsAdminUser(actor) {
		adminScopes := []string{
			"users:write",
			"catalog:read",
			"catalog:write",
			"admin:read",
			"admin:write",
		}
		for _, s := range adminScopes {
			scopes[s] = struct{}{}
		}
	}

	return scopes
}

// HasCapability checks if an actor possesses all required capability scopes.
func (a *RoleCompatibilityAdapter) HasCapability(actor entity.Actor, requiredScopes []string) bool {
	if len(requiredScopes) == 0 {
		return true
	}

	if actor.IsBlocked {
		return false
	}

	actorScopes := a.DeriveActorScopes(actor)

	for _, req := range requiredScopes {
		if _, has := actorScopes[req]; !has {
			return false
		}
	}

	return true
}
