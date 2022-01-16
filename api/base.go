package api

import (
	"errors"
	"strings"
	"time"

	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
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
	datas []constant.OHLC
}

// Next ...
func (l *DataLoader) Next() *constant.OHLC {
	nextPos := l.curr + 1
	if nextPos >= l.size {
		return nil
	}
	data := l.datas[l.curr]
	l.curr = nextPos
	return &data
}

func (l *DataLoader) Progress() int {
	if l.size == 0 {
		return 100
	}
	v := float64(l.curr+1) / float64(l.size) * 100
	return int(v)
}

// Dump ...
func (l *DataLoader) Dump() []constant.OHLC {
	return l.datas
}

// Load ...
func (l *DataLoader) Load(ohlcs []constant.OHLC) {
	l.datas = append(l.datas, ohlcs...)
	l.size += len(l.datas)
	// move one for first records, at least one ohlc
	// l.Next()
}

// BaseExchange ...
type BaseExchange struct {
	period             string
	size               int
	id                 int     // id of the exchange
	ioMode             string  // io mode for exchange
	direction          string  // trade type
	stockType          string  // stockType
	lever              float64 // lever
	recordsPeriodMap   map[string]int64
	recordsPeriodDbMap map[string]int64 // coin watched
	// recordsPeriod support
	minAmountMap map[string]float64 // minAmount of trade
	limit        int64
	lastSleep    int64
	lastTimes    int64
	currencyMap  map[string]float64

	coverRate    float64
	taker        float64
	maker        float64
	contractRate float64 // 合约每张价值

	host string

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
	orderID, err := e.father.buy(price, amount, msg)
	if err != nil {
		return "", err
	}
	return orderID, nil
}

func (e *BaseExchange) Sell(price, amount, msg string) (string, error) {
	orderID, err := e.father.sell(price, amount, msg)
	if err != nil {
		return "", err
	}
	return orderID, nil
}

func (e *BaseExchange) CancelOrder(orderID string) (bool, error) {
	status, err := e.father.cancelOrder(orderID)
	if err != nil {
		return false, err
	}
	return status, nil
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
func (e *BaseExchange) SetBackCommission(taker, maker, contractRate, coverRate float64) {
	e.contractRate = contractRate
	e.taker = taker
	e.maker = maker
	e.coverRate = coverRate
}

// GetBackCommission 获取回测手续费
func (e *BaseExchange) GetBackCommission() []float64 {
	return []float64{e.taker, e.maker, e.contractRate, e.coverRate}
}

// GetBackAccount ...
func (e *BaseExchange) GetBackAccount() map[string]float64 {
	return e.currencyMap
}

// SetBackAccount ...
func (e *BaseExchange) SetBackAccount(key string, val float64) {
	e.currencyMap[key] = val
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

// Set Period
func (e *BaseExchange) SetPeriod(period string) {
	e.period = period
}

// Get Period
func (e *BaseExchange) GetPeriod() string {
	return e.period
}

func (e *BaseExchange) Start() error {
	return e.father.start()
}

func (e *BaseExchange) Stop() error {
	return e.father.stop()
}

// Init ...
func (e *BaseExchange) Init(opt constant.Option) error {
	e.logger = model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type, Back: opt.BackLog}
	e.option = opt
	e.limit = opt.Limit
	e.lastSleep = time.Now().UnixNano()
	e.recordsPeriodDbMap = map[string]int64{
		"M1":  constant.Minute,
		"M5":  5 * constant.Minute,
		"M15": 15 * constant.Minute,
		"M30": 30 * constant.Minute,
		"H1":  constant.Hour,
		"H2":  2 * constant.Hour,
		"H4":  4 * constant.Hour,
		"D1":  constant.Day,
		"W1":  constant.Week,
	}
	e.SetPeriodSize(constant.RecordSize)
	e.currencyMap = make(map[string]float64)
	return nil
}

func (e *BaseExchange) SetID(id int) {
	e.id = id
}

// GetStockType ...
func (e *BaseExchange) getSymbol(symbol string) (string, string) {
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
	return e.father.getRecords(e.GetStockType())
}

// GetTicker  market ticker
func (e *BaseExchange) GetTicker() (*constant.Ticker, error) {
	return e.father.getTicker(e.GetStockType())
}

// GetDepth.father.get depth from exchange
func (e *BaseExchange) GetDepth() (*constant.Depth, error) {
	return e.father.getDepth(e.GetStockType())
}

// GetOrder ...
func (e *BaseExchange) GetOrder(id string) (*constant.Order, error) {
	return e.father.getOrder(e.GetStockType(), id)
}

// GetOrders.father.get all unfilled orders
func (e *BaseExchange) GetOrders() ([]constant.Order, error) {
	return e.father.getOrders(e.GetStockType())
}

// GetAccount ...
func (e *BaseExchange) GetAccount() (*constant.Account, error) {
	return e.father.getAccount()
}

// GetPosition.father.get position from exchange
func (e *BaseExchange) GetPosition() ([]constant.Position, error) {
	return e.father.getPosition(e.GetStockType())
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
