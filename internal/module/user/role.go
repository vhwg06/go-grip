package user

// RoleName represents backoffice permission levels.
type RoleName string

const (
	RoleAdministrator RoleName = "Administrator"
	RoleEditor        RoleName = "Editor"
	RoleAuthor        RoleName = "Author"
	RoleContributor   RoleName = "Contributor"
	RoleSubscriber    RoleName = "Subscriber"
)

// Role defines an administrative permission grouping.
type Role struct {
	ID   string   `json:"id"`
	Name RoleName `json:"name"`
}

// CanChangeRoles reports whether this role can assign roles or lock users.
func (r RoleName) CanChangeRoles() bool {
	return r == RoleAdministrator
}
