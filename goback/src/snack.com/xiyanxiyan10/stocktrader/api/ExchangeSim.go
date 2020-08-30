package api

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	goex "github.com/nntaoli-project/goex"
	dbtypes "snack.com/xiyanxiyan10/stockdb/types"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

var (
	DataFinishedError        = errors.New("depth data finished")
	InsufficientError        = errors.New("insufficient")
	CancelOrderFinishedError = errors.New("order finished")
	NotFoundOrderError       = errors.New("not found order")
	AssetSnapshotCsvFileName = "%s_asset_snapshot.csv"
)

type ExchangeSim struct {
	*sync.RWMutex
	acc                  *constant.Account
	name                 string
	makerFee             float64
	takerFee             float64
	supportCurrencyPairs []string
	quoteCurrency        string
	pendingOrders        map[string]*constant.Order
	finishedOrders       map[string]*constant.Order
	depthLoader          map[string]*DepthDataLoader
	currDepth            dbtypes.OHLC
	idGen                *util.IDGen

	sortedCurrencies constant.Account
}

func NewExchangeSim(config ExchangeSimConfig) *ExchangeSim {
	sim := &ExchangeSim{
		RWMutex:              new(sync.RWMutex),
		idGen:                NewIdGen(config.ExName),
		name:                 config.ExName,
		makerFee:             config.MakerFee,
		takerFee:             config.TakerFee,
		acc:                  &config.Account,
		supportCurrencyPairs: config.SupportCurrencyPairs,
		quoteCurrency:        config.QuoteCurrency,
		pendingOrders:        make(map[string]*constant.Order, 100),
		finishedOrders:       make(map[string]*constant.Order, 100),
		depthLoader:          make(map[string]*DepthDataLoader, 1),
	}

	for _, pair := range config.SupportCurrencyPairs {
		/*
			if !pair.CurrencyB.Eq(config.QuoteCurrency) {
				panic("the CurrencyPair only one quote currency per backtest")
			}
			sim.depthLoader[pair] = NewDepthDataLoader(DataConfig{
				Ex:       sim.name,
				Pair:     pair,
				StarTime: config.BackTestStartTime,
				EndTime:  config.BackTestEndTime,
				UnGzip:   config.UnGzip,
				Size:     config.DepthSize,
			})
		*/
	}

	for key, sub := range sim.acc.SubAccounts {
		sim.sortedCurrencies.SubAccounts[key] = sub
	}

	return sim
}

func (ex *ExchangeSim) fillOrder(isTaker bool, amount, price float64, ord *constant.Order) {
	ord.FinishedTime = ex.currDepth.UTime.UnixNano() / int64(time.Millisecond) //set filled time

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

func (ex *ExchangeSim) matchOrder(ord *constnat.Order, isTaker bool) {
	ticker := ex.currDepth
	switch ord.TradeType {
	case constant.TradeTypeSell:
		if ticker.Close >= ord.Price && ticker.Volume > 0 {
			ex.fillOrder(isTaker, ticker.Volume, ticker.Close, ord)
			if ord.Status == constant.ORDER_FINISH {
				delete(ex.pendingOrders, ord.OrderID2)
				ex.finishedOrders[ord.OrderID2] = ord
			}
		}
	case constant.TradeTypeBuy:
		if ticker.Close <= ord.Price && ticker.Volume > 0 {
			ex.fillOrder(isTaker, ask.Amount, ask.Price, ord)
			if ord.Status == goex.ORDER_FINISH {
				delete(ex.pendingOrders, ord.OrderID2)
				ex.finishedOrders[ord.OrderID2] = ord
			}
		}
	}
}

func (ex *ExchangeSim) match() {
	ex.Lock()
	defer ex.Unlock()
	for id, _ := range ex.pendingOrders {
		ex.matchOrder(ex.pendingOrders[id], false)
	}
}

func stockPair2Vec(pair string) []string {
	res := strings.Split(pair, "/")
	if len(res) < 2 {
		return []string{"", ""}
	}
}

func (ex *ExchangeSim) LimitBuy(amount, price, currency string) (*constant.Order, error) {
	ex.Lock()
	defer ex.Unlock()

	ord := constant.Order{
		Price:     goex.ToFloat64(price),
		Amount:    goex.ToFloat64(amount),
		Id:        ex.idGen.Get(),
		OrderTime: int(ex.currDepth.Time / int64(time.Millisecond)),
		Status:    constant.ORDER_UNFINISH,
		StockType: currency,
		TradeType: constant.TradeTypeBuy,
	}
	//ord.Cid = ord.OrderID2

	err := ex.frozenAsset(ord)
	if err != nil {
		return nil, err
	}

	ex.pendingOrders[ord.Id] = &ord

	ex.matchOrder(&ord, true)

	var result goex.Order
	DeepCopyStruct(ord, &result)
	return &result, nil
}

func (ex *ExchangeSim) LimitSell(amount, price, currency string) (*goex.Order, error) {
	ex.Lock()
	defer ex.Unlock()

	ord := constant.Order{
		Price:     goex.ToFloat64(price),
		Amount:    goex.ToFloat64(amount),
		Id:        ex.idGen.Get(),
		OrderTime: int(ex.currDepth.Time / int64(time.Millisecond)),
		Status:    constant.ORDER_UNFINISH,
		StockType: currency,
		TradeType: constant.TradeTypeSell,
	}
	//ord.Cid = ord.OrderID2

	err := ex.frozenAsset(ord)
	if err != nil {
		return nil, err
	}

	ex.pendingOrders[ord.Id] = &ord

	ex.matchOrder(&ord, true)

	var result goex.Order
	DeepCopyStruct(ord, &result)

	return &result, nil
}

func (ex *ExchangeSim) MarketBuy(amount, price, currency string) (*constant.Order, error) {
	panic("not support")
}

func (ex *ExchangeSim) MarketSell(amount, price, currency string) (*constant.Order, error) {
	panic("not support")
}

func (ex *ExchangeSim) CancelOrder(orderId string, currency string) (bool, error) {
	ex.Lock()
	defer ex.Unlock()

	ord := ex.finishedOrders[orderId]
	if ord != nil {
		return false, CancelOrderFinishedError
	}

	ord = ex.pendingOrders[orderId]
	if ord == nil {
		return false, NotFoundOrderError
	}

	delete(ex.pendingOrders, ord.OrderID2)

	ord.Status = constant.ORDER_CANCEL
	ex.finishedOrders[ord.Id] = ord

	ex.unFrozenAsset(0, 0, 0, *ord)

	return true, nil
}

func (ex *ExchangeSim) GetOneOrder(orderId, currency string) (*constant.Order, error) {
	ex.RLock()
	defer ex.RUnlock()

	ord := ex.finishedOrders[orderId]
	if ord == nil {
		ord = ex.pendingOrders[orderId]
	}

	if ord != nil {
		// deep copy
		var result constant.Order
		DeepCopyStruct(ord, &result)

		return &result, nil
	}

	return nil, NotFoundOrderError
}

func (ex *ExchangeSim) GetUnfinishOrders(currency string) ([]constant.Order, error) {
	ex.RLock()
	defer ex.RUnlock()

	var unfinishedOrders []constant.Order
	for _, ord := range ex.pendingOrders {
		unfinishedOrders = append(unfinishedOrders, *ord)
	}

	return unfinishedOrders, nil
}

func (ex *ExchangeSim) GetOrderHistorys(currency string, currentPage, pageSize int) ([]constant.Order, error) {
	ex.RLock()
	defer ex.RUnlock()

	var orders []goex.Order
	for _, ord := range ex.finishedOrders {
		if ord.StockType == currency {
			orders = append(orders, *ord)
		}
	}
	return orders, nil
}

func (ex *ExchangeSim) GetAccount() (*constant.Account, error) {
	ex.RLock()
	defer ex.RUnlock()

	var account constant.Account
	account.SubAccounts = make(map[string]constant.SubAccount)
	for key, sub := range ex.acc.SubAccounts {
		account.SubAccounts[key] = sub
	}
	return &account, nil
}

func (ex *ExchangeSim) GetTicker(currency string) (*constant.Ticker, error) {
	curr := ex.currDepth
	return &constant.Ticker{
		Last: curr.Close,
		Buy:  curr.Close,
		Sell: curr.Close,
		High: curr.High,
		Low:  curr.Low,
	}, nil
}

func (ex *ExchangeSim) GetDepth(size int, currency string) (*goex.Depth, error) {
	depth := ex.depthLoader[currency].Next()
	if depth == nil {
		return nil, DataFinishedError
	}
	ex.currDepth = *depth
	ex.match()
	return depth, nil
}

func (ex *ExchangeSim) GetKlineRecords(currency goex.CurrencyPair, period, size, since int) ([]goex.Kline, error) {
	panic("not support")
}

func (ex *ExchangeSim) GetTrades(currencyPair goex.CurrencyPair, since int64) ([]goex.Trade, error) {
	panic("not support")
}

func (ex *ExchangeSim) GetExchangeName() string {
	return ex.name
}

//冻结
func (ex *ExchangeSim) frozenAsset(order goex.Order) error {

	switch order.Side {
	case goex.SELL:
		avaAmount := ex.acc.SubAccounts[order.Currency.CurrencyA].Amount
		if avaAmount < order.Amount {
			return InsufficientError
		}
		ex.acc.SubAccounts[order.Currency.CurrencyA] = goex.SubAccount{
			Currency:     order.Currency.CurrencyA,
			Amount:       avaAmount - order.Amount,
			ForzenAmount: ex.acc.SubAccounts[order.Currency.CurrencyA].ForzenAmount + order.Amount,
			LoanAmount:   0,
		}
	case goex.BUY:
		avaAmount := ex.acc.SubAccounts[order.Currency.CurrencyB].Amount
		need := order.Amount * order.Price
		if avaAmount < need {
			return InsufficientError
		}
		ex.acc.SubAccounts[order.Currency.CurrencyB] = goex.SubAccount{
			Currency:     order.Currency.CurrencyB,
			Amount:       avaAmount - need,
			ForzenAmount: ex.acc.SubAccounts[order.Currency.CurrencyB].ForzenAmount + need,
			LoanAmount:   0,
		}
	}

	ex.assetSnapshot()

	return nil
}

//解冻
func (ex *ExchangeSim) unFrozenAsset(fee, matchAmount, matchPrice float64, order goex.Order) {
	assetA := ex.acc.SubAccounts[order.Currency.CurrencyA]
	assetB := ex.acc.SubAccounts[order.Currency.CurrencyB]

	switch order.Side {
	case goex.SELL:
		if order.Status == goex.ORDER_CANCEL {
			ex.acc.SubAccounts[assetA.Currency] = goex.SubAccount{
				Currency:     assetA.Currency,
				Amount:       assetA.Amount + order.Amount - order.DealAmount,
				ForzenAmount: assetA.ForzenAmount - (order.Amount - order.DealAmount),
				LoanAmount:   0,
			}
		} else {
			ex.acc.SubAccounts[assetA.Currency] = goex.SubAccount{
				Currency:     assetA.Currency,
				Amount:       assetA.Amount,
				ForzenAmount: assetA.ForzenAmount - matchAmount,
				LoanAmount:   0,
			}
			ex.acc.SubAccounts[assetB.Currency] = goex.SubAccount{
				Currency:     assetB.Currency,
				Amount:       assetB.Amount + matchAmount*matchPrice - fee,
				ForzenAmount: assetB.ForzenAmount,
			}
		}

	case goex.BUY:
		if order.Status == goex.ORDER_CANCEL {
			unFrozen := (order.Amount - order.DealAmount) * order.Price
			ex.acc.SubAccounts[assetB.Currency] = goex.SubAccount{
				Currency:     assetB.Currency,
				Amount:       assetB.Amount + unFrozen,
				ForzenAmount: assetB.ForzenAmount - unFrozen,
			}
		} else {
			ex.acc.SubAccounts[assetA.Currency] = goex.SubAccount{
				Currency:     assetA.Currency,
				Amount:       assetA.Amount + matchAmount - fee,
				ForzenAmount: assetA.ForzenAmount,
				LoanAmount:   0,
			}
			ex.acc.SubAccounts[assetB.Currency] = goex.SubAccount{
				Currency:     assetB.Currency,
				Amount:       assetB.Amount + matchAmount*(order.Price-matchPrice),
				ForzenAmount: assetB.ForzenAmount - matchAmount*order.Price,
			}
		}
	}

	ex.assetSnapshot()
}

func (ex *ExchangeSim) assetSnapshot() {
	csvFile := fmt.Sprintf(AssetSnapshotCsvFileName, ex.name)
	f, err := os.OpenFile(csvFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0744)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Println("close file error=", err)
		}
	}()

	csvW := csv.NewWriter(f)

	var (
		data []string
	)

	netAsset := 0.0
	for _, currency := range ex.sortedCurrencies {
		sub := ex.acc.SubAccounts[currency]
		data = append(data, goex.FloatToString(sub.Amount, 10))
		data = append(data, goex.FloatToString(sub.ForzenAmount, 10))
		if currency.Eq(ex.quoteCurrency) {
			netAsset += sub.Amount + sub.ForzenAmount
		} else {
			pair := goex.NewCurrencyPair(currency, ex.quoteCurrency)
			ticker, err := ex.GetTicker(pair)
			if err != nil {
				log.Println("[ERROR] GetTicker CurrencyPair=", pair.ToSymbol(""), ",error=", err)
				continue
			}
			netAsset += (sub.Amount + sub.ForzenAmount) * ticker.Buy
		}
	}
	data = append(data, goex.FloatToString(netAsset, 10))

	csvW.Write(data)
	csvW.Flush()
}
