package api

import (
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

const (
	tickerURL = "http://hq.sinajs.cn/list="
	depthURL  = "http://hq.sinajs.cn/list="
	recordURL = "http://money.finance.sina.com.cn/quotes_service/api/json_v2.php/CN_MarketData.getKLineData?"
)

// SZExchange the exchange struct of futureExchange.com
type SZExchange struct {
	BaseExchange

	tradeTypeMap        map[int]string
	tradeTypeMapReverse map[string]int
	exchangeTypeMap     map[string]string

	records map[string][]constant.Record
}

// NewSZExchange create an exchange struct of futureExchange.com
func NewSZExchange(opt constant.Option) *SZExchange {
	spotExchange := SZExchange{
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
func (e *SZExchange) Ready() interface{} {
	return "success"
}

// Init get the type of this exchange
func (e *SZExchange) Init(opt constant.Option) error {
	e.BaseExchange.Init(opt)
	return nil
}

// Log print something to console
func (e *SZExchange) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, e.GetStockType(), 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *SZExchange) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *SZExchange) GetName() string {
	return e.option.Name
}

// GetDepth ...
func (e *SZExchange) GetDepth() interface{} {
	return nil
}

// GetPosition ...
func (e *SZExchange) GetPosition() interface{} {
	return nil
}

// GetAccount get the account detail of this exchange
func (e *SZExchange) GetAccount() interface{} {
	return nil
}

// Buy ...
func (e *SZExchange) Buy(price, amount string, msg ...interface{}) interface{} {
	return nil
}

// Sell ...
func (e *SZExchange) Sell(price, amount string, msg ...interface{}) interface{} {
	return nil
}

// GetOrder get details of an order
func (e *SZExchange) GetOrder(id string) interface{} {
	return nil
}

// GetOrders get all unfilled orders
func (e *SZExchange) GetOrders() interface{} {
	return nil
}

// CancelOrder cancel an order
func (e *SZExchange) CancelOrder(orderID string) interface{} {
	return nil
}

// GetTicker get market ticker
func (e *SZExchange) GetTicker() interface{} {
	return nil
}

// GetRecords get candlestick data
func (e *SZExchange) GetRecords(periodStr string) interface{} {
	return nil
}
