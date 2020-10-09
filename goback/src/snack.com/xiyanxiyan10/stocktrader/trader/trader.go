package trader

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/qiniu/py"
	"github.com/qiniu/py/pyutil"
	"github.com/qiniu/x/log"
	"github.com/robertkrimen/otto"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/config"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/draw"
	"snack.com/xiyanxiyan10/stocktrader/goplugin"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/notice"
)

// Trader Variable
var (
	Executor      = make(map[int64]*Global) //保存正在运行的策略，防止重复运行
	errHalt       = fmt.Errorf("HALT")
	exchangeMaker = map[string]func(constant.Option) api.Exchange{ //保存所有交易所的构造函数
		constant.HuoBiDm:    api.NewHuoBiDmExchange,
		constant.HuoBi:      api.NewHuoBiExchange,
		constant.SZ:         api.NewSZExchange,
		constant.SpotBack:   api.NewSpotBackExchange,
		constant.FutureBack: api.NewFutureBackExchange,
	}
	pyexchangeMaker map[string]func(constant.Option) api.ExchangePython
)

// GetTraderStatus ...
func GetTraderStatus(id int64) (status int64) {
	if t, ok := Executor[id]; ok && t != nil {
		status = t.Status
		// show pending in status
		if t.Pending == 1 {
			status = -1
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
	if GetTraderStatus(id) > 0 {
		return stop(id)
	}
	return run(id)
}

// initializePy ...
func initializePy(trader *Global) (err error) {
	return
}

// run ...
func runPy(trader Global, id int64) (err error) {
	err = initializePy(&trader)
	if err != nil {
		return
	}
	runtime.GOMAXPROCS(1)

	go func() {
		gomod, err := py.NewGoModule("exchange", "", trader.espy[0])
		if err != nil {
			log.Fatal("NewGoModule failed:", err)
			return
		}
		defer gomod.Decref()

		gomode, err := py.NewGoModule("E", "", trader.espy[0])
		if err != nil {
			log.Fatal("NewGoModule failed:", err)
			return
		}
		defer gomode.Decref()

		var gpy GlobalPython
		gpy.global = &trader
		globalmode, err := py.NewGoModule("G", "", gpy)
		if err != nil {
			log.Fatal("NewGoModule failed:", err)
			return
		}
		defer globalmode.Decref()

		code, err := py.Compile(trader.Algorithm.Script, "", py.FileInput)
		if err != nil {
			log.Fatal("Compile failed:", err)
			return
		}
		defer code.Decref()

		mod, err := py.ExecCodeModule("pycode", code.Obj())
		if err != nil {
			log.Fatal("ExecCodeModule failed:", err)
			return
		}
		defer mod.Decref()

		defer func() {
			if err := recover(); err != nil && err != errHalt {
				trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
			}
			ret, err := pyutil.CallMethod(mod.Obj(), "exit")
			if err != nil {
				log.Fatal("exit failed:", err)
			}
			defer ret.Decref()
			trader.ws.Close()
			close(trader.ctx.Interrupt)
			trader.Status = 0
			trader.Pending = 0
		}()
		trader.LastRunAt = time.Now()
		trader.Status = 1
		ret, err := pyutil.CallMethod(mod.Obj(), "main")
		if err != nil {
			log.Fatal("Call main failed:", err)
		}
		defer ret.Decref()
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
	trader.ws = constant.NewWsPIP(20)

	// set the diagram path
	filePath := config.String(constant.FilePath)
	trader.draw.SetPath(filePath + "/" + strconv.FormatInt(trader.ID, 10) + ".html")

	goExtend := goplugin.NewGoPlugin()
	goExtend.AddMail(trader.mail)
	goExtend.AddDraw(trader.draw)
	goExtend.Logger = &trader.Logger
	for i, e := range es {
		if maker, ok := exchangeMaker[e.Type]; ok {
			opt := constant.Option{
				Index:     i,
				TraderID:  trader.ID,
				Type:      e.Type,
				Name:      e.Name,
				AccessKey: e.AccessKey,
				SecretKey: e.SecretKey,
			}
			exchange := maker(opt)
			goExtend.AddExchange(exchange)
			trader.es = append(trader.es, exchange)
		}
		if maker, ok := pyexchangeMaker[e.Type]; ok {
			opt := constant.Option{
				Index:     i,
				TraderID:  trader.ID,
				Type:      e.Type,
				Name:      e.Name,
				AccessKey: e.AccessKey,
				SecretKey: e.SecretKey,
			}
			exchange := maker(opt)
			trader.espy = append(trader.espy, exchange)
		}
	}
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
	trader.goplugin = goExtend
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

// runCheck ...
func runCheck(id int64, script string) (err error) {
	if script != constant.ScriptPython {
		return
	}
	for i := range Executor {
		t := Executor[i]
		if t != nil {
			continue
		}
		if t.Status < 1 {
			continue
		}
		if t.scriptType == constant.ScriptPython {
			err = errors.New("python scripts only run one")
			return
		}
	}
	return
}

// run ...
func run(id int64) (err error) {
	trader, err := initialize(id)
	if err != nil {
		return
	}
	err = runCheck(id, trader.scriptType)
	if err != nil {
		return
	}

	switch trader.scriptType {
	case constant.ScriptPython:
		return runPy(trader, id)
	case constant.ScriptJs:
	default:
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
			if exit, err := trader.ctx.Get("exit"); err == nil && exit.IsFunction() {
				if _, err := exit.Call(exit); err != nil {
					trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
				}
			}
			trader.ws.Close()
			close(trader.ctx.Interrupt)
			trader.Status = 0
			trader.Pending = 0
		}()
		trader.LastRunAt = time.Now()
		trader.Status = 1
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
	if Executor[id].Pending == 1 {
		return fmt.Errorf("pending Trader")
	}
	Executor[id].ctx.Interrupt <- func() { panic(errHalt) }
	Executor[id].Pending = 1
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

func init() {
	pyexchangeMaker = make(map[string]func(constant.Option) api.ExchangePython)
	for key, funcs := range exchangeMaker {
		pyexchangeMaker[key] = api.NewExchangePython(funcs)
	}
}
