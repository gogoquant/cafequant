package api

import (
	"fmt"
	"time"

	goex "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"

	"snack.com/xiyanxiyan10/conver"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
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
	spotExchange.SetMinAmountMap(map[string]float64{
		"BTC/USD": 0.001,
	})
	spotExchange.back = false
	return &spotExchange
}

// Stop ...
func (e *SpotExchange) Stop() error {
	return nil
}

// Start ...
func (e *SpotExchange) Start() error {
	proxyURL := config.String("proxy")
	if proxyURL == "" {
		e.apiBuilder = builder.NewAPIBuilder().HttpTimeout(2 * time.Second)
	} else {
		e.apiBuilder = builder.NewAPIBuilder().HttpProxy(proxyURL).HttpTimeout(2 * time.Second)
	}
	if e.apiBuilder == nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "api builder fail")
		return fmt.Errorf("api builder fail")
	}
	exchangeName := e.exchangeTypeMap[e.option.Type]
	e.api = e.apiBuilder.APIKey(e.option.AccessKey).APISecretkey(e.option.SecretKey).Build(exchangeName)
	return nil
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

// GetDepth ...
func (e *SpotExchange) GetDepth() (*constant.Depth, error) {
	var resDepth constant.Depth
	stockType := e.GetStockType()
	exchangeStockType, ok := e.stockTypeMap[stockType]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is stockType")
		return nil, fmt.Errorf("GetDepth() error, the error number is stockType")
	}
	depth, err := e.api.GetDepth(constant.DepthSize, exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetDepth() error, the error number is %s", err.Error())
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
	return &resDepth, nil
}

// GetPosition ...
func (e *SpotExchange) GetPosition() ([]constant.Position, error) {
	return nil, fmt.Errorf("not support")
}

// GetAccount get the account detail of this exchange
func (e *SpotExchange) GetAccount() (*constant.Account, error) {
	account, err := e.api.GetAccount()
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetAccount() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetAccount() error, the error number is %s", err.Error())
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
	return &resAccount, nil
}

// Buy ...
func (e *SpotExchange) Buy(price, amount string, msg ...interface{}) (string, error) {
	var err error
	var order *goex.Order
	stockType := e.GetStockType()
	exchangeStockType, ok := e.stockTypeMap[stockType]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Buy() error, the error number is stockType")
		return "", fmt.Errorf("Buy() error, the error number is stockType")
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
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Sell() error, the error number is ", err.Error())
		return "", fmt.Errorf("Sell() error, the error number is %s", err.Error())
	}
	priceFloat := conver.Float64Must(price)
	amountFloat := conver.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
	return order.Cid, nil
}

// Sell ...
func (e *SpotExchange) Sell(price, amount string, msg ...interface{}) (string, error) {
	var err error
	var order *goex.Order
	stockType := e.GetStockType()
	exchangeStockType, ok := e.stockTypeMap[stockType]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Sell() error, the error number is stockType")
		return "", fmt.Errorf("Sell() error, the error number is stockType")
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
		e.logger.Log(constant.ERROR, e.GetStockType(), conver.Float64Must(amount), conver.Float64Must(amount), "Sell() error, the error number is ", err.Error())
		return "", fmt.Errorf("Sell() error, the error number is %s", err.Error())
	}
	priceFloat := conver.Float64Must(price)
	amountFloat := conver.Float64Must(amount)
	e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
	return order.Cid, nil
}

// GetOrder get details of an order
func (e *SpotExchange) GetOrder(id string) (*constant.Order, error) {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, conver.Float64Must(id), "GetOrder() error, the error number is stockType")
		return nil, fmt.Errorf("GetOrder() error, the error number is stockType")
	}
	order, err := e.api.GetOneOrder(id, exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, conver.Float64Must(id), "GetOrder() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetOrder() error, the error number is %s", err.Error())
	}

	var TradeType string
	if order.OrderType == goex.OPEN_BUY {
		TradeType = constant.TradeTypeBuy
	} else {
		TradeType = constant.TradeTypeSell
	}

	return &constant.Order{
		Id:         order.OrderID2,
		Price:      order.Price,
		Amount:     order.Amount,
		DealAmount: order.DealAmount,
		TradeType:  TradeType,
		StockType:  e.GetStockType(),
	}, nil
}

// GetOrders get all unfilled orders
func (e *SpotExchange) GetOrders() ([]constant.Order, error) {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		e.logger.Log(constant.ERROR, "", 0, 0, "GetOrders() error, the error number is stockType")
		return nil, fmt.Errorf("GetOrders() error, the error number is stockType")
	}
	orders, err := e.api.GetUnfinishOrders(exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetOrders() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetOrders() error, the error number is %s", err.Error())
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
	return resOrders, nil
}

// CancelOrder cancel an order
func (e *SpotExchange) CancelOrder(orderID string) (bool, error) {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, conver.Float64Must(orderID), "CancelOrder() error, the error number is stockType")
		return false, fmt.Errorf("CancelOrder() error, the error number is stockType")
	}
	result, err := e.api.CancelOrder(orderID, exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, conver.Float64Must(orderID), "CancelOrder() error, the error number is ", err.Error())
		return false, fmt.Errorf("CancelOrder() error, the error number is %s", err.Error())
	}
	if !result {
		return false, fmt.Errorf("CancelOrder() error, the error number is false")
	}
	e.logger.Log(constant.TradeTypeCancel, e.GetStockType(), 0, conver.Float64Must(orderID), "CancelOrder() success")
	return true, nil
}

// GetTicker get market ticker
func (e *SpotExchange) GetTicker() (*constant.Ticker, error) {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	if !ok {
		e.logger.Log(constant.ERROR, "", 0, 0, "GetTicker() error, the error number is stockType")
		return nil, fmt.Errorf("GetTicker() error, the error number is stockType")
	}
	exTicker, err := e.api.GetTicker(exchangeStockType)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetTicker() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetTicker() error, the error number is %s", err.Error())
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
		Time: conver.Int64Must(tickStr),
	}
	return &ticker, nil
}

// GetRecords get candlestick data
func (e *SpotExchange) GetRecords() ([]constant.Record, error) {
	exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
	var period int64 = -1
	var since = 0
	periodStr := e.GetPeriod()
	size := e.GetSize()
	period, ok = e.recordsPeriodMap[periodStr]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, 0, "GetRecords() error, the error number is stockType")
		return nil, fmt.Errorf("GetRecords() error, the error number is stockType")
	}

	klineVec, err := e.api.GetKlineRecords(exchangeStockType, int(period), size, since)
	if err != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetRecords() error, the error number is ", err.Error())
		return nil, fmt.Errorf("GetRecords() error, the error number is %s", err.Error())
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
	return e.records[periodStr], nil
}
