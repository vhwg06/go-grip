# Implementation Plan: Grip Store Backend API

**Branch**: `002-grip-backend-api` | **Date**: 2026-05-23 | **Spec**: [specs/002-grip-backend-api/spec.md](specs/002-grip-backend-api/spec.md)

**Input**: Feature specification from `/specs/002-grip-backend-api/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Build the Grip Store backend as a Go modular monolith with DDD/Clean Architecture boundaries, REST-only transport, Fiber middleware-based authentication and authorization, GORM-backed PostgreSQL persistence, payment/order/stock correctness, and no AMQP RPC, gRPC, or NATS RPC runtime paths. The current codebase already follows `cmd`, `internal`, and `pkg` structure with entity/usecase/repository/controller separation; this feature replaces the current in-memory repository behavior with durable PostgreSQL repositories and expands the domain model to Grip Store purchasing, inventory, payment, points, notifications, wishlist, reviews, settings, and admin operations.

## Technical Context

**Language/Version**: Go 1.26

**Primary Dependencies**: Fiber v2 for REST, GORM with PostgreSQL driver for persistence, golang-jwt v5 for access tokens, validator v10 for request validation, migrate v4 for schema migration, zerolog for structured logging, swag for API docs, testify and mockgen for tests

**Storage**: PostgreSQL only for durable application state; no cached in-memory repository implementations for production behavior

**Testing**: `go test ./...`, usecase unit tests with generated mocks, repository tests against PostgreSQL, REST integration tests via `docker-compose-integration-test.yml`

**Target Platform**: Linux container runtime using Docker Compose

**Project Type**: REST API backend service implemented as a modular monolith

**Performance Goals**: 95% of catalog/profile/order status reads return in under 1 second; 95% of order writes/admin inventory writes complete in under 2 seconds; 0 oversold stock incidents in concurrent purchase tests

**Constraints**: REST is the only exposed application transport; AMQP RPC, gRPC, and NATS RPC are removed from startup, config, dependencies, and route wiring; auth and admin checks are enforced through HTTP middleware; DDD and Clean Architecture boundaries are mandatory; PostgreSQL is the source of truth; GORM is used instead of in-memory maps for repositories; frontend contract compatibility is required

**Scale/Scope**: One backend service, one PostgreSQL database, 8 business modules from the spec, scheduled maintenance jobs, payment gateway callbacks, and admin/buyer REST surfaces

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Clean Architecture Separation: PASS. Domain entities and usecases remain independent of Fiber, GORM, PostgreSQL, and payment/OAuth clients. Repositories and clients are injected behind interfaces.
- Transport Agnosticism: PASS WITH VARIANCE. Usecases remain transport-neutral, but this feature intentionally removes AMQP RPC, gRPC, and NATS RPC adapters because the user explicitly requires REST-only delivery.
- Idiomatic Go Quality: PASS. Existing `cmd`, `internal`, `pkg` layout remains; errors are categorized in `internal/entity/errors.go` and wrapped at boundaries.
- Test-Driven and Mocking: PASS. Usecase tests must be written against generated repository/client mocks before or alongside implementation.
- Robust Integration Testing: PASS WITH VARIANCE. Integration tests use real PostgreSQL and HTTP server. RabbitMQ/NATS are intentionally excluded because RPC transports are removed.
- Application Constraints: PASS WITH VARIANCE. Go and PostgreSQL remain aligned. RabbitMQ/NATS/gRPC constraints are superseded for this feature by the explicit REST-only requirement.

Post-Design Re-check: PASS WITH SAME VARIANCES. The design keeps Clean Architecture and test requirements intact while documenting removal of RPC transports as a deliberate scope decision.

## Project Structure

### Documentation (this feature)

```text
specs/002-grip-backend-api/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── rest-api.md
└── tasks.md             # Created later by /speckit-tasks
```

### Source Code (repository root)

```text
cmd/
└── app/
    └── main.go

config/
└── config.go            # remove RPC env requirements; add OAuth, payment, admin, cron settings

internal/
├── app/
│   ├── app.go           # initialize REST server only
│   └── migrate.go
├── controller/
│   └── restapi/
│       ├── router.go
│       ├── middleware/
│       │   ├── auth.go
│       │   ├── admin.go
│       │   ├── rate_limit.go
│       │   ├── logger.go
│       │   └── recovery.go
│       └── v1/
│           ├── auth.go
│           ├── catalog.go
│           ├── checkout.go
│           ├── orders.go
│           ├── profile.go
│           ├── wishlist.go
│           ├── notifications.go
│           ├── admin.go
│           └── response.go
├── entity/
│   ├── user.go
│   ├── product.go
│   ├── card.go
│   ├── order.go
│   ├── payment.go
│   ├── refund.go
│   ├── review.go
│   ├── wishlist.go
│   ├── notification.go
│   ├── setting.go
│   └── errors.go
├── repo/
│   ├── contracts.go
│   └── persistent/
│       ├── models/      # GORM persistence models, not domain entities
│       ├── user_postgres.go
│       ├── catalog_postgres.go
│       ├── checkout_postgres.go
│       ├── order_postgres.go
│       ├── profile_postgres.go
│       ├── wishlist_postgres.go
│       ├── notification_postgres.go
│       └── admin_postgres.go
└── usecase/
    ├── contracts.go
    ├── auth/
    ├── catalog/
    ├── checkout/
    ├── orders/
    ├── profile/
    ├── wishlist/
    ├── notification/
    └── admin/

migrations/
pkg/
├── httpserver/
├── jwt/
├── logger/
└── postgres/            # GORM-backed DB setup
```

**Structure Decision**: Extend the current Go Clean Architecture layout. Keep entities and usecases under `internal`, keep REST adapters under `internal/controller/restapi`, implement durable GORM repositories under `internal/repo/persistent`, and remove `internal/controller/amqp_rpc`, `internal/controller/grpc`, `internal/controller/nats_rpc`, `pkg/grpcserver`, `pkg/nats`, and `pkg/rabbitmq` from the running application and dependency graph.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| REST-only instead of all constitution-listed transports | User explicitly requested removal of AMQP RPC, gRPC, and NATS RPC for this feature | Keeping unused transports preserves template complexity, extra env requirements, and unsupported runtime surfaces |
| Excluding RabbitMQ/NATS from integration dependencies | RPC transports are out of scope and should not be required to run or test the service | Requiring unused brokers would make local and CI setup heavier without validating feature behavior |
