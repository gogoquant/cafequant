package api

// Option is an exchange option
type Option struct {
	TraderID  int64
	Type      string
	Name      string
	AccessKey string
	SecretKey string
}

// Exchange interface
type Exchange interface {
	Log(...interface{})                              //向管理台发送这个交易所的打印信息
	GetType() string                                 //获取交易所类型,是火币还是OKEY等。。。
	GetName() string                                 //获取交易所名称,自定义的
	SetLimit(times interface{}) float64              //设置交易所的API访问频率,和 E.AutoSleep() 配合使用
	AutoSleep()                                      //自动休眠以满足设置的交易所的API访问频率
	GetMinAmount(stock string) float64               //获取交易所的最小交易数量
	GetAccount() interface{}                         //获取交易所的账户资金信息
	GetDepth(size int, stockType string) interface{} //返回买卖深度表
	Buy(price, amount string, msg ...interface{}) interface{}
	Sell(price, amount string, msg ...interface{}) interface{}
	GetOrder(id string) interface{}                             //返回订单信息
	GetOrders() interface{}                                     //返回所有的未完成订单列表
	GetTrades() interface{}                                     //返回最近的已完成订单列表
	CancelOrder(orderID string) bool                            //取消一笔订单
	GetTicker(sizes ...interface{}) interface{}                 //获取交易所的最新市场行情数据
	GetRecords(period string, sizes ...interface{}) interface{} //返回交易所的最新K线数据列表
	SetContractType(contractType string)
	GetContractType() string
	SetDirection(direction string)
	GetDirection() string
	SetStockType(stockType string)
	GetStockType() string
}

var (
	constructor = map[string]func(Option) Exchange{}
)
