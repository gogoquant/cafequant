package api

import (
	"fmt"

	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

// ExchangeBroker define by every broker
type ExchangeBroker interface {
	getDepth(stockType string) (*constant.Depth, error)
	getOrder(symbol, id string) (*constant.Order, error)
	getAccount() (*constant.Account, error)
	getOrders(symbol string) ([]constant.Order, error)
	getTicker(symbol string) (*constant.Ticker, error)
	getPosition(stockType string) ([]constant.Position, error)
	buy(price, amount, msg string) (string, error)
	sell(price, amount, msg string) (string, error)
	cancelOrder(orderID string) (bool, error)
	start() error
	stop() error
}

// Exchange interface
type Exchange interface {
	Log(action, symbol string, price, amount float64, messages string)
	GetType() string
	GetName() string
	SetLimit(times int64) int64
	Sleep(intervals int64)
	AutoSleep()
	Buy(price, amount, msg string) (string, error)
	Sell(price, amount, msg string) (string, error)
	GetOrder(id string) (*constant.Order, error)
	GetOrders() ([]constant.Order, error)
	CancelOrder(orderID string) (bool, error)
	GetTicker() (*constant.Ticker, error)
	GetPosition() ([]constant.Position, error)
	GetAccount() (*constant.Account, error)
	GetDepth() (*constant.Depth, error)
	SetDirection(direction string)
	GetDirection() string
	SetMarginLevel(lever float64)
	GetMarginLevel() float64
	SetStockType(stockType string)
	GetStockType() string
	GetBackAccount() map[string]float64
	SetBackAccount(string, float64)
	SetBackCommission(float64, float64, float64, float64)
	GetBackCommission() []float64
	Start() error
	Stop() error
}

var (
	constructor = map[string]func(constant.Option) (Exchange, error){}
	// ExchangeMaker online exchange
	ExchangeMaker = map[string]func(constant.Option) (Exchange, error){
		constant.HuoBiDm: NewHuoBiDmExchange,
		constant.HuoBi:   NewHuoBiExchange,
		constant.SZ:      NewSZExchange,
	}
	// ExchangeBackerMaker backtest exchange
	ExchangeBackerMaker = map[string]func(constant.Option) (Exchange, error){
		constant.HuoBiDm: NewFutureBackExchange,
		constant.HuoBi:   NewSpotBackExchange,
		constant.SZ:      NewSpotBackExchange,
	}
)

// loadMaker ...
func loadMaker(exchangeType string) (func(constant.Option) (Exchange, error), error) {
	f, err := util.HotPlugin(config.String(fmt.Sprintf("/%s/%s.so",
		config.String(constant.GoPluginPath), exchangeType)), constant.GoHandler)
	if err != nil {
		return nil, err
	}
	makerplugin, ok := f.(func(constant.Option) (Exchange, error))
	if !ok {
		return nil, err
	}
	//register plugin into store
	ExchangeMaker[exchangeType] = makerplugin
	return makerplugin, nil
}

// GetExchange Maker
func getExchangeMaker(opt constant.Option) (maker func(constant.Option) (Exchange, error), ok bool) {
	exchangeType := opt.Type
	Back := opt.BackTest
	if !Back {
		maker, ok = ExchangeMaker[exchangeType]
		fmt.Printf("get online exchange %s\n", exchangeType)
		//if !ok {
		//	loadMaker(exchangeType)
		//}
		return
	}
	fmt.Printf("get back exchange %s\n", exchangeType)
	maker, ok = ExchangeBackerMaker[exchangeType]
	return
}

// GetExchange ...
func GetExchange(opt constant.Option) (Exchange, error) {
	maker, ok := getExchangeMaker(opt)
	if !ok {
		return nil, fmt.Errorf("get exchange maker fail")
	}
	exchange, err := maker(opt)
	if err != nil {
		return nil, err
	}
	return exchange, nil
}
