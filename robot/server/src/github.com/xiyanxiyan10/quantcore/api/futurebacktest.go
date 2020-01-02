package api

import (
	goex "github.com/nntaoli-project/GoEx"
	"github.com/xiyanxiyan10/quantcore/constant"
	"github.com/xiyanxiyan10/quantcore/model"
	"github.com/xiyanxiyan10/quantcore/util"
	"time"
)

// FutureBackTestExchange the exchange struct of futureExchange.com
type FutureBackTestExchange struct {
	BaseExchange
	stockTypeMap        map[string]string
	stockTypeMapReverse map[string]string

	tradeTypeMap        map[int]string
	tradeTypeMapReverse map[string]int
	exchangeTypeMap     map[string]string

	records map[string][]constant.Record
	host    string
	logger  model.Logger
	option  constant.Option

	limit     float64
	lastSleep int64
	lastTimes int64
}

// NewFutureExchange create an exchange struct of futureExchange.com
func NewFutureBackTestExchange(opt constant.Option) *FutureExchange {
	futureExchange := FutureExchange{
		stockTypeMap: map[string]goex.CurrencyPair{
			"BTC/USD": goex.BTC_USD,
		},
		stockTypeMapReverse: map[goex.CurrencyPair]string{},
		tradeTypeMapReverse: map[string]int{},
		tradeTypeMap: map[int]string{
			goex.OPEN_BUY:   constant.TradeTypeLong,
			goex.OPEN_SELL:  constant.TradeTypeShort,
			goex.CLOSE_BUY:  constant.TradeTypeLongClose,
			goex.CLOSE_SELL: constant.TradeTypeShortClose,
		},

		exchangeTypeMap: map[string]string{
			constant.Fmex:    goex.FMEX,
			constant.HuobiDm: goex.HBDM,
		},
		records:   make(map[string][]constant.Record),
		host:      "https://www.futureExchange.com/api/v1/",
		logger:    model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:    opt,
		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
		//apiBuilder: builder.NewAPIBuilder().HttpTimeout(5 * time.Second),
	}
	futureExchange.SetRecordsPeriodMap(map[string]int{
		"M1":  goex.KLINE_PERIOD_1MIN,
		"M5":  goex.KLINE_PERIOD_5MIN,
		"M15": goex.KLINE_PERIOD_15MIN,
		"M30": goex.KLINE_PERIOD_30MIN,
		"H1":  goex.KLINE_PERIOD_1H,
		"H2":  goex.KLINE_PERIOD_4H,
		"H4":  goex.KLINE_PERIOD_4H,
		"D1":  goex.KLINE_PERIOD_1DAY,
		"W1":  goex.KLINE_PERIOD_1WEEK,
	})
	futureExchange.SetMinAmountMap(map[string]float64{
		"BTC/USD": 0.001,
	})
	return &futureExchange
}

// Init Init this exchange
func (e *FutureBackTestExchange) Init() error {
	for k, v := range e.stockTypeMap {
		e.stockTypeMapReverse[v] = k
	}
	for k, v := range e.tradeTypeMap {
		e.tradeTypeMapReverse[v] = k
	}
	return nil
}

// SetStockTypeMap ...
func (e *FutureBackTestExchange) SetStockTypeMap(m map[string]string) {
	e.stockTypeMap = m
}

func (e *FutureBackTestExchange) GetStockTypeMap() map[string]string {
	return e.stockTypeMap
}

// Log print something to console
func (e *FutureBackTestExchange) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *FutureBackTestExchange) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *FutureBackTestExchange) GetName() string {
	return e.option.Name
}

func (e *FutureBackTestExchange) GetDepth(size int) interface{} {
	return false
}

func (e *FutureBackTestExchange) GetPosition() interface{} {
	return false
}

// SetLimit set the limit calls amount per second of this exchange
func (e *FutureBackTestExchange) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *FutureBackTestExchange) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amount of this exchange
func (e *FutureBackTestExchange) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *FutureBackTestExchange) GetAccount() interface{} {
	return false
}

func (e *FutureBackTestExchange) Buy(price, amount string, msg ...interface{}) interface{} {
	return false
}

func (e *FutureBackTestExchange) Sell(price, amount string, msg ...interface{}) interface{} {
	return false
}

// GetOrder get details of an order
func (e *FutureBackTestExchange) GetOrder(id string) interface{} {
	return false
}

// GetOrders get all unfilled orders
func (e *FutureBackTestExchange) GetOrders() interface{} {
	return false
}

// GetTrades get all filled orders recently
func (e *FutureBackTestExchange) GetTrades(params ...interface{}) interface{} {
	return false
}

// CancelOrder cancel an order
func (e *FutureBackTestExchange) CancelOrder(orderID string) bool {
	return true
}

// GetTicker get market ticker
func (e *FutureBackTestExchange) GetTicker() interface{} {
	return false
}

// GetRecords get candlestick data
// params[0] period
// params[1] size
// params[2] since
func (e *FutureBackTestExchange) GetRecords(params ...interface{}) interface{} {
	return false
}
