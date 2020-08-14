package api

import "snack.com/xiyanxiyan10/stocktrader/constant"

// NewHuoBiDmExchange create an exchange struct of futureExchange.com
func NewHuoBiDmExchange(opt constant.Option) Exchange {
	exchange := NewFutureExchange(opt)
	_ = exchange.Init(opt)
	return exchange
}
