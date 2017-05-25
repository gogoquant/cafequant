/*
策略出处: https://www.botvs.com/strategy/179
策略名称: 趋势跟踪震荡策略
策略作者: Zero
策略描述:

*/

//货币名称
var stockType           = ""
//行为安全间隔
var Interval            = 2000 


//浮点下取整
function adjustFloat(v) {
    return Math.floor(v*1000)/1000;
}

//关闭交易单
function CancelPendingOrders(stockType) {
    while (true) {
        var orders = null;
        while (!(orders = exchange.GetOrders(stockType))) {
            Sleep(Interval);
        }

        if (orders.length == 0) {
            return;
        }

        for (var j = 0; j < orders.length; j++) {
            E.CancelOrder(orders[j]);
        }
    }
}

//阻塞获取信息
function PendingAction(action) {
    var result;
    while (!(result = E.action())) {
        Sleep(Interval);
    }
    return account;
}

var STATE_WAIT_IDLE     = 0;
var STATE_WAIT_BUY      = 1;
var STATE_WAIT_SELL     = 2;
var STATE_BUY           = 3;
var STATE_SELL          = 4;

var State = STATE_WAIT_IDLE;
var InitAccount = null;

var LastBuyPrice = 0;
var LastSellPrice = 0;

var LastHighPrice = 0;
var LastLowPrice = 0;

var LastRecord = null;

var Goingshort = false;

function onTick(exchange) {
    var oldState = State;

    if (State == STATE_WAIT_IDLE) {
        var records = exchange.GetRecords();
        if (!records || records.length < (EMA_Slow + 3)) {
            return;
        }
        // Price not change
        var newLast = records[records.length-1];
        if ((!LastRecord) || (LastRecord.Time == newLast.Time && LastRecord.Close == newLast.Close)) {
            LastRecord = newLast;
            return;
        }
        LastRecord = newLast;

        var emaFast = TA.EMA(records, EMA_Fast);
        var emaSlow = TA.EMA(records, EMA_Slow);
        if (emaFast[emaFast.length-1] > emaSlow[emaSlow.length-1]) {
            Goingshort = false;
            State = STATE_WAIT_BUY;
        } else if (EnableGoingShort && (emaFast[emaFast.length-1] < emaSlow[emaSlow.length-1])) {
            Goingshort = true;
            State = STATE_WAIT_SELL;
        } else {
            return;
        }
    }

    var ticker = GetTicker();
    // 重新设置这两个参数
    if (oldState == STATE_WAIT_IDLE && State != STATE_WAIT_IDLE) {
        LastLowPrice = ticker.Last;
        LastHighPrice = ticker.Last;
    }

    // 做多
    if (!Goingshort) {
        if (State == STATE_WAIT_BUY) {
            var lastDownRatio = Math.abs((LastHighPrice - LastLowPrice) / LastHighPrice) * 100;
            var currentUpRatio = Math.abs((ticker.Last - LastLowPrice) / ticker.Last) * 100;
            if (lastDownRatio > RatioDown && currentUpRatio > (RatioDown*(UpWeightingVal/100))) {
                State = STATE_BUY;
            } else {
                LastHighPrice = Math.max(LastHighPrice, ticker.Last);
                LastLowPrice = Math.min(LastLowPrice, ticker.Last);
            }
        } else if (State == STATE_WAIT_SELL) {
            var ratioStopLoss = Math.abs((LastHighPrice - ticker.Last) / LastHighPrice) * 100;
            var ratioStopProfit = Math.abs((ticker.Last - LastBuyPrice) / LastBuyPrice) * 100;
            var ratioMaxUp = Math.abs((LastHighPrice - LastBuyPrice) / LastBuyPrice) * 100;
            if (ticker.Last < LastBuyPrice && ratioStopLoss >= StopLoss) {
                State = STATE_SELL;
                Log("开始止损, 当前下跌点数:", adjustFloat(ratioStopLoss), "当前价格", ticker.Last, "对比价格", adjustFloat(LastHighPrice));
            } else if (ticker.Last > LastBuyPrice && ticker.Last < LastHighPrice && ratioStopProfit <= (ratioMaxUp*StopProfitThreshold)) {
                State = STATE_SELL;
                Log("开始止赢, 当前上涨点数:", adjustFloat(ratioStopProfit), "当前价格", ticker.Last, "对比价格", adjustFloat(LastBuyPrice));
            }
            LastHighPrice = Math.max(LastHighPrice, ticker.Last);
        }
    } else {
        if (State == STATE_WAIT_SELL) {
            var lastUpRatio = Math.abs((LastHighPrice - LastLowPrice) / LastHighPrice) * 100;
            var currentDownRatio = Math.abs((LastHighPrice - ticker.Last) / LastHighPrice) * 100;
            if (lastUpRatio > RatioUp && currentDownRatio > (RatioUp*(DownWeightingVal/100))) {
                State = STATE_SELL;
            } else {
                LastHighPrice = Math.max(LastHighPrice, ticker.Last);
                LastLowPrice = Math.min(LastLowPrice, ticker.Last);
            }
        } else if (State == STATE_WAIT_BUY) {
            var ratioStopLoss = Math.abs((ticker.Last - LastLowPrice) / LastLowPrice) * 100;
            var ratioStopProfit = Math.abs((LastSellPrice - ticker.Last) / LastSellPrice) * 100;
            var ratioMaxDown = Math.abs((LastSellPrice - LastLowPrice) / LastSellPrice) * 100;
            if (ticker.Last > LastSellPrice && ratioStopLoss >= StopLoss) {
                State = STATE_BUY;
                Log("开始止损, 当前上涨点数:", adjustFloat(ratioStopLoss), "当前价格", ticker.Last, "对比价格", adjustFloat(LastSellPrice));
            } else if (ticker.Last < LastSellPrice && ticker.Last > LastLowPrice && ratioStopProfit <= (ratioMaxDown*StopProfitThreshold)) {
                State = STATE_BUY;
                Log("开始止盈, 当前下跌点数:", adjustFloat(ratioStopProfit), "当前价格", ticker.Last, "对比价格", adjustFloat(LastLowPrice));
            }
            LastHighPrice = Math.max(LastHighPrice, ticker.Last);
            LastLowPrice = Math.min(LastLowPrice, ticker.Last);
        }

    }

    if (State != STATE_BUY && State != STATE_SELL) {
        return;
    }

    // Buy or Sell, Cancel pending orders first
    CancelPendingOrders();

    var account = GetAccount();

    // 做多
    if (!Goingshort) {
        if (State == STATE_BUY) {
            var price = ticker.Last + SlidePrice;
            var amount = adjustFloat(account.Balance / price);
            if (amount >= exchange.GetMinStock()) {
                if (exchange.Buy(price, amount, "做多")) {
                    LastBuyPrice = LastHighPrice = price;
                }
            } else {
                State = STATE_WAIT_SELL;
            }
        } else {
            var sellAmount = account.Stocks - InitAccount.Stocks;
            if (sellAmount > exchange.GetMinStock()) {
                exchange.Sell(ticker.Last - SlidePrice, sellAmount);
            } else {
                // No stocks, wait buy and log profit
                LogProfit(account.Balance - InitAccount.Balance, account);
                State = STATE_WAIT_IDLE;
            }
        }
    } else {
        if (State == STATE_BUY) {
            var price = ticker.Last + SlidePrice;
            var amount = Math.min(adjustFloat(account.Balance / price), InitAccount.Stocks - account.Stocks);
            if (amount >= exchange.GetMinStock()) {
                exchange.Buy(price, amount);
            } else {
                LogProfit(account.Balance - InitAccount.Balance, account);
                State = STATE_WAIT_IDLE;
            }
        } else {
            var price = ticker.Last - SlidePrice;
            var sellAmount = account.Stocks;
            if (sellAmount > exchange.GetMinStock()) {
                exchange.Sell(ticker.Last - SlidePrice, sellAmount, "做空");
                LastSellPrice = LastLowPrice = price;
            } else {
                // No stocks, wait buy and log profit
                State = STATE_WAIT_BUY;
            }
        }
    }
}

function main() {
    InitAccount = GetAccount();
    Log(exchange.GetName(), exchange.GetCurrency(), InitAccount);
    EnableGoingShort = EnableGoingShort && (InitAccount.Stocks > exchange.GetMinStock());
    LoopInterval = Math.max(LoopInterval, 1);
    while (true) {
        onTick(exchange);
        Sleep(LoopInterval*1000);
    }
}