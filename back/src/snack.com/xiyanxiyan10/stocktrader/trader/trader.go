package trader

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"reflect"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"snack.com/xiyanxiyan10/stocktrader/goplugin"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/notice"
	"strconv"
	"time"
)

// Trader Variable
var (
	Executor = make(map[int64]*Global) //保存正在运行的策略，防止重复运行
	errHalt  = fmt.Errorf("HALT")
)

// GetTraderStatus ...
func GetTraderStatus(id int64) (status int64) {
	if t, ok := Executor[id]; ok && t != nil {
		status = t.Status
		// show pending in status
		if t.Pending == constant.Enable {
			status = constant.Pending
		}
	}
	return
}

// GetTraderLogStatus ...
func GetTraderLogStatus(id int64) (status string) {
	if t, ok := Executor[id]; ok && t != nil {
		return t.statusLog
	}
	return ""
}

// Switch ...
func Switch(id int64) (err error) {
	if GetTraderStatus(id) == constant.Running {
		return stop(id)
	}
	return run(id)
}

// initializeGo ...
func initializeGo(trader *Global) (err error) {
	return
}

// run ...
func runGo(trader Global, id int64) (err error) {
	err = initializeGo(&trader)
	if err != nil {
		return
	}

	go func() {
		defer func() {
			if err := recover(); err != nil && err != errHalt {
				trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
			}
			close(trader.ctx.Interrupt)
			trader.Status = constant.Stop
			trader.Pending = constant.Disable
		}()
		scripts := trader.Algorithm.Script
		p := make(map[string]string)
		err := json.Unmarshal([]byte(scripts), &p)
		if err != nil {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
			return
		}
		name := p["name"]
		trader.goplugin.SetStragey(name)
		trader.LastRunAt = time.Now()
		trader.Status = constant.Running
		err = trader.goplugin.LoadStragey()
		if err != nil {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
			return
		}
		err = trader.goplugin.Init(p)
		if err != nil {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
			return
		}
		err = trader.goplugin.Run(p)
		if err != nil {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err.Error())
			return
		}

	}()
	Executor[trader.ID] = &trader
	return
}

// initializeJs ...
func initializeJs(trader *Global) (err error) {
	if localErr := trader.ctx.Set("Go", trader.goplugin); localErr != nil {
		err = localErr
		return
	}
	if localErr := trader.ctx.Set("Global", GlobalHandler(trader)); localErr != nil {
		err = localErr
		return
	}
	if localErr := trader.ctx.Set("G", GlobalHandler(trader)); localErr != nil {
		err = localErr
		return
	}

	if localErr := trader.ctx.Set("Plugin", trader.goplugin); localErr != nil {
		err = localErr
		return
	}

	if localErr := trader.ctx.Set("Exchange", trader.es[0]); localErr != nil {
		err = localErr
		return
	}
	if localErr := trader.ctx.Set("E", trader.es[0]); localErr != nil {
		err = localErr
		return
	}
	if localErr := trader.ctx.Set("Exchanges", trader.es); localErr != nil {
		err = localErr
		return
	}
	if localErr := trader.ctx.Set("Es", trader.es); localErr != nil {
		err = localErr
		return
	}
	return
}

//initialize 核心是初始化js运行环境，及其可以调用的api
func initialize(id int64) (trader Global, err error) {
	if t := Executor[id]; t != nil && t.Status > 0 {
		return
	}
	err = model.DB.First(&trader.Trader, id).Error
	if err != nil {
		return
	}
	self, err := model.GetUserByID(trader.UserID)
	if err != nil {
		return
	}
	if trader.AlgorithmID <= 0 {
		err = fmt.Errorf("please select a algorithm")
		return
	}
	err = model.DB.First(&trader.Algorithm, trader.AlgorithmID).Error
	if err != nil {
		return
	}
	es, err := self.GetTraderExchanges(trader.ID)
	if err != nil {
		return
	}
	trader.Logger = model.Logger{
		TraderID:     trader.ID,
		ExchangeType: "global",
	}

	trader.scriptType = trader.Algorithm.Type
	trader.tasks = make(Tasks)
	trader.ctx = otto.New()
	trader.ctx.Interrupt = make(chan func(), 1)
	trader.mail = notice.NewMailHandler()
	trader.ding = notice.NewDingHandler()
	trader.draw = draw.NewDrawHandler()

	// set the diagram path
	filePath := config.String(constant.FilePath)
	trader.draw.SetPath(filePath + "/" + strconv.FormatInt(trader.ID, 10) + ".html")

	goExtend := goplugin.NewGoPlugin()
	goExtend.AddMail(trader.mail)
	goExtend.AddDraw(trader.draw)
	goExtend.AddDing(trader.ding)
	goExtend.AddLogStatus(&trader.statusLog)
	goExtend.AddLog(&trader.Logger)

	for i, e := range es {
		opt := constant.Option{
			Index:     i,
			TraderID:  trader.ID,
			Type:      e.Type,
			Name:      e.Name,
			AccessKey: e.AccessKey,
			SecretKey: e.SecretKey,
			LogBack:   false,
		}
		if maker, ok := api.ExchangeMaker[e.Type]; ok {
			exchange, errD := maker(opt)
			if errD != nil {
				err = errD
				return
			}
			goExtend.AddExchange(exchange)
			trader.es = append(trader.es, exchange)
		}
	}
	trader.goplugin = goExtend

	if len(trader.es) == 0 {
		err = fmt.Errorf("please add at least one exchange")
		return
	}
	var backtot = 0
	for i := range trader.es {
		if trader.es[i].IsBack() {
			backtot++
		}
	}
	if backtot == 0 {
		trader.back = false
	} else if len(trader.es) == backtot {
		trader.back = true
	} else {
		err = fmt.Errorf("please use exchanges all back or all online")
	}
	return
}

// err2String 捕获策略的错误信息并转化为对应的字符串
func err2String(err interface{}) string {
	switch err.(type) {
	case *otto.Error:
		return err.(*otto.Error).String()
	case error:
		return err.(error).Error()
	case string:
		return err.(string)
	default:
		return "err unknown type:" + reflect.TypeOf(err).Name()
	}
}

// run ...
func run(id int64) (err error) {
	trader, err := initialize(id)
	if err != nil {
		return
	}

	switch trader.scriptType {
	case constant.ScriptGo:
		return runGo(trader, id)
	case constant.ScriptJs:
		return runJs(trader, id)

	}
	return runJs(trader, id)
}

// runJs ...
func runJs(trader Global, id int64) (err error) {
	err = initializeJs(&trader)
	if err != nil {
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil && err != errHalt {
				trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
			}
			// stop the cache process
			for _, e := range trader.es {
				err = e.Stop()
				if err != nil {
					trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
				}
			}
			if exit, err := trader.ctx.Get("exit"); err == nil && exit.IsFunction() {
				if _, err := exit.Call(exit); err != nil {
					trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
				}
			}
			close(trader.ctx.Interrupt)
			trader.Status = constant.Stop
			trader.Pending = constant.Disable
		}()
		trader.LastRunAt = time.Now()
		trader.Status = constant.Running
		if _, err := trader.ctx.Run(trader.Algorithm.Script); err != nil {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
		}
		if main, err := trader.ctx.Get("main"); err != nil || !main.IsFunction() {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, "can not get the main function")
		} else {
			if _, err := main.Call(main); err != nil {
				trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
			}
		}
	}()
	Executor[trader.ID] = &trader
	return
}

// getStatus ...
//func getStatus(id int64) (status string) {
//	if t := Executor[id]; t != nil {
//		status = t.statusLog
//	}
//	return
//}

// stop ...
func stop(id int64) (err error) {
	if t, ok := Executor[id]; !ok || t == nil {
		return fmt.Errorf("can not found the Trader")
	}
	if Executor[id].Pending == constant.Enable {
		return fmt.Errorf("pending Trader")
	}
	Executor[id].Pending = constant.Enable
	trader := Executor[id]
	for _, e := range trader.es {
		err := e.Stop()
		if err != nil {
			return fmt.Errorf("stop exchange %s fail:%s", e.GetName(), err.Error())
		}
	}
	switch trader.scriptType {
	case constant.ScriptGo:
		return stopGo(id)
	case constant.ScriptJs:
		return stopJs(id)
	}
	return
}

// stop ...
func stopJs(id int64) (err error) {
	Executor[id].ctx.Interrupt <- func() { panic(errHalt) }
	return
}

// stop ...
func stopGo(id int64) (err error) {
	err = Executor[id].goplugin.Exit(nil)
	if err != nil {
		return err
	}
	return
}

// clean ...
//func clean(userID int64) {
//	for _, t := range Executor {
//		if t != nil && t.UserID == userID {
//			stop(t.ID)
//		}
//	}
//}
