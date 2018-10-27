package api

import (
	goback "github.com/xiyanxiyan10/gobacktest"
)

// Option is an exchange option
type Option struct {
	TraderID    int64  //trader id
	AlgorithmID int64  //algorithm id
	User        string //user token
	Type        string //exchange type
	Name        string //exchange name
	AccessKey   string //asscess key
	SecretKey   string //secret key
	Mode        string //run mode
	// Ctx       *otto.Otto

	In   *IncomingHandler // get incoming data for api
	Back *goback.Backtest // get back event support

}

// Exchange interface
type Exchange interface {
	Log(...interface{})
	GetType() string
	GetName() string
	SetLimit(times interface{}) float64
	AutoSleep()
	GetMinAmount(stock string) float64
	GetAccount() interface{}
	Trade(tradeType string, stockType string, price, amount interface{}, msgs ...interface{}) interface{}
	GetOrder(stockType, id string) interface{}
	GetOrders(stockType string) interface{}
	GetTrades(stockType string) interface{}
	CancelOrder(order Order) bool
	SetSubscribe(symbol string) error
	StockMap() map[string]string
	SetStockMap(m map[string]string)
	GetTicker(stockType string, sizes ...interface{}) interface{}
	GetRecords(stockType, period string, sizes ...interface{}) interface{}

	Draw(map[string]interface{}) interface{}

	Start() error
	Stop() error
	Status() int
}

//Entity transaction
type Etrade interface {
	Log(...interface{})
	GetList(name string, company string, minprice float64, maxprice float64) interface{}
}

var (
	constructor = map[string]func(Option) Exchange{}
)
