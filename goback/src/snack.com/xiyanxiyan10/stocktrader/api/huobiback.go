package api

import "snack.com/xiyanxiyan10/stocktrader/constant"

// NewHuoBiBackExchange create an exchange struct of futureExchange.com
func NewHuoBiBackExchange(opt constant.Option) Exchange {
	exchange := NewExchangeBackLink(opt)
	_ = exchange.Init(opt)
	return exchange
}
