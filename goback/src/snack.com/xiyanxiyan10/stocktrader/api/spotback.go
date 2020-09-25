package api

import (
	"errors"
	"math"
	"sync"
	"time"

	goex "github.com/nntaoli-project/goex"
	dbtypes "snack.com/xiyanxiyan10/stockdb/types"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

// ExchangeBackConfig ...
type ExchangeBackConfig struct {
	ExName               string
	TakerFee             float64
	MakerFee             float64
	SupportCurrencyPairs []string
	QuoteCurrency        string //净值币种
	Account              constant.Account
	BackTestStartTime    int64
	BackTestEndTime      int64
	DepthSize            int64 //回测多少档深度
	UnGzip               bool  //是否解压
}

// ExchangeBack ...
type ExchangeBack struct {
	BaseExchange
	*sync.RWMutex
	acc                  *constant.Account
	name                 string
	makerFee             float64
	takerFee             float64
	supportCurrencyPairs []string
	quoteCurrency        string
	pendingOrders        map[string]*constant.Order
	finishedOrders       map[string]*constant.Order
	dataLoader           map[string]*DataLoader
	currData             dbtypes.OHLC
	idGen                *util.IDGen
	contractRate         float64 // 合约每张价值
	CurrencyStandard     bool    // 是否为币本位
	sortedCurrencies     constant.Account
	longPosition         map[string]constant.Position // 多仓
	shortPosition        map[string]constant.Position // 空仓
}

// NewExchangeBack ...
func NewExchangeBack(config ExchangeBackConfig) *ExchangeBack {
	sim := &ExchangeBack{
		RWMutex:              new(sync.RWMutex),
		idGen:                util.NewIDGen(config.ExName),
		name:                 config.ExName,
		makerFee:             config.MakerFee,
		takerFee:             config.TakerFee,
		acc:                  &config.Account,
		supportCurrencyPairs: config.SupportCurrencyPairs,
		quoteCurrency:        config.QuoteCurrency,
		pendingOrders:        make(map[string]*constant.Order, 100),
		finishedOrders:       make(map[string]*constant.Order, 100),
		dataLoader:           make(map[string]*DataLoader, 1),
		longPosition:         make(map[string]constant.Position, 1),
		shortPosition:        make(map[string]constant.Position, 1),
	}

	for key, sub := range sim.acc.SubAccounts {
		sim.sortedCurrencies.SubAccounts[key] = sub
	}
	sim.back = true
	return sim
}

// Ready ...
func (e *ExchangeBack) Ready() interface{} {
	var account constant.Account
	e.RWMutex = new(sync.RWMutex)
	e.idGen = util.NewIDGen(e.GetExchangeName())
	e.name = e.GetExchangeName()
	e.makerFee = e.BaseExchange.maker
	e.takerFee = e.BaseExchange.taker
	e.acc = &account

	e.pendingOrders = make(map[string]*constant.Order, 100)
	e.finishedOrders = make(map[string]*constant.Order, 100)
	e.dataLoader = make(map[string]*DataLoader, 1)
	e.longPosition = make(map[string]constant.Position, 1)
	e.shortPosition = make(map[string]constant.Position, 1)
	for stock := range e.BaseExchange.subscribeMap {
		var loader DataLoader
		e.dataLoader[stock] = &loader
		val := e.BaseExchange.BackGetOHLCs(e.BaseExchange.start, e.BaseExchange.end, e.BaseExchange.period)
		if val == nil {
			return nil
		}
		ohlcs := val.([]dbtypes.OHLC)
		e.dataLoader[stock].Load(ohlcs)
	}
	currencyMap := e.BaseExchange.currencyMap
	for key, val := range currencyMap {
		var sub constant.SubAccount
		sub.Amount = val
		e.acc.SubAccounts[key] = sub
	}
	return "success"
}

func (ex *ExchangeBack) fillOrder(isTaker bool, amount, price float64, ord *constant.Order) {
	ord.FinishedTime = ex.currData.Time / int64(time.Millisecond) //set filled time
	dealAmount := 0.0
	remain := ord.Amount - ord.DealAmount
	if remain > amount {
		dealAmount = amount
	} else {
		dealAmount = remain
	}

	ratio := dealAmount / (ord.DealAmount + dealAmount)
	ord.AvgPrice = math.Round(ratio*price+(1-ratio)*ord.AvgPrice*100000000) / 100000000
	ord.DealAmount += dealAmount
	if ord.Amount == ord.DealAmount {
		ord.Status = constant.ORDER_FINISH
	} else {
		if ord.DealAmount > 0 {
			ord.Status = constant.ORDER_PART_FINISH
		}
	}

	fee := ex.makerFee
	if isTaker {
		fee = ex.takerFee
	}

	tradeFee := 0.0
	switch ord.TradeType {
	case constant.TradeTypeSell:
		tradeFee = dealAmount * price * fee
	case constant.TradeTypeBuy:
		tradeFee = dealAmount * fee
	}
	tradeFee = math.Floor(tradeFee*100000000) / 100000000

	ord.Fee += tradeFee

	ex.unFrozenAsset(tradeFee, dealAmount, price, *ord)
}

func (ex *ExchangeBack) matchOrder(ord *constant.Order, isTaker bool) {
	ticker := ex.currData
	switch ord.TradeType {
	case constant.TradeTypeSell:
		if ticker.Close >= ord.Price && ticker.Volume > 0 {
			ex.fillOrder(isTaker, ticker.Volume, ticker.Close, ord)
			if ord.Status == constant.ORDER_FINISH {
				delete(ex.pendingOrders, ord.Id)
				ex.finishedOrders[ord.Id] = ord
			}
		}
	case constant.TradeTypeBuy:
		if ticker.Close <= ord.Price && ticker.Volume > 0 {
			ex.fillOrder(isTaker, ticker.Volume, ticker.Close, ord)
			if ord.Status == constant.ORDER_FINISH {
				delete(ex.pendingOrders, ord.Id)
				ex.finishedOrders[ord.Id] = ord
			}
		}
	}
}

func (ex *ExchangeBack) match() {
	ex.Lock()
	defer ex.Unlock()
	for id := range ex.pendingOrders {
		ex.matchOrder(ex.pendingOrders[id], false)
	}
}

// LimitBuy ...
func (ex *ExchangeBack) LimitBuy(amount, price, currency string) (*constant.Order, error) {
	ex.Lock()
	defer ex.Unlock()

	ord := constant.Order{
		Price:     goex.ToFloat64(price),
		Amount:    goex.ToFloat64(amount),
		Id:        ex.idGen.Get(),
		Time:      ex.currData.Time / int64(time.Millisecond),
		Status:    constant.ORDER_UNFINISH,
		StockType: currency,
		TradeType: constant.TradeTypeBuy,
	}

	err := ex.frozenAsset(ord)
	if err != nil {
		return nil, err
	}

	ex.pendingOrders[ord.Id] = &ord

	ex.matchOrder(&ord, true)

	var result constant.Order
	util.DeepCopyStruct(ord, &result)
	return &result, nil
}

// LimitSell ...
func (ex *ExchangeBack) LimitSell(amount, price, currency string) (*constant.Order, error) {
	ex.Lock()
	defer ex.Unlock()

	ord := constant.Order{
		Price:     goex.ToFloat64(price),
		Amount:    goex.ToFloat64(amount),
		Id:        ex.idGen.Get(),
		Time:      ex.currData.Time / int64(time.Millisecond),
		Status:    constant.ORDER_UNFINISH,
		StockType: currency,
		TradeType: constant.TradeTypeSell,
	}

	err := ex.frozenAsset(ord)
	if err != nil {
		return nil, err
	}

	ex.pendingOrders[ord.Id] = &ord

	ex.matchOrder(&ord, true)

	var result constant.Order
	util.DeepCopyStruct(ord, &result)

	return &result, nil
}

// MarketBuy ...
func (ex *ExchangeBack) MarketBuy(amount, price, currency string) (*constant.Order, error) {
	panic("not support")
}

// MarketSell ...
func (ex *ExchangeBack) MarketSell(amount, price, currency string) (*constant.Order, error) {
	panic("not support")
}

// CancelOrder ...
func (ex *ExchangeBack) CancelOrder(orderID string, currency string) (bool, error) {
	ex.Lock()
	defer ex.Unlock()

	ord := ex.finishedOrders[orderID]
	if ord != nil {
		return false, ErrCancelOrderFinished
	}

	ord = ex.pendingOrders[orderID]
	if ord == nil {
		return false, ErrNotFoundOrder
	}

	delete(ex.pendingOrders, ord.Id)

	ord.Status = constant.ORDER_CANCEL
	ex.finishedOrders[ord.Id] = ord

	ex.unFrozenAsset(0, 0, 0, *ord)

	return true, nil
}

// GetOneOrder ...
func (ex *ExchangeBack) GetOneOrder(orderID, currency string) (*constant.Order, error) {
	ex.RLock()
	defer ex.RUnlock()

	ord := ex.finishedOrders[orderID]
	if ord == nil {
		ord = ex.pendingOrders[orderID]
	}

	if ord != nil {
		// deep copy
		var result constant.Order
		util.DeepCopyStruct(ord, &result)

		return &result, nil
	}

	return nil, ErrNotFoundOrder
}

// GetUnfinishOrders ...
func (ex *ExchangeBack) GetUnfinishOrders(currency string) ([]constant.Order, error) {
	ex.RLock()
	defer ex.RUnlock()

	var unfinishedOrders []constant.Order
	for _, ord := range ex.pendingOrders {
		unfinishedOrders = append(unfinishedOrders, *ord)
	}

	return unfinishedOrders, nil
}

// GetOrderHistorys ...
func (ex *ExchangeBack) GetOrderHistorys(currency string, currentPage, pageSize int) ([]constant.Order, error) {
	ex.RLock()
	defer ex.RUnlock()

	var orders []constant.Order
	for _, ord := range ex.finishedOrders {
		if ord.StockType == currency {
			orders = append(orders, *ord)
		}
	}
	return orders, nil
}

// GetAccount ...
func (ex *ExchangeBack) GetAccount() (*constant.Account, error) {
	ex.RLock()
	defer ex.RUnlock()
	var account constant.Account
	account.SubAccounts = make(map[string]constant.SubAccount)
	for key, sub := range ex.acc.SubAccounts {
		account.SubAccounts[key] = sub
	}
	return &account, nil
}

// GetTicker ...
func (ex *ExchangeBack) GetTicker(currency string) (*constant.Ticker, error) {
	loader := ex.dataLoader[currency]
	if loader == nil {
		return nil, errors.New("loader not found")
	}
	ohlc := loader.Next()
	if ohlc == nil {
		return nil, ErrDataFinished
	}
	ex.currData = *ohlc
	ex.match()
	return &constant.Ticker{
		Last: ohlc.Close,
		Buy:  ohlc.Close,
		Sell: ohlc.Close,
		High: ohlc.High,
		Low:  ohlc.Low,
	}, nil
}

// GetDepth ...
func (ex *ExchangeBack) GetDepth(size int, currency string) (*constant.Depth, error) {
	val := ex.BaseExchange.BackGetDepth(ex.currData.Time, ex.currData.Time, ex.currData.Time)
	if val == nil {
		return nil, errors.New("Get depth fail")
	}
	var depth constant.Depth
	dbdepth := val.(dbtypes.Depth)
	for _, ask := range dbdepth.Asks {
		var record constant.DepthRecord
		record.Amount = ask.Amount
		record.Price = ask.Price
		depth.Asks = append(depth.Asks, record)
	}

	for _, bid := range dbdepth.Bids {
		var record constant.DepthRecord
		record.Amount = bid.Amount
		record.Price = bid.Price
		depth.Bids = append(depth.Bids, record)
	}
	return &depth, nil
}

// GetExchangeName ...
func (ex *ExchangeBack) GetExchangeName() string {
	return ex.name
}

//冻结
func (ex *ExchangeBack) frozenAsset(order constant.Order) error {
	stocks := stockPair2Vec(order.StockType)
	CurrencyA := stocks[0]
	CurrencyB := stocks[1]
	switch order.TradeType {
	case constant.TradeTypeSell:
		avaAmount := ex.acc.SubAccounts[order.StockType].Amount
		if avaAmount < order.Amount {
			return ErrDataInsufficient
		}
		ex.acc.SubAccounts[CurrencyA] = constant.SubAccount{
			StockType:    CurrencyA,
			Amount:       avaAmount - order.Amount,
			FrozenAmount: ex.acc.SubAccounts[CurrencyA].FrozenAmount + order.Amount,
			LoanAmount:   0,
		}
	case constant.TradeTypeBuy:
		avaAmount := ex.acc.SubAccounts[CurrencyB].Amount
		need := order.Amount * order.Price
		if avaAmount < need {
			return ErrDataInsufficient
		}
		ex.acc.SubAccounts[CurrencyB] = constant.SubAccount{
			StockType:    CurrencyB,
			Amount:       avaAmount - need,
			FrozenAmount: ex.acc.SubAccounts[CurrencyB].FrozenAmount + need,
			LoanAmount:   0,
		}
	}
	return nil
}

// unFrozenAsset 解冻
func (ex *ExchangeBack) unFrozenAsset(fee, matchAmount, matchPrice float64, order constant.Order) {
	stocks := stockPair2Vec(order.StockType)
	CurrencyA := stocks[0]
	CurrencyB := stocks[1]
	assetA := ex.acc.SubAccounts[CurrencyA]
	assetB := ex.acc.SubAccounts[CurrencyB]

	switch order.TradeType {
	case constant.TradeTypeSell:
		if order.Status == constant.ORDER_CANCEL {
			ex.acc.SubAccounts[assetA.StockType] = constant.SubAccount{
				StockType:    assetA.StockType,
				Amount:       assetA.Amount + order.Amount - order.DealAmount,
				FrozenAmount: assetA.FrozenAmount - (order.Amount - order.DealAmount),
				LoanAmount:   0,
			}
		} else {
			ex.acc.SubAccounts[assetA.StockType] = constant.SubAccount{
				StockType:    assetA.StockType,
				Amount:       assetA.Amount,
				FrozenAmount: assetA.FrozenAmount - matchAmount,
				LoanAmount:   0,
			}
			ex.acc.SubAccounts[assetB.StockType] = constant.SubAccount{
				StockType:    assetB.StockType,
				Amount:       assetB.Amount + matchAmount*matchPrice - fee,
				FrozenAmount: assetB.FrozenAmount,
			}
		}

	case constant.TradeTypeBuy:
		if order.Status == constant.ORDER_CANCEL {
			unFrozen := (order.Amount - order.DealAmount) * order.Price
			ex.acc.SubAccounts[assetB.StockType] = constant.SubAccount{
				StockType:    assetB.StockType,
				Amount:       assetB.Amount + unFrozen,
				FrozenAmount: assetB.FrozenAmount - unFrozen,
			}
		} else {
			ex.acc.SubAccounts[assetA.StockType] = constant.SubAccount{
				StockType:    assetA.StockType,
				Amount:       assetA.Amount + matchAmount - fee,
				FrozenAmount: assetA.FrozenAmount,
				LoanAmount:   0,
			}
			ex.acc.SubAccounts[assetB.StockType] = constant.SubAccount{
				StockType:    assetB.StockType,
				Amount:       assetB.Amount + matchAmount*(order.Price-matchPrice),
				FrozenAmount: assetB.FrozenAmount - matchAmount*order.Price,
			}
		}
	}

}

// GetRecords get candlestick data
func (e *ExchangeBack) GetRecords(periodStr, maStr string) interface{} {
	var period int64 = -1
	var size = constant.RecordSize
	period, ok := e.recordsPeriodMap[periodStr]
	if !ok {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0, 0, "GetRecords() error, the error number is stockType")
		return nil
	}

	val := e.BaseExchange.BackGetOHLCs(e.currData.Time, e.BaseExchange.end, period)
	if val != nil {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0, "GetRecords() error")
		return nil
	}
	vec := val.([]dbtypes.OHLC)

	if len(vec) > size {
		vec = vec[0 : size-1]
	}
	var records []constant.Record
	for _, kline := range vec {
		records = append([]constant.Record{{
			Open:   kline.Open,
			High:   kline.High,
			Low:    kline.Low,
			Close:  kline.Close,
			Volume: kline.Volume,
		}}, records...)
	}
	return records
}
