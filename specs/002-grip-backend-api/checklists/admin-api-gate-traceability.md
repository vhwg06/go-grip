# Admin API Gate Traceability Matrix

Status: revised for the no-cards/no-points/no-checkin contract.

## Active Contract Areas

| Requirement group | Endpoint / contract area | Expected verification |
|---|---|---|
| Product list/create/edit Product Editor | `/v1/admin/products`, `/v1/admin/products/{id}`, `/v1/admin/products/{id}/form`, `/v1/admin/products/{id}/status` | Admin product API/UI suites prove the single Product Editor readback plus persisted commercial, media, specs, and intro-article state |
| Category structure | `/v1/admin/categories` | Admin product API/UI suites prove create, hierarchy, and reorder behavior |
| User account read/block | `/v1/admin/users`, `/v1/admin/users/{id}/block` | Admin user API/UI suites prove account summary, block/unblock, and customer handoff semantics |
| Admin orders | `/v1/admin/orders`, `/v1/admin/orders/{id}`, `/v1/admin/orders/{id}/status` | Admin order suites prove list/detail/transition behavior |
| Refund queue and decisions | `/v1/admin/refunds`, `/v1/admin/refunds/{id}`, `/approve`, `/reject` | Admin refund suites prove detail plus approve/reject decisions |
| Structured store settings | `/v1/admin/store-settings*`, `/v1/catalog/settings` | Store-settings suites prove admin persistence and public reflection |
| Profile/account read models | `/v1/profile`, `/v1/user/profile` | Profile API suite proves surviving auth/read/update behavior |

## Removed Contract Areas

These routes or fields are intentionally absent and should be verified as such:

| Removed area | Must stay absent |
|---|---|
| Admin cards | `/v1/admin/cards`, `/v1/admin/cards/import`, `/v1/admin/cards/replenish` |
| Admin user points mutation | `/v1/admin/users/{id}/points` |
| Profile check-in | `/v1/profile/checkin`, `/v1/user/profile/checkin-status` |
| Loyalty/store-setting fields | `points`, `pointsUsed`, `pointsToUse`, `checkinEnabled`, `checkinReward`, `refundReclaimCards` |
