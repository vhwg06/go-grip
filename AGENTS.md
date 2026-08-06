<!-- SPECKIT START -->
For additional context about technologies to be used, project structure,
shell commands, and other important information, read specs/002-grip-backend-api/plan.md
<!-- SPECKIT END -->

# Non-Violable Engineering Rules

- Load and read all applicable `AGENTS.md` files before coding or reviewing code.
- Do not change existing architecture, package boundaries, dependency direction, or domain ownership unless explicitly requested.
- Do not create new packages, subpackages, layers, modules, or abstractions when the change fits the existing structure.
- Do not move code between delivery, use-case, domain, and infrastructure layers without an explicit architectural requirement.
- Do not introduce parallel implementations for behavior already owned by an existing module.
- Do not duplicate domain models, policies, constants, validation rules, ownership contracts, or shared repository contracts across packages.
- Do not place business rules in handlers, repositories, database models, SDK adapters, bootstrap code, or generated code.
- Do not expose protobuf, HTTP framework, ORM, SQL driver, Kubernetes, or external SDK types beyond their adapter boundary.
- Do not expose ORM transaction handles or infrastructure persistence contracts through domain or application ports.
- Do not make domain or application packages import concrete infrastructure implementations.
- Do not construct repositories, infrastructure adapters, or external clients inside domain or use-case code.
- Do not access repositories directly from delivery handlers.
- Do not add generic repositories, base services, utility abstractions, or shared helpers without a proven repeated use case.
- Do not widen an existing interface when a smaller consumer-owned interface is sufficient.
- Do not add speculative fields, methods, extension points, configuration, feature flags, APIs, or abstractions for possible future use.
- Do not implement optional behavior that is not required by the current task.
- Do not change established data semantics, entity grain, identity rules, ownership rules, lifecycle, or source-of-truth behavior implicitly.
- Do not infer physical entities from logical resources, aggregate data, counts, indexes, names, positions, or other unstable identifiers, except where the transformation is explicitly defined by the owning domain contract.
- Do not leave reusable behavior’s contract, failure semantics, or canonical output implicit when they are not fully conveyed by its name and types.
- Do not narrate implementation steps; comments must state guarantees, edge cases, or non-obvious semantic decisions relied upon by callers.
- Do not add comments that merely restate code; comments must document contracts, invariants, edge cases, or non-obvious rationale.
- Do not leave any newly added or modified exported Go declaration without meaningful Go documentation describing its public contract; comments that merely restate the code do not satisfy this rule.
- Do not leave long, complex, reusable, or domain-critical code without comments that preserve its intent, contract, invariants, edge cases, and non-obvious rationale.
- Do not add comments that merely restate syntax or control flow; comments must prevent future readers from having to reconstruct why the code is designed that way.
- Do not introduce or modify domain-level semantic transformations without explicit confirmation when they are not already defined by the owning domain contract.
- Do not treat unstable identifiers as durable identity.
- Do not use timestamps as unique identity unless the domain explicitly defines them as identity.
- Do not use identifiers for ordering unless the domain explicitly defines them as ordered versions.
- Do not silently drop, coerce, repair, normalize, or partially persist invalid domain input unless the existing contract explicitly requires best-effort behavior.
- Do not weaken validation, transactionality, atomicity, idempotency, stale-write protection, concurrency control, or error handling to simplify implementation.
- Do not implement count-before-insert or check-before-write flows when database uniqueness or atomic conflict handling is required.
- Do not introduce non-deterministic behavior based on map iteration, unordered input, database return order, or unstable identifiers.
- Normalize, deduplicate, and sort unordered domain collections before comparison, hashing, or persistence when order is not semantically meaningful.
- Do not hardcode lifecycle values, categories, statuses, resource types, policy values, or error strings repeatedly; reuse the owning module’s constants or value types.
- Do not swallow errors or replace specific errors with successful empty results unless the existing failure policy explicitly permits it.
- Do not add no-op implementations, dummy methods, fake success paths, or placeholder behavior solely to satisfy compilation or interface conformance.
- Do not leave dead wiring, unused runtime loops, inactive command paths, or placeholder implementations that imply unsupported behavior exists.
- Do not change public contracts, protobuf field numbers, database column semantics, compatibility behavior, or externally observable semantics without explicit approval.
- Do not edit generated files manually; update the source definition and run the repository’s generation command.
- Do not modify existing applied migrations.
- Create a new migration following the repository’s current numbering and naming conventions.
- Do not choose or hardcode a migration number before checking the current migration history.
- Do not copy entire files from another branch when only specific behavior is required.
- Do not bring unrelated changes, dependencies, refactors, generated output, formatting changes, or cleanup from source commits.
- Do not perform repository-wide cleanup as part of a feature change.
- Limit cleanup to touched files and avoid formatting-only diff noise.
- Do not combine multiple struct fields, map key-value pairs, or slice elements across single lines within multi-line literal initializations; each field or entry must be on its own line.
- Do not rename unrelated symbols or reorganize unrelated files while implementing a functional change.
- Do not change constructor signatures without updating and verifying every call site, test setup, mock, bootstrap path, and runtime composition path.
- Do not add a dependency to the DI container unless it is required by runtime composition and follows the existing container convention.
- Keep constructors explicit.
- Do not hide required dependencies in globals, service locators, package state, implicit singletons, or lazy initialization.
- Do not start goroutines without explicit lifecycle ownership, context cancellation, failure handling, and shutdown behavior.
- Do not rediscover, rebuild, or mutate input during a retry when the operation requires retrying the same immutable request.
- Do not log credentials, tokens, secrets, private keys, authorization headers, cookies, or complete sensitive payloads.
- Do not change intended behavior merely to make tests pass.
- Do not delete, skip, disable, or weaken existing tests without an explicit reason.
- Add tests for changed behavior, edge cases, invalid input, compatibility, concurrency, and persistence constraints as applicable.
- Only add unit tests for use-case behavior unless another test type is explicitly requested.
- Do not add unit tests for application wiring, delivery handlers, infrastructure, repositories, or adapters unless explicitly requested.
- The final diff must contain only changes required by the task and its proven compile-time or runtime dependencies.
- Do not alter, weaken, or bypass test assertions to match incorrect backend responses; fix the backend Go handler, DTO projection, validation, or error mapping to match openapi.yaml.
- Do not diagnose API test failures or classify contract errors without inspecting exact HTTP request and response JSON payloads for each failing scenario.
- Do not apply superficial patches, dummy returns, or fallback wrappers to mask API contract violations; fix failures at the backend root cause.


## Non-Violable Repository and Aggregate Rules

- Do not create a repository without a clear entity or aggregate owner.
- Repository grain may follow the domain/application CRUD contract, including separately exposed repositories for aggregate children when the current use case requires direct child access.
- A database table alone must not define domain ownership, invariants, lifecycle, or public semantics.
- Aggregate-root orchestration may query and compose child entities through their repositories; child repositories must not invoke or coordinate the aggregate-root repository.
- Aggregate invariants remain in the domain/application layer; child repositories only persist the state they own and do not make cross-entity decisions.
- When a use case changes multiple components of one aggregate and those changes must be atomic, the use case must execute the repositories through a UnitOfWork.
- UnitOfWork owns opening, committing, and rolling back the transaction without exposing ORM transaction handles through application or domain ports.
- Repositories participating in one UnitOfWork operation must use the same UnitOfWork-bound transaction and must not open, commit, roll back, replace, or fall back to an independent transaction.
- Do not permit child persistence operations to continue on a base database connection when a UnitOfWork transaction is required.
- An error from any repository participating in a UnitOfWork operation must propagate and fail the entire transaction.
- Aggregate operation and dependency direction is only:

```text
Aggregate Root -> Child Entities
```

- Reverse operation or dependency direction from child entities or child persistence components to the Aggregate Root is forbidden.
- Child entities and child persistence components must not coordinate, load, save, or invoke the Aggregate Root repository.
- Infrastructure-private child persistence helpers must not be surfaced as domain or application ports.


