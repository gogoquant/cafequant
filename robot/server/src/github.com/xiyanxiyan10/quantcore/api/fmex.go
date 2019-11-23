package api

import (
	goex "github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/builder"
	"github.com/xiyanxiyan10/quantcore/constant"
	"github.com/xiyanxiyan10/quantcore/model"
	"github.com/xiyanxiyan10/quantcore/util"
	"strconv"

	"strings"
	"time"
)

// FMEX the exchange struct of fmex.com
type FMEX struct {
	stockTypeMap     map[string]goex.CurrencyPair
	tradeTypeMap     map[goex.TradeSide]string
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
	api        goex.API
}

// NewFMEX create an exchange struct of fmex.com
func NewFMEX(opt Option) Exchange {
	fmex := FMEX{
		stockTypeMap: map[string]goex.CurrencyPair{
			"BTC/USD": goex.BTC_USD,
		},
		tradeTypeMap: map[goex.TradeSide]string{
			goex.BUY:         constant.TradeTypeBuy,
			goex.SELL:        constant.TradeTypeSell,
			goex.BUY_MARKET:  constant.TradeTypeBuy,
			goex.SELL_MARKET: constant.TradeTypeSell,
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
		fmex.api = fmex.apiBuilder.APIKey(opt.AccessKey).APISecretkey(opt.SecretKey).Build(goex.FMEX)
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
	depth, err := e.api.GetDepth(size, goex.BTC_USD)
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
	account, err := e.api.GetAccount()
	if err != nil {
		return false
	}
	btcAccount, ok := account.SubAccounts[goex.BTC]
	if !ok {
		return false
	}
	return map[string]float64{
		"BTC":       conver.Float64Must(btcAccount.Amount),
		"FrozenBTC": conver.Float64Must(btcAccount.ForzenAmount),
	}
}

// Trade place an order
func (e *FMEX) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	stockType = strings.ToUpper(stockType)
	tradeType = strings.ToUpper(tradeType)
	price := conver.Float64Must(_price)
	amount := conver.Float64Must(_amount)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized stockType: ", stockType)
		return false
	}
	switch tradeType {
	case constant.TradeTypeBuy:
		return e.buy(stockType, price, amount, msgs...)
	case constant.TradeTypeSell:
		return e.sell(stockType, price, amount, msgs...)
	default:
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized tradeType: ", tradeType)
		return false
	}
}

func (e *FMEX) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	var order *goex.Order
	var err error
	amountStr := strconv.FormatFloat(amount, 'E', -1, 64)
	priceStr := strconv.FormatFloat(price, 'E', -1, 64)
	if amount <= 0 {
		order, err = e.api.MarketBuy(amountStr, "0", goex.BTC_USDT)
	} else {
		order, err = e.api.MarketBuy(amountStr, priceStr, goex.BTC_USDT)
	}
	if err != nil {
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return order.OrderID2
}

func (e *FMEX) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	var order *goex.Order
	var err error
	amountStr := strconv.FormatFloat(amount, 'E', -1, 64)
	priceStr := strconv.FormatFloat(price, 'E', -1, 64)
	if amount <= 0 {
		order, err = e.api.MarketSell(amountStr, "0", goex.BTC_USDT)
	} else {
		order, err = e.api.MarketSell(amountStr, priceStr, goex.BTC_USDT)
	}
	if err != nil {
		return false
	}
	e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
	return order.OrderID2
}

// GetOrder get details of an order
func (e *FMEX) GetOrder(stockType, id string) interface{} {
	order, err := e.api.GetOneOrder(id, goex.BTC_USD)
	if err != nil {
		return false
	}
	return Order{
		ID:         order.OrderID2,
		Price:      order.Price,
		Amount:     order.Amount,
		DealAmount: order.DealAmount,
		TradeType:  e.tradeTypeMap[order.Side],
		StockType:  stockType,
	}
}

// GetOrders get all unfilled orders
func (e *FMEX) GetOrders(stockType string) interface{} {
	orders, err := e.api.GetUnfinishOrders(goex.BTC_USD)
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
			TradeType:  e.tradeTypeMap[order.Side],
			StockType:  stockType,
		}
		resOrders = append(resOrders, resOrder)
	}
	return resOrders
}

// GetTrades get all filled orders recently
func (e *FMEX) GetTrades(stockType string) interface{} {
	return false
}

// CancelOrder cancel an order
func (e *FMEX) CancelOrder(order Order) bool {
	result, err := e.api.CancelOrder(order.ID, goex.BTC_USD)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, the error number is ", err.Error())
		return false
	}
	if !result {
		return false
	}
	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker get market ticker & depth
func (e *FMEX) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	return Ticker{}, nil
}

// GetTicker get market ticker & depth
func (e *FMEX) GetTicker(stockType string, sizes ...interface{}) interface{} {
	return false
}

// GetRecords get candlestick data
func (e *FMEX) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return false
}
