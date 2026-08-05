# Lead Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `Lead` business module.

---

## 1. Owned Symbols
- **Entities**: `LeadSubmission` (struct), `WorkflowStatus` (string enum: `new`, `in_progress`, `done`)

---

## 2. Owned Ports
- **Persistence Port**: `LeadRepo`
  - `Store(ctx context.Context, lead *LeadSubmission) error`
  - `Get(ctx context.Context, id string) (LeadSubmission, error)`

---

## 3. Application Services
- **Interface**: `LeadUseCase`
  - `Submit(ctx context.Context, lead LeadSubmission) (LeadSubmission, error)`
  - `Get(ctx context.Context, id string) (LeadSubmission, error)`
- **Implementation Struct**: `leadUseCase`
- **Constructor**: `NewLeadUseCase(r LeadRepo) LeadUseCase`

---

## 4. Infrastructure Adapters
- `internal/repo/persistent/lead_postgres.go` (`LeadRepo` struct implementing `LeadRepo`)
- `internal/repo/persistent/lead_postgres_test.go`

---

## 5. Delivery Endpoints & Mappers
- `internal/controller/restapi/v1/lead/handler.go`
- `internal/controller/restapi/v1/lead/mapper.go`
- `internal/controller/restapi/v1/lead/error_mapper.go`
- `internal/controller/restapi/v1/server.go`
- `internal/controller/restapi/router.go`

---

## 6. App Wiring
- `internal/app/app.go` (`useCases.lead *lead.UseCase` → `useCases.lead lead.LeadUseCase`)

---

## 7. Unit & Integration Tests
- `internal/usecase/lead/lead_test.go` → move to `internal/module/lead/lead_usecase_test.go`

---

## 8. Cross-Module Dependencies & Risks
- Independent single-service module.
- Uses `github.com/google/uuid` for ID generation.
