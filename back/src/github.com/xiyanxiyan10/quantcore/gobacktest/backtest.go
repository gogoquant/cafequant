package gobacktest

import "time"

// DP sets the the precision of rounded floating numbers
// used after calculations to format
const DP = 4 // DP

// Reset provides a resting interface.
type ResetHandler interface {
	Reset() error
}

// BackTest is the main struct which holds all elements.
type BackTest struct {
	symbols    []string
	data       DataHandler
	portfolio  PortfolioHandler
	exchange   ExecutionHandler
	statistic  StatisticHandler
	eventQueue []EventHandler
}

// New creates a default backTest with sensible defaults ready for use.
func New() *BackTest {
	return &BackTest{
		portfolio: &Portfolio{
			initialCash: 100000,
			sizeManager: &Size{DefaultSize: 100, DefaultValue: 1000},
			riskManager: &Risk{},
		},
		exchange: &Exchange{
			Symbol:      "TEST",
			Commission:  &FixedCommission{Commission: 0},
		},
		statistic: &Statistic{},
	}
}

// SetSymbols sets the symbols to include into the backTest.
func (t *BackTest) SetSymbols(symbols []string) {
	t.symbols = symbols
}

// SetData sets the data provider to be used within the backTest.
func (t *BackTest) SetData(data DataHandler) {
	t.data = data
}

// SetPortfolio sets the portfolio provider to be used within the backTest.
func (t *BackTest) SetPortfolio(portfolio PortfolioHandler) {
	t.portfolio = portfolio
}

// SetExchange sets the execution provider to be used within the backTest.
func (t *BackTest) SetExchange(exchange ExecutionHandler) {
	t.exchange = exchange
}

// SetStatistic sets the statistic provider to be used within the backTest.
func (t *BackTest) SetStatistic(statistic StatisticHandler) {
	t.statistic = statistic
}

// Reset the backTest into a clean state with loaded data.
func (t *BackTest) Reset() error {
	t.eventQueue = nil
	_ = t.data.Reset()
	_ = t.portfolio.Reset()
	_ = t.statistic.Reset()
	return nil
}

// Stats returns the statistic handler of the backTest.
func (t *BackTest) Stats() StatisticHandler {
	return t.statistic
}

// Run starts the backTest.
func (t *BackTest) Run() error {
	// setup before the backTest runs
	err := t.setup()
	if err != nil {
		return err
	}

	// poll event queue
	for event, ok := t.nextEvent(); true; event, ok = t.nextEvent() {
		// no event in the queue
		if !ok {
			// poll data stream
			data, ok := t.data.Next()
			// no more data, exit event loop
			if !ok {
				break
			}
			// found data event, add to event stream
			t.eventQueue = append(t.eventQueue, data)
			// start new event cycle
			continue
		}

		// processing event
		err := t.eventLoop(event)
		if err != nil {
			return err
		}
		// event in queue found, add to event history
		t.statistic.TrackEvent(event)
	}

	// teardown at the end of the backTest
	err = t.teardown()
	if err != nil {
		return err
	}

	return nil
}

// Run2Time starts the backTest. //
func (t *BackTest) Run2Next(next time.Time) (DataEvent, error) {
	for {
		// poll data from stream
		data, ok := t.data.Next()
		// no more data, exit event loop
		if !ok {
			break
		}
		// found data event, add to event stream
		t.eventQueue = append(t.eventQueue, data)
		// poll event queue
		for event, ok := t.nextEvent(); true; event, ok = t.nextEvent() {
			// no event in the queue
			if !ok {
				// start new event cycle
				break
			}

			// processing event
			err := t.eventLoop(event)
			if err != nil {
				return nil, err
			}
			// event in queue found, add to event history
			t.statistic.TrackEvent(event)
		}
		//return data if the time after next policy trigger
		if data.Time().After(next) {
			return data, nil
		}
	}
	return nil, nil
}

// setup runs at the beginning of the backTest to perform preparing operations.
func (t *BackTest) setup() error {
	// before first run, set portfolio cash
	t.portfolio.SetCash(t.portfolio.InitialCash())

	return nil
}

// teardown performs any cleaning operations at the end of the backTest.
func (t *BackTest) teardown() error {
	// no implementation yet
	return nil
}

// nextEvent gets the next event from the events queue.
func (t *BackTest) nextEvent() (e EventHandler, ok bool) {
	// if event queue empty return false
	if len(t.eventQueue) == 0 {
		return e, false
	}

	// return first element from the event queue
	e = t.eventQueue[0]
	t.eventQueue = t.eventQueue[1:]

	return e, true
}

// eventLoop directs the different events to their handler.
func (t *BackTest) eventLoop(e EventHandler) error {
	// check the order

	// type check for event type
	switch event := e.(type) {
	case DataEvent:
		// update portfolio to the last known price data
		t.portfolio.Update(event)
		// update statistics
		t.statistic.Update(event, t.portfolio)
		// check if any orders are filled before preceding
		t.exchange.OnData(event)
		orders, ok := t.portfolio.OnData(t.data)
		// add orders into queue which is married
		if ok {
			for _, order := range orders {
				t.eventQueue = append(t.eventQueue, order)
			}
		}
		// run strategy with this data event
		/*
			signals, err := t.strategy.OnData(event)
			if err != nil {
				break
			}
		*/
		/*
			for _, signal := range signals {
				t.eventQueue = append(t.eventQueue, signal)
			}
		*/

	case *Signal:
		order, err := t.portfolio.OnSignal(event, t.data)
		if err != nil {
			break
		}
		t.eventQueue = append(t.eventQueue, order)

	case *Order:
		fill, err := t.exchange.OnOrder(event, t.data)
		if err != nil {
			break
		}
		t.eventQueue = append(t.eventQueue, fill)

	case *Fill:
		transaction, err := t.portfolio.OnFill(event, t.data)
		if err != nil {
			break
		}
		t.statistic.TrackTransaction(transaction)
	}

	// Todo check the position
	// check the cost

	return nil
}
