package gobacktest

// Direction defines which direction a signal indicates
type Direction string

// Signal declares a basic signal event
type Signal struct {
	Event
	id        int       // order Id need to change
	price     float64   // if Price < 0, market order wanted
	qty       float64   // amount of buy , sld or ext
	direction Direction // long, short, exit close, or hold
}

// Direction returns the Direction of a Signal
func (s Signal) Direction() Direction {
	return s.direction
}

// SetDirection sets the Directions field of a Signal
func (s *Signal) SetDirection(dir Direction) {
	s.direction = dir
}

// Price returns the price of a Signal
func (s Signal) Price() float64 {
	return s.price
}

// SetPrice sets the price field of a Signal
func (s *Signal) SetPrice(price float64) {
	s.price = price
}

// Amount returns the price of a Signal
func (s Signal) Qty() float64 {
	return s.qty
}

// SetAmount sets the price field of a Signal
func (s *Signal) SetQty(qty float64) {
	s.qty = qty
}

// ID returns the ID of a Signal
func (s Signal) ID() int {
	return s.id
}

// SetID sets the ID field of a Signal
func (s *Signal) SetID(id int) {
	s.id = id
}
