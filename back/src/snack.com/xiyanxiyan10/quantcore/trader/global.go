package trader

import (
	"snack.com/xiyanxiyan10/quantcore/util"
	//"encoding/json"
	//"fmt"
	"log"
	//"reflect"
	"sync"
	"time"

	"github.com/robertkrimen/otto"
	"snack.com/xiyanxiyan10/quantcore/api"
	"snack.com/xiyanxiyan10/quantcore/constant"
	"snack.com/xiyanxiyan10/quantcore/draw"
	"snack.com/xiyanxiyan10/quantcore/model"
	"snack.com/xiyanxiyan10/quantcore/notice"
)

type Tasks map[string][]task

// Global ...
type Global struct {
	model.Trader
	Logger     model.Logger   //利用这个对象保存日志
	ctx        *otto.Otto     //js虚拟机
	es         []api.Exchange //交易所列表
	tasks      Tasks          //任务列表
	running    bool
	mailNotice notice.MailNotice // 邮件发送
	lineDrawer draw.LineDrawer   // 图标绘制
	//statusLog string
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

// SetMail ...
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

// LineDrawSetPath ...
func (g *Global) LineDrawSetPath(path string) interface{} {
	g.lineDrawer.SetPath(path)
	return true
}

// LineDrawGetPath ...
func (g *Global) LineDrawGetPath() interface{} {
	return g.lineDrawer.GetPath()
}

// LineDrawReset ...
func (g *Global) LineDrawReset() interface{} {
	g.lineDrawer.Reset()
	return true
}

// LineDrawKline ...
func (g *Global) LineDrawKline(data draw.KlineData) interface{} {
	g.lineDrawer.PlotKLine(data)
	return true
}

// LineDrawKline ...
func (g *Global) LineDrawLine(name string, data draw.LineData) interface{} {
	g.lineDrawer.PlotLine(name, data)
	return true
}

// LineDrawPlot ...
func (g *Global) LineDrawPlot() interface{} {
	if err := g.lineDrawer.Draw(); err != nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
		return false
	}
	return true
}

// Console ...
func (g *Global) Console(msgs ...interface{}) {
	log.Printf("%v %v\n", constant.INFO, msgs)
}

// Log ...
func (g *Global) Log(msgs ...interface{}) {
	g.Logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// LogProfit ...
func (g *Global) LogProfit(msgs ...interface{}) {
	profit := 0.0
	if len(msgs) > 0 {
		profit = util.Float64Must(msgs[0])
	}
	g.Logger.Log(constant.PROFIT, "", 0.0, profit, msgs[1:]...)
}

// LogStatus ...
//func (g *Global) LogStatus(msgs ...interface{}) {
//	go func() {
//		msg := ""
//		for _, m := range msgs {
//			v := reflect.ValueOf(m)
//			switch v.Kind() {
//			case reflect.Struct, reflect.Map, reflect.Slice:
//				if bs, err := json.Marshal(m); err == nil {
//					msg += string(bs)
//					continue
//				}
//			}
//			msg += fmt.Sprintf("%+v", m)
//		}
//		g.statusLog = msg
//	}()
//}

// AddTask ...
func (g *Global) AddTask(group otto.Value, fn otto.Value, args ...interface{}) bool {
	if g.running {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), tasks are running")
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
