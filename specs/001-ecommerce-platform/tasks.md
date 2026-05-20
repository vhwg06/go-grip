---

description: "Task list for E-Commerce Backend Platform"
---

# Tasks: E-Commerce Backend Platform

**Input**: Design documents from `/specs/001-ecommerce-platform/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are MANDATORY. Each user story must include unit and integration tests.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Go backend code: `internal/`
- REST controllers: `internal/controller/restapi/v1/`
- Repositories: `internal/repo/persistent/`
- Usecases: `internal/usecase/`
- Entities: `internal/entity/`
- Migrations: `migrations/`
- Integration tests: `integration-test/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Update ecommerce configuration defaults in config/config.go
- [X] T002 [P] Add REST response envelope helpers in internal/controller/restapi/response.go
- [X] T003 [P] Add shared validation and pagination DTOs in internal/entity/validation.go and internal/entity/pagination.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 [P] Implement admin auth and RBAC middleware in internal/controller/restapi/middleware/auth.go and internal/controller/restapi/middleware/rbac.go
- [X] T005 [P] Add shared error codes and REST error mapping in internal/entity/errors.go and internal/controller/restapi/error.go
- [X] T006 Update REST router scaffolding for new admin/public route groups in internal/controller/restapi/router.go and internal/controller/restapi/v1/router.go
- [ ] T007 [P] Add integration test helpers for admin auth and seeded data in integration-test/helpers_test.go
- [X] T008 Wire new module dependencies in internal/app/app.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Admin Authentication & RBAC (Priority: P1) 🎯 MVP

**Goal**: Secure admin access with role-based permissions, user management, and lock/unlock workflows

**Independent Test**: Create users, assign roles, and validate allowed and blocked admin actions per role

### Tests for User Story 1 (MANDATORY) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T009 [P] [US1] Add unit tests for admin user and role workflows in internal/usecase/user/user_test.go
- [ ] T010 [P] [US1] Add unit tests for RBAC enforcement in internal/controller/restapi/middleware/rbac_test.go
- [ ] T011 [P] [US1] Add repository unit tests for role and user persistence in internal/repo/persistent/role_postgres_test.go and internal/repo/persistent/user_postgres_test.go
- [X] T012 [US1] Regenerate repository and usecase mocks after user/role contract updates in internal/repo/mocks_repo_test.go and internal/usecase/mocks_usecase_test.go
- [ ] T013 [P] [US1] Add integration tests for admin auth and user management in integration-test/user_admin_test.go

### Implementation for User Story 1

- [X] T014 [P] [US1] Add Role entity and extend User entity with role/status in internal/entity/role.go and internal/entity/user.go
- [X] T015 [P] [US1] Create roles table migration in migrations/20260520000001_create_roles.up.sql and migrations/20260520000001_create_roles.down.sql
- [X] T016 [P] [US1] Add user role/status columns migration in migrations/20260520000002_add_role_status_to_users.up.sql and migrations/20260520000002_add_role_status_to_users.down.sql
- [X] T017 [US1] Extend repo contracts for roles/users in internal/repo/contracts.go
- [X] T018 [P] [US1] Implement role persistence and user role/status updates in internal/repo/persistent/role_postgres.go and internal/repo/persistent/user_postgres.go
- [X] T019 [US1] Extend usecase contracts and admin user operations in internal/usecase/contracts.go and internal/usecase/user/user.go
- [X] T020 [US1] Implement auth and user admin endpoints with RBAC checks in internal/controller/restapi/v1/auth.go and internal/controller/restapi/v1/user.go
- [X] T021 [US1] Register auth and user routes in internal/controller/restapi/v1/router.go

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Product Catalog Management (Priority: P1)

**Goal**: Manage products, categories, tags, attributes, and media assets through admin APIs

**Independent Test**: Create categories, products, tags, and media assets and verify they appear in admin views

### Tests for User Story 2 (MANDATORY) ⚠️

- [X] T022 [P] [US2] Add unit tests for catalog management in internal/usecase/catalog/catalog_test.go
- [X] T023 [P] [US2] Add unit tests for media validation and storage behavior in internal/usecase/media/media_test.go
- [X] T024 [P] [US2] Add repository unit tests for catalog persistence in internal/repo/persistent/product_postgres_test.go, internal/repo/persistent/category_postgres_test.go, and internal/repo/persistent/tag_postgres_test.go
- [X] T025 [P] [US2] Add repository unit tests for media and SEO persistence in internal/repo/persistent/media_postgres_test.go and internal/repo/persistent/seo_postgres_test.go
- [ ] T026 [US2] Regenerate repository and usecase mocks after catalog/media/SEO contract updates in internal/repo/mocks_repo_test.go and internal/usecase/mocks_usecase_test.go
- [ ] T027 [P] [US2] Add integration tests for catalog admin endpoints in integration-test/catalog_admin_test.go
- [ ] T028 [P] [US2] Add integration tests for media upload/list/delete in integration-test/media_admin_test.go

### Implementation for User Story 2

- [X] T029 [P] [US2] Add catalog entities in internal/entity/product.go, internal/entity/category.go, internal/entity/tag.go, internal/entity/product_category.go, and internal/entity/product_tag.go
- [X] T030 [P] [US2] Add media and SEO entities in internal/entity/media.go and internal/entity/seo_metadata.go
- [X] T031 [P] [US2] Create catalog core tables migration in migrations/20260520000100_create_catalog_core.up.sql and migrations/20260520000100_create_catalog_core.down.sql
- [X] T032 [P] [US2] Create media assets migration in migrations/20260520000101_create_media_assets.up.sql and migrations/20260520000101_create_media_assets.down.sql
- [X] T033 [P] [US2] Create SEO metadata migration in migrations/20260520000102_create_seo_metadata.up.sql and migrations/20260520000102_create_seo_metadata.down.sql
- [X] T034 [US2] Extend repo contracts for catalog/media/SEO in internal/repo/contracts.go
- [X] T035 [P] [US2] Implement catalog persistence in internal/repo/persistent/product_postgres.go, internal/repo/persistent/category_postgres.go, and internal/repo/persistent/tag_postgres.go
- [X] T036 [P] [US2] Implement media/SEO persistence in internal/repo/persistent/media_postgres.go and internal/repo/persistent/seo_postgres.go
- [X] T037 [US2] Implement catalog and media usecases in internal/usecase/catalog/catalog.go and internal/usecase/media/media.go
- [X] T038 [US2] Implement catalog/media admin endpoints in internal/controller/restapi/v1/catalog.go and internal/controller/restapi/v1/media.go
- [X] T039 [US2] Register catalog and media routes in internal/controller/restapi/v1/router.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - E-Commerce Storefront Browsing (Priority: P1)

**Goal**: Provide public browsing APIs for search, categories, product detail, homepage composition, and support channels

**Independent Test**: Seed product data and verify search, filter, sort, and homepage blocks return expected results

### Tests for User Story 3 (MANDATORY) ⚠️

- [X] T040 [P] [US3] Add unit tests for public catalog search/sort/filter in internal/usecase/catalog/search_test.go
- [X] T041 [P] [US3] Add unit tests for homepage/support configuration in internal/usecase/content/homepage_test.go
- [X] T042 [P] [US3] Add repository unit tests for homepage and support persistence in internal/repo/persistent/homepage_postgres_test.go and internal/repo/persistent/support_channel_postgres_test.go
- [ ] T043 [US3] Regenerate repository and usecase mocks after homepage/support contract updates in internal/repo/mocks_repo_test.go and internal/usecase/mocks_usecase_test.go
- [ ] T044 [P] [US3] Add integration tests for public storefront endpoints in integration-test/storefront_public_test.go

### Implementation for User Story 3

- [X] T045 [P] [US3] Add homepage and support entities in internal/entity/homepage_block.go and internal/entity/support_channel.go
- [X] T046 [P] [US3] Create homepage/support migrations in migrations/20260520000200_create_homepage_support.up.sql and migrations/20260520000200_create_homepage_support.down.sql
- [X] T047 [US3] Extend repo contracts for homepage/support in internal/repo/contracts.go
- [X] T048 [P] [US3] Implement homepage/support persistence in internal/repo/persistent/homepage_postgres.go and internal/repo/persistent/support_channel_postgres.go
- [X] T049 [US3] Implement public catalog browsing and homepage/support usecases in internal/usecase/catalog/catalog.go and internal/usecase/content/homepage.go
- [X] T050 [US3] Implement public storefront endpoints in internal/controller/restapi/v1/catalog.go and internal/controller/restapi/v1/homepage.go
- [X] T051 [US3] Implement admin homepage/support endpoints in internal/controller/restapi/v1/homepage.go and internal/controller/restapi/v1/support.go
- [X] T052 [US3] Register public/homepage/support routes in internal/controller/restapi/v1/router.go

**Checkpoint**: User Story 3 should now be independently functional

---

## Phase 6: User Story 4 - Shopping Cart & Order Requests (Priority: P2)

**Goal**: Provide anonymous cart management, order request submission, and lead capture

**Independent Test**: Create a cart, update items, submit an order request, verify lead creation, and verify configured notifications are dispatched

### Tests for User Story 4 (MANDATORY) ⚠️

- [X] T053 [P] [US4] Add unit tests for cart operations and order requests in internal/usecase/cart/cart_test.go
- [X] T054 [P] [US4] Add unit tests for lead submissions in internal/usecase/lead/lead_test.go
- [X] T055 [P] [US4] Add unit tests for order notification dispatch in internal/usecase/notification/notification_test.go
- [X] T056 [P] [US4] Add repository unit tests for cart and order persistence in internal/repo/persistent/cart_postgres_test.go and internal/repo/persistent/order_request_postgres_test.go
- [X] T057 [P] [US4] Add repository unit tests for lead persistence in internal/repo/persistent/lead_postgres_test.go
- [ ] T058 [US4] Regenerate repository and usecase mocks after cart/order/lead/notification contract updates in internal/repo/mocks_repo_test.go and internal/usecase/mocks_usecase_test.go
- [ ] T059 [P] [US4] Add integration tests for cart, order request, and notification flow in integration-test/cart_order_test.go
- [ ] T060 [P] [US4] Add integration tests for lead submissions in integration-test/lead_public_test.go

### Implementation for User Story 4

- [X] T061 [P] [US4] Add cart and order entities in internal/entity/cart.go and internal/entity/order_request.go
- [X] T062 [P] [US4] Add lead submission and notification entities in internal/entity/lead.go and internal/entity/notification.go
- [X] T063 [P] [US4] Create cart tables migration in migrations/20260520000300_create_cart_tables.up.sql and migrations/20260520000300_create_cart_tables.down.sql
- [X] T064 [P] [US4] Create order request and lead migrations in migrations/20260520000301_create_order_requests_and_leads.up.sql and migrations/20260520000301_create_order_requests_and_leads.down.sql
- [X] T065 [US4] Extend repo contracts for cart/order/lead in internal/repo/contracts.go
- [X] T066 [P] [US4] Implement cart and order request persistence in internal/repo/persistent/cart_postgres.go and internal/repo/persistent/order_request_postgres.go
- [X] T067 [P] [US4] Implement lead persistence in internal/repo/persistent/lead_postgres.go
- [X] T068 [US4] Add notification dispatch contracts and configuration wiring in internal/usecase/contracts.go and config/config.go
- [X] T069 [P] [US4] Implement notification dispatch usecase in internal/usecase/notification/notification.go
- [X] T070 [US4] Implement cart, order, lead, and notification usecase orchestration in internal/usecase/cart/cart.go and internal/usecase/lead/lead.go
- [X] T071 [US4] Implement cart/order/lead endpoints in internal/controller/restapi/v1/cart.go and internal/controller/restapi/v1/lead.go
- [X] T072 [US4] Register cart/order/lead routes in internal/controller/restapi/v1/router.go

**Checkpoint**: User Story 4 should now be independently functional

---

## Phase 7: User Story 5 - Content & Blog Management (Priority: P3)

**Goal**: Manage articles, static pages, and scheduled publishing with REST-only admin/public endpoints

**Independent Test**: Create draft and scheduled content, verify publish timing, and retrieve public pages

### Tests for User Story 5 (MANDATORY) ⚠️

- [X] T073 [P] [US5] Add unit tests for content article and page workflows in internal/usecase/content/content_test.go
- [X] T074 [P] [US5] Add unit tests for scheduled publishing catch-up in internal/usecase/content/scheduler_test.go
- [X] T075 [P] [US5] Add repository unit tests for content persistence in internal/repo/persistent/content_postgres_test.go
- [ ] T076 [US5] Regenerate repository and usecase mocks after content contract updates in internal/repo/mocks_repo_test.go and internal/usecase/mocks_usecase_test.go
- [ ] T077 [P] [US5] Add integration tests for content management endpoints in integration-test/content_admin_test.go
- [ ] T078 [P] [US5] Add integration tests for public content endpoints in integration-test/content_public_test.go

### Implementation for User Story 5

- [X] T079 [P] [US5] Add content entities in internal/entity/content.go and internal/entity/static_page.go
- [X] T080 [P] [US5] Create content tables migration in migrations/20260520000400_create_content_tables.up.sql and migrations/20260520000400_create_content_tables.down.sql
- [X] T081 [US5] Extend repo contracts for content in internal/repo/contracts.go
- [X] T082 [P] [US5] Implement content persistence in internal/repo/persistent/content_postgres.go
- [X] T083 [US5] Implement content usecases and scheduler in internal/usecase/content/content.go and internal/usecase/content/scheduler.go
- [X] T084 [US5] Implement content admin and public endpoints in internal/controller/restapi/v1/content.go
- [X] T085 [US5] Register content routes in internal/controller/restapi/v1/router.go

**Checkpoint**: All user stories should now be independently functional

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T086 [P] Add catalog search indexes in migrations/20260520000900_add_catalog_indexes.up.sql and migrations/20260520000900_add_catalog_indexes.down.sql
- [X] T087 [P] Add initial content import unit tests for up to 25 posts/products and partial failures in internal/usecase/importer/importer_test.go
- [ ] T088 [P] Add initial content import integration tests for successful and partial imports in integration-test/import_initial_content_test.go
- [X] T089 [P] Add initial import repository unit tests in internal/repo/persistent/import_postgres_test.go
- [X] T090 Add initial import entities and validation results in internal/entity/import_result.go
- [X] T091 Extend import repository and usecase contracts in internal/repo/contracts.go and internal/usecase/contracts.go
- [X] T092 [P] Implement initial import persistence for products and posts in internal/repo/persistent/import_postgres.go
- [X] T093 [P] Implement initial import usecase with partial failure reporting in internal/usecase/importer/importer.go
- [X] T094 Implement initial import admin endpoint in internal/controller/restapi/v1/import.go
- [X] T095 Register initial import routes in internal/controller/restapi/v1/router.go
- [ ] T096 Regenerate repository and usecase mocks after import contract updates in internal/repo/mocks_repo_test.go and internal/usecase/mocks_usecase_test.go
- [ ] T097 [P] Update REST API documentation in docs/swagger.yaml and docs/swagger.json
- [ ] T098 [P] Update feature documentation in docs/backend-specs.md and specs/001-ecommerce-platform/quickstart.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - No dependencies
- **User Story 3 (P1)**: Depends on User Story 2 for product data models, catalog persistence, and media references
- **User Story 4 (P2)**: Depends on User Story 2 for product data snapshots
- **User Story 5 (P3)**: Can start after Foundational (Phase 2) - No dependencies
- **Initial Import (Phase 8)**: Depends on User Story 2 and User Story 5 because it imports both products and posts

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Entities before repositories
- Repositories before usecases
- Usecases before endpoints
- Core implementation before integration

### Parallel Opportunities

- Setup tasks marked [P] can run in parallel after T001
- Foundational tasks marked [P] can run in parallel within Phase 2
- Once Foundational phase completes, User Stories 1, 2, and 5 can start in parallel (if team capacity allows)
- User Story 3 can start after User Story 2 catalog models and persistence are complete
- User Story 4 can start after User Story 2 product data snapshots are available
- Tests for a user story marked [P] can run in parallel
- Entity and migration tasks marked [P] can run in parallel
- Shared contract, mock generation, and router registration tasks are intentionally serialized

---

## Parallel Example: User Story 2

```bash
# Launch all tests for User Story 2 together:
Task: "Add unit tests for catalog management in internal/usecase/catalog/catalog_test.go"
Task: "Add unit tests for media validation and storage behavior in internal/usecase/media/media_test.go"
Task: "Add repository unit tests for catalog persistence in internal/repo/persistent/product_postgres_test.go, internal/repo/persistent/category_postgres_test.go, and internal/repo/persistent/tag_postgres_test.go"
Task: "Add repository unit tests for media and SEO persistence in internal/repo/persistent/media_postgres_test.go and internal/repo/persistent/seo_postgres_test.go"
Task: "Add integration tests for catalog admin endpoints in integration-test/catalog_admin_test.go"
Task: "Add integration tests for media upload/list/delete in integration-test/media_admin_test.go"

# Launch all entity and migration tasks for User Story 2 together:
Task: "Add catalog entities in internal/entity/product.go, internal/entity/category.go, internal/entity/tag.go, internal/entity/product_category.go, and internal/entity/product_tag.go"
Task: "Add media and SEO entities in internal/entity/media.go and internal/entity/seo_metadata.go"
Task: "Create catalog core tables migration in migrations/20260520000100_create_catalog_core.up.sql and migrations/20260520000100_create_catalog_core.down.sql"
Task: "Create media assets migration in migrations/20260520000101_create_media_assets.up.sql and migrations/20260520000101_create_media_assets.down.sql"
Task: "Create SEO metadata migration in migrations/20260520000102_create_seo_metadata.up.sql and migrations/20260520000102_create_seo_metadata.down.sql"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add User Story 4 → Test independently → Deploy/Demo
6. Add User Story 5 → Test independently → Deploy/Demo
7. Add initial import support → Test independently → Deploy/Demo

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 5
3. Start User Story 3 after User Story 2 catalog models and persistence are ready
4. Start User Story 4 after User Story 2 product snapshots are ready
5. Complete initial import tasks after User Stories 2 and 5 are ready

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- REST-only APIs and modular monolith constraints must be honored
- Avoid: vague tasks, shared-file conflicts, cross-story dependencies that break independence
