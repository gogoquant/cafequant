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

 872  - 0fmex/sort_quantbot.js
@ -0,0 +1,872 @@
// This is an example algorithm

var Skip = 8
var OutMax = 1
var SafeTime = 60000
var Amount = 400
var SafeAmount = 40000
var CoverSafeTime = SafeTime
var SortOrderTime = 2000
var CoverTime = 100
var ProfitTime = 9000
var WatchPos = 3
var Intervel = 100
var HighBox = 8500
var LowBox = 7500
var CloseMax = 200
var Lever = 1

function _N(num, pre){
    var str1=String(num);           //将类型转化为字符串类型
    if(str1.indexOf('.') <0){
        return num
    }
    var str=str1.substring(0,str1.indexOf('.')+pre+1);  //截取字符串
    var num=Number(str);            //转化为number类型
    return num
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

TradeRobot.prototype.Debug = function() {
    console.log("name: ", this.name)
    console.log("leftTime: ", this.leftTime)
    console.log("interval: ", this.interval)
    console.log("nextTime: ", this.nextTime)
};


TradeRobot.prototype.Name = function() {
    return this.name;
};

TradeRobot.prototype.SetInterval = function(num) {
    this.interval = num
    this.leftTime = 0 - num
    this.nextTime = (new Date()).valueOf() + num
};

TradeRobot.prototype.TriggerInterval = function() {
    this.nextTime = (new Date()).valueOf()
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
    account: null,
    orders: null,
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
buyRobot.run = 0
sellRobot.run = 0

function Status2Usr(){
    ordersInfo.avgprofile = util.Precision(ordersInfo.profile / ordersInfo.profileTime, 2)
    return {
           buyId: ordersInfo.buyId,
           buyPrice: ordersInfo.buyPrice,
           buyVol: ordersInfo.buyVol,
           buyWatch: buyWatchRobot.price,
           sellId: ordersInfo.sellId,
           sellPrice: ordersInfo.sellPrice,
           sellVol: ordersInfo.sellVol,
           sellWatch: sellWatchRobot.price,
           closeCnt: ordersInfo.closeCnt,
           profit: ordersInfo.avgprofile
    }
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
            if (E.CancelOrder(cancelId) != true) {
                ordersInfo.buyId = 0
                ordersInfo.buyPrice = -1
                G.Log("关闭买订单失败，放弃下单")
                return
            }
            G.Sleep(Intervel)
        }
        E.SetDirection('buy')
        var buyId = E.Buy(price, buyRobot.amount, '更新下买单')
        G.Sleep(Intervel)
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
            if (E.CancelOrder(cancelId) != true) {
                ordersInfo.sellId = 0
                ordersInfo.sellPrice = -1
                G.Log("关闭卖订单失败，放弃下单")
                return
            }
            G.Sleep(Intervel)
        }
        E.SetDirection('sell')
        var sellId = E.Sell(price, sellRobot.amount, '更新下卖单')
        G.Sleep(Intervel)
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
    //G.Log("ticker_curr:", ticker_curr)
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
        reset("buy")
        buyWatchRobot.out++
        buyRobot.run = 0
        buyRobot.TriggerInterval()
        if(price[0] != 0){
            buyWatchRobot.PriceOut += 1
            G.Log("多价格触发:", price[0], "->", buyWatchRobot.price)
        }
        if(buyVol < SafeAmount){
            buyWatchRobot.VolOut += 1
            G.Log("多深度触发:", buyVol, "->", SafeAmount)
        }
            //Log("多哨兵:", buyWatchRobot.out)
    }

    //卖哨兵
    if (price[1] != 0 || sellVol < SafeAmount) {
        reset("sell")
        sellWatchRobot.out++
        sellRobot.TriggerInterval()
        sellRobot.run = 0
        if(price[1] != 0){
            sellWatchRobot.PriceOut += 1
            G.Log("空价格触发:", price[1], "->", sellWatchRobot.price)
        }
        if(sellVol < SafeAmount){
            sellWatchRobot.VolOut += 1
            G.Log("空深度触发:", sellVol, "->", SafeAmount)
        }
            //Log("空哨兵:", sellWatchRobot.out)
    }

    oldBuyAmount = buyRobot.amount
    oldSellAmount = sellRobot.amount
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

    //当重置后订单量平衡发生便宜则需要重置订单
    if(oldBuyAmount != buyRobot.amount || oldSellAmount != sellRobot.amount){
        reset("")
        G.Log("重置订单量平衡 buy:", buyRobot.amount, " sell:", sellRobot.amount)
    }

    //是否到达检查周期
    if (sellRobot.CheckInterval()) {
        sellWatchRobot.SetWatch(WatchPos, depth.Asks[WatchPos].Amount, depth.Asks[WatchPos].Price)
        G.Log("设置空哨兵价格:", depth.Asks[WatchPos].Price)
        var outTot = sellWatchRobot.out
        sellWatchRobot.CleanCnt()
        G.Log("卖机器人尝试重启")
        if (sellRobot.run == 1) {
            G.Log("卖机器人运行中不用重启")
        } else if (outTot >= OutMax) {
            G.Log("卖机器重启失败:", outTot)
        } else {
            G.Log("卖机器人重启成功:", outTot)
            sellRobot.run = 1
        }
        sellRobot.SetInterval(SafeTime)
    } else {
        //Log("卖机器人还剩延时: ", sellRobot.Left())
    }

    //是否到达检查周期
    if (buyRobot.CheckInterval()) {
        buyWatchRobot.SetWatch(WatchPos, depth.Bids[WatchPos].Amount, depth.Bids[WatchPos].Price)
        G.Log("设置多哨兵价格:", depth.Bids[WatchPos].Price)
        var outTot = buyWatchRobot.out
        buyWatchRobot.CleanCnt()
        G.Log("买机器人尝试重启")
        if (buyRobot.run == 1) {
            G.Log("买机器人运行中不用重启")
        } else if (outTot >= OutMax) {
            G.Log("买机器重启失败:", outTot)
        } else {
            G.Log("买机器人重启成功:", outTot)
            buyRobot.run = 1
        }
        buyRobot.SetInterval(SafeTime)
    } else {
        //G.Log("买机器人还剩延时: ", buyRobot.Left())
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

function onDepth() {
    globalInfo.depth = E.GetDepth(11)
}

function onAccount() {
    globalInfo.account = E.GetAccount()
}

function onOrders() {
    globalInfo.orders = E.GetOrders()
}

function onTicker(){
    globalInfo.ticker = E.GetTicker()
}

function onPosition(){
    globalInfo.position = E.GetPosition()
}

function reset(orderType) { 
    var orders = globalInfo.orders
    if (orders) {
        // close all order when orderType == -1
        if (orderType == "") {
            for (var i = 0; i < orders.length; i++) {
                E.CancelOrder(orders[i].Id)
                G.Sleep(Intervel);
            }
            ordersInfo.buyId = 0
            ordersInfo.buyPrice = 0
            ordersInfo.sellId = 0
            ordersInfo.sellPrice = 0
            return
        }

        // close just one type order
        for (var i = 0; i < orders.length; i++) {
            if (orders[i].TradeType == orderType) {
                E.CancelOrder(orders[i].Id)
                G.Sleep(Intervel);
            }
        }
        if (orderType == "buy") {
            ordersInfo.buyId = 0
            ordersInfo.buyPrice = 0
        }

        if (orderType == "sell") {
            ordersInfo.sellId = 0
            ordersInfo.sellPrice = 0
        }
    }
}

// 平仓
function cover(pos) {
    var closeId = 0
    if (pos.length <= 0) {
        return
    }
    var leftAmount = pos[0].Amount
    if (pos[0].TradeType == "buy") { //平多仓，采用盘口吃单，会损失手续费，可改为盘口挂单，会增加持仓风险。
        G.Log("平多仓", leftAmount)
        E.SetDirection('sell')
        closeId = E.Sell(-1, leftAmount, '平多仓')
    } else {
        G.Log("平空仓", leftAmount)
        E.SetDirection('buy')
        closeId = E.Buy(-1, leftAmount, '平空仓')
    }
    G.Sleep(Intervel)
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

function VaildInfo(){
    if(typeof globalInfo.depth == undefined || globalInfo.depth === false || globalInfo.depth === null){
        G.Log("depth get fail:", globalInfo.depth)
        return false
    }
    /*
    if(typeof globalInfo.account == undefined || globalInfo.account === false || globalInfo.account === null){
        G.Log("account get fail:", globalInfo.account)
        return false
    }
    */

    if(typeof globalInfo.ticker == undefined || globalInfo.ticker === false || globalInfo.ticker === null){
        G.Log("ticker get fail:", globalInfo.ticker)
        return false
    }

    if(typeof globalInfo.position == undefined || globalInfo.position === false || globalInfo.position === null){
        G.Log("position get fail:", globalInfo.position)
        return false
    }
    return true
}

function onExit(){
   if(VaildInfo() === false){
        return 
    }
   
    if(globalInfo.position.length == 0 && globalInfo.orders.length == 0){
        globalInfo.run = 0
        return 
    }
    //clean 
    reset("")
    cover(globalInfo.position)
    G.Log("exit process running")
    G.Sleep(Intervel)
    return       
}

function PreExit() { //退出后撤销订单平仓
    //G.MailStop()

    queryQueue.clear()

    queryQueue.push(onDepth); 
    //queryQueue.push(onAccount);
    queryQueue.push(onTicker);
    queryQueue.push(onOrders);
    queryQueue.push(onPosition); 

    queryQueue.push(onExit);
    G.Log("enter into exit process") 
    return 
}

function exit() { //退出后撤销订单平仓
    G.Log("主动退出策略") 
    PreExit(); 
     
    while (true) {
        var fn = queryQueue.pop()
        fn()
        queryQueue.push(fn)
        if (globalInfo.run == 0) {
            break
        }
        G.Sleep(Intervel);
    }
    G.Log("exit process success")
    return 
}

function onProcess() {
    if(VaildInfo() === false){
        return 
    }
    //G.Log("orders:", globalInfo.orders)
    //return 
    var depth = globalInfo.depth
    var price = calcDepth(depth)
    
    // 箱体,安全检测
    if (watchPrice(price[0]) || watchPrice(price[1])) {
        G.Log("异常价格:" + JSON.stringify({
            buyPrice: price[0],
            sellPrice: price[1],
            HighBox: HighBox,
            LowBox: LowBox
        }))
        G.Log("异常价格退出策略")
        PreExit()
        return
    }

    if (coverRobot.CheckInterval()) {
        var pos = globalInfo.position
        if (pos && pos.length > 0 ? 1 : 0) {
            ordersInfo.closeCnt = ordersInfo.closeCnt + pos[0].Amount
            if (pos[0].Type == "buy") {
                G.Log("进入多平仓关闭机器人", pos[0].Amount)
            } else {
                G.Log("进入空平仓关闭机器人", pos[0].Amount)
            }
            PreExit()
            return
        }
        coverRobot.SetInterval(coverRobot.Interval())
    }

    if (sortOrderRobot.CheckInterval()) {
        trade(depth, price)
        sortOrderRobot.SetInterval(sortOrderRobot.Interval())
    }

    // 哨兵检测
    watchProcess(depth, price)

    // 止损
    if (ordersInfo.closeCnt >= CloseMax) {
        G.Log("平仓次数过多退出策略")
        //globalInfo.run = 0
        PreExit()
        return
    }

    // 收益展示
    if (profileRobot.CheckInterval()) {
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
        G.Log("status:", Status2Usr())
        lastProfile = ordersInfo.profile
        profileRobot.SetInterval(profileRobot.Interval())
    }

    if (resetRobot.CheckInterval()) {
        G.Log("重置订单")
        reset("")
        resetRobot.SetInterval(resetRobot.Interval())
    }

    //showTable()
    ordersInfo.profile += buyRobot.amount * 1.0
    ordersInfo.profile += sellRobot.amount * 1.0
    ordersInfo.profileTime += 1
}

function main() {
  E.SetMarginLevel(Lever)
  E.SetStockType("BTC/USD")
  E.SetContractType("this_week")
  /*
  G.MailSet("smtp.163.com", "465", "xiyanxiyan10@163.com", "xiyanxiyan10")
  G.MailStart()
  while(G.MailStatus() != "run"){
    G.Sleep(1000);
  }

  var mailCount = 0
  while(G.MailStatus() == "run"){
    mailCount++
    G.MailSend("hi:" + String(mailCount), "873706510@qq.com")
    G.Sleep(1000);
  }
  */

  queryQueue.push(onDepth); 
  //queryQueue.push(onAccount);
  queryQueue.push(onTicker);
  queryQueue.push(onOrders);
  queryQueue.push(onPosition);
  queryQueue.push(onProcess);
  

  while (true) {
        var fn = queryQueue.pop()
        fn()
        queryQueue.push(fn)
        if (globalInfo.run == 0) {
            break
        }
        G.Sleep(Intervel);
  }
  G.Log("exit process success")
}

 249  - 0framework/draw.js
@ -0,0 +1,249 @@
var chart = null
var series = []
var labelIdx = []
var preBarTime = 0
var preFlagTime = 0
var preDotTime = []
var hasPrimary = false;

var cfg = {
    tooltip: {
        xDateFormat: '%Y-%m-%d %H:%M:%S, %A'
    },
    legend: {
        enabled: true,
    },
    plotOptions: {
        candlestick: {
            color: '#d75442',
            upColor: '#6ba583'
        }
    },
    rangeSelector: {
        buttons: [{
            type: 'hour',
            count: 1,
            text: '1h'
        }, {
            type: 'hour',
            count: 3,
            text: '3h'
        }, {
            type: 'hour',
            count: 8,
            text: '8h'
        }, {
            type: 'all',
            text: 'All'
        }],
        selected: 2,
        inputEnabled: true
    },
    series: series,
}

$.GetCfg = function() {
    return cfg
}

$.PlotHLine = function(value, label, color, style) {
    if (typeof(cfg.yAxis) === 'undefined') {
        cfg.yAxis = {
            plotLines: []
        }
    } else if (typeof(cfg.yAxis.plotLines) === 'undefined') {
        cfg.yAxis.plotLines = []
    }
    /*
    Solid
    ShortDash
    ShortDot
    ShortDashDot
    ShortDashDotDot
    Dot
    Dash
    LongDash
    DashDot
    LongDashDot
    LongDashDotDot
    */
    var obj = {
        value: value,
        color: color || 'red',
        width: 2,
        dashStyle: style || 'Solid',
        label: {
            text: label || '',
            align: 'center'
        },
    }
    var found = false
    for (var i = 0; i < cfg.yAxis.plotLines.length; i++) {
        if (cfg.yAxis.plotLines[i].label.text == label) {
            cfg.yAxis.plotLines[i] = obj
            found = true
        }
    }
    if (!found) {
        cfg.yAxis.plotLines.push(obj)
    }
    if (!chart) {
        chart = Chart(cfg)
    } else {
        chart.update(cfg)
    }
}

$.PlotRecords = function(records, title) {
    var seriesIdx = labelIdx["candlestick"];
    if (!chart) {
        chart = Chart(cfg)
        chart.reset()
    }
    if (typeof(seriesIdx) == 'undefined') {
        cfg.__isStock = true
        seriesIdx = series.length
        series.push({
            type: 'candlestick',
            name: typeof(title) == 'undefined' ? '' : title,
            id: (hasPrimary ? 'records_' + seriesIdx : 'primary'),
            data: []
        });
        chart.update(cfg)
        labelIdx["candlestick"] = seriesIdx
    }
    hasPrimary = true;
    if (typeof(records.Time) !== 'undefined') {
        var Bar = records;
        if (Bar.Time == preBarTime) {
            chart.add(seriesIdx, [Bar.Time, Bar.Open, Bar.High, Bar.Low, Bar.Close], -1)
        } else if (Bar.Time > preBarTime) {
            preBarTime = Bar.Time
            chart.add(seriesIdx, [Bar.Time, Bar.Open, Bar.High, Bar.Low, Bar.Close])
        }
    } else {
        for (var i = 0; i < records.length; i++) {
            if (records[i].Time == preBarTime) {
                chart.add(seriesIdx, [records[i].Time, records[i].Open, records[i].High, records[i].Low, records[i].Close], -1)
            } else if (records[i].Time > preBarTime) {
                preBarTime = records[i].Time
                chart.add(seriesIdx, [records[i].Time, records[i].Open, records[i].High, records[i].Low, records[i].Close])
            }
        }
    }
    return chart
}

$.PlotLine = function(label, dot, time) {
    if (!chart) {
        cfg.xAxis = {
            type: 'datetime'
        }
        chart = Chart(cfg)
        chart.reset()
    }
    var seriesIdx = labelIdx[label]
    if (typeof(seriesIdx) === 'undefined') {
        seriesIdx = series.length
        preDotTime[seriesIdx] = 0
        labelIdx[label] = seriesIdx;
        series.push({
            type: 'line',
            id: (hasPrimary ? 'line_' + seriesIdx : 'primary'),
            yAxis: 0,
            showInLegend: true,
            name: label,
            data: [],
            tooltip: {
                valueDecimals: 5
            }
        })
        hasPrimary = true;
        chart.update(cfg)
    }
    if (typeof(time) == 'undefined') {
        time = new Date().getTime()
    }

    if (preDotTime[seriesIdx] != time) {
        preDotTime[seriesIdx] = time
        chart.add(seriesIdx, [time, dot])
    } else {
        chart.add(seriesIdx, [time, dot], -1)
    }

    return chart
}


$.PlotFlag = function(time, text, title, shape, color) {
    if (!chart) {
        chart = Chart(cfg)
        chart.reset()
    }
    label = "flag";
    var seriesIdx = labelIdx[label]
    if (typeof(seriesIdx) === 'undefined') {
        seriesIdx = series.length
        labelIdx[label] = seriesIdx;
        series.push({
            type: 'flags',
            onSeries: 'primary',
            data: []
        })
        chart.update(cfg)
    }

    var obj = {
        x: time,
        color: color,
        shape: shape,
        title: title,
        text: text
    }

    if (preFlagTime != time) {
        preFlagTime = time
        chart.add(seriesIdx, obj)
    } else {
        chart.add(seriesIdx, obj, -1)
    }

    return chart
}

$.PlotTitle = function(title, chartTitle) {
    cfg.subtitle = {
        text: title
    };
    if (typeof(chartTitle) !== 'undefined') {
        cfg.title = {
            text: chartTitle
        };
    }
    if (chart) {
        chart.update(cfg)
    }
}

function main() {
    var isFirst = true
    while (true) {
        var records = exchange.GetRecords();
        if (records && records.length > 0) {
            $.PlotRecords(records, 'BTC')
            if (isFirst) {
                $.PlotFlag(records[records.length - 1].Time, 'Start', 'S')
                isFirst = false
                $.PlotHLine(records[records.length - 1].Close, 'Close')
            }
        }
        var ticker = exchange.GetTicker()
        if (ticker) {
            $.PlotLine('Last', ticker.Last)
            $.PlotTitle('Last ' + ticker.Last)
        }

        Sleep(60000)
    }
}

 381  - 0framework/framework.js
@ -0,0 +1,381 @@
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

 88  - 0huobi/autocover.js
@ -0,0 +1,88 @@
// author mhw
// brief safe box 
// created 201910281001
// update  201912250000

// 平仓
function cover(pos) {
    var closeId = 0
    if (pos.length > 0) {
        var leftAmount = pos[0].Amount
        while(leftAmount > 0){
            //分批平仓
            var subAmount = leftAmount > MaxCover ? MaxCover : leftAmount
            if (pos[0].Type == 0) { //平多仓，采用盘口吃单，会损失手续费，可改为盘口挂单，会增加持仓风险。
                Log("平多仓", subAmount)  
                exchange.SetDirection('sell')
                var closeId = exchange.Sell(-1, subAmount, '平多仓')
            } else {
                Log("平空仓", subAmount)
                exchange.SetDirection('buy')
                var closeId = exchange.Buy(-1, subAmount, '平空仓')
            }
            Sleep(Intervel)
            leftAmount -= subAmount
        }
    }
}

// watchBox 箱体检测
function watchBox(price) {
    if (price > HighBox) {
        return true
    }
    if (price < LowBox) {
        return true
    }
    return false
}


function main() {
    exchange.SetContractType('swap')
    Log(exchange.GetAccount())
    while (true) {

        ticker_curr = exchange.GetTicker()
        Sleep(Intervel)
        last_price = ticker_curr.Last

        // 箱体,安全检测
        if (watchBox(last_price)) {
            Log("异常价格:" + JSON.stringify({
                lastPrice: last_price,
                HighBox: HighBox,
                LowBox: LowBox
            }))
            break
        }else{
            LogStatus(JSON.stringify({
                lastPrice: last_price,
                HighBox: HighBox,
                LowBox: LowBox
            }))
        }
        Sleep(Intervel)
       
    }

    // 平仓清理
    while (true) {
        Log("异常情况尝试平仓")
        // 定时
        var pos = exchange.GetPosition()
        if(typeof pos == "undefined" || pos == null){
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

