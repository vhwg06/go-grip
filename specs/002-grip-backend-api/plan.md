# Contract Triage Matrix v3 — 109 Failures

**Methodology**: OpenAPI → Implementation → Test  
**Evidence standard**: Each conclusion cites the specific code/schema reference that supports it.

---

## Confidence Legend

| Level | Meaning |
|---|---|
| **Confirmed** | Verified by both OpenAPI schema + source code |
| **High** | Strong code evidence; needs runtime confirmation |
| **Medium** | Multiple hypotheses remain; needs stack trace / log |
| **Unknown** | Needs deploy + log to determine |

---

## Mechanism 1 — Confirmed Test Bugs (sửa test, không sửa backend)

### 1A — Auth: `accessToken` field name mismatch  `[Confirmed]`

| | |
|---|---|
| **Test reads** | `data.token ?? data.access_token ?? ""` (line 186, auth/behavior.steps.ts) |
| **OpenAPI schema** | `TokenPairResponse.accessToken` — camelCase, required |
| **Backend emits** | `{ "accessToken": "...", "refreshToken": "..." }` — correct per schema |
| **Evidence** | `docs/api/openapi.yaml:5565-5574`; `auth/mapper.go:23-27` |

**Fix**: Sửa test line 186: `data.accessToken ?? ""`.  
Also fix refresh token extraction line 273-274: `data.accessToken ?? ""` và `data.refreshToken ?? ""`.

---

### 1B — Notification `POST /notifications/read-all`: 200 vs 204  `[Confirmed]`

| | |
|---|---|
| **Test expects** | status `204` |
| **OpenAPI declares** | `markAllNotificationsRead` → response `200` |
| **Handler returns** | `MarkAllNotificationsRead200Response{}` — correct per schema |
| **Evidence** | `openapi.yaml:2786-2800`; `notification/handler.go:93` |

**Fix**: Sửa test step `expect(status).toBe(204)` → `expect(status).toBe(200)`.

---

### 1C — Payment collection: wrong PUT path  `[Confirmed]`

| | |
|---|---|
| **Test calls** | `PUT /v1/admin/collect` (payment-collection/behavior.steps.ts:45) |
| **OpenAPI declares** | `PUT /admin/collect/setup` |
| **Router registers** | `router.Put(.../admin/collect/setup, ...)` (server.gen.go:5169) |
| **`/admin/collect`** | Only has `GET` → 405 Method Not Allowed |
| **Evidence** | `server.gen.go:5081,5169`; `openapi.yaml:2181,3294` |

**Fix**: Sửa test step đổi path từ `/v1/admin/collect` → `/v1/admin/collect/setup`.

---

### 1D — Banner/FAQ create: 200 vs 201  `[Confirmed]`

| | |
|---|---|
| **Test expects** | status `200` |
| **OpenAPI declares** | `adminSaveBanner` → `201`; `adminSaveFaq` → `201` |
| **Handler returns** | `AdminSaveBanner201JSONResponse`; `AdminSaveFaq201JSONResponse` |
| **Evidence** | `openapi.yaml:1947-1995,1996-2040`; `admin_content_handler.go:138,199` |

**Fix**: Sửa test expect `200` → `201` cho cả banner và FAQ create steps.

---

### 1E — Admin profile: assertions ngoài OpenAPI contract  `[Confirmed]`

| Test asserts | OpenAPI `AccountProfileResponse` declares | Status |
|---|---|---|
| `role \|\| role_id \|\| is_admin` | ❌ Không có | Test bug |
| `password_last_changed_at` | `ProfileSecurityResponse`: `email`, `hasPassword` only | Test bug |
| `two_factor_enabled` | ❌ Không có | Test bug |
| `backup_email` | ❌ Không có | Test bug |
| `device`, `location`, `last_seen_at` | `ProfileSessionResponse`: `id, userAgent, ip, createdAt, expiresAt` | Wrong field names |

**Evidence**: `openapi.yaml:5200-5228,5958-5985`.  
**Fix**: Sửa test step assertions để đọc đúng field names theo OpenAPI schema.

---

### 1F — Admin order queue: `totalAmount` field name  `[Confirmed]`

| | |
|---|---|
| **Test reads** | `totalAmount ?? total_amount ?? total` |
| **OpenAPI schema** | `AdminOrderDetailResponse.amount` (int64) |
| **Handler emits** | `Amount: &amount` — correct per schema |
| **Evidence** | `openapi.yaml:5592-5622`; `admin_orders_handler.go:167` |

**Fix**: Sửa test function `assertOrderRow` đọc field `amount` thay vì `totalAmount`.

---

### 1G — Logout without token: 400 vs 401  `[Confirmed]`

| | |
|---|---|
| **Test expects** | status `401` |
| **Observed** | status `400` |
| **Root cause** | `POST /auth/logout` có `requestBody: required: true`. Test gửi không có body → strict handler `BodyParser` fail → `fiber.NewError(400)` **trước khi** handler check auth. |
| **OpenAPI declares** | `401` response cho unauthenticated |
| **Evidence** | `openapi.yaml:118-144`; `server.gen.go` strict handler BodyParser pattern |

**Verdict**: Spec ambiguity — OpenAPI yêu cầu body nhưng kỳ vọng 401 cho no-auth.  
**Fix option A**: Test gửi kèm empty body `{}` → cho phép handler check auth → return 401.  
**Fix option B**: OpenAPI bỏ `required: true` trên logout body → body optional → handler check auth trước.

---

## Mechanism 2 — Confirmed Implementation Bugs (sửa backend)

### 2A — `GET /catalog/product-models`: result bị discard  `[Confirmed]`

**Code** (`catalog_models_handler.go`):
```go
_, err := h.catalogBase.ListPublicModels(ctx, filter)  // result explicitly discarded
items := make([]map[string]interface{}, 0)              // hardcoded empty
```

| | |
|---|---|
| **OpenAPI** | `ListCatalogProductModels200JSONResponse` → `PublicProductModelListResponse{items, total}` |
| **Test expects** | `items`, `page`, `limit`, `total` non-empty after catalog has data |
| **Evidence** | `catalog_models_handler.go:29`; `openapi.yaml` PublicProductModelListResponse |

**Confirmed**: Handler definitely discards use-case result. Whether this is the **sole** cause of empty list also depends on `ListPublicModels` returning data in test environment — this needs runtime verification after wiring.

**Fix**: Wire result from `ListPublicModels()` into response mapper.

> [!NOTE]
> Filter semantics (categoryId scope, minPrice/maxPrice granularity, search fields, sort order) must be traced OpenAPI → `PublicFilter` → repository. Do not implement filter logic not declared in OpenAPI.

---

### 2B — `GET /admin/orders/{orderId}` missing: 500 instead of 404  `[Confirmed]`

**Root cause chain** (fully traced):

1. `AdminGetOrder` calls `adminUC.GetOrder(ctx, actor, orderID)` → `admin_usecase.go:122`
2. `adminUC.GetOrder` calls `repo.GetOrderByID` without wrapping → `admin_usecase.go:126`
3. `admin_postgres.go` returns `entity.ErrOrderNotFound` when not found → `admin_postgres.go:241`
4. `mapAdminError` checks `entity.ErrNotFound | usermodule.ErrNotFound` — **does NOT check `entity.ErrOrderNotFound`** → falls to default → `http.StatusInternalServerError`
5. Handler returns `AdminGetOrder500JSONResponse{}`

| | |
|---|---|
| **Evidence** | `admin_usecase.go:126`; `admin_postgres.go:241`; `admin/error_mapper.go:42-52`; `entity/errors.go:13,21` |

**Fix**: Add `entity.ErrOrderNotFound` to `mapAdminError` → map to 404.

---

## Mechanism 3 — Catalog Base: 500 on All Write Operations  `[Unknown]`

46 failures across Catalog Base ProductModel and Variant features. HTTP 500 on all create/publish/update operations.

**What is confirmed**:
- `catalogBase` is properly wired: `app.go:112-119`
- Strict handler body parser uses `map[string]any` → cannot fail on valid JSON
- Handler returns `(400JSONResponse, nil)` for all service errors → Fiber writes 400, not 500
- Therefore: HTTP 500 must originate from a path **before** or **outside** the handler's error return

**Hypotheses (unranked, need stack trace)**:

| Hypothesis | Evidence | Probability |
|---|---|---|
| **A**: Migration `20260804000000_catalog_base.up.sql` not applied → `relation does not exist` → `load()` returns `fmt.Errorf(...)` → handler returns 400, but... | Confirmed migration exists locally | Unknown if deployed |
| **B**: `load()` fails → error is NOT `*APIError` → passes through handler differently in some code path | Needs trace | Medium |
| **C**: `ctx.JSON(&response)` serialization fails → propagates error up to Fiber → 500 | Possible if type mismatch | Medium |
| **D**: `nil` repo or `validateRepositories()` returns plain `errors.New()` → handler returns 400, not 500 | Would show 400, not 500 | Low |
| **E**: Strict handler generated by oapi-codegen returns error because handler returns a response object that doesn't implement the expected interface | Can happen after type regen | Medium |

**Required action**: Deploy → curl one endpoint → read response body + server log → identify exact error.

---

## Mechanism 4 — Announcement: Spec Ambiguity  `[Unknown]`

| | |
|---|---|
| **OpenAPI declares** | `getCatalogAnnouncement` → `type: array, items: type: object` (generic schema) |
| **Backend emits** | `[{ "enabled": true, "message": "..." }]` |
| **Test expects** | nullable object `{ id, content, active }` |

**Verdict**: OpenAPI schema is `type: object, additionalProperties: true` effectively. Both `{ enabled, message }` and `{ id, content, active }` are valid against this schema. Neither backend nor test can be definitively called wrong without spec clarification.

**Required action**: Confirm intent — what fields does the announcement object carry?

---

## Mechanism 5 — Cascade Failures (fix upstream first)

### 5A — Order operations: `createPendingOrder` fails  `[High]`

**Chain**:  
1. `createPendingOrder` → `GET /v1/catalog/products` → public catalog handler
2. Public catalog handler discards `ListPublicModels` result → `items: []`
3. `expect(product).toBeTruthy()` fails → scenario aborts

**Affects**: 4 scenarios in Order operations (Execute transition, Reject disallowed, Delete after cancel, Read refund relevance).  
**Fix**: Fix Mechanism 2A (public catalog list wiring) first.

---

### 5B — Checkout: no purchasable product  `[High]`

Same cascade as 5A — `GET /v1/catalog/products` returns empty.

---

### 5C — Catalog master data: category/definition creation failing  `[Medium]`

If catalog write path (Mechanism 3) is broken, test fixtures cannot create:
- `POST /admin/catalog/categories` → prerequisite for model creation
- `POST /admin/catalog/attribute-definitions` → prerequisite for variant dimensions

Affects: 5 Catalog master data scenarios.

---

## Mechanism 6 — Fixture Dependency Issues

### 6A — Refund approve/reject: empty queue fixture  `[High]`

**Root cause**:  
Scenario "Approve/Reject After Evidence Review" calls `GET /admin/refunds?status=pending` and takes `items[0].id`.  
If the pending refund queue is empty (no prior test created a refund in this run), `items[0]` is `undefined` → `refundId = "missing-refund"`.  
`strconv.ParseInt("missing-refund")` fails → handler returns `400`.  
Test expects `200` (assertAccepted).

**This is a fixture isolation issue** — scenario does not self-create the required refund.  
**Fix**: Scenario fixture must create its own order → request refund → then approve/reject.

---

## Mechanism 7 — Remaining Spec Ambiguity

### 7A — Homepage/footer: validation rules not in OpenAPI

| | |
|---|---|
| **Test expects** | `400` for duplicate block priority |
| **OpenAPI declares** | Request schema: `additionalProperties: true`; responses: `200` only, no `400` |

Validation rules not declared in spec. Cannot implement without PRD decision.  
**Required action**: If validation is a business requirement → add `400` response to OpenAPI first.

### 7B — Admin notification test: response body fields

| | |
|---|---|
| **OpenAPI declares** | `type: object` (generic, no required fields) |
| **Test expects** | `{ status: "queued", type: "email" }` |

Generic schema means both `{}` and `{ status: "queued" }` are valid against OpenAPI. Test may be over-specifying.  
**Required action**: Confirm intent — should response include `status` field?

---

## Confidence Summary Table

| Issue | Confidence | Fix location |
|---|---|---|
| Auth `accessToken` field name | **Confirmed** | Test |
| Notification read-all 200 vs 204 | **Confirmed** | Test |
| Payment collect path mismatch | **Confirmed** | Test |
| Banner/FAQ 200 vs 201 | **Confirmed** | Test |
| Admin profile field names | **Confirmed** | Test |
| Order queue `amount` field name | **Confirmed** | Test |
| Logout no-body → 400 not 401 | **Confirmed** | Test or OpenAPI |
| Public catalog list discarded | **Confirmed** | Backend |
| Admin order GET 500 → 404 | **Confirmed** | Backend (`mapAdminError`) |
| Order cascade (create fails) | **High** | Depends on 2A fix |
| Refund approve/reject fixture | **High** | Test fixture |
| Catalog write 500 | **Unknown** | Need stack trace |
| Announcement field names | **Unknown** | Need spec |
| Homepage/footer validation | **Unknown** | Need spec/PRD |
| Notification test response body | **Unknown** | Need spec |

---

## Phase Plan

```
Phase 1 — Fix confirmed test bugs (no deploy needed)
  • 1A auth token field
  • 1B notification read-all status
  • 1C payment collect path
  • 1D banner/FAQ status
  • 1E admin profile field names
  • 1F order queue amount field
  • 1G logout body handling (fix option A: send empty body)
  → Expected recovery: ~14 tests

Phase 2 — Fix confirmed backend bug (2B)
  • Add entity.ErrOrderNotFound to mapAdminError
  → Deploy

Phase 3 — Collect server logs
  • Curl POST /v1/admin/catalog/categories with valid payload + auth token
  • Read response body + application log
  • Identify root cause of Catalog Base 500

Phase 4 — Fix Catalog Base (after root cause confirmed)
  → Deploy

Phase 5 — Fix public catalog list wiring (2A)
  → Deploy
  → Cascade fixes: order createPendingOrder, checkout, refund fixtures

Phase 6 — Confirm spec ambiguities
  • Announcement fields
  • Notification test response body
  • Homepage/footer validation rules
  → Fix or update OpenAPI → implement → deploy

Phase 7 — Fix fixture isolation (6A refund)
  • Refund approve/reject fixture creates its own refund

Phase 8 — Full regression run
```
