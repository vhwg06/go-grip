# OpenAPI Contract Principles

`docs/api/openapi.yaml` defines the authoritative HTTP REST API contract for `go-grip`.

Backend handlers must implement it accurately, generated code must mirror it, and API tests (`grip-store/test`) must verify it.

---

## 3-Way Contract Alignment

```text
Business Semantics
        │
        ▼
OpenAPI Contract (docs/api/openapi.yaml)
        │
        ▼
Generated Types (api/gen/go/openapi/server.gen.go)
        │
        ▼
Backend Implementation (internal/controller/restapi/v1/...)
        │
        ▼
REST API Tests (grip-store/test/modules/...)
```

A passing test scenario proves all 5 layers agree.

A failing scenario means **at least one layer is out of sync**. The goal of diagnosis is to find which layer broke alignment.

---

## Non-Negotiable Contract Rules

### 1. OpenAPI is Source of Truth for HTTP Behavior
- HTTP status codes (200, 201, 204, 400, 401, 403, 404, 500) must match `openapi.yaml` response specs.
- Property casing (snake_case vs camelCase) and data types (array, boolean, string, integer, nullable) must strictly match schema components.
- Path parameters and query parameter names (`topic`, `tag`, `page`, `limit`) must be honored.

### 2. Never Weaken Tests to Fit Broken Responses
- **Do NOT** alter expected status codes (e.g. changing 400 to 500 or 200).
- **Do NOT** delete assertions checking mandatory response fields (e.g. `title`, `image_url`, `topic`, `priority`, `content`, `config`, `stats`).
- **Do NOT** accept 500 Internal Server Error when the API spec defines 400 Bad Request or 404 Not Found.

### 3. Handle Specification Errors Correctly
If business requirement / domain contract directly contradicts `openapi.yaml`:
- **STOP execution**.
- Report the specification discrepancy with exact file locations and line numbers.
- Do not patch backend or tests with invalid workarounds merely to force a green test.
