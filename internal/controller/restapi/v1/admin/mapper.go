package admin

import (
	"github.com/evrone/go-clean-template/api/gen/go/openapi"
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
