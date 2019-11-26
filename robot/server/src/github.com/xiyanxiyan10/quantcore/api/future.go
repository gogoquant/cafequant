package api

import (
	goex "github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/builder"
	"github.com/xiyanxiyan10/quantcore/constant"
	"github.com/xiyanxiyan10/quantcore/model"
	"github.com/xiyanxiyan10/quantcore/util"
	"time"
)

// FMEX the exchange struct of fmex.com
type FMEX struct {
	BaseExchange
	stockTypeMap     map[string]goex.CurrencyPair
	tradeTypeMap     map[int]string
	recordsPeriodMap map[string]string
	minAmountMap     map[string]float64
	records          map[string][]Record
	host             string
	logger           model.Logger
	option           Option

	limit     float64
	lastSleep int64
	lastTimes int64

	apiBuilder *builder.APIBuilder
	api        goex.FutureRestAPI
}

// NewFMEX create an exchange struct of fmex.com
func NewFMEX(opt Option) Exchange {
	fmex := FMEX{
		stockTypeMap: map[string]goex.CurrencyPair{
			"BTC/USD": goex.BTC_USD,
		},
		tradeTypeMap: map[int]string{
			goex.OPEN_BUY:   constant.TradeTypeLong,
			goex.OPEN_SELL:  constant.TradeTypeShort,
			goex.CLOSE_BUY:  constant.TradeTypeLongClose,
			goex.CLOSE_SELL: constant.TradeTypeShortClose,
		},
		recordsPeriodMap: map[string]string{
			"M":   "1min",
			"M5":  "5min",
			"M15": "15min",
			"M30": "30min",
			"H":   "1hour",
			"D":   "1day",
			"W":   "1week",
		},
		minAmountMap: map[string]float64{
			"BTC/USD": 0.001,
		},
		records:    make(map[string][]Record),
		host:       "https://www.fmex.com/api/v1/",
		logger:     model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:     opt,
		limit:      10.0,
		lastSleep:  time.Now().UnixNano(),
		apiBuilder: builder.NewAPIBuilder().HttpTimeout(5 * time.Second),
	}
	if fmex.apiBuilder != nil {
		fmex.api = fmex.apiBuilder.APIKey(opt.AccessKey).APISecretkey(opt.SecretKey).BuildFuture(goex.FMEX)
	}
	return &fmex
}

// Log print something to console
func (e *FMEX) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *FMEX) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *FMEX) GetName() string {
	return e.option.Name
}

func (e *FMEX) GetDepth(size int, stockType string) interface{} {
	exchangeStockType, ok := e.stockTypeMap[stockType]
	if !ok {
		return false
	}
	depth, err := e.api.GetFutureDepth(exchangeStockType, e.GetContractType(), size)
	if err != nil {
		return false
	}
	return depth
}

// SetLimit set the limit calls amount per second of this exchange
func (e *FMEX) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *FMEX) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount get the min trade amonut of this exchange
func (e *FMEX) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *FMEX) GetAccount() interface{} {
	userInfo := make(map[string]float64)
	account, err := e.api.GetFutureUserinfo()
	if err != nil {
		return false
	}

	for k, v := range account.FutureSubAccounts {
		stockType := k.Symbol
		userInfo[stockType+"_"+constant.AccountRights] = v.AccountRights
		userInfo[stockType+"_"+constant.KeepDeposit] = v.KeepDeposit
		userInfo[stockType+"_"+constant.ProfitReal] = v.ProfitReal
		userInfo[stockType+"_"+constant.ProfitUnreal] = v.ProfitUnreal
		userInfo[stockType+"_"+constant.RiskRate] = v.RiskRate
	}
	return userInfo
}

func (e *FMEX) Buy(price, amount string, msg ...interface{}) interface{} {
	var err error
	var openType int
	stockType := e.GetStockType()
	exchangeStockType, ok := e.stockTypeMap[stockType]
	if !ok {
		return false
	}
	level := e.GetMarginLevel()
	var matchPrice = 0
	if price == "0" {
		matchPrice = 1
	}
	if e.direction == constant.TradeTypeLong {
		openType = goex.OPEN_BUY
	} else if e.direction == constant.TradeTypeShortClose {
		openType = goex.CLOSE_SELL
	} else {
		return false
	}
	orderId, err := e.api.PlaceFutureOrder(exchangeStockType, e.GetContractType(),
		price, amount, openType, matchPrice, level)

	if err != nil {
		return false
	}
	priceFloat := conver.Float64Must(price)
	amountFloat := conver.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
	return orderId
}

func (e *FMEX) Sell(price, amount string, msg ...interface{}) interface{} {
	var err error
	var openType int
	stockType := e.GetStockType()
	exchangeStockType, ok := e.stockTypeMap[stockType]
	if !ok {
		return false
	}
	level := e.GetMarginLevel()
	var matchPrice = 0
	if price == "0" {
		matchPrice = 1
	}
	if e.direction == constant.TradeTypeShort {
		openType = goex.OPEN_SELL
	} else if e.direction == constant.TradeTypeLongClose {
		openType = goex.CLOSE_BUY
	} else {
		return false
	}
	orderId, err := e.api.PlaceFutureOrder(exchangeStockType, e.GetContractType(),
		price, amount, openType, matchPrice, level)

	if err != nil {
		return false
	}
	priceFloat := conver.Float64Must(price)
	amountFloat := conver.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
	return orderId
}

// GetOrder get details of an order
func (e *FMEX) GetOrder(id string) interface{} {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		return false
	}
	order, err := e.api.GetFutureOrder(id, exchangeStockType, e.GetContractType())
	if err != nil {
		return false
	}
	return Order{
		ID:         order.OrderID2,
		Price:      order.Price,
		Amount:     order.Amount,
		DealAmount: order.DealAmount,
		TradeType:  e.tradeTypeMap[order.OrderType],
		StockType:  e.GetStockType(),
	}
}

// GetOrders get all unfilled orders
func (e *FMEX) GetOrders() interface{} {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		return false
	}
	orders, err := e.api.GetUnfinishFutureOrders(exchangeStockType, e.contractType)
	if err != nil {
		return false
	}
	var resOrders []Order
	for _, order := range orders {
		resOrder := Order{
			ID:         order.OrderID2,
			Price:      order.Price,
			Amount:     order.Amount,
			DealAmount: order.DealAmount,
			TradeType:  e.tradeTypeMap[order.OrderType],
			StockType:  e.GetStockType(),
		}
		resOrders = append(resOrders, resOrder)
	}
	return resOrders
}

// GetTrades get all filled orders recently
func (e *FMEX) GetTrades() interface{} {
	return false
}

// CancelOrder cancel an order
func (e *FMEX) CancelOrder(orderID string) bool {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		return false
	}
	result, err := e.api.FutureCancelOrder(exchangeStockType, e.GetContractType(), orderID)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, the error number is ", err.Error())
		return false
	}
	if !result {
		return false
	}
	e.logger.Log(constant.TradeTypeCancelOrder, e.GetStockType(), 0, 0, orderID)
	return true
}

// getTicker get market ticker & depth
func (e *FMEX) getTicker(sizes ...interface{}) (ticker Ticker, err error) {
	return Ticker{}, nil
}

// GetTicker get market ticker & depth
func (e *FMEX) GetTicker(sizes ...interface{}) interface{} {
	return false
}

// GetRecords get candlestick data
func (e *FMEX) GetRecords(period string, sizes ...interface{}) interface{} {
	return false
}
