package goplugin

import (
	"snack.com/xiyanxiyan10/stocktrader/api"
)

// GoPlugin ...
type GoPlugin struct {
	exchanges []api.Exchange
}

// AddExchange ...
func (p *GoPlugin) AddExchange(e api.Exchange) {
	p.exchanges = append(p.exchanges, e)
}

// Handler ...
type Handler interface {
}
