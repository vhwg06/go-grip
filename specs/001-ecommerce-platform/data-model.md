# Data Model: E-Commerce Backend Platform

## Entities and Attributes

### User
- Fields: id, email, display_name, role_id, status (active/locked), created_at, updated_at
- Validation: email must be unique; status in {active, locked}

### Role
- Fields: id, name (Administrator, Editor, Author, Contributor, Subscriber)
- Rules: only Administrators may change roles

### Product
- Fields: id, title, sku, description, price, compare_price, status (active/draft/archived), brand, attributes (json), created_at, updated_at
- Validation: sku unique; price >= 0

### Category
- Fields: id, name, parent_id (nullable), position, is_active
- Relationship: parent-child tree

### Tag
- Fields: id, name
- Relationship: many-to-many with Product

### ProductCategory
- Fields: product_id, category_id

### ProductTag
- Fields: product_id, tag_id

### MediaAsset
- Fields: id, file_name, mime_type, size_bytes, url, alt_text, owner_type, owner_id, created_at
- Validation: mime_type in {image/jpeg, image/png, image/webp}; size_bytes <= 5MB

### Cart
- Fields: id, session_id, status (active/abandoned/converted), created_at, updated_at
- Relationship: one-to-many with CartItem

### CartItem
- Fields: id, cart_id, product_id, quantity, unit_price, product_snapshot (json)
- Validation: quantity >= 1

### OrderRequest
- Fields: id, cart_id, customer_name, customer_phone, customer_email, address, note, status (new/in_progress/done), created_at
- Relationship: references Cart snapshot at submission time

### ContentArticle
- Fields: id, title, slug, body, status (draft/scheduled/published), scheduled_at, published_at, author_id, created_at, updated_at
- Validation: title required; body required for publish

### StaticPage
- Fields: id, title, slug, body, template_key, status (draft/published), updated_at

### HomepageBlock
- Fields: id, block_type (banner, featured_category, category_block, footer, floating), config (json), position, is_active

### LeadSubmission
- Fields: id, source (contact/product), customer_name, customer_phone, customer_email, message, status (new/in_progress/done), created_at

### SupportChannel
- Fields: id, channel_type (live_chat/call/zalo/facebook), label, link, is_enabled

### SeoMetadata
- Fields: id, owner_type, owner_id, meta_title, meta_description, slug, alt_text

## Relationships

- User (1) -> (N) ContentArticle (author)
- Category (1) -> (N) Category (children)
- Product (N) <-> (N) Category via ProductCategory
- Product (N) <-> (N) Tag via ProductTag
- Product (1) -> (N) MediaAsset (owner_type=product)
- ContentArticle (1) -> (N) MediaAsset (owner_type=content)
- Cart (1) -> (N) CartItem
- Cart (1) -> (1) OrderRequest
- HomepageBlock (1) -> (0..N) MediaAsset (optional)
- SeoMetadata (1) -> (1) Product, ContentArticle, StaticPage (owner_type)

## State Transitions

- ContentArticle: draft -> scheduled -> published (scheduled can publish by time; catch-up after downtime)
- OrderRequest: new -> in_progress -> done
- Cart: active -> converted (order submitted) or abandoned
