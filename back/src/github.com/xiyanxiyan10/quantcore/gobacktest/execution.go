package gobacktest

import (
// "fmt"
)

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	OnData(DataEvent) (*Fill, error)
	OnOrder(OrderEvent, DataHandler) (*Fill, error)
}

// Exchange is a basic execution handler implementation
type Exchange struct {
	Symbol     string
	Commission CommissionHandler
}

// NewExchange creates a default exchange with sensible defaults ready for use.
func NewExchange() *Exchange {
	return &Exchange{
		Symbol:     "TEST",
		Commission: &FixedCommission{Commission: 0},
	}
}

// OnData executes any open order on new data
func (e *Exchange) OnData(data DataEvent) (*Fill, error) {
	return nil, nil
}

// OnOrder executes an order event
func (e *Exchange) OnOrder(order OrderEvent, data DataHandler) (*Fill, error) {
	// fetch latest known data event for the symbol
	latest := data.Latest(order.Symbol())

	// simple implementation, creates a direct fill from the order
	// based on the last known data price
	f := &Fill{
		Event:    Event{timestamp: order.Time(), symbol: order.Symbol()},
		Exchange: e.Symbol,
		qty:      order.Qty(),
		price:    latest.Price(), // last price from data event
	}

	f.direction = order.Direction()

	commission, err := e.Commission.Calculate(f.qty, f.price)
	if err != nil {
		return f, err
	}
	f.commission = commission
	// todo
	exchangeFee := 0.1
	if err != nil {
		return f, err
	}
	f.exchangeFee = exchangeFee

	f.cost = e.calculateCost(commission, exchangeFee)

	return f, nil
}

// calculateCost() calculates the total cost for a stock trade
func (e *Exchange) calculateCost(commission, fee float64) float64 {
	return commission + fee
}
