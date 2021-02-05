package api

import (
	"os"
	"snack.com/xiyanxiyan10/conver"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/notice"
)

// GlobalHandler ...
type GlobalHandler interface {
	LogStatus(messages string)
	DingSet(token, key string) error
	DingSend(msg string) error
	MailSet(to, server, portStr, username, password string) error
	MailSend(msg string) error
	DrawSetPath(path string)
	DrawGetPath() string
	DrawReset()
	DrawKLine(time string, a, b, c, d float32)
	DrawLine(name string, time string, data float32, shape string)
	DrawPlot() error
}

// Global ...
type Global struct {
	logger model.Logger // 利用这个对象保存日志

	backtest  bool // 是否为回测模式
	backlog   bool
	mail      notice.MailHandler // 邮件发送
	ding      notice.DingHandler // dingtalk
	draw      draw.DrawHandler   // 图标绘制
	statusLog string             // 状态日志
}

// NewGlobal get global struct
func NewGlobal(opt constant.Option) GlobalHandler {
	return NewGlobalStruct(opt)
}

func NewGlobalStruct(opt constant.Option) *Global {
	var trader Global
	trader.logger = model.Logger{
		TraderID:     opt.TraderID,
		ExchangeType: "global",
	}
	trader.backtest = opt.BackTest
	trader.backlog = opt.BackLog
	trader.mail = notice.NewMailHandler()
	trader.ding = notice.NewDingHandler()
	trader.draw = draw.NewDrawHandler()
	return &trader
}

// Sleep ...
func (g *Global) Sleep(intervals int64) {
	if g.backtest {
		return
	}
	g.Sleep(intervals)
}

// DingSet ...
func (g *Global) DingSet(token, key string) error {
	g.ding.Set(token, key)
	return nil
}

// DingSend ...
func (g *Global) DingSend(msg string) error {
	return g.ding.Send(msg)
}

// MailSet ...
func (g *Global) MailSet(to, server, portStr, username, password string) error {
	port, err := conver.Int(portStr)
	if err != nil {
		return err
	}
	g.mail.Set(to, server, port, username, password)
	return nil
}

// MailSend ...
func (g *Global) MailSend(msg string) error {
	return g.mail.Send(msg)
}

// DrawSetPath set file path for config map
func (g *Global) DrawSetPath(path string) {
	g.draw.SetPath(path)
}

// DrawGetPath get file path from config map
func (g *Global) DrawGetPath() string {
	// get the picture path
	path := g.draw.GetPath()
	if path == "" {
		path = config.String("filePath")
	}
	return path
}

// DrawReset ...
func (g *Global) DrawReset() {
	g.draw.Reset()
}

// DrawKLine ...
func (g *Global) DrawKLine(time string, a, b, c, d float32) {
	g.draw.PlotKLine(time, a, b, c, d)
}

// DrawLine ...
func (g *Global) DrawLine(name string, time string, data float32, shape string) {
	g.draw.PlotLine(name, time, data, shape)
}

// DrawPlot ...
func (g *Global) DrawPlot() error {
	if err := g.draw.Display(); err != nil {
		g.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return err
	}
	return nil
}

// Log ...
func (g *Global) Log(messages string) {
	g.logger.Log(constant.INFO, "", 0.0, 0.0, messages)
}

// LogStatus ...
func (g *Global) LogStatus(messages string) {
	g.statusLog = messages
}

// LogFile ...
func (g *Global) LogFile(name, strContent string) error {
	fd, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		g.logger.Log(constant.ERROR, "", 0.0, 0.0, "Can not open the file:", err)
		return err
	}
	fdContent := strContent
	buf := []byte(fdContent)
	_, err = fd.Write(buf)
	if err != nil {
		g.logger.Log(constant.ERROR, "", 0.0, 0.0, "Can not write the file:", err)
		return err
	}
	err = fd.Close()
	if err != nil {
		g.logger.Log(constant.ERROR, "", 0.0, 0.0, "Can not close the file:", err)
		return err
	}
	return nil
}
