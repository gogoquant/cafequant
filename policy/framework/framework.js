// 优先级队列
PriorityQueue = function(Fn) {
    var items = [];

    function QueueElement(element, priority) {
        this.element = element;
        this.priority = priority;
    }

    this.push = function(element, priority) {
        var qe = new QueueElement(element, priority);
        //若队列为空则直接将元素入列否则需要比较该元素与其他元素的优先级
        if (this.isEmpty()) {
            items.push(qe);
        } else {
            var added = false;
            //找出比要添加元素的优先级值更大的项目，就把新元素插入到它之前。
            for (var i = 0; i < items.length; i++) {
                //一旦找到priority值更大的元素就插入新元素并终止队列循环
                if (Fn(qe.priority, items[i].priority)) {
                    items.splice(i, 0, qe);
                    added = true;
                    break;
                }
            }
            if (!added) {
                items.push(qe);
            }
        }
    };
    //获取队首  
    this.getFront = function() {
        return arr[0];
    }
    //获取队尾  
    this.getRear = function() {
        return items[arr.length - 1]
    }
    //出队操作  
    this.pop = function() {
        return items.shift();
    }
    //清空队列  
    this.clear = function() {
        items = [];
    }

    this.isEmpty = function() {
        return items.length === 0;
    };

    this.size = function() {
        return items.length
    };
    this.print = function() {
        var str = '';
        for (var i = 0; i < items.length; i++) {
            str += items[i].priority + ' ' + items[i].element + '\n'
        }

        console.log(str.toString());
    }
}

// 队列////////////////////////////////////////////////////////////
ArrayQueue = function() {
    var arr = [];
    //入队操作  
    this.push = function(element) {
        arr.push(element);
        return true;
    }
    //出队操作  
    this.pop = function() {
        return arr.shift();
    }
    //获取队首  
    this.getFront = function() {
        return arr[0];
    }
    //获取队尾  
    this.getRear = function() {
        return arr[arr.length - 1]
    }
    //清空队列  
    this.clear = function() {
        arr = [];
    }
    //获取队长  
    this.size = function() {
        return arr.length;
    }
    this.printOne = function() {
        if (arr.length <= 0) {
            return 0
        }
        return this.getFront()
    }

}

// ticker管理用于检查交叉/////////////////////////////////////////////////////////////////
TickerRobot = function TickerRobot(name) {
    this.name = name;
    this.High = 0
    this.Low = 0
    this.Last = 0
    this.Volume = 0
    this.ma = 0
    this.time = 0
    this.cross = 0
    this.dir = "keep"
    this.Name = function() {
        return this.name;
    }
    this.Print = function() {
        return {
            ma: this.ma,
            last: this.Last,
            time: this.time,
            dir: this.dir
        }
    }
    this.SetTicker = function(ticker_curr) {
        this.High = ticker_curr.High,
            this.Low = ticker_curr.Low,
            this.Last = ticker_curr.Last,
            this.Volume = ticker_curr.Volume,
            this.time = _D(ticker_curr.Time)
    }
    this.SetMa = function(ma) {
        this.ma = ma
    }
    this.KlineCross = function() {
        if (this.ma <= this.High && this.Low <= this.ma) {
            this.cross = 1
        } else {
            this.cross = 0
        }
        return this.cross

    }
    this.PriceCross = function() {
        if (this.Last < this.ma) {
            this.cross = 0
        } else {
            this.cross = 1
        }
        return this.cross
    }
}

///////////观察者,用于记录某点的交易变化////////////////////////////////////////////////
WatchRobot = function(name) {
    this.name = name;
    this.pos = -1
    this.depth = 0
    this.price = 0
    this.offsetdepth = 0
    this.offsetprice = 0
    this.out = 0 //记录被突破的次数
    this.in = 0 //记录未被突破的次数
    this.Name = function() {
        return this.name;
    }
    this.Debug = function() {
        console.log("name: ", this.name)
        console.log("pos: ", this.pos)
        console.log("depth: ", this.depth)
        console.log("price: ", this.price)
    }
    this.CleanWatch = function() {
        this.pos = -1
        this.depth = 0
        this.price = 0
    }

    this.CleanCnt = function() {
        this.out = 0 //记录被突破的次数
        this.in = 0 //记录未被突破的次数
    }
    this.SetWatch = function(pos, depth, price) {
        this.pos = pos
        this.depth = depth
        this.price = price
    }

    this.CheckWatch = function(depth, price) {
        this.offsetdepth = depth - this.depth
        this.offsetprice = price - this.price
    }
}


TradeRobot = function(name) {
    this.name = name;
    this.price = 0 //下单价  
    this.amount = 0 //下单量
    this.leftTime = 0 //用于显示
    this.interval = -1
    this.nextTime = -1
    this.run = 1;
    this.Debug = function() {
        console.log("name: ", this.name)
        console.log("leftTime: ", this.leftTime)
        console.log("interval: ", this.interval)
        console.log("nextTime: ", this.nextTime)
    }
    this.Name = function() {
        return this.name;
    }
    this.SetInterval = function(num) {
        this.interval = num
        this.leftTime = 0 - num
        this.nextTime = (new Date()).valueOf() + num
    }

    this.AddInterval = function(num) {
        this.leftTime = (new Date()).valueOf() - this.nextTime
        if (this.leftTime < 0) {
            this.interval += num
            this.nextTime += num
            this.leftTime -= num
            return
        }
        this.interval = num
        this.leftTime = 0 - num
        this.nextTime = (new Date()).valueOf() + num
    }

    this.CheckInterval = function() {
        this.leftTime = (new Date()).valueOf() - this.nextTime
        Log({"curr":(new Date()).valueOf(), "next":this.nextTime, "left":this.leftTime})
        return this.leftTime >= 0 ? 1 : 0
    }

    this.Left = function() {
        this.leftTime = (new Date()).valueOf() - this.nextTime
        return this.leftTime
    }

    this.Interval = function() {
        return this.interval
    }

}

// 工具包, 主要提供完全通用的计算公式等
TradeUtil = function(name) {
    this.name = name;
    this.Name = function() {
        return this.name;
    }
    this.Tick2Sec = function(tick) {
        return tick / 1000.0;
    }
    this.Tick2Min = function(tick) {
        return tick / 1000.0 / 60.0;
    }
    this.Tick2Hour = function(tick) {
        return tick / 1000.0 / 60.0 / 60.0;
    }
    this.Precision = function(num, pre) {
        var str1 = String(num);
        if (str1.indexOf('.') < 0) {
            return num
        }
        var str = str1.substring(0, str1.indexOf('.') + pre + 1);
        var num = Number(str);
        return num
    }
    this.CheckNull = function(elem) {
        if (typeof elem == undefined || elem == null) {
            return true
        }
        return false
    }

    this.GetMaline = function(records, ma) {
        if (this.CheckNull(records)) {
            return null
        }
        if (records.length < Math.abs(ma)) {
            return null
        }
        var ma_line = TA.MA(records, ma);
        return ma_line
    }
    this.DeepCopy = function(obj) {
        var result = Array.isArray(obj) ? [] : {};
        for (var key in obj) {
            if (obj.hasOwnProperty(key)) {
                if (typeof obj[key] === 'object' && obj[key] !== null) {
                    result[key] = deepCopy(obj[key]); //递归复制
                } else {
                    result[key] = obj[key];
                }
            }
        }
        return result;
    }
    this.GetNewCycleRecords = function(sourceRecords, targetCycle) { // K线合成函数
        var ret = []

        // 首先获取源K线数据的周期
        if (!sourceRecords || sourceRecords.length < 2) {
            return null
        }
        var sourceLen = sourceRecords.length
        var sourceCycle = sourceRecords[sourceLen - 1].Time - sourceRecords[sourceLen - 2].Time

        if (targetCycle % sourceCycle != 0) {
            Log("targetCycle:", targetCycle)
            Log("sourceCycle:", sourceCycle)
            throw "targetCycle is not an integral multiple of sourceCycle."
        }

        if ((1000 * 60 * 60) % targetCycle != 0 && (1000 * 60 * 60 * 24) % targetCycle != 0) {
            Log("targetCycle:", targetCycle)
            Log("sourceCycle:", sourceCycle)
            Log((1000 * 60 * 60) % targetCycle, (1000 * 60 * 60 * 24) % targetCycle)
            throw "targetCycle cannot complete the cycle."
        }

        var multiple = targetCycle / sourceCycle


        var isBegin = false
        var count = 0
        var high = 0
        var low = 0
        var open = 0
        var close = 0
        var time = 0
        var vol = 0
        for (var i = 0; i < sourceLen; i++) {
            // 获取 时区偏移数值
            var d = new Date()
            var n = d.getTimezoneOffset()

            if (((1000 * 60 * 60 * 24) - sourceRecords[i].Time % (1000 * 60 * 60 * 24) + (n * 1000 * 60)) % targetCycle == 0) {
                isBegin = true
            }

            if (isBegin) {
                if (count == 0) {
                    high = sourceRecords[i].High
                    low = sourceRecords[i].Low
                    open = sourceRecords[i].Open
                    close = sourceRecords[i].Close
                    time = sourceRecords[i].Time
                    vol = sourceRecords[i].Volume

                    count++
                } else if (count < multiple) {
                    high = Math.max(high, sourceRecords[i].High)
                    low = Math.min(low, sourceRecords[i].Low)
                    close = sourceRecords[i].Close
                    vol += sourceRecords[i].Volume

                    count++
                }

                if (count == multiple || i == sourceLen - 1) {
                    ret.push({
                        High: high,
                        Low: low,
                        Open: open,
                        Close: close,
                        Time: time,
                        Volume: vol,
                    })
                    count = 0
                }
            }
        }

        return ret
    }

}