package api

import "snack.com/xiyanxiyan10/stocktrader/constant"

// NewHuoBiDmExchange create an exchange struct of futureExchange.com
func NewHuoBiDmExchange(opt constant.Option) Exchange {
	exchange := NewFutureExchange(opt)
	//exchange.SetTradeTypeMap(4, "closebuy")
	//exchange.SetTradeTypeMap(3, "closesell")
	_ = exchange.Init()
	return exchange
}
