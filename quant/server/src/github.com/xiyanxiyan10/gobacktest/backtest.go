package gobacktest

import (
	"fmt"
	"gopkg.in/logger.v1"
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

	// add Data
	AddData(data DataEvent) error

	// Next run engine
	Next() error
}

// Reseter provides a resting interface.
type Reseter interface {
	Reset() error
}

// BackTest 回测管理
type BackTest struct {
	config map[string]string

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
	orders, _ := back.exchange.OrdersBySymbol(symbol)
	return orders, nil
}

// AddOrder
func (back *BackTest) AddOrder(order Order) error {
	order.SetStatus(OrderNew)
	back.AddEvent(&order)
	return nil
}

// AddData ...
func (back *BackTest) AddData(data DataEvent) error {
	back.AddEvent(data)
	return nil
}

// Holds get all position
func (back *BackTest) Holds() (map[string]Position, error) {
	position := back.holds()
	return position, nil
}

// SetCash ...
func (back *BackTest) SetCash(cash float64) {
	back.Portfolio().SetInitialCash(cash)
	back.Portfolio().SetCash(cash)
}

// Status ...
func (back *BackTest) Status() int64 {
	return back.status
}

// NewBackTest
func NewBackTest(m map[string]string) Back {
	bt := &BackTest{
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
func (back *BackTest) holds() map[string]Position {
	m := back.portfolio.Holds()
	var p Position
	p.fqty = back.Portfolio().Cash()
	m["cash"] = p
	return m
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
	back.eventQueue = append(back.eventQueue, e)
	return nil
}

// Next run backTest process
func (back *BackTest) Next() error {
	for {
		err, end := back.activeEvent()
		if err != nil {
			return err
		}
		if end {
			return nil
		}
	}
}

// activeEvent directs the different events to their handler.
func (back *BackTest) activeEvent() (err error, end bool) {
	e, end := back.nextEvent()
	if end {
		return nil, end
	}

	// type check for event type
	switch event := e.(type) {

	case CmdEvent:
		log.Infof("Get cmd event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())
		err = nil
		break

	case DataEvent:
		log.Infof("Get data event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())
		//@todo Marry order here
		break
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

		// 关闭订单
		if event.status == OrderCancel {
			back.exchange.CancelOrder(*event)
			break
		}

		// 下订单
		if event.status == OrderNew {
			// 这里需要检查卖单的仓位和买单的资金
			back.exchange.AddOrder(event)
			if err != nil {
				break
			}
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
	return nil, false
}
