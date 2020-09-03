package pythonplugin

import (
	"github.com/sbinet/go-python"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
)

// PythonPlugin  ...
type PythonPlugin struct {
	name   string
	env    map[string]*python.PyObject
	store  map[string]*python.PyObject
	Logger *model.Logger // 利用这个对象保存日志
}

// SetStragey ...
func (p *PythonPlugin) SetStragey(name string) {
	p.name = name
}

// GetStragey ...
func (p *PythonPlugin) GetStragey() string {
	return p.name
}

// AddLogger ...
func (p *PythonPlugin) AddLogger(logger *model.Logger) {
	p.Logger = logger
}

// Handler ...
type Handler interface {
	LoadStragey() interface{}
	SetStragey(string)
	GetStragey() string
	Init(v ...interface{}) interface{}
	Call(name string, v ...interface{}) interface{}
	Exit(v ...interface{}) interface{}
}

// NewPythonPlugin ...
func NewPythonPlugin() *PythonPlugin {
	var goplugin PythonPlugin
	goplugin.store = make(map[string]*python.PyObject)
	goplugin.store = make(map[string]*python.PyObject)
	return &goplugin
}

// LoadStragey ...
func (g *PythonPlugin) LoadStragey() (res interface{}) {
	res = nil
	defer func() {
		if err := recover(); err != nil {
			g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "LoadStragey() fail:%v", err)
		}
	}()
	name := g.GetStragey()
	m := python.PyImport_ImportModule(name)
	if m != nil {
		g.env[name] = m
	}
	callback := m.GetAttrString("New")
	if callback == nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail")
		return nil
	}
	handler := callback.CallFunction()
	if handler == nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail")
		return nil
	}
	g.AddStragey(handler)
	return "success"
}

// AddStragey ...
func (g *PythonPlugin) AddStragey(v *python.PyObject) {
	g.store[g.GetStragey()] = v
}

// Exit ...
func (g *PythonPlugin) Exit(v ...interface{}) (res interface{}) {
	res = nil
	defer func() {
		if err := recover(); err != nil {
			g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail:%v", err)
		}
	}()
	m, ok := g.store[g.GetStragey()]
	if !ok {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail")
		return nil
	}
	callback := m.GetAttrString("Exit")
	if callback == nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail")
		return nil
	}
	out := callback.CallFunction()
	if out == nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail")
		return nil
	}
	return python.PyString_AsString(out)
}

// Init ...
func (g *PythonPlugin) Init(v ...interface{}) (res interface{}) {
	res = nil
	defer func() {
		if err := recover(); err != nil {
			g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail:%v", err)
		}
	}()
	m, ok := g.store[g.GetStragey()]
	if !ok {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail:%v")
		return nil
	}
	callback := m.GetAttrString("Init")
	if callback == nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail:%v")
		return nil
	}
	out := callback.CallFunction()
	if out == nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Init fail:%v")
		return nil
	}
	return python.PyString_AsString(out)
}

// Call ...
func (g *PythonPlugin) Call(name string, v ...interface{}) interface{} {
	defer func() {
		if err := recover(); err != nil {
			g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Call fail:%v", err)
		}
	}()
	m, ok := g.store[g.GetStragey()]
	if !ok {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Call fail")
		return nil
	}
	callback := m.GetAttrString("Call")
	if callback == nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Call fail")
		return nil
	}
	out := callback.CallFunction(name, v)
	if out == nil {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Stragey Call fail")
		return nil
	}
	return python.PyString_AsString(out)
}

func init() {
	if err := python.Initialize(); err != nil {
		panic(err)
	}
}
