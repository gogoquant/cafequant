package goplugin

import (
	"snack.com/xiyanxiyan10/stocktrader/api"
)

// Gofunc ...
type Gofunc func(...interface{}) interface{}

// GoPlugin ...
type GoPlugin struct {
	goPluginGrid
	exchanges []api.Exchange
	funcs     map[string]Gofunc
}

// AddExchange ...
func (p *GoPlugin) AddExchange(e api.Exchange) {
	p.exchanges = append(p.exchanges, e)
}

// AddFunc ...
func (p *GoPlugin) AddFunc(name string, callback Gofunc) {
	p.funcs[name] = callback
}

// Call ...
func (p *GoPlugin) Call(name string, v ...interface{}) interface{} {
	callback, ok := p.funcs[name]
	if !ok {
		return false
	}
	return callback(v)
}

// Handler ...
type Handler interface {
	// string is the stragey name, interface is the param and return
	Call(string, ...interface{}) interface{}
}

// NewGoPlugin ...
func NewGoPlugin() *GoPlugin {
	var goplugin GoPlugin
	goplugin.funcs = make(map[string]Gofunc)

	//register go func here
	goplugin.AddFunc("hello", goplugin.Hello)

	return &goplugin
}
