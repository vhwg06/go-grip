# Research: E-Commerce Backend Platform

## Decisions

### REST-only API surface
- **Decision**: Implement all external interfaces as REST APIs using the existing Fiber server.
- **Rationale**: The feature requires REST-only protocols and a frontend-agnostic backend; Fiber is already the HTTP stack in the repo.
- **Alternatives considered**: gRPC for internal calls, GraphQL for flexible queries.

### System heritage and data reuse
- **Decision**: Extend the existing modular monolith and PostgreSQL database using new migrations rather than creating a new service or database.
- **Rationale**: Requirements mandate reuse and extension of the current system and database.
- **Alternatives considered**: New microservice or separate database for the ecommerce domain.

### Anonymous cart persistence
- **Decision**: Store carts by anonymous session identifier in the database without login merge in scope.
- **Rationale**: Requirement explicitly prefers anonymous shopping and order requests without payment complexity.
- **Alternatives considered**: Auth-only carts or guest-to-user cart merge.

### Media upload limits
- **Decision**: Restrict uploads to JPG, PNG, WebP with a 5MB maximum and store metadata alongside media assets.
- **Rationale**: Matches clarified constraints and supports a stable, low-risk storage footprint.
- **Alternatives considered**: Allow larger files or additional media types.

### Scheduled content catch-up
- **Decision**: On service recovery, publish any scheduled items whose time has passed.
- **Rationale**: Ensures content consistency after downtime and meets acceptance criteria.
- **Alternatives considered**: Skip missed schedules or manual reconciliation only.
