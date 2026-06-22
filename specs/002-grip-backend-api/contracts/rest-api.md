# REST API Contract: Grip Store Backend API

## Common Rules

- Base path: `/v1`
- Protected routes use `Authorization: Bearer <access_token>`.
- Admin routes require authenticated admin context.
- Success responses may use envelope-style payloads, but removed routes/fields listed below must stay absent from the published contract.

## Auth

| Method | Path | Purpose |
|---|---|---|
| POST | `/v1/auth/register` | Register local account |
| POST | `/v1/auth/login` | Authenticate |
| POST | `/v1/auth/refresh` | Rotate access and refresh tokens |
| POST | `/v1/auth/logout` | Revoke current session |
| GET | `/v1/auth/me` | Return current actor/profile summary |

## Catalog

| Method | Path | Purpose |
|---|---|---|
| GET | `/v1/catalog/products` | List visible products |
| GET | `/v1/catalog/products/{id}` | Product detail with optional linked intro article |
| GET | `/v1/catalog/products/{id}/buy-meta` | Product purchase/review metadata |
| GET | `/v1/catalog/search` | Search visible products |
| GET | `/v1/catalog/categories` | List categories |
| GET | `/v1/catalog/settings` | Read public store settings |
| GET | `/v1/catalog/announcement` | Read public announcement |

## Checkout and Orders

| Method | Path | Purpose |
|---|---|---|
| GET | `/v1/checkout/preview` | Calculate subtotal/final payable amount |
| POST | `/v1/checkout/orders` | Create order |
| GET | `/v1/checkout/orders/{id}/status` | Read checkout-owned order status |
| POST | `/v1/checkout/orders/{id}/cancel` | Cancel pending order |
| POST | `/v1/checkout/payment/notify` | Payment callback |
| GET | `/v1/orders` | Owner order list |
| GET | `/v1/orders/{id}` | Owner order detail |
| POST | `/v1/orders/{id}/refund-request` | Create refund request |

Removed from this contract:

- `pointsToUse`
- `pointsUsed`
- zero-price points-only checkout semantics
- card delivery and card-key payload sections

## Profile

| Method | Path | Purpose |
|---|---|---|
| GET | `/v1/profile` | Current profile |
| PUT | `/v1/profile/email` | Update profile email |
| PUT | `/v1/profile/notifications` | Update notification preference |
| GET | `/v1/user/profile` | Legacy-compatible profile read model without points/check-in |

Removed from this contract:

- `POST /v1/profile/checkin`
- `GET /v1/user/profile/checkin-status`
- points/check-in payloads or read models

## Admin Catalog

| Method | Path | Purpose |
|---|---|---|
| GET | `/v1/admin/products` | List products for admin |
| POST | `/v1/admin/products` | Create product |
| PATCH | `/v1/admin/products/{id}` | Update product |
| PATCH | `/v1/admin/products/{id}/status` | Toggle active state |
| GET | `/v1/admin/products/{id}/form` | Read full Product Editor model including linked intro article metadata |
| GET | `/v1/admin/categories` | List categories |
| POST | `/v1/admin/categories` | Create/update/reorder category |

Contract rule:

- Product media/editorial flow remains inside `/v1/admin/products` and `{id}/form`.
- Product row `Edit` / `Quick edit` both target the same Product Editor contract.
- Product Editor scope includes commercial fields, product images, detail/spec content, and linked intro article attach/replace/clear state.
- There is no `/v1/admin/cards` contract.

## Admin Users, Orders, Refunds, Settings

| Method | Path | Purpose |
|---|---|---|
| GET | `/v1/admin/users` | List/search account rows |
| PATCH | `/v1/admin/users/{id}/block` | Block/unblock account |
| GET | `/v1/admin/orders` | Admin order list |
| GET | `/v1/admin/orders/{id}` | Admin order detail |
| PATCH | `/v1/admin/orders/{id}/status` | Update order status |
| DELETE | `/v1/admin/orders/{id}` | Delete order |
| GET | `/v1/admin/refunds` | Refund queue/history |
| GET | `/v1/admin/refunds/{id}` | Refund detail |
| POST | `/v1/admin/refunds/{id}/approve` | Approve refund |
| POST | `/v1/admin/refunds/{id}/reject` | Reject refund |
| GET | `/v1/admin/store-settings` | Structured admin settings read model |
| PUT | `/v1/admin/store-settings/brand` | Update brand section |
| PUT | `/v1/admin/store-settings/contact` | Update contact section |
| PUT | `/v1/admin/store-settings/homepage` | Update homepage section |
| PUT | `/v1/admin/store-settings/footer` | Update footer section |
| PUT | `/v1/admin/store-settings/floating-support` | Update floating support section |
| PUT | `/v1/admin/store-settings/visibility` | Update visibility section |
| PUT | `/v1/admin/store-settings/presence` | Update banner/about presence |
| PUT | `/v1/admin/store-settings/registry` | Update registry section |

Removed from this contract:

- `PATCH /v1/admin/users/{id}/points`
- store-settings fields `checkinEnabled`, `checkinReward`, `refundReclaimCards`

## Explicit Absence Verification

The published contract must not expose:

- `/v1/admin/cards`
- `/v1/admin/cards/import`
- `/v1/admin/cards/replenish`
- `/v1/admin/users/{id}/points`
- `/v1/profile/checkin`
- `/v1/user/profile/checkin-status`
- response fields `points`, `pointsUsed`, `pointsToUse`, `checkinEnabled`, `checkinReward`
