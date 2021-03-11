package trader

import (
	"context"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"snack.com/xiyanxiyan10/stocktrader/api"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
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
	return ""
}

// Switch ...
func Switch(id int64) (err error) {
	if GetTraderStatus(id) == constant.Running {
		return stop(id)
	}
	return run(id)
}

// initializeJs ...
func initializeJs(trader *Global) (err error) {
	if localErr := trader.ctx.Set("Global", api.GlobalHandler(trader)); localErr != nil {
		err = localErr
		return
	}
	if localErr := trader.ctx.Set("G", api.GlobalHandler(trader)); localErr != nil {
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
func initialize(id int64, backlog, backtest bool) (trader Global, err error) {
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

	var opt constant.Option
	opt.BackLog = backlog
	opt.BackTest = backtest
	opt.TraderID = id

	trader.Global = *(api.NewGlobalStruct(opt))
	trader.scriptType = trader.Algorithm.Type
	trader.tasks = make(Tasks)
	trader.ctx = otto.New()
	trader.ctx.Interrupt = make(chan func(), 1)

	for i, e := range es {
		opt := constant.Option{
			Index:     i,
			TraderID:  trader.ID,
			Type:      e.Type,
			Name:      e.Name,
			AccessKey: e.AccessKey,
			SecretKey: e.SecretKey,
			BackLog:   backlog,
			BackTest:  backtest,
		}
		if exchange, errD := api.GetExchange(opt); errD != nil {
			trader.es = append(trader.es, exchange)
		} else {
			fmt.Printf("make exchange fail:%s\n", errD.Error())
		}
	}

	if len(trader.es) == 0 {
		err = fmt.Errorf("please add at least one exchange")
		return
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
	trader, err := initialize(id, false, false)
	if err != nil {
		return
	}
	return runJs(trader, id)
}

// runPy ...
func runPy(trader Global, id int64) (err error) {
	go func() {
		defer func() {
			trader.Status = constant.Stop
			trader.Pending = constant.Disable
		}()
		ctx, cancel := context.WithCancel(context.Background())
		trader.cancel = cancel
		script := trader.Algorithm.Script
		str := []byte(script)
		filename := fmt.Sprintf("/tmp/%d_%s.py", trader.ID, trader.Name)
		err = ioutil.WriteFile(filename, str, 0644)
		if err != nil {
			return
		}
		cmd := exec.CommandContext(ctx, "python3", filename)
		cmd.Stdout = os.Stdout
		cmd.Start()
		trader.LastRunAt = time.Now()
		trader.Status = constant.Running

		cmd.Wait()
		fmt.Println("退出程序中...", cmd.Process.Pid)
	}()
	Executor[trader.ID] = &trader
	return
	//cancel()
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
				trader.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
			}
			// stop the cache process
			for _, e := range trader.es {
				err = e.Stop()
				if err != nil {
					trader.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
				}
			}
			if exit, err := trader.ctx.Get("exit"); err == nil && exit.IsFunction() {
				if _, err := exit.Call(exit); err != nil {
					trader.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
				}
			}
			close(trader.ctx.Interrupt)
			trader.Status = constant.Stop
			trader.Pending = constant.Disable
		}()
		trader.LastRunAt = time.Now()
		trader.Status = constant.Running
		if _, err := trader.ctx.Run(trader.Algorithm.Script); err != nil {
			trader.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
		}
		if main, err := trader.ctx.Get("main"); err != nil || !main.IsFunction() {
			trader.Log(constant.ERROR, "", 0.0, 0.0, "can not get the main function")
		} else {
			if _, err := main.Call(main); err != nil {
				trader.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
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
	return stopJs(id)
}

// stop ...
func stopJs(id int64) (err error) {
	trader := Executor[id]
	for _, e := range trader.es {
		err := e.Stop()
		if err != nil {
			return fmt.Errorf("stop exchange %s fail:%s", e.GetName(), err.Error())
		}
	}
	Executor[id].ctx.Interrupt <- func() { panic(errHalt) }
	return
}

// stop ...
func stopPy(id int64) (err error) {
	trader := Executor[id]
	trader.cancel()
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
