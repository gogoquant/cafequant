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



// 全局变量
var currency0 = exchanges[0].GetCurrency();
var ChartObj = null;
var perRecordsTime = 0;
var ma30 = 30;
var last_ticker_time = (new Date()).valueOf();
var h4_maline = 0
var h1_maline = 0
var m15_maline = 0
var m15LastPos = -1
var h4LastPos = -1
var h1LastPos = -1
var last_price = 0
var ticker_time = -1
var getfail_cnt = 0
var notice_cnt = 0
var noticefail_cnt = 0
var mail_pwd = MailPwd
var mail_user = MailUser
var mail_host = MailHost
//var mail_host = 'http://66.112.223.53/mail'


var m15TickerRobot = new TickerRobot("M15")
var h4TickerRobot = new TickerRobot("H4")
var h1TickerRobot = new TickerRobot("H1")
var util = new TradeUtil("util")
var msgSendRobot = new TradeRobot("msgsender")
var msgQueue = new ArrayQueue()
var queryQueue = new ArrayQueue()


function BuildMsg(user, pwd, msg){
    var msg = 'admin_user=' + user+ '&' +'admin_pwd=' + pwd + '&' + 'msg=' + msg 
    return {method:'POST', data:msg, timeout:1000}
}

function onMakeM15() {
    //M15
    records = exchange.GetRecords(PERIOD_M15)
    if (util.CheckNull(records)) {
        Log("获取15M klines 失败")
        getfail_cnt++
        return true;
    }
    //M15 K line
    $.PlotRecords(records, currency0);
    var tmp_line = util.GetMaline(records, ma30)
    if (util.CheckNull(tmp_line)) {
        Log("转换15M klines 失败")
        return true;
    }
    m15_maline = tmp_line[tmp_line.length - 1]
    return true;
}

function onMakeH4() {
    //H1
    records = exchange.GetRecords(PERIOD_H1);
    if (util.CheckNull(records)) {
        getfail_cnt++
        Log("获取1H klines 失败")
        return true;
    }
    //convert h1-->h4
    records = util.GetNewCycleRecords(records, 1000 * 60 * 60 * 4)
    var tmp_line = util.GetMaline(records, ma30)
    if (util.CheckNull(tmp_line)) {
        Log("转换4h klines 失败")
        return true;
    }
    h4_maline = tmp_line[tmp_line.length - 1]
    return true;
}

function onMakeH1() {
    //H1
    records = exchange.GetRecords(PERIOD_H1);
    if (util.CheckNull(records)) {
        getfail_cnt++
        Log("获取1H klines 失败")
        return true;
    }
    var tmp_line = util.GetMaline(records, ma30)
    if (util.CheckNull(tmp_line)) {
        Log("转换1h klines 失败")
        return true;
    }
    h1_maline = tmp_line[tmp_line.length - 1]
    return true;
}


function onTick() {
    ticker_curr = exchange.GetTicker()
    if (util.CheckNull(ticker_curr)) {
        getfail_cnt++
        Log("获取tick 失败")
        return true;
    }
    ticker_time = ticker_curr.Time
    last_price = ticker_curr.Last
    var nowTime = (new Date()).valueOf();
    

    m15TickerRobot.SetTicker(ticker_curr)
    m15TickerRobot.SetMa(m15_maline)
    m15TickerRobot.PriceCross()

    h4TickerRobot.SetTicker(ticker_curr)
    h4TickerRobot.SetMa(h4_maline)
    h4TickerRobot.PriceCross()

    h1TickerRobot.SetTicker(ticker_curr)
    h1TickerRobot.SetMa(h1_maline)
    h1TickerRobot.PriceCross()

    var table = {
        type: 'table',
        title: '综合信息',
        cols: ['键', '值'],
        rows: []
    }

    table.rows.push(["15分钟ma30", m15_maline.toString()])
    table.rows.push(["4小时ma30" , h4_maline.toString()])
    table.rows.push(["1小时ma30" , h1_maline.toString()])
    table.rows.push(["现价", last_price.toString()])
    table.rows.push(["交易所异常返回次数", getfail_cnt.toString()])
    table.rows.push(["MA30-M15调试数据", JSON.stringify(m15TickerRobot.Print())])
    table.rows.push(["MA30-H4调试数据", JSON.stringify(h4TickerRobot.Print())])
    table.rows.push(["MA30-H1调试数据", JSON.stringify(h1TickerRobot.Print())])
    LogStatus('`' + JSON.stringify(table) + '`\n')

    var m15CurrPos = m15TickerRobot.cross
    var h4CurrPos = h4TickerRobot.cross
    var h1CurrPos = h1TickerRobot.cross

    $.PlotLine('ticker-price', ticker_curr.Last, ticker_curr.Time);
    if (M15) {
        $.PlotLine('ma30-m15', m15_maline, ticker_curr.Time);
    }
    if (H4) {
        $.PlotLine('ma30-h4', h4_maline, ticker_curr.Time);
    }
    if(H1){
        $.PlotLine('ma30-h1', h1_maline, ticker_curr.Time);   
    }

    //反转时告警   
    if (m15CurrPos != m15LastPos && m15LastPos != -1 && M15) {
        if (m15LastPos == 0) {
            m15TickerRobot.dir = "up"
        } else {
            m15TickerRobot.dir = "down"
        }
        //Log("Ticker M15 Cross:" + JSON.stringify(m15TickerRobot))
        var logMsg = "ticker cross ma30m15:" + JSON.stringify(m15TickerRobot.Print())
        Log(logMsg)
        msgQueue.clear()
        if(msgSendRobot.CheckInterval()){
            msgQueue.push(logMsg)
            onNotice()
        }
    } else {
        m15TickerRobot.dir = "keep"
    }


    if (h4CurrPos != h4LastPos && h4LastPos != -1 && H4) {
        if (h4LastPos == 0) {
            h4TickerRobot.dir = "up"
        } else {
            h4TickerRobot.dir = "down"
        }
        var logMsg = "ticker cross ma30h4:" + JSON.stringify(h4TickerRobot.Print()) 
        Log(logMsg)
        msgQueue.clear()
        if(msgSendRobot.CheckInterval()){
            msgQueue.push(logMsg)
            onNotice()
        }
    } else {
        h4TickerRobot.dir = "keep"
    }

    if (h1CurrPos != h1LastPos && h1LastPos != -1 && H4) {
        if (h1LastPos == 0) {
            h1TickerRobot.dir = "up"
        } else {
            h1TickerRobot.dir = "down"
        }
        var logMsg = "ticker cross ma30h1:" + JSON.stringify(h1TickerRobot.Print())
        Log(logMsg)
        msgQueue.clear()
        if(msgSendRobot.CheckInterval()){
            msgQueue.push(logMsg)
            onNotice()
        }
    } else {
        h1TickerRobot.dir = "keep"
    }

    last_ticker_time = nowTime
    m15LastPos = m15CurrPos
    h4LastPos = h4CurrPos
    h1LastPos = h1CurrPos

    return true;
}

// 微信通知可以设置通知频率，即x分钟内1次通知
function onNotice() {
    if(msgQueue.size()==0){
        return true;
    }
    if (msgSendRobot.CheckInterval()) {
        curr_msg = msgQueue.getRear()
        Log(curr_msg + "@")
        notice_cnt++
        LogProfit(notice_cnt, '&') 
        send_msg = BuildMsg(mail_user, mail_pwd, curr_msg)
        Log("Send mail:", send_msg)
        HttpQuery(mail_host, send_msg)           
        msgQueue.clear()
        msgSendRobot.SetInterval(NoticeTime)
    }
    return true;
}


function main() {
    var curr_msg = "Test notice mail"
    var send_msg = BuildMsg(mail_user, mail_pwd, curr_msg)
    Log("Send mail:", send_msg)
    HttpQuery(mail_host, send_msg)           
    
    LogReset(1);
    ChartObj = Chart(null);
    ChartObj.reset();
    ChartObj = $.GetCfg();
    // 处理 指标轴------------------------
    ChartObj.yAxis = [{
            title: {
                text: 'K线'
            }, //标题
            style: {
                color: '#4572A7'
            }, //样式 
            opposite: false //生成右边Y轴
        },
        {
            title: {
                text: "指标轴"
            },
            opposite: true, //生成右边Y轴  ceshi
        }
    ];
    // 初始化指标线
    var chart = null;
    ticker_curr = null

    while(util.CheckNull(ticker_curr)) {
        ticker_curr = exchange.GetTicker()  
        Sleep(1000)
    }
    ticker_time = ticker_curr.Time

    $.PlotLine('ticker-price', 0, ticker_time);
    if (M15) {
        chart = $.PlotLine('ma30-m15', 0, ticker_time);
    }
    if (H4) {
        chart = $.PlotLine('ma30-h4', 0, ticker_time);
    }
    if (H1) {
        chart = $.PlotLine('ma30-h1', 0, ticker_time);
    }
    msgSendRobot.SetInterval(NoticeTime)
    chart.update(ChartObj);
    chart.reset();

    if(M15){
        queryQueue.push(onMakeM15);
    }

    if(H4){
        queryQueue.push(onMakeH4);
    }
    if(H1){
        queryQueue.push(onMakeH1); 
    }
    queryQueue.push(onTick); 
    //queryQueue.push(onNotice);
   while(true){
        var fn = queryQueue.pop()
        fn()
        queryQueue.push(fn)
        Sleep(Interval)
   }
}

/*
m15TickerRobot.High = 7640
m15TickerRobot.Low = 7365
m15TickerRobot.SetMa(7641)
m15TickerRobot.CheckCross()
console.log(m15TickerRobot.cross)
*/
