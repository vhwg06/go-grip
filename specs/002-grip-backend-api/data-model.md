# Data Model: Grip Store Backend API

## Modeling Rules

- Monetary values are represented as integer-compatible amounts in backend/domain code.
- Removed domains `Card` and `DailyCheckin` are not part of the active model.
- Removed loyalty fields are not carried on `User` or `Order`.

## User

**Purpose**: Authentication, account identity, admin eligibility, blocked state, and profile preferences.

**Fields**:

- ID
- Username
- Email
- Trust level / role
- Admin eligibility
- Blocked status
- Desktop notification preference
- Last login time
- Created at
- Updated at

**Notes**:

- No `points`
- No `last_checkin_at`
- No check-in streak/reward state

## Product

**Purpose**: Sellable catalog entity and product-editor source of truth.

**Fields**:

- Product ID
- Title / name
- SKU
- Description
- Category ID
- Price
- Active/visibility state
- Purchase limit
- Stock summary
- Media references
- Editorial metadata
- Created at
- Updated at

## Category

**Purpose**: Catalog grouping and ordering.

**Fields**:

- Category ID
- Name
- Slug
- Parent ID
- Position
- Created at
- Updated at

## Order

**Purpose**: Purchase lifecycle record.

**Fields**:

- Order ID
- Product snapshot
- Quantity
- Amount
- Buyer user ID or email
- Payment reference
- Status
- Paid at
- Delivered at
- Created at
- Updated at

**Notes**:

- No `pointsUsed`
- No card-key delivery payload

## Payment

**Purpose**: External payment interaction attached to an order.

**Fields**:

- Payment ID
- Order ID
- Provider metadata
- Status
- Callback payload / audit metadata
- Created at
- Updated at

## Refund Request

**Purpose**: Buyer support request and admin decision record.

**Fields**:

- Refund ID
- Order ID
- Reason
- Status
- Admin note
- Processed by
- Processed at
- Created at
- Updated at

**Notes**:

- Approval/rejection changes refund state only; it does not reclaim cards or restore points.

## Review

**Purpose**: Product review moderation and public reflection.

**Fields**:

- Review ID
- Product ID
- Order ID
- Author identity
- Rating
- Comment
- Moderation status
- Created at
- Updated at

## Wishlist Item

**Purpose**: Optional engagement item with user voting.

## Notification / Admin Message

**Purpose**: Inbox and broadcast communication.

## Setting

**Purpose**: Structured storefront configuration.

**Sections**:

- Brand
- Contact
- Homepage
- Footer
- Floating support
- Visibility
- Registry
- Banner/About presence

**Notes**:

- No `checkinEnabled`
- No `checkinReward`
- No `refundReclaimCards`
