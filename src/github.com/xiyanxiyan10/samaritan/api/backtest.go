package api

import (
	"github.com/bitly/go-simplejson"
)

// Backtest backtest struct
type Backtest struct {
}

// NewBacktest create a backtest
func NewBacktest(opt Option) Exchange {
	return &Backtest{}
}

// Log print something to console
func (e *Backtest) Log(msgs ...interface{}) {

}

// GetType get the type of this exchange
func (e *Backtest) GetType() string {
	return ""
}

// GetName get the name of this exchange
func (e *Backtest) GetName() string {
	return ""
}

// SetLimit set the limit calls amount per second of this exchange
func (e *Backtest) SetLimit(times interface{}) float64 {
	return 0.0
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *Backtest) AutoSleep() {

}

// GetMinAmount get the min trade amonut of this exchange
func (e *Backtest) GetMinAmount(stock string) float64 {
	return 0.0
}

// getAuthJSON
func (e *Backtest) getAuthJSON(method string, params ...interface{}) (jsoner *simplejson.Json, err error) {
	data := []byte{}
	return simplejson.NewJson(data)
}

// GetAccount get the account detail of this exchange
func (e *Backtest) GetAccount() interface{} {
	return map[string]float64{}
}

// Trade place an order
func (e *Backtest) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	return map[string]float64{}
}

func (e *Backtest) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	return map[string]float64{}
}

func (e *Backtest) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	return map[string]float64{}
}

// GetOrder get details of an order
func (e *Backtest) GetOrder(stockType, id string) interface{} {
	return map[string]float64{}
}

// GetOrders get all unfilled orders
func (e *Backtest) GetOrders(stockType string) interface{} {
	return map[string]float64{}
}

// GetTrades get all filled orders recently
func (e *Backtest) GetTrades(stockType string) interface{} {
	return map[string]float64{}
}

// CancelOrder cancel an order
func (e *Backtest) CancelOrder(order Order) bool {
	return true
}

// getTicker get market ticker & depth
func (e *Backtest) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	return
}

// GetTicker get market ticker & depth
func (e *Backtest) GetTicker(stockType string, sizes ...interface{}) interface{} {
	return map[string]float64{}
}

// GetRecords get candlestick data
func (e *Backtest) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return map[string]float64{}
}
