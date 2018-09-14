package gobacktest

// OrderStatus defines an order status
type OrderStatus int

// different types of order status
const (
	OrderNone OrderStatus = iota // 0
	OrderNew
	OrderSubmitted
	OrderPartiallyFilled
	OrderFilled
	OrderCanceled
	OrderCancelPending
	OrderInvalid
)

// OrderType defines which type an order is
type OrderType int

// different types of orders
const (
	MarketOrder OrderType = iota // 0
	MarketOnOpenOrder
	MarketOnCloseOrder
	StopMarketOrder
	LimitOrder
	StopLimitOrder
)

// Order declares a basic order event.
type Order struct {
	QuantifierType
	Event
	id           int
	status       OrderStatus
	assetType    string
	qtyFilled    int64
	avgFillPrice float64
	limitPrice   float64 // limit for the order
	stopPrice    float64
}

// ID returns the id of the Order.
func (o Order) ID() int {
	return o.id
}

// SetID of the Order.
func (o *Order) SetID(id int) {
	o.id = id
}

// Status returns the status of an Order
func (o Order) Status() OrderStatus {
	return o.status
}

// Limit returns the limit price of an Order
func (o Order) Limit() float64 {
	return o.limitPrice
}

// Stop returns the stop price of an Order
func (o Order) Stop() float64 {
	return o.stopPrice
}

// Cancel cancels an order
func (o *Order) Cancel() {
	o.status = OrderCancelPending
}

// Update updates an order on a fill event
func (o *Order) Update(fill FillEvent) {
	// not implemented
}
