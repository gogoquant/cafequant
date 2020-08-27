package goplugin

import (
	"plugin"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/notice"
)

// Gofunc ...
type Gofunc func(...interface{}) interface{}

// StrageyNew ...
type StrageyNew func(...interface{}) (GoStrageyHandler, error)

// GoStragey  ...
type GoStragey struct {
	Exchanges []api.Exchange
	Mail      notice.MailHandler // 邮件发送
	Draw      draw.DrawHandler   // 图标绘制
	Logger    *model.Logger      // 利用这个对象保存日志
}

// AddExchange ...
func (p *GoStragey) AddExchange(e ...api.Exchange) {
	p.Exchanges = append(p.Exchanges, e...)
}

// AddDraw ...
func (p *GoStragey) AddDraw(draw draw.DrawHandler) {
	p.Draw = draw
}

// AddMail ...
func (p *GoStragey) AddMail(mail notice.MailHandler) {
	p.Mail = mail
}

// AddLogger ...
func (p *GoStragey) AddLogger(logger *model.Logger) {
	p.Logger = logger
}

// GoStrageyHandler ...
type GoStrageyHandler interface {
	AddExchange(...api.Exchange)
	AddDraw(draw draw.DrawHandler)
	AddMail(mail notice.MailHandler)
	AddLogger(logger *model.Logger)

	Init(...interface{}) interface{}
	Call(...interface{}) interface{}
	Exit(...interface{}) interface{}
}

// GoPlugin ...
type GoPlugin struct {
	Name      string
	Logger    *model.Logger // 利用这个对象保存日志
	Exchanges []api.Exchange
	Mail      notice.MailHandler // 邮件发送
	Draw      draw.DrawHandler   // 图标绘制
	strageys  map[string]GoStrageyHandler
}

// SetStragey ...
func (p *GoPlugin) SetStragey(name string) {
	p.Name = name
}

// GetStragey ...
func (p *GoPlugin) GetStragey() string {
	return p.Name
}

// AddExchange ...
func (p *GoPlugin) AddExchange(e ...api.Exchange) {
	p.Exchanges = append(p.Exchanges, e...)
}

// AddDraw ...
func (p *GoPlugin) AddDraw(draw draw.DrawHandler) {
	p.Draw = draw
}

// AddMail ...
func (p *GoPlugin) AddMail(mail notice.MailHandler) {
	p.Mail = mail
}

// AddStragey ...
func (p *GoPlugin) AddStragey(v GoStrageyHandler) {
	p.strageys[p.GetStragey()] = v
	v.AddExchange(p.Exchanges...)
	v.AddMail(p.Mail)
	v.AddDraw(p.Draw)
	v.AddLogger(p.Logger)
}

// LoadStragey ...
func (p *GoPlugin) LoadStragey() (res interface{}) {
	res = nil
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		}
	}()
	name := p.GetStragey()
	handler, err := plugin.Open(config.String(constant.GoPluginPath+"/") + name + ".so")
	if err != nil {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		return nil
	}
	s, err := handler.Lookup(constant.GoHandler)
	if err != nil {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		return nil
	}
	newHandler, ok := s.(StrageyNew)
	if !ok {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail convert handler")
		return nil
	}
	strageyHandler, err := newHandler()
	if err != nil {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		return nil
	}
	p.AddStragey(strageyHandler)
	return "success"
}

// Init ...
func (p *GoPlugin) Init(v ...interface{}) (res interface{}) {
	res = nil
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail:%v", err)
		}
	}()
	handler, ok := p.strageys[p.GetStragey()]
	if !ok {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail:get name")
		return nil
	}
	return handler.Init(v)
}

// Call ...
func (p *GoPlugin) Call(v ...interface{}) interface{} {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Call fail:%v", err)
		}
	}()
	handler, ok := p.strageys[p.GetStragey()]
	if !ok {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Call fail:get name")
		return nil
	}
	return handler.Call(v)
}

// Exit ...
func (p *GoPlugin) Exit(v ...interface{}) interface{} {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Exit fail:%v", err)
		}
	}()
	handler, ok := p.strageys[p.GetStragey()]
	if !ok {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Exit fail:get name")
		return nil
	}
	return handler.Exit(v)
}

// Handler ...
type Handler interface {
	AddStragey(v GoStrageyHandler)
	LoadStragey() interface{}
	SetStragey(string)
	GetStragey() string
	Init(v ...interface{}) interface{}
	Call(v ...interface{}) interface{}
	Exit(v ...interface{}) interface{}
}

// NewGoPlugin ...
func NewGoPlugin() *GoPlugin {
	var goplugin GoPlugin
	goplugin.SetStragey("echo")
	grid, _ := NewEchoHandler()
	goplugin.AddStragey(grid)
	return &goplugin
}
