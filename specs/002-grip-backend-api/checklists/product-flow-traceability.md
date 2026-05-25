# Product-Flow Traceability Matrix

This matrix keeps `@product-flow` Playwright as the executable source of truth while tracing each scenario to Figma node(s) and the feature spec intent.

## Figma node anchors

- Homepage: `27:1404`
- Product Listing: `58:861`
- Product Detail: `62:2672`
- Cart: `114:3466`
- Checkout: `117:4153`

## Product-flow API

- `PF-API-001` list active products
  - Spec intent: catalog listing returns active products with paging.
  - Figma node: `58:861` (listing data source).
- `PF-API-002` category filter
  - Spec intent: listing can be filtered by category.
  - Figma node: `58:861` (filter panel/listing grid).
- `PF-API-003` keyword search
  - Spec intent: search returns matching catalog results.
  - Figma node: `58:861` (search/result area).
- `PF-API-004` detail returns specs
  - Spec intent: product detail exposes spec entries.
  - Figma node: `62:2672` (spec section).
- `PF-API-005` inactive product hidden
  - Spec intent: non-visible product cannot be viewed by buyer route.
  - Figma node: `62:2672` (detail entry gate).
- `PF-API-006` admin create persists specs
  - Spec intent: admin create writes persistent product details.
  - Figma node: `62:2672` (detail fields depend on saved data).
- `PF-API-007` admin update replaces specs transactionally
  - Spec intent: update is consistent and stale details are removed.
  - Figma node: `62:2672`.
- `PF-API-008` admin delete cascades details
  - Spec intent: deleting product removes dependent detail rows.
  - Figma node: `58:861` and `62:2672` (no stale list/detail rendering).

## Product-flow UI catalog

- `PF-CATALOG-001` list cards visible
  - Spec intent: listing renders discoverable product cards.
  - Figma node: `58:861`.
- `PF-CATALOG-002` category query navigation
  - Spec intent: category state filters listing view.
  - Figma node: `58:861`.
- `PF-CATALOG-003` search result visibility
  - Spec intent: query-to-result flow.
  - Figma node: `58:861`.
- `PF-CATALOG-004` sort behavior
  - Spec intent: user-controlled ordering.
  - Figma node: `58:861`.
- `PF-CATALOG-005` empty result state
  - Spec intent: no-match UX contract.
  - Figma node: `58:861`.
- `PF-CATALOG-006` card-to-detail navigation
  - Spec intent: listing CTA routes to detail.
  - Figma node: `58:861` -> `62:2672`.

## Product-flow UI detail

- `PF-DETAIL-001` core info visible
  - Spec intent: detail exposes title + price baseline.
  - Figma node: `62:2672`.
- `PF-DETAIL-002` specs table visible
  - Spec intent: backend details are rendered in spec table.
  - Figma node: `62:2672`.
- `PF-DETAIL-003` add to cart increments badge
  - Spec intent: detail purchase CTA mutates cart.
  - Figma node: `62:2672` -> `114:3466`.
- `PF-DETAIL-004` quantity respected
  - Spec intent: quantity control is honored at add-to-cart boundary.
  - Figma node: `62:2672` -> `114:3466`.
- `PF-DETAIL-005` non-existent/inactive cannot purchase
  - Spec intent: unavailable product must not expose buy CTA.
  - Figma node: `62:2672`.

## Product-flow UI homepage

- `PF-HOME-001` hero/category/featured blocks
  - Spec intent: guest discovery landing baseline.
  - Figma node: `27:1404`.
- `PF-HOME-002` category icon to listing
  - Spec intent: homepage category navigation.
  - Figma node: `27:1404` -> `58:861`.
- `PF-HOME-003` featured CTA to detail
  - Spec intent: featured item drill-down.
  - Figma node: `27:1404` -> `62:2672`.
- `PF-HOME-004` homepage CTA does not mutate cart
  - Spec intent: discovery CTA is navigation only.
  - Figma node: `27:1404` and `114:3466`.
- `PF-HOME-005` no follow action on featured card
  - Spec intent: avoid non-specified social action in featured context.
  - Figma node: `27:1404`.
- `PF-HOME-006` no add-to-cart on featured card
  - Spec intent: homepage card remains browse-first.
  - Figma node: `27:1404`.
