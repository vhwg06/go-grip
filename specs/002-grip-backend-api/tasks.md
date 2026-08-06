# Tasks — API Test Failure Fix Plan

Source: [Contract Triage Matrix v3](../../../.gemini/antigravity-ide/brain/55f3c568-798e-4c14-85c3-3cdc15ee93f9/error_group_analysis.md)  
Current state: 127/236 passing. Goal: 236/236 green.

---

## Batch 1 — Test Contract Fixes (no deploy required)

**Constraint**: Changes only in `grip-store/test/`. No backend code touched.  
No deploy needed. Can be verified by re-running the test suite against the live API immediately.  
**Expected recovery**: ~14 scenarios.

---

### T1.1 — Fix auth token field extraction

**Why**: `POST /auth/login` returns `{ accessToken, refreshToken }` per OpenAPI `TokenPairResponse`.  
Test reads `data.token ?? data.access_token` — neither field exists → `accessToken = ""` → all shopper-auth-dependent scenarios fail downstream (8 failures).

**How**:  
- File: `test/modules/auth/behavior.steps.ts`  
- Line 186: `data.token ?? data.access_token ?? ""` → `data.accessToken ?? ""`  
- Line 187: `data.refresh_token ?? data.refreshToken ?? ""` → `data.refreshToken ?? ""`  
- Line 273 (refresh flow): `data.token ?? ""` → `data.accessToken ?? ""`  
- Line 274 (refresh flow): `data.refresh_token ?? ""` → `data.refreshToken ?? ""`

**DoD**: `apiLogin()` sets non-empty `accessToken` and `refreshToken` after a successful login call. Scenarios "Login API returns access and refresh tokens" and all dependent shopper token scenarios pass.

- [ ] T1.1

---

### T1.2 — Fix notification read-all expected status

**Why**: `POST /notifications/read-all` — OpenAPI declares `200`, handler returns `200`. Test asserts `204` — test is wrong.

**How**:  
- File: `test/modules/engagement/behavior.steps.ts`  
- Find step: `"the engagement response status is \`204\`"` and the `markAllNotificationsRead` step  
- Change `expect(status).toBe(204)` → `expect(status).toBe(200)`

**DoD**: "Mark all notifications as read" scenario passes.

- [ ] T1.2

---

### T1.3 — Fix payment collection PUT path

**Why**: Test calls `PUT /v1/admin/collect` — that route only has `GET` → 405 Method Not Allowed.  
OpenAPI and router both define `PUT /admin/collect/setup`.

**How**:  
- File: `test/modules/admin/payment-collection/behavior.steps.ts`  
- Lines 45, 59: `adminPut(this, "/v1/admin/collect", ...)` → `adminPut(this, "/v1/admin/collect/setup", ...)`

**DoD**: "Change Payee Identity" and "Update QR Or Transfer Setup" scenarios return 200 (not 405).

- [ ] T1.3

---

### T1.4 — Fix banner and FAQ create expected status

**Why**: OpenAPI declares `adminSaveBanner → 201` and `adminSaveFaq → 201`. Handler returns 201. Test asserts 200.

**How**:  
- File: `test/modules/admin/content/behavior.steps.ts` (or wherever banner/faq steps live)  
- Find assertions: `expect(status).toBe(200)` for banner and FAQ create steps  
- Change to `expect(status).toBe(201)`

**DoD**: "Manage banners" and "Manage FAQs" scenarios pass the create step with status 201.

- [ ] T1.4

---

### T1.5 — Fix admin profile field name assertions

**Why**: Test asserts fields `role`/`role_id`/`is_admin`, `password_last_changed_at`, `two_factor_enabled`, `backup_email`, `device`, `location`, `last_seen_at` — none of these exist in OpenAPI schemas `AccountProfileResponse`, `ProfileSecurityResponse`, or `ProfileSessionResponse`.

**How**:  
- File: `test/modules/admin/admin-profile/behavior.steps.ts`  
- `GET /profile` step: remove assertion on `role`/`role_id`/`is_admin`; assert only `id`, `username`, `email` which are declared required  
- `GET /profile/security` step: assert `hasPassword` (boolean) and `email` per `ProfileSecurityResponse`; remove `password_last_changed_at`, `two_factor_enabled`, `backup_email`  
- `GET /profile/sessions` step: assert `userAgent`, `ip`, `createdAt`, `expiresAt` per `ProfileSessionResponse`; remove `device`, `location`, `last_seen_at`

**DoD**: "Read Self Identity", "Review Security Posture", "Inspect Recent Access" scenarios pass with correct field assertions.

- [ ] T1.5

---

### T1.6 — Fix order queue field name assertion

**Why**: `AdminOrderDetailResponse` declares field `amount` (int64). Test reads `totalAmount ?? total_amount ?? total` — none match → assertion fails.

**How**:  
- File: `test/modules/admin/order/behavior.steps.ts`  
- Function `assertOrderRow()` (line ~144): change `field(row, "totalAmount", "total_amount", "total")` → `field(row, "amount")`

**DoD**: "Read the order queue" scenario — `assertOrderRow` passes with `amount` field present.

- [ ] T1.6

---

### T1.7 — Fix logout without-auth test to send empty body

**Why**: `POST /auth/logout` requires `requestBody`. Test sends no body → strict handler BodyParser fires before auth check → returns 400 (not 401). Test expects 401.  
Fix: send empty body to let auth check execute.

**How**:  
- File: `test/modules/auth/behavior.steps.ts`  
- Step `"the client logs out without a token"` (line 307):  
  Change `this.getApiClient()).post("/v1/auth/logout")` → `this.getApiClient()).post("/v1/auth/logout", {})`

**DoD**: "Reject logout without authentication" scenario returns 401.

- [ ] T1.7

---

## Batch 2 — Backend: Fix order error mapping (deploy required)

**Constraint**: Only touches `internal/controller/restapi/v1/admin/error_mapper.go`. Run `go test ./...` before signalling deploy. No other files changed in this batch.  
**Expected recovery**: 1 scenario ("Read a missing order").

---

### T2.1 — Add `entity.ErrOrderNotFound` to `mapAdminError`

**Why**: `admin_postgres.GetOrderByID` returns `entity.ErrOrderNotFound` when row not found.  
`mapAdminError` checks `entity.ErrNotFound` and `usermodule.ErrNotFound` — these are separate sentinels (`errors.go:13` vs `errors.go:21`). `entity.ErrOrderNotFound` falls to default → 500. Test expects 404.

**How**:  
- File: `internal/controller/restapi/v1/admin/error_mapper.go`  
- In the `errors.Is(err, entity.ErrNotFound) || errors.Is(err, usermodule.ErrNotFound)` branch: add `|| errors.Is(err, entity.ErrOrderNotFound)`  
- No other changes.

**DoD**:  
- `go test ./internal/...` passes  
- After deploy: `GET /admin/orders/{nonexistent-id}` returns 404, not 500

- [ ] T2.1
- [ ] Signal user to deploy → wait for confirmation

---

## Batch 3 — Diagnose Catalog Base 500 (needs deployed Batch 2 + live log)

**Constraint**: NO code changes in this batch. This batch is pure diagnosis.  
Must have a deployed backend with application logs accessible.  
Result determines Batch 4 tasks. Cannot skip.

---

### T3.1 — Curl POST /admin/catalog/categories and capture log

**Why**: 46 catalog failures all return HTTP 500 but our handler code returns `(400JSONResponse, nil)`. The 500 must originate from a path outside the handler's error return. Root cause is unknown without a stack trace.

**How**:  
```bash
# Get admin token first, then:
curl -s -X POST https://<API_HOST>/v1/admin/catalog/categories \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Diagnosis Test","slug":"diagnosis-test-001"}' \
  -v 2>&1

# Then read application log for this request timestamp
```

**DoD**: Have one of:
- A) Full stack trace / panic log showing exact panic site  
- B) Response body with specific error message (e.g. "relation does not exist")  
- C) Confirmed 400 (not 500) → handler is working, issue is elsewhere

Record exact finding before starting Batch 4.

- [ ] T3.1 — Run diagnosis curl
- [ ] T3.2 — Document root cause from log output

---

## Batch 4 — Fix Catalog Base write path (tasks TBD after T3.2)

**Constraint**: Cannot start until T3.2 is documented.  
Tasks in this batch depend entirely on the root cause found in T3.  
Expected scope (hypotheses, will be confirmed/adjusted after T3):

---

### T4.x — Fix based on T3 diagnosis result

Tasks will be defined after T3.2. Placeholder sub-tasks based on leading hypotheses:

**If Hypothesis A (migration not applied)**:
- T4.A: Confirm migration `20260804000000_catalog_base.up.sql` status via DB  
- T4.A: Run migration on test environment; verify tables exist

**If Hypothesis C (JSON serialization)**:
- T4.C: Identify which response type fails serialization  
- T4.C: Fix the struct/type mismatch causing serialization error

**If Hypothesis E (generated interface mismatch)**:
- T4.E: Regenerate openapi types; verify handler return types implement interfaces  
- T4.E: Fix return type mismatch in affected handlers

**DoD for Batch 4**:  
- `POST /admin/catalog/categories` with valid payload returns `201` (not 500)  
- `POST /admin/catalog/product-models` with valid payload returns `201` (not 500)  
- `go test ./...` passes  
- Deploy → verify in test environment

- [ ] T4.x — Implement fix (tasks TBD after T3.2)
- [ ] Signal user to deploy → wait for confirmation

---

## Batch 5 — Wire public catalog list result (deploy required)

**Constraint**: Batch 4 must be deployed first — `ListPublicModels` only returns data if catalog write path works (products exist).  
Changes in one file only: `internal/controller/restapi/v1/catalog/catalog_models_handler.go`.  
**Expected recovery**: 13 Public catalog query scenarios + 4 order cascade + 2 checkout cascade = ~19 scenarios.

---

### T5.1 — Wire ListPublicModels result into response

**Why**: Handler explicitly discards the service result:  
```go
_, err := h.catalogBase.ListPublicModels(ctx, filter)
items := make([]map[string]interface{}, 0)  // always empty
```  
Test expects `items`, `page`, `limit`, `total` from actual catalog data.

**How**:  
- File: `internal/controller/restapi/v1/catalog/catalog_models_handler.go`  
- Replace `_, err :=` with `result, err :=`  
- Map `result` to `PublicProductModelListResponse` fields: `items`, `total`  
- Pass `filter.Limit` and page from params into `page`/`limit` fields  
- Only wire fields that are declared in OpenAPI `PublicProductModelListResponse` schema — do not add extra fields

**DoD**:  
- `GET /v1/catalog/product-models` with seeded catalog data returns `items` array with actual products and non-zero `total`  
- `go test ./internal/...` passes  
- After deploy: public catalog query scenarios pass

- [ ] T5.1
- [ ] Signal user to deploy → wait for confirmation

---

## Batch 6 — Spec ambiguity resolution (requires PRD decisions)

**Constraint**: Each task in this batch requires explicit clarification from spec/PRD before any code change. Do NOT implement based on test assumptions.  
Tasks are blocked until answers are provided.

---

### T6.1 — Clarify announcement schema fields

**Why**: OpenAPI schema for `GET /catalog/announcement` is `type: array, items: type: object` (no field names). Backend emits `{ enabled, message }`. Test expects `{ id, content, active }`. Both are valid against the generic schema — neither is provably wrong without spec.

**How**: Ask: What fields should the announcement object carry? Options:  
- A) `{ id, content, active }` — sửa backend mapper  
- B) `{ enabled, message }` — sửa test assertions  
- C) Tighten OpenAPI schema to declare exact fields — regen + implement + fix test

**DoD**: PRD decision recorded → implement accordingly → announcement scenario passes.

- [ ] T6.1 — Get spec decision
- [ ] T6.1 — Implement based on decision

---

### T6.2 — Clarify admin notification test response body

**Why**: OpenAPI `queueAdminNotificationTest` → `type: object` (no required fields). Handler returns `{}`. Test asserts `{ status: "queued", type: "email" }`. OpenAPI doesn't require these fields — test may be over-specifying.

**How**: Ask: Should the response include `status` and `type` fields?  
- Yes → add fields to OpenAPI schema → regen → implement in handler  
- No → fix test to only assert 200 status

**DoD**: "Queue an admin notification test send" scenario passes.

- [ ] T6.2 — Get spec decision
- [ ] T6.2 — Implement based on decision

---

### T6.3 — Clarify homepage/footer validation rules

**Why**: OpenAPI `PUT /admin/store-settings/homepage` declares `additionalProperties: true` and no `400` response. Test expects 400 for duplicate block priority. Validation cannot be implemented without it being declared in the spec.

**How**: Ask: Should homepage/footer settings validate ordering uniqueness and non-negative counts?  
- Yes → add `400` response to OpenAPI + declare validation semantics → implement + regen  
- No → remove the 400 assertions from test

**DoD**: Homepage and footer scenarios pass (either validation implemented or test expectations aligned to spec).

- [ ] T6.3 — Get spec decision
- [ ] T6.3 — Implement based on decision

---

## Batch 7 — Fix test fixture isolation (deploy required, after Batch 5)

**Constraint**: Requires Batch 5 to be deployed — fixture must be able to create products and orders.  
Changes only in `grip-store/test/`.

---

### T7.1 — Fix refund approve/reject fixture to self-create refund

**Why**: "Approve After Evidence Review" and "Reject After Evidence Review" read from the global pending refund queue without creating their own refund. If queue is empty, `refundId = "missing-refund"` → `strconv.ParseInt` fails → 400. Test expects 200.

**How**:  
- File: `test/modules/admin/refund/behavior.steps.ts`  
- Step `"evidence supports a positive refund outcome"`: before reading the queue, call `createBrowserRefund(this, ...)` to create an isolated refund  
- Step `"evidence does not support a positive refund outcome"`: same — create a refund first  
- Ensure each scenario creates and uses its own refund (isolated from global state)

**DoD**: "Approve After Evidence Review" and "Reject After Evidence Review" pass without depending on pre-existing queue state.

- [ ] T7.1
- [ ] Signal user to run test suite → read report

---

## Batch 8 — Full regression + close-out

**Constraint**: All previous batches deployed and test-verified.

---

### T8.1 — Run full test suite and triage remaining failures

**Why**: After all fixes, verify overall pass rate. Any remaining failures are either new regressions or spec ambiguities not yet resolved (T6.x).

**How**:  
- Signal user to run `npm run test:api` against latest deployed backend  
- Read `report.json`  
- For each remaining failure: apply same triage methodology (OpenAPI → Implementation → Test)

**DoD**: 236/236 scenarios passing, or all remaining failures have a documented spec decision pending.

- [ ] T8.1 — Request test run from user
- [ ] T8.1 — Triage remaining failures if any

---

## Batch dependency graph

```
Batch 1 (test fixes, no deploy)
    │
    ├──► Batch 2 (mapAdminError fix)
    │         │
    │         └──► [deploy B2]
    │                   │
    │                   └──► Batch 3 (diagnose catalog 500)
    │                             │
    │                             └──► Batch 4 (fix catalog write)
    │                                       │
    │                                       └──► [deploy B4]
    │                                                 │
    │                                                 └──► Batch 5 (wire public list)
    │                                                           │
    │                                                           └──► [deploy B5]
    │                                                                     │
    │                                                                     ├──► Batch 7 (refund fixture)
    │                                                                     └──► Batch 8 (regression)
    │
    └──► Batch 6 (spec decisions — parallel, unblocked but needs PRD)
```
