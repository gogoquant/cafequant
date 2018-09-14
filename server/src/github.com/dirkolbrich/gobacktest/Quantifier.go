package gobacktest

// Direction defines which direction a signal indicates
type Direction int

// Qty type
type QtyType int

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

const (
	// float64 qte
	FLOAT64_QTY QtyType = iota
	// int64 qte
	INT64_QTY
)

type QuantifierType struct {
	orderType OrderType // orderType order type
	qtyType   QtyType   // qty type
	qty       int64     // qte of the trader as int64
	fqty      float64   // qte of the trader as float64
	price     float64   // price of the Signal
	direction Direction // long, short, exit or hold
}

// Direction returns the Direction of a Signal
func (s QuantifierType) Direction() Direction {
	return s.direction
}

// SetDirection sets the Directions field of a Signal
func (s *QuantifierType) SetDirection(dir Direction) {
	s.direction = dir
}

// Qty returns the Qty field of a Signal
func (s *QuantifierType) Qty() int64 {
	return s.qty
}

// SetQty sets the Qty field of a Signal
func (s *QuantifierType) SetQty(i int64) {
	s.qty = i
	//s.qtyType = INT64_QTY
}

// Price returns the Price field of a Signal
func (s *QuantifierType) Price() float64 {
	return s.price
}

// SetPrice sets the Price field of a Signal
func (s *QuantifierType) SetPrice(i float64) {
	s.price = i
}

// FQty returns the FQty field of a Signal
func (s *QuantifierType) FQty() float64 {
	return s.fqty
}

// SetFQty sets the FQty field of a Signal
func (s *QuantifierType) SetFQty(i float64) {
	//s.qtyType = FLOAT64_QTY
	s.fqty = i
}

// GetQtyType return the Qty type field of a Signal
func (s *QuantifierType) QtyType() QtyType {
	return s.qtyType
}

// GetQtyType set the Qty type field of a Signal
func (s *QuantifierType) SetQtyType(i QtyType) {
	s.qtyType = i
}

// OrderType returns the OrderType field of a Signal
func (s *QuantifierType) OrderType() OrderType {
	return s.orderType
}

// SetOrderType sets the OrderType field of a Signal
func (s *QuantifierType) SetOrderType(i OrderType) {
	s.orderType = i
}

// SetOrderType sets the OrderType field of a Signal
func (s *QuantifierType) SetQuantifier(orderType OrderType, qtyType QtyType, qty int64, fqty float64, direction Direction, price float64) {
	s.qtyType = qtyType
	s.qty = qty
	s.fqty = fqty
	s.orderType = orderType
	s.direction = direction
	s.price = price
}

// SetOrderType sets the OrderType field of a Signal
func (s *QuantifierType) Quantifier() (orderType OrderType, qtyType QtyType, qty int64, fqty float64, direction Direction, price float64) {
	return s.orderType, s.qtyType, s.qty, s.fqty, s.direction, s.price
}
