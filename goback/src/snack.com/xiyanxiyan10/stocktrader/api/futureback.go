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

// ExchangeFutureBackConfig ...
type ExchangeFutureBackConfig struct {
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

// ExchangeFutureBack ...
type ExchangeFutureBack struct {
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
	sortedCurrencies     constant.Account
	longPosition         map[string]constant.Position // 多仓
	shortPosition        map[string]constant.Position // 空仓
}

// NewExchangeFutureBack ...
func NewExchangeFutureBack(config ExchangeBackConfig) *ExchangeFutureBack {
	sim := &ExchangeFutureBack{
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
	return sim
}

func (ex *ExchangeFutureBack) fillOrder(isTaker bool, amount, price float64, ord *constant.Order) {
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

func (ex *ExchangeFutureBack) matchOrder(ord *constant.Order, isTaker bool) {
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

func (ex *ExchangeFutureBack) match() {
	ex.Lock()
	defer ex.Unlock()
	for id := range ex.pendingOrders {
		ex.matchOrder(ex.pendingOrders[id], false)
	}
}

// LimitBuy ...
func (ex *ExchangeFutureBack) LimitBuy(amount, price, currency string) (*constant.Order, error) {
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
func (ex *ExchangeFutureBack) LimitSell(amount, price, currency string) (*constant.Order, error) {
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
func (ex *ExchangeFutureBack) MarketBuy(amount, price, currency string) (*constant.Order, error) {
	panic("not support")
}

// MarketSell ...
func (ex *ExchangeFutureBack) MarketSell(amount, price, currency string) (*constant.Order, error) {
	panic("not support")
}

// CancelOrder ...
func (ex *ExchangeFutureBack) CancelOrder(orderID string, currency string) (bool, error) {
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
func (ex *ExchangeFutureBack) GetOneOrder(orderID, currency string) (*constant.Order, error) {
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
func (ex *ExchangeFutureBack) GetUnfinishOrders(currency string) ([]constant.Order, error) {
	ex.RLock()
	defer ex.RUnlock()

	var unfinishedOrders []constant.Order
	for _, ord := range ex.pendingOrders {
		unfinishedOrders = append(unfinishedOrders, *ord)
	}

	return unfinishedOrders, nil
}

// GetOrderHistorys ...
func (ex *ExchangeFutureBack) GetOrderHistorys(currency string, currentPage, pageSize int) ([]constant.Order, error) {
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
func (ex *ExchangeFutureBack) GetAccount() (*constant.Account, error) {
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
func (ex *ExchangeFutureBack) GetTicker(currency string) (*constant.Ticker, error) {
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
func (ex *ExchangeFutureBack) GetDepth(size int, currency string) (*constant.Depth, error) {
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

// GetTrades ...
func (ex *ExchangeFutureBack) GetTrades(currencyPair string, since int64) ([]constant.Trader, error) {
	panic("not support")
}

// GetExchangeName ...
func (ex *ExchangeFutureBack) GetExchangeName() string {
	return ex.name
}

//冻结
func (ex *ExchangeFutureBack) frozenAsset(order constant.Order) error {
	stocks := stockPair2Vec(order.StockType)
	CurrencyA := stocks[0]
	ticker := ex.currData
	price := ticker.Close
	avaAmount := ex.acc.SubAccounts[CurrencyA].Amount
	lever := ex.BaseExchange.lever
	switch order.TradeType {
	case constant.TradeTypeLong, constant.TradeTypeShort:
		if avaAmount*lever*price < order.Amount*ex.BaseExchange.maintenanceRate {
			return ErrDataInsufficient
		}
		costAmount := (order.Amount * ex.BaseExchange.maintenanceRate) / (lever * order.Price)
		ex.acc.SubAccounts[CurrencyA] = constant.SubAccount{
			StockType:    CurrencyA,
			Amount:       avaAmount - costAmount,
			ForzenAmount: ex.acc.SubAccounts[CurrencyA].ForzenAmount + costAmount,
			LoanAmount:   0,
		}
	case constant.TradeTypeLongClose, constant.TradeTypeShortClose:
		if order.TradeType == constant.TradeTypeLongClose {
			position := ex.longPosition[CurrencyA]
			if position.Amount < order.Amount {
				return ErrDataInsufficient
			}
			position.Amount = position.Amount - order.Amount
			position.ForzenAmount = position.ForzenAmount + order.Amount
			ex.longPosition[CurrencyA] = position
		}

		if order.TradeType == constant.TradeTypeShortClose {
			position := ex.shortPosition[CurrencyA]
			if position.Amount < order.Amount {
				return ErrDataInsufficient
			}
			position.Amount = position.Amount - order.Amount
			position.ForzenAmount = position.ForzenAmount + order.Amount
			ex.shortPosition[CurrencyA] = position
		}
	}
	return nil
}

//解冻
func (ex *ExchangeFutureBack) unFrozenAsset(fee, matchAmount, matchPrice float64, order constant.Order) {
	stocks := stockPair2Vec(order.StockType)
	CurrencyA := stocks[0]
	assetA := ex.acc.SubAccounts[CurrencyA]
	lever := ex.BaseExchange.lever
	switch order.TradeType {
	case constant.TradeTypeLong, constant.TradeTypeShort:
		if order.Status == constant.ORDER_CANCEL {
			costAmount := (order.Amount * ex.BaseExchange.maintenanceRate) / (lever * order.Price)
			ex.acc.SubAccounts[assetA.StockType] = constant.SubAccount{
				StockType:    assetA.StockType,
				Amount:       assetA.Amount + costAmount - order.DealAmount,
				ForzenAmount: assetA.ForzenAmount - (costAmount - order.DealAmount),
				LoanAmount:   0,
			}
		}
	case constant.TradeTypeLongClose, constant.TradeTypeShortClose:
		if order.Status == constant.ORDER_CANCEL {
			if order.TradeType == constant.TradeTypeLongClose {
				position := ex.longPosition[CurrencyA]
				position.Amount = position.Amount + order.Amount
				position.ForzenAmount = position.ForzenAmount - order.Amount
				ex.longPosition[CurrencyA] = position
			}

			if order.TradeType == constant.TradeTypeShortClose {
				position := ex.shortPosition[CurrencyA]
				position.Amount = position.Amount + order.Amount
				position.ForzenAmount = position.ForzenAmount - order.Amount
				ex.shortPosition[CurrencyA] = position
			}
		}
	}
}
