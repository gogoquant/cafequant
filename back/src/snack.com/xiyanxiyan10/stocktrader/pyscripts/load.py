## simple example
# Symbol from https://www.fmz.com/m/database
import fmz
import sys
import time
import calendar
# fmz.get_bars('futures_huobidm.btc.quarter', 15m, '2020-01-01', '2021-1-10')
#
def GetRecords():
    if len(sys.argv) < 3:
        print('error param!')
        return
    symbol = sys.argv[1]
    period = '5m'
    start = sys.argv[2]
    begin = start
    timeArray = time.strptime(start, "%Y-%m")
    print('symbol %s period %s start %s' % (symbol, period, start))
    print(timeArray)
    monthRange = calendar.monthrange(timeArray.tm_year, timeArray.tm_mon)
    print(monthRange)
    records_vec = []
    for i in range(monthRange[1]):
        start = begin + '-' + str(i+1).zfill(2)
        # next month
        if i == monthRange[1] - 1:
               if timeArray.tm_mon < 12:
                    end = str(timeArray.tm_year).zfill(4) + '-' + str(timeArray.tm_mon + 1).zfill(2) + "-" + '01'
               else:
                    end = str(timeArray.tm_year+1).zfill(4) + "-01-01"
        else:
            end = begin + '-' + str(i+2).zfill(2)
        print('load symbol %s period %s start %s end %s' % (symbol, period, start, end))
        records = fmz.get_bars(symbol, period, start, end)
        records_vec.append(records)
        time.sleep(0.1)

    print(len(records_vec))

GetRecords()
