## simple example
# Symbol from https://www.fmz.com/m/database
import fmz
import sys

# 返回值为数组，值顺序为 开 高 低 收 成交量, 每次最大1000条
# fmz.get_bars('futures_huobidm.btc.quarter', '15m', '2020-01-01', '2021-1-10')
#
def GetRecords():
    if len(sys.argv) < 5:
        print('error param!')
        return
    symbol = sys.argv[1]
    period = sys.argv[2]
    start = sys.argv[3]
    end = sys.argv[4]
    print('symbol %s period %s start %s end %s' % (symbol, period, start, end))

    records = fmz.get_bars(symbol, period, start, end)
    print(records)

GetRecords()
