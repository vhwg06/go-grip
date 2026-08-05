# Importer Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `Importer` business module.

---

## 1. Owned Symbols
- **Entities**: `ImportItemType` (string enum: `product`, `post`), `ImportItem`, `ImportFailure`, `ImportResult`, `ImportProductDraft`, `ImportPostDraft`
- **Validation Rules**: `MaxInitialImportItems` (const = 25)

---

## 2. Owned Ports
- **Persistence Port**: `ImportRepo`
  - `StoreImportedProduct(ctx context.Context, draft ImportProductDraft) error`
  - `StoreImportedPost(ctx context.Context, draft ImportPostDraft) error`

---

## 3. Application Services
- **Interface**: `ImporterUseCase`
  - `Import(ctx context.Context, items []ImportItem) (ImportResult, error)`
- **Implementation Struct**: `importerUseCase`
- **Constructor**: `NewImporterUseCase(r ImportRepo, maxItems int) ImporterUseCase`

---

## 4. Infrastructure Adapters
- `internal/repo/persistent/import_postgres.go` (`ImportRepo` struct implementing `ImportRepo`)
- `internal/repo/persistent/import_postgres_test.go`

---

## 5. Delivery Endpoints & Mappers
- `internal/controller/restapi/v1/importer/handler.go`
- `internal/controller/restapi/v1/importer/error_mapper.go`
- `internal/controller/restapi/v1/server.go`
- `internal/controller/restapi/router.go`

---

## 6. App Wiring
- `internal/app/app.go` (`useCases.importer *importer.UseCase` → `useCases.importer importer.ImporterUseCase`)

---

## 7. Unit & Integration Tests
- `internal/usecase/importer/importer_test.go` → move to `internal/module/importer/importer_usecase_test.go`

---

## 8. Cross-Module Dependencies & Risks
- Module delegates persistence of product/post data via its own `ImportRepo` port.
- Risk: Low.
