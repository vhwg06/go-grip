# Quickstart: Grip Store Backend API

## Prerequisites

- Go 1.26 toolchain
- Docker and Docker Compose
- PostgreSQL reachable through `PG_URL`
- No RabbitMQ, NATS, or gRPC service is required for this feature

## Configuration

Required environment groups:

- App: `APP_NAME`, `APP_VERSION`, `LOG_LEVEL`
- HTTP: `HTTP_PORT`, optional `HTTP_USE_PREFORK_MODE`
- PostgreSQL: `PG_URL`, `PG_POOL_MAX`
- Auth: `JWT_SECRET`, access/refresh token durations
- Admin: `ADMIN_USERS`
- Payment: Epay partner credentials and callback URLs
- Feature settings: check-in reward/defaults, wishlist default enablement, scheduler intervals
- Optional integrations: SMTP, Telegram, Bark, external card supplier

Remove required config for `GRPC_PORT`, `RMQ_*`, and `NATS_*` as part of implementation.

## Local Run

1. Start PostgreSQL:

```bash
docker compose up -d postgres
```

2. Run migrations:

```bash
make migrate-up
```

3. Start the REST API:

```bash
go run ./cmd/app
```

4. Verify health:

```bash
curl http://localhost:${HTTP_PORT}/healthz
```

## Test Workflow

1. Generate mocks after repository/usecase interface changes:

```bash
go run go.uber.org/mock/mockgen -source=internal/repo/contracts.go -destination=internal/usecase/mocks_repo_test.go -package=usecase_test
go run go.uber.org/mock/mockgen -source=internal/usecase/contracts.go -destination=internal/usecase/mocks_usecase_test.go -package=usecase_test
```

2. Run unit tests:

```bash
go test ./internal/entity ./internal/usecase/...
```

3. Run repository and REST integration tests with PostgreSQL:

```bash
docker compose -f docker-compose-integration-test.yml up --abort-on-container-exit --build
```

4. Run full verification:

```bash
go test ./...
```

## Validation Notes

- Verified on 2026-05-23:
  - `go test ./...` passes.
  - `docker compose -f docker-compose-integration-test.yml up --abort-on-container-exit --build` passes.
  - No `GRPC_PORT`, `RMQ_*`, or `NATS_*` configuration is required for startup.

## Implementation Order

1. Remove AMQP RPC, gRPC, and NATS RPC startup/config/dependencies from the running app.
2. Add GORM-backed PostgreSQL setup while preserving the existing `pkg/postgres` boundary.
3. Add domain entities and repository/usecase contracts for Grip Store modules.
4. Implement auth middleware, admin middleware, blocked-user checks, and route-level rate limits.
5. Implement catalog/profile/settings read paths first.
6. Implement checkout transaction paths for preview, reservation, zero-price delivery, payment callback, cancellation, and status polling.
7. Implement orders, refunds, wishlist, reviews, notifications, and admin operations.
8. Add scheduled maintenance jobs for expired pending orders, expired cards, and aggregate repair.
9. Add contract/integration tests that verify frontend-compatible response shapes.

## Acceptance Smoke Tests

- Anonymous catalog only returns public visible products.
- Signed-in buyer sees expanded visibility and can preview checkout.
- Concurrent limited-stock purchase attempts never oversell.
- Zero-price points order delivers immediately.
- Replayed payment confirmation delivers once.
- Pending order timeout releases stock and points.
- Card keys are hidden before delivery and visible after delivery to owner/admin.
- Blocked users cannot create orders, wishlist items, votes, or reviews.
- Admin can import cards, approve refund, block user, and repair aggregates.
- Service starts with PostgreSQL and HTTP config only; no broker or gRPC environment variables are required.
