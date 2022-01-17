package trader

import (
	"context"
	"sync"

	talib "github.com/markcheno/go-talib"
	"github.com/robertkrimen/otto"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
)

// Tasks ...
type Tasks map[string][]task

//js中的一个任务,目的是可以并发工作
type task struct {
	ctx  *otto.Otto    //js虚拟机
	fn   otto.Value    //代表该任务的js函数
	args []interface{} //函数的参数
}

// Global ...
type Global struct {
	api.Global
	model.Trader

	ctx        *otto.Otto         // js虚拟机
	es         []api.Exchange     // 交易所列表
	tasks      Tasks              // 任务列表
	running    bool               // 运行中
	scriptType string             // 脚本语言
	cancel     context.CancelFunc // 执行python脚本
}

// GroupCandles group candles
func (g *Global) GroupCandles(highs []float64, opens []float64, closes []float64, lows []float64, groupingFactor int) ([]float64, []float64, []float64, []float64, error) {
	return talib.GroupCandles(highs, opens, closes, lows, groupingFactor)
}

// AddTask ...
func (g *Global) AddTask(group otto.Value, fn otto.Value, args ...interface{}) bool {
	if g.running {
		g.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), Tasks are running")
		return false
	}
	if !group.IsString() {
		g.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), Invalid group name")
		return false
	}
	if !fn.IsString() {
		g.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), Invalid function name")
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
		g.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), tasks are running")
		return false
	}
	if !group.IsString() {
		g.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), Invalid group name")
		return false
	}
	if !fn.IsString() {
		g.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), Invalid function name")
		return false
	}
	if _, ok := g.tasks[group.String()]; !ok {
		g.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), group not exist")
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
	g.Log(constant.ERROR, "", 0.0, 0.0, "BindTaskParam(), function not exist")
	return false
}

// ExecTasks ...
func (g *Global) ExecTasks(group otto.Value) (results []interface{}) {
	if !group.IsString() {
		g.Log(constant.ERROR, "", 0.0, 0.0, "ExecTasks(), Invalid group name")
		return
	}
	if _, ok := g.tasks[group.String()]; !ok {
		g.Log(constant.ERROR, "", 0.0, 0.0, "ExecTasks(), group not exist")
		return
	}
	if g.running {
		g.Log(constant.ERROR, "", 0.0, 0.0, "ExecTasks(), tasks are running")
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
				g.Log(constant.ERROR, "", 0.0, 0.0, "Can not get the task function")
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
