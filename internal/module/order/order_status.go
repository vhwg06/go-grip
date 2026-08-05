package order

// OrderStatus represents lifecycle states of an Order.
type OrderStatus string

const (
	OrderStatusPending       OrderStatus = "pending"
	OrderStatusPaid          OrderStatus = "paid"
	OrderStatusDelivered     OrderStatus = "delivered"
	OrderStatusCancelled     OrderStatus = "cancelled"
	OrderStatusFailed        OrderStatus = "failed"
	OrderStatusRefundPending OrderStatus = "refund_pending"
	OrderStatusRefunded      OrderStatus = "refunded"
)
