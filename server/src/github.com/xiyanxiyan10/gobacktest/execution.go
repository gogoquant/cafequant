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
	Symbol      string
	Commission  CommissionHandler
	ExchangeFee ExchangeFeeHandler
}

// NewExchange creates a default exchange with sensible defaults ready for use.
func NewExchange() *Exchange {
	return &Exchange{
		Symbol:      "TEST",
		Commission:  &FixedCommission{Commission: 0},
		ExchangeFee: &FixedExchangeFee{ExchangeFee: 0},
	}
}

// OnData executes any open order on new data
func (e *Exchange) OnData(data DataEvent) (*Fill, error) {
	return nil, nil
}

// OnOrder executes an order event
func (e *Exchange) OnOrder(order OrderEvent, data DataHandler) (*Fill, error) {
	if order.OrderType() == LimitOrder || order.OrderType() == MarketOrder {
		return e.createOrder(order, data)
	}
	return nil, nil
}

// OnOrder executes an order event
func (e *Exchange) createOrder(order OrderEvent, data DataHandler) (*Fill, error) {

	// simple implementation, creates a direct fill from the order
	// based on the last known data price
	f := &Fill{
		Event:    Event{timestamp: order.Time(), symbol: order.Symbol()},
		Exchange: e.Symbol,
		//qty:      order.Qty(),
		//price:    latest.Price(), // last price from data event
	}
	f.SetQuantifier(order.Quantifier())

	fqty := order.FQty()

	commission, err := e.Commission.Calculate(fqty, f.price)
	if err != nil {
		return f, err
	}
	f.commission = commission

	exchangeFee, err := e.ExchangeFee.Fee()
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
