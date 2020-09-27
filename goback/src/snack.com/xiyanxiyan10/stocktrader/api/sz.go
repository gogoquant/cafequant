package api

import "snack.com/xiyanxiyan10/stocktrader/constant"

// NewSZExchange create an exchange struct of futureExchange.com
func NewSZExchange(opt constant.Option) Exchange {
	exchange := NewSZSpotExchange(opt)
	_ = exchange.Init(opt)
	return exchange
}
