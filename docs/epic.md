# Epic PRD: Grip Store Backend REST API Server

## 1. Epic Name

**Grip Store Backend REST API Server** — Xây dựng backend server tách biệt phục vụ toàn bộ nghiệp vụ cho Grip Store (virtual goods e-commerce platform).

---

## 2. Goal

### Problem

Grip Store hiện đang chạy dưới dạng Next.js full-stack monolith: server actions + Drizzle ORM + SQLite trực tiếp từ frontend. Kiến trúc này gây ra:
- Frontend chứa toàn bộ logic nghiệp vụ (stock reservation, payment, points calculation) → khó maintain, không scale được.
- Không thể tách team FE/BE, deploy độc lập, hoặc phục vụ multiple clients (mobile app, admin panel riêng).
- Security risk: business logic expose trên client-side bundle.

### Solution

Tách toàn bộ business logic ra một REST API backend server độc lập. Frontend (đã refactor xong) chỉ còn là pure client app gọi REST endpoints. Backend chịu 100% trách nhiệm: auth, stock management, payment processing, order state machine, points/gamification, notifications, admin operations.

### Impact

- **Separation of Concerns**: FE team và BE team làm việc độc lập, deploy riêng.
- **Security**: Business logic không leak ra client, card keys bảo mật server-side.
- **Scalability**: Backend có thể horizontal scale, cache, replicate độc lập.
- **Multi-client ready**: Mobile app, third-party integrations có thể dùng cùng API.
- **Time to market**: FE đã mock xong contract → BE chỉ cần implement đúng spec.

---

## 3. User Personas

### Buyer (Người mua)
- Truy cập shop, xem catalog, mua thẻ ảo/digital goods.
- Thanh toán qua Epay gateway hoặc điểm tích lũy.
- Nhận card key sau khi thanh toán, xem lịch sử orders.
- Điểm danh hàng ngày để tích điểm, vote wishlist.

### Admin (Quản trị viên)
- Quản lý sản phẩm, thẻ kho (cards inventory).
- Giám sát đơn hàng, duyệt/từ chối hoàn tiền.
- Quản lý users (block/unblock, chỉnh points).
- Cấu hình hệ thống (theme, announcements, notifications).
- Gửi broadcast messages.

### Anonymous Visitor (Khách vãng lai)
- Xem catalog công khai (products với visibilityLevel ≤ -1).
- Mua hàng bằng email (không cần đăng nhập) nếu sản phẩm cho phép.

---

## 4. High-Level User Journeys

### Journey 1: Browse & Purchase
```
Khách truy cập → Xem catalog (filtered by visibility) → Chọn sản phẩm 
→ Xem chi tiết (stock, price, reviews) → Điền form mua (quantity, email, use points?) 
→ Backend tính toán giá (preview) → Xác nhận mua → Backend khóa thẻ kho 
→ [Nếu > 0 đồng] Redirect Epay gateway → Thanh toán → Webhook callback → Giao thẻ
→ [Nếu = 0 đồng] Giao thẻ ngay lập tức
→ Xem card key tại trang order detail
```

### Journey 2: Auth & Account
```
User bấm Login → Redirect OAuth (LinuxDO hoặc GitHub) → Backend xử lý callback
→ Merge account nếu trùng email/username → Issue JWT tokens → Frontend lưu tokens
→ Mọi request kèm Bearer token → Auto-refresh khi hết hạn
```

### Journey 3: Daily Check-in & Points
```
User đăng nhập → Vào Profile → Bấm Điểm danh → Backend atomic update streak
→ Cộng points → Hiển thị streak count → Dùng points khi mua hàng
```

### Journey 4: Admin Management
```
Admin login → Dashboard (revenue, pending orders, expired cards)
→ Manage products (CRUD, toggle, reorder)
→ Import cards (bulk text)
→ Process orders (mark paid/delivered, cancel)
→ Handle refunds (approve → reclaim cards + refund points)
→ Configure settings (shop name, theme, notifications)
```

### Journey 5: Order Lifecycle
```
pending → [payment success] → paid → [cards delivered] → delivered
pending → [user cancel / timeout 5min] → cancelled (release cards + refund points)
pending → [payment failure] → failed
delivered → [user request] → refund_pending → [admin approve] → refunded
```

---

## 5. Business Requirements

### Functional Requirements

#### Module 1: Authentication & Authorization
- OAuth login via LinuxDO (`connect.linux.do`) và GitHub
- Issue JWT access token (15min TTL) + refresh token (7 days)
- Token refresh rotation (single-use refresh tokens)
- Logout (invalidate tokens)
- `GET /api/auth/me` trả về user profile kèm `isAdmin` flag

#### Module 2: Catalog & Products
- List products filtered by `visibilityLevel ≤ threshold` (threshold = -1 cho anonymous, max(0, trustLevel) cho logged-in)
- Stock hiển thị từ aggregate cache fields (`stockCount - lockedCount`), KHÔNG đếm live từ bảng cards
- Shared product stock = `999999` (INFINITE) nếu có ≥ 1 card khả dụng, ngược lại = 0
- Product detail trả kèm `maxPurchaseableQuantity` (tính server-side dựa trên purchaseLimit và lịch sử mua)
- Search: full-text match trên name + description, filter by category, sort (hot/price/stock)
- Categories, settings (theme, features), announcement endpoints
- API `403 Forbidden` nếu user access product có visibilityLevel cao hơn threshold

#### Module 3: Checkout & Payment
- **Checkout Preview**: Server tính `finalPrice = max(0, quantity * price - pointsToUse)` với `pointsToUse = min(userPoints, ceil(totalAmount))`
- **Create Order**: Atomic stock reservation (5 phút TTL), fallback grab expired reservations
- **Shared Product**: Không lock card, lấy random card key khi deliver, không set `isUsed = 1`
- **Zero-Price Flow**: Nếu `finalPrice ≤ 0` → trừ points atomic → deliver ngay → không qua gateway
- **Epay Integration**: MD5 signature generation, payment URL params
- **Webhook Handler**: Verify MD5 signature, idempotent processing, deliver cards on success
- **Order Cancel**: Release reserved cards, refund points if used
- Points deduction với atomic guard (`WHERE points >= :amount`), rollback on failure

#### Module 4: Orders
- List user orders (by userId OR email), paginated
- Order detail: card key chỉ hiển thị khi status = `delivered`
- Status polling cho payment waiting page
- Refund request submission (only for `delivered` orders)
- **Dumb Frontend mapping**: Backend trả `statusText` (Vietnamese label) + `statusColor` (CSS color) cho mọi order

#### Module 5: Profile & Gamification
- Profile dashboard: user info + points balance + order stats + check-in status
- Update email, toggle desktop notifications
- **Daily Check-in**: Atomic SQL với CASE WHEN để tính streak, UTC timezone normalization
- Anti-spam: WHERE clause prevents double check-in same day
- Configurable reward amount từ settings (`checkin_reward`)
- System setting `checkin_enabled` gates the feature

#### Module 6: Wishlist & Reviews
- Wishlist: CRUD items + toggle vote (INSERT/DELETE in `wishlist_votes`)
- Delete wishlist: owner hoặc admin
- Reviews: tied to delivered order + product, 1 review per order
- Rating 1-5, aggregate recalc on product (`AVG(rating)`, `COUNT(*)`)
- Blocked users cannot create wishlist/reviews
- System setting `wishlist_enabled` gates the feature

#### Module 7: Notifications
- **Unified Inbox**: Merge `user_notifications` (personal) + `broadcast_messages` (system-wide)
- Unread count formula: `personal_unread + max(0, broadcast_total - broadcast_read)` (limit 10 latest broadcasts)
- Mark read individual, mark all read, clear all
- Clear operation: delete personal + record `broadcast_cleared_at` timestamp + INSERT reads for all current broadcasts
- Notification payload: translation keys + JSON data (frontend handles i18n rendering)

#### Module 8: Admin
- **Authorization**: `ADMIN_USERS` env var (comma-separated usernames, case-insensitive match)
- Products CRUD + toggle active + reorder (sort_order)
- Cards: bulk import (text → individual rows), delete, pull from external API
- Orders: view all, mark paid/delivered, cancel, delete
- Refunds: list pending, approve (refund points + reclaim cards + recalc stock), reject
- Users: list, update points, block/unblock
- Settings: key-value CRUD for all system configs
- Categories: CRUD with icon + sort_order
- Messages: send broadcast, send targeted (by username/userId/all)
- Notifications: test Telegram/Bark/Email
- Data: import SQLite backup, repair aggregates

#### Module 9: Background Jobs (Cron)
- **Cancel Expired Orders** (every 1 min): `pending` orders older than 5 min → cancel + release cards + refund points
- **Cleanup Expired Cards** (every 10 min): delete cards with `expires_at < now`, recalc product stock
- **Sync Aggregates** (every 5 min): recalculate `stock_count`, `locked_count`, `sold_count` for all products

### Non-Functional Requirements

- **Database**: SQLite WAL mode, single-file, tương thích 100% schema hiện tại (không đổi column names/types)
- **Performance**: Read endpoints < 100ms, write endpoints < 500ms
- **Concurrency**: Atomic UPDATE...RETURNING cho stock reservation, conditional WHERE cho points deduction
- **Idempotency**: Payment webhook processing phải idempotent (không duplicate delivery)
- **Security**: 
  - Parameterized SQL queries (no injection)
  - Card keys masked in API responses unless order is `delivered`
  - Rate limiting: checkout/auth = 5 req/s, default = 100 req/s
  - Visibility level enforced at query level
  - Blocked users cannot perform mutating operations
- **Reliability**: Points rollback on order creation failure
- **Observability**: Structured JSON logging (pino), error tracking
- **Deployment**: Docker container, volume mount cho SQLite file, env-based config
- **CORS**: Allow requests from frontend origin only

---

## 6. Success Metrics

| KPI | Target | Measurement |
|-----|--------|-------------|
| API Coverage | 100% endpoints implemented | All FE adapter calls succeed |
| Contract Compliance | 0 breaking changes vs FE contracts | FE build + smoke test pass |
| Response Time (p95) | < 200ms reads, < 1s writes | Application metrics |
| Payment Success Rate | ≥ 99% webhook processing | Webhook logs |
| Stock Consistency | 0 oversold incidents | Aggregate sync + atomic locks |
| Uptime | 99.9% | Health check monitoring |
| Zero Data Loss | 0 lost orders/points | Transaction logs audit |

---

## 7. Out of Scope

- **Frontend changes**: FE đã refactor xong, contract đã ổn định. Không thay đổi FE code.
- **Database schema changes**: Giữ nguyên schema hiện tại. Chỉ thêm indexes nếu cần.
- **Mobile app**: Backend phục vụ web frontend trước. Mobile client là epic sau.
- **Real-time features**: WebSocket/SSE cho live notifications. Hiện tại dùng polling.
- **Multi-language API responses**: `statusText` trả tiếng Việt hardcoded (matching current behavior). i18n API là epic sau.
- **Payment gateway migration**: Vẫn dùng Epay. Đổi gateway là epic riêng.
- **CDN/Image upload**: Product images là URL string. Upload/resize là scope riêng.
- **Analytics/Reporting**: Dashboard stats cơ bản. BI/advanced reporting là scope riêng.
- **Email template design**: Gửi plain text/simple HTML email. Template engine là scope riêng.

---

## 8. Business Value

**High** — Đây là blocker chính để production deploy. Frontend đã sẵn sàng (adapter layer mock), nhưng không có backend thì không ship được. Epic này unblock:
- Production deployment
- Team separation (FE/BE parallel work)
- Security hardening (business logic off-client)
- Future scalability (mobile, third-party)

---

## Appendix A: Frontend Contracts (Đã Implement)

Dưới đây là TypeScript types và API calls mà frontend đã code xong. Backend **PHẢI** trả về response tương thích với các types này.

### Auth Domain Types

```typescript
// src/domain/auth.ts
interface User {
  id: string
  username: string | null
  email: string | null
  avatar_url: string | null
  trustLevel: number
  isAdmin: boolean
  points: number
  desktopNotificationsEnabled: boolean
}

interface AuthTokens {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: User
}
```

### Catalog Domain Types

```typescript
// src/domain/catalog.ts
interface CatalogProduct {
  id: string
  name: string
  description: string | null
  price: string                    // Money as string
  compareAtPrice: string | null
  image: string | null
  category: string | null
  isHot: boolean
  isShared: boolean
  purchaseLimit: number | null
  purchaseWarning: string | null
  visibilityLevel: number
  stock: number                    // Computed: stockCount - lockedCount (or INFINITE for shared)
  sold: number                     // soldCount
  rating: number
  reviewCount: number
}

interface CatalogCategory {
  id?: number
  name: string
  icon: string | null
  sortOrder: number
}

interface CatalogSettings {
  shopName: string
  shopDescription: string | null
  shopLogo: string | null
  shopFooter: string | null
  themeColor: string
  noindexEnabled: boolean
  wishlistEnabled: boolean
  checkinEnabled: boolean
  checkinReward: number
  lowStockThreshold: number
}

interface CatalogProductsResponse {
  items: CatalogProduct[]
  page: number
  limit: number
  total: number
}

interface CatalogSearchParams {
  q?: string
  category?: string
  page?: number
  limit?: number
  sort?: string   // "default" | "price_asc" | "price_desc" | "sales" | "rating"
}

// GET /api/catalog/products/:id response
interface CatalogProductViewState {
  product: CatalogProduct | null
  requiredLevel?: number | null   // If user lacks permission
}
```

### Checkout Domain Types

```typescript
// src/domain/checkout.ts
interface CheckoutOrderInput {
  productId: string
  quantity: number
  email?: string
  usePoints?: boolean
}

interface CheckoutPaymentResult {
  success: boolean
  url?: string                     // Epay gateway URL
  params?: Record<string, string | number | boolean | null | undefined>  // Form params for POST
  orderId?: string
  status?: string
  amount?: string
  pointsUsed?: number
  isZeroPrice?: boolean            // true = instant delivery, no payment redirect
  error?: string
}

interface BuyPageMeta {
  reviews: BuyMetaReview[]
  averageRating: number
  reviewCount: number
  canReview: boolean
  reviewOrderId?: string           // Order ID eligible for review
  emailConfigured: boolean         // User has email set
}

interface PaymentOrderInput {
  amount: number | string
  payee?: string | null
}
```

### Orders Domain Types

```typescript
// src/domain/orders.ts
type OrderStatus = "pending" | "paid" | "delivered" | "failed" | "refunded" | "cancelled"

interface OrderSummary {
  orderId: string
  productId: string | null
  productName: string
  amount: string
  status: OrderStatus | null
  createdAt: string | null
  canReview?: boolean              // Has delivered order without review
}

interface OrderDetail extends OrderSummary {
  cardKey: string | null           // NULL unless status = "delivered"
  payee?: string | null
  paidAt: string | null
}

interface OrderDetailResponse {
  order: OrderDetail | null
  canViewKey: boolean
  isOwner: boolean
  refundRequest: { status: string | null; reason: string | null } | null
}

interface OrdersListResponse {
  items: OrderSummary[]
  page: number
  total: number
}

interface OrderActionResult {
  success: boolean
  status?: OrderStatus
  error?: string
}
```

### Profile Domain Types

```typescript
// src/domain/profile.ts
interface ProfileView {
  user: {
    id: string
    name: string
    username: string | null
    avatar: string | null
    email: string | null
    trustLevel?: number
  }
  points: number
  checkinEnabled: boolean
  orderStats: { total: number; pending: number; delivered: number }
  notifications: ProfileNotification[]
  desktopNotificationsEnabled: boolean
}

interface CheckinStatus {
  checkedIn: boolean
  disabled?: boolean
  consecutiveDays?: number
  lastCheckinAt?: number | null
}

interface CheckinResult {
  success: boolean
  error?: string
  points?: number
  checkedIn?: boolean
  consecutiveDays?: number
}
```

### Wishlist & Reviews Domain Types

```typescript
// src/domain/wishlist.ts
interface WishlistItem {
  id: number
  title: string
  description?: string | null
  username?: string | null
  createdAt?: number | null
  votes: number
  voted: boolean                   // Whether current user has voted
}

interface WishlistResponse {
  items: WishlistItem[]
  enabled: boolean
}

interface ToggleWishlistVoteResult {
  success: boolean
  voted?: boolean
  count?: number
  error?: string
}

interface ProductReview {
  id: number
  username: string
  userId?: string | null
  rating: number
  comment: string | null
  createdAt: string | null
}

interface ProductReviewsResponse {
  reviews: ProductReview[]
  averageRating: number
  reviewCount: number
}

interface SubmitReviewInput {
  orderId: string
  rating: number
  comment: string
}
```

### Notifications Domain Types

```typescript
// src/domain/notifications.ts
interface NotificationItem {
  id: number
  type: string
  titleKey: string
  contentKey: string
  data: string | null             // JSON string with template variables
  isRead: boolean | null
  createdAt: number | null
}

interface NotificationsResponse {
  success: boolean
  items: NotificationItem[]
  error?: string
}

interface NotificationCountResponse {
  success: boolean
  count: number
  error?: string
}

interface SendUserMessageInput {
  title: string
  body: string
}
```

### Admin Domain Types

```typescript
// src/domain/admin.ts
type AdminActionResult = {
  success: boolean
  error?: string
  count?: number
  [key: string]: unknown
}

interface AdminDashboardPayload {
  stats: any
  settingsMap: Record<string, string>
  visitorCount: number
  registryEnabled?: boolean
}

interface AdminProductsPayload {
  products: any[]
  lowStockThreshold: number
}

interface AdminMessagesPayload {
  history: any[]    // Sent messages
  inbox: any[]      // User messages to admin
}

interface AdminNotificationsSettings {
  telegramBotToken: string
  telegramChatId: string
  telegramLanguage: string
  telegramEnabled: boolean
  barkEnabled: boolean
  barkServerUrl: string
  barkDeviceKey: string
  resendApiKey: string
  resendFromEmail: string
  resendFromName: string
  resendEnabled: boolean
  emailLanguage: string
}

interface AdminCollectPayload {
  payLink: string
  payee: string | null
}
```

---

## Appendix B: Complete API Endpoint Registry

Frontend đã implement adapter calls tới tất cả endpoints dưới đây. Backend phải expose **chính xác** các paths này.

### Auth (`src/adapters/api/auth.api.ts`)
| Method | Path | FE Function |
|--------|------|-------------|
| GET | `/api/auth/oauth/linuxdo` | `loginWithLinuxDO()` → redirect |
| GET | `/api/auth/oauth/github` | `loginWithGitHub()` → redirect |
| POST | `/api/auth/logout` | `logout()` |
| POST | `/api/auth/refresh` | `refreshToken()` → `AuthResponse` |
| GET | `/api/auth/me` | `getMe()` → `User` |

### Catalog (`src/adapters/api/catalog.api.ts`)
| Method | Path | FE Function |
|--------|------|-------------|
| GET | `/api/catalog/products` | `getActiveProducts(params)` → `CatalogProductsResponse` |
| GET | `/api/catalog/products/:id` | `getProduct(id)` → `CatalogProductViewState` |
| GET | `/api/catalog/products/:id/buy-meta` | `getBuyPageMeta(id)` → `BuyPageMeta` |
| GET | `/api/catalog/search` | `getActiveProducts({q})` → `CatalogProductsResponse` |
| GET | `/api/catalog/categories` | `getCategories()` → `CatalogCategory[]` |
| GET | `/api/catalog/settings` | `getPublicSettings()` → `CatalogSettings` |
| GET | `/api/catalog/announcement` | `getAnnouncement()` |

### Checkout (`src/adapters/api/checkout.api.ts`)
| Method | Path | FE Function |
|--------|------|-------------|
| POST | `/api/checkout/orders` | `createOrder(input)` → `CheckoutPaymentResult` |
| POST | `/api/checkout/payment-orders` | `createPaymentOrder(input)` → `CheckoutPaymentResult` |
| GET | `/api/checkout/orders/:id/payment-params` | `getRetryPaymentParams(id)` → `CheckoutPaymentResult` |
| POST | `/api/checkout/notify` | _(webhook, no FE call)_ |
| GET | `/api/checkout/callback/:id` | _(redirect, no FE call)_ |

### Orders (`src/adapters/api/orders.api.ts`)
| Method | Path | FE Function |
|--------|------|-------------|
| GET | `/api/orders` | `getMyOrders(page)` → `OrdersListResponse` |
| GET | `/api/orders/:id` | `getOrder(id)` → `OrderDetailResponse` |
| GET | `/api/orders/:id/status` | `checkOrderStatus(id)` → `OrderActionResult` |
| POST | `/api/orders/:id/cancel` | `cancelPendingOrder(id)` → `OrderActionResult` |
| POST | `/api/orders/:id/refund-request` | `submitRefundRequest(id, reason)` → `OrderActionResult` |

### Profile (`src/adapters/api/profile.api.ts`)
| Method | Path | FE Function |
|--------|------|-------------|
| GET | `/api/profile` | `getProfile()` → `ProfileView` |
| PATCH | `/api/profile/email` | `updateProfileEmail(email)` |
| PATCH | `/api/profile/notifications` | `updateDesktopNotifications(enabled)` |
| GET | `/api/profile/points` | `getUserPoints()` → `number` |
| POST | `/api/profile/checkin` | `checkIn()` → `CheckinResult` |
| GET | `/api/profile/checkin/status` | `getCheckinStatus()` → `CheckinStatus` |

### Wishlist & Reviews (`src/adapters/api/wishlist.api.ts`)
| Method | Path | FE Function |
|--------|------|-------------|
| GET | `/api/wishlist` | `getWishlistItems(limit)` → `WishlistResponse` |
| POST | `/api/wishlist` | `submitWishlistItem(title, desc)` |
| POST | `/api/wishlist/:id/vote` | `toggleWishlistVote(id)` → `ToggleWishlistVoteResult` |
| DELETE | `/api/wishlist/:id` | `deleteWishlistItem(id)` |
| GET | `/api/products/:id/reviews` | `getProductReviews(id)` → `ProductReviewsResponse` |
| POST | `/api/products/:id/reviews` | `submitReview(input)` |

### Notifications (`src/adapters/api/notifications.api.ts`)
| Method | Path | FE Function |
|--------|------|-------------|
| GET | `/api/notifications` | `getMyNotifications()` → `NotificationsResponse` |
| GET | `/api/notifications/unread-count` | `getUnreadCount()` → `NotificationCountResponse` |
| POST | `/api/notifications/:id/read` | `markNotificationRead(id)` |
| POST | `/api/notifications/read-all` | `markAllNotificationsRead()` |
| POST | `/api/notifications/clear` | `clearMyNotifications()` |
| POST | `/api/messages/user` | `sendUserMessage(title, body)` |

### Admin (`src/adapters/api/admin.api.ts`)
| Method | Path | FE Function |
|--------|------|-------------|
| GET | `/api/admin/me` | `checkAdmin()` |
| GET | `/api/admin/dashboard` | `getAdminDashboard()` → `AdminDashboardPayload` |
| GET | `/api/admin/products` | `getAdminProducts()` → `AdminProductsPayload` |
| GET | `/api/admin/products/:id` | `getAdminProduct(id)` |
| GET | `/api/admin/products/:id/form` | `getAdminProductForm(id)` |
| GET | `/api/admin/products/new` | `getAdminProductForm()` |
| POST | `/api/admin/products` | `saveProduct(formData)` _(multipart)_ |
| DELETE | `/api/admin/products/:id` | `deleteProduct(id)` |
| PATCH | `/api/admin/products/:id/status` | `toggleProductStatus(id, isActive)` |
| PATCH | `/api/admin/products/:id/order` | `reorderProduct(id, sortOrder)` |
| GET | `/api/admin/products/:id/cards` | `getAdminCards(productId)` |
| POST | `/api/admin/cards` | `addCards(formData)` _(multipart)_ |
| DELETE | `/api/admin/cards/:id` | `deleteCard(cardId)` |
| DELETE | `/api/admin/cards` | `deleteCards(cardIds[])` |
| PATCH | `/api/admin/products/:id/cards/api` | `saveCardsApiConfig(...)` |
| PATCH | `/api/admin/products/:id/cards/api/enabled` | `setCardsApiEnabled(...)` |
| POST | `/api/admin/products/:id/cards/pull` | `pullCardFromApi(productId)` |
| GET | `/api/admin/orders` | `getAdminOrders(params)` |
| GET | `/api/admin/orders/:id` | `getAdminOrder(id)` |
| POST | `/api/admin/orders/:id/mark-paid` | `markOrderPaid(orderId)` |
| POST | `/api/admin/orders/:id/mark-delivered` | `markOrderDelivered(orderId)` |
| POST | `/api/admin/orders/:id/cancel` | `cancelOrder(orderId)` |
| PATCH | `/api/admin/orders/:id/email` | `updateOrderEmail(orderId, email)` |
| DELETE | `/api/admin/orders/:id` | `deleteOrder(orderId)` |
| DELETE | `/api/admin/orders` | `deleteOrders(orderIds[])` |
| GET | `/api/admin/orders/:id/refund-status` | `verifyOrderRefundStatus(orderId)` |
| POST | `/api/admin/orders/:id/mark-refunded` | `markOrderRefunded(orderId)` |
| POST | `/api/admin/orders/:id/proxy-refund` | `proxyRefund(orderId)` |
| GET | `/api/admin/refunds` | `getAdminRefunds()` |
| GET | `/api/admin/refunds/pending-count` | `getPendingRefundRequestCount()` |
| POST | `/api/admin/refunds/:id/approve` | `adminApproveRefund(id, note?)` |
| POST | `/api/admin/refunds/:id/reject` | `adminRejectRefund(id, note?)` |
| GET | `/api/admin/users` | `getAdminUsers(params)` |
| PATCH | `/api/admin/users/:id/points` | `saveUserPoints(userId, points)` |
| PATCH | `/api/admin/users/:id/block` | `toggleBlock(userId, isBlocked)` |
| GET | `/api/admin/categories` | `getAdminCategories()` |
| POST | `/api/admin/categories` | `saveCategory(formData)` |
| DELETE | `/api/admin/categories/:id` | `deleteCategory(id)` |
| GET | `/api/admin/messages` | `getAdminMessages()` → `AdminMessagesPayload` |
| POST | `/api/admin/messages` | `sendAdminMessage(params)` |
| DELETE | `/api/admin/messages/:id` | `deleteAdminMessage(id)` |
| POST | `/api/admin/messages/clear` | `clearAdminMessages()` |
| GET | `/api/admin/user-messages/unread-count` | `getUnreadUserMessageCount()` |
| POST | `/api/admin/user-messages/:id/read` | `markUserMessageRead(id)` |
| GET | `/api/admin/notifications` | `getAdminNotificationSettings()` |
| POST | `/api/admin/notifications` | `saveNotificationSettings(formData)` |
| POST | `/api/admin/notifications/test/telegram` | `testNotification()` |
| POST | `/api/admin/notifications/test/bark` | `testBarkNotification()` |
| POST | `/api/admin/notifications/test/email` | `testEmailNotification(to)` |
| GET | `/api/admin/reviews` | `getAdminReviews()` |
| DELETE | `/api/admin/reviews/:id` | `deleteReview(reviewId)` |
| GET | `/api/admin/settings` | _(via dashboard)_ |
| PATCH | `/api/admin/settings/:key` | `saveSetting(key, value)` |
| GET | `/api/admin/announcement` | `getAnnouncementConfig()` |
| POST | `/api/admin/announcement` | `saveAnnouncement(config)` |
| GET | `/api/admin/collect` | `getAdminCollect()` → `AdminCollectPayload` |
| GET | `/api/admin/data` | `getAdminData()` |
| POST | `/api/admin/data/repair` | `repairDataAction()` |
| POST | `/api/admin/data/import` | `importData(formData)` _(multipart)_ |
| POST | `/api/admin/registry/dismiss` | `dismissRegistryPrompt()` |
| POST | `/api/admin/registry/join` | `joinRegistry(origin)` |
| POST | `/api/admin/registry/leave` | `leaveRegistry()` |
| GET | `/api/admin/registry/status` | `getRegistryStatus()` |

---

## Appendix C: Critical Business Rules Reference

Chi tiết đầy đủ tại: `/docs/migration/legacy_logic/domain-rules.md`

| Rule | Constants | Notes |
|------|-----------|-------|
| Stock Reservation TTL | `300,000 ms` (5 min) | Auto-cancel + release sau 5 phút |
| Infinite Stock | `999,999` | Shared products với ≥ 1 available card |
| Points Conversion | 1 point = 1 currency unit | `min(userPoints, ceil(totalAmount))` |
| Visibility Threshold | -1 = public, 0 = logged-in, 1-3 = trust level | SQL WHERE filter |
| Admin Check | `ADMIN_USERS` env var | Comma-separated, case-insensitive |
| Check-in Timezone | UTC start of day | `Date.UTC(year, month, date)` |
| Broadcast Limit | 10 | Max recent broadcasts in unread count |
| Cleanup Throttle | 600,000 ms (10 min) | Min interval between card cleanup runs |
