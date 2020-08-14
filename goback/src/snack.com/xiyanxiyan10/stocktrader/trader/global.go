package trader

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/robertkrimen/otto"

	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/notice"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

type Tasks map[string][]task

// Global ...
type Global struct {
	model.Trader
	Logger     model.Logger      // 利用这个对象保存日志
	ctx        *otto.Otto        // js虚拟机
	es         []api.Exchange    // 交易所列表
	tasks      Tasks             // 任务列表
	running    bool              // 运行中
	ws         *constant.WsPiP   // 全局异步通道
	mailNotice notice.MailNotice // 邮件发送
	lineDrawer draw.LineDrawer   // 图标绘制
	statusLog  string            // 状态日志
}

//js中的一个任务,目的是可以并发工作
type task struct {
	ctx  *otto.Otto    //js虚拟机
	fn   otto.Value    //代表该任务的js函数
	args []interface{} //函数的参数
}

// Sleep ...
func (g *Global) Sleep(intervals ...interface{}) {
	interval := int64(0)
	if len(intervals) > 0 {
		interval = util.Int64Must(intervals[0])
	}
	if interval > 0 {
		time.Sleep(time.Duration(interval) * time.Millisecond)
	} else {
		for _, e := range g.es {
			e.AutoSleep()
		}
	}
}

// Wait ch ...
func (g *Global) Wait() interface{} {
	return g.ws.Pop()
}

// MailSet ...
func (g *Global) MailSet(server, portStr, username, password string) interface{} {
	port, err := util.Int(portStr)
	if err != nil {
		return false
	}
	g.mailNotice.Set(server, port, username, password)
	return true
}

// MailSend ...
func (g *Global) MailSend(msg, to string) interface{} {
	err := g.mailNotice.Send(msg, to)
	if err != nil {
		return false
	}
	return true
}

// MailStart ...
func (g *Global) MailStart() interface{} {
	err := g.mailNotice.Start()
	if err != nil {
		return false
	}
	return true
}

// SetMail ...
func (g *Global) MailStop() interface{} {
	err := g.mailNotice.Stop()
	if err != nil {
		return false
	}
	return true
}

// MailStatus ...
func (g *Global) MailStatus() interface{} {
	return g.mailNotice.Status()
}

// LineDrawSetPath set file path for config map
func (g *Global) DrawSetPath(path string) interface{} {
	g.lineDrawer.SetPath(path)
	return true
}

// LineDrawGetPath get file path from config map
func (g *Global) DrawGetPath() interface{} {
	// get the picture path
	path := g.lineDrawer.GetPath()
	if path == "" {
		path = config.String("filePath")
	}
	return path
}

// LineDrawReset ...
func (g *Global) DrawReset() interface{} {
	g.lineDrawer.Reset()
	return true
}

// LineDrawKline ...
func (g *Global) DrawKline(time string, data [4]float32) interface{} {
	var kline draw.KlineData
	kline.Time = time
	kline.Data = data
	g.lineDrawer.PlotKLine(kline)
	return true
}

// LineDrawKline ...
func (g *Global) DrawLine(name string, time string, data float32) interface{} {
	var line draw.LineData
	line.Time = time
	line.Data = data
	g.lineDrawer.PlotLine(name, line)
	return true
}

// LineDrawPlot ...
func (g *Global) DrawPlot() interface{} {
	if err := g.lineDrawer.Display(); err != nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
		return false
	}
	return true
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
		profit = util.Float64Must(messages[0])
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
