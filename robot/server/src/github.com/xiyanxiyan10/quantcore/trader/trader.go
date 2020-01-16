package trader

import (
	"fmt"
	"github.com/xiyanxiyan10/quantcore/draw"
	conver "github.com/xiyanxiyan10/quantcore/util"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/xiyanxiyan10/quantcore/api"
	"github.com/xiyanxiyan10/quantcore/constant"
	"github.com/xiyanxiyan10/quantcore/model"
)

// Trader Variable
var (
	Executor      = make(map[int64]*Global) //保存正在运行的策略，防止重复运行
	errHalt       = fmt.Errorf("HALT")
	exchangeMaker = map[string]func(constant.Option) api.Exchange{ //保存所有交易所的构造函数
		constant.HuoBiDm: api.NewHuoBiDmExchange,
		constant.Fmex:    api.NewFmexExchange,
		constant.HuoBi:   api.NewHuoBiExchange,
	}
)

// GetTraderStatus ...
func GetTraderStatus(id int64) (status int64) {
	if t, ok := Executor[id]; ok && t != nil {
		status = t.Status
	}
	return
}

// Switch ...
func Switch(id int64) (err error) {
	if GetTraderStatus(id) > 0 {
		return stop(id)
	}
	return run(id)
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
	trader.tasks = make(Tasks)
	trader.ctx = otto.New()
	trader.ctx.Interrupt = make(chan func(), 1)
	trader.mailServer = conver.NewMailServer(5, 3)
	trader.lineDrawer = draw.GetLineDrawer()
	for _, e := range es {
		if maker, ok := exchangeMaker[e.Type]; ok {
			opt := constant.Option{
				TraderID:  trader.ID,
				Type:      e.Type,
				Name:      e.Name,
				AccessKey: e.AccessKey,
				SecretKey: e.SecretKey,
			}
			trader.es = append(trader.es, maker(opt))
		}
	}
	if len(trader.es) == 0 {
		err = fmt.Errorf("please add at least one exchange")
		return
	}
	trader.ctx.Set("Global", &trader)
	trader.ctx.Set("G", &trader)
	trader.ctx.Set("Exchange", trader.es[0])
	trader.ctx.Set("E", trader.es[0])
	trader.ctx.Set("Exchanges", trader.es)
	trader.ctx.Set("Es", trader.es)
	return
}

// run ...
func run(id int64) (err error) {
	trader, err := initialize(id)
	if err != nil {
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil && err != errHalt {
				trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
			}
			if exit, err := trader.ctx.Get("exit"); err == nil && exit.IsFunction() {
				if _, err := exit.Call(exit); err != nil {
					trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
				}
			}
			trader.Status = 0
		}()
		trader.LastRunAt = time.Now()
		trader.Status = 1
		if _, err := trader.ctx.Run(trader.Algorithm.Script); err != nil {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		}
		if main, err := trader.ctx.Get("main"); err != nil || !main.IsFunction() {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, "can not get the main function")
		} else {
			if _, err := main.Call(main); err != nil {
				trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
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
	Executor[id].ctx.Interrupt <- func() { panic(errHalt) }
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
