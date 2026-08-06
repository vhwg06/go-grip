# Failure Cluster Taxonomy

When API test scenarios fail, do **not** investigate them individually. Group failing scenarios into behavior clusters based on common root causes across handlers and modules.

---

## 1. DTO Projection & Field Mapping
- **Symptom**: Step assertion expects a field (`image_url`, `topic`, `tags`, `priority`, `gallery`, `template_key`), but received `undefined` or `null`.
- **Root Cause**: Handler helper function (e.g. `toArticleResponse()`, `toStaticPageResponse()`) failed to copy domain entity fields into the OpenAPI response struct.
- **Fix**: Update mapper in `extended_handler.go` or `mapper.go` to project all fields.

---

## 2. PATCH Partial Update Merge
- **Symptom**: `PATCH` request causes `500 Internal Server Error` (e.g. PostgreSQL `UNIQUE constraint` violation on `slug`), or overwrites unspecified fields with zero-values (`""`, `0`, `false`).
- **Root Cause**: Handler creates a blank entity struct and only populates fields present in the request body, erasing unsupplied fields on store update.
- **Fix**: Handler must call `GetEntity(id)` first, overlay only non-nil request body fields onto stored entity, then pass to `UpdateEntity()`.

---

## 3. Query Parameter Forwarding (Filtering & Search)
- **Symptom**: Filtering query (e.g. `?topic=engineering`, `?tag=announcement`, `?limit=5`) returns all items unfiltered or wrong page size.
- **Root Cause**: Delivery handler ignores `request.Params` query parameters and passes default or empty filter to usecase.
- **Fix**: Extract `request.Params` (e.g. `Topic`, `Tag`, `Limit`, `Page`) in handler and forward into usecase `ArticleFilter` or `Pagination` struct.

---

## 4. Response Shape & Wrapper Projection
- **Symptom**: Test expects a custom wrapped JSON shape (e.g. `{ config: { brand, contact }, stats: { settingsCount }, visitorCount }` or `{ data: [ blocks ] }`), but received flat or missing keys.
- **Root Cause**: Endpoint returns raw unprojected database maps/slices instead of expected top-level response objects.
- **Fix**: Implement dedicated projection helper (e.g. `settingsProjection()`, `homepageBlocks()`) in delivery handler.

---

## 5. Validation & Error Status Mapping (400 vs 500)
- **Symptom**: Invalid request payload (e.g. invalid email format, invalid target URL, negative quantity) returns `500` instead of `400 Bad Request`.
- **Root Cause**: Handler or `error_mapper.go` does not validate request input or fails to map `ErrInvalidInput` to `http.StatusBadRequest`.
- **Fix**: Validate input in handler (or `validEmail()` check) and update `error_mapper.go` to return status `400`.

---

## 6. Sorting & Ordering
- **Symptom**: Test expects items ordered by `sortOrder` ascending or `priority` descending, but returned list is unsorted.
- **Root Cause**: Handler or usecase omits `sort.SliceStable()` or database `ORDER BY`.
- **Fix**: Add stable sorting in handler or usecase query.

---

## 7. Unwired / Missing Handler Stub
- **Symptom**: Endpoint returns `500 AdminUpdateProductEditorial not implemented` or `404 Not Found` for valid route.
- **Root Cause**: Generated interface method in `server.go` is stubbed with `fmt.Errorf("not implemented")` and not wired to a handler method.
- **Fix**: Implement handler method in `admin_content_handler.go` / `extended_handler.go` and wire call in `server.go`.

---

## 8. Fixtures & Authentication
- **Symptom**: Scenario fails with `401 Unauthorized` or `403 Forbidden` or `missing-review` / empty array.
- **Root Cause**: Test step executed without acquiring JWT token (`getAdminToken()`, `checkoutToken()`) or attempt to query review/checkout on an empty database without provisioning fixture first.
- **Fix**: Call API fixture setup (e.g. `createReviewViaApi()`) in Given step.
