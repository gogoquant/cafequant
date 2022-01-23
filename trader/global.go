package trader

import (
	"fmt"
	"sync"

	"github.com/robertkrimen/otto"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/util"
)

// Tasks ...
type Tasks map[string][]task

//js中的一个任务,目的是可以并发工作
type task struct {
	ctx  *otto.Otto    //js虚拟机
	fn   otto.Value    //代表该任务的js函数
	args []interface{} //函数的参数
}

type Plugin interface {
	Start(map[string]interface{}) error
	Process(map[string]interface{}) (interface{}, error)
	Stop(map[string]interface{}) error
}

// Global ...
type Global struct {
	api.Global
	model.Trader
	plugins map[string]Plugin

	ctx        *otto.Otto     // js虚拟机
	es         []api.Exchange // 交易所列表
	tasks      Tasks          // 任务列表
	running    bool           // 运行中
	scriptType string         // 脚本语言
}

// LoadPlugin register plugin struct
func (g *Global) LoadPlugin(pluginName, pluginType string) error {
	f, err := util.HotPlugin(config.String(fmt.Sprintf("/%s/%s.so",
		config.String(constant.GoPluginPath), pluginType)), constant.GoHandler)
	if err != nil {
		return err
	}
	makerplugin, ok := f.(func() Plugin)
	if !ok {
		return fmt.Errorf("load plugin %s fail", pluginName)
	}
	g.plugins[pluginName] = makerplugin()
	return nil
}

// ProcessPlugin call plugin after load plugin
func (g *Global) ProcessPlugin(name string, params map[string]interface{}) (interface{}, error) {
	plugin, ok := g.plugins[name]
	if !ok {
		return nil, fmt.Errorf("find plugin %s, fail", name)
	}
	return plugin.Process(params)
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
