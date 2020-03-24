package api

import "github.com/xiyanxiyan10/quantcore/constant"

// NewFutureExchange create an exchange struct of futureExchange.com
func NewHuoBiDmExchange(opt constant.Option) Exchange {
	exchange := NewFutureExchange(opt)
	_ = exchange.Init()
	return exchange
}
