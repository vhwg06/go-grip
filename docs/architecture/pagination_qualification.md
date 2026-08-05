# Shared Pagination Library Qualification Manifest

This document records the consumer inventory, semantic reconciliation, and qualification decision for pagination across the repository.

---

## 1. Consumer Inventory Matrix

| Consumer Component | File Paths | Current Type | Indexing Base | Default Limit | Max Limit | Offset Formula | Invalid Input Policy |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Media** | `usecase/media`, `repo/persistent/media`, `controller/restapi/v1/media` | `entity.Pagination` | 0-based | 20 | 100 | `max(0, Offset)` | Clamp bounds via `Normalize()` |
| **Orders** | `usecase/orders`, `repo/persistent/order`, `controller/restapi/v1/orders` | `entity.Pagination` | 0-based | 20 | 100 | `max(0, Offset)` | Clamp bounds via `Normalize()` |
| **Wishlist** | `usecase/wishlist`, `repo/persistent/wishlist`, `controller/restapi/v1/wishlist` | `entity.Pagination` | 0-based | 20 | 100 | `max(0, Offset)` | Clamp bounds via `Normalize()` |
| **Notification** | `usecase/notification`, `repo/persistent/notification`, `controller/restapi/v1/notification` | `entity.Pagination` | 0-based | 20 | 100 | `max(0, Offset)` | Clamp bounds via `Normalize()` |
| **Admin** | `usecase/admin`, `repo/persistent/admin`, `controller/restapi/v1/admin` | `entity.Pagination` | 0-based | 20 | 100 | `max(0, Offset)` | Clamp bounds via `Normalize()` |
| **Content** | `usecase/content`, `repo/persistent/content` | `ArticleFilter.Pagination` | 0-based | 20 | 100 | `max(0, Offset)` | Clamp bounds via `Normalize()` |
| **Catalog** | `usecase/catalog`, `repo/persistent/catalog`, `entity/product` | `ProductFilter.Pagination` | 0-based | 20 | 100 | `max(0, Offset)` | Clamp bounds via `Normalize()` |
| **User** | `usecase/user`, `controller/restapi/v1/user` | `limit, offset int` | 0-based | 20 | 100 | `max(0, Offset)` | Raw limit/offset parameters |
| **Task** | `usecase/task`, `controller/restapi/v1/task` | `limit, offset int` | 0-based | 20 | 100 | `max(0, Offset)` | Raw limit/offset parameters |

---

## 2. Semantic Reconciliation & Invariants

1. **Page Indexing**: All list endpoints across the repository use 0-based offset pagination with limit size (`OFFSET ... LIMIT ...`).
2. **Bounds**:
   - Default Limit: `20`
   - Maximum Limit: `100`
   - Minimum Offset: `0`
3. **Invalid Input Behavior**:
   - Limit `<= 0` → default to `20`
   - Limit `> 100` → clamp to `100`
   - Offset `< 0` → clamp to `0`
4. **Zero-Value Behavior**:
   - `Pagination{}` zero value represents un-normalized input (`Limit=0, Offset=0`).
   - Calling `Normalize()` on zero-value `Pagination{}` safely returns `Pagination{Limit: 20, Offset: 0}`.

---

## 3. Qualification Decision

**Outcome 1 Selected**: Genuinely Shared Library Primitive (`internal/shared/pagination`).

### Rationale
Every consumer across domain use cases, repositories, and delivery handlers shares identical offset/limit semantics, normalization boundaries (20 default, 100 max), and response page metadata representation. Placing this in `internal/shared/pagination` provides a clean, domain-agnostic dependency leaf without creating a vague catch-all `common` package.

---

## 4. Target Library API Surface

```go
package pagination

// Request metadata for offset/limit pagination.
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// New creates a normalized Pagination with bounded limits.
func New(limit, offset int) Pagination

// Normalize applies bounded defaults (Default: 20, Max: 100, Min Offset: 0).
func (p Pagination) Normalize() Pagination

// Page describes response pagination metadata.
type Page struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// NewPage constructs response pagination metadata.
func NewPage(limit, offset, total int) Page
```
