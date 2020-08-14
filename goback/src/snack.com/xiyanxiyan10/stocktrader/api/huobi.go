package api

import "snack.com/xiyanxiyan10/stocktrader/constant"

// NewHuoBiExchange create an exchange struct of futureExchange.com
func NewHuoBiExchange(opt constant.Option) Exchange {
	exchange := NewSpotExchange(opt)
	_ = exchange.Init(opt)
	return exchange
}
