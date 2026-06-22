# Implementation Plan: Grip Store Backend API

## Summary

Build and maintain the Grip Store backend as a Go modular monolith with REST-only transport, PostgreSQL persistence, and published contracts aligned to the current no-cards/no-points/no-checkin scope.

## Current Scope Lock

- Product editorial/media remains on `/v1/admin/products` and product form routes.
- Orders, payments, refunds, profile, settings, review, wishlist, notification, and messaging flows remain active.
- Cards, points mutation, checkout points fields, and check-in routes are removed from the active domain and must stay absent from contracts, Swagger artifacts, and generated mocks.

## Delivery Constraints

- Keep Clean Architecture boundaries under `internal/entity`, `internal/usecase`, `internal/repo`, and `internal/controller/restapi`.
- Keep generated mocks and Swagger artifacts in sync with source interfaces and handler annotations.
- Prefer focused verification on checkout/orders/profile/admin slices when changing closeout surfaces.
