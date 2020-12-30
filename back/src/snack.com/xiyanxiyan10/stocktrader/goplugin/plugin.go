package goplugin

import (
	"fmt"
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
	Ding      notice.DingHandler // dingding
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

// AddDing ...
func (p *GoStragey) AddDing(ding notice.DingHandler) {
	p.Ding = ding
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
	AddDing(ding notice.DingHandler)
	AddLogger(logger *model.Logger)

	Init(...interface{}) interface{}
	Run(...interface{}) interface{}
	Exit(...interface{}) interface{}
}

// GoPlugin ...
type GoPlugin struct {
	name      string
	Logger    *model.Logger               // 利用这个对象保存日志
	Exchanges []api.Exchange              // 交易所列表
	Mail      notice.MailHandler          // 邮件发送
	Draw      draw.DrawHandler            // 图标绘制
	Ding      notice.DingHandler          // dingding
	strageys  map[string]GoStrageyHandler // 策略集合
}

// SetStragey ...
func (p *GoPlugin) SetStragey(name string) {
	p.name = name
}

// GetStragey ...
func (p *GoPlugin) GetStragey() string {
	return p.name
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

// AddDing ...
func (p *GoPlugin) AddDing(ding notice.DingHandler) {
	p.Ding = ding
}

// AddStragey ...
func (p *GoPlugin) AddStragey(v GoStrageyHandler) {
	p.strageys[p.GetStragey()] = v
	v.AddExchange(p.Exchanges...)
	v.AddMail(p.Mail)
	v.AddDraw(p.Draw)
	v.AddLogger(p.Logger)
	v.AddDing(p.Ding)
}

// LoadStragey ...
func (p *GoPlugin) LoadStragey() error {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		}
	}()
	name := p.GetStragey()
	handler, err := plugin.Open(config.String(constant.GoPluginPath+"/") + name + ".so")
	if err != nil {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		return err
	}
	s, err := handler.Lookup(constant.GoHandler)
	if err != nil {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		return err
	}
	newHandler, ok := s.(StrageyNew)
	if !ok {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail convert handler")
		return fmt.Errorf("LoadStragey() fail convert handler")
	}
	strageyHandler, err := newHandler()
	if err != nil {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		return err
	}
	p.AddStragey(strageyHandler)
	return nil
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

// Run ...
func (p *GoPlugin) Run(v ...interface{}) interface{} {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Run fail:%v", err)
		}
	}()
	handler, ok := p.strageys[p.GetStragey()]
	if !ok {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Run fail:get name")
		return nil
	}
	return handler.Run(v)
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

// GoHandler ...
type GoHandler interface {
	LoadStragey() error
	SetStragey(string)
	GetStragey() string
	Init(v ...interface{}) interface{}
	Run(v ...interface{}) interface{}
	Exit(v ...interface{}) interface{}
}

// NewGoPlugin ...
func NewGoPlugin() *GoPlugin {
	var goplugin GoPlugin
	goplugin.strageys = make(map[string]GoStrageyHandler)
	goplugin.SetStragey("echo")
	grid, _ := NewEchoHandler()
	goplugin.AddStragey(grid)
	return &goplugin
}
