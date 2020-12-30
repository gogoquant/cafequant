package api

import (
	"github.com/qiniu/py"
	"snack.com/xiyanxiyan10/stocktrader/constant"
)

// ExchangePython interface
type ExchangePython interface {
	Ready(args *py.Tuple) (ret *py.Base, err error)           //设置mode
	SetIO(args *py.Tuple) (ret *py.Base, err error)           //设置IO
	GetIO(args *py.Tuple) (ret *py.Base, err error)           //获取IO
	Subscribe(args *py.Tuple) (ret *py.Base, err error)       //订阅
	Log(args *py.Tuple) (ret *py.Base, err error)             //向管理台发送这个交易所的打印信息
	GetType(args *py.Tuple) (ret *py.Base, err error)         //获取交易所类型
	GetName(args *py.Tuple) (ret *py.Base, err error)         //获取交易所名称,自定义的
	SetLimit(args *py.Tuple) (ret *py.Base, err error)        //设置交易所的API访问频率
	Sleep(args *py.Tuple) (ret *py.Base, err error)           //延时
	GetAccount(args *py.Tuple) (ret *py.Base, err error)      //获取交易所的账户资金信息
	GetDepth(args *py.Tuple) (ret *py.Base, err error)        //返回买卖深度表
	Buy(args *py.Tuple) (ret *py.Base, err error)             //买
	Sell(args *py.Tuple) (ret *py.Base, err error)            //卖
	GetOrder(args *py.Tuple) (ret *py.Base, err error)        //返回订单信息
	GetOrders(args *py.Tuple) (ret *py.Base, err error)       //返回所有的未完成订单列表
	CancelOrder(args *py.Tuple) (ret *py.Base, err error)     //取消一笔订单
	GetTicker(args *py.Tuple) (ret *py.Base, err error)       //获取交易所的最新市场行情数据
	GetRecords(args *py.Tuple) (ret *py.Base, err error)      //返回交易所的最新K线数据列表
	SetContractType(args *py.Tuple) (ret *py.Base, err error) //设置合约周期
	GetContractType(args *py.Tuple) (ret *py.Base, err error) //获取合约周期
	SetDirection(args *py.Tuple) (ret *py.Base, err error)    //设置交易方向
	GetDirection(args *py.Tuple) (ret *py.Base, err error)    //获取交易方向
	SetMarginLevel(args *py.Tuple) (ret *py.Base, err error)  //杠杆设置
	GetMarginLevel(args *py.Tuple) (ret *py.Base, err error)  //获取杠杆
	SetStockType(args *py.Tuple) (ret *py.Base, err error)    //设置货币类型
	GetStockType(args *py.Tuple) (ret *py.Base, err error)    //获取货币类型
	GetPosition(args *py.Tuple) (ret *py.Base, err error)     //持仓量
	//AutoSleep(args *py.Tuple) (ret *py.Base, err error)     //自动休眠以满足设置的交易所的API访问频率

	// 回测
	GetBackAccount(args *py.Tuple) (ret *py.Base, err error)    //获取账号原货币量
	SetBackAccount(args *py.Tuple) (ret *py.Base, err error)    //设置账号原货币量
	SetBackCommission(args *py.Tuple) (ret *py.Base, err error) //设置回测手续费
	GetBackCommission(args *py.Tuple) (ret *py.Base, err error) //获取回测手续费
	SetBackTime(args *py.Tuple) (ret *py.Base, err error)       //设置回测周期
	GetBackTime(args *py.Tuple) (ret *py.Base, err error)       //设置回测周期
}

var (
	constructorPython = map[string]func(constant.Option) ExchangePython{}
)
