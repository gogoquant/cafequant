package api

import (
	"github.com/bitly/go-simplejson"
	goback "github.com/xiyanxiyan10/gobacktest"
	"github.com/xiyanxiyan10/samaritan/constant"
	"github.com/xiyanxiyan10/samaritan/conver"
	"github.com/xiyanxiyan10/samaritan/model"
	"strconv"
	"time"
)

func init() {
	constructor["coinbacktest"] = NewCoinBacktest
}

// Backtest backtest struct
type BtBacktest struct {
	// one Trader pointer to one backtest with some exchanges
	back *goback.Backtest

	stockTypeMap     map[string]string
	tradeTypeMap     map[string]string
	recordsPeriodMap map[string]string
	minAmountMap     map[string]float64
	records          map[string][]Record

	logger model.Logger
	option Option

	limit     float64
	lastSleep int64
	lastTimes int64

	exchangeHandler ExchangeHandler
}

// ExchangeHandler api used for backtest
type ExchangeHandler interface {
	GetTicker(stockType string, sizes ...interface{}) interface{}
	SetLimit(times interface{}) float64
	GetMinAmount(stock string) float64
	Start(back *goback.Backtest) error
	Stop(back *goback.Backtest) error
	Status() int
	AutoSleep()
}

// NewCoinBacktest create a coin backtest
func NewCoinBacktest(opt Option) Exchange {
	back := BtBacktest{logger: model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:       opt,
		limit:        10.0,
		lastSleep:    time.Now().UnixNano(),
		stockTypeMap: map[string]string{},
	}

	maker, ok := constructor[opt.Type]
	if !ok {
		return nil
	}
	handler := maker(opt)
	back.SetStockMap(handler.StockMap())

	//back.SetPortfolio(goback.NewPortfolio())
	//back.SetMarry(&BtMarry{})
	back.exchangeHandler = maker(opt)
	return &back
}

// StockMap ...
func (e *BtBacktest) StockMap() map[string]string {
	return e.stockTypeMap
}

// SetStockMap ...
func (e *BtBacktest) SetStockMap(m map[string]string) {
	e.stockTypeMap = m
}

// SetGoback ...
func (e *BtBacktest) SetGoback(back *goback.Backtest) {
	e.back = back
}

// Log print something to console
func (e *BtBacktest) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType get the type of this exchange
func (e *BtBacktest) GetType() string {
	return e.option.Type
}

// GetName get the name of this exchange
func (e *BtBacktest) GetName() string {
	return e.option.Name
}

// SetLimit set the limit calls amount per second of this exchange
func (e *BtBacktest) SetLimit(times interface{}) float64 {
	return e.exchangeHandler.SetLimit(times)
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *BtBacktest) AutoSleep() {
	e.exchangeHandler.AutoSleep()
}

// GetMinAmount get the min trade amonut of this exchange
func (e *BtBacktest) GetMinAmount(stock string) float64 {
	return e.exchangeHandler.GetMinAmount(stock)
}

// getAuthJSON
func (e *BtBacktest) getAuthJSON(method string, params ...interface{}) (jsoner *simplejson.Json, err error) {
	data := []byte{}
	return simplejson.NewJson(data)
}

// GetAccount get the account detail of this exchange
func (e *BtBacktest) GetAccount() interface{} {
	return map[string]float64{}
}

// Trade place an order
func (e *BtBacktest) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	price := conver.Float64Must(_price)
	amount := conver.Float64Must(_amount)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, stockType, 0.0, 0.0, "Trade() error, unrecognized stockType: ", stockType)
		return false
	}
	switch tradeType {
	case constant.TradeTypeBuy:
		return e.buy(stockType, price, amount, msgs...)
	case constant.TradeTypeSell:
		return e.sell(stockType, price, amount, msgs...)
	default:
		e.logger.Log(constant.ERROR, stockType, 0.0, 0.0, "Trade() error, unrecognized tradeType: ", tradeType)
		return false
	}
}

func (e *BtBacktest) buy(stockType string, price, fqty float64, msgs ...interface{}) interface{} {
	event := &goback.Event{}
	event.SetTime(time.Now())
	event.SetSymbol(stockType)
	signal := &goback.Signal{
		Event: *event,
	}
	signal.SetFQty(fqty)
	signal.SetQtyType(goback.FLOAT64_QTY)
	if price < 0 {
		signal.SetOrderType(goback.MarketOrder)
	} else {
		signal.SetPrice(price)
		signal.SetOrderType(goback.LimitOrder)
	}
	e.back.AddSignal(signal)
	e.logger.Log(constant.BUY, stockType, price, fqty, msgs...)
	return "success"
}

func (e *BtBacktest) sell(stockType string, price, fqty float64, msgs ...interface{}) interface{} {
	event := &goback.Event{}
	event.SetTime(time.Now())
	event.SetSymbol(stockType)
	signal := &goback.Signal{
		Event: *event,
	}
	signal.SetFQty(fqty)
	signal.SetQtyType(goback.FLOAT64_QTY)
	if price < 0 {
		signal.SetOrderType(goback.MarketOrder)
	} else {
		signal.SetPrice(price)
		signal.SetOrderType(goback.LimitOrder)
	}
	e.back.AddSignal(signal)
	e.logger.Log(constant.SELL, stockType, price, fqty, msgs...)
	return "success"
}

// GetOrder get details of an order
func (e *BtBacktest) GetOrder(stockType, id string) interface{} {
	backorders, ok := e.back.OrdersBySymbol(stockType)
	if !ok {
		return false
	}

	for _, backorder := range backorders {
		if strconv.Itoa(backorder.ID()) == id {

			order := Order{
				ID:         strconv.Itoa(backorder.ID()),
				Price:      backorder.Price(),
				Amount:     backorder.FQty(),
				DealAmount: backorder.FQty(),

				StockType: stockType,
			}
			if backorder.Direction() == 0 {
				order.StockType = constant.TradeTypeBuy
			} else {
				order.StockType = constant.TradeTypeSell
			}
			return order
		}

	}
	return false
}

// GetOrders get all unfilled orders
func (e *BtBacktest) GetOrders(stockType string) interface{} {
	orders := []Order{}
	backorders, ok := e.back.OrdersBySymbol(stockType)
	if !ok {
		return false
	}

	for _, backorder := range backorders {
		if backorder.Status() != goback.OrderNew {
			continue
		}
		order := Order{
			ID:         strconv.Itoa(backorder.ID()),
			Price:      backorder.Price(),
			Amount:     backorder.FQty(),
			DealAmount: backorder.FQty(),
			StockType:  stockType,
		}
		if backorder.Direction() == 0 {
			order.StockType = constant.TradeTypeBuy
		} else {
			order.StockType = constant.TradeTypeSell
		}
		orders = append(orders, order)
	}
	return orders
}

// GetTrades get all filled orders recently
func (e *BtBacktest) GetTrades(stockType string) interface{} {
	orders := []Order{}
	backorders, ok := e.back.OrdersBySymbol(stockType)
	if !ok {
		return false
	}

	for _, backorder := range backorders {
		if backorder.Status() != goback.OrderSubmitted {
			continue
		}
		order := Order{
			ID:         strconv.Itoa(backorder.ID()),
			Price:      backorder.Price(),
			Amount:     backorder.FQty(),
			DealAmount: backorder.FQty(),
			StockType:  stockType,
		}
		if backorder.Direction() == 0 {
			order.StockType = constant.TradeTypeBuy
		} else {
			order.StockType = constant.TradeTypeSell
		}
		orders = append(orders, order)
	}
	return orders
}

// CancelOrder cancel an order
func (e *BtBacktest) CancelOrder(order Order) bool {
	id, err := strconv.Atoi(order.ID)
	if err != nil {
		e.logger.Log(constant.ERROR, order.StockType, 0.0, 0.0, err)
	}
	e.back.CancelOrder(id)
	return false
}

// EnableSubscribe ...
func (e *BtBacktest) EnableSubscribe(symbol string) error {
	return e.back.Portfolio().EnableSubscribe(symbol)
}

// DisableSubscribe ...
func (e *BtBacktest) DisableSubscribe(symbol string) error {
	return e.back.Portfolio().DisableSubscribe(symbol)
}

// Draw point
func (e *BtBacktest) Draw(val map[string]interface{}) interface{} {

	if e.back.DataGram() == nil {
		return false
	}

	datagram := goback.NewDataGram()
	datagram.SetTime(time.Now())
	datagram.SetFields(val)
	datagram.SetId("data_" + strconv.FormatInt(e.option.TraderID, 10))
	e.back.AddEvent(datagram)
	return true
}

// SetDrawMode
func (e *BtBacktest)SetDrawMode(mode string)interface{}{
	if e.back.DataGram() == nil {
		return false
	}
	datagram := goback.NewDataGram()
	datagram.SetSymbol(mode)
	e.back.AddEvent(datagram)
	return true
}

// GetTicker get market ticker & depth
func (e *BtBacktest) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker := e.exchangeHandler.GetTicker(stockType, sizes...)
	if ticker == false {
		return false
	}

	return ticker
}

// GetRecords get candlestick data
func (e *BtBacktest) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return false
}

// Start
func (bt *BtBacktest) Start(back *goback.Backtest) error {
	return bt.exchangeHandler.Start(back)
}

// Stop
func (bt *BtBacktest) Stop(back *goback.Backtest) error {
	return bt.exchangeHandler.Stop(back)
}

// Status
func (bt *BtBacktest) Status() int {
	return bt.exchangeHandler.Status()
}
