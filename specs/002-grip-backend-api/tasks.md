# Tasks: Grip Store Backend API

**Input**: Design documents from `/specs/002-grip-backend-api/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/rest-api.md, quickstart.md

**Tests**: Tests are MANDATORY by project constitution and task template. Each user story includes unit and integration or contract tests that should fail before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it touches different files and has no dependency on incomplete tasks in the same phase.
- **[Story]**: User story label from the feature specification.
- Every task includes an exact target file path.

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare the project for REST-only, GORM-backed implementation.

- [x] T001 Remove AMQP RPC, gRPC, and NATS RPC server initialization from internal/app/app.go
- [x] T002 Remove RPC configuration requirements and add OAuth, admin, payment, and scheduler settings in config/config.go
- [x] T003 Replace pgx-only database setup with GORM PostgreSQL setup while preserving the public constructor in pkg/postgres/postgres.go
- [ ] T004 Add GORM PostgreSQL dependencies and remove unused RPC transport dependencies in go.mod
- [x] T005 Update local and integration compose environment to remove broker/RPC variables and include new Grip Store settings in docker-compose.yml
- [x] T006 Update integration compose environment for REST-only PostgreSQL testing in docker-compose-integration-test.yml
- [x] T007 [P] Update example environment values for REST-only Grip Store configuration in .env.example
- [x] T008 [P] Update Swagger title, base metadata, and REST-only description in internal/controller/restapi/router.go
- [x] T009 [P] Add Grip Store migration placeholders and ordering notes in migrations/20260523000000_grip_store_core.up.sql
- [x] T010 [P] Add rollback placeholders for Grip Store schema in migrations/20260523000000_grip_store_core.down.sql

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish shared contracts, entities, middleware, persistence models, and test fixtures required by all user stories.

**Critical**: No user story implementation should begin until this phase is complete.

### Tests for Foundation

- [x] T011 [P] Add failing repository transaction fixture tests for PostgreSQL setup in internal/repo/persistent/postgres_gorm_test.go
- [x] T012 [P] Add failing REST middleware tests for auth, admin, blocked-user, and rate-limit behavior in internal/controller/restapi/middleware/auth_test.go
- [x] T013 [P] Add failing application startup test proving no RPC servers are started in internal/app/app_test.go

### Implementation for Foundation

- [x] T014 Define shared actor, money, pagination, and status value objects in internal/entity/common.go
- [x] T015 Define Product, Category, Card, and Setting domain entities in internal/entity/product.go
- [x] T016 [P] Define Order, Payment, and RefundRequest domain entities with state transition helpers in internal/entity/order.go
- [x] T017 [P] Define User, RefreshSession, DailyCheckin, and provider identity fields in internal/entity/user.go
- [x] T018 [P] Define Review, WishlistItem, and WishlistVote domain entities in internal/entity/wishlist.go
- [x] T019 [P] Define Notification, BroadcastMessage, BroadcastRead, and AdminMessage domain entities in internal/entity/notification.go
- [x] T020 Add Grip Store business errors and authorization error mapping in internal/entity/errors.go
- [x] T021 Define AuthRepository, CatalogRepository, CheckoutRepository, OrderRepository, ProfileRepository, WishlistRepository, NotificationRepository, and AdminRepository interfaces in internal/repo/contracts.go
- [x] T022 Define Auth, Catalog, Checkout, Orders, Profile, Wishlist, Notification, Admin, and Maintenance usecase interfaces in internal/usecase/contracts.go
- [x] T023 Add GORM persistence models for users, sessions, products, cards, orders, payments, refunds, check-ins, reviews, wishlist, notifications, and settings in internal/repo/persistent/models/grip_store.go
- [x] T024 Implement entity-to-GORM mapping helpers in internal/repo/persistent/models/mappers.go
- [x] T025 Implement shared GORM transaction helper and row-lock helper in internal/repo/persistent/transaction.go
- [x] T026 Implement REST response envelope and error mapper for Grip Store routes in internal/controller/restapi/v1/response.go
- [x] T027 Implement auth actor extraction middleware using bearer tokens in internal/controller/restapi/middleware/auth.go
- [x] T028 Implement admin username authorization middleware in internal/controller/restapi/middleware/admin.go
- [x] T029 Implement blocked-user and mutating-route guard middleware in internal/controller/restapi/middleware/user_status.go
- [x] T030 Implement route group rate limiting for auth and checkout in internal/controller/restapi/middleware/rate_limit.go
- [x] T031 Wire REST-only v1 route groups for auth, catalog, checkout, orders, profile, wishlist, notifications, and admin in internal/controller/restapi/v1/router.go
- [x] T032 Update usecase mock generation targets for new interfaces in internal/repo/contracts.go
- [x] T033 Regenerate repository and usecase mocks after interface changes in internal/usecase/mocks_repo_test.go

**Checkpoint**: Foundation ready. REST app starts without RPC dependencies, shared entities/contracts exist, and user story work can proceed.

---

## Phase 3: User Story 1 - Browse and Buy Digital Goods (Priority: P1) - MVP

**Goal**: Buyers and anonymous visitors can browse visible products, preview checkout, create paid or points-covered orders, complete payment, and receive card keys only after delivery.

**Independent Test**: Seed products and cards, complete one paid order and one points-covered order through REST, and verify stock, order state, and card key visibility.

### Tests for User Story 1

- [x] T034 [P] [US1] Add contract tests for catalog product, search, settings, checkout preview, order creation, payment notify, and order detail responses in integration-test/grip_store_us1_contract_test.go
- [x] T035 [P] [US1] Add usecase tests for catalog visibility, display stock, purchase limits, checkout preview, zero-price delivery, and card-key masking in internal/usecase/checkout/checkout_test.go
- [x] T036 [P] [US1] Add repository integration tests for product listing, card reservation, order creation, payment idempotency, and delivered card reads in internal/repo/persistent/checkout_postgres_test.go

### Implementation for User Story 1

- [x] T037 [P] [US1] Implement GORM catalog repository for visible product list, product detail, search, categories, settings, and announcements in internal/repo/persistent/catalog_postgres.go
- [x] T038 [P] [US1] Implement GORM checkout repository for preview inputs, transactional order creation, card reservation, point deduction, payment records, and zero-price delivery in internal/repo/persistent/checkout_postgres.go
- [x] T039 [P] [US1] Implement GORM order read repository for owner lookup, order status, delivered card masking, and purchase limit history in internal/repo/persistent/order_postgres.go
- [x] T040 [US1] Implement catalog usecase for visibility threshold, display stock, shared product stock, buy metadata, and product search in internal/usecase/catalog/catalog.go
- [x] T041 [US1] Implement checkout usecase for price preview, stock reservation, order creation, zero-price delivery, payment instruction generation, callback handling, and idempotent payment notification in internal/usecase/checkout/checkout.go
- [x] T042 [US1] Implement order read usecase for order detail, order list, status text, status color, and card-key masking in internal/usecase/orders/orders.go
- [x] T043 [US1] Implement catalog REST handlers matching contracts/rest-api.md in internal/controller/restapi/v1/catalog.go
- [x] T044 [US1] Implement checkout and payment REST handlers matching contracts/rest-api.md in internal/controller/restapi/v1/checkout.go
- [x] T045 [US1] Implement order detail and order list REST handlers used by purchase flows in internal/controller/restapi/v1/orders.go
- [x] T046 [US1] Wire catalog, checkout, and order usecases into application construction in internal/app/app.go
- [x] T047 [US1] Add database migrations for products, categories, cards, orders, payments, settings, and critical checkout indexes in migrations/20260523000100_grip_store_catalog_checkout.up.sql
- [x] T048 [US1] Add rollback migration for catalog and checkout tables in migrations/20260523000100_grip_store_catalog_checkout.down.sql
- [x] T049 [US1] Add Swagger annotations for catalog, checkout, and order purchase routes in internal/controller/restapi/v1/catalog.go

**Checkpoint**: User Story 1 is independently functional as the MVP.

---

## Phase 4: User Story 2 - Authenticate and Manage Account Benefits (Priority: P1)

**Goal**: Users can sign in through supported providers, refresh sessions, access profile data, check in daily, spend points, and see account-linked order history.

**Independent Test**: Simulate OAuth callbacks for LinuxDO and GitHub, refresh a session, perform one daily check-in, and verify profile and order history values.

### Tests for User Story 2

- [x] T050 [P] [US2] Add contract tests for OAuth start, OAuth callback, refresh, logout, auth/me, profile, profile update, and check-in routes in integration-test/grip_store_us2_auth_profile_test.go
- [x] T051 [P] [US2] Add usecase tests for account merge, refresh rotation, admin flag resolution, profile dashboard, daily check-in uniqueness, and points updates in internal/usecase/auth/auth_test.go
- [x] T052 [P] [US2] Add repository integration tests for provider identity merge, refresh sessions, check-ins, and profile point updates in internal/repo/persistent/auth_profile_postgres_test.go

### Implementation for User Story 2

- [x] T053 [P] [US2] Implement OAuth client interfaces and LinuxDO/GitHub adapters in internal/repo/webapi/oauth.go
- [x] T054 [P] [US2] Implement GORM auth repository for users, provider identities, refresh sessions, account merge, and admin username lookup in internal/repo/persistent/auth_postgres.go
- [x] T055 [P] [US2] Implement GORM profile repository for dashboard data, profile update, daily check-in, and points balance in internal/repo/persistent/profile_postgres.go
- [x] T056 [US2] Implement auth usecase for OAuth redirect, callback, account merge, token pair issue, refresh rotation, logout, and current user profile in internal/usecase/auth/auth.go
- [x] T057 [US2] Implement profile usecase for dashboard, email update, notification preference update, daily check-in reward, and check-in feature gate in internal/usecase/profile/profile.go
- [x] T058 [US2] Implement auth REST handlers matching contracts/rest-api.md in internal/controller/restapi/v1/auth.go
- [x] T059 [US2] Implement profile REST handlers matching contracts/rest-api.md in internal/controller/restapi/v1/profile.go
- [x] T060 [US2] Wire auth and profile usecases plus OAuth clients into application construction in internal/app/app.go
- [x] T061 [US2] Add database migrations for provider identities, refresh sessions, daily check-ins, and user profile fields in migrations/20260523000200_grip_store_auth_profile.up.sql
- [x] T062 [US2] Add rollback migration for auth and profile tables in migrations/20260523000200_grip_store_auth_profile.down.sql
- [x] T063 [US2] Add Swagger annotations for auth and profile routes in internal/controller/restapi/v1/auth.go

**Checkpoint**: User Story 2 works independently and can support authenticated purchase flows.

---

## Phase 5: User Story 3 - Protect Stock, Payments, and Order Lifecycle (Priority: P1)

**Goal**: The backend prevents overselling, safely processes payment confirmations, cancels stale orders, refunds points where needed, and exposes clear lifecycle state.

**Independent Test**: Race limited-stock purchases, replay payment callbacks, cancel pending orders, run maintenance, and verify no duplicate delivery or lost points.

### Tests for User Story 3

- [x] T064 [P] [US3] Add concurrent checkout and oversell prevention integration tests in integration-test/grip_store_us3_concurrency_test.go
- [x] T065 [P] [US3] Add usecase tests for order state transitions, cancellation, payment replay, delivery idempotency, and refund point handling in internal/usecase/orders/orders_test.go
- [x] T066 [P] [US3] Add repository integration tests for row-level card locking, expired reservation stealing, cleanup jobs, and aggregate repair in internal/repo/persistent/order_lifecycle_postgres_test.go

### Implementation for User Story 3

- [x] T067 [P] [US3] Implement payment verifier and Epay signature adapter behind a usecase interface in internal/repo/webapi/epay.go
- [x] T068 [P] [US3] Extend checkout repository transaction logic for row-level card locks, expired reservation handling, and payment replay protection in internal/repo/persistent/checkout_postgres.go
- [x] T069 [P] [US3] Extend order repository for cancellation, delivery, refund pending state, failed payment state, and status polling in internal/repo/persistent/order_postgres.go
- [x] T070 [P] [US3] Implement maintenance repository methods for expired order cancellation, expired card cleanup, and aggregate sync in internal/repo/persistent/maintenance_postgres.go
- [x] T071 [US3] Extend checkout usecase for Epay verification, duplicate callback handling, and atomic delivery in internal/usecase/checkout/checkout.go
- [x] T072 [US3] Extend orders usecase for cancellation, refund request submission, status polling, and lifecycle validation in internal/usecase/orders/orders.go
- [x] T073 [US3] Implement maintenance usecase for expired pending orders, expired cards, and aggregate repair in internal/usecase/admin/maintenance.go
- [x] T074 [US3] Extend checkout and order REST handlers for cancel, notify, callback, status polling, and refund request routes in internal/controller/restapi/v1/checkout.go
- [x] T075 [US3] Add in-process scheduler startup and graceful shutdown for maintenance jobs in internal/app/app.go
- [x] T076 [US3] Add database migrations for payment idempotency keys, order lifecycle indexes, card reservation indexes, and aggregate fields in migrations/20260523000300_grip_store_lifecycle.up.sql
- [x] T077 [US3] Add rollback migration for lifecycle indexes and payment idempotency additions in migrations/20260523000300_grip_store_lifecycle.down.sql

**Checkpoint**: User Story 3 protects money, stock, and lifecycle correctness under concurrency.

---

## Phase 6: User Story 4 - Admin Operates the Store (Priority: P2)

**Goal**: Admins can manage products, categories, card inventory, orders, refunds, users, settings, messages, notifications, and data repair from REST APIs.

**Independent Test**: Sign in as an admin, perform each management action, and verify non-admin requests are rejected.

### Tests for User Story 4

- [x] T078 [P] [US4] Add contract tests for admin products, categories, cards, orders, refunds, users, settings, messages, and data repair routes in integration-test/grip_store_us4_admin_contract_test.go
- [x] T079 [P] [US4] Add usecase tests for admin authorization, product management, card import, refund approval/rejection, user block, points adjustment, and settings updates in internal/usecase/admin/admin_test.go
- [x] T080 [P] [US4] Add repository integration tests for admin product CRUD, bulk card import, refund approval transaction, user updates, settings CRUD, and aggregate repair in internal/repo/persistent/admin_postgres_test.go

### Implementation for User Story 4

- [x] T081 [P] [US4] Implement GORM admin repository for product, category, card, order, refund, user, setting, message, and aggregate repair operations in internal/repo/persistent/admin_postgres.go
- [x] T082 [P] [US4] Implement admin notification sender interfaces for Telegram, Bark, Email, and no-op test mode in internal/repo/webapi/admin_notifications.go
- [x] T083 [US4] Implement admin usecase for product/category CRUD, card import, card delete, order status changes, refund decisions, user block, point adjustment, settings CRUD, messages, integration tests, and data repair in internal/usecase/admin/admin.go
- [x] T084 [US4] Implement admin REST handlers matching contracts/rest-api.md in internal/controller/restapi/v1/admin.go
- [x] T085 [US4] Wire admin middleware across all `/v1/admin` route groups in internal/controller/restapi/v1/router.go
- [x] T086 [US4] Wire admin usecase and admin notification clients into application construction in internal/app/app.go
- [x] T087 [US4] Add database migrations for refund request admin fields, admin messages, settings constraints, and card import indexes in migrations/20260523000400_grip_store_admin.up.sql
- [x] T088 [US4] Add rollback migration for admin tables and indexes in migrations/20260523000400_grip_store_admin.down.sql
- [x] T089 [US4] Add Swagger annotations for admin routes in internal/controller/restapi/v1/admin.go

**Checkpoint**: User Story 4 gives admins operational control without exposing admin behavior to regular users.

---

## Phase 7: User Story 5 - Engage Users with Wishlist, Reviews, and Notifications (Priority: P3)

**Goal**: Signed-in users can create wishlist ideas, vote, review delivered orders, and use a unified inbox for personal and broadcast messages.

**Independent Test**: Create a delivered order, submit one review, vote on wishlist items, send personal and broadcast messages, and verify unread counts and clear behavior.

### Tests for User Story 5

- [x] T090 [P] [US5] Add contract tests for wishlist, votes, reviews, notifications, unread count, mark read, mark all read, and clear inbox routes in integration-test/grip_store_us5_engagement_contract_test.go
- [x] T091 [P] [US5] Add usecase tests for wishlist feature gate, owner/admin deletion, vote uniqueness, review eligibility, review aggregate recalculation, unified inbox count, and clear behavior in internal/usecase/wishlist/wishlist_test.go
- [x] T092 [P] [US5] Add repository integration tests for wishlist votes, reviews, product aggregate update, notifications, broadcasts, reads, and clear timestamp behavior in internal/repo/persistent/engagement_postgres_test.go

### Implementation for User Story 5

- [x] T093 [P] [US5] Implement GORM wishlist repository for wishlist items, votes, reviews, and review aggregate recalculation in internal/repo/persistent/wishlist_postgres.go
- [x] T094 [P] [US5] Implement GORM notification repository for personal notifications, broadcast messages, broadcast reads, unread counts, and clear state in internal/repo/persistent/notification_postgres.go
- [x] T095 [US5] Implement wishlist usecase for item CRUD, owner/admin authorization, vote toggle, feature gate checks, review creation, and product aggregate updates in internal/usecase/wishlist/wishlist.go
- [x] T096 [US5] Implement notification usecase for unified inbox, unread count, mark read, mark all read, clear inbox, and translation-key payloads in internal/usecase/notification/notification.go
- [x] T097 [US5] Implement wishlist and review REST handlers matching contracts/rest-api.md in internal/controller/restapi/v1/wishlist.go
- [x] T098 [US5] Implement notification REST handlers matching contracts/rest-api.md in internal/controller/restapi/v1/notifications.go
- [x] T099 [US5] Wire wishlist and notification route groups into v1 routing in internal/controller/restapi/v1/router.go
- [x] T100 [US5] Wire wishlist and notification usecases into application construction in internal/app/app.go
- [x] T101 [US5] Add database migrations for wishlist items, wishlist votes, reviews, personal notifications, broadcasts, broadcast reads, and clear timestamps in migrations/20260523000500_grip_store_engagement.up.sql
- [x] T102 [US5] Add rollback migration for engagement tables and indexes in migrations/20260523000500_grip_store_engagement.down.sql
- [x] T103 [US5] Add Swagger annotations for wishlist, review, and notification routes in internal/controller/restapi/v1/wishlist.go

**Checkpoint**: User Story 5 adds engagement features without changing purchase or admin correctness.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Validate the whole feature, remove leftover template transports, and improve operational quality.

- [ ] T104 [P] Update quickstart and REST-only runtime notes in README.md
- [ ] T105 [P] Update generated Swagger artifacts for all Grip Store REST routes in docs/swagger.yaml
- [ ] T106 Remove obsolete RPC controller packages from internal/controller/amqp_rpc, internal/controller/grpc, and internal/controller/nats_rpc
- [ ] T107 Remove obsolete RPC server packages from pkg/grpcserver, pkg/nats, and pkg/rabbitmq
- [ ] T108 Remove obsolete RPC tests, generated protobuf tooling, and unused RPC dependencies from go.mod
- [ ] T109 Run mock generation and commit generated mock updates in internal/usecase/mocks_repo_test.go
- [ ] T110 Run formatting and import ordering across Go files with gofmt/goimports-equivalent tooling in internal
- [ ] T111 Run full unit test suite and fix failures reported by go test ./... in go.mod
- [ ] T112 Run integration test suite and fix failures reported by docker-compose-integration-test.yml
- [ ] T113 Validate quickstart smoke tests and document any required environment changes in specs/002-grip-backend-api/quickstart.md
- [ ] T114 [P] Add final security checklist notes for card-key masking, admin-only routes, blocked-user guards, and payment signature verification in specs/002-grip-backend-api/checklists/security.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 Setup**: No dependencies.
- **Phase 2 Foundational**: Depends on Phase 1 and blocks all user stories.
- **Phase 3 US1**: Depends on Phase 2 and is the MVP purchase flow.
- **Phase 4 US2**: Depends on Phase 2; can proceed in parallel with US1 after shared auth contracts are stable, but purchase-point integration is validated with US1.
- **Phase 5 US3**: Depends on Phase 2 and uses checkout/order components from US1; prioritize immediately after US1 for production safety.
- **Phase 6 US4**: Depends on Phase 2 and uses entities from US1-US3; can begin after admin middleware and repository contracts exist.
- **Phase 7 US5**: Depends on Phase 2 and delivered-order behavior from US1/US3 for review eligibility.
- **Phase 8 Polish**: Depends on selected user stories being complete.

### User Story Dependencies

- **US1 Browse and Buy Digital Goods**: Foundation only; recommended MVP.
- **US2 Authenticate and Manage Account Benefits**: Foundation only for auth/profile; points-spend validation integrates with US1.
- **US3 Protect Stock, Payments, and Order Lifecycle**: Foundation plus US1 checkout/order primitives.
- **US4 Admin Operates the Store**: Foundation plus product/order/user entities; can be developed alongside US3 with interface coordination.
- **US5 Engage Users with Wishlist, Reviews, and Notifications**: Foundation plus delivered-order behavior for reviews.

### Within Each User Story

- Tests first and expected to fail.
- Repository contracts and models before repository implementations.
- Usecases before REST handlers.
- REST handlers before route wiring.
- Migrations before repository integration tests can pass.
- Story complete only when unit, contract, and integration tests for that story pass.

---

## Parallel Execution Examples

### User Story 1

```text
Parallel tests:
- T034 integration-test/grip_store_us1_contract_test.go
- T035 internal/usecase/checkout/checkout_test.go
- T036 internal/repo/persistent/checkout_postgres_test.go

Parallel repository work after entities/contracts:
- T037 internal/repo/persistent/catalog_postgres.go
- T038 internal/repo/persistent/checkout_postgres.go
- T039 internal/repo/persistent/order_postgres.go
```

### User Story 2

```text
Parallel tests:
- T050 integration-test/grip_store_us2_auth_profile_test.go
- T051 internal/usecase/auth/auth_test.go
- T052 internal/repo/persistent/auth_profile_postgres_test.go

Parallel adapter/repository work:
- T053 internal/repo/webapi/oauth.go
- T054 internal/repo/persistent/auth_postgres.go
- T055 internal/repo/persistent/profile_postgres.go
```

### User Story 3

```text
Parallel tests:
- T064 integration-test/grip_store_us3_concurrency_test.go
- T065 internal/usecase/orders/orders_test.go
- T066 internal/repo/persistent/order_lifecycle_postgres_test.go

Parallel infrastructure work:
- T067 internal/repo/webapi/epay.go
- T068 internal/repo/persistent/checkout_postgres.go
- T069 internal/repo/persistent/order_postgres.go
- T070 internal/repo/persistent/maintenance_postgres.go
```

### User Story 4

```text
Parallel tests:
- T078 integration-test/grip_store_us4_admin_contract_test.go
- T079 internal/usecase/admin/admin_test.go
- T080 internal/repo/persistent/admin_postgres_test.go

Parallel adapters:
- T081 internal/repo/persistent/admin_postgres.go
- T082 internal/repo/webapi/admin_notifications.go
```

### User Story 5

```text
Parallel tests:
- T090 integration-test/grip_store_us5_engagement_contract_test.go
- T091 internal/usecase/wishlist/wishlist_test.go
- T092 internal/repo/persistent/engagement_postgres_test.go

Parallel repositories:
- T093 internal/repo/persistent/wishlist_postgres.go
- T094 internal/repo/persistent/notification_postgres.go
```

---

## Implementation Strategy

### MVP First

1. Complete Phase 1 and Phase 2.
2. Complete Phase 3 User Story 1.
3. Stop and validate anonymous catalog, checkout preview, order creation, payment notification, zero-price delivery, and card-key visibility.

### Production Safety Increment

1. Complete User Story 2 for real users, sessions, profile, points, and check-ins.
2. Complete User Story 3 before production launch so stock/payment/order lifecycle behavior is safe under concurrency and callback replay.

### Operational Increment

1. Complete User Story 4 for admin product, card, order, refund, user, setting, and data repair operations.
2. Complete User Story 5 for engagement and notifications.
3. Complete Phase 8 to remove leftover RPC code, run full verification, and update docs.

## Task Summary

- Total tasks: 114
- Setup tasks: 10
- Foundational tasks: 23
- US1 tasks: 16
- US2 tasks: 14
- US3 tasks: 14
- US4 tasks: 12
- US5 tasks: 14
- Polish tasks: 11
- Parallel opportunities: setup docs/migrations, foundation tests/entities, story tests, story repositories/adapters, docs/security checklist
