package api

import "snack.com/xiyanxiyan10/stocktrader/constant"

// NewFutureExchange create an exchange struct of futureExchange.com
func NewFmexExchange(opt constant.Option) Exchange {
	exchange := NewFutureExchange(opt)
	_ = exchange.Init()
	return exchange
}
