package api

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// CTPExchange the exchange struct of futureExchange.com
type CTPExchange struct {
	BaseExchange

	tradeTypeMap        map[int]string
	tradeTypeMapReverse map[string]int
	exchangeTypeMap     map[string]string

	records map[string][]constant.Record
}

// NewCTPExchange create an exchange struct of futureExchange.com
func NewCTPExchange(opt constant.Option) *CTPExchange {
	spotExchange := CTPExchange{
		records: make(map[string][]constant.Record),
		//apiBuilder: builder.NewAPIBuilder().HttpTimeout(5 * time.Second),
	}
	spotExchange.SetRecordsPeriodMap(map[string]int64{
		"M1":  60,
		"M5":  300,
		"M15": 900,
		"M30": 1800,
		"H1":  3600,
		"H2":  7200,
	})
	spotExchange.SetMinAmountMap(map[string]float64{
		"BTC/USD": 0.001,
	})
	spotExchange.back = false
	return &spotExchange
}

// Ready ...
func (e *CTPExchange) Ready() interface{} {
	return "success"
}

// Init get the type of this exchange
func (e *CTPExchange) Init(opt constant.Option) error {
	e.BaseExchange.Init(opt)
	return nil
}

// Log print something to console
func (e *CTPExchange) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, e.GetStockType(), 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *CTPExchange) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *CTPExchange) GetName() string {
	return e.option.Name
}

// GetDepth ...
func (e *CTPExchange) GetDepth() interface{} {
	return nil
}

// GetPosition ...
func (e *CTPExchange) GetPosition() interface{} {
	return nil
}

// GetAccount get the account detail of this exchange
func (e *CTPExchange) GetAccount() interface{} {
	return nil
}

// Buy ...
func (e *CTPExchange) Buy(price, amount string, msg ...interface{}) interface{} {
	return nil
}

// Sell ...
func (e *CTPExchange) Sell(price, amount string, msg ...interface{}) interface{} {
	return nil
}

// GetOrder get details of an order
func (e *CTPExchange) GetOrder(id string) interface{} {
	return nil
}

// GetOrders get all unfilled orders
func (e *CTPExchange) GetOrders() interface{} {
	return nil
}

// CancelOrder cancel an order
func (e *CTPExchange) CancelOrder(orderID string) interface{} {
	return nil
}

// GetTicker get market ticker
func (e *CTPExchange) GetTicker() interface{} {
	return nil
}

// GetRecords get candlestick data
func (e *CTPExchange) GetRecords(periodStr string) interface{} {
	return nil
}
