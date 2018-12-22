package gobacktest

import (
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"gopkg.in/logger.v1"
	"time"
)

var errHalt = fmt.Errorf("HALT")

// Back 回测引擎对外接口
type Back interface {
	// 获取回测引擎名称
	Name() string

	// 设置回测引擎的名称
	SetName(name string)

	// 开始回测
	Start() (err error)

	// 运行状态
	Status() int64

	// 停止回测
	Stop() (err error)

	// 设置撮合脚本
	SetScripts(scripts string)

	// 获取仓位
	Holds() (map[string]Position, error)

	// 获取某类订单
	OrdersBySymbol(symbol string) ([]OrderEvent, error)

	// 关闭订单
	CancelOrder(order Order) error

	// 新开订单
	AddOrder(order Order) error

	// 特殊命令，关闭引擎等
	Cmd(cmd string) error
}

// ScriptsApi 撮合脚本使用的api
type ScriptsApi interface {
	// 提交订单
	CommitOrder(id int) error

	// 驱动引擎并尝试获取下一份数据
	NextEvent() (err error, status string, data DataEvent)

	// 获取某类订单
	EOrdersBySymbol(symbol string) ([]OrderEvent, error)
}

// Reseter provides a resting interface.
type Reseter interface {
	Reset() error
}

// Backtest 回测管理
type Backtest struct {
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
	scripts    string

	Ctx *otto.Otto
}

// SignalAdd Add signal event into event queue
func (back *Backtest) AddSignal(signals ...SignalEvent) error {
	for _, signal := range signals {
		back.AddEvent(signal)
	}
	return nil
}

// Stats returns the statistic handler of the backtest.
func (back *Backtest) Stats() StatisticHandler {
	return back.statistic
}

// OrdersBySymbol get order by symbol
func (back *Backtest) OrdersBySymbol(symbol string) ([]OrderEvent, error) {
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
func (back *Backtest) EOrdersBySymbol(symbol string) ([]OrderEvent, error) {
	orders, _ := back.exchange.OrdersBySymbol(symbol)
	return orders, nil
}

// AddOrder
func (back *Backtest) AddOrder(order Order) error {
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
func (back *Backtest) Holds() (map[string]Position, error) {
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

// SetScripts ...
func (back *Backtest) SetScripts(scripts string) {
	back.scripts = scripts
}

// SetCash ...
func (back *Backtest) SetCash(cash float64) {
	back.Portfolio().SetInitialCash(cash)
	back.Portfolio().SetCash(cash)
}

// getResult ...
func (back *Backtest) getResult() (ResultEvent, error) {
	select {
	case data, _ := <-back.out:
		return data, nil
	case <-time.After(1 * time.Second):
		return nil, errors.New("Timeout")
	}
}

// Start ...
func (back *Backtest) Start() (err error) {
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
func (back *Backtest) SetId(id string) {
	back.Id = id
}

// Status ...
func (back *Backtest) Status() int64 {
	return back.status
}

// Stop ...
func (back *Backtest) Stop() (err error) {
	back.Ctx.Interrupt <- func() { panic(errHalt) }
	return
}

// NewBackTest
func NewBackTest(m map[string]string) Back {
	bt := Backtest{
		in:     make(chan EventHandler, 50),
		out:    make(chan ResultEvent, 50),
		status: 0,
		config: m,
	}
	err := bt.initialize()
	if err != nil {
		return nil
	} else {
		return &bt
	}
}

// Name ...
func (back *Backtest) Name() string {
	return back.name
}

// SetName ...
func (back *Backtest) SetName(name string) {
	back.name = name
}

// CommitOrder ...
func (back *Backtest) CommitOrder(id int) error {
	fill, err := back.exchange.CommitOrder(id)
	if err == nil && fill != nil {
		back.AddEvent(fill)
	}
	return err
}

// CancelOrder ...
func (back *Backtest) CancelOrder(order Order) error {
	order.SetStatus(OrderCancel)
	back.AddEvent(&order)
	return nil
}

// initialize ...
func (back *Backtest) initialize() (err error) {
	back.Ctx = otto.New()
	back.Ctx.Interrupt = make(chan func(), 1)
	back.Ctx.Set("Exchange", ScriptsApi(back))
	return
}

// holds
func (back *Backtest) holds() (res Result) {
	m := back.portfolio.Holds()
	var p Position
	p.fqty = back.Portfolio().Cash()
	m["cash"] = p
	res.SetData(m)
	return
}

// Cmd ...
func (back *Backtest) Cmd(cmd string) error {
	var event Cmd
	event.SetCmd(cmd)
	back.AddEvent(&event)
	return nil
}

// New creates a default backtest with sensible defaults ready for use.
func New() *Backtest {
	return &Backtest{
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
}

// SetSymbols sets the symbols to include into the backtest.
func (back *Backtest) SetSymbols(symbols []string) {
	back.symbols = symbols
}

// SetData sets the data provider to be used within the backtest.
func (back *Backtest) SetData(data DataHandler) {
	back.data = data
}

// SetPortfolio sets the portfolio provider to be used within the backtest.
func (back *Backtest) SetPortfolio(portfolio PortfolioHandler) {
	back.portfolio = portfolio
}

// SetExchange sets the execution provider to be used within the backtest.
func (back *Backtest) SetExchange(exchange ExchangeHandler) {
	back.exchange = exchange
}

// SetStatistic sets the statistic provider to be used within the backtest.
func (back *Backtest) SetStatistic(statistic StatisticHandler) {
	back.statistic = statistic
}

// Portfolio sets the Portfolio provider to be used within the backtest.
func (back *Backtest) Portfolio() PortfolioHandler {
	return back.portfolio
}

func (back *Backtest) Exchange() ExchangeHandler {
	return back.exchange
}

// Reset ...
func (back *Backtest) Reset() error {
	back.eventQueue = nil
	back.data.Reset()
	back.portfolio.Reset()
	back.statistic.Reset()
	return nil
}

// NextEvent process the event before return the data event back to user' scripts
func (back *Backtest) NextEvent() (err error, status string, data DataEvent) {
	event := <-back.in
	return back.activeEvent(event)
}

// setup runs at the beginning of the backtest to perfom preparing operations.
func (back *Backtest) setup() error {
	// before first run, set portfolio cash
	back.portfolio.SetCash(back.portfolio.InitialCash())
	return nil
}

// teardown performs any cleaning operations at the end of the backtest.
func (back *Backtest) teardown() error {
	// no implementation yet
	return nil
}

// nextEvent gets the next event from the events queue.
func (back *Backtest) nextEvent() (e EventHandler, ok bool) {

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
func (back *Backtest) AddEvent(e EventHandler) error {
	back.in <- e
	return nil
}

// activeEvent directs the different events to their handler.
func (back *Backtest) activeEvent(e EventHandler) (err error, status string, data DataEvent) {
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

		// 获取在挂的单
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

		// 关闭订单
		if event.status == OrderCancel {
			back.exchange.CancelOrder(*event)
			break
		}

		// 获取仓位
		if event.status == OrderFillByAll {
			fills := back.holds()
			back.out <- &fills
			break
		}

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
		log.Infof("Unknow order type")
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
