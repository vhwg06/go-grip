# Implementation Plan: E-Commerce Backend Platform

**Branch**: `001-backend-ecommerce` | **Date**: 2026-05-20 | **Spec**: [specs/001-ecommerce-platform/spec.md](specs/001-ecommerce-platform/spec.md)

**Input**: Feature specification from `/specs/001-ecommerce-platform/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Deliver REST-only e-commerce and content-management capabilities by extending the existing modular monolith and PostgreSQL schema. Reuse current REST controllers, usecase interfaces, repository patterns, and migrations, while keeping the backend fully decoupled from frontend rendering and aligned with system heritage constraints.

## Technical Context

**Language/Version**: Go 1.26

**Primary Dependencies**: Fiber v2 (REST), pgx v5 (PostgreSQL), golang-jwt v5, validator v10, migrate v4, swag, zerolog, testify, mockgen

**Storage**: PostgreSQL (existing database, migrations in `migrations/`)

**Testing**: `go test`, `testify`, mockgen unit tests, integration tests via `docker-compose-integration-test.yml`

**Target Platform**: Linux server (Docker compose)

**Project Type**: Modular monolith REST API service

**Performance Goals**: 95% storefront searches and listings return under 2 seconds; scheduled content published within 5 minutes; 99.5% uptime

**Constraints**: System heritage required (reuse existing modular monolith + database); REST-only (no gRPC/GraphQL); backend fully decoupled from frontend; HTTPS/SSL required; storage budget around 4GB; media uploads limited to 5MB for JPG/PNG/WebP

**Scale/Scope**: Single service + single database + media storage; initial import up to 25 posts/products; public + admin operations in the same service

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Clean Architecture Separation: PASS. Usecase/entity layers remain independent with repository interfaces.
- Transport Agnosticism: PASS. REST-only at controller layer; usecases remain transport-neutral for future adapters.
- Idiomatic Go Quality: PASS. Existing formatting/linting standards applied; errors handled in `entity` and propagated with context.
- Test-Driven and Mocking: PASS. Unit tests required for new usecases/repos using mockgen.
- Robust Integration Testing: PASS. Integration tests run via `docker-compose-integration-test.yml` with Postgres and other dependencies as needed.
- Application Constraints: PASS. Go + PostgreSQL preserved; REST-only constraint honored for this feature.

Post-Design Re-check: PASS. Data model, contracts, and quickstart remain aligned with REST-only and clean architecture constraints.

## Project Structure

### Documentation (this feature)

```text
specs/001-ecommerce-platform/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/
└── app/
    └── main.go

internal/
├── app/
├── controller/
│   ├── restapi/
│   │   ├── router.go
│   │   ├── middleware/
│   │   └── v1/
│   │       ├── auth.go
│   │       ├── user.go
│   │       ├── catalog.go         # new
│   │       ├── cart.go            # new
│   │       ├── content.go         # new
│   │       ├── media.go           # new
│   │       ├── homepage.go        # new
│   │       └── lead.go            # new
│   ├── grpc/                      # existing, not used for this feature
│   ├── amqp_rpc/
│   └── nats_rpc/
├── entity/
│   ├── user.go
│   ├── product.go                 # new
│   ├── category.go                # new
│   ├── cart.go                    # new
│   ├── order_request.go           # new
│   ├── content.go                 # new
│   ├── media.go                   # new
│   └── lead.go                    # new
├── repo/
│   ├── contracts.go
│   └── persistent/
│       ├── user_postgres.go
│       ├── product_postgres.go    # new
│       ├── cart_postgres.go       # new
│       └── content_postgres.go    # new
└── usecase/
    ├── contracts.go
    ├── user/
    ├── catalog/                   # new
    ├── cart/                      # new
    ├── content/                   # new
    ├── media/                     # new
    └── lead/                      # new

migrations/
pkg/
```

**Structure Decision**: Extend the existing modular monolith. REST-only controllers live under `internal/controller/restapi/v1`, with domain logic in `internal/usecase` and entities in `internal/entity`, backed by PostgreSQL repositories in `internal/repo/persistent`.

## Complexity Tracking

No constitution violations requiring justification.
