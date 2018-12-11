package gobacktest

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"gopkg.in/logger.v1"
)

var errHalt       = fmt.Errorf("HALT")

// Back Manager handler for back
type Back interface {
	Name() string
	SetName(name string)
	Start() (err error)
	Status() (int64)
	Stop() (err error)
	SetScripts(scripts string)

	AddEvent(e EventHandler) error
	OrdersBySymbol(stockType string) ([]OrderEvent, bool)
	CancelOrder(id int) error
	Cmd(cmd string) error

}

// ScriptsApi api for back scripts
type ScriptsApi interface {
	CommitOrder(id int) (error)
	GetEvent() (err error, status string, data DataEvent)
}

// Reseter provides a resting interface.
type Reseter interface {
	Reset() error
}

// Backtest is the main struct which holds all back event for users
type Backtest struct {
	Id string

	config map[string]string

	in chan EventHandler
	out chan ResultEvent

	status  int64
	name    string

	symbols []string

	data DataHandler

	portfolio PortfolioHandler
	exchange  ExchangeHandler
	statistic StatisticHandler

	eventQueue []EventHandler
	scripts string

	Ctx    *otto.Otto
}

// initialize ...
func (back *Backtest)initialize() (err error) {
	back.Ctx = otto.New()
	back.Ctx.Interrupt = make(chan func(), 1)
	back.Ctx.Set("Exchange", ScriptsApi(back))
	return
}

// SetScripts ...
func (back *Backtest)SetScripts(scripts string){
	back.scripts = scripts
}

// Start ...
func (back *Backtest)Start() (err error) {
	go func() {
		defer func() {
			if err := recover(); err != nil && err != errHalt {
				log.Error(err)
			}
			if exit, err := back.Ctx.Get("exit"); err == nil && exit.IsFunction() {
				if _, err := exit.Call(exit); err != nil {
					log.Error(err)
				}
			}
			back.status = 0
		}()

		back.status = 1
		if _, err := back.Ctx.Run(back.scripts); err != nil {
			log.Error(err)
		}
		if main, err := back.Ctx.Get("main"); err != nil || !main.IsFunction() {
			log.Error("Can not get the main function")
		} else {
			if _, err := main.Call(main); err != nil {
				log.Error(err)
			}
		}
	}()
	return
}

// SetId
func (b *Backtest)SetId(id string){
	b.Id = id
}

// Status ...
func (back *Backtest)Status() (int64) {
	return back.status
}

// Stop ...
func (back *Backtest)Stop() (err error) {
	back.Ctx.Interrupt <- func() { panic(errHalt) }
	return
}

// NewBackTest
func NewBackTest(m map[string]string) Back {
	bt := Backtest{
		in: make(chan EventHandler, 50),
		out: make(chan ResultEvent, 50),
		status:  0,
		config:  m,
	}
	err := bt.initialize()
	if err != nil{
		return nil
	}else{
		return &bt
	}
}

// Name ...
func (e *Backtest) Name() string {
	return e.name
}

// Name ...
func (e *Backtest) SetName(name string) {
	e.name = name
}

// CommitOrder ...
func (t *Backtest) CommitOrder(id int) (error) {
	fill, err := t.exchange.CommitOrder(id)
	if err == nil && fill != nil {
		t.AddEvent(fill)
	}
	return err
}

// OrdersBySymbol ...
func (t *Backtest) OrdersBySymbol(stockType string) ([]OrderEvent, bool) {
	return t.exchange.OrdersBySymbol(stockType)
}

// CancelOneOrder ...
func (t *Backtest) CancelOrder(id int) error {
	return t.exchange.CancelOrder(id)
}

// Cmd ...
func (t *Backtest) Cmd(cmd string) error {
	var event Cmd
	event.SetCmd(cmd)
	t.AddEvent(&event)
	return nil
}

// New creates a default backtest with sensible defaults ready for use.
func New() *Backtest {
	return &Backtest{
		portfolio: &Portfolio{
			initialCash: 100000,
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
}

// SetSymbols sets the symbols to include into the backtest.
func (t *Backtest) SetSymbols(symbols []string) {
	t.symbols = symbols
}

// SetData sets the data provider to be used within the backtest.
func (t *Backtest) SetData(data DataHandler) {
	t.data = data
}

// SetPortfolio sets the portfolio provider to be used within the backtest.
func (t *Backtest) SetPortfolio(portfolio PortfolioHandler) {
	t.portfolio = portfolio
}

// SetExchange sets the execution provider to be used within the backtest.
func (t *Backtest) SetExchange(exchange ExchangeHandler) {
	t.exchange = exchange
}

// SetStatistic sets the statistic provider to be used within the backtest.
func (t *Backtest) SetStatistic(statistic StatisticHandler) {
	t.statistic = statistic
}

// Portfolio sets the Portfolio provider to be used within the backtest.
func (t *Backtest) Portfolio() PortfolioHandler {
	return t.portfolio
}

func (t *Backtest) Exchange() ExchangeHandler {
	return t.exchange
}

// Reset ...
func (t *Backtest) Reset() error {
	t.eventQueue = nil
	t.data.Reset()
	t.portfolio.Reset()
	t.statistic.Reset()
	return nil
}

// SignalAdd Add signal event into event queue
func (t *Backtest) AddSignal(signals ...SignalEvent) error {
	for _, signal := range signals {
		t.AddEvent(signal)
	}
	return nil
}

// Stats returns the statistic handler of the backtest.
func (t *Backtest) Stats() StatisticHandler {
	return t.statistic
}

// GetEvent process the event before return the data event back to user' scripts
func (t *Backtest) GetEvent() (err error, status string, data DataEvent){
	event := <-t.in
	return  t.eventActive(event)
}

// setup runs at the beginning of the backtest to perfom preparing operations.
func (t *Backtest) setup() error {
	// before first run, set portfolio cash
	t.portfolio.SetCash(t.portfolio.InitialCash())

	return nil
}

// teardown performs any cleaning operations at the end of the backtest.
func (t *Backtest) teardown() error {
	// no implementation yet
	return nil
}

// nextEvent gets the next event from the events queue.
func (t *Backtest) nextEvent() (e EventHandler, ok bool) {

	// if event queue empty return false
	if len(t.eventQueue) == 0 {
		return e, false
	}

	// return first element from the event queue
	e = t.eventQueue[0]
	t.eventQueue = t.eventQueue[1:]

	return e, true
}

// AddEvent
func (t *Backtest) AddEvent(e EventHandler) error {
	t.in <- e
	return nil
}

// eventActive directs the different events to their handler.
func (t *Backtest) eventActive(e EventHandler) (err error, status string, data DataEvent) {
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
		order, err := t.portfolio.OnSignal(event, t.data)
		if err != nil {
			break
		}
		t.AddEvent(order)

	case *Order:
		log.Infof("Get order event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())

		if event.status == OrdersBySymbol{
			orders, _ := t.exchange.OrdersBySymbol(event.symbol)
			var result Result
			result.SetTime(event.timestamp)
			result.SetData(orders)
			var rs ResultEvent
			rs = &result
			t.out <- ResultEvent(rs)
			break
		}

		t.exchange.AddOrder(event)
		if err != nil {
			break
		}
		//t.AddEvent(fill)

	case *Fill:
		log.Infof("Get fill event symbol (%s) timestamp (%s)", event.Symbol(), event.Time())
		t.exchange.OnFill(event)
		_, err := t.portfolio.OnFill(event)
		if err != nil {
			break
		}
		//t.AddEvent(transaction)
	}
	return
}
