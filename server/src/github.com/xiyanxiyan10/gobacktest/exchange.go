package gobacktest

import "github.com/deckarep/golang-set"

// ExecutionHandler is the basic interface for executing orders
type ExchangeHandler interface {
	Booker

	OnData(DataEvent) (*Fill, error)
	OnOrder(OrderEvent, DataHandler) (*Fill, error)
	OnFill(*Fill) error
}

// Booker defines methods for handling the order book of the portfolio
type Booker interface {
	AddOrder(OrderEvent) error
	Orders() ([]OrderEvent, bool)
	CommitOrder(id int) (*Fill, error)
	OrdersBySymbol(symbol string) ([]OrderEvent, bool)
	CancelOrder(id int) error
}

// Exchange is a basic execution handler implementation
type Exchange struct {
	orderManager OrderBook
}

// NewExchange creates a default exchange with sensible defaults ready for use.
func NewExchange() *Exchange {
	return &Exchange{
		orderManager: NewOrderBook(),
	}
}

// OnData executes any open order on new data
func (e *Exchange) OnData(data DataEvent) (*Fill, error) {
	return nil, nil
}

// OnOrder executes an order event
func (e *Exchange) OnOrder(order OrderEvent, data DataHandler) (*Fill, error) {
	return nil, nil
}

// OnOrder executes an order event
func (e *Exchange) OnFill(fill *Fill) error {
	freeHandler := fill.FeeHandler()
	if freeHandler == nil {
		return nil
	}
	cost := freeHandler.Fee(fill.Exchange(), fill.FQty(), fill.Price(), fill.OrderType(), fill.Symbol(), fill.Direction())
	fill.SetCost(cost)
	return nil
}

// OrderBook returns the order book of the portfolio
func (p Exchange) Orders() ([]OrderEvent, bool) {
	return p.orderManager.Orders()
}

// OrdersBySymbol returns the order of a specific symbol from the order book.
func (p Exchange) OrdersBySymbol(symbol string) ([]OrderEvent, bool) {
	return p.orderManager.OrdersBySymbol(symbol)
}

// CancelOrder ...
func (p *Exchange) CancelOrder(id int) error {
	p.orderManager.CancelOrder(id)
	return nil
}

// CommitOrder ...
func (p *Exchange) CommitOrder(id int) (*Fill, error) {
	return p.orderManager.CommitOrder(id)
}

// AddOrder
func (p *Exchange) AddOrder(o OrderEvent) error {
	return p.orderManager.Add(o)
}

//  SetSubscribes
func (p *Exchange) SetSubscribe(symbol string) error {
	return p.orderManager.SetSubscribe(symbol)
}

// Subscribes
func (p *Exchange) Subscribes() mapset.Set {
	return p.orderManager.Subscribes()
}
