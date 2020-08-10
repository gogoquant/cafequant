package api

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// FutureBacktest define api not support backtest
type FutureBacktest struct {
	FutureExchange
}

// GetAccount ...
func (e *FutureBacktest) GetAccount() interface{} {
	e.logger.Log(constant.ERROR, "GetAccount", 0, 0, "GetAccount not support")
	return nil
}

// GetDepth ...
func (e *FutureBacktest) GetDepth(size int) interface{} {
	e.logger.Log(constant.ERROR, "GetDepth", 0, 0, "GetDepth not support")
	return nil
}

// Buy ...
func (e *FutureBacktest) Buy(price, amount string, msg ...interface{}) interface{} {
	e.logger.Log(constant.ERROR, "Buy", 0, 0, "Buy not support")
	return nil
}

// Sell ...
func (e *FutureBacktest) Sell(price, amount string, msg ...interface{}) interface{} {
	e.logger.Log(constant.ERROR, "Sell", 0, 0, "Sell not support")
	return nil
}

// GetOrder ...
func (e *FutureBacktest) GetOrder(id string) interface{} {
	e.logger.Log(constant.ERROR, "GetOrder", 0, 0, "GetOrder not support")
	return nil
}

// GetOrders ...
func (e *FutureBacktest) GetOrders() interface{} {
	e.logger.Log(constant.ERROR, "GetOrders", 0, 0, "not support")
	return nil
}

// GetTrades ...
func (e *FutureBacktest) GetTrades(params ...interface{}) interface{} {
	e.logger.Log(constant.ERROR, "GetTrades", 0, 0, "GetTrades not support")
	return nil
}

// CancelOrder ...
func (e *FutureBacktest) CancelOrder(orderID string) interface{} {
	e.logger.Log(constant.ERROR, "CancelOrder", 0, 0, "CancelOrder not support")
	return nil
}

// GetTicker ...
func (e *FutureBacktest) GetTicker() interface{} {
	e.logger.Log(constant.ERROR, "GetTicker", 0, 0, "GetTicker not support")
	return nil
}

// GetPosition ...
func (e *FutureBacktest) GetPosition() interface{} {
	e.logger.Log(constant.ERROR, "GetPosition", 0, 0, "GetPosition not support")
	return nil
}
