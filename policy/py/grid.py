#!/usr/bin/python3
'''
policy for grid trade
'''
import config
import api
import constant

config.Init("../config.ini")
opt = constant.Option()
opt.io = constant.IONONE
opt.BackLog = True
opt.BackTest = True
opt.Index = 1
opt.TraderID = 1
opt.Name = "trend"
opt.BackExit = True
opt.Type = constant.HuoBiDm
symbol = 'BTC/USD.quarter'
period = 'M15'
periodsize  = 3
debug = True
limit = 1000

try:
    ex = api.GetExchange(opt)
except Exception as err:
    print(err)


try :
    tool = api.NewGlobal(opt)
except Exception as err:
    print(err)

ex.SetStockType(symbol)
ex.SetPeriod(period)
ex.SetLimit(limit)
ex.SetPeriodSize(periodsize)
ex.SetIO(opt.io)
ex.SetSubscribe(symbol, constant.CacheAccount)
ex.SetSubscribe(symbol, constant.CacheRecord)
#ex.SetSubscribe(symbol, constant.CachePosition)
ex.SetSubscribe(symbol, constant.CacheOrder)
ex.SetSubscribe(symbol, constant.CacheTicker)

time_range = ex.BackGetTimeRange()

ex.SetBackCommission(0.0, 0.0, 100, 100)
ex.SetBackTime(time_range[0], time_range[1], ex.GetPeriod())
ex.SetBackAccount('BTC', 10000)
ex.SetBackAccount('USD', 10000)
ex.SetMarginLevel(1.0)
ex.Start()

def openFunc(price, amount, msg, d):
    if d == 0:
        if debug:
            print("open long %s %s %s" % (price, amount, msg))
            return
        ex.SetDirection(constant.TradeTypeLong)
        ex.Buy(str(price), str(amount), msg)
    else:
        if debug:
            print("open short %s %s %s" % (price, amount, msg))
            return
        ex.SetDirection(constant.TradeTypeShort)
        ex.Sell(str(price), str(amount), msg)

def closeFunc(price, amount, msg, d):
    if d == 0:
        if debug:
            print("close long %s %s %s" %  (price, amount, msg))
            return
        ex.SetDirection(constant.TradeTypeLongClose)
        ex.Sell(str(price), str(amount), msg)
    else:
        if debug:
            print("close short %s %s %s" % (price, amount, msg))
            return
        ex.SetDirection(constant.TradeTypeShortClose)
        ex.Buy(str(price), str(amount), msg)

while True:
    records = ex.GetRecords()
    if len(records) > 1:
        print(records[-1])
    ticker = ex.GetTicker()
    price = 3000
    amount = 0.001
    openFunc(str(price),str(amount), 'open long', 0)
    closeFunc(str(price),str(amount), 'close long', 0)
    print(ticker.Last)
