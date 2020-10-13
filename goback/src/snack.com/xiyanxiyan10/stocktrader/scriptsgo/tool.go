package main

import (
	"fmt"
	"snack.com/xiyanxiyan10/stocktrader/constant"
	"snack.com/xiyanxiyan10/stocktrader/model"
	"snack.com/xiyanxiyan10/stocktrader/trader"
	/*
		"errors"
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
		"snack.com/xiyanxiyan10/stocktrader/draw"
		"snack.com/xiyanxiyan10/stocktrader/goplugin"
		"snack.com/xiyanxiyan10/stocktrader/model"
		"snack.com/xiyanxiyan10/stocktrader/notice"
	*/)

func main() {
	var logger model.Logger
	logger.Back = true
	var opt constant.Option
	var Constract = "quarter"
	var Symbol = "BTC/USD"
	var Period = "M30"
	var IO = "online"
	var global trader.Global
	global.Logger = logger
	opt.AccessKey = ""
	opt.SecretKey = ""
	opt.Name = constant.HuoBiDm
	opt.TraderID = 1
	opt.Type = constant.HuoBiDm
	opt.Index = 1
	opt.LogBack = true

	maker := trader.ExchangeMaker[opt.Type]
	huobiExchange := maker(opt)

	huobiExchange.SetIO(IO)
	huobiExchange.SetContractType(Constract)
	huobiExchange.SetStockType(Symbol)
	huobiExchange.Ready()

	for {
		records, err := huobiExchange.GetRecords(Period, "")
		if err != nil {
			fmt.Printf("get records fail:%v", err)
		}
		fmt.Printf("get records %v", records)
		global.Sleep(1000)
	}

}
