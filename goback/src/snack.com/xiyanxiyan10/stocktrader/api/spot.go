package api

import (
	"errors"
	"fmt"
	"time"

	goex "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"

	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

// SpotExchange the exchange struct of futureExchange.com
type SpotExchange struct {
	BaseExchange
	stockTypeMap        map[string]goex.CurrencyPair
	stockTypeMapReverse map[goex.CurrencyPair]string

	tradeTypeMap        map[int]string
	tradeTypeMapReverse map[string]int
	exchangeTypeMap     map[string]string

	records map[string][]constant.Record

	apiBuilder *builder.APIBuilder
	api        goex.API
}

// NewSpotExchange create an exchange struct of futureExchange.com
func NewSpotExchange(opt constant.Option) *SpotExchange {
	spotExchange := SpotExchange{
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
			constant.HuoBi: goex.HUOBI_PRO,
		},
		records: make(map[string][]constant.Record),
		//apiBuilder: builder.NewAPIBuilder().HttpTimeout(5 * time.Second),
	}
	spotExchange.SetRecordsPeriodMap(map[string]int64{
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
	spotExchange.SetMinAmountMap(map[string]float64{
		"BTC/USD": 0.001,
	})
	return &spotExchange
}

// Init get the type of this exchange
func (e *SpotExchange) Init(opt constant.Option) error {
	e.BaseExchange.Init(opt)
	for k, v := range e.stockTypeMap {
		e.stockTypeMapReverse[v] = k
	}
	for k, v := range e.tradeTypeMap {
		e.tradeTypeMapReverse[v] = k
	}
	proxyURL := config.String("proxy")
	if proxyURL == "" {
		e.apiBuilder = builder.NewAPIBuilder().HttpTimeout(2 * time.Second)
	} else {
		e.apiBuilder = builder.NewAPIBuilder().HttpProxy(proxyURL).HttpTimeout(2 * time.Second)
	}
	if e.apiBuilder == nil {
		return errors.New("api builder fail")
	}
	exchangeName := e.exchangeTypeMap[e.option.Type]
	e.api = e.apiBuilder.APIKey(e.option.AccessKey).APISecretkey(e.option.SecretKey).Build(exchangeName)
	return nil
}

// GetBackBar ...
func (e *SpotExchange) GetBackBar() interface{} {
	return nil
}

// SetStockTypeMap ...
func (e *SpotExchange) SetStockTypeMap(m map[string]goex.CurrencyPair) {
	e.stockTypeMap = m
}

// GetStockTypeMap ...
func (e *SpotExchange) GetStockTypeMap() map[string]goex.CurrencyPair {
	return e.stockTypeMap
}

// Log print something to console
func (e *SpotExchange) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, e.GetStockType(), 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *SpotExchange) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *SpotExchange) GetName() string {
	return e.option.Name
}

// GetDepth ...
func (e *SpotExchange) GetDepth(size int) interface{} {
	var resDepth constant.Depth
	stockType := e.GetStockType()
	exchangeStockType, ok := e.stockTypeMap[stockType]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is stockType")
		return nil
	}
	depth, err := e.api.GetDepth(size, exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is ", err.Error())
		return nil
	}
	resDepth.Time = depth.UTime.Unix()
	for _, ask := range depth.AskList {
		var resAsk constant.DepthRecord
		resAsk.Amount = ask.Amount
		resAsk.Price = ask.Price
		resDepth.Asks = append(resDepth.Asks, resAsk)
	}
	for _, bid := range depth.BidList {
		var resBid constant.DepthRecord
		resBid.Amount = bid.Amount
		resBid.Price = bid.Price
		resDepth.Bids = append(resDepth.Bids, resBid)
	}
	return resDepth
}

// GetPosition ...
func (e *SpotExchange) GetPosition() interface{} {
	return nil
}

// GetMinAmount get the min trade amount of this exchange
func (e *SpotExchange) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *SpotExchange) GetAccount() interface{} {
	account, err := e.api.GetAccount()
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetAccount() error, the error number is ", err.Error())
		return nil
	}
	var resAccount constant.Account
	resAccount.SubAccounts = make(map[string]constant.SubAccount)
	for k, v := range account.SubAccounts {
		var subAccount constant.SubAccount
		stockType := k.Symbol
		subAccount.Amount = v.Amount
		subAccount.FrozenAmount = v.ForzenAmount
		subAccount.LoanAmount = v.LoanAmount
		resAccount.SubAccounts[stockType] = subAccount
	}
	return resAccount
}

func (e *SpotExchange) Buy(price, amount string, msg ...interface{}) interface{} {
	var err error
	var order *goex.Order
	stockType := e.GetStockType()
	exchangeStockType, ok := e.stockTypeMap[stockType]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Buy() error, the error number is stockType")
		return nil
	}
	var matchPrice = 0
	if price == "-1" {
		matchPrice = 1
	}
	if matchPrice == 1 {
		order, err = e.api.MarketBuy(amount, price, exchangeStockType)
	} else {
		order, err = e.api.LimitBuy(amount, price, exchangeStockType)
	}

	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Sell() error, the error number is ", err.Error())
		return nil
	}
	priceFloat := util.Float64Must(price)
	amountFloat := util.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
	return order.Cid
}

func (e *SpotExchange) Sell(price, amount string, msg ...interface{}) interface{} {
	var err error
	var order *goex.Order
	stockType := e.GetStockType()
	exchangeStockType, ok := e.stockTypeMap[stockType]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Sell() error, the error number is stockType")
		return nil
	}
	var matchPrice = 0
	if price == "-1" {
		matchPrice = 1
	}
	if matchPrice == 1 {
		order, err = e.api.MarketSell(amount, price, exchangeStockType)
	} else {
		order, err = e.api.LimitSell(amount, price, exchangeStockType)
	}

	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Sell() error, the error number is ", err.Error())
		return nil
	}
	priceFloat := util.Float64Must(price)
	amountFloat := util.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
	return order.Cid
}

// GetOrder get details of an order
func (e *SpotExchange) GetOrder(id string) interface{} {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, util.Float64Must(id), "GetOrder() error, the error number is stockType")
		return nil
	}
	order, err := e.api.GetOneOrder(id, exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, util.Float64Must(id), "GetOrder() error, the error number is ", err.Error())
		return nil
	}

	var TradeType string
	if order.OrderType == goex.OPEN_BUY {
		TradeType = constant.TradeTypeBuy
	} else {
		TradeType = constant.TradeTypeSell
	}

	return constant.Order{
		Id:         order.OrderID2,
		Price:      order.Price,
		Amount:     order.Amount,
		DealAmount: order.DealAmount,
		TradeType:  TradeType,
		StockType:  e.GetStockType(),
	}
}

// GetOrders get all unfilled orders
func (e *SpotExchange) GetOrders() interface{} {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		e.logger.Log(constant.ERROR, "", 0, 0, "GetOrders() error, the error number is stockType")
		return nil
	}
	orders, err := e.api.GetUnfinishOrders(exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetOrders() error, the error number is ", err.Error())
		return nil
	}
	resOrders := []constant.Order{}
	for _, order := range orders {
		var TradeType string
		if order.OrderType == goex.OPEN_BUY {
			TradeType = constant.TradeTypeBuy
		} else {
			TradeType = constant.TradeTypeSell
		}
		resOrder := constant.Order{
			Id:         order.OrderID2,
			Price:      order.Price,
			Amount:     order.Amount,
			DealAmount: order.DealAmount,
			TradeType:  TradeType,
			StockType:  e.GetStockType(),
		}
		resOrders = append(resOrders, resOrder)
	}
	return resOrders
}

// GetTrades get all filled orders recently
func (e *SpotExchange) GetTrades(params ...interface{}) interface{} {
	var traders []constant.Trader
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, 0, "GetTrades() error, the error number is stockType")
		return nil
	}
	APITraders, err := e.api.GetTrades(exchangeStockType, 0)
	if err != nil {
		return nil
	}
	for _, APITrader := range APITraders {
		trader := constant.Trader{
			Id:        APITrader.Tid,
			TradeType: e.tradeTypeMap[int(APITrader.Type)],
			Amount:    APITrader.Amount,
			Price:     APITrader.Price,
			StockType: e.stockTypeMapReverse[APITrader.Pair],
			Time:      APITrader.Date,
		}
		traders = append(traders, trader)
	}
	return traders
}

// CancelOrder cancel an order
func (e *SpotExchange) CancelOrder(orderID string) interface{} {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, util.Float64Must(orderID), "CancelOrder() error, the error number is stockType")
		return nil
	}
	result, err := e.api.CancelOrder(orderID, exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, util.Float64Must(orderID), "CancelOrder() error, the error number is ", err.Error())
		return nil
	}
	if !result {
		return nil
	}
	e.logger.Log(constant.TradeTypeCancel, e.GetStockType(), 0, util.Float64Must(orderID), "CancelOrder() success")
	return true
}

// GetTicker get market ticker
func (e *SpotExchange) GetTicker() interface{} {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		e.logger.Log(constant.ERROR, "", 0, 0, "GetTicker() error, the error number is stockType")
		return nil
	}
	exTicker, err := e.api.GetTicker(exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ", err.Error())
		return nil
	}
	//force covert
	tickStr := fmt.Sprint(exTicker.Date)
	ticker := constant.Ticker{
		Last: exTicker.Last,
		Buy:  exTicker.Buy,
		Sell: exTicker.Sell,
		High: exTicker.High,
		Low:  exTicker.Low,
		Vol:  exTicker.Vol,
		Time: util.Int64Must(tickStr),
	}
	return ticker
}

// GetRecords get candlestick data
func (e *SpotExchange) GetRecords(periodStr string) interface{} {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	var period int64 = -1
	var size = constant.RecordSize
	var since = 0

	period, ok = e.recordsPeriodMap[periodStr]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, 0, "GetRecords() error, the error number is stockType")
		return nil
	}

	klineVec, err := e.api.GetKlineRecords(exchangeStockType, int(period), size, since)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetRecords() error, the error number is ", err.Error())
		return nil
	}
	timeLast := int64(0)
	if len(e.records[periodStr]) > 0 {
		timeLast = e.records[periodStr][len(e.records[periodStr])-1].Time
	}
	var recordsNew []constant.Record
	for i := len(klineVec); i > 0; i-- {
		kline := klineVec[i-1]
		recordTime := kline.Timestamp
		if recordTime > timeLast {
			recordsNew = append([]constant.Record{{
				Time:   recordTime,
				Open:   kline.Open,
				High:   kline.High,
				Low:    kline.Low,
				Close:  kline.Close,
				Volume: kline.Vol,
			}}, recordsNew...)
		} else if timeLast > 0 && recordTime == timeLast {
			e.records[periodStr][len(e.records[periodStr])-1] = constant.Record{
				Time:   recordTime,
				Open:   kline.Open,
				High:   kline.High,
				Low:    kline.Low,
				Close:  kline.Close,
				Volume: kline.Vol,
			}
		} else {
			break
		}
	}
	e.records[periodStr] = append(e.records[periodStr], recordsNew...)
	if len(e.records[periodStr]) > size {
		e.records[periodStr] = e.records[periodStr][len(e.records[periodStr])-size : len(e.records[periodStr])]
	}
	return e.records[periodStr]
}
