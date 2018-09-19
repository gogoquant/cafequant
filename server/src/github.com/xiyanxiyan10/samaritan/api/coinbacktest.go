package api

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	goback "github.com/dirkolbrich/gobacktest"
	"github.com/xiyanxiyan10/samaritan/constant"
	"github.com/xiyanxiyan10/samaritan/conver"
	"github.com/xiyanxiyan10/samaritan/model"
	"strconv"
	"strings"
	"time"
)

func init() {
	constructor["coinbacktest"] = NewBacktest
}

// Backtest backtest struct
type BtBacktest struct {
	goback.Backtest

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
}

// BacktestData online back test support
type BacktestData interface {
	getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error)
}

// NewBacktest create a backtest
func NewBacktest(opt Option) Exchange {
	back := BtBacktest{logger: model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option: opt,

		limit:     10.0,
		lastSleep: time.Now().UnixNano()}

	// @todo check if exists
	maker, _ := constructor[opt.Type]

	back.exchangeHandler = maker(opt)
	return &back
}

// Log print something to console
func (e *BtBacktest) Log(msgs ...interface{}) {

}

// GetType get the type of this exchange
func (e *BtBacktest) GetType() string {
	return ""
}

// GetName get the name of this exchange
func (e *BtBacktest) GetName() string {
	return ""
}

// SetLimit set the limit calls amount per second of this exchange
func (e *BtBacktest) SetLimit(times interface{}) float64 {
	return 0.0
}

// AutoSleep auto sleep to achieve the limit calls amount per second of this exchange
func (e *BtBacktest) AutoSleep() {

}

// GetMinAmount get the min trade amonut of this exchange
func (e *BtBacktest) GetMinAmount(stock string) float64 {
	return 0.0
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
	stockType = strings.ToUpper(stockType)
	tradeType = strings.ToUpper(tradeType)
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
	//Todo time and symbol test set
	event := &goback.Event{}
	event.SetTime(time.Now())
	event.SetSymbol(stockType)
	signal := &goback.Signal{
		Event: *event,
	}
	signal.SetPrice(price)
	signal.SetFQty(fqty)
	signal.SetQtyType(goback.FLOAT64_QTY)
	if price < 0 {
		signal.SetOrderType(goback.MarketOrder)
	} else {
		signal.SetOrderType(goback.LimitOrder)
	}
	e.AddSignal(signal)
	e.logger.Log(constant.BUY, stockType, price, fqty, msgs...)
	return fmt.Sprint("%s", "success")
}

func (e *BtBacktest) sell(stockType string, price, fqty float64, msgs ...interface{}) interface{} {
	//Todo time and symbol test set
	event := &goback.Event{}
	event.SetTime(time.Now())
	event.SetSymbol(stockType)
	signal := &goback.Signal{
		Event: *event,
	}
	signal.SetPrice(price)
	signal.SetFQty(fqty)
	signal.SetQtyType(goback.FLOAT64_QTY)
	if price < 0 {
		signal.SetOrderType(goback.MarketOrder)
	} else {
		signal.SetOrderType(goback.LimitOrder)
	}
	e.AddSignal(signal)
	e.logger.Log(constant.SELL, stockType, price, fqty, msgs...)
	return fmt.Sprint("%s", "success")
}

// GetOrder get details of an order
func (e *BtBacktest) GetOrder(stockType, id string) interface{} {
	backorders, ok := e.OrdersBySymbol(stockType)
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
			// @todo more trade type
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
	backorders, ok := e.OrdersBySymbol(stockType)
	if !ok {
		return false
	}

	for _, backorder := range backorders {
		order := Order{
			ID:         strconv.Itoa(backorder.ID()),
			Price:      backorder.Price(),
			Amount:     backorder.FQty(),
			DealAmount: backorder.FQty(),
			StockType:  stockType,
		}
		// @todo more trade type
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
	return false
}

// CancelOrder cancel an order
func (e *BtBacktest) CancelOrder(order Order) bool {
	id, err := strconv.Atoi(order.ID)
	if err != nil {
		e.logger.Log(constant.ERROR, order.StockType, 0.0, 0.0, err)
	}
	e.CancelOneOrder(id)
	return false
}

// getTicker get market ticker & depth
func (e *BtBacktest) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, end bool, err error) {
	// @Todo convert data to tick
	_, end, err = e.Run2Data()
	if err != nil {
		return ticker, false, err
	}
	if end == true {
		return ticker, true, errors.New("backTest end")
	}
	return
}

// GetTicker get market ticker & depth
func (e *BtBacktest) GetTicker(stockType string, sizes ...interface{}) interface{} {
	if e.option.Mode == constant.OFFLINE {
		ticker, end, err := e.getTicker(stockType, sizes...)
		if err != nil {
			if end {
				e.logger.Log(constant.INFO, stockType, 0.0, 0.0, err)
			} else {
				e.logger.Log(constant.ERROR, stockType, 0.0, 0.0, err)
			}
			return false
		}
		return ticker
	}
	ticker := e.exchangeHandler.GetTicker(stockType, sizes...)
	//process all event every tick
	err := e.Run2Event()
	if err != nil {
		e.logger.Log(constant.ERROR, stockType, 0.0, 0.0, err)
	}
	return ticker
}

// GetRecords get candlestick data
func (e *BtBacktest) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	return false
}
