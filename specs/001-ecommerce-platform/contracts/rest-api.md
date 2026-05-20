# REST API Contracts: E-Commerce Backend Platform

Base URL: `/api/v1`

## Auth & Users (Admin)

- `POST /auth/login` -> authenticate admin user
- `POST /auth/logout` -> invalidate session
- `GET /users` -> list users
- `POST /users` -> create user
- `GET /users/{id}` -> user detail
- `PATCH /users/{id}` -> update user
- `POST /users/{id}/lock` -> lock user
- `POST /users/{id}/unlock` -> unlock user

## Catalog (Admin)

- `GET /catalog/products` -> list products (filters, pagination)
- `POST /catalog/products` -> create product
- `GET /catalog/products/{id}` -> product detail
- `PATCH /catalog/products/{id}` -> update product
- `DELETE /catalog/products/{id}` -> delete or archive product

- `GET /catalog/categories` -> list categories (tree)
- `POST /catalog/categories` -> create category
- `PATCH /catalog/categories/{id}` -> update category
- `DELETE /catalog/categories/{id}` -> delete category

- `GET /catalog/tags` -> list tags
- `POST /catalog/tags` -> create tag
- `DELETE /catalog/tags/{id}` -> delete tag

## Media Library (Admin)

- `POST /media` -> upload media (JPG/PNG/WebP, <= 5MB)
- `GET /media` -> list media assets
- `DELETE /media/{id}` -> delete media asset

## Content Management (Admin)

- `GET /content/articles` -> list articles
- `POST /content/articles` -> create article (draft)
- `PATCH /content/articles/{id}` -> update article
- `POST /content/articles/{id}/schedule` -> schedule publish
- `POST /content/articles/{id}/publish` -> publish now
- `DELETE /content/articles/{id}` -> delete article

- `GET /content/pages` -> list static pages
- `POST /content/pages` -> create static page
- `PATCH /content/pages/{id}` -> update static page
- `DELETE /content/pages/{id}` -> delete static page

## Homepage Configuration (Admin)

- `GET /homepage/blocks` -> list homepage blocks
- `POST /homepage/blocks` -> create block
- `PATCH /homepage/blocks/{id}` -> update block
- `DELETE /homepage/blocks/{id}` -> delete block

## Support Channels (Admin)

- `GET /support/channels` -> list channels
- `PATCH /support/channels/{id}` -> enable/disable or update

## Public Storefront APIs

- `GET /public/search` -> keyword search with pagination/sort
- `GET /public/categories` -> category tree
- `GET /public/categories/{id}/products` -> products by category
- `GET /public/products/{id}` -> product detail
- `GET /public/homepage` -> homepage composition blocks
- `GET /public/content/articles` -> list articles
- `GET /public/content/articles/{id}` -> article detail
- `GET /public/content/pages/{slug}` -> static page detail
- `GET /public/footer` -> footer configuration
- `GET /public/support` -> support channels

## Cart & Order Requests (Public)

- `POST /cart` -> create cart (session id)
- `GET /cart/{session_id}` -> view cart
- `POST /cart/{session_id}/items` -> add item
- `PATCH /cart/{session_id}/items/{item_id}` -> update quantity
- `DELETE /cart/{session_id}/items/{item_id}` -> remove item
- `POST /order-requests` -> submit order request from cart

## Lead/Contact Forms (Public)

- `POST /leads` -> submit contact or consultation form
- `GET /leads/{id}` -> detail (admin or internal use)

## Response Envelope (Guideline)

- Success: `{ "data": ..., "meta": { "pagination": ... } }`
- Error: `{ "error": { "code": "VALIDATION_ERROR", "message": "...", "details": [...] } }`
