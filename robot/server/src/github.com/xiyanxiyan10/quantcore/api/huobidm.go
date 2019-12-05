package api

import "github.com/xiyanxiyan10/quantcore/constant"

// NewFutureExchange create an exchange struct of futureExchange.com
func NewHuobiDmExchange(opt constant.Option) Exchange {
	exchange := NewFutureExchange(opt)
	exchange.Init()
	return exchange
}
