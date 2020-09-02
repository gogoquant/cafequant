package trader

import (
	"fmt"
	"strconv"
	"time"

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
		constant.HuoBiDm:     api.NewHuoBiDmExchange,
		constant.HuoBiDmBack: api.NewHuoBiDmBackExchange,
		constant.HuoBi:       api.NewHuoBiExchange,
		constant.HuoBiBack:   api.NewHuoBiBackExchange,
	}
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
	trader.mail = notice.NewMailHandler(5, 3)
	trader.draw = draw.NewDrawHandler()
	trader.ws = constant.NewWsPip(20)

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
				Ws:        trader.ws,
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
	}
	if len(trader.es) == 0 {
		err = fmt.Errorf("please add at least one exchange")
		return
	}
	trader.goplugin = goExtend
	if localErr := trader.ctx.Set("Go", &trader.goplugin); localErr != nil {
		err = localErr
		return
	}
	if localErr := trader.ctx.Set("Global", &trader); localErr != nil {
		err = localErr
		return
	}
	if localErr := trader.ctx.Set("G", &trader); localErr != nil {
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

// err2String 捕获策略的错误信息并转化为对应的字符串
func err2String(err interface{}) string {
	switch err.(type) {
	case error:
		return err.(error).Error()
	case *otto.Error:
		return err.(*otto.Error).String()
	case string:
		return err.(string)
	default:
		return "err unknown"
	}
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
				trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
			}
			if exit, err := trader.ctx.Get("exit"); err == nil && exit.IsFunction() {
				if _, err := exit.Call(exit); err != nil {
					trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err2String(err))
				}
			}
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
