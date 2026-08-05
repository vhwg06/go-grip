# Cart Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `Cart` business module.

---

## 1. Owned Symbols
- **Entities**: `Cart`, `CartItem`, `CartStatus` (enum: `active`, `abandoned`, `converted`), `OrderRequest`, `WorkflowStatus` (enum: `new`, `in_progress`, `done`)

---

## 2. Owned Ports
- **Persistence Ports**:
  - `CartRepo`: `Store`, `GetBySession`, `AddItem`, `UpdateItem`, `RemoveItem`, `Convert`
  - `OrderRequestRepo`: `Store`, `GetByID`, `UpdateStatus`

---

## 3. Application Services
- **Interfaces**: `CartUseCase`: `Create`, `Get`, `AddItem`, `UpdateItem`, `RemoveItem`, `SubmitOrder`
- **Implementation Structs**: `cartUseCase`
- **Constructors**: `NewCartUseCase`

---

## 4. Infrastructure Adapters
- `internal/repo/persistent/cart_postgres.go`
- `internal/repo/persistent/order_request_postgres.go`

---

## 5. Delivery Endpoints & Mappers
- `internal/controller/restapi/v1/cart/`

---

## 6. App Wiring
- `internal/app/app.go` (`useCases.cart`)

---

## 7. Unit & Integration Tests
- `internal/usecase/cart/`

---

## 8. Cross-Module Dependencies & Risks
- Cart relies on Notification service for sending order confirmation.
- Risk: Low.
