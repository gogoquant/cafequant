package trader

import (
	"errors"
	"github.com/qiniu/py"
	"github.com/qiniu/py/pyutil"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"time"
)

// GlobalPython ...
type GlobalPython struct {
	global *Global
}

// Sleep ...
func (g *GlobalPython) Sleep(args *py.Tuple) (ret *py.Base, err error) {
	var i int64
	err = py.ParseV(args, &i)
	if err != nil {
		return
	}
	time.Sleep(time.Duration(i) * time.Millisecond)
	return py.IncNone(), nil
}

// DingSet ...
func (g *GlobalPython) DingSet(args *py.Tuple) (ret *py.Base, err error) {
	var token, key string
	err = py.ParseV(args, &token, &key)
	if err != nil {
		return
	}
	g.global.DingSet(token, key)
	return py.IncNone(), nil
}

// DingSend ...
func (g *GlobalPython) DingSend(args *py.Tuple) (ret *py.Base, err error) {
	var msg string
	err = py.ParseV(args, &msg)
	if err != nil {
		return
	}
	err = g.global.ding.Send(msg)
	if err != nil {
		return
	}
	return py.IncNone(), nil
}

// MailSet ...
func (g *GlobalPython) MailSet(args *py.Tuple) (ret *py.Base, err error) {
	var to, server, port, username, password string
	err = py.ParseV(args, &to, &server, &port, &username, &password)
	if err != nil {
		return
	}
	g.global.MailSet(to, server, port, username, password)
	return py.IncNone(), nil
}

// MailSend ...
func (g *GlobalPython) MailSend(args *py.Tuple) (ret *py.Base, err error) {
	var msg string
	err = py.ParseV(args, &msg)
	if err != nil {
		return
	}
	err = g.global.mail.Send(msg)
	if err != nil {
		return
	}
	return py.IncNone(), nil
}

// Log ...
func (g *GlobalPython) Log(args *py.Tuple) (ret *py.Base, err error) {
	var vars []interface{}
	err = py.ParseV(args, &vars)
	if err != nil {
		return
	}
	g.global.Logger.Log(constant.INFO, "", 0.0, 0.0, vars...)
	return py.IncNone(), nil
}

// DrawSetPath set file path for config map
func (g *GlobalPython) DrawSetPath(args *py.Tuple) (ret *py.Base, err error) {
	var path string
	err = py.ParseV(args, &path)
	if err != nil {
		return
	}
	g.global.draw.SetPath(path)
	return py.IncNone(), nil
}

// DrawGetPath get file path from config map
func (g *GlobalPython) DrawGetPath(args *py.Tuple) (ret *py.Base, err error) {
	path := g.global.draw.GetPath()
	if path == "" {
		path = config.String("filePath")
	}
	val, ok := pyutil.NewVar(path)
	if !ok {
		return py.IncNone(), errors.New("get newvar fail")
	}
	return val, nil
}

// DrawReset ...
func (g *GlobalPython) DrawReset(args *py.Tuple) (ret *py.Base, err error) {
	g.global.draw.Reset()
	return py.IncNone(), nil
}

// DrawKline ...
func (g *GlobalPython) DrawKline(args *py.Tuple) (ret *py.Base, err error) {
	var time string
	var kline draw.KlineData
	var a, b, c, d float32
	err = py.ParseV(args, &time, &a, &b, &c, &d)
	if err != nil {
		return
	}
	kline.Data[0] = a
	kline.Data[1] = b
	kline.Data[2] = c
	kline.Data[3] = d
	kline.Time = time
	g.global.draw.PlotKLine(kline)
	return py.IncNone(), nil
}

// DrawLine ...
func (g *GlobalPython) DrawLine(args *py.Tuple) (ret *py.Base, err error) {
	var line draw.LineData
	var time, name string
	var data float32
	err = py.ParseV(args, &name, &time, &data)
	if err != nil {
		return
	}
	line.Time = time
	line.Data = data
	g.global.draw.PlotLine(name, line)
	return py.IncNone(), nil
}

// DrawPlot ...
func (g *GlobalPython) DrawPlot(args *py.Tuple) (ret *py.Base, err error) {
	err = g.global.draw.Display()
	if err != nil {
		g.global.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return
	}
	return py.IncNone(), nil
}

// LogFile ...
func (g *GlobalPython) LogFile(args *py.Tuple) (ret *py.Base, err error) {
	var name, content string
	err = py.ParseV(args, &name, &content)
	if err != nil {
		return
	}
	ans := g.global.LogFile(name, content)
	if ans == nil {
		err = errors.New("Logfile fail")
		return
	}
	return
}
