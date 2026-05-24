# REST API Contract: Grip Store Backend API

## Common Rules

- Base path: `/v1`
- Request and response body format: JSON unless noted for external payment callback compatibility.
- Authentication: `Authorization: Bearer <access_token>` for protected routes.
- Actor context is resolved by REST middleware and passed to controllers/usecases.
- Admin routes require both valid authentication and admin middleware approval.
- Mutating routes reject blocked users unless an admin is performing an allowed admin action.
- Error responses use a consistent shape:

```json
{
  "error": "machine.readable.code",
  "message": "Human readable message"
}
```

## Auth

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/v1/auth/register` | Public | Register a local account |
| POST | `/v1/auth/login` | Public | Authenticate with local credentials |
| POST | `/v1/auth/refresh` | Refresh credential | Rotate access and refresh credentials |
| POST | `/v1/auth/logout` | Bearer | Revoke current session |
| GET | `/v1/auth/me` | Bearer | Return current profile and admin flag |

## Catalog

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/v1/catalog/products` | Public | List visible active products |
| GET | `/v1/catalog/products/{id}` | Public | Product detail with stock and purchase limit |
| GET | `/v1/catalog/products/{id}/buy-meta` | Public | Reviews and buyer review eligibility |
| GET | `/v1/catalog/search` | Public | Search/filter/sort visible products |
| GET | `/v1/catalog/categories` | Public | List categories |
| GET | `/v1/catalog/settings` | Public | Read public store settings |
| GET | `/v1/catalog/announcement` | Public | Read active announcement |

Product detail responses include:

```json
{
  "id": "prod_123",
  "name": "Example card",
  "price": 10000,
  "displayStock": 12,
  "maxPurchaseableQuantity": 2,
  "rating": 4.8,
  "reviewCount": 20
}
```

## Checkout and Payment

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/v1/checkout/preview` | Bearer | Calculate price, points, and final payable amount |
| POST | `/v1/checkout/orders` | Bearer | Create order and reserve stock |
| POST | `/v1/checkout/payment-orders` | Bearer | Create direct payment order |
| GET | `/v1/checkout/orders/{id}/payment-params` | Bearer | Recreate payment instructions |
| GET | `/v1/checkout/orders/{id}/status` | Bearer | Poll order status |
| POST | `/v1/checkout/orders/{id}/cancel` | Bearer | Cancel pending order |
| POST | `/v1/checkout/notify` | Provider signature | Process payment confirmation |
| GET | `/v1/checkout/callback/{id}` | Public | Payment return redirect handler |

Checkout preview response:

```json
{
  "productId": "prod_123",
  "quantity": 2,
  "subtotal": 20000,
  "pointsToUse": 5000,
  "finalPrice": 15000
}
```

Order status responses include `statusText` and `statusColor` for client rendering.

## Cart

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/v1/cart` | Bearer | Create current user cart |
| GET | `/v1/cart/{session_id}` | Bearer | Read current user cart |
| POST | `/v1/cart/{session_id}/items` | Bearer | Add item to current user cart |
| PATCH | `/v1/cart/{session_id}/items/{item_id}` | Bearer | Update current user cart item quantity |
| DELETE | `/v1/cart/{session_id}/items/{item_id}` | Bearer | Remove item from current user cart |
| POST | `/v1/order-requests` | Bearer | Submit order request from current user cart |

The `session_id` path field is retained for compatibility; cart ownership is resolved from the authenticated actor.

## Orders

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/v1/orders` | Bearer | List owner orders |
| GET | `/v1/orders/{id}` | Owner/admin | Read order detail |
| POST | `/v1/orders/{id}/refund-request` | Bearer owner | Request refund for delivered order |

Card key fields are omitted or masked unless the order is delivered and the requester is authorized.

## Profile

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/v1/profile` | Bearer | Profile dashboard |
| PATCH | `/v1/profile` | Bearer | Update email and notification preference |
| POST | `/v1/profile/check-in` | Bearer | Daily check-in |

## Wishlist and Reviews

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/v1/wishlist` | Public | List wishlist items |
| POST | `/v1/wishlist` | Bearer | Create wishlist item |
| PATCH | `/v1/wishlist/{id}` | Owner/admin | Update wishlist item |
| DELETE | `/v1/wishlist/{id}` | Owner/admin | Delete wishlist item |
| POST | `/v1/wishlist/{id}/vote` | Bearer | Toggle vote |
| POST | `/v1/reviews` | Bearer | Create review for delivered order |

## Notifications

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/v1/notifications` | Bearer | Unified inbox |
| GET | `/v1/notifications/unread-count` | Bearer | Personal plus broadcast unread count |
| POST | `/v1/notifications/{id}/read` | Bearer | Mark one notification read |
| POST | `/v1/notifications/read-all` | Bearer | Mark all current messages read |
| DELETE | `/v1/notifications` | Bearer | Clear inbox |

## Admin

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET/POST/PATCH/DELETE | `/v1/admin/products` | Admin | Product management |
| GET/POST/PATCH/DELETE | `/v1/admin/categories` | Admin | Category management |
| GET/POST/DELETE | `/v1/admin/cards` | Admin | Card inventory management |
| POST | `/v1/admin/cards/import` | Admin | Bulk card import |
| POST | `/v1/admin/cards/replenish` | Admin | Pull cards from configured supplier |
| GET/PATCH/DELETE | `/v1/admin/orders` | Admin | Order management |
| GET | `/v1/admin/refunds` | Admin | List pending refunds |
| POST | `/v1/admin/refunds/{id}/approve` | Admin | Approve refund |
| POST | `/v1/admin/refunds/{id}/reject` | Admin | Reject refund |
| GET/PATCH | `/v1/admin/users` | Admin | User list, points, block/unblock |
| GET/PUT/DELETE | `/v1/admin/settings` | Admin | Setting management |
| POST | `/v1/admin/messages/broadcast` | Admin | Send broadcast message |
| POST | `/v1/admin/messages/targeted` | Admin | Send targeted message |
| POST | `/v1/admin/notifications/test` | Admin | Test notification integrations |
| POST | `/v1/admin/data/import` | Admin | Import existing data backup |
| POST | `/v1/admin/data/repair-aggregates` | Admin | Recalculate stock aggregates |

## Health and Operations

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/healthz` | Public | Liveness probe |
| GET | `/metrics` | Configured | Metrics endpoint when enabled |
