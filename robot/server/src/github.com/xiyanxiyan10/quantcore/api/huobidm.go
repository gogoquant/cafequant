package api

// NewFutureExchange create an exchange struct of futureExchange.com
func NewHuobiDmExchange(opt Option) Exchange {
	exchange := NewFutureExchange(opt)
	exchange.Init()
	return exchange
}
