package gobacktest

import (
	"errors"
	"fmt"
	"gopkg.in/logger.v1"
	"time"
)

var errHalt = fmt.Errorf("HALT")

// Back backTest engine all Api support
type Back interface {

	// set the name of the engine
	Name() string

	// set the name of the engine
	SetName(name string)

	// get the holds of the engine
	Holds() (map[string]Position, error)

	// get order from the symbol
	OrdersBySymbol(symbol string) ([]OrderEvent, error)

	// cancel the order
	CancelOrder(order Order) error

	// add order
	AddOrder(order Order) error
}

// Reseter provides a resting interface.
type Reseter interface {
	Reset() error
}

// BackTest 回测管理
type BackTest struct {
	Id string

	config map[string]string

	in  chan EventHandler
	out chan ResultEvent

	status int64
	name   string

	symbols []string

	data DataHandler

	portfolio PortfolioHandler
	exchange  ExchangeHandler
	statistic StatisticHandler

	eventQueue []EventHandler
}

// SignalAdd Add signal event into event queue
func (back *BackTest) AddSignal(signals ...SignalEvent) error {
	for _, signal := range signals {
		back.AddEvent(signal)
	}
	return nil
}

// Stats returns the statistic handler of the backtest.
func (back *BackTest) Stats() StatisticHandler {
	return back.statistic
}

// OrdersBySymbol get order by symbol
func (back *BackTest) OrdersBySymbol(symbol string) ([]OrderEvent, error) {
	var order Order
	order.SetSymbol(symbol)
	order.SetStatus(OrdersBySymbol)
	back.AddEvent(&order)
	res, err := back.getResult()
	if err != nil {
		return nil, err
	}
	orders, err := res.Data().([]OrderEvent)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// EOrdersBySymbol used just for engine
func (back *BackTest) EOrdersBySymbol(symbol string) ([]OrderEvent, error) {
	orders, _ := back.exchange.OrdersBySymbol(symbol)
	return orders, nil
}

// AddOrder
func (back *BackTest) AddOrder(order Order) error {
	order.SetStatus(OrderNew)
	back.AddEvent(&order)
	res, err := back.getResult()
	if err != nil {
		return err
	}
	err = res.Data().(error)
	if err != nil {
		return err
	}
	return nil
}

// Holds get all position
func (back *BackTest) Holds() (map[string]Position, error) {
	var order Order
	order.SetStatus(OrderFillByAll)
	back.AddEvent(&order)
	res, err := back.getResult()
	if err != nil {
		return nil, err
	}
	position, err := res.Data().(map[string]Position)
	if err != nil {
		return nil, err
	}
	return position, nil
}

// SetCash ...
func (back *BackTest) SetCash(cash float64) {
	back.Portfolio().SetInitialCash(cash)
	back.Portfolio().SetCash(cash)
}

// getResult ...
func (back *BackTest) getResult() (ResultEvent, error) {
	select {
	case data, _ := <-back.out:
		return data, nil
	case <-time.After(1 * time.Second):
		return nil, errors.New("TimeOut")
	}
}

// Status ...
func (back *BackTest) Status() int64 {
	return back.status
}

// NewBackTest
func NewBackTest(m map[string]string) Back {
	bt := &BackTest{
		in:     make(chan EventHandler, 50),
		out:    make(chan ResultEvent, 50),
		status: 0,
		config: m,
		portfolio: &Portfolio{
			initialCash: 0,
			cash:        0,
			sizeManager: &Size{DefaultSize: 100, DefaultValue: 1000},
			riskManager: &Risk{},
		},
		exchange: &Exchange{
			//Symbol: "TEST",
			//Commission:  &FixedCommission{Commission: 0},
			//ExchangeFee: &FixedExchangeFee{ExchangeFee: 0},
		},
		statistic: &Statistic{},
	}
	err := bt.initialize()
	if err != nil {
		return nil
	} else {
		return bt
	}
}

// Name ...
func (back *BackTest) Name() string {
	return back.name
}

// SetName ...
func (back *BackTest) SetName(name string) {
	back.name = name
}

// CommitOrder ...
func (back *BackTest) CommitOrder(id int) error {
	fill, err := back.exchange.CommitOrder(id)
	if err == nil && fill != nil {
		back.AddEvent(fill)
	}
	return err
}

// CancelOrder ...
func (back *BackTest) CancelOrder(order Order) error {
	order.SetStatus(OrderCancel)
	back.AddEvent(&order)
	return nil
}

// initialize ...
func (back *BackTest) initialize() (err error) {
	return
}

// holds
func (back *BackTest) holds() (res Result) {
	m := back.portfolio.Holds()
	var p Position
	p.fqty = back.Portfolio().Cash()
	m["cash"] = p
	res.SetData(m)
	return
}

// Cmd ...
func (back *BackTest) Cmd(cmd string) error {
	var event Cmd
	event.SetCmd(cmd)
	back.AddEvent(&event)
	return nil
}

// SetSymbols sets the symbols to include into the backtest.
func (back *BackTest) SetSymbols(symbols []string) {
	back.symbols = symbols
}

// SetData sets the data provider to be used within the backtest.
func (back *BackTest) SetData(data DataHandler) {
	back.data = data
}

// SetPortfolio sets the portfolio provider to be used within the backtest.
func (back *BackTest) SetPortfolio(portfolio PortfolioHandler) {
	back.portfolio = portfolio
}

// SetExchange sets the execution provider to be used within the backtest.
func (back *BackTest) SetExchange(exchange ExchangeHandler) {
	back.exchange = exchange
}

// SetStatistic sets the statistic provider to be used within the backtest.
func (back *BackTest) SetStatistic(statistic StatisticHandler) {
	back.statistic = statistic
}

// Portfolio sets the Portfolio provider to be used within the backtest.
func (back *BackTest) Portfolio() PortfolioHandler {
	return back.portfolio
}

func (back *BackTest) Exchange() ExchangeHandler {
	return back.exchange
}

// Reset ...
func (back *BackTest) Reset() error {
	back.eventQueue = nil
	back.data.Reset()
	back.portfolio.Reset()
	back.statistic.Reset()
	return nil
}

// NextEvent process the event before return the data event back to user' scripts
func (back *BackTest) NextEvent() (err error, status string, data DataEvent) {
	event := <-back.in
	return back.activeEvent(event)
}

// setup runs at the beginning of the backtest to perfom preparing operations.
func (back *BackTest) setup() error {
	// before first run, set portfolio cash
	back.portfolio.SetCash(back.portfolio.InitialCash())
	return nil
}

// teardown performs any cleaning operations at the end of the backtest.
func (back *BackTest) teardown() error {
	// no implementation yet
	return nil
}

// nextEvent gets the next event from the events queue.
func (back *BackTest) nextEvent() (e EventHandler, ok bool) {

	// if event queue empty return false
	if len(back.eventQueue) == 0 {
		return e, false
	}

	// return first element from the event queue
	e = back.eventQueue[0]
	back.eventQueue = back.eventQueue[1:]
	return e, true
}

// AddEvent
func (back *BackTest) AddEvent(e EventHandler) error {
	back.in <- e
	return nil
}

// activeEvent directs the different events to their handler.
func (back *BackTest) activeEvent(e EventHandler) (err error, status string, data DataEvent) {
	status = "continue"

	// type check for event type
	switch event := e.(type) {

	case CmdEvent:
		log.Infof("Get cmd event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())
		err = nil
		status = "end"
		break

	case DataEvent:
		log.Infof("Get data event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())
		return nil, "data", event
		// update portfolio to the last known price data
		//t.portfolio.Update(event)
		// update statistics
		//t.statistic.Update(event, t.portfolio)
		// check if any orders are filled before proceding
		//t.exchange.OnData(event)

	case *Signal:
		log.Infof("Get signal event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())
		order, err := back.portfolio.OnSignal(event, back.data)
		if err != nil {
			break
		}
		back.AddEvent(order)

	case *Order:
		log.Infof("Get order event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())

		/*
			if event.status == OrdersBySymbol {
				orders, _ := back.exchange.OrdersBySymbol(event.symbol)
				var result Result
				result.SetTime(event.timestamp)
				result.SetData(orders)
				var rs ResultEvent
				rs = &result
				back.out <- ResultEvent(rs)
				break
			}
		*/

		// 关闭订单
		if event.status == OrderCancel {
			back.exchange.CancelOrder(*event)
			break
		}

		// 获取仓位
		/*
			if event.status == OrderFillByAll {
				fills := back.holds()
				back.out <- &fills
				break
			}
		*/

		// 下订单
		if event.status == OrderNew {
			// 这里需要检查卖单的仓位和买单的资金
			back.exchange.AddOrder(event)
			if err != nil {
				break
			}
			var res Result
			res.SetData("success")
			back.out <- &res
			break
		}

		// 未知行为
		log.Infof("UnKnow order type")
		break
		//t.AddEvent(fill)

	case *Fill:

		log.Infof("Get fill event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())
		back.exchange.OnFill(event)
		_, err := back.portfolio.OnFill(event)
		if err != nil {
			break
		}
		//t.AddEvent(transaction)
	}
	return
}
