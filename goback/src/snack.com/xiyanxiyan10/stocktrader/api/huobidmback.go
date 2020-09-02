package api

import "snack.com/xiyanxiyan10/stocktrader/constant"

// NewHuoBiDmBackExchange create an exchange struct of futureExchange.com
func NewHuoBiDmBackExchange(opt constant.Option) Exchange {
	exchange := NewExchangeFutureBackLink(opt)
	_ = exchange.Init(opt)
	return exchange
}
