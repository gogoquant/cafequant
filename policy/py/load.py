# Simple example
# Symbol from https://www.fmz.com/m/database
# Demo python3 load.py futures_huobidm.btc.quarter 2020-01 2 ~/Desktop/huobidm.csv

import fmz
import sys
import time
import calendar
import csv

headers = ['time','open','high','low','close', 'vol']

def main():
    if len(sys.argv) < 4:
        print('error param!')
        return

    symbol = sys.argv[1]
    period = '15m'
    start = sys.argv[2]
    tot = sys.argv[3]
    csv_file = sys.argv[4]

    timeArray = time.strptime(start, "%Y-%m")
    tot = int(tot)

    rows = []
    for i in range(0,tot):
        month = timeArray.tm_mon + i
        year = timeArray.tm_year
        if month > 12:
            month = month - 12
            year = year + 1

        print("start load month %s %s" % (str(year), str(month)))
        start = str(year).zfill(4) + '-' + str(month).zfill(2)
        print("start load month %s" % start)
        records_vec = getRecords(symbol, period, start, csv_file)
        if len(records_vec) == 0:
            print('month %s is empty!' % start)
            break

        for records in records_vec:
            for record in records:
                rows.append(record)

    with open(csv_file,'w')as f:
        f_csv = csv.writer(f)
        f_csv.writerow(headers)
        f_csv.writerows(rows)

    print(len(rows))

def getRecords(symbol, period, start, csv_file):
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

    return records_vec

main()
