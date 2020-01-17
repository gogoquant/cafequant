package api

import "github.com/xiyanxiyan10/quantcore/constant"

// NewHuoBiExchange create an exchange struct of futureExchange.com
func NewHuoBiExchange(opt constant.Option) Exchange {
	exchange := NewSpotExchange(opt)
	exchange.Init()
	return exchange
}
