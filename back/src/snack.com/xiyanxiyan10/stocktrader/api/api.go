package api

import "snack.com/xiyanxiyan10/stocktrader/constant"

// Exchange interface
type Exchange interface {
	Log(...interface{})                                        //向管理台发送这个交易所的打印信息
	GetType() string                                           //获取交易所类型
	GetName() string                                           //获取交易所名称,自定义的
	SetLimit(times interface{}) float64                        //设置交易所的API访问频率,和 E.AutoSleep() 配合使用
	AutoSleep()                                                //自动休眠以满足设置的交易所的API访问频率
	GetMinAmount(stock string) float64                         //获取交易所的最小交易数量
	GetAccount() interface{}                                   //获取交易所的账户资金信息
	GetDepth(size int) interface{}                             //返回买卖深度表
	Buy(price, amount string, msg ...interface{}) interface{}  //买
	Sell(price, amount string, msg ...interface{}) interface{} //卖
	GetOrder(id string) interface{}                            //返回订单信息
	GetOrders() interface{}                                    //返回所有的未完成订单列表
	GetTrades(params ...interface{}) interface{}               //返回最近的已完成订单列表
	CancelOrder(orderID string) bool                           //取消一笔订单
	GetTicker() interface{}                                    //获取交易所的最新市场行情数据
	GetRecords(params ...interface{}) interface{}              //返回交易所的最新K线数据列表
	SetContractType(contractType string)                       //设置合约周期
	GetContractType() string                                   //获取合约周期
	SetDirection(direction string)                             //设置交易方向
	GetDirection() string                                      //获取交易方向
	SetMarginLevel(lever int)                                  //杠杆设置
	GetMarginLevel() int                                       //获取杠杆
	SetStockType(stockType string)                             //设置货币类型
	GetStockType() string                                      //获取货币类型
	GetPosition() interface{}                                  //持仓量
}

var (
	constructor = map[string]func(constant.Option) Exchange{}
)