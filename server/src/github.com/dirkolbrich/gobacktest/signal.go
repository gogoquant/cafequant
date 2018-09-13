package gobacktest

// Direction defines which direction a signal indicates
type Direction int

// different types of order directions
const (
	// Buy
	BOT Direction = iota // 0
	// Sell
	SLD
	// Hold
	HLD
	// Exit
	EXT
)

// Signal declares a basic signal event
type Signal struct {
	Event
	orderType OrderType // orderType order type
	qty   int64         // qte of the trader
	price float64       // price of the Signal
	amount float64      // amount of the Signal
	direction Direction // long, short, exit or hold
}

// Direction returns the Direction of a Signal
func (s Signal) Direction() Direction {
	return s.direction
}

// SetDirection sets the Directions field of a Signal
func (s *Signal) SetDirection(dir Direction) {
	s.direction = dir
}

// Qty returns the Qty field of a Signal
func (s *Signal) Qty() int64 {
	return s.qty
}

// SetQty sets the Qty field of a Signal
func (s *Signal) SetQty(i int64) {
	s.qty = i
}

// Price returns the Price field of a Signal
func (s *Signal) Price() float64 {
	return s.price
}

// SetPrice sets the Price field of a Signal
func (s *Signal) SetPrice(i float64) {
	s.price = i
}

// Amount returns the Amount field of a Signal
func (s *Signal) Amount() float64 {
	return s.amount
}

// SetAmount sets the Amount field of a Signal
func (s *Signal) SetAmount(i float64) {
	s.amount = i
}

// OrderType returns the OrderType field of a Signal
func (s *Signal) OrderType() OrderType {
	return s.orderType
}

// SetOrderType sets the OrderType field of a Signal
func (s *Signal) SetOrderType(i OrderType) {
	s.orderType = i
}

