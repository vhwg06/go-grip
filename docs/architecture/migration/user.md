# User & Identity Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `User` business module.

---

## 1. Owned Symbols
- **Entities**: `User` (struct), `UserStatus` (enum: `active`, `locked`), `Role` (struct), `RoleName` (enum: `Administrator`, `Editor`, `Author`, `Contributor`, `Subscriber`), `RefreshSession` (struct)

---

## 2. Owned Ports
- **Persistence Ports**:
  - `UserRepo`:
    - `Store(ctx context.Context, user *User) error`
    - `GetByID(ctx context.Context, id string) (User, error)`
    - `GetByEmail(ctx context.Context, email string) (User, error)`
    - `GetByUsername(ctx context.Context, username string) (User, error)`
    - `List(ctx context.Context, limit, offset int) ([]User, int, error)`
    - `Update(ctx context.Context, user *User) error`
  - `RoleRepo`:
    - `GetRole(ctx context.Context, id string) (Role, error)`
  - `AuthRepo`:
    - `StoreRefreshSession(ctx context.Context, session RefreshSession) error`
    - `GetRefreshSession(ctx context.Context, tokenID string) (RefreshSession, error)`
    - `RevokeRefreshSession(ctx context.Context, tokenID string) error`

---

## 3. Application Services
- **Interfaces**:
  - `UserUseCase`: `Register`, `Login`, `GetUser`, `List`, `CreateAdminUser`, `UpdateProfile`, `Lock`, `Unlock`
  - `AuthUseCase`: `Login`, `Refresh`, `Logout`, `Me`
  - `ProfileUseCase`: `Get`, `Update`
  - `AdminUseCase`: `ListUsers`, `UpdateUserStatus`, `ListOrders`, `GetOrder`, `RepairAggregates`
- **Implementation Structs**: `userUseCase`, `authUseCase`, `profileUseCase`, `adminUseCase`
- **Constructors**: `NewUserUseCase`, `NewAuthUseCase`, `NewProfileUseCase`, `NewAdminUseCase`

---

## 4. Infrastructure Adapters
- `internal/repo/persistent/user_postgres.go`
- `internal/repo/persistent/auth_postgres.go`
- `internal/repo/persistent/profile_postgres.go`
- `internal/repo/persistent/admin_postgres.go`

---

## 5. Delivery Endpoints & Mappers
- `internal/controller/restapi/v1/user/`
- `internal/controller/restapi/v1/auth/`
- `internal/controller/restapi/v1/profile/`
- `internal/controller/restapi/v1/admin/`

---

## 6. App Wiring
- `internal/app/app.go` (`useCases.user`, `useCases.auth`, `useCases.profile`, `useCases.admin`)

---

## 7. Unit & Integration Tests
- `internal/usecase/user/`
- `internal/usecase/auth/`
- `internal/usecase/profile/`
- `internal/usecase/admin/`

---

## 8. Cross-Module Dependencies & Risks
- Admin module also references orders/aggregates (cross-module behavior consume via ports or shared snapshots).
- Risk: Moderate due to multi-service scope (`User`, `Auth`, `Profile`, `Admin`).
