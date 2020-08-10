package api

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
)

// BackTestCover define api not support backtest
type BackTestCover struct {
	logger model.Logger
}

// GetAccount ...
func (e *BackTestCover) GetAccount() interface{} {
	e.logger.Log(constant.ERROR, "GetAccount", 0, 0, "GetAccount not support")
	return nil
}

// GetDepth ...
func (e *BackTestCover) GetDepth(size int) interface{} {
	e.logger.Log(constant.ERROR, "GetDepth", 0, 0, "GetDepth not support")
	return nil
}

// Buy ...
func (e *BackTestCover) Buy(price, amount string, msg ...interface{}) interface{} {
	e.logger.Log(constant.ERROR, "Buy", 0, 0, "Buy not support")
	return nil
}

// Sell ...
func (e *BackTestCover) Sell(price, amount string, msg ...interface{}) interface{} {
	e.logger.Log(constant.ERROR, "Sell", 0, 0, "Sell not support")
	return nil
}

// GetOrder ...
func (e *BackTestCover) GetOrder(id string) interface{} {
	e.logger.Log(constant.ERROR, "GetOrder", 0, 0, "GetOrder not support")
	return nil
}

// GetOrders ...
func (e *BackTestCover) GetOrders() interface{} {
	e.logger.Log(constant.ERROR, "GetOrders", 0, 0, "not support")
	return nil
}

// GetTrades ...
func (e *BackTestCover) GetTrades(params ...interface{}) interface{} {
	e.logger.Log(constant.ERROR, "GetTrades", 0, 0, "GetTrades not support")
	return nil
}

// CancelOrder ...
func (e *BackTestCover) CancelOrder(orderID string) interface{} {
	e.logger.Log(constant.ERROR, "CancelOrder", 0, 0, "CancelOrder not support")
	return nil
}

// GetTicker ...
func (e *BackTestCover) GetTicker() interface{} {
	e.logger.Log(constant.ERROR, "GetTicker", 0, 0, "GetTicker not support")
	return nil
}

// GetPosition ...
func (e *BackTestCover) GetPosition() interface{} {
	e.logger.Log(constant.ERROR, "GetPosition", 0, 0, "GetPosition not support")
	return nil
}

// BaseExchange ...
type BaseExchange struct {
	BaseExchangeCachePool // cache for exchange
	ID                    int
	IOMode                int                // io mode for exchange
	contractType          string             // contractType
	direction             string             // trade type
	stockType             string             // stockType
	lever                 float64            // lever
	recordsPeriodMap      map[string]int64   // recordsPeriod support
	minAmountMap          map[string]float64 // minAmount of trade
}

// SetID set ID
func (e *BaseExchange) SetID(mode int) {
	e.ID = mode
}

// GetID get ID
func (e *BaseExchange) GetID() int {
	return e.ID
}

// SetIO set IO mode
func (e *BaseExchange) SetIO(mode int) {
	e.IOMode = mode
}

// GetIO get IO mode
func (e *BaseExchange) GetIO() int {
	return e.IOMode
}

// SetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) SetContractType(contractType string) {
	e.contractType = contractType
}

// GetContractType set the limit calls amount per second of this exchange
func (e *BaseExchange) GetContractType() string {
	return e.contractType
}

// SetDirection set the limit calls amount per second of this exchange
func (e *BaseExchange) SetDirection(direction string) {
	e.direction = direction
}

// GetDirection set the limit calls amount per second of this exchange
func (e *BaseExchange) GetDirection() string {
	return e.direction
}

// SetMarginLevel set the limit calls amount per second of this exchange
func (e *BaseExchange) SetMarginLevel(lever float64) {
	e.lever = lever
}

// GetMarginLevel set the limit calls amount per second of this exchange
func (e *BaseExchange) GetMarginLevel() float64 {
	return e.lever
}

// GetStockType set the limit calls amount per second of this exchange
func (e *BaseExchange) GetStockType() string {
	return e.stockType
}

// SetStockType set the limit calls amount per second of this exchange
func (e *BaseExchange) SetStockType(stockType string) {
	e.stockType = stockType
}

// SetMinAmountMap ...
func (e *BaseExchange) SetMinAmountMap(m map[string]float64) {
	e.minAmountMap = m
}

// GetMinAmountMap ...
func (e *BaseExchange) GetMinAmountMap() map[string]float64 {
	return e.minAmountMap
}

// SetRecordsPeriodMap ...
func (e *BaseExchange) SetRecordsPeriodMap(m map[string]int64) {
	e.recordsPeriodMap = m
}

// GetRecordsPeriodMap ...
func (e *BaseExchange) GetRecordsPeriodMap() map[string]int64 {
	return e.recordsPeriodMap
}
