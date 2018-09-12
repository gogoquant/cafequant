package api

import (
	"errors"
	"github.com/bitly/go-simplejson"
	goback "github.com/dirkolbrich/gobacktest"
	"github.com/xiyanxiyan10/samaritan/constant"
	"github.com/xiyanxiyan10/samaritan/model"
)

// Backtest backtest struct
type BtBacktest struct {
	goback.Backtest

	logger           model.Logger
	option           Option
}

// NewBacktest create a backtest
func NewBacktest(opt Option) Exchange {
	return &BtBacktest{}
}

// Log print something to console
func (e *BtBacktest) Log(msgs ...interface{}) {

}

// GetType get the type of this exchange
func (e *BtBacktest) GetType() string {
	return ""
}

// GetName get the name of this exchange
func (e *BtBacktest) GetName() string {
	return ""
}

// SetLimit set the limit calls amount per second of this exchange
func (e *BtBacktest) SetLimit(times interface{}) float64 {
	return 0.0
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *BtBacktest) AutoSleep() {

}

// GetMinAmount get the min trade amonut of this exchange
func (e *BtBacktest) GetMinAmount(stock string) float64 {
	return 0.0
}

// getAuthJSON
func (e *BtBacktest) getAuthJSON(method string, params ...interface{}) (jsoner *simplejson.Json, err error) {
	data := []byte{}
	return simplejson.NewJson(data)
}

// GetAccount get the account detail of this exchange
func (e *BtBacktest) GetAccount() interface{} {
	return map[string]float64{}
}

// Trade place an order
func (e *BtBacktest) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	return map[string]float64{}
}

func (e *BtBacktest) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	return map[string]float64{}
}

func (e *BtBacktest) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	return map[string]float64{}
}

// GetOrder get details of an order
func (e *BtBacktest) GetOrder(stockType, id string) interface{} {
	return map[string]float64{}
}

// GetOrders get all unfilled orders
func (e *BtBacktest) GetOrders(stockType string) interface{} {
	return map[string]float64{}
}

// GetTrades get all filled orders recently
func (e *BtBacktest) GetTrades(stockType string) interface{} {
	return map[string]float64{}
}

// CancelOrder cancel an order
func (e *BtBacktest) CancelOrder(order Order) bool {
	return true
}

// getTicker get market ticker & depth
func (e *BtBacktest) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, end bool, err error) {
	// @Todo convert data to tick
	_, end, err = e.Run2Data()
	if err != nil{
		return ticker, false, err
	}
	if end == true {
		return ticker, true, errors.New("backTest end")
	}
	return
}

// GetTicker get market ticker & depth
func (e *BtBacktest) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, end, err := e.getTicker(stockType, sizes...)
	if err != nil {
		if end {
			e.logger.Log(constant.INFO, "", 0.0, 0.0, err)
		}else {
			e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		}
		return false
	}
	return ticker
}

// GetRecords get candlestick data
func (e *BtBacktest) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return map[string]float64{}
}
