# Verification & Completion Checklist

Before declaring a failure cluster resolved, verify every item on this checklist:

---

## Contract & Integrity Checklist

- [ ] **OpenAPI Unchanged Unless Approved**: `docs/api/openapi.yaml` was NOT altered unless a Specification Issue was formally verified.
- [ ] **Tests Not Weakened**: No step assertions, status code checks, or response field assertions in `grip-store/test` were deleted, relaxed, or bypassed.
- [ ] **Root Cause Fixed**: The architectural defect (DTO projection, PATCH merge, error mapping, parameter forwarding) was repaired in delivery handlers rather than using an endpoint-specific hack.
- [ ] **Entire Cluster Green**: All scenarios in the failure cluster turn **GREEN** when re-tested.
- [ ] **Clean Compilation**: `go build ./...` succeeds without compile errors or unused import warnings.
- [ ] **No Duplicated Mapping**: Helper mappers (`toArticleResponse`, `toStaticPageResponse`) are reused across endpoints rather than duplicating struct initializations.
- [ ] **No Contract Drift**: Added or modified handler responses conform strictly to OpenAPI generated types in `server.gen.go`.

---

## Command Verification

```bash
# 1. Compile backend
go build ./...

# 2. Run app directly locally
make dev

# 3. Verify target test module
TEST_API_BASE_URL=http://localhost:8080 npx tsx tools/run-cucumber.ts module <MODULE_NAME>
```
