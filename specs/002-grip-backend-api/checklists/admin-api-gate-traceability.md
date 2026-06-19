# Admin API Gate Traceability Matrix

Status: implementation-phase completion evidence for the FE-ready `004..010` pass.

This matrix treats the following as the authoritative source set:

- `../grip-store/specs/004-admin-media-management/spec.md`
- `../grip-store/specs/004-admin-media-management/contracts/media-api.md`
- `../grip-store/specs/004-admin-media-management/tasks.md`
- `../grip-store/specs/005-admin-store-settings/spec.md`
- `../grip-store/specs/005-admin-store-settings/contracts/store-settings-api.md`
- `../grip-store/specs/005-admin-store-settings/test-plan.md`
- `../grip-store/specs/006-admin-reviews/spec.md`
- `../grip-store/specs/006-admin-reviews/contracts/reviews-api.md`
- `../grip-store/specs/006-admin-reviews/test-plan.md`
- `../grip-store/specs/007-admin-content-ops/spec.md`
- `../grip-store/specs/007-admin-content-ops/contracts.md`
- `../grip-store/specs/007-admin-content-ops/test-plan.md`
- `../grip-store/specs/008-admin-catalog-ops/spec.md`
- `../grip-store/specs/008-admin-catalog-ops/contracts.md`
- `../grip-store/specs/008-admin-catalog-ops/test-plan.md`
- `../grip-store/specs/009-admin-order-ops/spec.md`
- `../grip-store/specs/009-admin-order-ops/contracts.md`
- `../grip-store/specs/009-admin-order-ops/test-plan.md`
- `../grip-store/specs/010-admin-user-engagement-ops/use-cases.md`
- `../grip-store/specs/010-admin-user-engagement-ops/spec.md`
- `../grip-store/specs/010-admin-user-engagement-ops/contracts.md`
- `../grip-store/specs/010-admin-user-engagement-ops/test-plan.md`

## Gate Evidence

- Backend regression: `go test ./...` passed on `2026-06-19`.
- Delivery gate: `cd ../grip-store && npm run test:api -- --reporter=line` passed on `2026-06-19` with `115 passed`.

## Module Mapping

### `004-admin-media-management` (legacy, absorbed)

`004` is superseded by `005`, `007`, and `008`. The legacy requirements are absorbed as follows:

| Legacy requirement area | Current owner | Endpoint / contract area | Playwright API coverage |
|---|---|---|---|
| Media library list/presign/register/delete | `007` | `/v1/media`, `/v1/admin/media/presigned` | `media.api.spec.ts`: `should require auth for media listing`, `should reject non-admin media listing`, `should list media assets for admin token`, `should return a presigned upload contract for admin token`; `admin.api.spec.ts`: `should register media metadata successfully with auth`, `should list registered media assets`, `should delete registered media asset` |
| Banner CRUD/order/public reflection | `007` | `/v1/admin/banners`, `/v1/public/homepage` | `content.api.spec.ts`: `should allow admin to manage banners and reflect only active slides publicly in sort order` |
| Storefront footer/floating/homepage settings | `005` | `/v1/admin/store-settings/*`, `/v1/catalog/settings`, `/v1/site-config` | `store-settings.api.spec.ts`: admin read/update and public reflection cases |
| Article/news CRUD and public reflection | `007` | `/v1/content/articles`, `/v1/public/content/articles` | `content.api.spec.ts`: article CRUD and public filtering cases |
| About-us content persistence/reflection | `007` | `/v1/content/pages`, `/v1/public/content/pages/:slug` | `content.api.spec.ts`: `should persist about page content and gallery for public reflection` |
| Product content/media inside catalog editor | `008` | `/v1/admin/products`, related form payloads | `admin.api.spec.ts`: product CRUD baseline; catalog public projection in `catalog.api.spec.ts` |

### `005-admin-store-settings`

| Requirement group | Endpoint / contract area | Playwright API coverage |
|---|---|---|
| Structured admin read model | `/v1/admin/store-settings` | `store-settings.api.spec.ts`: `reads structured store settings payload for admin` |
| Public reflection from same source of truth | `/v1/catalog/settings`, `/v1/site-config` | `store-settings.api.spec.ts`: `reads public catalog settings from the same source of truth` |
| Auth / admin boundary | `/v1/admin/store-settings` | `store-settings.api.spec.ts`: unauthenticated and non-admin rejection cases |
| Brand section validation/persistence | `/v1/admin/store-settings/brand` | `store-settings.api.spec.ts`: `updates brand section with validated payload` |
| Homepage composition validation | `/v1/admin/store-settings/homepage` | `store-settings.api.spec.ts`: `rejects invalid homepage configuration payload` |
| Footer/social persistence | `/v1/admin/store-settings/footer` | `store-settings.api.spec.ts`: `updates footer and social settings with nested structured payload` |
| Floating support validation/persistence | `/v1/admin/store-settings/floating-support` | `store-settings.api.spec.ts`: `updates floating support actions with per-channel validation`, `rejects malformed social link or floating target` |

### `006-admin-reviews`

| Requirement group | Endpoint / contract area | Playwright API coverage |
|---|---|---|
| Moderation queue and stats | `/v1/admin/reviews` | `reviews-moderation.api.spec.ts`: queue auth and stats case |
| Approve / hide / feature transitions | `/v1/admin/reviews/:id/approve`, `/hide`, `/feature` | `reviews-moderation.api.spec.ts`: approve, hide, feature cases |
| Bulk publish | `/v1/admin/reviews/publish-selected` | `reviews-moderation.api.spec.ts`: bulk publish case |
| Delete moderation action | `/v1/admin/reviews/:id` | `reviews-moderation.api.spec.ts`: delete auth case; `reviews.api.spec.ts`: admin delete auth boundary |
| Public reflection only of allowed states | `/v1/catalog/products/:id/reviews` | `reviews.api.spec.ts`: `should return reviews list (public)` |
| Authenticated customer review creation | `/v1/reviews` | `reviews.api.spec.ts`: `should create review with auth`, plus unauthenticated rejection |

### `007-admin-content-ops`

| Requirement group | Endpoint / contract area | Playwright API coverage |
|---|---|---|
| Media library CRUD/presign/register | `/v1/media`, `/v1/admin/media/presigned` | `media.api.spec.ts` and `admin.api.spec.ts` media metadata cases |
| Banner CRUD/order/active public projection | `/v1/admin/banners`, `/v1/public/homepage` | `content.api.spec.ts`: banner management case |
| Article CRUD/publish/public list/detail | `/v1/content/articles`, `/v1/public/content/articles` | `content.api.spec.ts`: article CRUD and filtering cases |
| FAQ CRUD/toggle/reorder/public fetch | `/v1/admin/faqs`, `/v1/faqs/active` | `content.api.spec.ts`: FAQ CRUD/public reflection case |
| About-us content persistence/reflection | `/v1/content/pages`, `/v1/public/content/pages/about` | `content.api.spec.ts`: about page persistence/reflection case |

### `008-admin-catalog-ops`

| Requirement group | Endpoint / contract area | Playwright API coverage |
|---|---|---|
| Admin product list/create/auth boundaries | `/v1/admin/products` | `admin.api.spec.ts`: list/create/401/403 cases |
| Category admin list/create | `/v1/admin/categories` | `admin.api.spec.ts`: category list/create cases |
| Card / inventory-linked operations | `/v1/admin/cards`, `/v1/admin/cards/import`, `/v1/admin/cards/replenish` | `admin.api.spec.ts`: list/import/pull cases |
| Public catalog list/detail/search/categories/settings | `/v1/catalog/products`, `/v1/catalog/products/:id`, `/v1/catalog/search`, `/v1/catalog/categories`, `/v1/catalog/settings` | `catalog.api.spec.ts`: list/detail/buy-meta/search/categories/settings cases |

### `009-admin-order-ops`

| Requirement group | Endpoint / contract area | Playwright API coverage |
|---|---|---|
| Admin order list | `/v1/admin/orders` | `admin.api.spec.ts`: list orders and non-admin rejection |
| Customer order surfaces | `/v1/orders`, `/v1/orders/:id/status`, `/v1/orders/:id/refund-request` | `orders.api.spec.ts`: list/status/refund auth and pagination cases |
| Checkout-owned order creation / status / cancel flow | `/v1/checkout/orders`, `/v1/checkout/orders/:id/status`, `/v1/checkout/orders/:id/cancel`, `/v1/checkout/preview` | `checkout.api.spec.ts`: order create and boundary cases |
| Refund queue + approve / reject | `/v1/admin/refunds`, `/v1/admin/refunds/:id/approve`, `/v1/admin/refunds/:id/reject` | `admin.api.spec.ts`: `should list pending refunds and allow admin rejection`; complementary approve/reject branch coverage remains in `internal/controller/restapi/v1/admin_order_refund_test.go` |

### `010-admin-user-engagement-ops`

| Requirement group | Endpoint / contract area | Playwright API coverage |
|---|---|---|
| User list / moderation entry point | `/v1/admin/users`, `/v1/admin/users/:id` | `admin.api.spec.ts`: `should list users with admin token`, `should update user points with admin token` |
| Broadcast / messaging entry point | `/v1/admin/messages/broadcast`, `/v1/admin/messages/targeted` | `admin.api.spec.ts`: `should broadcast message with admin token`, `should send targeted message with admin token`, plus non-admin rejection |
| Notification admin test-send surface | `/v1/admin/notifications/test` | `admin.api.spec.ts`: `should queue notification test send with admin token` |
| Notification read models for customer engagement | `/v1/notifications`, `/v1/notifications/unread-count`, mark-read/read-all/clear | `notifications.api.spec.ts`: all notification auth and read-model cases |
| Wishlist public/auth surfaces | `/v1/wishlist`, `/v1/wishlist/:id/vote`, `/v1/wishlist/:id` | `wishlist.api.spec.ts`: public list, auth list, mutation boundary cases |
| Profile/check-in side surfaces owned by engagement package | `/v1/profile`, `/v1/user/profile`, `/v1/profile/checkin`, `/v1/user/profile/checkin-status` | `profile.api.spec.ts`: profile/points/check-in auth and read-model cases |

## Shared Gate Files Outside Primary Module Ownership

These files remain part of the required green gate even when they are not the primary owner of `004..010` module intent:

| File | Reason it stays in the gate |
|---|---|
| `auth.api.spec.ts` | foundational auth dependency for every admin/customer contract |
| `catalog.api.spec.ts` | public reflection dependency for catalog/reviews/settings |
| `checkout.api.spec.ts` | order lifecycle dependency for `009` |
| `notifications.api.spec.ts` | customer-visible reflection dependency for `010` |
| `orders.api.spec.ts` | owner-visible order/refund surfaces for `009` |
| `profile.api.spec.ts` | engagement/profile reflection dependency for `010` |
| `wishlist.api.spec.ts` | public/auth engagement reflection dependency for `010` |

## Completion Notes

- No remaining red tests exist in the current Playwright API delivery gate.
- `004` legacy coverage is explicitly remapped rather than treated as a standalone active module.
- Refund moderation evidence in `009` is split between the shared `admin.api.spec.ts` delivery gate and Go controller/usecase tests for complementary transition branches.
