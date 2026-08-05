# Catalog Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `Catalog` business module.

---

## 1. Owned Symbols
- **Entities & Value Objects**: `Product`, `ProductFilter`, `Category`, `Tag`, `ProductCategory`, `ProductTag`, `SEOMetadata`
- **CatalogBase Aggregate Symbols**: CatalogBase domain contracts and types defined in `internal/usecase/catalogbase/contracts.go`.

---

## 2. Owned Ports
- **Persistence Ports**:
  - `CatalogRepo` (in `internal/repo/`)
  - `GripCatalogRepo` (in `internal/repo/`)
  - `CatalogRepositories` & `UnitOfWork` (in `internal/usecase/catalogbase/contracts.go`)

---

## 3. Application Services
- **Interfaces**:
  - `CatalogUseCase`: `CreateProduct`, `ListProducts`, `GetProduct`, `UpdateProduct`, `DeleteProduct`, `CreateCategory`, `ListCategories`, `CreateTag`, `ListTags`
  - `CatalogBaseUseCase`: `catalogbase.UseCase`
- **Implementation Structs**: `catalogUseCase`, `catalogbase.Service`
- **Constructors**: `NewCatalogUseCase`, `NewCatalogBaseUseCase`

---

## 4. Infrastructure Adapters
- `internal/repo/persistent/catalog_postgres.go`
- `internal/repo/persistent/catalogbase_*.go`

---

## 5. Delivery Endpoints & Mappers
- `internal/controller/restapi/v1/catalog/`

---

## 6. App Wiring
- `internal/app/app.go` (`useCases.catalog`, `useCases.catalogBase`)

---

## 7. Unit & Integration Tests
- `internal/usecase/catalog/`
- `internal/usecase/catalogbase/`

---

## 8. Cross-Module Dependencies & Risks
- Referenced by Cart, Orders, Checkout, Wishlist, Admin, Importer.
- Risk: Moderate to High due to cross-module references.
