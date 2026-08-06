---
name: api-test-fix
description: Diagnose and resolve REST API contract failures by identifying the true owner before making changes.
---

# Purpose

Restore alignment between:

- OpenAPI Specification (`docs/api/openapi.yaml`)
- Generated API Code (`api/gen/go/openapi/server.gen.go`)
- Backend Delivery & Logic (`internal/controller/restapi/v1/...`)
- REST API Tests (`grip-store/test`)

Never optimize for green tests.

Optimize for contract correctness.

---

# Workflow

1. Run the failing test module or suite.
2. Cluster failures by root cause / behavior pattern.
3. Audit OpenAPI contract definition.
4. Audit generated OpenAPI types.
5. Audit backend handler implementation.
6. Determine failure owner (Spec, Generated, Backend, Test, Environment).
7. Repair the true owner at the architectural root cause.
8. Re-run the entire failure cluster to confirm all scenarios pass.

See:

- [diagnosis-workflow.md](file:///Users/cynus/Desktop/go-grip/.agents/skills/api-test-fix/references/diagnosis-workflow.md)
- [failure-clusters.md](file:///Users/cynus/Desktop/go-grip/.agents/skills/api-test-fix/references/failure-clusters.md)

---

# Decision Rule

Failure Owner may be:

- **Specification**: Business semantics contradict `openapi.yaml`. (STOP & report issue)
- **Generated API**: Generated code out of sync with spec. (Regenerate pipeline)
- **Backend**: Delivery handler / DTO mapper violates contract. (Fix backend)
- **API Test**: Test setup or step assertion is invalid. (Fix test without weakening)
- **Environment / Fixture**: Database seed or port binding issue. (Repair environment)

Only modify the identified owner.

---

# Before Editing

Read:

- [api-contract.md](file:///Users/cynus/Desktop/go-grip/.agents/skills/api-test-fix/references/api-contract.md)

---

# When Fixing Backend

Read:

- [backend-fix-patterns.md](file:///Users/cynus/Desktop/go-grip/.agents/skills/api-test-fix/references/backend-fix-patterns.md)

---

# Before Finishing

Read:

- [completion-checklist.md](file:///Users/cynus/Desktop/go-grip/.agents/skills/api-test-fix/references/completion-checklist.md)