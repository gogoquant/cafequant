package api

import (
	"fmt"
	dbtypes "snack.com/xiyanxiyan10/stockdb/types"
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
	getRecords(stockType string) ([]constant.Record, error)
	getPosition(stockType string) ([]constant.Position, error)
	buy(price, amount, msg string) (string, error)
	sell(price, amount, msg string) (string, error)
	cancelOrder(orderID string) (bool, error)
	start() error
	stop() error
}

// Exchange interface
type Exchange interface {
	// 初始化完毕start run
	Start() error
	// 初始化完毕stop run
	Stop() error
	// 设置IO
	SetIO(mode string)
	// 获取IO
	GetIO() string
	// Set Period
	SetPeriod(string)
	// Get Period
	GetPeriod() string
	// Set Period size
	SetPeriodSize(int)
	// Get Period Size
	GetPeriodSize() int
	// 获取订阅
	GetSubscribe() map[string][]string
	// 订阅
	SetSubscribe(string, string)
	// 向管理台发送这个交易所的打印信息
	Log(action, symbol string, price, amount float64, messages string)
	// 获取交易所类型
	GetType() string
	// 获取交易所名称,自定义的
	GetName() string
	// 设置交易所的API访问频率,和 E.AutoSleep() 配合使用
	SetLimit(times int64) int64
	// 延时
	Sleep(intervals int64)
	// 自动休眠以满足设置的交易所的API访问频率
	AutoSleep()
	// 买
	Buy(price, amount, msg string) (string, error)
	// 卖
	Sell(price, amount, msg string) (string, error)
	// 返回订单信息
	GetOrder(id string) (*constant.Order, error)
	// 返回所有的未完成订单列表
	GetOrders() ([]constant.Order, error)
	// 取消一笔订单
	CancelOrder(orderID string) (bool, error)
	// 获取交易所的最新市场行情数据
	GetTicker() (*constant.Ticker, error)
	// 返回交易所的最新K线数据列表, 部分平台可以直接获取计算好的均线
	GetRecords() ([]constant.Record, error)
	// 持仓量
	GetPosition() ([]constant.Position, error)
	// 获取交易所的账户资金信息
	GetAccount() (*constant.Account, error)
	// 返回买卖深度表
	GetDepth() (*constant.Depth, error)
	// 设置交易方向
	SetDirection(direction string)
	// 获取交易方向
	GetDirection() string
	// 杠杆设置
	SetMarginLevel(lever float64)
	// 获取杠杆
	GetMarginLevel() float64
	// 设置货币类型
	SetStockType(stockType string)
	// 获取货币类型
	GetStockType() string
	// 获取回测账号
	GetBackAccount() map[string]float64
	// 账号原货币量
	SetBackAccount(string, float64)
	// 设置回测手续费
	SetBackCommission(float64, float64, float64, float64)
	// 获取回测手续费
	GetBackCommission() []float64
	// 设置回测周期
	SetBackTime(start, end int64, period string)
	//设置回测周期
	GetBackTime() constant.BackTime
	//推送数据到数据仓库
	BackPutOHLC(time int64, open, high, low, closed, volume float64, ext, period string) error
	//推送数据 [] 到数据仓库
	BackPutOHLCs(datums []dbtypes.OHLC, period string) error
	//获取货币种类
	BackGetSymbols() ([]string, error)
	//获取交易所种类
	BackGetMarkets() ([]string, error)
	//获取数据中心数据
	BackGetStats() ([]dbtypes.Stats, error)
	//获取周期范围
	BackGetPeriodRange() ([2]int64, error)
	//获取时间范围
	BackGetTimeRange() ([2]int64, error)
	//获取OHLC
	BackGetOHLCs(begin, end int64, period string) ([]dbtypes.OHLC, error)
	//获取Depth
	BackGetDepth(begin, end int64, period string) (dbtypes.Depth, error)
}

var (
	constructor = map[string]func(constant.Option) (Exchange, error){}
	// ExchangeMaker online exchange
	ExchangeMaker = map[string]func(constant.Option) (Exchange, error){ //保存所有交易所的构造函数
		constant.HuoBiDm: NewHuoBiDmExchange,
		constant.HuoBi:   NewHuoBiExchange,
		constant.SZ:      NewSZExchange,
	}
	// ExchangeBackerMaker backtest exchange
	ExchangeBackerMaker = map[string]func(constant.Option) (Exchange, error){
		//保存所有交易所的构造函数
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
		_, ok = ExchangeMaker[exchangeType]
		if !ok {
			loadMaker(exchangeType)
		}
	}
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
