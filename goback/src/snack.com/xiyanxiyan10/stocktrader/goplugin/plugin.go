package goplugin

import (
	"errors"
	"plugin"
	"snack.com/xiyanxiyan10/stocktrader/api"
)

// Gofunc ...
type Gofunc func(...interface{}) interface{}

// GoStragey  ...
type GoStragey struct {
	exchanges []api.Exchange
}

// AddExchange ...
func (p *GoStragey) AddExchange(e ...api.Exchange) {
	p.exchanges = append(p.exchanges, e...)
}

// GoStrageyHandler ...
type GoStrageyHandler interface {
	AddExchange(...api.Exchange)
	Init(...interface{}) interface{}
	Call(...interface{}) interface{}
	Exit(...interface{}) interface{}
}

// GoPlugin ...
type GoPlugin struct {
	exchanges []api.Exchange
	strageys  map[string]GoStrageyHandler
}

// AddExchange ...
func (p *GoPlugin) AddExchange(e ...api.Exchange) {
	p.exchanges = append(p.exchanges, e...)
}

// AddStragey ...
func (p *GoPlugin) AddStragey(name string, v GoStrageyHandler) {
	p.strageys[name] = v
	v.AddExchange(p.exchanges...)
}

// LoadStragey ...
func (p *GoPlugin) LoadStragey(name string) error {
	handler, err := plugin.Open("/tmp/" + name + ".so")
	if err != nil {
		return err
	}
	s, err := handler.Lookup("NewHandler")
	if err != nil {
		return err
	}
	v, ok := s.(GoStrageyHandler)
	if !ok {
		return errors.New("load stragey handler fail")
	}
	p.AddStragey(name, v)
	return nil
}

// Init ...
func (p *GoPlugin) Init(name string, v ...interface{}) interface{} {
	handler, ok := p.strageys[name]
	if !ok {
		return false
	}
	return handler.Init(v)
}

// Call ...
func (p *GoPlugin) Call(name string, v ...interface{}) interface{} {
	handler, ok := p.strageys[name]
	if !ok {
		return false
	}
	return handler.Call(v)
}

// Exit ...
func (p *GoPlugin) Exit(name string, v ...interface{}) interface{} {
	handler, ok := p.strageys[name]
	if !ok {
		return false
	}
	return handler.Exit(v)
}

// Handler ...
type Handler interface {
	AddStragey(name string, v GoStrageyHandler)
	LoadStragey(name string) error
	Init(name string, v ...interface{}) interface{}
	Call(name string, v ...interface{}) interface{}
	Exit(name string, v ...interface{}) interface{}
}

// NewGoPlugin ...
func NewGoPlugin() *GoPlugin {
	var goplugin GoPlugin
	return &goplugin
}
