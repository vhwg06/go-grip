package admin

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
	usermodule "github.com/evrone/go-clean-template/internal/module/user"
)

// toAdminDashboardStatsResponse maps stats to openapi.AdminDashboardStatsResponse DTO.
func toAdminDashboardStatsResponse(users, orders, revenue, products int) openapi.AdminDashboardStatsResponse {
	return openapi.AdminDashboardStatsResponse{
		TotalUsers:    users,
		TotalOrders:   orders,
		TotalRevenue:  revenue,
		TotalProducts: products,
	}
}

// toAdminAuditLogResponse maps audit log entry to openapi.AdminAuditLogResponse DTO.
func toAdminAuditLogResponse(id, userID, action, resource string) openapi.AdminAuditLogResponse {
	uID := userID
	res := resource
	return openapi.AdminAuditLogResponse{
		Id:       id,
		UserId:   &uID,
		Action:   action,
		Resource: &res,
	}
}

// userStatusFromString maps a raw string status value to the typed UserStatus constant.
// Only "locked" maps to locked; all other values default to active.
func userStatusFromString(s string) usermodule.UserStatus {
	if s == "locked" {
		return usermodule.UserStatusLocked
	}
	return usermodule.UserStatusActive
}
