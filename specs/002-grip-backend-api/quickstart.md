# Quickstart: Grip Store Backend API

## Preconditions

- Go toolchain installed
- PostgreSQL available
- Environment configured for the app

## Core Commands

```bash
go test ./... -run '^$'
go test ./internal/... ./integration-test/...
go run go.uber.org/mock/mockgen -source=internal/repo/contracts.go -destination=internal/usecase/mocks_repo_test.go -package=usecase_test
go run go.uber.org/mock/mockgen -source=internal/usecase/contracts.go -destination=internal/usecase/mocks_usecase_test.go -package=usecase_test
go tool swag init --parseDependency --parseInternal -g internal/controller/restapi/router.go
```

## Smoke Expectations

1. `/v1/auth/login` returns access and refresh tokens.
2. `/v1/profile` and `/v1/user/profile` return profile data without points/check-in fields.
3. `/v1/checkout/preview` returns subtotal/final payable values without points fields.
4. `/v1/admin/products` and `/v1/admin/products/{id}/form` cover product edit/media flows.
5. `/v1/admin/users` and `/v1/admin/users/{id}/block` cover current account-management scope.
6. `/v1/admin/store-settings` returns the structured settings contract without check-in or refund-reclaim fields.

## Absence Verification

These routes or fields must not reappear after regeneration or refactors:

- `/v1/admin/cards`
- `/v1/admin/users/{id}/points`
- `/v1/profile/checkin`
- `/v1/user/profile/checkin-status`
- `points`
- `pointsUsed`
- `pointsToUse`
- `checkinEnabled`
- `checkinReward`
- `refundReclaimCards`
