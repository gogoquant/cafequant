package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	goex "github.com/nntaoli-project/goex"
	log "github.com/sirupsen/logrus"
	"github.com/zhnxin/csvreader"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"

	//"strconv"
	"sync"
	"time"
)

// ExchangeFutureBackConfig ...
type ExchangeFutureBackConfig struct {
	ExName               string
	TakerFee             float64
	MakerFee             float64
	SupportCurrencyPairs []string
	QuoteCurrency        string //净值币种
	Account              constant.Account
	DepthSize            int64 //回测多少档深度
	UnGzip               bool  //是否解压
}

// ExchangeFutureBack ...
type ExchangeFutureBack struct {
	BaseExchange
	progress int
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
	stockTypeMap         map[string]goex.CurrencyPair
	currData             map[string]constant.OHLC
	idGen                *util.IDGen
	sortedCurrencies     constant.Account
	longPosition         map[string]constant.Position // 多仓
	shortPosition        map[string]constant.Position // 空仓

	recordsMap   map[string]map[int64]int
	recordsCache map[string][]constant.Record //records store as cache
}

// NewExchangeFutureBack2Config ...
func NewExchangeFutureBack2Config(config ExchangeBackConfig) *ExchangeFutureBack {
	sim := &ExchangeFutureBack{
		RWMutex:              new(sync.RWMutex),
		idGen:                util.NewIDGen(config.ExName),
		name:                 config.ExName,
		makerFee:             config.MakerFee,
		takerFee:             config.TakerFee,
		acc:                  &config.Account,
		supportCurrencyPairs: config.SupportCurrencyPairs,
		quoteCurrency:        config.QuoteCurrency,
		pendingOrders:        make(map[string]*constant.Order, 0),
		finishedOrders:       make(map[string]*constant.Order, 0),
		dataLoader:           make(map[string]*DataLoader, 0),
		longPosition:         make(map[string]constant.Position, 0),
		shortPosition:        make(map[string]constant.Position, 0),
	}

	for key, sub := range sim.acc.SubAccounts {
		sim.sortedCurrencies.SubAccounts[key] = sub
	}
	return sim
}

// NewExchangeFutureBack ...
func NewExchangeFutureBack(config ExchangeBackConfig) *ExchangeFutureBack {
	sim := &ExchangeFutureBack{}
	return sim
}

// Debug ..,
func (e *ExchangeFutureBack) Debug() error {
	log.Infof("---FutureBack info start---\n")
	log.Infof("currTicker:\n")
	v, err := json.Marshal(e.currData[e.stockType])
	if err != nil {
		log.Infof("utilt ticker err :%s\n", err.Error())
		return err
	}
	log.Infof("%s\n", string(v))
	marginRatio, lft, rht := e.marginRatio()
	log.Infof("marginRatio %f lft %f  rht %f\n", marginRatio, lft, rht)
	log.Infof("longPosition:\n")
	v, err = json.Marshal(e.longPosition)
	if err != nil {
		log.Infof("utilt longPosition err :%s\n", err.Error())
		return err
	}
	log.Infof("%s\n", string(v))
	log.Infof("shortPosition:\n")
	v, err = json.Marshal(e.shortPosition)
	if err != nil {
		log.Infof("utilt shortPosition err :%s\n", err.Error())
		return err
	}
	log.Infof("%s\n", string(v))
	log.Infof("account:\n")
	if e.acc != nil {
		v, err = json.Marshal(e.acc)
		if err != nil {
			log.Infof("utilt account err :%s\n", err.Error())
			return err
		}
		log.Infof("%s\n", string(v))
	}
	log.Infof("pendingOrders:\n")
	v, err = json.Marshal(e.pendingOrders)
	if err != nil {
		log.Infof("utilt pending orders err :%s\n", err.Error())
		return err
	}
	log.Infof("%s\n", string(v))
	log.Infof("---FutureBack info end---\n")
	return nil
}

func isContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func (e *ExchangeFutureBack) Stop() error {
	return nil
}

// Start ...
func (e *ExchangeFutureBack) Start() error {
	var account constant.Account
	e.RWMutex = new(sync.RWMutex)
	e.idGen = util.NewIDGen(e.GetExchangeName())
	e.name = e.GetExchangeName()
	e.makerFee = e.BaseExchange.maker
	e.takerFee = e.BaseExchange.taker
	e.acc = &account
	e.acc.SubAccounts = make(map[string]constant.SubAccount)
	e.pendingOrders = make(map[string]*constant.Order, 0)
	e.finishedOrders = make(map[string]*constant.Order, 0)
	e.dataLoader = make(map[string]*DataLoader, 0)
	e.longPosition = make(map[string]constant.Position, 0)
	e.shortPosition = make(map[string]constant.Position, 0)

	//@ load ohlc here
	historyDir := config.String("history")
	for _, name := range e.option.WatchList {
		dataPath := historyDir + "/" + strings.Replace((e.GetExchangeName()+name), "/", ".", -1) + ".csv"

		var ohlcs []constant.OHLC
		err := csvreader.New().UnMarshalFile(dataPath, &ohlcs)
		if err != nil {
			log.Errorf("Load data from %s to %s error %s", dataPath, name, err.Error())
			return err
		} else {
			log.Infof("Load data from %s to %s success", dataPath, name)
		}
		e.dataLoader[name] = new(DataLoader)
		e.dataLoader[name].Load(ohlcs)
	}
	currencyMap := e.BaseExchange.currencyMap
	for key, val := range currencyMap {
		var sub constant.SubAccount
		sub.Amount = val
		e.acc.SubAccounts[key] = sub
	}
	return nil
}

// position2ValDiff ...
func (ex *ExchangeFutureBack) position2ValDiff(last float64, position constant.Position) float64 {
	amount := position.Amount + position.FrozenAmount
	price := position.Price
	priceDiff := last - price
	priceRate := util.SafefloatDivide(priceDiff, price)
	valDiff := priceRate * amount
	return valDiff
}

func (e *ExchangeFutureBack) settlePositionProfit(last float64, position *constant.Position, asset *constant.SubAccount, dir int) {
	valdiff := e.position2ValDiff(last, *position)
	amountdiff := util.SafefloatDivide(valdiff*e.contractRate, last)
	if dir == 1 {
		amountdiff = 0 - amountdiff
	}
	position.Profit = amountdiff
	position.ProfitRate = util.SafefloatDivide(valdiff, position.Amount+position.FrozenAmount)
}

// settlePosition ...
func (ex *ExchangeFutureBack) settlePosition(stockType string) {
	ticker := ex.currData[stockType]
	last := ticker.Close
	stocks := stockPair2Vec(stockType)
	CurrencyA := stocks[0]
	asset := ex.acc.SubAccounts[CurrencyA]

	if _, ok := ex.longPosition[CurrencyA]; ok {
		position := ex.longPosition[CurrencyA]
		ex.settlePositionProfit(last, &position, &asset, 0)
		ex.acc.SubAccounts[CurrencyA] = asset
		ex.longPosition[CurrencyA] = position
	}

	if _, ok := ex.shortPosition[CurrencyA]; ok {
		position := ex.shortPosition[CurrencyA]
		ex.settlePositionProfit(last, &position, &asset, 1)
		ex.acc.SubAccounts[CurrencyA] = asset
		ex.shortPosition[CurrencyA] = position
	}
}

// fillOrder ...
func (ex *ExchangeFutureBack) fillOrder(isTaker bool, amount, price float64, ord *constant.Order) {
	ord.FinishedTime = ex.currData[ord.StockType].Time / int64(time.Millisecond) //set filled time
	ord.DealAmount = ord.Amount
	dealAmount := ord.DealAmount
	ord.Status = constant.ORDER_FINISH

	fee := ex.makerFee
	if isTaker {
		fee = ex.takerFee
	}

	tradeFee := 0.0
	tradeFee = dealAmount * ex.contractRate * fee
	tradeFee = math.Floor(tradeFee*100000000) / 100000000

	ord.Fee += tradeFee

	ex.unFrozenAsset(tradeFee, dealAmount, price, *ord)
}

func (ex *ExchangeFutureBack) matchOrder(ord *constant.Order, isTaker bool) {
	ticker := ex.currData[ord.StockType]
	switch ord.TradeType {
	case constant.TradeTypeLong, constant.TradeTypeShortClose:
		if ticker.Close <= ord.Price && ticker.Volume > 0 {
			ex.fillOrder(isTaker, ticker.Volume, ticker.Close, ord)
			if ord.Status == constant.ORDER_FINISH {
				delete(ex.pendingOrders, ord.Id)
				ex.finishedOrders[ord.Id] = ord
			}
		}
	case constant.TradeTypeShort, constant.TradeTypeLongClose:
		if ticker.Close >= ord.Price && ticker.Volume > 0 {
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

func (ex *ExchangeFutureBack) coverPosition(stockType string) {
	stocks := stockPair2Vec(stockType)
	CurrencyA := stocks[0]
	marginRatio, _, rht := ex.marginRatio()
	if marginRatio < 0 || rht > 0.0 && marginRatio < ex.coverRate {

		log.Infof("force cover %f -> %f\n", marginRatio, ex.coverRate)
		ex.Debug()
		//Force cover position
		ex.longPosition[CurrencyA] = constant.Position{}
		ex.shortPosition[CurrencyA] = constant.Position{}
		ex.acc.SubAccounts[CurrencyA] = constant.SubAccount{}
		for _, order := range ex.pendingOrders {
			order.Status = constant.ORDER_CANCEL_ING
			ex.finishedOrders[order.Id] = order
		}
		ex.pendingOrders = make(map[string]*constant.Order)
	}
	if position, ok := ex.longPosition[CurrencyA]; ok && position.Amount == 0 {
		delete(ex.longPosition, CurrencyA)
	}

	if position, ok := ex.shortPosition[CurrencyA]; ok && position.Amount == 0 {
		delete(ex.shortPosition, CurrencyA)
	}
}

func (ex *ExchangeFutureBack) marginRatio() (float64, float64, float64) {
	stockType := ex.BaseExchange.GetStockType()
	stocks := stockPair2Vec(stockType)
	CurrencyA := stocks[0]
	asset := ex.acc.SubAccounts[CurrencyA]
	longposition := ex.longPosition[CurrencyA]
	shortposition := ex.shortPosition[CurrencyA]
	lft := asset.Amount + longposition.Profit + shortposition.Profit
	rht := 0.0
	rht = rht + util.SafefloatDivide(longposition.Amount*ex.BaseExchange.contractRate, longposition.Price)
	rht = rht + util.SafefloatDivide(shortposition.Amount*ex.BaseExchange.contractRate, shortposition.Price)
	rht = rht + asset.FrozenAmount*ex.lever
	return util.SafefloatDivide(lft, rht), lft, rht
}

// LimitBuy ...
func (ex *ExchangeFutureBack) LimitBuy(amount, price, currency string) (*constant.Order, error) {
	ex.Lock()
	defer ex.Unlock()

	ord := constant.Order{
		Price:     goex.ToFloat64(price),
		Amount:    goex.ToFloat64(amount),
		Id:        ex.idGen.Get(),
		Time:      ex.currData[currency].Time / int64(time.Millisecond),
		Status:    constant.ORDER_UNFINISH,
		StockType: currency,
		TradeType: ex.BaseExchange.GetDirection(),
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
		Price: goex.ToFloat64(price),
		//OpenPrice: ex.currData.Close,
		Amount:    goex.ToFloat64(amount),
		Id:        ex.idGen.Get(),
		Time:      ex.currData[currency].Time / int64(time.Millisecond),
		Status:    constant.ORDER_UNFINISH,
		StockType: currency,
		TradeType: ex.BaseExchange.GetDirection(),
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
		longposition := ex.longPosition[key]
		shortposition := ex.shortPosition[key]
		sub.ProfitUnreal += longposition.Profit
		sub.ProfitUnreal += shortposition.Profit
		sub.Amount = sub.ProfitUnreal
		account.SubAccounts[key] = sub
	}
	return &account, nil
}

// GetTicker ...
func (ex *ExchangeFutureBack) GetTicker(currency string) (*constant.Ticker, error) {
	var ohlc *constant.OHLC
	for symbol, loader := range ex.dataLoader {
		if loader == nil {
			return nil, errors.New("loader not found")
		}
		curr := loader.Next()
		if curr == nil {
			return nil, nil
		}
		if symbol == currency {
			ohlc = curr
		}
		ex.currData[currency] = *curr
		ex.match()
		ex.settlePosition(currency)
		ex.coverPosition(currency)
	}
	ex.Debug()
	if ohlc == nil {
		//backtest end
		return nil, nil
	}
	return &constant.Ticker{
		Vol:  ohlc.Volume,
		Time: ohlc.Time,
		Last: ohlc.Close,
		Buy:  ohlc.Close,
		Sell: ohlc.Close,
		High: ohlc.High,
		Low:  ohlc.Low,
	}, nil
}

// GetDepth ...
func (ex *ExchangeFutureBack) GetDepth(size int, currency string) (*constant.Depth, error) {
	return nil, fmt.Errorf("not support")
}

// GetExchangeName ...
func (ex *ExchangeFutureBack) GetExchangeName() string {
	return ex.name
}

// frozenAsset 冻结
func (ex *ExchangeFutureBack) frozenAsset(order constant.Order) error {
	stocks := stockPair2Vec(order.StockType)
	CurrencyA := stocks[0]
	ticker := ex.currData[order.StockType]
	var price float64 = 1
	price = ticker.Close
	avaAmount := ex.acc.SubAccounts[CurrencyA].Amount
	//avaAmount
	longposition := ex.longPosition[CurrencyA]
	shortposition := ex.shortPosition[CurrencyA]
	if longposition.Profit < 0 {
		avaAmount += longposition.Profit
	}

	if shortposition.Profit < 0 {
		avaAmount += shortposition.Profit
	}
	lever := ex.BaseExchange.lever
	switch order.TradeType {
	case constant.TradeTypeLong, constant.TradeTypeShort:
		if avaAmount*lever*price < order.Amount*ex.BaseExchange.contractRate {
			return fmt.Errorf("open %s not insufficient : %f->%f", CurrencyA, avaAmount, order.Amount)
		}
		costAmount := util.SafefloatDivide(order.Amount*ex.BaseExchange.contractRate, lever*price)
		ex.acc.SubAccounts[CurrencyA] = constant.SubAccount{
			StockType:    CurrencyA,
			Amount:       avaAmount - costAmount,
			FrozenAmount: ex.acc.SubAccounts[CurrencyA].FrozenAmount + costAmount,
			LoanAmount:   0,
		}
	case constant.TradeTypeLongClose, constant.TradeTypeShortClose:
		if order.TradeType == constant.TradeTypeLongClose {
			position := ex.longPosition[CurrencyA]
			if position.Amount < order.Amount {
				return fmt.Errorf("close %s not insufficient : %f->%f", CurrencyA, avaAmount, order.Amount)
			}
			position.Amount = position.Amount - order.Amount
			position.FrozenAmount = position.FrozenAmount + order.Amount
			ex.longPosition[CurrencyA] = position
		}

		if order.TradeType == constant.TradeTypeShortClose {
			position := ex.shortPosition[CurrencyA]
			if position.Amount < order.Amount {
				return fmt.Errorf("close %s not insufficient : %f->%f", CurrencyA, avaAmount, order.Amount)
			}
			position.Amount = position.Amount - order.Amount
			position.FrozenAmount = position.FrozenAmount + order.Amount
			ex.shortPosition[CurrencyA] = position
		}
	}
	return nil
}

// unFrozenAsset 解冻
func (ex *ExchangeFutureBack) unFrozenAsset(fee, matchAmount, matchPrice float64, order constant.Order) {
	stockType, _ := ex.getSymbol(order.StockType)
	stocks := stockPair2Vec(stockType)
	CurrencyA := stocks[0]
	assetA := ex.acc.SubAccounts[CurrencyA]
	lever := ex.BaseExchange.lever
	switch order.TradeType {
	case constant.TradeTypeLong, constant.TradeTypeShort:
		order.OpenPrice = ex.currData[order.StockType].Close
		costAmount := util.SafefloatDivide(order.Amount*ex.BaseExchange.contractRate, lever*order.OpenPrice)
		if order.Status == constant.ORDER_CANCEL {
			ex.acc.SubAccounts[assetA.StockType] = constant.SubAccount{
				StockType:    assetA.StockType,
				Amount:       assetA.Amount + costAmount - order.DealAmount,
				FrozenAmount: assetA.FrozenAmount - (costAmount - order.DealAmount),
				LoanAmount:   0,
			}
		} else {
			if order.TradeType == constant.TradeTypeLong {
				position := ex.longPosition[CurrencyA]
				position.Price = util.SafefloatDivide(position.Price*position.Amount+order.OpenPrice*order.Amount,
					position.Amount+order.Amount)
				position.Amount = position.Amount + order.Amount
				ex.longPosition[CurrencyA] = position
				//log.Infof("set long position as:%v\n", position)
			} else {
				position := ex.shortPosition[CurrencyA]
				position.Price = util.SafefloatDivide(position.Price*position.Amount+order.OpenPrice*order.Amount,
					position.Amount+order.Amount)
				position.Amount = position.Amount + order.Amount
				ex.shortPosition[CurrencyA] = position
				//log.Infof("set short position as:%v\n", position)
			}
			ex.acc.SubAccounts[assetA.StockType] = constant.SubAccount{
				StockType:    assetA.StockType,
				FrozenAmount: assetA.FrozenAmount - costAmount,
				Amount:       assetA.Amount,
				LoanAmount:   0,
			}
		}
	case constant.TradeTypeLongClose, constant.TradeTypeShortClose:
		var position constant.Position
		order.OpenPrice = ex.currData[order.StockType].Close
		if order.TradeType == constant.TradeTypeLongClose {
			position = ex.longPosition[CurrencyA]
		} else {
			position = ex.shortPosition[CurrencyA]
		}
		if order.Status == constant.ORDER_CANCEL {
			position.Amount = position.Amount + order.Amount
			position.FrozenAmount = position.FrozenAmount - order.Amount
		} else {
			//ex.longPosition[CurrencyA] = position
			position.FrozenAmount = position.FrozenAmount - order.Amount
			ex.longPosition[CurrencyA] = position
			costAmount := (order.Amount * ex.BaseExchange.contractRate) / (lever * order.OpenPrice)
			ex.acc.SubAccounts[assetA.StockType] = constant.SubAccount{
				StockType:    assetA.StockType,
				Amount:       assetA.Amount + costAmount + fee,
				FrozenAmount: assetA.FrozenAmount,
				LoanAmount:   0,
			}
		}
		if order.TradeType == constant.TradeTypeLongClose {
			ex.longPosition[CurrencyA] = position
		} else {
			ex.shortPosition[CurrencyA] = position
		}

	}
}

// GetRecords get candlestick data
func (e *ExchangeFutureBack) GetRecords() ([]constant.Record, error) {
	size := e.GetPeriodSize()

	ticker, err := e.GetTicker(e.GetStockType())
	if err != nil {
		return nil, err
	}
	if ticker == nil {
		return nil, nil
	}
	curr := e.currData[e.GetStockType()].Time

	if e.recordsCache == nil {
		e.recordsCache = make(map[string][]constant.Record)
	}
	if e.recordsMap == nil {
		e.recordsMap = make(map[string]map[int64]int)
	}
	key := e.GetStockType()
	// try to store records in cache at first
	if len(e.recordsCache[key]) == 0 {
		vec := e.dataLoader[key].datas
		var records []constant.Record
		for i := 0; i < len(vec); i++ {
			kline := vec[i]
			records = append([]constant.Record{{
				Open:   kline.Open,
				High:   kline.High,
				Low:    kline.Low,
				Close:  kline.Close,
				Volume: kline.Volume,
				Time:   kline.Time,
			}}, records...)
			_, ok := e.recordsMap[key]
			if !ok {
				e.recordsMap[key] = make(map[int64]int)
			}
			e.recordsMap[key][kline.Time] = i
		}
		e.recordsCache[key] = records
	}
	end := e.recordsMap[key][curr]
	start := 0
	if end > size && size != 0 {
		start = end - size
	}
	records := e.recordsCache[key][start:end]

	return records, nil
}
