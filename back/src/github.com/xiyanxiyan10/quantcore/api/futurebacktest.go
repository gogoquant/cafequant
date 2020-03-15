package api

/*
import (
	"errors"
	"github.com/nntaoli-project/goex"
	"github.com/xiyanxiyan10/quantcore/constant"
	"github.com/xiyanxiyan10/quantcore/model"
	"github.com/xiyanxiyan10/quantcore/util"
	sdkConstant "github.com/xiyanxiyan10/stockdb/constant"
	"github.com/xiyanxiyan10/stockdb/sdk"
	sdkType "github.com/xiyanxiyan10/stockdb/types"
	"time"
)

// FutureBackTest the exchange struct of futureExchange.com
type FutureBackTest struct {
	BaseExchange
	stockTypeMap        map[string]goex.CurrencyPair
	stockTypeMapReverse map[goex.CurrencyPair]string

	tradeTypeMap        map[int]string
	tradeTypeMapReverse map[string]int
	exchangeTypeMap     map[string]string

	records map[string][]constant.Record
	host    string
	logger  model.Logger
	option  constant.Option

	backTestCurr    int64
	backTestBegin   int64
	backTestEnd     int64
	backTestIndex   int64
	backTestTickers []sdkType.OHLC
	limit           float64
	lastSleep       int64
	lastTimes       int64

	client *sdk.Client
}

// NewFutureExchange create an exchange struct of futureExchange.com
func NewFutureBackTest(opt constant.Option) *FutureBackTest {
	futureExchange := FutureBackTest{
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
			constant.HuoBiDm: goex.HBDM,
		},
		records:   make(map[string][]constant.Record),
		host:      "https://www.futureExchange.com/api/v1/",
		logger:    model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:    opt,
		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
		//apiBuilder: builder.NewAPIBuilder().HttpTimeout(5 * time.Second),
	}
	futureExchange.SetRecordsPeriodMap(map[string]int64{
		"M1":  sdkConstant.Minute,
		"M5":  sdkConstant.Minute * 5,
		"M15": sdkConstant.Minute * 15,
		"M30": sdkConstant.Minute * 30,
		"H1":  sdkConstant.Hour,
		"H2":  sdkConstant.Hour * 2,
		"H4":  sdkConstant.Hour * 4,
		"D1":  sdkConstant.Day,
		"W1":  sdkConstant.Week,
	})
	futureExchange.SetMinAmountMap(map[string]float64{
		"BTC/USD": 0.001,
	})
	return &futureExchange
}

// ValidBuy ...
func (e *FutureBackTest) ValidBuy() error {
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
func (e *FutureBackTest) ValidSell() error {
	dir := e.GetDirection()
	if dir == constant.TradeTypeSell {
		return nil
	}
	if dir == constant.TradeTypeLongClose {
		return nil
	}
	return errors.New("错误sell交易方向:" + e.GetDirection())
}

// GetType get the type of this exchange
func (e *FutureBackTest) Init() error {
	for k, v := range e.stockTypeMap {
		e.stockTypeMapReverse[v] = k
	}
	for k, v := range e.tradeTypeMap {
		e.tradeTypeMapReverse[v] = k
	}
	e.client = sdk.NewClient("uri", "auth")
	return nil
}

// SetStockTypeMap ...
func (e *FutureBackTest) SetStockTypeMap(m map[string]goex.CurrencyPair) {
	e.stockTypeMap = m
}

// GetStockTypeMap ...
func (e *FutureBackTest) GetStockTypeMap() map[string]goex.CurrencyPair {
	return e.stockTypeMap
}

// Log print something to console
func (e *FutureBackTest) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, e.GetStockType(), 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *FutureBackTest) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *FutureBackTest) GetName() string {
	return e.option.Name
}

// GetDepth ...
func (e *FutureBackTest) GetDepth(size int) interface{} {
	var resDepth constant.Depth
	var opt sdkType.Option
	opt.Market = e.GetName()
	opt.Symbol = e.GetStockType()
	response := e.client.GetDepth(opt)
	if !response.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetDepth() error, the error number is ", response.Message)
		return false
	}
	askList := response.Data.Asks
	bidList := response.Data.Bids
	for _, ask := range askList {
		var resAsk constant.DepthRecord
		resAsk.Amount = ask.Amount
		resAsk.Price = ask.Price
		resDepth.Asks = append(resDepth.Asks, resAsk)
	}
	for _, bid := range bidList {
		var resBid constant.DepthRecord
		resBid.Amount = bid.Amount
		resBid.Price = bid.Price
		resDepth.Bids = append(resDepth.Bids, resBid)
	}
	return resDepth
}

// GetPosition ...
func (e *FutureBackTest) GetPosition() interface{} {
			resPositionVec := []constant.Position{}
			stockType := e.GetStockType()
			exchangeStockType, ok := e.stockTypeMap[stockType]
			if !ok {
				e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetPosition() error, the error number is stockType")
				return false
			}
			positions, err := e.api.GetFuturePosition(exchangeStockType, e.GetContractType())
			if err != nil {
				e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetPosition() error, the error number is ", err.Error())
				return false
			}
			for _, position := range positions {
				var resPosition constant.Position
				if position.BuyAmount > 0 {
					resPosition.Price = position.BuyPriceAvg
					resPosition.Amount = position.BuyAmount
					resPosition.MarginLevel = position.LeverRate
					resPosition.Profit = position.BuyProfitReal
					resPosition.ForcePrice = position.ForceLiquPrice
					resPosition.TradeType = constant.TradeTypeBuy
					resPosition.ContractType = position.ContractType
					resPosition.StockType = position.Symbol.CurrencyA.Symbol + "/" + position.Symbol.CurrencyB.Symbol
					resPositionVec = append(resPositionVec, resPosition)
				}
				if position.SellAmount > 0 {
					resPosition.Price = position.SellPriceAvg
					resPosition.Amount = position.SellAmount
					resPosition.MarginLevel = position.LeverRate
					resPosition.ForcePrice = position.ForceLiquPrice
					resPosition.TradeType = constant.TradeTypeSell
					resPosition.ContractType = e.contractType
					resPosition.StockType = position.Symbol.CurrencyA.Symbol + "/" + position.Symbol.CurrencyB.Symbol
					resPositionVec = append(resPositionVec, resPosition)
				}
			}
			return resPositionVec
		}

		// SetLimit set the limit calls amount per second of this exchange
		func (e *FutureBackTest) SetLimit(times interface{}) float64 {
			e.limit = util.Float64Must(times)
			return e.limit
		}

		// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
		func (e *FutureBackTest) AutoSleep() {
			now := time.Now().UnixNano()
			interval := 1e+9/e.limit*util.Float64Must(e.lastTimes) - util.Float64Must(now-e.lastSleep)
			if interval > 0.0 {
				time.Sleep(time.Duration(util.Int64Must(interval)))
			}
			e.lastTimes = 0
			e.lastSleep = now

	return false
}

// GetMinAmount get the min trade amount of this exchange
func (e *FutureBackTest) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

// GetAccount get the account detail of this exchange
func (e *FutureBackTest) GetAccount() interface{} {
		account, err := e.api.GetFutureUserinfo()
		if err != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetAccount() error, the error number is ", err.Error())
			return false
		}
		var resAccount constant.Account
		resAccount.SubAccounts = make(map[string]constant.SubAccount)
		for _, v := range account.FutureSubAccounts {
			var subAccount constant.SubAccount
			stockType := v.Currency.Symbol
			subAccount.AccountRights = v.AccountRights
			subAccount.KeepDeposit = v.KeepDeposit
			subAccount.ProfitReal = v.ProfitReal
			subAccount.ProfitUnreal = v.ProfitUnreal
			subAccount.RiskRate = v.RiskRate
			resAccount.SubAccounts[stockType] = subAccount
		}
		return resAccount

	return false
}

func (e *FutureBackTest) Buy(price, amount string, msg ...interface{}) interface{} {
		var err error
		var openType int
		stockType := e.GetStockType()
		exchangeStockType, ok := e.stockTypeMap[stockType]
		if !ok {
			e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Buy() error, the error number is stockType")
			return false
		}
		if err := e.ValidBuy(); err != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Buy() error, the error number is ", err.Error())
			return false
		}
		level := e.GetMarginLevel()
		var matchPrice = 0
		if price == "-1" {
			matchPrice = 1
		}
		openType = e.tradeTypeMapReverse[e.GetDirection()]
		orderId, err := e.api.PlaceFutureOrder(exchangeStockType, e.GetContractType(),
			price, amount, openType, matchPrice, level)

		if err != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Sell() error, the error number is ", err.Error())
			return false
		}
		priceFloat := util.Float64Must(price)
		amountFloat := util.Float64Must(amount)
		e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
		return orderId

	return false
}

func (e *FutureBackTest) Sell(price, amount string, msg ...interface{}) interface{} {
			var err error
			var openType int
			stockType := e.GetStockType()
			exchangeStockType, ok := e.stockTypeMap[stockType]
			if !ok {
				e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Sell() error, the error number is stockType")
				return false
			}
			if err := e.ValidSell(); err != nil {
				e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Sell() error, the error number is ", err.Error())
				return false
			}
			level := e.GetMarginLevel()
			var matchPrice = 0
			if price == "-1" {
				matchPrice = 1
			}
			openType = e.tradeTypeMapReverse[e.GetDirection()]
			orderId, err := e.api.PlaceFutureOrder(exchangeStockType, e.GetContractType(),
				price, amount, openType, matchPrice, level)

			if err != nil {
				e.logger.Log(constant.ERROR, e.GetStockType(), util.Float64Must(amount), util.Float64Must(amount), "Sell() error, the error number is ", err.Error())
				return false
			}
			priceFloat := util.Float64Must(price)
			amountFloat := util.Float64Must(amount)
			e.logger.Log(e.direction, stockType, priceFloat, amountFloat, msg...)
			return orderId
		}

		// GetOrder get details of an order
		func (e *FutureBackTest) GetOrder(id string) interface{} {
			exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
			if !ok {
				e.logger.Log(constant.ERROR, e.GetStockType(), 0, util.Float64Must(id), "GetOrder() error, the error number is stockType")
				return false
			}
			orders, err := e.api.GetUnfinishFutureOrders(exchangeStockType, e.GetContractType())
			if err != nil {
				e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, util.Float64Must(id), "GetOrder() error, the error number is ", err.Error())
				return false
			}
			for _, order := range orders {
				if id != order.OrderID2 {
					continue
				}
				return constant.Order{
					Id:         order.OrderID2,
					Price:      order.Price,
					Amount:     order.Amount,
					DealAmount: order.DealAmount,
					TradeType:  e.tradeTypeMap[order.OrderType],
					StockType:  e.GetStockType(),
				}
			}
			return false

	return false
}

// GetOrders get all unfilled orders
func (e *FutureBackTest) GetOrders() interface{} {
			exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
			if !ok {
				e.logger.Log(constant.ERROR, "", 0, 0, "GetOrders() error, the error number is stockType")
				return false
			}
			orders, err := e.api.GetUnfinishFutureOrders(exchangeStockType, e.GetStockType())
			if err != nil {
				e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetOrders() error, the error number is ", err.Error())
				return false
			}
			resOrders := []constant.Order{}
			for _, order := range orders {
				resOrder := constant.Order{
					Id:         order.OrderID2,
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
		func (e *FutureBackTest) GetTrades(params ...interface{}) interface{} {
			var traders []constant.Trader
			exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
			if !ok {
				e.logger.Log(constant.ERROR, e.GetStockType(), 0, 0, "GetTrades() error, the error number is stockType")
				return false
			}
			APITraders, err := e.api.GetTrades(e.GetContractType(), exchangeStockType, 0)
			if err != nil {
				return false
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

	return false
}

// CancelOrder cancel an order
func (e *FutureBackTest) CancelOrder(orderID string) bool {
		exchangeStockType, ok := e.stockTypeMap[e.GetStockType()]
		if !ok {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0, util.Float64Must(orderID), "CancelOrder() error, the error number is stockType")
			return false
		}
		result, err := e.api.FutureCancelOrder(exchangeStockType, e.GetContractType(), orderID)
		if err != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, util.Float64Must(orderID), "CancelOrder() error, the error number is ", err.Error())
			return false
		}
		if !result {
			return false
		}
		e.logger.Log(constant.TradeTypeCancel, e.GetStockType(), 0, util.Float64Must(orderID), "CancelOrder() success")
		return true

	return false
}

// GetTicker get market ticker
func (e *FutureBackTest) GetTicker() interface{} {
	if e.backTestCurr > e.backTestEnd {
		//reload the ticker
		var opt sdkType.Option
		opt.BeginTime = e.backTestBegin
		opt.EndTime = e.backTestEnd
		opt.Symbol = e.GetStockType()
		opt.Market = e.GetName()
		response := e.client.GetOHLCs(opt)
		if !response.Success {
			return false
		}
		data := response.Data
		if len(data) <= 0 {
			return false
		}
		e.backTestBegin = data[0].Time
		e.backTestEnd = data[len(data)-1].Time
		e.backTestIndex = 0
		e.backTestTickers = data
	}
	backTestTicker := e.backTestTickers[e.backTestIndex]
	ticker := constant.Ticker{
		Last: backTestTicker.Close,
		Buy:  backTestTicker.Open,
		Sell: backTestTicker.Close,
		High: backTestTicker.High,
		Low:  backTestTicker.Low,
		Vol:  backTestTicker.Volume,
		Time: backTestTicker.Time,
	}
	return ticker
}

// GetRecords get candlestick data
// params[0] period
// params[1] size
// params[2] since
func (e *FutureBackTest) GetRecords(params ...interface{}) interface{} {
	var period int64 = -1
	var size = 0
	var periodStr = "M15"

	if len(params) >= 1 && util.StringMust(params[0]) != "" {
		periodStr = util.StringMust(params[0])
	}

	period, ok := e.recordsPeriodMap[periodStr]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, 0, "GetRecords() error, the error number is stockType")
		return false
	}

	if len(params) >= 2 && util.IntMust(params[1]) > 0 {
		size = util.IntMust(params[1])
	}

	var opt sdkType.Option
	opt.BeginTime = e.backTestBegin
	opt.EndTime = e.backTestEnd
	opt.Symbol = e.GetStockType()
	opt.Market = e.GetName()
	opt.Period = period
	response := e.client.GetOHLCs(opt)
	if !response.Success {
		return false
	}
	klineVec := response.Data
	if len(klineVec) <= 0 {
		return false
	}
	timeLast := int64(0)
	if len(e.records[periodStr]) > 0 {
		timeLast = e.records[periodStr][len(e.records[periodStr])-1].Time
	}
	var recordsNew []constant.Record
	for i := len(klineVec); i > 0; i-- {
		kline := klineVec[i-1]
		recordTime := kline.Time
		if recordTime > timeLast {
			recordsNew = append([]constant.Record{{
				Time:   recordTime,
				Open:   kline.Open,
				High:   kline.High,
				Low:    kline.Low,
				Close:  kline.Close,
				Volume: kline.Volume,
			}}, recordsNew...)
		} else if timeLast > 0 && recordTime == timeLast {
			e.records[periodStr][len(e.records[periodStr])-1] = constant.Record{
				Time:   recordTime,
				Open:   kline.Open,
				High:   kline.High,
				Low:    kline.Low,
				Close:  kline.Close,
				Volume: kline.Volume,
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
*/
