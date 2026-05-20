<!--
Sync Impact Report:
- Version change: [TEMPLATE] -> 1.0.0 (Initial Creation)
- Added sections: Clean Architecture Separation, Transport Agnosticism, Idiomatic Go Quality, Test-Driven and Mocking, Robust Integration Testing
- Removed sections: Template placeholders
- Templates requiring updates:
  ✅ .specify/templates/tasks-template.md (Tests marked mandatory)
- Follow-up TODOs: None
-->
# Go Grip Constitution

## Core Principles

### I. Clean Architecture Separation
Domain logic (Entities and UseCases) **MUST** be independent of frameworks, transport layers (gRPC, REST, AMQP), and databases. All outer layer dependencies (repositories, web APIs) **MUST** be hidden behind interfaces defined in the `entity` or `usecase` layers, strictly adhering to dependency inversion.

### II. Transport Agnosticism
Business logic **MUST** be reusable across all supported transport layers (REST API, gRPC, NATS, RabbitMQ). Controllers are strictly limited to handling transport-specific serialization/deserialization, auth wiring, and delegating execution directly to the underlying UseCases.

### III. Idiomatic Go Quality
All code **MUST** follow standard Go formatting and idioms. Errors **MUST** be handled explicitly without `panic` in business logic, properly categorized within `entity/errors.go`, and propagated with context wrapping where appropriate. Proper project layout (like `cmd`, `internal`, `pkg`) is mandatory.

### IV. Test-Driven and Mocking
Unit tests are **MANDATORY** for all UseCases and Repositories. Generated mocks (e.g. `mocks_repo_test.go`, `mocks_usecase_test.go`) **MUST** be utilized to achieve high test coverage for core business logic. Testing should lead the implementation process where feasible.

### V. Robust Integration Testing
Integration tests **MUST** verify true end-to-end behavior by spinning up actual dependencies (such as Postgres and RabbitMQ) via Docker (e.g., `docker-compose-integration-test.yml`). Core external dependencies **MUST NOT** be heavily mocked during integration tests to guarantee infrastructure reliability.

## Application Constraints

Technology stack is restricted to Go for backend services. Data persistence relies on PostgreSQL, and asynchronous communication/RPC utilizes RabbitMQ and NATS. gRPC is the primary protocol for strict synchronous internal communication.

## Development Workflow

All new features start with defining interfaces within the UseCase/Entity layers before moving onto outer layer implementations. Ensure mock generation scripts are run whenever repository interfaces change. PR reviews must include validation of test coverage and integration test successful runs.

## Governance

This constitution supersedes ad-hoc coding practices. All Pull Requests and code reviews **MUST** verify compliance with these architectural and testing principles. Architectural changes require an explicit amendment to this document.

**Version**: 1.0.0 | **Ratified**: 2026-05-20 | **Last Amended**: 2026-05-20
