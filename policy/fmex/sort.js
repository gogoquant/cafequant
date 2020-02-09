// author mhw
// brief fmex policy
// created 201910281001
// update  201912261141

/*
{
    "timestamp": 1577761436313,
    "period": 0,
    "args": [
        ["Skip", 8],
        ["Amount", 6000],
        ["Intervel", 100],
        ["HighBox", 7500],
        ["LowBox", 6800],
        ["CloseMax", 25000],
        ["Lever", 20]
    ]
}
*/
//exchange.IO("base", Url) //切换基地址，方便切换实盘和模拟盘，实盘地址：https://api.fmex.com

var OutMax = 1
//var SafeAmount = 20000
var CoverSafeTime = SafeTime
var SortOrderTime = 2000
//var SafeTime = 120000
var CoverTime = 100
var ProfitTime = 9000
var WatchPos = 3

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

//观察者,用于记录某点的交易变化
function WatchRobot(name) {
    this.name = name;
    this.pos = -1
    this.depth = 0
    this.price = 0
    this.offsetdepth = 0
    this.offsetprice = 0
    this.out = 0 //记录被突破的次数
    this.VolOut = 0
    this.PriceOut = 0
    this.in = 0 //记录未被突破的次数
}

WatchRobot.prototype.Name = function() {
    return this.name;
};

WatchRobot.prototype.Debug = function() {
    console.log("name: ", this.name)
    console.log("pos: ", this.pos)
    console.log("depth: ", this.depth)
    console.log("price: ", this.price)
};

WatchRobot.prototype.CleanWatch = function() {
    this.pos = -1
    this.depth = 0
    this.price = 0
};

WatchRobot.prototype.CleanCnt = function() {
    this.out = 0 //记录被突破的次数
    this.in = 0 //记录未被突破的次数
    this.PriceOut = 0
    this.VolOut = 0
};

WatchRobot.prototype.SetWatch = function(pos, depth, price) {
    this.pos = pos
    this.depth = depth
    this.price = price
};

WatchRobot.prototype.CheckWatch = function(depth, price) {
    this.offsetdepth = depth - this.depth
    this.offsetprice = price - this.price
};

TradeRobot.prototype.Debug = function() {
    console.log("name: ", this.name)
    console.log("leftTime: ", this.leftTime)
    console.log("interval: ", this.interval)
    console.log("nextTime: ", this.nextTime)
};


//交易者,用于交易配置以及延时处理等
function TradeRobot(name) {
    this.name = name;
    this.price = 0 //下单价  
    this.amount = 0 //下单量
    this.leftTime = 0 //用于显示
    this.interval = -1
    this.nextTime = (new Date()).valueOf()
    this.run = 1;
}


TradeRobot.prototype.Name = function() {
    return this.name;
};

TradeRobot.prototype.SetInterval = function(num) {
    this.interval = num
    this.leftTime = 0 - num
    this.nextTime = (new Date()).valueOf() + num
};

TradeRobot.prototype.AddInterval = function(num) {
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
};

TradeRobot.prototype.CheckInterval = function() {
    this.leftTime = (new Date()).valueOf() - this.nextTime
    return this.leftTime >= 0 ? 1 : 0
};

TradeRobot.prototype.Left = function() {
    this.leftTime = (new Date()).valueOf() - this.nextTime
    return this.leftTime
};

TradeRobot.prototype.Interval = function() {
    return this.interval
};

// 工具包, 主要提供完全通用的计算公式等
function TradeUtil(name) {
    this.name = name;
}

TradeUtil.prototype.Name = function() {
    return this.name;
};


TradeUtil.prototype.Tick2Sec = function(tick) {
    return tick / 1000.0;
};

TradeUtil.prototype.Tick2Min = function(tick) {
    return tick / 1000.0 / 60.0;
};

TradeUtil.prototype.Tick2Hour = function(tick) {
    return tick / 1000.0 / 60.0 / 60.0;
};

TradeUtil.prototype.Precision = function(num, pre) {
    var str1 = String(num);
    if (str1.indexOf('.') < 0) {
        return num
    }
    var str = str1.substring(0, str1.indexOf('.') + pre + 1);
    var num = Number(str);
    return num
}

TradeUtil.prototype.Sleep = function(d) {
    for (var t = (new Date()).valueOf();
        (new Date()).valueOf() - t <= d;);
}

var ordersInfo = {
    buyId: 0,
    buyPrice: 0,
    sellId: 0,
    sellPrice: 0,
    closeCnt: 0,
    profile: 0.0,
    avgprofile: 0.0,
    profileTime: 1,
    buyVol: 0.0,
    sellVol: 0.0
}


var globalInfo = {
    ticker: null,
    depth: null,
    position: null,
    run: 1
}

var depthInfo = {
    asks: [],
    bids: []
}

var buyRobot = new TradeRobot("buyRobot")
var sellRobot = new TradeRobot("sellRobot")
var sortOrderRobot = new TradeRobot("sortOrderRobot")
var coverRobot = new TradeRobot("coverRobot")
var profileRobot = new TradeRobot("profileRobot")
var resetRobot = new TradeRobot("resetRobot")
var util = new TradeUtil("util")
buyRobot.SetInterval(SafeTime)
sellRobot.SetInterval(SafeTime)
coverRobot.SetInterval(CoverTime)
profileRobot.SetInterval(ProfitTime)
resetRobot.SetInterval(10 * 50 * 1000)
sortOrderRobot.SetInterval(SortOrderTime)
var queryQueue = new ArrayQueue()
var buyWatchRobot = new WatchRobot("buy")
var sellWatchRobot = new WatchRobot("sell")
var lastProfile = 0.0

function Status2Usr(){
    return {
           buyId: ordersInfo.buyId,
           buyPrice: ordersInfo.buyPrice,
           buyVol: ordersInfo.buyVol,
           sellId: ordersInfo.sellId,
           sellPrice: ordersInfo.sellPrice,
           sellVol: ordersInfo.sellVol,
           closeCnt: ordersInfo.closeCnt,
           profit: ordersInfo.avgprofile
    }
}

// 检查安全哨兵
function checkWatch(price, orders, pos) {
    sellPrice = 0
    buyPrice = 0
    if (sellWatchRobot.pos != -1 && sellWatchRobot.price < price) {
        sellPrice = price
    }

    if (buyWatchRobot.pos != -1 && buyWatchRobot.price > price) {
        buyPrice = price
    }
    return [buyPrice, sellPrice]
}

// @brief 精度调整函数小数点后0位
// @param num 待调整值
function adjustNum0(num) { //控制价格精度
    return num >= _N(num, 0) + 0.5 ? _N(num, 0) + 0.5 : _N(num, 0)
}

// @brief 评分函数 
// @param depth 深度
// @param amount 待挂单量, 目前使用外部传入的固定值
// @param currAmount 自己在该区间的挂单量
// @param fact 该区间的奖励系数
function getRatio(depth, amount, currAmount, fact) {
    ratio = _N(10000 * amount * fact / Math.max(depth + amount - currAmount, amount), 5)
}

// cut asks
function cutVec(vec, pos, len) {
    newAsks = []
    newAsks.push(arr.slice(pos, pos + len))
    return newAsks
}

// check depth rate
function checkVol(depth, price, type) {
    buyDepth = 0.0
    sellDepth = 0.0
    if(type == 0){
        for(var i = 0;  i < depth.Bids.length && depth.Bids[i].Price > price; i++){
            buyDepth += depth.Bids[i].Amount
        }
        return buyDepth
    }else{
        for(var i = 0;  i < depth.Asks.length && depth.Asks[i].Price < price; i++){
            sellDepth += depth.Asks[i].Amount
        }
        return sellDepth
    }
}

// 获取当前最佳深度
function calcDepth(depth) { //计算最佳挂单位置
    depthInfo = {
        asks: [],
        bids: []
    }
    var max_asks = Skip
    var max_bids = Skip
    var ask_price = depth.Asks[0].Price
    var bid_price = depth.Bids[0].Price
    var ask_tot_amount = 0.0
    var bid_tot_amount = 0.0

    for (var i = 0; i < 11; i++) {
        var fact = i == 0 ? 1 / 4 : 1 / 40 //官方的挖矿系数，可按需设置，如离盘口越远越大，减少成交风险
        var my_ask_amount = depth.Asks[i].Price == ordersInfo.sellPrice ? Amount : 0 //排除掉自己订单的干扰
        var my_bid_amount = depth.Bids[i].Price == ordersInfo.buyPrice ? Amount : 0
        while (ask_price <= depth.Asks[i].Price) { //考虑到未被占用的深度位置
            var ask_amount = ask_price == depth.Asks[i].Price ? depth.Asks[i].Amount : 0
            // 计算评分
            var ratio = _N(10000 * Amount * fact / Math.max(ask_amount + Amount - my_ask_amount, Amount), 5) //避免因深度更新延时导致除0
            depthInfo.asks.push(['sell_' + (i + 1), ask_price, ask_amount, ratio, 0, i == WatchPos ? 1 : 0])

            if (i >= Skip && ask_tot_amount < SafeAmount) {
                if (ratio >= depthInfo.asks[max_asks][3]) {
                    max_asks = depthInfo.asks.length - 1
                } //大于等于保证相同挖矿效率下远离盘口的优先
            }
            ask_price += 0.5
        }
        while (bid_price >= depth.Bids[i].Price) {
            var bid_amount = bid_price == depth.Bids[i].Price ? depth.Bids[i].Amount : 0
            var ratio = _N(10000 * Amount * fact / Math.max(bid_amount + Amount - my_bid_amount, Amount), 5)
            depthInfo.bids.push(['buy_' + (i + 1), bid_price, bid_amount, ratio, 0, i == WatchPos ? 1 : 0])
            if (i >= Skip && bid_tot_amount < SafeAmount) {
                if (ratio >= depthInfo.bids[max_bids][3]) {
                    max_bids = depthInfo.bids.length - 1
                }
            }
            bid_price -= 0.5
        }
        ask_tot_amount += depth.Asks[i].Amount
        bid_tot_amount += depth.Bids[i].Amount
    }
    depthInfo.asks[max_asks][4] = sellRobot.amount
    depthInfo.bids[max_bids][4] = buyRobot.amount
    return [depthInfo.bids[max_bids][1], depthInfo.asks[max_asks][1]]
}

function showTable() {
    var table = {
        type: 'table',
        title: '挂单信息',
        cols: ['位置', '价格', '数量', '额度占比（万分之一)', '挂单数量', '哨兵位置'],
        rows: []
    }
    for (var i = 0; i < depthInfo.asks.length; i++) {
        table.rows.push(depthInfo.asks[i])
    }
    for (var i = 0; i < depthInfo.bids.length; i++) {
        table.rows.push(depthInfo.bids[i])
    }
    //计算平均收益
    ordersInfo.avgprofile = util.Precision(ordersInfo.profile / ordersInfo.profileTime, 2)
    LogStatus('`' + JSON.stringify(table) + '`\n' + JSON.stringify(Status2Usr()))
}

function reset(orderType) { //重置策略，防止一些订单卡住，可能会影响其它正在运行的策略
    var orders = exchange.GetOrders()
    if (orders) {
        // close all order when orderType == -1
        if (orderType == -1) {
            for (var i = 0; i < orders.length; i++) {
                exchange.CancelOrder(orders[i].Id)
                Sleep(Intervel)
            }
            ordersInfo.buyId = 0
            ordersInfo.buyPrice = 0
            ordersInfo.sellId = 0
            ordersInfo.sellPrice = 0
            return
        }

        // close just one type order
        for (var i = 0; i < orders.length; i++) {
            if (orders[i].Type == orderType) {
                exchange.CancelOrder(orders[i].Id)
                Sleep(Intervel)
            }
        }
        if (orderType == 0) {
            ordersInfo.buyId = 0
            ordersInfo.buyPrice = 0
        }

        if (orderType == 1) {
            ordersInfo.sellId = 0
            ordersInfo.sellPrice = 0
        }
    }
}

function onexit() { //退出后撤销订单平仓
    reset(-1)
    var pos = exchange.GetPosition()
    if (pos) {
        cover(pos)
    }
}

// 平仓
function cover(pos) {
    var closeId = 0
    if (pos.length <= 0) {
        return
    }
    var leftAmount = pos[0].Amount
    if (pos[0].Type == 0) { //平多仓，采用盘口吃单，会损失手续费，可改为盘口挂单，会增加持仓风险。
        Log("平多仓", leftAmount)
        exchange.SetDirection('sell')
        closeId = exchange.Sell(-1, leftAmount, '平多仓')
    } else {
        Log("平空仓", leftAmount)
        exchange.SetDirection('buy')
        closeId = exchange.Buy(-1, leftAmount, '平空仓')
    }
    Sleep(Intervel)
}

//安全函数watchPrice
function watchPrice(price) {
    if (typeof price == "undefined") {
        return true
    }
    if (price > HighBox) {
        return true
    }
    if (price < LowBox) {
        return true
    }
    return false
}

function tradeBuy(depth, price) {
    if (price != ordersInfo.buyPrice) {
        var cancelId = ordersInfo.buyId
        if (cancelId) {
            if (exchange.CancelOrder(cancelId) != true) {
                Log("关闭买订单失败，放弃下单")
                return
            }
            Sleep(Intervel)
        }
        exchange.SetDirection('buy')
        var buyId = exchange.Buy(price, buyRobot.amount, '更新下买单')
        Sleep(Intervel)
        if (buyId) {
            ordersInfo.buyId = buyId
            ordersInfo.buyPrice = price
        } else {
            ordersInfo.buyId = 0
            ordersInfo.buyPrice = -1
        }

    }
    //Log("tradeBuy Success")  
}

function tradeSell(depth, price) {
    if (price != ordersInfo.sellPrice) {
        var cancelId = ordersInfo.sellId
        if (cancelId) {
            if (exchange.CancelOrder(cancelId) != true) {
                Log("关闭卖订单失败，放弃下单")
                return
            }
            Sleep(Intervel)
        }
        exchange.SetDirection('sell')
        var sellId = exchange.Sell(price, sellRobot.amount, '更新下卖单')
        Sleep(Intervel)
        //Log("下卖单")
        if (sellId) {
            ordersInfo.sellId = sellId
            ordersInfo.sellPrice = price
        } else {
            ordersInfo.sellId = 0
            ordersInfo.sellPrice = -1
        }
    }
    //Log("tradeSell Success")    
}


function watchProcess(depth, nextprice) {
    ticker_curr = globalInfo.ticker
    last_price = ticker_curr.Last
    var buyVol = -1
    var sellVol = -1
    if(ordersInfo.buyPrice > 0){      
        buyVol = checkVol(depth, ordersInfo.buyPrice, 0)
        ordersInfo.buyVol = buyVol
    }else{
        buyVol = checkVol(depth, nextprice[0], 0)
        ordersInfo.buyVol = buyVol
    }

    if(ordersInfo.sellPrice > 0){
        sellVol = checkVol(depth, ordersInfo.sellPrice, 1)
        ordersInfo.sellVol = sellVol
    }else{
        sellVol = checkVol(depth, nextprice[1], 1)
        ordersInfo.sellVol = sellVol
    }
    
    price = checkWatch(last_price, null, WatchPos)
 
    //买哨兵
    if (price[0] != 0 || buyVol < SafeAmount) {
        reset(0)
        buyWatchRobot.out++
        buyRobot.run = 0
        buyRobot.SetInterval(Intervel)
        if(price[0] != 0){
            buyWatchRobot.PriceOut += 1
            Log("多价格触发:", price[0])
        }
        if(buyVol < SafeAmount){
            buyWatchRobot.VolOut += 1
            Log("多深度触发:", buyVol)
        }
            //Log("多哨兵:", buyWatchRobot.out)
    }


    //卖哨兵
    if (price[1] != 0 || sellVol < SafeAmount) {
        reset(1)
        sellWatchRobot.out++
        sellRobot.run = 0
        sellRobot.SetInterval(Intervel)
        if(price[0] != 0){
            sellWatchRobot.PriceOut += 1
            Log("空价格触发:", price[1])
        }
        if(sellVol < SafeAmount){
            sellWatchRobot.VolOut += 1
            Log("空深度触发:", sellVol)
        }
            //Log("空哨兵:", sellWatchRobot.out)
    }

    //重新计算订单分配
    if (buyRobot.run + sellRobot.run == 0) {
        buyRobot.amount = 0
        sellRobot.amount = 0
    } else {
       //挂单量再平衡
        buyRobot.amount = Amount * (buyRobot.run /(buyRobot.run + sellRobot.run))
        sellRobot.amount = Amount * (sellRobot.run  /(buyRobot.run + sellRobot.run))
    }

    if (price[0] == 0) {
        buyWatchRobot.in++
    }
    if (price[1] == 0) {
        sellWatchRobot.in++
    }

    //是否到达检查周期
    if (sellRobot.CheckInterval()) {
        sellWatchRobot.SetWatch(WatchPos, depth.Asks[WatchPos].Amount, depth.Asks[WatchPos].Price)
        var outTot = sellWatchRobot.out
        sellWatchRobot.CleanCnt()
        Log("卖机器人尝试重启")
        if (sellRobot.run == 1) {
            Log("卖机器人运行中不用重启")
        } else if (outTot >= OutMax) {
            Log("卖机器重启失败:", outTot)
        } else {
            Log("卖机器人重启成功:", outTot)
            sellRobot.run = 1
        }
        sellRobot.SetInterval(SafeTime)
    } else {
        //Log("卖机器人还剩延时: ", sellRobot.Left())
    }

    //是否到达检查周期
    if (buyRobot.CheckInterval()) {
        buyWatchRobot.SetWatch(WatchPos, depth.Bids[WatchPos].Amount, depth.Bids[WatchPos].Price)
        var outTot = buyWatchRobot.out
        buyWatchRobot.CleanCnt()
        Log("买机器人尝试重启")
        if (buyRobot.run == 1) {
            Log("买机器人运行中不用重启")
        } else if (outTot >= OutMax) {
            Log("买机器重启失败:", outTot)
        } else {
            Log("买机器人重启成功:", outTot)
            buyRobot.run = 1
        }
        buyRobot.SetInterval(SafeTime)
    } else {
        //Log("买机器人还剩延时: ", buyRobot.Left())
    }

}

// trade Api 
function trade(depth, price) {
    var buyPrice = price[0]
    var sellPrice = price[1] 
    var buyVol = checkVol(depth, buyPrice, 0)
    var sellVol = checkVol(depth, sellPrice, 1)
    
    if (buyRobot.run && buyRobot.amount != 0 && buyVol > SafeAmount) {
        tradeBuy(depth, buyPrice)
    }
    if (sellRobot.run && sellRobot.amount != 0 && sellVol > SafeAmount) {
        tradeSell(depth, sellPrice)
    }
}

// debug trade
function debug_trade(price) {
    var sellPrice = price[0]
    var buyPrice = price[1]
    if (buyPrice != ordersInfo.buyPrice) {
        ordersInfo.buyPrice = buyPrice
    }
    if (sellPrice != ordersInfo.sellPrice) {
        ordersInfo.sellPrice = sellPrice
    }
}

function onDepth() {
    globalInfo.depth = _C(exchange.GetDepth)
}

function onPosition() {
    globalInfo.position = exchange.GetPosition()
}

function onTicker() {
    globalInfo.ticker = exchange.GetTicker()
}

function onProcess() {
    if (globalInfo.depth == null || globalInfo.position == null || globalInfo.ticker == null) {
        return
    }
    var depth = globalInfo.depth
    var price = calcDepth(depth)
    
    // 箱体,安全检测
    if (watchPrice(price[0]) || watchPrice(price[1])) {
        Log("异常价格:" + JSON.stringify({
            buyPrice: price[0],
            sellPrice: price[1],
            HighBox: HighBox,
            LowBox: LowBox
        }))
        globalInfo.run = 0
        return
    }

    var onTime = coverRobot.CheckInterval()
    if (onTime) {
        var pos = globalInfo.position
        if (pos && pos.length > 0 ? 1 : 0) {
            ordersInfo.closeCnt = ordersInfo.closeCnt + pos[0].Amount
            reset(-1)
            cover(pos)

            if (pos[0].Type == 0) {
                Log("多平仓后关闭买机器人")
                buyWatchRobot.out += 1
                buyRobot.run = 0
                buyRobot.SetInterval(CoverSafeTime)
            } else {
                Log("空平仓后关闭卖机器人")
                sellWatchRobot.out += 1
                sellRobot.run = 0
                //被击穿需要超长睡眠
                sellRobot.SetInterval(CoverSafeTime)
            }
        }
        coverRobot.SetInterval(coverRobot.Interval())
    }

    var onTime = sortOrderRobot.CheckInterval()
    if (onTime) {
        trade(depth, price)
        sortOrderRobot.SetInterval(sortOrderRobot.Interval())
    }

    // 哨兵检测
    watchProcess(depth, price)

    // 止损
    if (ordersInfo.closeCnt >= CloseMax) {
        Log("平仓次数过多，停止策略")
        globalInfo.run = 0
        return
    }

    var onTime = profileRobot.CheckInterval()
    // 收益展示
    if (onTime) {
        //Log("买观察者", JSON.stringify(buyWatchRobot))
        //Log("卖观察者", JSON.stringify(sellWatchRobot))
        //Log("买机器人", JSON.stringify(buyRobot))
        //Log("卖机器人", JSON.stringify(sellRobot))
        //Log("买哨兵", JSON.stringify(buyWatchRobot))
        //Log("卖哨兵", JSON.stringify(sellWatchRobot))
        //Log("挂单机器人", JSON.stringify(sortOrderRobot))
        //Log("平仓机器人", JSON.stringify(coverRobot))
        //Log("收益机器人", JSON.stringify(profileRobot))
        //Log("重置机器人", JSON.stringify(resetRobot))
        LogProfit(ordersInfo.avgprofile, '&')
        lastProfile = ordersInfo.profile
        profileRobot.SetInterval(profileRobot.Interval())
    }

    // 防止订单卡住
    var onTime = resetRobot.CheckInterval()
    if (onTime) {
        Log("重置订单")
        reset(-1)
        resetRobot.SetInterval(resetRobot.Interval())
    }

    showTable()
    ordersInfo.profile += buyRobot.amount * 1.0
    ordersInfo.profile += sellRobot.amount * 1.0
    ordersInfo.profileTime += 1
}

function main() {
    exchange.SetContractType('swap')
    exchange.SetMarginLevel(Lever)

    //检查到合适位置再启动
    buyRobot.run = 0
    sellRobot.run = 0
    buyRobot.amount = 0
    sellRobot.amount = 0
    var depth = _C(exchange.GetDepth)
    buyWatchRobot.SetWatch(WatchPos, depth.Bids[WatchPos].Amount, depth.Bids[WatchPos].Price)
    sellWatchRobot.SetWatch(WatchPos, depth.Asks[WatchPos].Amount, depth.Asks[WatchPos].Price)

    // register callback
    queryQueue.push(onDepth);
    queryQueue.push(onPosition);
    queryQueue.push(onTicker);
    queryQueue.push(onProcess);

    reset(-1)
    while (true) {
        var fn = queryQueue.pop()
        fn()
        queryQueue.push(fn)
        if (globalInfo.run == 0) {
            break
        }
        Sleep(Intervel)
    }

    // 平仓清理
    while (true) {
        Log("异常情况尝试平仓")
        // 定时
        var pos = exchange.GetPosition()
        if (typeof pos == "undefined" || pos == null) {
            continue
        }

        if (pos.length <= 0) {
            Log("异常情况清理完毕程序退出")
            break
        }
        //尝试平仓
        cover(pos)

        Sleep(Intervel)
    }
}

/*
function _N(num, pre){
    var str1=String(num);           //将类型转化为字符串类型
    if(str1.indexOf('.') <0){
        return num
    }
    var str=str1.substring(0,str1.indexOf('.')+pre+1);  //截取字符串
    var num=Number(str);            //转化为number类型
    return num
}
*/
