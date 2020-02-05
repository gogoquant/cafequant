// 滑点
SlidePrice = 0.1
// 箱体上沿
HighBox = 9900
// 箱体下沿
LowBox = 8000
// 网格方向
BuyFirst = 1
// 网格数量
GridNum = 10
// 网格价格距离
PriceGrid = 50
// 价格精度
Precision = 1
// 开仓保护价差
OpenProtectDiff = 10
// 买单数量
BAmountOnce = 10
// 卖单数量
SAmountOnce = 10
// 是否止损
EnableStopLoss = 1
// 止损模式
StopLossMode = 0
// 止损盈亏损率
StopLoss = 20
// 最小量
MinStock = 0.1
// 是否自动移动价格
AutoMove = 1
// 仓位当前价格最大价差
MaxDistance = 300
// 最大空仓时间
HoldTime = 1000 * 60 * 10
// 收网周期
FishCheckTime = 1000 * 60 * 5
// 最小周期
Interval = 1000 * 60 * 2

// local
var globalInfo = {

}

var ProfitCount = 0;
var globalInfo = {}
var STATE_WAIT_OPEN = "wait_open";
var STATE_WAIT_COVER = "wait_cover";
var STATE_WAIT_CLOSE = "wait_close";
var STATE_END_CLOSE = "end_close";
var ORDER_TYPE_BUY = 0;
var ORDER_TYPE_SELL = 1;

/*
 * Only used for test
 */

function Log(val) {
    console.log(val)
}

Exchange = function() {
    cnt = 1
    this.Buy = function(price, amount, extra) {
        Log("Buy price: " + price + ", amount: " + amount + ", extra:" + extra)
        cnt += 1
        return cnt
    }

    this.Sell = function(price, amount, extra) {
        Log("Sell price: " + price + ", amount: " + amount + ", extra:" + extra)
        cnt++;
        return cnt
    }
}

exchange = new Exchange()

function _N(num, pre) {
    var str1 = String(num); //将类型转化为字符串类型
    if (str1.indexOf('.') < 0) {
        return num
    }
    var str = str1.substring(0, str1.indexOf('.') + pre + 1); //截取字符串
    var num = Number(str); //转化为number类型
    return num
}


// 1. 准备撒网阶段首先以当前买一或卖一位置预留一段保护价格算出网格起使价格
// 2. 判断资金量是否充足以及利用起使价格计算网格位置
// 3. 撒网阶段 注册计算出的网格, 注册的网格开始状态都为待挂单状态
// 4. 注册的网格点距离当前价格接近一定触发范围后才从待挂单转为挂单，过远的单转化为待挂单状态，避免距离过远的无意义单占用可挂单数以及产生具有一定价差的主动成交
// 5. 当所挂单被被吃掉后，立即根据设置的盈利价差挂反向单
// 6. 在网格周期中当在超过指定时间未持仓和当前价格偏离撒网价格过大时则收网进入下一轮重新撒网， ??平仓撒网获取的仓位进入下一轮撒网
// 7. 当持仓超过指定盈利阈值时则止盈关闭所有订单 ??平仓撒网获取的仓位进入下一轮撒网
// 8. 当持仓超过指定盈利阈值时则止损关闭所有订单 ??平仓撒网获取的仓位进入下一轮撒网
// 9. 平仓单目前是直接加上固定差价算出位置，是否需要更灵活的布局方案
// 10. 多空仓位对冲的双向网格,对冲时多空比例的控制
// 11. 实际使用中发现，是否各个网格的量应该距离fish开始价格越远则量按照一定比例更为离散

// ArrayQueue 队列
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

// TradeRobot 交易者,用于交易配置以及延时处理等
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


function ValidItem(val) {
    if (typeof val == undefined || val == null) {
        return false
    }
    return true
}

function initInfo() {
    globalInfo = {}
}

function onDepth() {
    globalInfo.depth = exchange.GetDepth()
}

function onPosition() {
    globalInfo.position = exchange.GetPosition()
}

function onTicker() {
    globalInfo.ticker = exchange.GetTicker()
}

function onOrders() {
    globalInfo.orders = exchange.GetOrders()
}

function onAccount() {
    globalInfo.account = exchange.GetAccount()
}

function checkInfo() {
    var obj = globalInfo
    for (var key in obj) {
        console.log("get key :" + key)
        if (!ValidItem(obj[key])) {
            Log(key + ": get fail from exchange")
            return false
        }
    }
    return true
}

function blockGetInfo() {
    initInfo()
    do {
        for (var i = 0; i < arguments.length; i++) {
            var fn = arguments[i]
            fn()
        }
    } while (!checkInfo());
}

function hasOrder(orders, orderId) {
    for (var i = 0; i < orders.length; i++) {
        if (orders[i].Id == orderId) {
            return true;
        }
    }
    return false;
}

function foundOrder(orders, orderId) {
    for (var i = 0; i < orders.length; i++) {
        if (orders[i].Id == orderId) {
            return orders[i];
        }
    }
    return null;
}

//阻塞关闭订单
function cancelPending() {
    var ret = false;
    while (true) {
        if (ret) {
            Sleep(Interval);
        }
        blockGetInfo(onOrders)
        var orders = globalInfo.orders
        if (orders.length == 0) {
            break;
        }
        for (var j = 0; j < orders.length; j++) {
            exchange.CancelOrder(orders[j].Id, orders[j]);
            ret = true;
        }
    }
    return ret;
}

//阻塞关闭一个订单
function cancelOnePending(Id) {
    var ret = false;
    while (true) {
        if (ret) {
            Sleep(Interval);
        }
        blockGetInfo(onOrders)
        var orders = globalInfo.orders
        order = foundOrder(orders, Id)
        if (order == null) {
            break;
        } else {
            exchange.CancelOrder(order.Id, order);
            ret = true;
        }
    }
    return ret;
}

function valuesToString(values, pos) {
    var result = '';
    if (typeof(pos) === 'undefined') {
        pos = 0;
    }
    for (var i = pos; i < values.length; i++) {
        if (i > pos) {
            result += ' ';
        }
        if (values[i] === null) {
            result += 'null';
        } else if (typeof(values[i]) == 'undefined') {
            result += 'undefined';
        } else {
            switch (values[i].constructor.name) {
                case 'Date':
                case 'Number':
                case 'String':
                case 'Function':
                    result += values[i].toString();
                    break;
                default:
                    result += JSON.stringify(values[i]);
                    break;
            }
        }
    }
    return result;
}


function GridTrader() {
    var vId = 0;
    var orderBooks = [];
    var hisBooks = [];
    var orderBooksLen = 0;
    var hisBooksLen = 0;
    var profitPrice = -1

    this.SetProfitPrice = function(price) {
        profitPrice = price
    }

    this.GetProfitPrice = function() {
        return profitPrice
    }

    this.Buy = function(price, amount, extra) {
        if (typeof(extra) === 'undefined') {
            extra = '';
        } else {
            extra = valuesToString(arguments, 2);
        }
        vId++;
        var orderId = "V" + vId;
        orderBooks[orderId] = {
            Type: ORDER_TYPE_BUY,
            Status: STATE_WAIT_OPEN,
            OpenId: 0,
            CoverId: 0,
            Price: price,
            VID: orderId,
            CoverPrice: price + this.GetProfitPrice(),
            Amount: amount,
            Extra: extra
        };
        orderBooksLen++;
        return orderId;
    };
    this.Sell = function(price, amount, extra) {
        if (typeof(extra) === 'undefined') {
            extra = '';
        } else {
            extra = valuesToString(arguments, 2);
        }
        vId++;
        var orderId = "V" + vId;
        orderBooks[orderId] = {
            Type: ORDER_TYPE_SELL,
            Status: STATE_WAIT_OPEN,
            OpenId: 0,
            CoverId: 0,
            Price: price,
            VID: orderId,
            CoverPrice: price - this.GetProfitPrice(),
            Amount: amount,
            Extra: extra
        };
        orderBooksLen++;
        return orderId;
    };

    this.GetOrders = function() {
        return orderBooks;
    }

    this.GetOrder = function(orderId) {
        if (typeof(orderId) === 'number') {
            return exchange.GetOrder(orderId);
        }
        if (typeof(hisBooks[orderId]) !== 'undefined') {
            return hisBooks[orderId];
        }
        if (typeof(orderBooks[orderId]) !== 'undefined') {
            return orderBooks[orderId];
        }
        return null;
    };

    this.Len = function() {
        return orderBooksLen;
    };

    this.RealLen = function() {
        var n = 0;
        for (orderId in orderBooks) {
            if (orderBooks[orderId].Id > 0) {
                n++;
            }
        }
        return n;
    };


    // 一个格子的状态机转换
    this.PollOne = function(order, ticker, exchangeOrders, priceDiff) {
        var pfn = order.Type == ORDER_TYPE_BUY ? exchange.Buy : exchange.Sell;
        var coverPfn = order.Type == ORDER_TYPE_BUY ? exchange.Sell : exchange.Buy;

        // 等待开仓的订单
        if (order.Status == STATE_WAIT_OPEN) {
            var diff = _N(order.Type == ORDER_TYPE_BUY ? (ticker.Buy - order.Price) : (order.Price - ticker.Sell));
            if (diff > 0 && diff <= priceDiff) {
                var realId = pfn(order.Price, order.Amount, order.Extra + "(距离: " + diff + (order.Type == ORDER_TYPE_BUY ? (" 买一: " + ticker.Buy) : (" 卖一: " + ticker.Sell)) + ")" + " VID:" + order.VID);
                if (typeof(realId) === 'number') {
                    order.OpenId = realId;
                    order.Status = STATE_WAIT_COVER
                }
            }
            return
        }

        // 等待平仓的订单
        if (order.Status == STATE_WAIT_COVER) {
            var found = hasOrder(exchangeOrders, order.OpenId)
            if (!found) {
                var realId = coverPfn(order.CoverPrice, order.Amount, order.Extra + " 平仓 价格:" + order.CoverPrice + " VID:" + order.VID);
                if (typeof(realId) === 'number') {
                    order.CoverId = realId;
                    order.Status = STATE_WAIT_CLOSE
                }
            }
            return
        }

        //  等待完结的订单
        if (order.Status == STATE_WAIT_CLOSE) {
            Log("close order:" + order.CoverId)
            var found = hasOrder(exchangeOrders, order.CoverId)
            if (!found) {
                order.Status = STATE_END_CLOSE;
            }
            return
        }
    }

    // 遍历所有各自尝试转换状态机
    this.Poll = function(ticker, orders, priceDiff) {
        for (orderId in orderBooks) {
            var order = orderBooks[orderId];
            this.PollOne(order, ticker, orders, priceDiff)
        }
    }
}

// 动态再平衡
function balanceAccount(orgAccount, initAccount) {
    cancelPending();
    var slidePrice = SlidePrice
    //var slidePrice = 0.1;
    var ok = true;
    while (true) {
        blockGetInfo(onAccount, onDepth)
        var orders = globalInfo.orders
        var depth = globalInfo.depth
        var diff = _N(nowAccount.Stocks - initAccount.Stocks);
        if (Math.abs(diff) < MinStock) {
            break;
        }
        var books = diff > 0 ? depth.Bids : depth.Asks;
        var n = 0;
        var price = 0;
        for (var i = 0; i < books.length; i++) {
            n += books[i].Amount;
            if (n >= Math.abs(diff)) {
                price = books[i].Price;
                break;
            }
        }
        var pfn = diff > 0 ? exchange.Sell : exchange.Buy;
        var amount = Math.abs(diff);
        var price = diff > 0 ? (price - slidePrice) : (price + slidePrice);
        Log("开始平衡", (diff > 0 ? "卖出" : "买入"), amount, "个币");
        if (diff > 0) {
            amount = Math.min(nowAccount.Stocks, amount);
        } else {
            amount = Math.min(nowAccount.Balance / price, amount);
        }
        if (amount < MinStock) {
            Log("资金不足, 无法平衡到初始状态");
            ok = false;
            break;
        }
        pfn(price, amount);
        Sleep(2000);
        cancelPending();
    }
    if (ok) {
        LogProfit(_N(nowAccount.Balance - orgAccount.Balance));
        Log("平衡完成", nowAccount);
    }
}

function onexit() {
    cancelPending();
    Log("策略成功停止");
    blockGetInfo(onAccount)
    var account = globalInfo.account
    Log(account);
}


function fishing(orgAccount, fishCount) {
    // 撒网
    var gridTrader = new GridTrader();
    // 持仓定时器
    var holdTimer = new TradeRobot("hold")
    var fishCheckTimer = new TradeRobot("check")
    holdTimer.SetInterval(HoldTime)
    fishCheckTimer.SetInterval(FishCheckTime)

    blockGetInfo(onAccount, onTicker)
    var account = globalInfo.account
    var InitAccount = account;
    var ticker = globalInfo.ticker
    Log(account);

    firstPrice = BuyFirst ? _N(ticker.Buy - PriceGrid, Precision) : _N(ticker.Sell + PriceGrid, Precision);

    var needStocks = 0;
    var needMoney = 0;
    var actualNeedMoney = 0;
    var actualNeedStocks = 0;
    var notEnough = false;
    var canNum = 0;

    // 计算需要的量
    for (var idx = 0; idx < GridNum; idx++) {
        var price = _N((BuyFirst ? firstPrice - (idx * PriceGrid) : firstPrice + (idx * PriceGrid)), Precision);
        needStocks += SAmountOnce;
        needMoney += price * BAmountOnce;
        if (BuyFirst) {
            if (_N(needMoney) <= _N(account.Balance)) {
                actualNeedMondy = needMoney;
                actualNeedStocks = needStocks;
                canNum++;
            } else {
                notEnough = true;
                Log("资金不足, 需要", _N(needMoney), "元, 程序只做", canNum, "个网格 #ff0000");
            }
        } else {
            if (_N(needStocks) <= _N(account.Stocks)) {
                actualNeedMondy = needMoney;
                actualNeedStocks = needStocks;
                canNum++;
                Log("资金不足, 需要", _N(needStocks), "个币, 程序只做", canNum, "个网格 #ff0000");
            } else {
                notEnough = true;
            }
        }
        if (!notEnough) {
            if (BuyFirst) {
                gridTrader.Buy(price, BAmountOnce, "")
            } else {
                gridTrader.Sell(price, SAmountOnce, "")
            }
        }
    }

    if (canNum < GridNum) {
        Log("警告, 当前资金只可做", canNum, "个网格, 全网共需", (BuyFirst ? needMoney : needStocks), "请保持资金充足");
    }
    if (BuyFirst) {
        Log('预计动用资金: ', _N(needMoney), "元");
    } else {
        Log('预计动用币数: ', _N(needStocks), "个, 约", _N(needMoney), "元");
    }

    var profitMax = 0;
    var preMsg = ""
    while (true) {

        if (fishCheckTimer.CheckInterval()) {
            fishCheckTimer.SetInterval(FishCheckTime)

            blockGetInfo(onPosition, onAccount)
            positions = globalInfo.position
            account = globalInfo.account

            var isHold = positions.length > 0;
            if (isHold) {
                holdTimer.SetInterval(HoldTime)
            }

            if (isHold) {
                msg += "持仓: " + positions[0].Amount + " 持仓均价: " + _N(positions[0].Price) + " 浮动盈亏: " + _N(positions[0].Profit);
                if (EnableStopLoss && -positions[0].Profit >= StopLoss) {
                    Log("当前浮动盈亏", positions[0].Profit, "开始止损");
                    balanceAccount(orgAccount, InitAccount);
                    if (StopLossMode === 0) {
                        throw "止损退出";
                    } else {
                        return true;
                    }
                }
            } else {
                msg += "空仓";
            }
            msg += " 可用保证金: " + nowAccount.Stocks;

            var distance = 0;
            if (AutoMove) {
                if (BuyFirst) {
                    distance = ticker.Last - firstPrice;
                } else {
                    distance = firstPrice - ticker.Last;
                }
                var refish = false;
                if (!isHold && holdTimer.CheckInterval()) {
                    Log("空仓过久, 开始移动网格");
                    refish = true;
                }
                if (distance > MaxDistance) {
                    Log("价格超出网格区间过多, 开始移动网格, 当前距离: ", _N(distance, Precision), "当前价格:", ticker.Last);
                    refish = true;
                }
                if (refish) {
                    balanceAccount(orgAccount, InitAccount);
                    return true;
                }
            }

            if (msg != preMsg) {
                LogStatus(msg);
                preMsg = msg;
            }

        }
        blockGetInfo(onOrders, onTicker)
        var ticker = globalInfo.ticker
        var orders = globalInfo.orders
        gridTrader.Poll(ticker, orders, PriceDiff)
        Sleep(CheckInterval);
    }
    return true;
}

function main() {
    BuyFirst = (OpType == 0);
    blockGetInfo(onAccount)
    var orgAccount = globalInfo.account
    var fishCount = 1;
    return
    while (true) {
        if (!fishing(orgAccount, fishCount)) {
            break;
        }
        fishCount++;
        Log("第", fishCount, "次重新撒网...");
        Sleep(Interval);
    }
}


/* 
 *Test function for checkInfo
 */

/*
initInfo()
globalInfo.empty = null
globalInfo.undefine = undefined
globalInfo.val = 1 
console.log(global)
console.log(checkInfo() ? "checkInfo Success" : "checkInfo Fail")
*/


/*
 * Test function for gridTrader
 */


console.log(exchange.Buy(9000, 100, "buy"))
console.log(exchange.Sell(9000, 100, "sell"))
grid = new GridTrader()
orders = []
var buyPrice = 9000
var sellPrice = 9001
var amount = 100
var ticker = Object()
ticker.Buy = 9000
ticker.Sell = 9000
ticker.Last = 8001
grid.SetProfitPrice(50)
Log("profit price is :" + grid.GetProfitPrice())

Log("init process")
grid.Sell(sellPrice, amount, "Sell")
grid.Sell(sellPrice, amount, "Sell")
grid.Sell(sellPrice, amount, "Sell")
grid.Sell(sellPrice, amount, "Sell")
Log(grid.GetOrders())

Log("open process")
grid.Poll(ticker, orders, 10)
Log(grid.GetOrders())

Log("cover process")
orders = [{
        Id: 5,
    },
    {
        Id: 6
    }
]
grid.Poll(ticker, orders, 10)
Log(grid.GetOrders())

orders = [{
        Id: 5,
    },
    {
        Id: 6
    }
]

grid.Poll(ticker, orders, 10)
Log(grid.GetOrders())
