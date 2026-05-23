# Feature Specification: Grip Store Backend API

**Feature Branch**: `002-grip-backend-api`

**Created**: 2026-05-23

**Status**: Draft

**Input**: User description: "implement based on docs/epic.md and docs/mock-plan.md"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Browse and Buy Digital Goods (Priority: P1)

Buyers and anonymous visitors browse eligible products, review product details, preview the final purchase price, create an order, complete payment when required, and receive card keys only after successful delivery.

**Why this priority**: This is the primary revenue flow and the minimum usable release for Grip Store.

**Independent Test**: Can be fully tested by publishing a product with stock, completing a paid and a points-covered order, and confirming the buyer receives the correct order status and card key visibility.

**Acceptance Scenarios**:

1. **Given** an anonymous visitor and a public active product, **When** the visitor views the catalog and product detail, **Then** only products visible to anonymous users are shown with correct price, stock, review, and purchase limit information.
2. **Given** a signed-in buyer with points, **When** the buyer previews checkout, **Then** the system returns the total price, points applied, and final payable amount before order creation.
3. **Given** an order with a payable balance, **When** payment is confirmed, **Then** the order advances through paid and delivered states and reveals the purchased card key.
4. **Given** an order fully covered by points, **When** the buyer creates the order, **Then** points are deducted and delivery occurs without sending the buyer to an external payment step.

---

### User Story 2 - Authenticate and Manage Account Benefits (Priority: P1)

Users sign in with supported external accounts, maintain their profile, check in daily for points, spend points during checkout, and see order history tied to their account or purchase email.

**Why this priority**: Account identity, points, and order ownership are required for trust-level catalog access, loyalty incentives, and purchase recovery.

**Independent Test**: Can be fully tested by signing in with each supported provider, refreshing a session, checking in once per day, and verifying points and order history update correctly.

**Acceptance Scenarios**:

1. **Given** a new visitor, **When** they sign in with a supported account provider, **Then** a user profile is created or merged with an existing matching account and the client receives authenticated access.
2. **Given** a signed-in user, **When** they request their profile, **Then** the response includes identity, points balance, notification preference, check-in status, order statistics, and admin eligibility.
3. **Given** a user who has not checked in today, **When** they check in, **Then** their streak and points are updated exactly once for that day.
4. **Given** an expired access session with a valid renewal credential, **When** the client renews access, **Then** the old renewal credential can no longer be reused.

---

### User Story 3 - Protect Stock, Payments, and Order Lifecycle (Priority: P1)

The store prevents overselling, reserves stock for a limited time, processes payment confirmations safely, cancels stale orders, refunds points when appropriate, and exposes clear order states to the frontend.

**Why this priority**: Stock and payment correctness directly protect revenue, customer trust, and card key secrecy.

**Independent Test**: Can be fully tested by racing multiple purchases against limited stock, replaying payment confirmations, cancelling pending orders, and verifying no duplicate delivery or lost points occurs.

**Acceptance Scenarios**:

1. **Given** multiple buyers competing for limited cards, **When** they create orders concurrently, **Then** no more orders reserve cards than available stock.
2. **Given** a pending order older than the reservation window, **When** cleanup runs or a buyer cancels, **Then** reserved cards are released and used points are returned.
3. **Given** the same successful payment confirmation is received more than once, **When** it is processed repeatedly, **Then** the order is delivered once and no duplicate points, stock, or card changes occur.
4. **Given** an order is not delivered, **When** order details are viewed, **Then** card keys remain hidden while status text and color remain available for client display.

---

### User Story 4 - Admin Operates the Store (Priority: P2)

Admins manage products, categories, card inventory, orders, refunds, users, settings, messages, and data repair operations from a privileged management experience.

**Why this priority**: Admin capabilities are needed to launch, maintain stock, handle support issues, and keep business configuration current.

**Independent Test**: Can be fully tested by signing in as an admin, performing each management action, and confirming non-admin users are denied access.

**Acceptance Scenarios**:

1. **Given** an admin user, **When** they create or update a product, category, setting, or announcement, **Then** the change is persisted and reflected in buyer-facing responses.
2. **Given** an admin importing card inventory, **When** they submit a bulk card list, **Then** each valid card is stored and stock counts reflect the imported inventory.
3. **Given** a delivered order with a refund request, **When** an admin approves the refund, **Then** points are refunded, cards are reclaimed when eligible, and stock summaries are updated.
4. **Given** a blocked user, **When** that user attempts a mutating buyer action, **Then** the action is rejected.

---

### User Story 5 - Engage Users with Wishlist, Reviews, and Notifications (Priority: P3)

Signed-in users create wishlist ideas, vote, review purchased products, and receive a unified inbox that combines personal and broadcast notifications.

**Why this priority**: These features improve engagement and trust but are secondary to the core buying and administration flows.

**Independent Test**: Can be fully tested by creating a delivered order, submitting one review for it, creating and voting on wishlist items, sending personal and broadcast messages, and verifying unread counts.

**Acceptance Scenarios**:

1. **Given** a delivered order that has not been reviewed, **When** the buyer submits a rating and comment, **Then** the review is accepted and product review summaries are recalculated.
2. **Given** a wishlist item, **When** a signed-in user toggles their vote, **Then** the vote count changes once for that user.
3. **Given** personal and broadcast notifications exist, **When** a user opens their inbox, **Then** the latest relevant messages and accurate unread count are returned together.
4. **Given** wishlist functionality is disabled by store settings, **When** a user attempts to create or vote on wishlist items, **Then** the action is rejected with a clear disabled-feature result.

### Edge Cases

- Products with restricted visibility are excluded from listings and rejected on direct access when the user does not meet the required trust level.
- Shared products show effectively unlimited stock only when at least one usable card exists, and delivery does not consume the shared source card.
- Purchase limits account for prior paid and delivered orders for the same user or purchase email.
- Payment callbacks with invalid signatures, unknown orders, failed payments, or repeated success confirmations do not deliver cards.
- Orders cancelled by users, timeout, or admin action release reserved stock and return unused points where applicable.
- Expired cards are removed from availability and stock summaries are recalculated.
- Blocked users can read allowed public data but cannot create wishlist items, reviews, orders tied to their account, or other mutating user content.
- Notification clearing handles personal messages and broadcast read tracking without resurrecting previously cleared broadcasts.
- Account merges preserve points, orders, reviews, notifications, messages, wishlist activity, and blocked status without duplicate votes or reads.
- Admin-only operations are rejected for non-admin users even if they are signed in.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST authenticate users through LinuxDO and GitHub sign-in and support account merge when provider identities match an existing email or username.
- **FR-002**: System MUST issue renewable authenticated sessions, rotate renewal credentials after each use, and support logout that invalidates future renewal.
- **FR-003**: System MUST identify admin users by configured usernames and include admin eligibility in the current-user profile.
- **FR-004**: System MUST expose buyer profile information including points balance, order statistics, notification preference, check-in status, and blocked status effects.
- **FR-005**: System MUST list and search active products using catalog visibility rules for anonymous and signed-in users.
- **FR-006**: System MUST provide product detail data including display stock, reviews summary, purchase warnings, and maximum purchasable quantity.
- **FR-007**: System MUST provide public category, settings, feature flag, theme, and announcement data needed by the client.
- **FR-008**: System MUST calculate checkout previews using product price, quantity, buyer points, and the rule that points applied cannot exceed the buyer balance or rounded order total.
- **FR-009**: System MUST create orders only when requested quantity, product visibility, product availability, purchase limits, buyer status, and stock reservation rules all pass.
- **FR-010**: System MUST reserve non-shared card stock for a 5-minute pending order window and release reservations on cancellation, timeout, failed payment, or refund where applicable.
- **FR-011**: System MUST support shared products by allowing delivery from a reusable card key without marking that card as consumed.
- **FR-012**: System MUST deliver orders immediately when the final payable amount is zero after points are applied.
- **FR-013**: System MUST generate payment instructions for payable orders and process payment confirmations only after validating the confirmation authenticity.
- **FR-014**: System MUST process payment confirmations idempotently so repeated confirmations cannot duplicate delivery, points changes, or stock changes.
- **FR-015**: System MUST expose order lists and order details for the owning user or purchase email, with pagination where lists may grow.
- **FR-016**: System MUST hide card keys unless the order is delivered and visible to the requesting owner or authorized admin.
- **FR-017**: System MUST include client-ready order status labels in Vietnamese and status colors for all order responses.
- **FR-018**: System MUST allow pending orders to be cancelled by eligible users and refund any points used on the cancelled order.
- **FR-019**: System MUST allow buyers to request refunds for delivered orders and allow admins to approve or reject those requests.
- **FR-020**: System MUST process approved refunds by returning eligible points, reclaiming eligible card stock, updating order state, and recording admin decision details.
- **FR-021**: System MUST allow signed-in users to update email and desktop notification preference.
- **FR-022**: System MUST allow daily check-in only once per user per calendar day and update streak and rewards according to store settings.
- **FR-023**: System MUST disable check-in, wishlist, or other configurable features when the corresponding store setting is off.
- **FR-024**: System MUST allow wishlist item creation, update, deletion by owner or admin, and one vote toggle per signed-in user.
- **FR-025**: System MUST allow reviews only for delivered orders and limit reviews to one per eligible order.
- **FR-026**: System MUST recalculate product rating and review count whenever reviews are created or changed.
- **FR-027**: System MUST provide a unified inbox combining personal notifications and broadcast messages with accurate unread counts.
- **FR-028**: System MUST support marking individual notifications read, marking all read, and clearing inbox state.
- **FR-029**: System MUST allow admins to create, update, reorder, enable, disable, and delete products and categories.
- **FR-030**: System MUST allow admins to import, delete, inspect, and replenish card inventory and repair aggregate stock counts.
- **FR-031**: System MUST allow admins to view all orders, mark orders paid or delivered, cancel orders, delete orders, and manage pending refunds.
- **FR-032**: System MUST allow admins to list users, adjust points, block or unblock users, and send targeted or broadcast messages.
- **FR-033**: System MUST provide system health visibility and scheduled maintenance for expired orders, expired cards, and stock summaries.
- **FR-034**: System MUST rate limit sensitive flows including sign-in and checkout to reduce abuse while preserving normal buyer usage.
- **FR-035**: System MUST reject requests that would expose protected business data, card keys, restricted products, or admin functions to unauthorized users.
- **FR-036**: System MUST keep responses compatible with the already-refactored client contract for all buyer, profile, order, notification, and admin flows.

### Key Entities *(include if feature involves data)*

- **User**: A buyer or admin identity with username, email, points balance, trust level, blocked status, notification preference, check-in streak, and sign-in provider associations.
- **Product**: A digital good for sale with name, description, price, category, image, active state, hot flag, shared-stock flag, purchase limit, visibility level, stock summary, rating, and review count.
- **Card**: A deliverable digital key associated with a product, including used state, reservation state, expiration, and delivery eligibility.
- **Order**: A purchase record containing product snapshot, buyer identity or email, quantity, amount, points used, payment reference, lifecycle status, delivery timestamp, and card delivery data.
- **Payment Confirmation**: A gateway result tied to an order, including authenticity proof, success or failure state, and processing outcome.
- **Refund Request**: A buyer support request for a delivered order, including reason, status, admin decision, and processing timestamps.
- **Category**: A product grouping with display name, icon, ordering, and active catalog usage.
- **Review**: A buyer rating and comment tied to one delivered order and product.
- **Wishlist Item**: A user-created request or suggestion with owner, description, creation time, and vote activity.
- **Notification**: A personal or broadcast message with translation key, payload data, read state, and clear state.
- **Setting**: Store configuration controlling theme, announcement, feature availability, rewards, payment behavior, and integrations.
- **Admin Message**: A targeted or broadcast communication sent by an admin to users.
- **Maintenance Job**: A recurring operational task that cancels expired orders, removes expired cards, or repairs aggregate summaries.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of client purchase, account, order, notification, and admin contract calls used by the refactored frontend return compatible success or error shapes in smoke testing.
- **SC-002**: 95% of catalog, product detail, profile, and order status views return user-visible results in under 1 second under normal production traffic.
- **SC-003**: 95% of order creation, cancellation, refund approval, and admin inventory changes complete in under 2 seconds under normal production traffic.
- **SC-004**: Oversold stock incidents remain at 0 across concurrent purchase tests with at least 50 simultaneous buyers competing for limited inventory.
- **SC-005**: Replayed successful payment confirmations produce exactly one delivered order in 100% of replay tests.
- **SC-006**: Pending orders older than 5 minutes are cancelled and release stock within 2 minutes of the next maintenance cycle in 99% of cases.
- **SC-007**: Card keys are hidden before delivery and visible after delivery in 100% of authorization tests.
- **SC-008**: At least 99% of daily check-in attempts produce the correct once-per-day result, including duplicate attempts on the same day.
- **SC-009**: Admin users can complete core operations for product update, card import, order cancellation, refund decision, user block, and broadcast message in under 3 minutes each during acceptance testing.
- **SC-010**: Monthly availability for buyer browsing, checkout, and order retrieval flows is at least 99.9%.
- **SC-011**: No lost orders, points, or card reservations are found in audit checks after payment, cancellation, timeout, and refund test runs.

## Assumptions

- The existing frontend contract is stable and this feature must satisfy it without requiring frontend behavior changes.
- The first release serves the existing web client; mobile and third-party clients are future consumers but not separate launch deliverables.
- Existing store data and column meanings are retained so current products, users, orders, cards, and settings continue to work after migration.
- Epay remains the payment provider for this feature; payment provider replacement is outside scope.
- Product images continue to be referenced as URLs; upload, resize, and CDN workflows are outside scope.
- Real-time notifications are outside scope; clients may poll for status and inbox updates.
- Buyer-facing order status labels are Vietnamese for this release to match current client behavior.
- Advanced analytics, BI reporting, and redesigned email templates are outside scope.
