// Coin type
var Coin = "BTC/USD";
// 箱体上沿
var HighBox = 12500;
// 箱体下沿
var LowBox = 11500;
// 计划持仓量
var HighPosition = 12;
// 网格价格距离
var GridOffset = 25;
// 价格精度
var Precision = 1;
// 开仓保护价差
var OpenProtect = 3;
// 对冲仓位差
var PositionDiff = 4;
// 买单数量
var AmountOnce = 1;
// 止损后模式
var StopLossAfterMode = 0;
// 止损盈亏损率
var StopLoss = -5.0;
// 止盈率
var StopWin = -5.0;
// 最大空仓时间
var HoldTime = 1000 * 60 * 5;
// 收网检测周期
var FishCheckTime = 1000 * 20 * 1;
// 最小周期
var Interval = 1000 * 10 * 1;
// 盈利滑动
var ProfitPrice = 12;
// 合约
var ContractType = "quarter";
// 杠杆
var MarginLevel = 20;
// 合约列表
var ContractVec = ["this_week", "next_week", "quarter"];
// 货币支持类型
var CoinVec = ["BTC/USD"];
// tot married order
var marryTot = 0;
// local
var globalInfo = {};

var STATE_WAIT_OPEN = "wait_open";
var STATE_WAIT_COVER = "wait_cover";
var STATE_WAIT_CLOSE = "wait_close";
var STATE_END_CLOSE = "end_close";
var ORDER_TYPE_BUY = 0;
var ORDER_TYPE_SELL = 1;

ExchangeDebug = function() {
  cnt = 1;
  this.Buy = function(price, amount, extra) {
    G.Log("Buy price: " + price + ", amount: " + amount + ", extra:" + extra);
    cnt += 1;
    return cnt;
  };
  this.Sell = function(price, amount, extra) {
    G.Log("Sell price: " + price + ", amount: " + amount + ", extra:" + extra);
    cnt++;
    return cnt;
  };
  this.CancelOrder = function(dir) {
    G.Log("cancel order:", dir);
    return cnt;
  };
  this.SetContractType = function(contract) {
    return E.SetContractType(contract);
  };
  this.SetMarginLevel = function(contract) {
    return E.SetMarginLevel(contract);
  };
  this.GetOrder = function(orderId) {
    return E.GetOrder(orderId);
  };
  this.GetOrders = function() {
    return E.GetOrders();
  };
  this.SetDirection = function(dir) {
    return E.SetDirection(dir);
  };
  this.SetStockType = function(dir) {
    return E.SetStockType(dir);
  };
};

ExchangeReal = function() {
  this.Buy = function(price, amount, extra) {
    return E.Buy(price, amount, extra);
  };
  this.Sell = function(price, amount, extra) {
    return E.Sell(price, amount, extra);
  };
  this.SetDirection = function(dir) {
    return E.SetDirection(dir);
  };
  this.CancelOrder = function(dir) {
    return E.CancelOrder(dir);
  };
  this.SetContractType = function(contract) {
    return E.SetContractType(contract);
  };
  this.SetMarginLevel = function(contract) {
    return E.SetMarginLevel(contract);
  };
  this.GetOrder = function(orderId) {
    return E.GetOrder(orderId);
  };
  this.GetOrders = function() {
    return E.GetOrders();
  };
  this.SetStockType = function(dir) {
    return E.SetStockType(dir);
  };
};

exchange = new ExchangeReal();

function _N(num, pre) {
  var str1 = String(num); //将类型转化为字符串类型
  if (str1.indexOf(".") < 0) {
    return num;
  }
  var str = str1.substring(0, str1.indexOf(".") + pre + 1); //截取字符串
  var num = Number(str); //转化为number类型
  return num;
}

// 1. 空方向测试
// 2. 多空双向使用一套程序管理两个仓位?
// 3. 平仓挂单和开仓挂单本质无强关联，是否分批平仓单策略只通过当前仓为平均价格来推算挂平仓单的方式?
// 4. 考虑使用其他算法，调整网格的离散度?

// ArrayQueue 队列
function ArrayQueue() {
  var arr = [];
  //入队操作
  this.push = function(element) {
    arr.push(element);
    return true;
  };
  //出队操作
  this.pop = function() {
    return arr.shift();
  };
  //获取队首
  this.getFront = function() {
    return arr[0];
  };
  //获取队尾
  this.getRear = function() {
    return arr[arr.length - 1];
  };
  //清空队列
  this.clear = function() {
    arr = [];
  };
  //获取队长
  this.size = function() {
    return arr.length;
  };
  this.printOne = function() {
    if (arr.length <= 0) {
      return 0;
    }
    return this.getFront();
  };
}

// TradeRobot 交易者,用于交易配置以及延时处理等
function TradeRobot(name) {
  this.name = name;
  this.price = 0; //下单价
  this.amount = 0; //下单量
  this.leftTime = 0; //用于显示
  this.interval = -1;
  this.nextTime = new Date().valueOf();
  this.run = 1;
}

TradeRobot.prototype.Name = function() {
  return this.name;
};

TradeRobot.prototype.SetInterval = function(num) {
  this.interval = num;
  this.leftTime = 0 - num;
  this.nextTime = new Date().valueOf() + num;
};

TradeRobot.prototype.AddInterval = function(num) {
  this.leftTime = new Date().valueOf() - this.nextTime;
  if (this.leftTime < 0) {
    this.interval += num;
    this.nextTime += num;
    this.leftTime -= num;
    return;
  }
  this.interval = num;
  this.leftTime = 0 - num;
  this.nextTime = new Date().valueOf() + num;
};

TradeRobot.prototype.Timeout = function() {
  this.leftTime = this.nextTime - new Date().valueOf();
  return this.leftTime >= 0 ? 0 : 1;
};

TradeRobot.prototype.TimeLeft = function() {
  return this.nextTime - new Date().valueOf();
};

TradeRobot.prototype.Left = function() {
  this.leftTime = new Date().valueOf() - this.nextTime;
  return this.leftTime;
};

TradeRobot.prototype.Interval = function() {
  return this.interval;
};

function ValidItem(val) {
  if (typeof val == undefined || val == null) {
    return false;
  }
  return true;
}

function position2Rate(position, price) {
  return position.ProfitRate * 100.0;
}

function orderClean(orders) {
  for (var i = 0; i < orders.length; i += 1) {
    orders[i].Info = null;
  }
  return orders;
}

// order2DirOrder normal filter by dir
function order2DirOrder(orders, dir) {
  if (dir === -1) {
    return orders;
  }
  var ordervec = [];
  for (var i = 0; i < orders.length; i += 1) {
    //@ Todo closebuy and closesell reverse in orders
    if (
      dir === 0 &&
      (orders[i].TradeType === "buy" || orders[i].TradeType === "closebuy")
    ) {
      ordervec.push(orders[i]);
    }
    if (
      dir === 1 &&
      (orders[i].TradeType === "sell" || orders[i].TradeType === "closesell")
    ) {
      ordervec.push(orders[i]);
    }
  }
  return ordervec;
}

function initInfo() {
  globalInfo = {};
}

//function onDepth() {
//  globalInfo.depth = exchange.GetDepth();
//}

function onPosition() {
  globalInfo.positions = E.GetPosition();
}

function onTicker() {
  globalInfo.ticker = E.GetTicker();
}

function onOrders() {
  globalInfo.orders = E.GetOrders();
}

function onAccount() {
  globalInfo.account = E.GetAccount();
}

function checkInfo() {
  var obj = globalInfo;
  for (var key in obj) {
    console.log("get key :" + key);
    if (!ValidItem(obj[key])) {
      G.Log(key + ": get fail from exchange");
      return false;
    }
  }
  return true;
}

function blockGetInfo() {
  initInfo();
  do {
    for (var i = 0; i < arguments.length; i++) {
      var fn = arguments[i];
      fn();
    }
  } while (!checkInfo());
}

function Order2Cost(price, amount, last) {
  return _N((1.0 * price * amount * 0.01) / last);
}

function hasOrder(orders, orderId) {
  var orders = order2DirOrder(orders, BuyFirst ? 0 : 1);
  for (var i = 0; i < orders.length; i++) {
    if (orders[i].Id == orderId) {
      return true;
    }
  }
  return false;
}

function foundOrder(orders, orderId) {
  var orders = order2DirOrder(orders, BuyFirst ? 0 : 1);
  for (var i = 0; i < orders.length; i++) {
    if (orders[i].Id == orderId) {
      return orders[i];
    }
  }
  return null;
}

// 阻塞关闭订单
function cancelPending(dir) {
  var ret = false;
  var cycle = true;
  while (cycle) {
    if (ret) {
      G.Sleep(Interval);
    }
    blockGetInfo(onOrders);
    var orders = order2DirOrder(globalInfo.orders, dir);
    if (orders.length == 0) {
      break;
    }
    for (var j = 0; j < orders.length; j++) {
      exchange.CancelOrder(orders[j].Id);
      ret = true;
    }
  }
  return ret;
}

// 阻塞关闭一个订单
function cancelOnePending(Id) {
  var ret = false;
  var cycle = true;
  while (cycle) {
    if (ret) {
      G.Sleep(Interval);
    }
    blockGetInfo(onOrders);
    var orders = globalInfo.orders;
    var order = foundOrder(orders, Id);
    if (order == null) {
      break;
    } else {
      exchange.CancelOrder(order.Id);
      ret = true;
    }
  }
  return ret;
}

function valuesToString(values, pos) {
  var result = "";
  if (typeof pos === "undefined") {
    pos = 0;
  }
  for (var i = pos; i < values.length; i++) {
    if (i > pos) {
      result += " ";
    }
    if (values[i] === null) {
      result += "null";
    } else if (typeof values[i] == "undefined") {
      result += "undefined";
    } else {
      switch (values[i].constructor.name) {
        case "Date":
        case "Number":
        case "String":
        case "Function":
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

function account2balance(account, symbol) {
  var accounts = account.SubAccounts;
  if (accounts.hasOwnProperty(symbol)) {
    return accounts[symbol].AccountRights;
  }
  return null;
}

function GridTrader() {
  var vId = 0;
  var realId;
  var found;
  var orderBooks = new Object();
  var hisBooks = new Object();
  var orderBooksLen = 0;
  var hisBooksLen = 0;
  var profitPrice = -1;
  var orderId;

  this.SetProfitPrice = function(price) {
    profitPrice = price;
  };

  this.GetProfitPrice = function() {
    return profitPrice;
  };

  this.BooksLen = function() {
    var obj = new Object();
    obj.wait_open = 0;
    obj.wait_cover = 0;
    obj.wait_close = 0;
    obj.history = hisBooksLen;
    obj.curr = orderBooksLen;
    for (orderId in orderBooks) {
      var order = orderBooks[orderId];
      if (order.Status == STATE_WAIT_OPEN) {
        obj.wait_open++;
      }
      if (order.Status == STATE_WAIT_COVER) {
        obj.wait_cover++;
      }

      if (order.Status == STATE_WAIT_CLOSE) {
        obj.wait_close++;
      }
    }
    return obj;
  };

  this.Debug = function() {
    G.Log("Orders List:");
    for (orderId in orderBooks) {
      G.Log(orderBooks[orderId]);
    }
    G.Log("HisOrders List:");
    for (orderId in hisBooks) {
      G.Log(hisBooks[orderId]);
    }
    return;
  };

  this.Buy = function(price, amount, extra) {
    G.Log("Buy price: " + price + ", amount: " + amount + ", extra:" + extra);
    if (typeof extra === "undefined") {
      extra = "";
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
    G.Log("Sell price: " + price + ", amount: " + amount + ", extra:" + extra);
    if (typeof extra === "undefined") {
      extra = "";
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
  };

  this.GetHistoryOrders = function() {
    return hisBooks;
  };

  this.GetOrder = function(orderId) {
    if (typeof orderId === "string") {
      return exchange.GetOrder(orderId);
    }
    if (typeof hisBooks[orderId] !== "undefined") {
      return hisBooks[orderId];
    }
    if (typeof orderBooks[orderId] !== "undefined") {
      return orderBooks[orderId];
    }
    return null;
  };

  // 一个格子的状态机转换
  this.PollOne = function(order, ticker, exchangeOrders, ext) {
    var pfn = order.Type == ORDER_TYPE_BUY ? exchange.Buy : exchange.Sell;
    var pfnDir = order.Type == ORDER_TYPE_BUY ? "buy" : "sell";
    var coverPfn = order.Type == ORDER_TYPE_BUY ? exchange.Sell : exchange.Buy;
    var coverPfnDir = order.Type == ORDER_TYPE_BUY ? "closebuy" : "closesell";

    // 等待开仓的订单
    if (order.Status == STATE_WAIT_OPEN) {
      var diff = _N(
        order.Type == ORDER_TYPE_BUY
          ? ticker.Buy - order.Price
          : order.Price - ticker.Sell
      );
      // 不主动成交
      if (diff < 0) {
        return;
      }
      exchange.SetDirection(pfnDir);
      realId = pfn(
        order.Price,
        order.Amount,
        order.Extra +
          "(距离: " +
          diff +
          (order.Type == ORDER_TYPE_BUY
            ? " 买一: " + ticker.Buy
            : " 卖一: " + ticker.Sell) +
          ")" +
          " VID:" +
          order.VID
      );
      if (realId != null) {
        order.OpenId = realId;
        order.Status = STATE_WAIT_COVER;
      } else {
        return;
      }
      return;
    }

    // 等待平仓的订单
    if (order.Status == STATE_WAIT_COVER) {
      found = hasOrder(exchangeOrders, order.OpenId);
      if (!found) {
        exchange.SetDirection(coverPfnDir);
        realId = coverPfn(
          order.CoverPrice,
          order.Amount,
          order.Extra + " 平仓 价格:" + order.CoverPrice + " VID:" + order.VID
        );
        if (realId != null) {
          order.CoverId = realId;
          order.Status = STATE_WAIT_CLOSE;
        } else {
          return;
        }
      }
      return;
    }

    //  等待完结的订单
    if (order.Status == STATE_WAIT_CLOSE) {
      found = hasOrder(exchangeOrders, order.CoverId);
      if (!found) {
        G.Log("close order:" + order.CoverId);
        order.Status = STATE_END_CLOSE;
      }
      return;
    }
  };

  // 遍历所有各自尝试转换状态机
  this.Poll = function(ticker, orders, ext) {
    var married = 0;
    var deleteBooks = new Object();
    for (orderId in orderBooks) {
      var order = orderBooks[orderId];
      this.PollOne(order, ticker, orders, ext);
      //record order wait to convert to history
      if (order.Status == STATE_END_CLOSE) {
        deleteBooks[orderId] = orderId;
      }
    }
    for (orderId in deleteBooks) {
      hisBooks[orderId] = deleteBooks[orderId];
      hisBooksLen++;
      married++;
      delete orderBooks[orderId];
      orderBooksLen--;
    }
    return married;
  };
}

// getHoldPosition get the position hold
function getHoldPosition(positions, dir) {
  var len = positions.length;
  if (len == 0) {
    return null;
  }
  for (var i = 0; i < len; i += 1) {
    if (dir == 0 && positions[i].TradeType == "buy") {
      return positions[i];
    }
    if (dir == 1 && positions[i].TradeType == "sell") {
      return positions[i];
    }
  }
  return null;
}

// 持仓定时器
var holdTimer = new TradeRobot("hold");
//var gridTrader = new GridTrader();
var fishCheckTimer = new TradeRobot("check");
var firstPrice = -1;

// 动态再平衡, 注意不会平仓调hold部分的仓位
function resetAccount(dir) {
  G.Log("平衡账户mode:", dir);
  cancelPending(dir);
  while (true) {
    blockGetInfo(onOrders, onPosition);
    var orders = order2DirOrder(globalInfo.orders, BuyFirst ? 0 : 1);
    var positions = globalInfo.positions;
    G.Log("positions:", positions);
    G.Log("orders:", orders);
    var pos = BuyFirst
      ? getHoldPosition(positions, 0)
      : getHoldPosition(positions, 1);
    G.Log("pos:", pos);
    var holdAmount = pos == null ? 0 : pos.Amount;
    if (holdAmount == 0 && orders.length == 0) {
      break;
    }
    if (holdAmount > 0) {
      var leftAmount = holdAmount;
      if (BuyFirst) {
        G.Log("平多仓", leftAmount);
        exchange.SetDirection("closebuy");
        var closeId = exchange.Sell(-1, leftAmount, "平多仓");
      } else {
        G.Log("平空仓", leftAmount);
        exchange.SetDirection("closesell");
        var closeId = exchange.Buy(-1, leftAmount, "平空仓");
      }
    }
    G.Sleep(Interval);
    if (orders.length != 0) {
      cancelPending(dir);
    }
  }
  G.Log("平衡完成");
}

function positionDiff(positions) {
  var buyPos = getHoldPosition(positions, 0);
  var sellPos = getHoldPosition(positions, 1);
  var buyPosAmount = buyPos == null ? 0 : buyPos.Amount;
  var sellPosAmount = sellPos == null ? 0 : sellPos.Amount;
  var diffAmount = buyPosAmount - sellPosAmount;
  return diffAmount;
}

function exit() {
  G.Log("策略开始停止");
  resetAccount(BuyFirst ? 0 : 1);
  G.Log("策略成功停止");
}

// return 0-continue 1-fish again 2-exit app 3-continue
function fishingCheck(orgAccount, gridTrader, position, totpositions, ticker) {
  var msg = "";
  var isHold = false;

  var holdAmount = 0;
  if (position != null) {
    isHold = true;
    holdAmount = position.Amount;
  }

  if (fishCheckTimer.Timeout()) {
    fishCheckTimer.SetInterval(FishCheckTime);

    if (isHold) {
      var profitRate = position2Rate(position, ticker.Last);
      msg +=
        "持仓: " +
        position.Amount +
        " 持仓均价: " +
        _N(position.Price, Precision) +
        " 浮动盈亏量: " +
        String(position.Profit) +
        " 浮动盈亏率:" +
        String(_N(profitRate, Precision)) +
        "%";

      if (StopLoss > 0 && profitRate + StopLoss < 0) {
        G.Log("当前浮动盈亏", profitRate, "开始止损");
        resetAccount(BuyFirst ? 0 : 1);
        if (StopLossAfterMode === 0) {
          return 2;
        }
        return 1;
      }

      if (StopWin > 0 && profitRate - StopWin > 0) {
        G.Log("当前浮动盈亏", profitRate, "开始止盈");
        resetAccount(BuyFirst ? 0 : 1);
        return 1;
      }
    } else {
      msg += "空仓";
    }

    var holdAmount = position === null ? 0 : position.Amount;
    if (holdAmount != 0) {
      holdTimer.SetInterval(HoldTime);
    }

    if (ticker.Last < LowBox || ticker.Last > HighBox) {
      G.Log(
        "当前价格超过箱体 last:",
        ticker.Last,
        " lowbox:",
        LowBox,
        " HighBox:",
        HighBox
      );
      return 3;
    }

    if (AutoMove) {
      var refish = false;
      if (holdTimer.Timeout()) {
        G.Log("空仓过久未变化, 开始移动网格");
        refish = 1;
      }

      if (refish) {
        resetAccount(BuyFirst ? 0 : 1);
        return 1;
      }
    }

    msg += "\n";
    var account = globalInfo.account;
    var oldStock = orgAccount.TotBalance + 0.00000001;
    //var currStock = account2balance(account, Coin);
    //var diffStock = currStock - oldStock;
    var diffAmount = positionDiff(totpositions);
    msg += "总原货币量:" + String(_N(oldStock, 10)) + "\n";
    //msg += "总现货币量:" + String(_N(currStock, 10)) + "\n";
    //msg += "总盈亏量:" + String(_N(diffStock, 10)) + "\n";
    msg += "箱体上沿:" + String(_N(HighBox)) + "\n";
    msg += "箱体下沿:" + String(_N(LowBox)) + "\n";
    msg += "止损百分比:" + String(_N(StopLoss)) + "%\n";
    msg += "止盈百分比:" + String(_N(StopWin)) + "%\n";
    msg += "仓位上沿:" + String(_N(HighPosition)) + "\n";
    msg += "当前价格:" + String(_N(ticker.Last)) + "\n";
    msg += "totpositions:" + JSON.stringify(totpositions) + "\n";
    //msg += "当前仓位差:" + String(_N(diffAmount)) + "\n";
    msg +=
      "已撮合单数:" +
      String(marryTot * (BuyFirst ? AmountOnce : AmountOnce)) +
      "\n";
    msg +=
      "hold定时器剩余时间:" + String(holdTimer.TimeLeft() / 1000.0) + "s\n";
    //G.LogStatus(msg);
    G.Log(msg);
    //gridTrader.Debug();
  }

  // 检查后发现持仓达到最大仓位后不需要继续追加持仓
  if (isHold && holdAmount >= HighPosition) {
    return 3;
  }

  var orderLen = gridTrader.BooksLen();
  //G.Log("order Len:", JSON.stringify(orderLen));
  if (orderLen.wait_open + orderLen.wait_cover > 0) {
    return 3;
  }

  return 0;
}

function nextGridPrice(ticker, lastPrice) {
  var nextPrice = lastPrice;
  var cycle = true;
  while (cycle) {
    nextPrice = _N(
      BuyFirst ? nextPrice - GridOffset : nextPrice + GridOffset,
      Precision
    );
    if (nextPrice < LowBox || nextPrice > HighBox) {
      return -1;
    }
    if (BuyFirst) {
      if (nextPrice < ticker.Buy) {
        break;
      }
    } else {
      if (nextPrice > ticker.Sell) {
        break;
      }
    }
  }
  return nextPrice;
}

function totProfit(positions, ticker) {
  var buyPos = getHoldPosition(positions, 0);
  var sellPos = getHoldPosition(positions, 1);
  var buyProfitRate = 0;
  var sellProfitRate = 0;
  if (buyPos != null) {
    buyProfitRate = position2Rate(buyPos, ticker.Last);
  }

  if (sellPos != null) {
    sellProfitRate = position2Rate(sellPos, ticker.Last);
  }
  return buyProfitRate + sellProfitRate;
}

function fishing(orgAccount, fishCount) {
  var gridTrader = new GridTrader();
  gridTrader.SetProfitPrice(ProfitPrice);
  holdTimer.SetInterval(HoldTime);

  fishCheckTimer.SetInterval(FishCheckTime);
  var lastPrice = -1;
  var cycle = true;
  while (cycle) {
    blockGetInfo(onOrders, onTicker, onPosition, onAccount);
    var ticker = globalInfo.ticker;
    var orders = globalInfo.orders;
    var account = globalInfo.account;
    var positions = globalInfo.positions;
    //gridTrader.Debug();
    var ext = new Object();
    var pos = BuyFirst
      ? getHoldPosition(positions, 0)
      : getHoldPosition(positions, 1);
    var holdAmount = pos === null ? 0 : pos.Amount;

    var checkFlag = fishingCheck(
      orgAccount,
      gridTrader,
      pos,
      positions,
      ticker
    );

    //G.Log("checkflag is:", checkFlag);
    if (checkFlag == 0) {
    }

    if (checkFlag == 1) {
      return true;
    }

    if (checkFlag == 2) {
      return false;
    }

    if (checkFlag == 3) {
      marryTot += gridTrader.Poll(ticker, orders, ext);
      G.Sleep(Interval);
      continue;
    }

    var nextPrice = -1;
    if (lastPrice < 0) {
      firstPrice = BuyFirst
        ? _N(ticker.Buy - OpenProtect, Precision)
        : _N(ticker.Sell + OpenProtect, Precision);
      nextPrice = firstPrice;
      G.Log(
        "计算下单位置:" +
          "ticker.Buy:" +
          String(ticker.Buy) +
          " ticker.Sell" +
          String(ticker.Sell) +
          " nextPrice:" +
          String(nextPrice)
      );
      // need to open new one
    } else {
      nextPrice = nextGridPrice(ticker, lastPrice);
      // out of box
      if (nextPrice < 0) {
        G.Log("尝试挂单位置超过箱体，放弃挂单");
        marryTot += gridTrader.Poll(ticker, orders, ext);
        G.Sleep(Interval);
        continue;
      }

      var posdiff = positionDiff(positions);
      if (!BuyFirst) {
        posdiff = 0 - posdiff;
      }
      if (PositionDiff > 0 && posdiff > PositionDiff) {
        G.Log("仓位差过大，放弃挂单 " + String(posdiff));
        marryTot += gridTrader.Poll(ticker, orders, ext);
        G.Sleep(Interval);
        continue;
      } else {
        G.Log("仓位差合理，继续挂单 " + String(posdiff));
      }
    }
    /*
    var needStocks = Order2Cost(
      nextPrice,
      BuyFirst ? AmountOnce : AmountOnce,
      ticker.Last
    );
    G.Log("下单需要stock:", needStocks);
    if (needStocks >= account.Stocks * MarginLevel) {
      G.Log("需要的stock不足:", needStocks);
      marryTot += gridTrader.Poll(ticker, orders, ext);
      G.Sleep(Interval);
      continue;
    }
    */

    if (BuyFirst) {
      gridTrader.Buy(nextPrice, AmountOnce, "");
    } else {
      gridTrader.Sell(nextPrice, AmountOnce, "");
    }

    marryTot += gridTrader.Poll(ticker, orders, ext);
    lastPrice = nextPrice;
    G.Sleep(Interval);
  }
  //return true;
}

function IsParameterInvalid() {
  if (AmountOnce <= 0) {
    return "AmountOnce invalid:" + AmountOnce.toString();
  }
  if (AmountOnce <= 0) {
    return "AmountOnce invalid:" + AmountOnce.toString();
  }

  if (GridOffset <= 0) {
    return "GridOffset invalid:" + GridOffset.toString();
  }

  if (ProfitPrice <= 0) {
    return "ProfitPrice invalid:" + ProfitPrice.toString();
  }

  if (OpenProtect < 0) {
    return "OpenProtect invalid:" + OpenProtect.toString();
  }

  if (Precision < 0) {
    return "Precision invalid:" + Precision.toString();
  }
  if (HoldTime <= 0) {
    return "HoldTime invalid:" + HoldTime.toString();
  }

  if (HighPosition <= 0) {
    return "HighPosition invalid:" + HighPosition.toString();
  }

  if (LowBox < 0) {
    return "LowBox invalid:" + LowBox.toString();
  }

  if (HighBox <= 0) {
    return "HighBox invalid:" + HighBox.toString();
  }
  if (LowBox > HighBox) {
    return "box range invalid:" + LowBox.toString() + "->" + HighBox.toString();
  }
  if (-1 == ContractVec.indexOf(ContractType)) {
    return "contractType not support:" + ContractType;
  }

  if (-1 == CoinVec.indexOf(Coin)) {
    return "Coin not support:" + Coin;
  }
  if (MarginLevel < 1) {
    return "marginLevel not support:" + MarginLevel.toString();
  }
}

function main() {
  var invalid = IsParameterInvalid();
  if (invalid != null) {
    G.Log(invalid);
    return 0;
  }
  exchange.SetStockType(Coin);
  exchange.SetContractType(ContractType); // 设置合约
  exchange.SetMarginLevel(MarginLevel); // 设置杠杆
  blockGetInfo(onOrders, onAccount, onPosition, onTicker);
  var orgAccount = JSON.parse(JSON.stringify(globalInfo.account));
  var fishCount = 1;
  var totBalance = account2balance(orgAccount, Coin);
  orgAccount.TotBalance = totBalance;
  G.Log(
    "Stocks:",
    orgAccount.Stocks,
    "FrozenStocks:",
    orgAccount.FrozenStocks,
    "TotBalance:",
    orgAccount.TotBalance
  );
  var position = BuyFirst
    ? getHoldPosition(globalInfo.positions, 0)
    : getHoldPosition(globalInfo.positions, 1);
  if (position != null) {
    G.Log("仓位 amount:", position.Amount);
    G.Log("仓位 price:", position.Price);
    G.Log("仓位 profit:", position.Profit);
    G.Log("仓位 rate:", position2Rate(position, globalInfo.ticker.Last));
  }
  var orders = orderClean(globalInfo.orders);
  G.Log("all orders:", JSON.stringify(orders));
  G.Log("buy orders:", JSON.stringify(order2DirOrder(orders, 0)));
  G.Log("sell orders:", JSON.stringify(order2DirOrder(orders, 1)));
  var cycle = true;
  while (cycle) {
    if (!fishing(orgAccount, fishCount)) {
      break;
    }
    fishCount++;
    G.Log("第", fishCount, "次重新撒网...");
    G.Sleep(Interval);
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

//console.log(exchange.Buy(9000, 100, "buy"))
//console.log(exchange.Sell(9000, 100, "sell"))
/*
grid = new GridTrader();
orders = [];
var buyPrice = 9000;
var sellPrice = 9001;
var amount = 100;
var ticker = Object();
var ext = new Object();
ext.hold = 0;
ticker.Buy = 9000;
ticker.Sell = 8090;
ticker.Last = 8001;
grid.SetProfitPrice(50);
G.Log("profit price is :" + grid.GetProfitPrice());

G.Log("init process");
grid.Sell(sellPrice, amount, "Sell");
grid.Sell(sellPrice, amount, "Sell");
grid.Sell(sellPrice, amount, "Sell");
grid.Sell(sellPrice, amount, "Sell");
grid.Debug();
G.Log("len function");
G.Log(grid.BooksLen());

G.Log("open process");
grid.Poll(ticker, orders, 10, ext);
grid.Debug();

G.Log("cover process");
orders = [
  {
    Id: 5
  },
  {
    Id: 6
  }
];
grid.Poll(ticker, orders, 10, ext);
G.Log("Orders List");
G.Log(grid.GetOrders());
G.Log("History List");
G.Log(grid.GetHistoryOrders());
G.Log("len function");
G.Log(grid.BooksLen());

orders = [
  {
    Id: 5
  },
  {
    Id: 6
  }
];

grid.Poll(ticker, orders, 10, ext);
G.Log("Orders List");
G.Log(grid.GetOrders());
G.Log("History List");
G.Log(grid.GetHistoryOrders());

HighBox = 9900;
LowBox = 7001;
ticker.Buy = 7050;
G.Log(nextGridPrice(ticker, 8000));
*/
