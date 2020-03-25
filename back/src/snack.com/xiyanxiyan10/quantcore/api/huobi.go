package api

import "snack.com/xiyanxiyan10/quantcore/constant"

// NewHuoBiExchange create an exchange struct of futureExchange.com
func NewHuoBiExchange(opt constant.Option) Exchange {
	exchange := NewSpotExchange(opt)
	_ = exchange.Init()
	return exchange
}
