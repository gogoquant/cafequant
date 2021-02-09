package api

import (
	"errors"
	"fmt"
	dbconstant "snack.com/xiyanxiyan10/stockdb/constant"
	dbsdk "snack.com/xiyanxiyan10/stockdb/sdk"
	dbtypes "snack.com/xiyanxiyan10/stockdb/types"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"strings"
	"time"
)

var (
	// ErrDataFinished ...
	ErrDataFinished = errors.New("depth data finished")
	// ErrDataInsufficient ...
	ErrDataInsufficient = errors.New("insufficient")
	// ErrCancelOrderFinished ...
	ErrCancelOrderFinished = errors.New("order finished")
	// ErrNotFoundOrder ...
	ErrNotFoundOrder = errors.New("not found order")
)

// DataConfig ...
type DataConfig struct {
	Ex     string
	Pair   string
	Size   int //多少档深度数据
	UnGzip bool
}

// DataLoader ...
type DataLoader struct {
	curr  int
	size  int
	datas []dbtypes.OHLC
}

// Next ...
func (l *DataLoader) Next() *dbtypes.OHLC {
	nextPos := l.curr + 1
	if nextPos >= l.size {
		return nil
	}
	data := l.datas[l.curr]
	l.curr = nextPos
	return &data
}

// Load ...
func (l *DataLoader) Load(ohlcs []dbtypes.OHLC) {
	l.datas = append(l.datas, ohlcs...)
	l.size = len(l.datas)
	// move one for first records, at least one ohlc
	l.Next()
}

// BaseExchange ...
type BaseExchange struct {
	BaseExchangeCaches // cache for exchange
	// period for.father.get records
	periodVal string
	// period for backtest
	period string
	size   int
	id     int    // id of the exchange
	ioMode string // io mode for exchange
	//contractType       string  // contractType
	direction          string  // trade type
	stockType          string  // stockType
	lever              float64 // lever
	recordsPeriodMap   map[string]int64
	recordsPeriodDbMap map[string]int64
	// recordsPeriod support
	minAmountMap map[string]float64 // minAmount of trade
	limit        int64
	lastSleep    int64
	lastTimes    int64
	subscribeMap map[string][]string
	currencyMap  map[string]float64

	coverRate    float64
	taker        float64
	maker        float64
	coin         bool
	contractRate float64 // 合约每张价值
	//currencyStandard bool    // 是否为币本位

	start  int64
	end    int64
	host   string
	logger model.Logger
	option constant.Option

	father ExchangeBroker
}

func stockPair2Vec(pair string) []string {
	res := strings.Split(pair, "/")
	if len(res) < 2 {
		return []string{"", ""}
	}
	return res
}

func (e *BaseExchange) Buy(price, amount, msg string) (string, error) {
	return e.father.buy(price, amount, msg)
}

func (e *BaseExchange) Sell(price, amount, msg string) (string, error) {
	return e.father.sell(price, amount, msg)
}

func (e *BaseExchange) CancelOrder(orderID string) (bool, error) {
	return e.father.cancelOrder(orderID)
}

// SetPeriodSize Set size
func (e *BaseExchange) SetPeriodSize(size int) {
	e.size = size

}

// GetPeriodSize Get Size
func (e *BaseExchange) GetPeriodSize() int {
	return e.size
}

// SetBackCommission 设置回测手续费
func (e *BaseExchange) SetBackCommission(taker, maker, contractRate, coverRate float64, coin bool) {
	e.contractRate = contractRate
	e.taker = taker
	e.maker = maker
	e.coin = coin
	e.coverRate = coverRate
}

// GetBackCommission 获取回测手续费
/*
func (e *BaseExchange) GetBackCommission() (float64, float64, float64, float64, bool) {
	return e.taker, e.maker, e.contractRate, e.coverRate, e.coin
}
*/

// SetBackTime ...
func (e *BaseExchange) SetBackTime(start, end int64, period string) {
	e.start = start
	e.end = end
	e.period = period
}

// GetBackAccount ...
func (e *BaseExchange) GetBackAccount() map[string]float64 {
	return e.currencyMap
}

// SetBackAccount ...
func (e *BaseExchange) SetBackAccount(key string, val float64) {
	e.currencyMap[key] = val
}

// GetBackTime ...
/*
func (e *BaseExchange) GetBackTime() (int64, int64, string) {
	return e.start, e.end, e.period
}
*/

// GetSubscribe ...
func (e *BaseExchange) GetSubscribe() map[string][]string {
	return e.subscribeMap
}

// IsSubscribe ...
func (e *BaseExchange) IsSubscribe(source, action string) bool {
	actions := e.subscribeMap[source]
	for _, tmp := range actions {
		if tmp == action {
			return true
		}
	}
	return false
}

// SetSubscribe ...
func (e *BaseExchange) SetSubscribe(source, action string) {
	if e.subscribeMap == nil {
		e.subscribeMap = make(map[string][]string)
	}
	e.subscribeMap[source] = append(e.subscribeMap[source], action)
}

// SetLimit set the limit calls amount per second of this exchange
func (e *BaseExchange) SetLimit(times int64) int64 {
	e.limit = times
	return e.limit
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *BaseExchange) AutoSleep() {
	if e.option.BackTest {
		return
	}
	if e.limit == 0 {
		e.limit = 1000
	}
	time.Sleep(time.Duration(e.limit) * time.Millisecond)
	e.lastTimes = 0
	e.lastSleep = time.Now().UnixNano()
}

// Sleep ...
func (e *BaseExchange) Sleep(intervals int64) {
	if e.option.BackTest {
		return
	}
	time.Sleep(time.Duration(intervals) * time.Millisecond)
}

// BackGetStats ...
func (e *BaseExchange) BackGetStats() ([]dbtypes.Stats, error) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
				fmt.Sprintf("GetStats error, the error number is %s", r))
		}
	}()
	client := dbsdk.NewClient(config.String(constant.STOCKDBURL), config.String(constant.STOCKDBAUTH))
	ohlc := client.GetStats()
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			fmt.Sprintf("GetStats error, the error number is %s", ohlc.Message))
		return nil, fmt.Errorf("GetStats error, the error number is %s", ohlc.Message)
	}
	return ohlc.Data, nil
}

// BackGetMarkets ...
func (e *BaseExchange) BackGetMarkets() ([]string, error) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
				fmt.Sprintf("GetMarkets error, the error number is %s", r))
		}
	}()
	client := dbsdk.NewClient(config.String(constant.STOCKDBURL), config.String(constant.STOCKDBAUTH))
	ohlc := client.GetMarkets()
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			fmt.Sprintf("GetMarkets error, the error number is %s", ohlc.Message))
		return nil, fmt.Errorf("GetMarkets error, the error number is %s", ohlc.Message)
	}
	return ohlc.Data, nil
}

// BackGetSymbols  ...
func (e *BaseExchange) BackGetSymbols() ([]string, error) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
				fmt.Sprintf("GetSymbol error, the error number is %s", r))
		}
	}()
	client := dbsdk.NewClient(config.String(constant.STOCKDBURL), config.String(constant.STOCKDBAUTH))
	ohlc := client.GetSymbols(e.option.Type)
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			fmt.Sprintf("GetSymbols error, the error number is %s", ohlc.Message))
		return nil, fmt.Errorf("GetSymbols error, the error number is %s", ohlc.Message)
	}
	return ohlc.Data, nil
}

// BackGetOHLCs ...
func (e *BaseExchange) BackGetOHLCs(begin, end int64, period string) ([]dbtypes.OHLC, error) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
				fmt.Sprintf("GetOHLCs error, the error number is %s", r))
		}
	}()
	var opt dbtypes.Option
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	opt.Period = e.recordsPeriodDbMap[period]
	opt.BeginTime = begin
	opt.EndTime = end
	client := dbsdk.NewClient(config.String(constant.STOCKDBURL), config.String(constant.STOCKDBAUTH))
	ohlc := client.GetOHLCs(opt)
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			fmt.Sprintf("GetOHLCs error, the error number is %s", ohlc.Message))
		return nil, fmt.Errorf("GetOHLCs error, the error number is %s", ohlc.Message)
	}
	return ohlc.Data, nil
}

// BackPutOHLC ...
func (e *BaseExchange) BackPutOHLC(time int64, open, high, low, closed, volume float64, ext string, period string) error {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
				fmt.Sprintf("PutOHLC error, the error number is %s", r))
		}
	}()
	var opt dbtypes.Option
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	opt.Period = e.recordsPeriodDbMap[period]
	client := dbsdk.NewClient(config.String(constant.STOCKDBURL), config.String(constant.STOCKDBAUTH))
	var datum dbtypes.OHLC
	datum.Time = time
	datum.Open = open
	datum.High = high
	datum.Low = low
	datum.Close = closed
	datum.Volume = volume
	ohlc := client.PutOHLC(datum, opt)
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			fmt.Sprintf("PutOHLC error, the error number is %s\n", ohlc.Message))
		return fmt.Errorf("PutOHLC error, the error number is %s", ohlc.Message)
	}
	return nil
}

// BackPutOHLCs ...
func (e *BaseExchange) BackPutOHLCs(datums []dbtypes.OHLC, period string) error {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
				fmt.Sprintf("PutOHLCs error, the error number is %s", r))
		}
	}()
	var opt dbtypes.Option
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	opt.Period = e.recordsPeriodDbMap[period]
	client := dbsdk.NewClient(config.String(constant.STOCKDBURL), config.String(constant.STOCKDBAUTH))
	ohlc := client.PutOHLCs(datums, opt)
	if !ohlc.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			fmt.Sprintf("PutOHLCs error, the error number is %s\n", ohlc.Message))
		return fmt.Errorf("PutOHLCs error, the error number is %s", ohlc.Message)
	}
	return nil
}

// BackGetTimeRange ...
func (e *BaseExchange) BackGetTimeRange() ([2]int64, error) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
				fmt.Sprintf("GeTimeRanege error, the error number is %s", r))
		}
	}()
	var opt dbtypes.Option
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	client := dbsdk.NewClient(config.String(constant.STOCKDBURL), config.String(constant.STOCKDBAUTH))
	timeRange := client.GetTimeRange(opt)
	if !timeRange.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			fmt.Sprintf("GetTimeRange, the error number is %s", timeRange.Message))
	}
	return timeRange.Data, nil
}

// Set Period
func (e *BaseExchange) SetPeriod(period string) {
	e.periodVal = period
}

// Get Period
func (e *BaseExchange) GetPeriod() string {
	return e.periodVal
}

// BackGetPeriodRange ...
func (e *BaseExchange) BackGetPeriodRange() ([2]int64, error) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
				fmt.Sprintf("GetPeriodRange error, the error number is %s", r))
		}
	}()
	var opt dbtypes.Option
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	client := dbsdk.NewClient(config.String(constant.STOCKDBURL), config.String(constant.STOCKDBAUTH))
	timeRange := client.GetPeriodRange(opt)
	if !timeRange.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			fmt.Sprint("GetPeriodRange, the error number is %s"+timeRange.Message))
		return [2]int64{}, fmt.Errorf("GetPeriodRange, the error number is " + timeRange.Message)
	}
	return timeRange.Data, nil
}

// BackGetDepth ...
func (e *BaseExchange) BackGetDepth(begin, end int64, period string) (dbtypes.Depth, error) {

	defer func() {
		if r := recover(); r != nil {
			e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
				fmt.Sprintf("GetDepth error, the error number is %s", r))
		}
	}()
	var opt dbtypes.Option
	opt.Market = e.option.Type
	opt.Symbol = e.GetStockType()
	opt.Period = e.recordsPeriodDbMap[period]
	opt.BeginTime = begin
	opt.EndTime = end
	client := dbsdk.NewClient(config.String(constant.STOCKDBURL), config.String(constant.STOCKDBAUTH))
	depth := client.GetDepth(opt)
	if !depth.Success {
		e.logger.Log(constant.ERROR, e.GetStockType(), 0.0, 0.0,
			fmt.Sprint("GetDepth error, the error number is %s"+depth.Message))
		return dbtypes.Depth{}, fmt.Errorf("GetDepth error, the error number not in backtest")
	}
	return depth.Data, nil
}

// Init ...
func (e *BaseExchange) Init(opt constant.Option) error {
	e.logger = model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type, Back: opt.BackLog}
	e.option = opt
	e.limit = opt.Limit
	e.ch = make(chan [2]string)
	e.lastSleep = time.Now().UnixNano()
	e.recordsPeriodDbMap = map[string]int64{
		"M1":  dbconstant.Minute,
		"M5":  5 * dbconstant.Minute,
		"M15": 15 * dbconstant.Minute,
		"M30": 30 * dbconstant.Minute,
		"H1":  dbconstant.Hour,
		"H2":  2 * dbconstant.Hour,
		"H4":  4 * dbconstant.Hour,
		"D1":  dbconstant.Day,
		"W1":  dbconstant.Week,
	}
	e.currencyMap = make(map[string]float64)
	return nil
}

// Stop ...
func (e *BaseExchange) Stop() error {
	close(e.ch)
	return nil
}

// Start ...
func (e *BaseExchange) Start() error {
	return nil
}

// SetID set ID
func (e *BaseExchange) SetID(id int) {
	e.id = id
}

// GetID.father.get ID
func (e *BaseExchange) GetID() int {
	return e.id
}

// SetIO set IO mode
func (e *BaseExchange) SetIO(mode string) {
	e.ioMode = mode
}

// GetIO.father.get IO mode
func (e *BaseExchange) GetIO() string {
	return e.ioMode
}

// GetStockType ...
func (e *BaseExchange) GetSymbol(symbol string) (string, string) {
	vec := strings.Split(symbol, ".")
	if len(vec) < 2 {
		return vec[0], ""
	}
	return vec[0], vec[1]
}

// SetDirection set the limit calls amount per second of this exchange
func (e *BaseExchange) SetDirection(direction string) {
	e.direction = direction
}

// GetDirection set the limit calls amount per second of this exchange
func (e *BaseExchange) GetDirection() string {
	return e.direction
}

// SetMarginLevel set the limit calls amount per second of this exchange
func (e *BaseExchange) SetMarginLevel(lever float64) {
	e.lever = lever
}

// GetMarginLevel set the limit calls amount per second of this exchange
func (e *BaseExchange) GetMarginLevel() float64 {
	return e.lever
}

// GetStockType set the limit calls amount per second of this exchange
func (e *BaseExchange) GetStockType() string {
	return e.stockType
}

// SetStockType set the limit calls amount per second of this exchange
func (e *BaseExchange) SetStockType(stockType string) {
	e.stockType = stockType
}

// SetMinAmountMap ...
func (e *BaseExchange) SetMinAmountMap(m map[string]float64) {
	e.minAmountMap = m
}

// GetMinAmountMap ...
func (e *BaseExchange) GetMinAmountMap() map[string]float64 {
	return e.minAmountMap
}

// SetRecordsPeriodMap ...
func (e *BaseExchange) SetRecordsPeriodMap(m map[string]int64) {
	e.recordsPeriodMap = m
}

// GetRecordsPeriodMap ...
func (e *BaseExchange) GetRecordsPeriodMap() map[string]int64 {
	return e.recordsPeriodMap
}

// GetRecords ...
func (e *BaseExchange) GetRecords() ([]constant.Record, error) {
	stockType := e.GetStockType()
	io := e.GetIO()
	refresh := false
	if io == constant.IOBLOCK {
		refresh = true
	}
	if io == constant.IOCACHE || io == constant.IOBLOCK {
		val := e.GetCache(constant.CacheRecord, e.GetStockType(), refresh)
		if val.Data == nil {
			return nil, fmt.Errorf("record not load")
		}
		return val.Data.([]constant.Record), nil
	}
	return e.father.getRecords(stockType)
}

func (e *BaseExchange) isRefresh() bool {
	io := e.GetIO()
	refresh := false
	if io == constant.IOBLOCK {
		refresh = true
	}
	return refresh
}

// GetTicker.father.get market ticker
func (e *BaseExchange) GetTicker() (*constant.Ticker, error) {
	stockType := e.GetStockType()
	io := e.GetIO()
	if io == constant.IOCACHE || io == constant.IOBLOCK {
		e.wait(stockType, constant.CacheTicker)
		val := e.GetCache(constant.CacheTicker, e.GetStockType(), e.isRefresh())
		if val.Data == nil {
			return nil, fmt.Errorf("ticker not load ")
		}
		dst := val.Data.(constant.Ticker)
		return &dst, nil
	}
	return e.father.getTicker(stockType)
}

// GetDepth.father.get depth from exchange
func (e *BaseExchange) GetDepth() (*constant.Depth, error) {
	stockType := e.GetStockType()
	io := e.GetIO()
	if io == constant.IOCACHE || io == constant.IOBLOCK {
		val := e.GetCache(constant.CacheDepth, stockType, e.isRefresh())
		if val.Data == nil {
			return nil, fmt.Errorf("depth not load ")
		}
		dst := val.Data.(constant.Depth)
		return &dst, nil
	}
	return e.father.getDepth(stockType)
}

// GetOrder ...
func (e *BaseExchange) GetOrder(id string) (*constant.Order, error) {
	stockType := e.GetStockType()
	io := e.GetIO()
	if io == constant.IOCACHE || io == constant.IOBLOCK {
		orders, err := e.GetOrders()
		if err != nil {
			return nil, err
		}
		for _, order := range orders {
			if order.Id == id {
				return &order, nil
			}
		}
		return nil, fmt.Errorf("order not found")
	}
	return e.father.getOrder(stockType, id)
}

// GetOrders.father.get all unfilled orders
func (e *BaseExchange) GetOrders() ([]constant.Order, error) {
	stockType := e.GetStockType()
	io := e.GetIO()
	if io == constant.IOCACHE || io == constant.IOBLOCK {
		e.wait(stockType, constant.CacheOrder)

		val := e.GetCache(constant.CacheOrder, e.GetStockType(), e.isRefresh())
		if val.Data == nil {
			return nil, fmt.Errorf("account not load")
		}
		dst := val.Data.([]constant.Order)
		return dst, nil
	}
	return e.father.getOrders(stockType)
}

// GetAccount ...
func (e *BaseExchange) GetAccount() (*constant.Account, error) {
	io := e.GetIO()
	if io == constant.IOCACHE || io == constant.IOBLOCK {
		e.wait("", constant.CacheAccount)
		val := e.GetCache(constant.CacheAccount, e.GetStockType(), e.isRefresh())
		if val.Data == nil {
			return nil, fmt.Errorf("account not load")
		}
		dst := val.Data.(constant.Account)
		return &dst, nil
	}
	return e.father.getAccount()
}

// GetPosition.father.get position from exchange
func (e *BaseExchange) GetPosition() ([]constant.Position, error) {
	stockType := e.GetStockType()
	io := e.GetIO()
	if io == constant.IOCACHE || io == constant.IOBLOCK {
		e.wait(stockType, constant.CachePosition)
		val := e.GetCache(constant.CachePosition, e.GetStockType(), e.isRefresh())
		if val.Data == nil {
			return nil, fmt.Errorf("position not load ")
		}
		return val.Data.([]constant.Position), nil
	}
	return e.father.getPosition(stockType)
}

// ValidBuy ...
func (e *BaseExchange) ValidBuy() error {
	dir := e.GetDirection()
	if dir == constant.TradeTypeBuy {
		return nil
	}
	if dir == constant.TradeTypeShortClose {
		return nil
	}
	return errors.New("buy direction error:" + e.GetDirection())
}

// ValidSell ...
func (e *BaseExchange) ValidSell() error {
	dir := e.GetDirection()
	if dir == constant.TradeTypeSell {
		return nil
	}
	if dir == constant.TradeTypeLongClose {
		return nil
	}
	return errors.New("sell direction error:" + e.GetDirection())
}

// Log print something to console
func (e *BaseExchange) Log(action, symbol string, price, amount float64, messages string) {
	e.logger.Log(action, symbol, price, amount, messages)
}

// GetType.father.get the type of this exchange
func (e *BaseExchange) GetType() string {
	return e.option.Type
}

// GetName.father.get the name of this exchange
func (e *BaseExchange) GetName() string {
	return e.option.Name
}
