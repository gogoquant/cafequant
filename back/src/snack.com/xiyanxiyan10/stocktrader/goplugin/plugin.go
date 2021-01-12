package goplugin

import (
	"encoding/json"
	"fmt"
	"plugin"
	"reflect"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/notice"
)

// Gofunc ...
type Gofunc func(...interface{}) error

// StrageyNew ...
type StrageyNew func(...interface{}) (GoStrageyHandler, error)

// GoStragey  ...
type GoStragey struct {
	Exchanges []api.Exchange
	Mail      notice.MailHandler // 邮件发送
	Ding      notice.DingHandler // dingding
	Draw      draw.DrawHandler   // 图标绘制
	Logger    *model.Logger      // 利用这个对象保存日志
	Status    *string            // 利用该字段改写状态日志
}

// AddExchange ...
func (p *GoStragey) AddExchange(e ...api.Exchange) {
	p.Exchanges = append(p.Exchanges, e...)
}

// AddLogStatus ...
func (p *GoStragey) AddLogStatus(s *string) {
	p.Status = s
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

// LogStatus ...
func (g *GoStragey) LogStatus(messages ...interface{}) {
	go func() {
		msg := ""
		for _, m := range messages {
			v := reflect.ValueOf(m)
			switch v.Kind() {
			case reflect.Struct, reflect.Map, reflect.Slice:
				if bs, err := json.Marshal(m); err == nil {
					msg += string(bs)
					continue
				}
			}
			msg += fmt.Sprintf("%+v", m)
		}
		*(g.Status) = msg
	}()
}

// GoStrageyHandler ...
type GoStrageyHandler interface {
	AddExchange(...api.Exchange)
	AddDraw(draw draw.DrawHandler)
	AddMail(mail notice.MailHandler)
	AddDing(ding notice.DingHandler)
	AddLogger(logger *model.Logger)
	AddLogStatus(s *string)

	Init(map[string]string) error
	Run(map[string]string) error
	Exit(map[string]string) error
}

// GoPlugin ...
type GoPlugin struct {
	name      string
	Logger    *model.Logger               // 利用这个对象保存日志
	Exchanges []api.Exchange              // 交易所列表
	Mail      notice.MailHandler          // 邮件发送
	Draw      draw.DrawHandler            // 图标绘制
	Ding      notice.DingHandler          // dingding
	LogStatus *string                     // status
	strageys  map[string]GoStrageyHandler // 策略集合
}

// SetStragey ...
func (p *GoPlugin) SetStragey(name string) {
	p.name = name
}

// AddLogStatus ...
func (p *GoPlugin) AddLogStatus(s *string) {
	p.LogStatus = s
}

// AddLog ...
func (p *GoPlugin) AddLog(l *model.Logger) {
	p.Logger = l
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
	v.AddLogStatus(p.LogStatus)
	v.AddDing(p.Ding)
}

// LoadStragey ...
func (p *GoPlugin) LoadStragey() error {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		}
	}()
	// name is the so and handler new key
	name := p.GetStragey()
	handler, err := plugin.Open(config.String(constant.GoPluginPath) + "/" + name + ".so")
	if err != nil {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		return err
	}
	s, err := handler.Lookup(constant.GoHandler)
	fmt.Printf("lookup hanlder %s :%v\n", constant.GoHandler, s)
	if err != nil {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		return err
	}
	if s == nil {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail interface is nil")
		return fmt.Errorf("LoadStragey() fail interface is nil")
	}
	newHandler, ok := s.(func() (GoStrageyHandler, error))
	if !ok {
		t := reflect.TypeOf(s)
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail convert handler type:"+t.Name())
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
func (p *GoPlugin) Init(v map[string]string) error {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail:%v", err)
		}
	}()
	handler, ok := p.strageys[p.GetStragey()]
	if !ok {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail "+p.GetStragey())
		return fmt.Errorf("Stragey Exit fail get name " + p.GetStragey())
	}
	return handler.Init(v)
}

// Run ...
func (p *GoPlugin) Run(v map[string]string) error {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Run fail")
		}
	}()
	handler, ok := p.strageys[p.GetStragey()]
	if !ok {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Run fail "+p.GetStragey())
		return fmt.Errorf("Stragey Run fail get name " + p.GetStragey())
	}
	return handler.Run(v)
}

// Exit ...
func (p *GoPlugin) Exit(v map[string]string) error {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Exit fail")
		}
	}()
	handler, ok := p.strageys[p.GetStragey()]
	if !ok {
		p.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Exit fail get name "+p.GetStragey())
		return fmt.Errorf("Stragey Exit fail get name " + p.GetStragey())
	}
	return handler.Exit(v)
}

// GoHandler ...
type GoHandler interface {
	LoadStragey() error
	SetStragey(string)
	GetStragey() string
	Init(v map[string]string) error
	Run(v map[string]string) error
	Exit(v map[string]string) error
}

// NewGoPlugin ...
func NewGoPlugin() *GoPlugin {
	var goplugin GoPlugin
	goplugin.strageys = make(map[string]GoStrageyHandler)
	return &goplugin
}
