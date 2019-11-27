package api

// NewFutureExchange create an exchange struct of futureExchange.com
func NewFmexExchange(opt Option) Exchange {
	exchange := NewFutureExchange(opt)
	exchange.Init()
	return exchange
}
