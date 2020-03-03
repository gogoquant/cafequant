// coin type
Coin = "BTC";
// 滑点
Slide = 0.1;
// 箱体上沿
HighBox = 10000;
// 箱体下沿
LowBox = 8000;
// 网格方向
BuyFirst = 1;
// 计划持仓量
MaxPosition = 20;
// 网格价格距离
GridOffset = 50;
// 价格精度
Precision = 1;
// 开仓保护价差
OpenProtect = 5;
// 买单数量
BAmountOnce = 2;
// 卖单数量
SAmountOnce = 2;
// 是否止损
EnableStopLoss = 1;
// 止损模式
StopLossMode = 0;
// 是否止盈
EnableStopWin = 1;
// 止损盈亏损率
StopLoss = 60;
// 止盈率
StopWin = 15;
// 最小量
MinStock = 0.1;
// 是否自动移动价格
AutoMove = 1;
// 仓位当前价格最大价差
MaxDistance = 300;
// 最大空仓时间
HoldTime = 1000 * 60 * 5;
// 收网检测周期
FishCheckTime = 1000 * 60 * 1;
// 最小周期
Interval = 1000 * 10 * 1;
// 盈利滑动
ProfitPrice = 30;
// 合约
ContractType = "quarter";
// 杠杆
MarginLevel = 20;
// 合约列表
ContractVec = ["this_week", "next_week", "quarter"];

// local
var globalInfo = {};

var STATE_WAIT_OPEN = "wait_open";
var STATE_WAIT_COVER = "wait_cover";
var STATE_WAIT_CLOSE = "wait_close";
var STATE_END_CLOSE = "end_close";
var ORDER_TYPE_BUY = 0;
var ORDER_TYPE_SELL = 1;

// 持仓定时器
var holdTimer = new TradeRobot("hold");
var fishCheckTimer = new TradeRobot("check");
var firstPrice = -1;
var orgAccount = new Object();

/*
 * Only used for test
 */

/*
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
*/

// @Todo
// 1. 空方向测试
// 2. 多空双向使用一套程序管理两个仓位?
// 3. 平仓挂单和开仓挂单本质无强关联，是否分批平仓单策略只通过当前仓为平均价格来推算挂平仓单的方式?
// 4. 考虑使用其他算法，调整网格的离散度?

// ArrayQueue 队列
ArrayQueue = function() {
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
};

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
  return ((position.Profit * price) / position.Amount) * MarginLevel;
}

function initInfo() {
  globalInfo = {};
}

function onDepth() {
  globalInfo.depth = exchange.GetDepth();
}

function onPosition() {
  globalInfo.positions = exchange.GetPosition();
}

function onTicker() {
  globalInfo.ticker = exchange.GetTicker();
}

function onOrders() {
  globalInfo.orders = exchange.GetOrders();
}

function onAccount() {
  globalInfo.account = exchange.GetAccount();
}

function checkInfo() {
  var obj = globalInfo;
  for (var key in obj) {
    console.log("get key :" + key);
    if (!ValidItem(obj[key])) {
      Log(key + ": get fail from exchange");
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
    blockGetInfo(onOrders);
    var orders = globalInfo.orders;
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
    blockGetInfo(onOrders);
    var orders = globalInfo.orders;
    order = foundOrder(orders, Id);
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
  data = account.Info.data;
  for (j = 0, len = data.length; j < len; j++) {
    info = data[j];
    if (info.symbol == symbol) {
      return info.margin_balance;
    }
  }
  return null;
}

function GridTrader() {
  var vId = 0;
  var orderBooks = new Object();
  var hisBooks = new Object();
  var orderBooksLen = 0;
  var hisBooksLen = 0;
  var openLen = 0;
  var profitPrice = -1;
  var lastOrderPrice = -1;

  this.SetProfitPrice = function(price) {
    profitPrice = price;
  };

  this.GetProfitPrice = function() {
    return profitPrice;
  };

  this.Debug = function() {
    Log("open len:", openLen);
    Log("Orders List:");
    for (orderId in orderBooks) {
      Log(orderBooks[orderId]);
    }
    Log("HisOrders List:");
    for (orderId in hisBooks) {
      Log(hisBooks[orderId]);
    }
    return;
  };

  this.OpenLen = function() {
    return openLen;
  };

  this.Buy = function(price, amount, extra) {
    Log("Buy price: " + price + ", amount: " + amount + ", extra:" + extra);
    //return

    if (typeof extra === "undefined") {
      extra = "";
    } else {
      extra = valuesToString(arguments, 2);
    }
    lastOrderPrice = price;
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
    openLen++;
    orderBooksLen++;
    return orderId;
  };
  this.Sell = function(price, amount, extra) {
    Log("Sell price: " + price + ", amount: " + amount + ", extra:" + extra);
    //return

    if (typeof extra === "undefined") {
      extra = "";
    } else {
      extra = valuesToString(arguments, 2);
    }
    lastOrderPrice = price;
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
    openLen++;
    orderBooksLen++;
    return orderId;
  };

  this.GetOrders = function() {
    return orderBooks;
  };

  this.GetHistoryOrders = function() {
    return hisBooks;
  };

  this.GetLastOrderPrice = function() {
    return this.lastOrderPrice;
  };

  this.GetOrder = function(orderId) {
    if (typeof orderId === "number") {
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
  this.PollOne = function(order, ticker, exchangeOrders) {
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
      var realId = pfn(
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
      }
      return;
    }

    // 等待平仓的订单
    if (order.Status == STATE_WAIT_COVER) {
      var found = hasOrder(exchangeOrders, order.OpenId);
      if (!found) {
        exchange.SetDirection(coverPfnDir);
        var realId = coverPfn(
          order.CoverPrice,
          order.Amount,
          order.Extra + " 平仓 价格:" + order.CoverPrice + " VID:" + order.VID
        );
        if (realId != null) {
          order.CoverId = realId;
          order.Status = STATE_WAIT_CLOSE;
          openLen--;
        }
      }
      return;
    }

    //  等待完结的订单
    if (order.Status == STATE_WAIT_CLOSE) {
      var found = hasOrder(exchangeOrders, order.CoverId);
      if (!found) {
        Log("close order:" + order.CoverId);
        order.Status = STATE_END_CLOSE;
      }
      return;
    }
  };

  // 遍历所有各自尝试转换状态机
  this.Poll = function(ticker, orders) {
    var deleteBooks = new Object();
    for (orderId in orderBooks) {
      var order = orderBooks[orderId];
      this.PollOne(order, ticker, orders);
      //record order wait to convert to history
      if (order.Status == STATE_END_CLOSE) {
        deleteBooks[orderId] = orderId;
      }
    }
    for (orderId in deleteBooks) {
      hisBooks[orderId] = deleteBooks[orderId];
      hisBooksLen++;
      delete orderBooks[orderId];
      orderBooksLen--;
    }
  };
}


function getHoldPosition(positions, dir) {
  var len = positions.length;
  if (len == 0) {
    return null;
  }
  for (var i = 0; i < len; i += 1) {
    if (positions[i].Type == dir) {
      return positions[i];
    }
  }
  return null;
}

// 动态再平衡
function balanceAccount() {
  Log("平衡账户");
  cancelPending();
  while (true) {
    blockGetInfo(onOrders, onPosition);
    var orders = globalInfo.orders;
    var positions = globalInfo.positions;
    var pos = BuyFirst
      ? getHoldPosition(positions, 0)
      : getHoldPosition(positions, 1);

    if (pos == null && orders.length == 0) {
      break;
    }

    if (pos != null) {
      var leftAmount = pos.Amount;
      if (pos.Type == 0) {
        //平多仓，采用盘口吃单，会损失手续费，可改为盘口挂单，会增加持仓风险。
        Log("平多仓", leftAmount);
        exchange.SetDirection("closebuy");
        var closeId = exchange.Sell(-1, leftAmount, "平多仓");
      } else {
        Log("平空仓", leftAmount);
        exchange.SetDirection("closesell");
        var closeId = exchange.Buy(-1, leftAmount, "平空仓");
      }
    }
    Sleep(Interval);
    if (orders.length != 0) {
      cancelPending();
    }
  }
  Log("平衡完成");
}

function onexit() {
  cancelPending();
  Log("策略成功停止");
  blockGetInfo(onAccount);
  var account = globalInfo.account;
  Log(account);
}

// return 0-continue 1-fish again 2-exit app 3-continue
function fishingCheck(orgAccount, grid) {
  var msg = "";
  var ticker = globalInfo.ticker;
  var positions = globalInfo.positions;
  var position = BuyFirst
    ? getHoldPosition(globalInfo.positions, 0)
    : getHoldPosition(globalInfo.positions, 1);
  var isHold = false;
  var holdAmount = 0;
  if (position != null) {
    isHold = true;
    holdAmount = position.Amount;
  }

  if (fishCheckTimer.Timeout()) {
    fishCheckTimer.SetInterval(FishCheckTime);

    if (isHold) {
      holdTimer.SetInterval(HoldTime);
    }

    if (isHold) {
      var profitRate = position2Rate(position, globalInfo.ticker.Last);
      Log("仓位 profit:", profitRate);
      msg +=
        "持仓: " +
        position.Amount +
        " 持仓均价: " +
        _N(position.Price, Precision) +
        " 浮动盈亏量: " +
        _N(position.Profit) +
        " 浮动盈亏率" +
        _N(profitRate, Precision);

      if (EnableStopLoss && profitRate + StopLoss < 0) {
        Log("当前浮动盈亏", profitRate, "开始止损");
        balanceAccount();
        if (StopLossMode === 0) {
          return 2;
        }
        return 1;
      }

      if (EnableStopWin && profitRate - StopWin > 0) {
        Log("当前浮动盈亏", profitRate, "开始止盈");
        balanceAccount();
        return 1;
      }
    } else {
      msg += "空仓";
    }
    var distance = 0;
    if (AutoMove) {
      if (BuyFirst) {
        distance = ticker.Last - firstPrice;
      } else {
        distance = firstPrice - ticker.Last;
      }
      var refish = false;
      if (!isHold && holdTimer.Timeout()) {
        Log("空仓过久, 开始移动网格");
        refish = 1;
      }
      if (distance > MaxDistance) {
        Log(
          "价格超出网格区间过多, 开始移动网格, 当前距离: ",
          _N(distance, Precision),
          "当前价格:",
          ticker.Last
        );
        //refish = 1;
      }
      if (refish) {
        balanceAccount();
        return 1;
      }
    }
    msg += "\n";
    var account = globalInfo.account;
    oldStock = orgAccount.totBalance + 0.00000001;
    currStock = account2balance(account, Coin);
    diffStock = currStock - oldStock;
    msg += "总原货币量:" + String(_N(oldStock, 6)) + "\n";
    msg += "总现货币量:" + String(_N(currStock, 6)) + "\n";
    msg += "总盈亏量:" + String(_N(diffStock)) + "\n";
    msg +=
      "总盈亏率" + String(_N(((diffStock * 1.0) / oldStock) * 100, 6)) + "%\n";
    LogStatus(msg);
    grid.Debug();
  }

  // 检查后发现持仓达到期望后不需要继续追加持仓
  if (isHold && holdAmount > 0 && holdAmount >= MaxPosition) {
    return 3;
  }
  return 0;
}

function nextGridPrice(ticker, lastPrice) {
  var nextPrice = lastPrice;
  while (true) {
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

function fishing(orgAccount, fishCount) {
  var gridTrader = new GridTrader();
  gridTrader.SetProfitPrice(ProfitPrice);

  holdTimer.SetInterval(HoldTime);
  fishCheckTimer.SetInterval(FishCheckTime);
  var lastPrice = -1;
  while (true) {
    blockGetInfo(onOrders, onTicker, onPosition, onAccount);
    var isHold = false;
    var ticker = globalInfo.ticker;
    var orders = globalInfo.orders;
    var account = globalInfo.account;

    //超出网格则停止机器人
    if (ticker.Last < LowBox || ticker.Last > HighBox) {
      Log(
        "当前价格超过箱体 last:",
        ticker.Last,
        " lowbox:",
        LowBox,
        " HighBox:",
        HighBox
      );
      return false;
    }

    var checkFlag = fishingCheck(orgAccount, gridTrader);

    Log("checkflag is:", checkFlag);

    if (checkFlag == 0) {
    }

    if (checkFlag == 1) {
      return true;
    }

    if (checkFlag == 2) {
      return false;
    }

    if (checkFlag == 3) {
      gridTrader.Poll(ticker, orders);
      Sleep(Interval);
      continue;
    }

    if (gridTrader.OpenLen() > 0) {
      gridTrader.Poll(ticker, orders);
      Sleep(Interval);
      continue;
    }

    var nextPrice = -1;
    //lastPrice = gridTrader.GetLastOrderPrice()
    if (lastPrice < 0) {
      firstPrice = BuyFirst
        ? _N(ticker.Buy - OpenProtect, Precision)
        : _N(ticker.Sell + OpenProtect, Precision);
      nextPrice = firstPrice;
      Log("ticker.Buy:", ticker.Buy, "ticker.Sell:", ticker.Sell);
      Log("初始fish nextPrice:", nextPrice);
      // need to open new one
    } else {
      nextPrice = nextGridPrice(ticker, lastPrice);
      // out of box
      if (nextPrice < 0) {
        Log("尝试挂单位置超过箱体，放弃挂单");
        gridTrader.Poll(ticker, orders);
        Sleep(Interval);
        continue;
      }
    }
    var needStocks = Order2Cost(
      nextPrice,
      BuyFirst ? BAmountOnce : SAmountOnce,
      ticker.Last
    );
    Log("下单需要stock:", needStocks);
    if (needStocks >= account.Stocks * MarginLevel) {
      Log("需要的stock不足:", needStocks);
      gridTrader.Poll(ticker, orders);
      Sleep(Interval);
      continue;
    }

    if (BuyFirst) {
      gridTrader.Buy(nextPrice, BAmountOnce, "");
    } else {
      gridTrader.Sell(nextPrice, SAmountOnce, "");
    }

    gridTrader.Poll(ticker, orders);

    lastPrice = nextPrice;
    Sleep(Interval);
  }
  return true;
}

function main() {
  exchange.SetContractType(ContractType); // 设置合约
  exchange.SetMarginLevel(MarginLevel); // 设置杠杆
  blockGetInfo(onAccount, onPosition, onTicker);
  var orgAccount = Object.assign({}, globalInfo.account);
  var fishCount = 1;
  var totBalance = account2balance(orgAccount, Coin);
  orgAccount.totBalance = totBalance;
  Log(
    "Stocks:",
    orgAccount.Stocks,
    "FrozenStocks:",
    orgAccount.FrozenStocks,
    "totBalance:",
    orgAccount.totBalance
  );
  var position = BuyFirst
    ? getHoldPosition(globalInfo.positions, 0)
    : getHoldPosition(globalInfo.positions, 1);
  if (position != null) {
    Log("仓位 amount:", position.Amount);
    Log("仓位 price:", position.Price);
    Log("仓位 profit:", position.Profit);
    Log("仓位 rate:", position2Rate(position, globalInfo.ticker.Last));
  }

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

/*
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
Log("Orders List")
Log(grid.GetOrders())
Log("History List")
Log(grid.GetHistoryOrders())

Log("open process")
grid.Poll(ticker, orders, 10)
Log("Orders List")
Log(grid.GetOrders())
Log("History List")
Log(grid.GetHistoryOrders())

Log("cover process")
orders = [{
        Id: 5,
    },
    {
        Id: 6
    }
]
grid.Poll(ticker, orders, 10)
Log("Orders List")
Log(grid.GetOrders())
Log("History List")
Log(grid.GetHistoryOrders())

orders = [{
        Id: 5,
    },
    {
        Id: 6
    }
]

grid.Poll(ticker, orders, 10)
Log("Orders List")
Log(grid.GetOrders())
Log("History List")
Log(grid.GetHistoryOrders())

HighBox = 9900
LowBox = 7001
ticker.Buy = 7050
Log(nextGridPrice(ticker, 8000))
*/

