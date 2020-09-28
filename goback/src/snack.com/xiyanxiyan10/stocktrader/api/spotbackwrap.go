package api

import (
	"errors"

	goex "github.com/nntaoli-project/goex"
	"snack.com/xiyanxiyan10/conver"
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// NewSpotBackExchange create an exchange struct of futureExchange.com
func NewSpotBackExchange(opt constant.Option) Exchange {
	exchange := NewExchangeBackLink(opt)
	_ = exchange.Init(opt)
	return exchange
}

// ExchangeBackLink ...
type ExchangeBackLink struct {
	ExchangeBack

	stockTypeMap        map[string]goex.CurrencyPair
	stockTypeMapReverse map[goex.CurrencyPair]string

	tradeTypeMap        map[int]string
	tradeTypeMapReverse map[string]int
	exchangeTypeMap     map[string]string

	records map[string][]constant.Record
}

// SetTradeTypeMap ...
func (e *ExchangeBackLink) SetTradeTypeMap(key int, val string) {
	e.tradeTypeMap[key] = val
}

// SetTradeTypeMapReverse ...
func (e *ExchangeBackLink) SetTradeTypeMapReverse(key string, val int) {
	e.tradeTypeMapReverse[key] = val
}

// NewExchangeBackLink create an exchange struct of futureExchange.com
func NewExchangeBackLink(opt constant.Option) *ExchangeBackLink {
	futureExchange := ExchangeBackLink{
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
			constant.HuoBiDm: goex.HBDM,
		},
		records: make(map[string][]constant.Record),
		//apiBuilder: builder.NewAPIBuilder().HttpTimeout(5 * time.Second),
	}
	opt.Limit = 10.0
	//futureExchange.BaseExchange.Init(opt)
	futureExchange.SetRecordsPeriodMap(map[string]int64{
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
	futureExchange.SetID(opt.Index)
	return &futureExchange
}

// ValidBuy ...
func (e *ExchangeBackLink) ValidBuy() error {
	dir := e.GetDirection()
	if dir == constant.TradeTypeBuy {
		return nil
	}
	if dir == constant.TradeTypeShortClose {
		return nil
	}
	return errors.New("错误buy交易方向: " + e.GetDirection())
}

// ValidSell ...
func (e *ExchangeBackLink) ValidSell() error {
	dir := e.GetDirection()
	if dir == constant.TradeTypeSell {
		return nil
	}
	if dir == constant.TradeTypeLongClose {
		return nil
	}
	return errors.New("错误sell交易方向:" + e.GetDirection())
}

// SetMode ...
func (e *ExchangeBackLink) SetMode(mode int) interface{} {
	return "success"
}

// Init init the instance of this exchange
func (e *ExchangeBackLink) Init(opt constant.Option) error {
	e.BaseExchange.Init(opt)
	for k, v := range e.stockTypeMap {
		e.stockTypeMapReverse[v] = k
	}
	for k, v := range e.tradeTypeMap {
		e.tradeTypeMapReverse[v] = k
	}
	return nil
}

// SetStockTypeMap set stock type map
func (e *ExchangeBackLink) SetStockTypeMap(m map[string]goex.CurrencyPair) {
	e.stockTypeMap = m
}

// GetStockTypeMap get stock type map
func (e *ExchangeBackLink) GetStockTypeMap() map[string]goex.CurrencyPair {
	return e.stockTypeMap
}

// Log print something to console
func (e *ExchangeBackLink) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, e.GetStockType(), 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *ExchangeBackLink) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *ExchangeBackLink) GetName() string {
	return e.option.Name
}

// GetDepth get depth from exchange
func (e *ExchangeBackLink) GetDepth() interface{} {
	stockType := e.GetStockType()
	depth, err := e.ExchangeBack.GetDepth(constant.DepthSize, stockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is stockType")
		return nil
	}
	return depth
}

// GetPosition get position from exchange
func (e *ExchangeBackLink) GetPosition() interface{} {
	stockType := e.GetStockType()
	position, err := e.ExchangeBack.GetDepth(100, stockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetPosition() error, the error number is ", err.Error())
		return nil
	}
	return position
}

// GetMinAmount get the min trade amount of this exchange
func (e *ExchangeBackLink) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *ExchangeBackLink) GetAccount() interface{} {
	account, err := e.ExchangeBack.GetAccount()
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetAccount() error, the error number is ", err.Error())
		return nil
	}
	return account
}

// Buy buy from exchange
func (e *ExchangeBackLink) Buy(price, amount string, msg ...interface{}) interface{} {
	var err error
	var ord *constant.Order
	stockType := e.GetStockType()
	if err := e.ValidBuy(); err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Buy() error, the error number is ", err.Error())
		return nil
	}
	if price == "-1" {
		ord, err = e.ExchangeBack.MarketBuy(amount, price, stockType)
	} else {
		ord, err = e.ExchangeBack.LimitBuy(amount, price, stockType)
	}

	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Buy() error, the error number is ", err.Error())
		return nil
	}
	priceFloat := conver.Float64Must(price)
	amountFloat := conver.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
	return ord.Id

}

// Sell sell from exchange
func (e *ExchangeBackLink) Sell(price, amount string, msg ...interface{}) interface{} {
	var err error
	var ord *constant.Order
	stockType := e.GetStockType()
	if err := e.ValidSell(); err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Sell() error, the error number is ", err.Error())
		return nil
	}
	if price == "-1" {
		ord, err = e.ExchangeBack.MarketSell(amount, price, stockType)
	} else {
		ord, err = e.ExchangeBack.LimitSell(amount, price, stockType)
	}

	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Buy() error, the error number is ", err.Error())
		return nil
	}
	priceFloat := conver.Float64Must(price)
	amountFloat := conver.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
	return ord.Id
}

// GetOrder get detail of an order
func (e *ExchangeBackLink) GetOrder(id string) interface{} {
	order, err := e.ExchangeBack.GetOneOrder(id, e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, conver.Float64Must(id), "GetOrder() error, the error number is ", err.Error())
	}
	return order
}

// CompareOrders ...
func (e *ExchangeBackLink) CompareOrders(lft, rht []constant.Order) bool {
	mp := make(map[string]bool)
	if len(lft) != len(rht) {
		return false
	}
	for _, order := range lft {
		mp[order.Id] = true
	}

	for _, order := range rht {
		_, ok := mp[order.Id]
		if !ok {
			return false
		}
	}
	return true
}

// GetOrders get all unfilled orders
func (e *ExchangeBackLink) GetOrders() interface{} {
	orders, err := e.ExchangeBack.GetUnfinishOrders(e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetOrders() error, the error number is ", err.Error())
		return nil
	}
	return orders
}

// CancelOrder cancel an order
func (e *ExchangeBackLink) CancelOrder(orderID string) interface{} {
	result, err := e.ExchangeBack.CancelOrder(orderID, e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, conver.Float64Must(orderID), "CancelOrder() error, the error number is ", err.Error())
		return nil
	}
	if !result {
		return nil
	}
	e.logger.Log(constant.TradeTypeCancel, e.GetStockType(), 0, conver.Float64Must(orderID), "CancelOrder() success")
	return true
}

// GetTicker get market ticker
func (e *ExchangeBackLink) GetTicker() interface{} {
	ticker, err := e.ExchangeBack.GetTicker(e.GetStockType())
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ", err.Error())
		return nil
	}
	return *ticker
}

// GetRecords get candlestick data
func (e *ExchangeBackLink) GetRecords(periodStr, maStr string) interface{} {
	e.logger.Log(constant.ERROR, e.GetStockType(), 0, 0, "GetRecords() error, the error number is stockType")
	return nil
}
