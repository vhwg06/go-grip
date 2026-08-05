# Order Module Migration Manifest

This document records the inventory of all symbols, contracts, adapters, delivery endpoints, consumers, and migration risks for the `Order` business module.

---

## 1. Owned Symbols
- **Entities**: `Order`, `Payment`, `RefundRequest`, `RefundStatus` (`pending`, `approved`, `rejected`), `OrderStatus` (`pending`, `paid`, `delivered`, `cancelled`, `failed`, `refund_pending`, `refunded`), `Amount`
- **Methods**: `(Order) CanRequestRefund() bool`, `(Order) IsTerminal() bool`

---

## 2. Owned Ports
- **Persistence Ports**:
  - `OrderRepo`: `ListOrdersByOwner`, `GetOrderByID`, `SubmitRefundRequest`, `CancelPendingOrder`
  - `CheckoutRepo`: `CreateOrderWithReservation`, `AttachPayment`, `UpdateOrderStatus`, `ReleaseReservation`

---

## 3. Application Services
- **Interfaces**:
  - `OrdersUseCase`: `List`, `Get`, `RequestRefund`
  - `CheckoutUseCase`: `Preview`, `CreateOrder`, `PaymentParams`, `PaymentNotify`, `PaymentStatus`, `Cancel`
- **Implementation Structs**: `ordersUseCase`, `checkoutUseCase`
- **Constructors**: `NewOrdersUseCase`, `NewCheckoutUseCase`

---

## 4. Infrastructure Adapters
- `internal/repo/persistent/order_postgres.go`

---

## 5. Delivery Endpoints & Mappers
- `internal/controller/restapi/v1/orders/`
- `internal/controller/restapi/v1/checkout/`

---

## 6. App Wiring
- `internal/app/app.go` (`useCases.orders`, `useCases.checkout`)

---

## 7. Unit & Integration Tests
- `internal/usecase/orders/`
- `internal/usecase/checkout/`

---

## 8. Cross-Module Dependencies & Risks
- Checkout interacts with `webapi.PaymentVerifier`.
- Risk: Low.
