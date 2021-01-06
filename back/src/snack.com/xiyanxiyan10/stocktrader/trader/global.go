package trader

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"log"
	"os"
	"reflect"
	"snack.com/xiyanxiyan10/conver"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"snack.com/xiyanxiyan10/stocktrader/goplugin"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/notice"
	"sync"
	"time"
)

// Tasks ...
type Tasks map[string][]task

// GlobalHandler ...
type GlobalHandler interface {
	DingSet(token, key string) error
	DingSend(msg string) error
	MailSet(to, server, portStr, username, password string) error
	MailSend(msg string) error
	DrawSetPath(path string)
	DrawGetPath() string
	DrawReset()
	DrawKLine(time string, a, b, c, d float32)
	DrawLine(name string, time string, data float32)
	DrawPlot() error
	Wait() int
}

// Global ...
type Global struct {
	model.Trader
	back       bool                // 是否为回测模式
	Logger     model.Logger        // 利用这个对象保存日志
	ctx        *otto.Otto          // js虚拟机
	es         []api.Exchange      // 交易所列表
	tasks      Tasks               // 任务列表
	running    bool                // 运行中
	scriptType string              // 脚本语言
	ws         constant.PIPHandler // websocket
	mail       notice.MailHandler  // 邮件发送
	ding       notice.DingHandler  // dingtalk
	draw       draw.DrawHandler    // 图标绘制
	goplugin   goplugin.GoHandler  // go 插件
	statusLog  string              // 状态日志
}

// Sleep ...
func (g *Global) Sleep(intervals ...interface{}) {
	if g.back {
		return
	}
	interval := int64(0)
	if len(intervals) > 0 {
		interval = conver.Int64Must(intervals[0])
	}
	if interval > 0 {
		time.Sleep(time.Duration(interval) * time.Millisecond)
	} else {
		for _, e := range g.es {
			e.AutoSleep()
		}
	}
}

// Wait ...
func (g *Global) Wait() int {
	val := g.ws.Pop()
	return val
}

// GetMail ...
func (g *Global) GetMail() notice.MailHandler {
	return g.mail
}

// GetDraw ...
func (g *Global) GetDraw() draw.DrawHandler {
	return g.draw
}

//js中的一个任务,目的是可以并发工作
type task struct {
	ctx  *otto.Otto    //js虚拟机
	fn   otto.Value    //代表该任务的js函数
	args []interface{} //函数的参数
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
func (g *Global) DrawLine(name string, time string, data float32) {
	g.draw.PlotLine(name, time, data)
}

// DrawPlot ...
func (g *Global) DrawPlot() error {
	if err := g.draw.Display(); err != nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return err
	}
	return nil
}

// Console ...
func (g *Global) Console(messages ...interface{}) {
	log.Printf("%v %v\n", constant.INFO, messages)
}

// Log ...
func (g *Global) Log(messages ...interface{}) {
	g.Logger.Log(constant.INFO, "", 0.0, 0.0, messages...)
}

// LogProfit ...
func (g *Global) LogProfit(messages ...interface{}) {
	profit := 0.0
	if len(messages) > 0 {
		profit = conver.Float64Must(messages[0])
	}
	g.Logger.Log(constant.PROFIT, "", 0.0, profit, messages[1:]...)
}

// LogStatus ...
func (g *Global) LogStatus(messages ...interface{}) {
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
		g.statusLog = msg
	}()
}

// AddTask ...
func (g *Global) AddTask(group otto.Value, fn otto.Value, args ...interface{}) bool {
	if g.running {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), Tasks are running")
		return false
	}
	if !group.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), Invalid group name")
		return false
	}
	if !fn.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), Invalid function name")
		return false
	}
	if _, ok := g.tasks[group.String()]; !ok {
		g.tasks[group.String()] = []task{}
	}
	t := task{ctx: g.ctx.Copy(), fn: fn, args: args}
	t.ctx.Interrupt = make(chan func(), 1)
	g.tasks[group.String()] = append(g.tasks[group.String()], t)
	return true
}

// BindTaskParam ...
func (g *Global) BindTaskParam(group otto.Value, fn otto.Value, args ...interface{}) bool {
	if g.running {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), tasks are running")
		return false
	}
	if !group.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), Invalid group name")
		return false
	}
	if !fn.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), Invalid function name")
		return false
	}
	if _, ok := g.tasks[group.String()]; !ok {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), group not exist")
		return false
	}
	ts := g.tasks[group.String()]
	for i := 0; i < len(ts); i++ {
		t := &ts[i]
		if t.fn.String() == fn.String() {
			t.args = args
			return true
		}
	}
	g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), function not exist")
	return false
}

// ExecTasks ...
func (g *Global) ExecTasks(group otto.Value) (results []interface{}) {
	if !group.IsString() {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "ExecTasks(), Invalid group name")
		return
	}
	if _, ok := g.tasks[group.String()]; !ok {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "ExecTasks(), group not exist")
		return
	}
	if g.running {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "ExecTasks(), tasks are running")
		return
	}
	g.running = true
	ts := g.tasks[group.String()]
	for range ts {
		results = append(results, false)
	}
	wg := sync.WaitGroup{}
	for i, t := range ts {
		wg.Add(1)
		go func(i int, t task) {
			if f, err := t.ctx.Get(t.fn.String()); err != nil || !f.IsFunction() {
				g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Can not get the task function")
			} else {
				result, err := f.Call(f, t.args...)
				if err != nil || result.IsUndefined() || result.IsNull() {
					results[i] = false
				} else {
					results[i] = result
				}
			}
			wg.Done()
		}(i, t)
	}
	wg.Wait()
	g.running = false
	return
}

// LogFile ...
func (g *Global) LogFile(name, strContent string) error {
	fd, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Can not open the file:", err)
		return err
	}
	fdContent := strContent
	buf := []byte(fdContent)
	_, err = fd.Write(buf)
	if err != nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Can not write the file:", err)
		return err
	}
	err = fd.Close()
	if err != nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Can not close the file:", err)
		return err
	}
	return nil
}
