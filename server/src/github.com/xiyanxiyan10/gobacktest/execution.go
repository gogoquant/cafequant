package gobacktest

import (
// "fmt"
)

// ExecutionHandler is the basic interface for executing orders
type ExecutionHandler interface {
	OnData(DataEvent) (*Fill, error)
	OnOrder(OrderEvent, DataHandler) (*Fill, error)
	OnFill(*Fill) error
}

// Exchange is a basic execution handler implementation
type Exchange struct {
	Symbol string
	//Commission  CommissionHandler
	ExchangeFee ExchangeFeeHandler
}

// NewExchange creates a default exchange with sensible defaults ready for use.
func NewExchange() *Exchange {
	return &Exchange{
		Symbol: "TEST",
		//Commission:  &FixedCommission{Commission: 0},
		ExchangeFee: &FixedExchangeFee{ExchangeFee: 0},
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
