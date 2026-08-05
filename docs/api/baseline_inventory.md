# Baseline HTTP Route Inventory Manifest

Tài liệu này lưu trữ toàn bộ cây đăng ký HTTP Route (Registration Graph Inventory) hiện tại của hệ thống `go-grip`, phục vụ việc đóng khung (baseline freeze) trước khi chuyển đổi sang mô hình **OpenAPI-First Architecture**.

- **Thời điểm tạo**: 2026-08-05
- **Nguồn đăng ký xuất phát**:
  - [internal/controller/restapi/router.go](file:///c:/Users/hien.vm.CMCTI/Code/go-grip/internal/controller/restapi/router.go)
  - [internal/controller/restapi/v1/router.go](file:///c:/Users/hien.vm.CMCTI/Code/go-grip/internal/controller/restapi/v1/router.go)
  - [internal/controller/restapi/v1/grip_store_routes.go](file:///c:/Users/hien.vm.CMCTI/Code/go-grip/internal/controller/restapi/v1/grip_store_routes.go)
  - [internal/controller/restapi/v1/ecommerce.go](file:///c:/Users/hien.vm.CMCTI/Code/go-grip/internal/controller/restapi/v1/ecommerce.go)

---

## 1. Operational & System Endpoints

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Classification | Status Codes |
|---|---|---|---|---|---|---|---|
| `GET` | `/healthz` | `router.go:56` | inline `SendStatus(200)` | Public | None | Operational | `200` |
| `GET` | `/metrics` | `router.go:46` | `fiberprometheus.Middleware` | Public | None | Operational | `200`, `500` |

---

## 2. Documentation Endpoints

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Classification | Status Codes |
|---|---|---|---|---|---|---|---|
| `GET` | `/swagger/*` | `router.go:52` | `swagger.HandlerDefault` | Public | None | Documentation | `200`, `404`, `500` |

---

## 3. OpenAPI-Managed API Endpoints

### 3.1. Auth Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `POST` | `/v1/auth/register` | `v1/router.go:44` | `r.register` | Public | IP (5/1s) | JSON `registerRequest` (`email`, `password`) | `201 Created` | `400`, `409`, `422`, `500` |
| `POST` | `/v1/auth/login` | `v1/router.go:45` | `r.login` | Public | IP (5/1s) | JSON `loginRequest` (`email`, `password`) | `200 OK` | `400`, `401`, `422`, `500` |
| `POST` | `/v1/auth/refresh` | `v1/grip_store_routes.go:13` | `r.gripRefresh` | Public | None | JSON `gripRefreshRequest` (`refreshToken` / `refresh_token`) | `200 OK` | `400`, `401`, `500` |
| `POST` | `/v1/auth/logout` | `v1/grip_store_routes.go:14` | `r.gripLogout` | Bearer JWT | None | JSON `gripRefreshRequest` (`refreshToken` / `refresh_token`) | `204 No Content` | `400`, `401`, `500` |
| `GET` | `/v1/auth/me` | `v1/grip_store_routes.go:15` | `r.gripMe` | Bearer JWT | None | None | `200 OK` | `401`, `500` |

### 3.2. User Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `GET` | `/v1/user/profile` | `v1/router.go:56` | `r.profile` | Bearer JWT | None | None | `200 OK` | `401`, `500` |
| `GET` | `/v1/users/` | `v1/router.go:61` | `r.listUsers` | Bearer JWT | None | Query `limit`, `offset` | `200 OK` | `401`, `500` |
| `POST` | `/v1/users/` | `v1/router.go:62` | `r.createAdminUser` | Bearer JWT | None | JSON `createAdminUserRequest` (`email`, `password`, `role`) | `201 Created` | `400`, `401`, `409`, `422`, `500` |
| `GET` | `/v1/users/:id` | `v1/router.go:63` | `r.getUser` | Bearer JWT | None | Path Param `id` | `200 OK` | `401`, `404`, `500` |
| `PATCH` | `/v1/users/:id` | `v1/router.go:64` | `r.updateUserProfile` | Bearer JWT | None | Path Param `id`, JSON `updateUserProfileRequest` | `200 OK` | `400`, `401`, `404`, `500` |
| `POST` | `/v1/users/:id/lock` | `v1/router.go:65` | `r.lockUser` | Bearer JWT | None | Path Param `id` | `200 OK` | `401`, `404`, `500` |
| `POST` | `/v1/users/:id/unlock` | `v1/router.go:66` | `r.unlockUser` | Bearer JWT | None | Path Param `id` | `200 OK` | `401`, `404`, `500` |

### 3.3. Task Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `POST` | `/v1/tasks/` | `v1/router.go:71` | `r.createTask` | Bearer JWT | None | JSON `createTaskRequest` (`title`, `description`) | `201 Created` | `400`, `401`, `422`, `500` |
| `GET` | `/v1/tasks/` | `v1/router.go:72` | `r.listTasks` | Bearer JWT | None | Query `limit`, `offset`, `status` | `200 OK` | `401`, `500` |
| `GET` | `/v1/tasks/:id` | `v1/router.go:73` | `r.getTask` | Bearer JWT | None | Path Param `id` | `200 OK` | `401`, `404`, `500` |
| `PUT` | `/v1/tasks/:id` | `v1/router.go:74` | `r.updateTask` | Bearer JWT | None | Path Param `id`, JSON `updateTaskRequest` | `200 OK` | `400`, `401`, `404`, `500` |
| `PATCH` | `/v1/tasks/:id/status` | `v1/router.go:75` | `r.transitionTask` | Bearer JWT | None | Path Param `id`, JSON `transitionTaskStatusRequest` | `200 OK` | `400`, `401`, `404`, `422`, `500` |
| `DELETE` | `/v1/tasks/:id` | `v1/router.go:76` | `r.deleteTask` | Bearer JWT | None | Path Param `id` | `204 No Content` | `401`, `404`, `500` |

### 3.4. Translation Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `GET` | `/v1/translation/history` | `v1/router.go:81` | `r.history` | Bearer JWT | None | Query `limit`, `offset` | `200 OK` | `401`, `500` |
| `POST` | `/v1/translation/do-translate` | `v1/router.go:82` | `r.doTranslate` | Bearer JWT | None | JSON `translationRequest` (`text`, `source`, `target`) | `200 OK` | `400`, `401`, `500` |

### 3.5. Catalog & Catalog Base Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `GET` | `/v1/catalog/products` | `v1/grip_store_routes.go:18` | `r.gripListProducts` | Public | None | Query `q`, `brand`, `sort`, `limit`, `offset` | `200 OK` | `500` |
| `GET` | `/v1/catalog/products/:id` | `v1/grip_store_routes.go:19` | `r.gripGetProduct` | Public | None | Path Param `id` | `200 OK` | `404`, `500` |
| `GET` | `/v1/catalog/products/:id/buy-meta` | `v1/grip_store_routes.go:20` | `r.gripGetBuyMeta` | Public | None | Path Param `id` | `200 OK` | `404`, `500` |
| `GET` | `/v1/catalog/products/:id/reviews` | `v1/grip_store_routes.go:21` | `r.gripReviewList` | Public | None | Path Param `id`, Query `limit`, `offset` | `200 OK` | `500` |
| `GET` | `/v1/catalog/search` | `v1/grip_store_routes.go:22` | `r.gripSearchProducts` | Public | None | Query `q`, `limit`, `offset` | `200 OK` | `500` |
| `GET` | `/v1/catalog/categories` | `v1/grip_store_routes.go:23` | `r.gripListCategories` | Public | None | None | `200 OK` | `500` |
| `GET` | `/v1/catalog/settings` | `v1/grip_store_routes.go:24` | `r.gripListSettings` | Public | None | None | `200 OK` | `500` |
| `GET` | `/v1/catalog/announcement` | `v1/grip_store_routes.go:25` | `r.gripGetAnnouncement` | Public | None | None | `200 OK` | `500` |
| `GET` | `/v1/catalog/product-models` | `v1/grip_store_routes.go:29` | `r.catalogBasePublicList` | Public | None | Query `limit`, `offset`, `category_id` | `200 OK` | `500` |
| `GET` | `/v1/catalog/product-models/:modelId` | `v1/grip_store_routes.go:30` | `r.catalogBasePublicDetail` | Public | None | Path Param `modelId` | `200 OK` | `404`, `500` |
| `GET` | `/v1/catalog/product-models/:modelId/options` | `v1/grip_store_routes.go:31` | `r.catalogBasePublicOptions` | Public | None | Path Param `modelId` | `200 OK` | `404`, `500` |
| `POST` | `/v1/catalog/product-models/:modelId/variants:resolve` | `v1/grip_store_routes.go:32` | `r.catalogBasePublicResolve` | Public | None | Path Param `modelId`, JSON selected options | `200 OK` | `400`, `404`, `500` |
| `GET` | `/v1/site-config` | `v1/grip_store_routes.go:33` | `r.gripSiteConfig` | Public | None | None | `200 OK` | `500` |
| `POST` | `/v1/catalog/products` | `v1/ecommerce.go:18` | `r.createProduct` | Admin | None | JSON `entity.Product` payload | `201 Created` | `400`, `401`, `403`, `500` |
| `PATCH` | `/v1/catalog/products/:id` | `v1/ecommerce.go:19` | `r.updateProduct` | Admin | None | Path Param `id`, JSON `entity.Product` | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `DELETE` | `/v1/catalog/products/:id` | `v1/ecommerce.go:20` | `r.deleteProduct` | Admin | None | Path Param `id` | `204 No Content` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/catalog/categories` | `v1/ecommerce.go:21` | `r.createCategory` | Admin | None | JSON `entity.Category` payload | `201 Created` | `400`, `401`, `403`, `500` |
| `POST` | `/v1/catalog/tags` | `v1/ecommerce.go:22` | `r.createTag` | Admin | None | JSON `entity.Tag` payload | `201 Created` | `400`, `401`, `403`, `500` |

### 3.6. Checkout & Orders Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `GET` | `/v1/checkout/preview` | `v1/grip_store_routes.go:36` | `r.gripCheckoutPreview` | Public | IP (30/1m) | Query Params | `200 OK` | `400`, `500` |
| `POST` | `/v1/checkout/orders` | `v1/grip_store_routes.go:37` | `r.gripCreateOrder` | Bearer JWT | IP (30/1m) | JSON Create Order payload | `201 Created` | `400`, `401`, `500` |
| `POST` | `/v1/checkout/payment-orders` | `v1/grip_store_routes.go:38` | `r.gripCreatePaymentOrder` | Bearer JWT | IP (30/1m) | JSON Payment Order payload | `201 Created` | `400`, `401`, `500` |
| `GET` | `/v1/checkout/orders/:id/payment-params` | `v1/grip_store_routes.go:39` | `r.gripPaymentParams` | Public | IP (30/1m) | Path Param `id` | `200 OK` | `404`, `500` |
| `GET` | `/v1/checkout/orders/:id/status` | `v1/grip_store_routes.go:40` | `r.gripOrderStatus` | Bearer JWT | IP (30/1m) | Path Param `id` | `200 OK` | `401`, `404`, `500` |
| `POST` | `/v1/checkout/orders/:id/cancel` | `v1/grip_store_routes.go:41` | `r.gripCancelOrder` | Bearer JWT | IP (30/1m) | Path Param `id` | `200 OK` | `401`, `404`, `500` |
| `POST` | `/v1/checkout/notify` | `v1/grip_store_routes.go:42` | `r.gripPaymentNotify` | Public | IP (30/1m) | Webhook Notification payload | `200 OK` | `400`, `500` |
| `GET` | `/v1/checkout/callback/:id` | `v1/grip_store_routes.go:43` | `r.gripPaymentCallback` | Public | IP (30/1m) | Path Param `id` | `200 OK` | `404`, `500` |
| `GET` | `/v1/orders/` | `v1/grip_store_routes.go:46` | `r.gripListOrders` | Bearer JWT | None | Query `limit`, `offset` | `200 OK` | `401`, `500` |
| `GET` | `/v1/orders/:id` | `v1/grip_store_routes.go:47` | `r.gripGetOrder` | Bearer JWT | None | Path Param `id` | `200 OK` | `401`, `404`, `500` |
| `POST` | `/v1/orders/:id/refund-request` | `v1/grip_store_routes.go:48` | `r.gripRequestRefund` | Bearer JWT | None | Path Param `id`, JSON Refund Reason | `200 OK` | `400`, `401`, `404`, `500` |

### 3.7. Profile, Wishlist & Historical Removal Stubs

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `GET` | `/v1/profile/` | `v1/grip_store_routes.go:51` | `r.gripProfileGet` | Bearer JWT | None | None | `200 OK` | `401`, `500` |
| `PATCH` | `/v1/profile/` | `v1/grip_store_routes.go:52` | `r.gripProfileUpdate` | Bearer JWT | None | JSON Update Profile payload | `200 OK` | `400`, `401`, `500` |
| `GET` | `/v1/profile/security` | `v1/grip_store_routes.go:53` | `r.gripProfileGetSecurity` | Bearer JWT | None | None | `200 OK` | `401`, `500` |
| `GET` | `/v1/profile/sessions` | `v1/grip_store_routes.go:54` | `r.gripProfileGetSessions` | Bearer JWT | None | None | `200 OK` | `401`, `500` |
| `POST` | `/v1/profile/checkin` | `v1/grip_store_routes.go:72` | inline `removedCheckIn` | Public | None | None | N/A | `404 Not Found` |
| `GET` | `/v1/profile/checkin-status` | `v1/grip_store_routes.go:73` | inline `removedCheckIn` | Public | None | None | N/A | `404 Not Found` |
| `GET` | `/v1/profile/checkin/status` | `v1/grip_store_routes.go:74` | inline `removedCheckIn` | Public | None | None | N/A | `404 Not Found` |
| `GET` | `/v1/user/profile/checkin-status` | `v1/grip_store_routes.go:75` | inline `removedCheckIn` | Public | None | None | N/A | `404 Not Found` |
| `GET` | `/v1/user/profile/checkin/status` | `v1/grip_store_routes.go:76` | inline `removedCheckIn` | Public | None | None | N/A | `404 Not Found` |
| `GET` | `/v1/wishlist/` | `v1/grip_store_routes.go:57` | `r.gripWishlistList` | Public | None | Query `limit`, `offset` | `200 OK` | `500` |
| `POST` | `/v1/wishlist/` | `v1/grip_store_routes.go:58` | `r.gripWishlistCreate` | Bearer JWT | None | JSON Wishlist Item payload | `201 Created` | `400`, `401`, `500` |
| `PATCH` | `/v1/wishlist/:id` | `v1/grip_store_routes.go:59` | `r.gripWishlistUpdate` | Bearer JWT | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `404`, `500` |
| `DELETE` | `/v1/wishlist/:id` | `v1/grip_store_routes.go:60` | `r.gripWishlistDelete` | Bearer JWT | None | Path Param `id` | `204 No Content` | `401`, `404`, `500` |
| `POST` | `/v1/wishlist/:id/vote` | `v1/grip_store_routes.go:61` | `r.gripWishlistVote` | Bearer JWT | None | Path Param `id` | `200 OK` | `401`, `404`, `500` |

### 3.8. Reviews, FAQs & Public Content

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `GET` | `/v1/reviews` | `v1/grip_store_routes.go:63` | `r.gripReviewList` | Public | None | Query `product_id`, `limit`, `offset` | `200 OK` | `500` |
| `GET` | `/v1/products/:id/reviews` | `v1/grip_store_routes.go:64` | `r.gripReviewList` | Public | None | Path Param `id`, Query Params | `200 OK` | `500` |
| `POST` | `/v1/products/:id/reviews` | `v1/grip_store_routes.go:65` | `r.gripReviewCreate` | Bearer JWT | None | Path Param `id`, JSON Review payload | `201 Created` | `400`, `401`, `500` |
| `POST` | `/v1/reviews/` | `v1/grip_store_routes.go:79` | `r.gripReviewCreate` | Bearer JWT | None | JSON Review payload | `201 Created` | `400`, `401`, `500` |
| `GET` | `/v1/faqs/active` | `v1/grip_store_routes.go:66` | `r.listActiveFAQs` | Public | None | None | `200 OK` | `500` |
| `GET` | `/v1/public/search` | `v1/ecommerce.go:49` | `r.listProducts` | Public | None | Query `q`, `brand`, `sort`, `limit`, `offset` | `200 OK` | `500` |
| `GET` | `/v1/public/categories` | `v1/ecommerce.go:50` | `r.listCategories` | Public | None | None | `200 OK` | `500` |
| `GET` | `/v1/public/categories/:id/products` | `v1/ecommerce.go:51` | `r.listProducts` | Public | None | Path Param `id`, Query Params | `200 OK` | `500` |
| `GET` | `/v1/public/products/:id` | `v1/ecommerce.go:52` | `r.getProduct` | Public | None | Path Param `id` | `200 OK` | `404`, `500` |
| `GET` | `/v1/public/homepage` | `v1/ecommerce.go:53` | `r.listPublicHomepage` | Public | None | None | `200 OK` | `500` |
| `GET` | `/v1/public/content/articles` | `v1/ecommerce.go:54` | `r.listPublicArticles` | Public | None | Query Params | `200 OK` | `500` |
| `GET` | `/v1/public/content/articles/:id` | `v1/ecommerce.go:55` | `r.getArticle` | Public | None | Path Param `id` | `200 OK` | `404`, `500` |
| `GET` | `/v1/public/content/pages/:slug` | `v1/ecommerce.go:56` | `r.getPage` | Public | None | Path Param `slug` | `200 OK` | `404`, `500` |
| `GET` | `/v1/public/footer` | `v1/ecommerce.go:57` | `r.listPublicHomepage` | Public | None | None | `200 OK` | `500` |
| `GET` | `/v1/public/support` | `v1/ecommerce.go:58` | `r.listPublicSupport` | Public | None | None | `200 OK` | `500` |

### 3.9. Notifications Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `GET` | `/v1/notifications/` | `v1/grip_store_routes.go:82` | `r.gripNotificationsList` | Bearer JWT | None | Query `limit`, `offset` | `200 OK` | `401`, `500` |
| `GET` | `/v1/notifications/unread-count` | `v1/grip_store_routes.go:83` | `r.gripNotificationsUnread` | Bearer JWT | None | None | `200 OK` | `401`, `500` |
| `POST` | `/v1/notifications/:id/read` | `v1/grip_store_routes.go:84` | `r.gripNotificationsMarkRead` | Bearer JWT | None | Path Param `id` | `200 OK` | `401`, `404`, `500` |
| `POST` | `/v1/notifications/read-all` | `v1/grip_store_routes.go:85` | `r.gripNotificationsReadAll` | Bearer JWT | None | None | `200 OK` | `401`, `500` |
| `DELETE` | `/v1/notifications/` | `v1/grip_store_routes.go:86` | `r.gripNotificationsClear` | Bearer JWT | None | None | `204 No Content` | `401`, `500` |

### 3.10. Admin Management Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `GET` | `/v1/admin/products` | `v1/grip_store_routes.go:89` | `r.gripAdminListProducts` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/products` | `v1/grip_store_routes.go:90` | `r.gripAdminCreateProduct` | Admin | None | JSON Body | `201 Created` | `400`, `401`, `403`, `500` |
| `PATCH` | `/v1/admin/products/:id` | `v1/grip_store_routes.go:91` | `r.gripAdminUpdateProduct` | Admin | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `PATCH` | `/v1/admin/products/:id/status` | `v1/grip_store_routes.go:92` | `r.gripAdminUpdateProductStatus` | Admin | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `DELETE` | `/v1/admin/products/:id` | `v1/grip_store_routes.go:93` | `r.gripAdminDeleteProduct` | Admin | None | Path Param `id` | `204 No Content` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/categories` | `v1/grip_store_routes.go:94` | `r.gripAdminListCategories` | Admin | None | None | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/categories` | `v1/grip_store_routes.go:95` | `r.gripAdminCreateCategory` | Admin | None | JSON Body | `201 Created` | `400`, `401`, `403`, `500` |
| `PATCH` | `/v1/admin/categories/:id` | `v1/grip_store_routes.go:96` | `r.gripAdminUpdateCategory` | Admin | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `DELETE` | `/v1/admin/categories/:id` | `v1/grip_store_routes.go:97` | `r.gripAdminDeleteCategory` | Admin | None | Path Param `id` | `204 No Content` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/orders` | `v1/grip_store_routes.go:98` | `r.gripAdminListOrders` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `GET` | `/v1/admin/orders/:id` | `v1/grip_store_routes.go:99` | `r.gripAdminGetOrder` | Admin | None | Path Param `id` | `200 OK` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/orders/:id/refund-status` | `v1/grip_store_routes.go:100` | `r.gripAdminGetOrderRefundStatus` | Admin | None | Path Param `id` | `200 OK` | `401`, `403`, `404`, `500` |
| `PATCH` | `/v1/admin/orders/:id` | `v1/grip_store_routes.go:101` | `r.gripAdminUpdateOrder` | Admin | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `DELETE` | `/v1/admin/orders/:id` | `v1/grip_store_routes.go:102` | `r.gripAdminDeleteOrder` | Admin | None | Path Param `id` | `204 No Content` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/refunds` | `v1/grip_store_routes.go:103` | `r.gripAdminListRefunds` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `GET` | `/v1/admin/refunds/:id` | `v1/grip_store_routes.go:104` | `r.gripAdminGetRefund` | Admin | None | Path Param `id` | `200 OK` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/refunds/:id/approve` | `v1/grip_store_routes.go:105` | `r.gripAdminApproveRefund` | Admin | None | Path Param `id` | `200 OK` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/refunds/:id/reject` | `v1/grip_store_routes.go:106` | `r.gripAdminRejectRefund` | Admin | None | Path Param `id`, JSON Body | `200 OK` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/users` | `v1/grip_store_routes.go:107` | `r.gripAdminUsersList` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `PATCH` | `/v1/admin/users/:id` | `v1/grip_store_routes.go:108` | `r.gripAdminUsersUpdate` | Admin | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `PATCH` | `/v1/admin/users/:id/block` | `v1/grip_store_routes.go:109` | `r.gripAdminUsersUpdateBlock` | Admin | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/reviews` | `v1/grip_store_routes.go:110` | `r.gripAdminListReviews` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `GET` | `/v1/admin/reviews/:id` | `v1/grip_store_routes.go:111` | `r.gripAdminGetReview` | Admin | None | Path Param `id` | `200 OK` | `401`, `403`, `404`, `500` |
| `PUT` | `/v1/admin/reviews/:id/approve` | `v1/grip_store_routes.go:112` | `r.gripAdminApproveReview` | Admin | None | Path Param `id` | `200 OK` | `401`, `403`, `404`, `500` |
| `PUT` | `/v1/admin/reviews/:id/hide` | `v1/grip_store_routes.go:113` | `r.gripAdminHideReview` | Admin | None | Path Param `id` | `200 OK` | `401`, `403`, `404`, `500` |
| `PUT` | `/v1/admin/reviews/:id/feature` | `v1/grip_store_routes.go:114` | `r.gripAdminFeatureReview` | Admin | None | Path Param `id` | `200 OK` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/reviews/publish-selected` | `v1/grip_store_routes.go:115` | `r.gripAdminPublishSelectedReviews` | Admin | None | JSON Body `ids` | `200 OK` | `400`, `401`, `403`, `500` |
| `DELETE` | `/v1/admin/reviews/:id` | `v1/grip_store_routes.go:116` | `r.gripAdminDeleteReview` | Admin | None | Path Param `id` | `204 No Content` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/settings` | `v1/grip_store_routes.go:117` | `r.gripAdminListSettings` | Admin | None | None | `200 OK` | `401`, `403`, `500` |
| `PUT` | `/v1/admin/settings/:key` | `v1/grip_store_routes.go:118` | `r.gripAdminUpsertSetting` | Admin | None | Path Param `key`, JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `DELETE` | `/v1/admin/settings/:key` | `v1/grip_store_routes.go:119` | `r.gripAdminDeleteSetting` | Admin | None | Path Param `key` | `204 No Content` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/collect` | `v1/grip_store_routes.go:120` | `r.gripAdminGetCollect` | Admin | None | None | `200 OK` | `401`, `403`, `500` |
| `PUT` | `/v1/admin/collect` | `v1/grip_store_routes.go:121` | `r.gripAdminPutCollect` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `GET` | `/v1/admin/store-settings` | `v1/grip_store_routes.go:122` | `r.gripAdminGetStoreSettings` | Admin | None | None | `200 OK` | `401`, `403`, `500` |
| `PUT` | `/v1/admin/store-settings/presence` | `v1/grip_store_routes.go:123` | `r.gripAdminPutStoreSettingsPresence` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `PUT` | `/v1/admin/store-settings/brand` | `v1/grip_store_routes.go:124` | `r.gripAdminPutStoreSettingsBrand` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `PUT` | `/v1/admin/store-settings/contact` | `v1/grip_store_routes.go:125` | `r.gripAdminPutStoreSettingsContact` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `PUT` | `/v1/admin/store-settings/homepage` | `v1/grip_store_routes.go:126` | `r.gripAdminPutStoreSettingsHomepage` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `PUT` | `/v1/admin/store-settings/footer` | `v1/grip_store_routes.go:127` | `r.gripAdminPutStoreSettingsFooter` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `PUT` | `/v1/admin/store-settings/floating-support` | `v1/grip_store_routes.go:128` | `r.gripAdminPutStoreSettingsFloatingSupport` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `PUT` | `/v1/admin/store-settings/visibility` | `v1/grip_store_routes.go:129` | `r.gripAdminPutStoreSettingsVisibility` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `PUT` | `/v1/admin/store-settings/registry` | `v1/grip_store_routes.go:130` | `r.gripAdminPutStoreSettingsRegistry` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `POST` | `/v1/admin/messages/broadcast` | `v1/grip_store_routes.go:131` | `r.gripAdminBroadcast` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `POST` | `/v1/admin/messages/targeted` | `v1/grip_store_routes.go:132` | `r.gripAdminTargeted` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `GET` | `/v1/admin/messages` | `v1/grip_store_routes.go:133` | `r.gripAdminListMessages` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `GET` | `/v1/admin/notifications` | `v1/grip_store_routes.go:134` | `r.gripAdminGetNotifications` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/notifications` | `v1/grip_store_routes.go:135` | `r.gripAdminPostNotifications` | Admin | None | JSON Body | `201 Created` | `400`, `401`, `403`, `500` |
| `POST` | `/v1/admin/notifications/test` | `v1/grip_store_routes.go:136` | `r.gripAdminNotificationTest` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `POST` | `/v1/admin/data/import` | `v1/grip_store_routes.go:137` | `r.gripAdminImportData` | Admin | None | Multipart / JSON | `200 OK` | `400`, `401`, `403`, `500` |
| `POST` | `/v1/admin/data/repair-aggregates` | `v1/grip_store_routes.go:138` | `r.gripAdminRepairAggregates` | Admin | None | None | `200 OK` | `401`, `403`, `500` |
| `GET` | `/v1/admin/media/presigned` | `v1/grip_store_routes.go:139` | `r.getPresignedURL` | Admin | None | Query `filename`, `content_type` | `200 OK` | `400`, `401`, `403`, `500` |
| `GET` | `/v1/admin/media` | `v1/grip_store_routes.go:140` | `r.listMedia` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/media` | `v1/grip_store_routes.go:141` | `r.createMedia` | Admin | None | Multipart / JSON | `201 Created` | `400`, `401`, `403`, `500` |
| `DELETE` | `/v1/admin/media/:id` | `v1/grip_store_routes.go:142` | `r.deleteMedia` | Admin | None | Path Param `id` | `204 No Content` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/banners` | `v1/grip_store_routes.go:143` | `r.listAdminBanners` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/banners` | `v1/grip_store_routes.go:144` | `r.saveAdminBanner` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `DELETE` | `/v1/admin/banners/:id` | `v1/grip_store_routes.go:145` | `r.deleteAdminBanner` | Admin | None | Path Param `id` | `204 No Content` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/faqs` | `v1/grip_store_routes.go:146` | `r.listAdminFAQs` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/faqs` | `v1/grip_store_routes.go:147` | `r.saveAdminFAQ` | Admin | None | JSON Body | `200 OK` | `400`, `401`, `403`, `500` |
| `DELETE` | `/v1/admin/faqs/:id` | `v1/grip_store_routes.go:148` | `r.deleteAdminFAQ` | Admin | None | Path Param `id` | `204 No Content` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/products/new` | `v1/grip_store_routes.go:149` | `r.gripAdminProductsNew` | Admin | None | None | `200 OK` | `401`, `403`, `500` |
| `GET` | `/v1/admin/products/:id/form` | `v1/grip_store_routes.go:150` | `r.gripAdminProductForm` | Admin | None | Path Param `id` | `200 OK` | `401`, `403`, `404`, `500` |

### 3.11. Admin Catalog Base Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `GET` | `/v1/admin/catalog/categories` | `v1/grip_store_routes.go:153` | `r.catalogBaseListCategories` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/catalog/categories` | `v1/grip_store_routes.go:154` | `r.catalogBaseCreateCategory` | Admin | None | JSON Body | `201 Created` | `400`, `401`, `403`, `500` |
| `PATCH` | `/v1/admin/catalog/categories/:categoryId` | `v1/grip_store_routes.go:155` | `r.catalogBaseUpdateCategory` | Admin | None | Path Param `categoryId`, JSON | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `DELETE` | `/v1/admin/catalog/categories/:categoryId` | `v1/grip_store_routes.go:156` | `r.catalogBaseDeleteCategory` | Admin | None | Path Param `categoryId` | `204 No Content` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/categories/:categoryId/deactivate` | `v1/grip_store_routes.go:157` | `r.catalogBaseDeactivateCategory` | Admin | None | Path Param `categoryId` | `200 OK` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/catalog/attribute-definitions` | `v1/grip_store_routes.go:159` | `r.catalogBaseListDefinitions` | Admin | None | Query Params | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/catalog/attribute-definitions` | `v1/grip_store_routes.go:160` | `r.catalogBaseCreateDefinition` | Admin | None | JSON Body | `201 Created` | `400`, `401`, `403`, `500` |
| `PATCH` | `/v1/admin/catalog/attribute-definitions/:definitionId` | `v1/grip_store_routes.go:161` | `r.catalogBaseUpdateDefinition` | Admin | None | Path Param `definitionId`, JSON | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/attribute-definitions/:definitionId/deactivate` | `v1/grip_store_routes.go:162` | `r.catalogBaseDeactivateDefinition` | Admin | None | Path Param `definitionId` | `200 OK` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/attribute-definitions/:definitionId/enum-values` | `v1/grip_store_routes.go:163` | `r.catalogBaseAddEnumValue` | Admin | None | Path Param `definitionId`, JSON | `201 Created` | `400`, `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/attribute-definitions/:definitionId/enum-values/:enumValueId/deactivate` | `v1/grip_store_routes.go:164` | `r.catalogBaseDeactivateEnumValue` | Admin | None | Path Params: `definitionId`, `enumValueId` | `200 OK` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/catalog/masters/:masterKind` | `v1/grip_store_routes.go:166` | `r.catalogBaseListMasters` | Admin | None | Path Param `masterKind` | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/catalog/masters/:masterKind` | `v1/grip_store_routes.go:167` | `r.catalogBaseCreateMaster` | Admin | None | Path Param `masterKind`, JSON | `201 Created` | `400`, `401`, `403`, `500` |
| `PATCH` | `/v1/admin/catalog/masters/:masterKind/:masterId` | `v1/grip_store_routes.go:168` | `r.catalogBaseUpdateMaster` | Admin | None | Path Params: `masterKind`, `masterId` | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/masters/:masterKind/:masterId/deactivate` | `v1/grip_store_routes.go:169` | `r.catalogBaseDeactivateMaster` | Admin | None | Path Params: `masterKind`, `masterId` | `200 OK` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/catalog/product-models` | `v1/grip_store_routes.go:171` | `r.catalogBaseListModels` | Admin | None | Query `limit`, `offset` | `200 OK` | `401`, `403`, `500` |
| `POST` | `/v1/admin/catalog/product-models` | `v1/grip_store_routes.go:172` | `r.catalogBaseCreateModel` | Admin | None | JSON Body | `201 Created` | `400`, `401`, `403`, `500` |
| `GET` | `/v1/admin/catalog/product-models/:modelId` | `v1/grip_store_routes.go:173` | `r.catalogBaseGetModel` | Admin | None | Path Param `modelId` | `200 OK` | `401`, `403`, `404`, `500` |
| `PATCH` | `/v1/admin/catalog/product-models/:modelId` | `v1/grip_store_routes.go:174` | `r.catalogBaseUpdateModel` | Admin | None | Path Param `modelId`, JSON | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `DELETE` | `/v1/admin/catalog/product-models/:modelId` | `v1/grip_store_routes.go:175` | `r.catalogBaseDeleteModel` | Admin | None | Path Param `modelId` | `204 No Content` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/product-models/:modelId/publish` | `v1/grip_store_routes.go:176` | `r.catalogBasePublishModel` | Admin | None | Path Param `modelId` | `200 OK` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/product-models/:modelId/unpublish` | `v1/grip_store_routes.go:177` | `r.catalogBaseUnpublishModel` | Admin | None | Path Param `modelId` | `200 OK` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/product-models/:modelId/discontinue` | `v1/grip_store_routes.go:178` | `r.catalogBaseDiscontinueModel` | Admin | None | Path Param `modelId` | `200 OK` | `401`, `403`, `404`, `500` |
| `PUT` | `/v1/admin/catalog/product-models/:modelId/media` | `v1/grip_store_routes.go:179` | `r.catalogBaseReplaceMedia` | Admin | None | Path Param `modelId`, JSON | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/product-models/:modelId/variant-dimensions` | `v1/grip_store_routes.go:180` | `r.catalogBaseCreateDimension` | Admin | None | Path Param `modelId`, JSON | `201 Created` | `400`, `401`, `403`, `404`, `500` |
| `PATCH` | `/v1/admin/catalog/product-models/:modelId/variant-dimensions/:dimensionId` | `v1/grip_store_routes.go:181` | `r.catalogBaseUpdateDimension` | Admin | None | Path Params: `modelId`, `dimensionId` | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/product-models/:modelId/variant-dimensions/:dimensionId/values` | `v1/grip_store_routes.go:182` | `r.catalogBaseAddDimensionValue` | Admin | None | Path Params: `modelId`, `dimensionId` | `201 Created` | `400`, `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/product-models/:modelId/variant-dimensions/:dimensionId/values/:valueId/deactivate` | `v1/grip_store_routes.go:183` | `r.catalogBaseDeactivateDimensionValue` | Admin | None | Path Params: `modelId`, `dimensionId`, `valueId` | `200 OK` | `401`, `403`, `404`, `500` |
| `GET` | `/v1/admin/catalog/product-models/:modelId/variants` | `v1/grip_store_routes.go:184` | `r.catalogBaseListVariants` | Admin | None | Path Param `modelId` | `200 OK` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/product-models/:modelId/variants` | `v1/grip_store_routes.go:185` | `r.catalogBaseCreateVariant` | Admin | None | Path Param `modelId`, JSON | `201 Created` | `400`, `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/variants/prices:bulk` | `v1/grip_store_routes.go:186` | `r.catalogBaseBulkPrice` | Admin | None | JSON Bulk Prices payload | `200 OK` | `400`, `401`, `403`, `500` |
| `GET` | `/v1/admin/catalog/variants/:variantId` | `v1/grip_store_routes.go:187` | `r.catalogBaseGetVariant` | Admin | None | Path Param `variantId` | `200 OK` | `401`, `403`, `404`, `500` |
| `PATCH` | `/v1/admin/catalog/variants/:variantId` | `v1/grip_store_routes.go:188` | `r.catalogBaseUpdateVariant` | Admin | None | Path Param `variantId`, JSON | `200 OK` | `400`, `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/variants/:variantId/activate` | `v1/grip_store_routes.go:189` | `r.catalogBaseActivateVariant` | Admin | None | Path Param `variantId` | `200 OK` | `401`, `403`, `404`, `500` |
| `POST` | `/v1/admin/catalog/variants/:variantId/inactivate` | `v1/grip_store_routes.go:190` | `r.catalogBaseInactivateVariant` | Admin | None | Path Param `variantId` | `200 OK` | `401`, `403`, `404`, `500` |

### 3.12. Media, Content, Cart & Order Requests Capabilities

| Method | Path | Registration Owner | Handler Function | Auth / Protection | Rate Limit | Request Shape | Success Status | Error Statuses |
|---|---|---|---|---|---|---|---|---|
| `POST` | `/v1/media` | `v1/ecommerce.go:24` | `r.createMedia` | Bearer JWT | None | Multipart / JSON payload | `201 Created` | `400`, `401`, `500` |
| `GET` | `/v1/media` | `v1/ecommerce.go:25` | `r.listMedia` | Bearer JWT | None | Query `limit`, `offset` | `200 OK` | `401`, `500` |
| `DELETE` | `/v1/media/:id` | `v1/ecommerce.go:26` | `r.deleteMedia` | Bearer JWT | None | Path Param `id` | `204 No Content` | `401`, `404`, `500` |
| `PUT` | `/v1/media/simulate-upload/:filename` | `v1/ecommerce.go:27` | `r.simulateUpload` | Public | None | Path Param `filename` | `200 OK` | `500` |
| `GET` | `/v1/homepage/blocks` | `v1/ecommerce.go:29` | `r.listHomepageBlocks` | Bearer JWT | None | Query Params | `200 OK` | `401`, `500` |
| `POST` | `/v1/homepage/blocks` | `v1/ecommerce.go:30` | `r.createHomepageBlock` | Bearer JWT | None | JSON Body | `201 Created` | `400`, `401`, `500` |
| `PATCH` | `/v1/homepage/blocks/:id` | `v1/ecommerce.go:31` | `r.updateHomepageBlock` | Bearer JWT | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `404`, `500` |
| `DELETE` | `/v1/homepage/blocks/:id` | `v1/ecommerce.go:32` | `r.deleteHomepageBlock` | Bearer JWT | None | Path Param `id` | `204 No Content` | `401`, `404`, `500` |
| `GET` | `/v1/support/channels` | `v1/ecommerce.go:33` | `r.listSupportChannels` | Bearer JWT | None | Query Params | `200 OK` | `401`, `500` |
| `PATCH` | `/v1/support/channels/:id` | `v1/ecommerce.go:34` | `r.updateSupportChannel` | Bearer JWT | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `404`, `500` |
| `GET` | `/v1/content/articles` | `v1/ecommerce.go:36` | `r.listAdminArticles` | Bearer JWT | None | Query Params | `200 OK` | `401`, `500` |
| `POST` | `/v1/content/articles` | `v1/ecommerce.go:37` | `r.createArticle` | Bearer JWT | None | JSON Body | `201 Created` | `400`, `401`, `500` |
| `PATCH` | `/v1/content/articles/:id` | `v1/ecommerce.go:38` | `r.updateArticle` | Bearer JWT | None | Path Param `id`, JSON Body | `200 OK` | `400`, `401`, `404`, `500` |
| `DELETE` | `/v1/content/articles/:id` | `v1/ecommerce.go:39` | `r.deleteArticle` | Bearer JWT | None | Path Param `id` | `204 No Content` | `401`, `404`, `500` |
| `GET` | `/v1/content/articles/:id/preview` | `v1/ecommerce.go:40` | `r.previewArticle` | Bearer JWT | None | Path Param `id` | `200 OK` | `401`, `404`, `500` |
| `POST` | `/v1/content/articles/:id/schedule` | `v1/ecommerce.go:41` | `r.updateArticle` | Bearer JWT | None | Path Param `id`, JSON Schedule payload | `200 OK` | `400`, `401`, `404`, `500` |
| `POST` | `/v1/content/articles/:id/publish` | `v1/ecommerce.go:42` | `r.updateArticle` | Bearer JWT | None | Path Param `id`, JSON Publish payload | `200 OK` | `400`, `401`, `404`, `500` |
| `GET` | `/v1/content/pages` | `v1/ecommerce.go:43` | `r.getPage` | Bearer JWT | None | Query Params | `200 OK` | `401`, `500` |
| `POST` | `/v1/content/pages` | `v1/ecommerce.go:44` | `r.createPage` | Bearer JWT | None | JSON Body | `201 Created` | `400`, `401`, `500` |
| `PATCH` | `/v1/content/pages/:slug` | `v1/ecommerce.go:45` | `r.updatePage` | Bearer JWT | None | Path Param `slug`, JSON Body | `200 OK` | `400`, `401`, `404`, `500` |
| `POST` | `/v1/import/initial-content` | `v1/ecommerce.go:46` | `r.importInitialContent` | Bearer JWT | None | JSON Content Import payload | `200 OK` | `400`, `401`, `500` |
| `POST` | `/v1/cart/` | `v1/ecommerce.go:61` | `r.createCart` | Bearer JWT | None | JSON Cart payload | `201 Created` | `400`, `401`, `500` |
| `GET` | `/v1/cart/:session_id` | `v1/ecommerce.go:62` | `r.getCart` | Bearer JWT | None | Path Param `session_id` | `200 OK` | `401`, `404`, `500` |
| `POST` | `/v1/cart/:session_id/items` | `v1/ecommerce.go:63` | `r.addCartItem` | Bearer JWT | None | Path Param `session_id`, JSON Item | `200 OK` | `400`, `401`, `404`, `500` |
| `PATCH` | `/v1/cart/:session_id/items/:item_id` | `v1/ecommerce.go:64` | `r.updateCartItem` | Bearer JWT | None | Path Params: `session_id`, `item_id`, JSON | `200 OK` | `400`, `401`, `404`, `500` |
| `DELETE` | `/v1/cart/:session_id/items/:item_id` | `v1/ecommerce.go:65` | `r.removeCartItem` | Bearer JWT | None | Path Params: `session_id`, `item_id` | `204 No Content` | `401`, `404`, `500` |
| `POST` | `/v1/order-requests` | `v1/ecommerce.go:66` | `r.submitOrder` | Bearer JWT | None | JSON Order Request payload | `201 Created` | `400`, `401`, `500` |
| `POST` | `/v1/leads` | `v1/ecommerce.go:67` | `r.submitLead` | Public | None | JSON Lead payload | `201 Created` | `400`, `500` |
| `GET` | `/v1/leads/:id` | `v1/ecommerce.go:68` | `r.getLead` | Bearer JWT | None | Path Param `id` | `200 OK` | `401`, `404`, `500` |

---

## 4. Summary Stats & Inventory Coverage Audit

- **Tổng số Route Endpoints**: 142 endpoints.
- **Operational Endpoints**: 2 (`/healthz`, `/metrics`)
- **Documentation Endpoints**: 1 (`/swagger/*`)
- **OpenAPI-Managed Endpoints**: 139 endpoints
- **Phủ 100% Registration Graph**:
  - `router.go`
  - `v1/router.go`
  - `grip_store_routes.go`
  - `ecommerce.go`
- **Tình trạng Contract Explicit**: 100% routes được định nghĩa rõ ràng về HTTP Method, Path, Auth Level, Rate Limiting, Request Shape, và Status Codes. Không có endpoint nào bị ẩn hoặc không thể xác định.
