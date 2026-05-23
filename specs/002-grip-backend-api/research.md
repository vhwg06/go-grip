# Research: Grip Store Backend API

## Decision: Use REST as the only application transport

**Rationale**: The feature request explicitly requires REST and removal of AMQP RPC, gRPC, and NATS RPC. The existing Fiber REST layer already provides middleware, versioned routes, health checks, metrics, Swagger wiring, and error response patterns. Keeping only REST reduces runtime dependencies and aligns with the frontend's client contract.

**Alternatives considered**:

- Keep all current transports: rejected because the user asked to remove RPC transports and because the feature contract is HTTP/JSON.
- Keep gRPC internally: rejected because it still creates an extra synchronous transport surface and environment/config burden.

## Decision: Preserve DDD and Clean Architecture boundaries

**Rationale**: The existing project already separates `entity`, `usecase`, `repo`, and `controller`. New Grip Store behavior should fit that model: entities model business invariants, usecases orchestrate workflows, repository interfaces describe persistence needs, and REST controllers only bind HTTP input/output and middleware.

**Alternatives considered**:

- Put GORM calls in controllers: rejected because it violates Clean Architecture and makes business rules harder to test.
- Put Fiber context in usecases: rejected because it couples domain logic to HTTP.

## Decision: Use GORM with PostgreSQL for persistent repositories

**Rationale**: The user explicitly requested GORM with PostgreSQL instead of cached in-memory repositories. Current persistent repository files include map-backed behavior in several places; that is insufficient for payments, stock reservation, idempotency, admin operations, and integration tests. GORM provides transaction handling, row locking, model mapping, and mature PostgreSQL support while fitting behind repository interfaces.

**Alternatives considered**:

- Keep pgx and hand-written SQL only: rejected because the requested persistence library is GORM.
- Keep in-memory maps with PostgreSQL placeholders: rejected because it cannot satisfy durability, concurrency, or production behavior.

## Decision: Keep domain entities separate from GORM models

**Rationale**: GORM tags and relationship loading are infrastructure concerns. Domain entities should remain clean and usable in usecase tests without a database. Repository implementations translate between `internal/entity` types and `internal/repo/persistent/models` types.

**Alternatives considered**:

- Add GORM tags directly to entities: rejected because it leaks persistence into the domain layer.
- Use anonymous DTO maps: rejected because it loses compile-time validation and is error-prone for a large domain.

## Decision: Implement stock reservation and order delivery inside database transactions

**Rationale**: Stock and payment correctness are the highest-risk flows. Order creation, point deduction, card reservation, payment confirmation, delivery, cancellation, and refund approval need transaction boundaries with row-level locking and idempotency checks.

**Alternatives considered**:

- Use process-level locks: rejected because they fail across multiple service instances.
- Reserve stock in memory and persist later: rejected because crashes or concurrency would cause overselling.

## Decision: Use Fiber middleware for auth, admin authorization, user status, and rate limits

**Rationale**: The existing REST layer already uses middleware for logging, recovery, and auth. Auth middleware can parse bearer tokens, attach actor context, reject blocked users for mutating routes, and enforce admin access before controllers run.

**Alternatives considered**:

- Auth checks only inside handlers: rejected because it duplicates logic and risks missing protected routes.
- Auth checks only inside usecases: rejected because transport-level request rejection and rate limiting belong at the adapter boundary.

## Decision: Use scheduled jobs inside the service for launch scope

**Rationale**: The spec requires cancellation of expired pending orders, cleanup of expired cards, and aggregate repair. Running these as in-process scheduled jobs keeps deployment simple for the first release while still using repository/usecase interfaces.

**Alternatives considered**:

- External worker service: rejected for initial scope because it adds deployment complexity.
- Manual admin-only cleanup: rejected because stale reservations must be released promptly.

## Decision: Maintain frontend contract compatibility through REST contract tests

**Rationale**: The frontend is already refactored around expected response shapes such as `statusText`, `statusColor`, `maxPurchaseableQuantity`, and checkout preview totals. Contract tests should verify route responses, not just usecase behavior.

**Alternatives considered**:

- Only unit-test usecases: rejected because it does not prove HTTP payload compatibility.
- Defer contract checks to manual frontend testing: rejected because regressions would be late and expensive.
