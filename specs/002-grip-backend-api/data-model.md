# Data Model: Grip Store Backend API

## Modeling Rules

- Domain entities live in `internal/entity` without GORM tags.
- GORM persistence models live under `internal/repo/persistent/models`.
- Repository implementations translate between persistence models and domain entities.
- PostgreSQL is the source of truth for production and integration tests.
- Monetary values and points are represented as integer minor units where possible to avoid floating-point drift.

## Entities

### User

**Fields**: ID, provider identities, username, email, points balance, trust level, admin eligibility, blocked status, desktop notification preference, last login time, last check-in date, consecutive check-in days, created time, updated time.

**Relationships**: Has orders, reviews, refund requests, check-ins, notifications, wishlist items, wishlist votes, broadcast reads, refresh sessions.

**Validation**: Username is unique case-insensitively for login/admin matching. Email is normalized where available. Blocked users cannot perform mutating buyer actions.

### Refresh Session

**Fields**: ID, user ID, token identifier, expires at, revoked at, created at.

**Relationships**: Belongs to User.

**Validation**: Refresh credentials are single-use. Expired or revoked credentials cannot renew access.

### Product

**Fields**: ID, name, description, price, compare-at price, category ID, image URL, hot flag, active flag, shared-stock flag, sort order, purchase limit, purchase warning, visibility level, stock count, locked count, sold count, rating, review count, created at, updated at.

**Relationships**: Has cards, reviews, orders, category.

**Validation**: Active products only appear in buyer catalog. Visibility level must be enforced on list and direct detail views. Purchase limit must be enforced per user or purchase email.

### Category

**Fields**: ID, name, icon, sort order, created at, updated at.

**Relationships**: Has products.

**Validation**: Name is required. Sort order controls public display.

### Card

**Fields**: ID, product ID, card key, used flag, reserved order ID, reserved at, expires at, used at, created at.

**Relationships**: Belongs to Product and optionally to reserved/delivered Order.

**Validation**: Non-shared cards may be reserved by at most one pending order. Expired, used, or currently reserved cards are unavailable. Shared product cards are reusable for delivery and are not marked used.

### Order

**Fields**: Order ID, product ID, product name snapshot, amount, email, status, payment trade number, card key or card IDs, user ID, username, payee, points used, quantity, current payment ID, paid at, delivered at, created at, updated at.

**Relationships**: Belongs to Product and optionally User; has refund requests and reviews.

**Validation**: Card keys are visible only in delivered state to owner or admin. Pending orders may be cancelled before payment. Delivered orders may receive refund requests.

**State Transitions**:

```text
pending -> paid -> delivered
pending -> cancelled
pending -> failed
delivered -> refund_pending -> refunded
refund_pending -> delivered
```

### Payment

**Fields**: ID, order ID, provider, provider payment ID, amount, status, request payload summary, confirmation payload summary, signature status, processed at, created at.

**Relationships**: Belongs to Order.

**Validation**: Successful confirmations must be authentic and idempotent. A payment success cannot deliver the same order twice.

### Refund Request

**Fields**: ID, order ID, user ID, username, reason, status, admin username, admin note, processed at, created at, updated at.

**Relationships**: Belongs to Order and User.

**Validation**: Only delivered orders can enter refund pending. Approval returns eligible points and reclaims eligible card stock.

### Daily Check-in

**Fields**: ID, user ID, check-in date, reward amount, streak after check-in, created at.

**Relationships**: Belongs to User.

**Validation**: Unique per user per calendar day.

### Review

**Fields**: ID, product ID, order ID, user ID, username, rating, comment, created at, updated at.

**Relationships**: Belongs to Product, Order, and User.

**Validation**: Rating is 1 through 5. One review per eligible delivered order.

### Wishlist Item

**Fields**: ID, title, description, user ID, username, vote count, created at, updated at.

**Relationships**: Belongs to User; has wishlist votes.

**Validation**: Owner or admin may delete. Blocked users cannot create or vote.

### Wishlist Vote

**Fields**: ID, item ID, user ID, created at.

**Relationships**: Belongs to Wishlist Item and User.

**Validation**: Unique per item and user.

### Notification

**Fields**: ID, user ID, type, title key, content key, payload data, read flag, created at.

**Relationships**: Belongs to User.

**Validation**: User can mark own notifications read or clear own inbox.

### Broadcast Message

**Fields**: ID, title key, content key, payload data, sender, created at.

**Relationships**: Has broadcast reads.

**Validation**: Unread count combines personal unread and unread broadcasts after clear timestamp.

### Broadcast Read

**Fields**: ID, message ID, user ID, created at.

**Relationships**: Belongs to Broadcast Message and User.

**Validation**: Unique per message and user.

### Setting

**Fields**: Key, value, updated at.

**Relationships**: Referenced by feature gates and runtime business rules.

**Validation**: Known settings validate type and allowed range, including check-in enablement, wishlist enablement, reward amount, theme, announcement, and payment configuration.

### Admin Message

**Fields**: ID, target type, target value, title, body, sender, created at.

**Relationships**: May create personal notifications or broadcast messages.

**Validation**: Sender must be admin. Target type must be all, user ID, or username.

## Repository Boundaries

- **AuthRepository**: users, provider identities, refresh sessions, account merge.
- **CatalogRepository**: products, categories, visibility filtering, stock summaries, review summaries.
- **CheckoutRepository**: transactional order creation, card reservation, point deduction, payment records.
- **OrderRepository**: owner/admin order reads, status transitions, cancellation, delivery, refund state.
- **ProfileRepository**: profile updates, check-ins, points balance.
- **WishlistRepository**: wishlist items, votes, reviews.
- **NotificationRepository**: personal notifications, broadcasts, read/clear state.
- **AdminRepository**: product/category/card/user/settings/message management and aggregate repair.

## Critical Constraints

- Use database transactions for order creation, payment success, cancellation, refund approval, and stock repair.
- Use row-level locking or equivalent transaction-safe selection for card reservation.
- Use unique constraints for refresh token identifiers, daily check-ins, reviews per order, wishlist votes, and broadcast reads.
- Do not expose `card_key` outside delivered owner/admin responses.
- Do not rely on process-local maps for production state.
