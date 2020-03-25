package backtest

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
	LimitOrder
)

// Order declares a basic order event.
type Order struct {
	Event
	id           int
	orderType    OrderType // market or limit
	status       OrderStatus
	direction    Direction // buy or sell
	assetType    string
	qty          float64 // quantity of the order
	price        float64 // if Price < 0, market order wanted
	qtyFilled    float64
	avgFillPrice float64
	limitPrice   float64 // limit for the order
	stopPrice    float64
}

// Price returns the price of a Signal
func (s Order) OrderType() OrderType {
	return s.orderType
}

// SetPrice sets the price field of a Signal
func (s *Order) SetOrderType(t OrderType) {
	s.orderType = t
}

// Price returns the price of a Signal
func (s Order) Price() float64 {
	return s.price
}

// SetPrice sets the price field of a Signal
func (s *Order) SetPrice(price float64) {
	s.price = price
}

// ID returns the id of the Order.
func (o Order) ID() int {
	return o.id
}

// SetID of the Order.
func (o *Order) SetID(id int) {
	o.id = id
}

// Direction returns the Direction of an Order
func (o Order) Direction() Direction {
	return o.direction
}

// SetDirection sets the Directions field of an Order
func (o *Order) SetDirection(dir Direction) {
	o.direction = dir
}

// Qty returns the Qty field of an Order
func (o Order) Qty() float64 {
	return o.qty
}

// SetQty sets the Qty field of an Order
func (o *Order) SetQty(i float64) {
	o.qty = i
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
