# Diagnosis & Diagnostic Workflow

---

## Failure Clustering Execution

Do **not** analyze failing scenarios one by one blindly. Run the failing suite or module and parse `report.json` to cluster failures by feature area and error pattern:

```bash
# 1. Run failing module against local server
TEST_API_BASE_URL=http://localhost:8080 npx tsx tools/run-cucumber.ts module <MODULE_NAME>

# 2. Cluster failing scenarios from report.json
node -e '
const fs = require("fs");
const data = JSON.parse(fs.readFileSync("/Users/cynus/Desktop/grip-store/test/artifacts/report.json", "utf8"));
const clusters = {};
data.forEach(f => {
  f.elements.forEach(s => {
    const err = s.steps.find(st => st.result?.status === "failed");
    if (err) {
      const msg = (err.result.error_message || "").split("\n")[0];
      clusters[f.name] = clusters[f.name] || [];
      clusters[f.name].push({ scenario: s.name, step: err.name, error: msg });
    }
  });
});
console.log(JSON.stringify(clusters, null, 2));
'
```

---

## 8-Step Sequential Diagnosis

```text
Run Suite & Group Failures
          │
          ▼
   Audit openapi.yaml
          │
          ▼
Audit Generated Code (server.gen.go)
          │
          ▼
Audit Backend Handler (extended_handler.go, etc.)
          │
          ▼
   Determine Owner
    ├── Spec Wrong      ──> STOP & Report
    ├── Backend Wrong   ──> Repair Handler / Mapper
    ├── Generated Wrong ──> Regenerate Pipeline
    ├── Test Wrong      ──> Fix Step (No Weakening)
    └── Env / Fixture   ──> Repair Seed / Auth
          │
          ▼
  Rebuild & Verify Cluster
```

---

## Decision Table: Determining Failure Owner

| Diagnostic Evidence | Owner | Corrective Action |
|---|---|---|
| OpenAPI spec contradicts agreed business requirements | **Specification** | **STOP**. Report specification issue. Do not modify code. |
| Generated `server.gen.go` out of sync with `openapi.yaml` | **Generated Code** | Re-run OpenAPI code generator. |
| Backend returns 500 on invalid input, missing fields, or 404 on valid route | **Backend** | Fix delivery handler, DTO mapper, or error mapper. |
| Test step calls wrong path prefix (e.g. `/v1/public/...` vs `/v1/...`) or misinterprets spec | **API Test** | Fix step definition path without weakening assertions. |
| Test fails due to empty DB or unauthenticated request | **Environment / Fixture** | Provision test fixture (e.g. `createReviewViaApi`) in Given step. |
