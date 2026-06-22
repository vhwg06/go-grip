# Feature Specification: Grip Store Backend API

**Feature Branch**: `002-grip-backend-api`

**Created**: 2026-05-23

**Status**: Revised after `remove cards / points / checkin`

## User Scenarios & Testing

### User Story 1 - Browse and Buy Catalog Products (Priority: P1)

Buyers browse visible products, preview order totals, create orders, complete payment when required, and read order status without any cards/points loyalty side flow.

**Independent Test**: Seed a visible product, create a normal paid order, complete payment, and verify order status/read models update through the public and owner routes.

**Acceptance Scenarios**:

1. **Given** an anonymous visitor, **When** they view catalog and product detail, **Then** only visible products are returned with public pricing, stock, and review summary data.
2. **Given** a signed-in buyer, **When** they preview checkout, **Then** the system returns subtotal and final payable amount with no points-related fields.
3. **Given** a payable order, **When** payment is confirmed, **Then** the order advances through the expected lifecycle exactly once.
4. **Given** a delivered or pending order, **When** order details are viewed, **Then** the client receives order state, payment metadata, and buyer-safe detail only.

### User Story 2 - Authenticate and Manage Account Profile (Priority: P1)

Users sign in with supported credentials, refresh sessions, view profile data, and update basic profile settings. Check-in and points are not part of the contract.

**Independent Test**: Sign in, refresh, read `/v1/profile`, update profile fields, and confirm removed loyalty/check-in routes stay absent.

**Acceptance Scenarios**:

1. **Given** a valid account, **When** the user signs in, **Then** the client receives authenticated access and refresh tokens.
2. **Given** a signed-in user, **When** they request profile data, **Then** the response includes identity, notification preference, and admin eligibility without points/check-in fields.
3. **Given** an expired access token with a valid refresh token, **When** refresh is requested, **Then** a new access/refresh pair is issued and the old refresh token cannot be reused.
4. **Given** a removed loyalty endpoint, **When** it is requested, **Then** the backend does not expose that route.

### User Story 3 - Protect Orders, Payments, and Lifecycle (Priority: P1)

The store prevents invalid order transitions, processes payment notifications safely, cancels stale orders, and handles refunds without card reclaim or points restoration behavior.

**Independent Test**: Create orders, replay payment notifications, cancel stale orders, and process refunds while verifying lifecycle correctness and absence of removed side effects.

**Acceptance Scenarios**:

1. **Given** a pending order, **When** payment succeeds, **Then** the order moves through the allowed status path only once.
2. **Given** a stale pending order, **When** maintenance runs or the buyer cancels, **Then** the order is cancelled cleanly.
3. **Given** a refund request, **When** an admin approves or rejects it, **Then** refund state is persisted with decision metadata and no card/points behavior occurs.
4. **Given** duplicate payment confirmations, **When** they are replayed, **Then** no duplicate order effects are created.

### User Story 4 - Admin Operates Catalog, Users, Orders, and Settings (Priority: P2)

Admins manage products, categories, orders, refunds, users, settings, reviews, and messages from privileged REST and UI flows. Product editorial/media stays on `/admin/products`; there is no `/admin/cards` surface.

**Independent Test**: Sign in as admin, exercise product edit/media save, block/unblock users, manage settings, process refunds, and confirm removed admin cards/user-points routes stay absent.

**Acceptance Scenarios**:

1. **Given** an admin, **When** they create or edit a product, **Then** commercial and editorial changes persist through `/v1/admin/products` and its form/readback routes.
2. **Given** an admin, **When** they manage user state, **Then** block/unblock and account read-model behavior work without points mutation fields or endpoints.
3. **Given** an admin, **When** they manage store settings, **Then** the structured settings contract persists without check-in or refund-reclaim flags.
4. **Given** a removed admin route such as `/v1/admin/cards`, **When** it is requested, **Then** the backend does not expose that route.

## Requirements

### Functional Requirements

- **FR-001**: System MUST authenticate users through local account credentials and JWT-backed sessions.
- **FR-002**: System MUST issue renewable authenticated sessions, rotate refresh credentials after each use, and support logout.
- **FR-003**: System MUST identify admin users by configured usernames and surface admin eligibility in the current-user profile.
- **FR-004**: System MUST expose buyer profile information needed by the current client without points or check-in fields.
- **FR-005**: System MUST list and search visible products for anonymous and signed-in users.
- **FR-006**: System MUST expose product detail data including price, visibility, stock, and review summary.
- **FR-007**: System MUST provide public settings and announcement data required by the client.
- **FR-008**: System MUST calculate checkout previews without points deduction fields.
- **FR-009**: System MUST create orders only when product availability, visibility, purchase rules, and buyer permissions pass.
- **FR-010**: System MUST attach payment metadata and process payment notifications idempotently.
- **FR-011**: System MUST expose owner/admin order read models with lifecycle-safe detail.
- **FR-012**: System MUST allow eligible users to cancel pending orders.
- **FR-013**: System MUST allow buyers to request refunds and admins to approve or reject them.
- **FR-014**: System MUST allow admins to create, update, reorder, enable, disable, and delete products and categories.
- **FR-015**: System MUST keep product editorial/media save flows on `/v1/admin/products` and related form routes.
- **FR-016**: System MUST allow admins to list users, read account state, and block or unblock users.
- **FR-017**: System MUST provide structured admin/public store-settings read and write contracts without check-in or refund-reclaim fields.
- **FR-018**: System MUST reject removed routes including `/v1/admin/cards`, `/v1/admin/users/:id/points`, `/v1/profile/checkin`, and `/v1/user/profile/checkin-status`.
- **FR-019**: System MUST support review, wishlist, notification, and admin messaging flows that remain in the current product scope.
- **FR-020**: System MUST keep published REST and Swagger contracts aligned with the removed no-cards/no-points scope.

### Key Entities

- **User**: Buyer or admin identity with username, email, trust/access state, notification preference, admin eligibility, and blocked status.
- **Product**: Sellable catalog entity with pricing, visibility, stock summary, media/editorial data, and category membership.
- **Order**: Purchase record containing product snapshot, buyer identity or email, quantity, amount, payment reference, and lifecycle status.
- **Payment Confirmation**: Gateway result tied to an order and processed idempotently.
- **Refund Request**: Support request tied to an order with reason, status, and admin decision metadata.
- **Category**: Catalog grouping with ordering and parent-child structure.
- **Review**: Product review moderation entity.
- **Wishlist Item**: Optional engagement entity for user suggestions/votes.
- **Notification / Admin Message**: User-facing or broadcast communication records.
- **Setting**: Structured storefront configuration for brand, homepage, footer, floating support, and presence controls.

## Success Criteria

- **SC-001**: Published REST and Swagger contracts contain no `/v1/admin/cards`, no user-points mutation route, and no check-in routes.
- **SC-002**: Profile and user/admin read models contain no `points`, `pointsUsed`, `pointsToUse`, `checkinEnabled`, or `checkinReward` fields.
- **SC-003**: Product create/edit/media flows are fully represented by `/v1/admin/products` and product form readback routes.
- **SC-004**: `go test ./... -run '^$'` completes successfully after the removal pass.
- **SC-005**: Focused checkout/orders/profile/admin suites pass against the no-cards/no-points contract.
- **SC-006**: Frontend Playwright product/user/customer/profile/admin-product/admin-user slices pass against the same contract.

## Assumptions

- No backward-compatibility read model is required for cards, points, or check-in fields.
- Refund processing no longer attempts card reclaim or points restoration.
- Historical planning documents mentioning cards/points/check-in are archival only unless explicitly refreshed by new evidence.
