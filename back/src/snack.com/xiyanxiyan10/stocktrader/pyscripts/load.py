# Simple example
# Symbol from https://www.fmz.com/m/database
# Demo python3 load.py futures_huobidm.btc.quarter 2020-01 ~/Desktop/huobidm.csv

import fmz
import sys
import time
import calendar
import csv

headers = ['time','open','high','low','close', 'vol']

def GetRecords():
    if len(sys.argv) < 4:
        print('error param!')
        return
    symbol = sys.argv[1]
    period = '5m'
    start = sys.argv[2]
    csv_file = sys.argv[3]
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

    rows = []
    for records in records_vec:
        for record in records:
            rows.append(record)

    with open(csv_file,'w')as f:
        f_csv = csv.writer(f)
        f_csv.writerow(headers)
        f_csv.writerows(rows)


    print(len(rows))

GetRecords()
