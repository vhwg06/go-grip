# Security Checklist: Grip Store Backend API

- [x] Card keys are masked or omitted for non-delivered orders in buyer and public responses.
- [x] Card keys are only returned when order state is delivered and requester is owner or admin.
- [x] All `/v1/admin/*` routes are protected by authentication middleware and admin-username authorization middleware.
- [x] Blocked users are rejected for mutating buyer actions (checkout, wishlist, reviews, notifications updates).
- [x] Payment callback processing validates provider signature before applying paid/delivered transitions.
- [x] Payment callback handling is idempotent for replayed success notifications.
- [x] Refresh token/session rotation invalidates previously used refresh credentials.
- [x] Sensitive order/payment errors are mapped to safe API error responses without leaking secrets.
