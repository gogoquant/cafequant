package api

import (
	dbtypes "snack.com/xiyanxiyan10/stockdb/types"
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// Exchange interface
type Exchange interface {
	Ready() error                                                  //设置mode
	IsBack() bool                                                  //is back exchange
	SetIO(mode string)                                             //设置IO
	GetIO() string                                                 //获取IO
	GetSubscribe() map[string][]string                             //获取订阅
	SetSubscribe(string, string)                                   //订阅
	Log(...interface{})                                            //向管理台发送这个交易所的打印信息
	GetType() string                                               //获取交易所类型
	GetName() string                                               //获取交易所名称,自定义的
	SetLimit(times interface{}) float64                            //设置交易所的API访问频率,和 E.AutoSleep() 配合使用
	Sleep(intervals ...interface{})                                //延时
	AutoSleep()                                                    //自动休眠以满足设置的交易所的API访问频率
	GetAccount() (*constant.Account, error)                        //获取交易所的账户资金信息
	GetDepth() (*constant.Depth, error)                            //返回买卖深度表
	Buy(price, amount string, msg ...interface{}) (string, error)  //买
	Sell(price, amount string, msg ...interface{}) (string, error) //卖
	GetOrder(id string) (*constant.Order, error)                   //返回订单信息
	GetOrders() ([]constant.Order, error)                          //返回所有的未完成订单列表
	CancelOrder(orderID string) (bool, error)                      //取消一笔订单
	GetTicker() (*constant.Ticker, error)                          //获取交易所的最新市场行情数据
	GetRecords(periodStr, maStr string) ([]constant.Record, error) //返回交易所的最新K线数据列表
	SetContractType(contractType string)                           //设置合约周期
	GetContractType() string                                       //获取合约周期
	SetDirection(direction string)                                 //设置交易方向
	GetDirection() string                                          //获取交易方向
	SetMarginLevel(lever float64)                                  //杠杆设置
	GetMarginLevel() float64                                       //获取杠杆
	SetStockType(stockType string)                                 //设置货币类型
	GetStockType() string                                          //获取货币类型
	GetPosition() ([]constant.Position, error)                     //持仓量

	// backtest
	GetBackAccount() map[string]float64
	SetBackAccount(string, float64)                                                                   //账号原货币量
	SetBackCommission(float64, float64, float64, float64, bool)                                       //设置回测手续费
	GetBackCommission() (float64, float64, float64, float64, bool)                                    //获取回测手续费
	SetBackTime(start, end int64, period string)                                                      //设置回测周期
	GetBackTime() (int64, int64, string)                                                              //设置回测周期
	BackPutOHLC(time int64, open, high, low, closed, volume float64, ext string, period string) error //推送数据到数据仓库
	BackGetSymbols() ([]string, error)                                                                //获取货币种类
	BackGetMarkets() ([]dbtypes.Stats, error)                                                         //获取交易所种类
	BackGetStats() error                                                                              //获取数据中心数据

	BackGetPeriodRange() ([2]int64, error)                                //获取周期范围
	BackGetTimeRange() ([2]int64, error)                                  //获取事件范围
	BackGetOHLCs(begin, end int64, period string) ([]dbtypes.OHLC, error) //获取OHLC
	BackGetDepth(begin, end int64, period string) (dbtypes.Depth, error)  //获取Depth
}

var (
	constructor = map[string]func(constant.Option) (Exchange, error){}
)
