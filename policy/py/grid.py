import config
import api
import constant
import pdb

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

try:
    ex = api.GetExchange(opt)
except Exception as err:
    print(err)


try :
    tool = api.NewGlobal(opt)
except Exception as err:
    print(err)

ex.SetStockType('BTC/USD.quarter')
ex.SetPeriod('M15')
ex.SetPeriodSize(3)
ex.SetIO(constant.IONONE)
key = 'BTC/USD.quarter'
ex.SetSubscribe(key, constant.CacheAccount)
ex.SetSubscribe(key, constant.CacheRecord)
ex.SetSubscribe(key, constant.CachePosition)
ex.SetSubscribe(key, constant.CacheOrder)

time_range = ex.BackGetTimeRange()

ex.SetBackTime(time_range[0], time_range[1], ex.GetPeriod())
ex.Start()

while True:
    records = ex.GetRecords()
    #if len(records) > 1:
    #    print(records[-1])
