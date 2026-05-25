# Product-Flow Admin Refinement Notes

Context: Playwright is treated as source of truth for current business checks. Current suite passes, but `playwright/specs/product-flow` does not fully express admin product-flow usecases yet.

## Current coverage status

- `playwright/specs/product-flow/*`
  - Covers buyer-facing product flow (home/catalog/detail) and admin CRUD at API level (`PF-API-006..008`).
  - Does **not** cover admin product flow at UI level end-to-end under `@product-flow`.
- Admin UI checks currently live in `playwright/specs/admin/products.spec.ts` and `playwright/specs/admin/admin-specs.spec.ts`.

## Gaps to refine into product-flow requirements

1. Admin create product UI contract
   - Required business checks not explicit in `@product-flow`:
   - category required/validation
   - slug uniqueness + normalization behavior
   - price format/rounding consistency from UI -> backend
2. Admin update product UI contract
   - persist + reload consistency for title/price/category/visibility/specs/images
   - optimistic/stale write behavior (if any) not defined
3. Admin product list behavior
   - search/filter/sort/pagination semantics not defined as product-flow requirements
4. Admin delete safety
   - confirm dialog semantics + post-delete list/detail consistency not defined
5. Media upload flow contract
   - presigned URL failure path, invalid mime, oversized file behavior not captured in product-flow
6. Role/permission business boundary
   - non-admin access behavior for admin product pages is tested in admin/api sets but not linked to product-flow requirement IDs

## Suggested requirement IDs (for future spec/test generation)

- `PF-ADMIN-UI-001` create product with mandatory fields + validation errors
- `PF-ADMIN-UI-002` create product with specs + image and verify storefront detail rendering
- `PF-ADMIN-UI-003` update product and verify replaced specs + updated media
- `PF-ADMIN-UI-004` toggle visibility and verify catalog visibility boundary
- `PF-ADMIN-UI-005` delete product and verify not reachable on catalog/detail
- `PF-ADMIN-UI-006` admin list query/sort/pagination consistency
- `PF-ADMIN-UI-007` permission guard for non-admin

## Rule compliance note

- No test was modified in this refinement step.
- This file is for requirement refinement only; implementation remains code-first when tests fail.
