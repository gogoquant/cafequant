// 对冲双网格
//
// Coin type
var Coin = 'BTC/USD';
// 箱体上沿
var HighBox = 14000;
// 箱体下沿
var LowBox = 10000;
// 计划持仓量
var HighPosition = 20;
// 网格价格距离
var GridOffset = 50;
// 价格精度
var Precision = 1;
// 开仓保护价差
var OpenProtect = 1;
// 对冲仓位差
var PositionDiff = 4;
// 数量
var GridAmount = 2;
// 止损后模式
var StopLossAfterMode = 0;
// 止损盈亏损率
var StopLoss = -5.0;
// 止盈率
var StopWin = 25.0;
// 是否自动移动价格
var AutoMove = true;
// 最大空仓时间
var HoldTime = 1000 * 60 * 2;
// 收网检测周期
var FishCheckTime = 1000 * 20 * 1;
// 最小周期
var Interval = 1000 * 10 * 1;
// 盈利滑动
var ProfitPrice = 50;
// 合约
var ContractType = 'quarter';
// //
var ContractVal = 100;
// 杠杆
var MarginLevel = 20;
// 合约列表
//var ContractVec = ['this_week', 'next_week', 'quarter'];
// 货币支持类型
var CoinVec = ['BTC/USD'];

// local
var GlobalInfo = {};

var TRADE_TYPE_LONG = 'buy';
var TRADE_TYPE_SHORT = 'sell';
var TRADE_TYPE_LONGCLOSE = 'closebuy';
var TRADE_TYPE_SHORTCLOSE = 'closesell';
var STATE_WAIT_OPEN = 'wait_open';
var STATE_WAIT_COVER = 'wait_cover';
var STATE_WAIT_CLOSE = 'wait_close';
var STATE_END_CLOSE = 'end_close';
var ORDER_TYPE_BUY = 0;
var ORDER_TYPE_SELL = 1;

var ExchangeReal = function() {
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

function floatReverse(num, pre) {
  var str1 = '';
  var str = '';
  var res = -1;

  str1 = String(num); //将类型转化为字符串类型
  if (str1.indexOf('.') < 0) {
    return num;
  }
  str = str1.substring(0, str1.indexOf('.') + pre + 1); //截取字符串
  res = Number(str); //转化为number类型
  return res;
}

// TradeTimer 交易者,用于交易配置以及延时处理等
function TradeTimer(name) {
  this.name = name;
  this.price = 0; //下单价
  this.amount = 0; //下单量
  this.leftTime = 0; //用于显示
  this.interval = -1;
  this.nextTime = new Date().valueOf();
  this.run = 1;
}

TradeTimer.prototype.Name = function() {
  return this.name;
};

TradeTimer.prototype.SetInterval = function(num) {
  this.interval = num;
  this.leftTime = 0 - num;
  this.nextTime = new Date().valueOf() + num;
};

TradeTimer.prototype.AddInterval = function(num) {
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

TradeTimer.prototype.Timeout = function() {
  this.leftTime = this.nextTime - new Date().valueOf();
  return this.leftTime >= 0 ? 0 : 1;
};

TradeTimer.prototype.TimeLeft = function() {
  return this.nextTime - new Date().valueOf();
};

TradeTimer.prototype.Left = function() {
  this.leftTime = new Date().valueOf() - this.nextTime;
  return this.leftTime;
};

TradeTimer.prototype.Interval = function() {
  return this.interval;
};

function emptyItem(val) {
  if (typeof val === 'undefined' || val === null) {
    return false;
  }
  return true;
}

function position2Val(position) {
  return position.ProfitRate * ContractVal;
}

// order2DirOrder normal filter by dir
function order2DirOrder(orders, dirs) {
  var i;
  var ordervec = [];
  var dir;

  if (dirs === '') {
    return orders;
  }
  for (i = 0; i < orders.length; i += 1) {
    for (dir in dirs) {
      if (orders[i].TradeType === dir) {
        ordervec.push(orders[i]);
      }
    }
  }
  return ordervec;
}

function initInfo() {
  GlobalInfo = {};
}

//function onDepth() {
//  GlobalInfo.depth = exchange.GetDepth();
//}

function onPosition() {
  GlobalInfo.positions = E.GetPosition();
}

function onTicker() {
  GlobalInfo.ticker = E.GetTicker();
}

function onOrders() {
  GlobalInfo.orders = E.GetOrders();
}

function onAccount() {
  GlobalInfo.account = E.GetAccount();
}

function checkInfo() {
  var obj = GlobalInfo;
  var key = '';

  for (key in obj) {
    if (!emptyItem(obj[key])) {
      G.Log('get fail from exchange:' + key);
      return false;
    }
  }
  return true;
}

function blockGetInfo() {
  var i;
  var fn;

  initInfo();
  do {
    for (i = 0; i < arguments.length; i++) {
      fn = arguments[i];

      fn();
    }
  } while (!checkInfo());
}

function hasOrder(orders, orderId, dirs) {
  var i = 0;
  var tmporders = order2DirOrder(orders, dirs);

  for (i = 0; i < tmporders.length; i++) {
    if (orders[i].Id === orderId) {
      return true;
    }
  }
  return false;
}

function foundOrder(orders, orderId, dirs) {
  var tmporders = order2DirOrder(orders, dirs);
  var i = 0;

  for (i = 0; i < tmporders.length; i++) {
    if (orders[i].Id === orderId) {
      return orders[i];
    }
  }
  return null;
}

// 阻塞关闭all订单
function cancelPending(dir) {
  var ret = false;
  var cycle = true;
  var orders;
  var j;

  while (cycle) {
    if (ret) {
      G.Sleep(Interval);
    }
    blockGetInfo(onOrders);
    orders = order2DirOrder(GlobalInfo.orders, dir);

    if (orders.length === 0) {
      break;
    }
    for (j = 0; j < orders.length; j++) {
      exchange.CancelOrder(orders[j].Id);
      ret = true;
    }
  }
  return ret;
}

// 阻塞关闭一个订单
function cancelOnePending(Id, dirs) {
  var ret = false;
  var cycle = true;
  var orders;
  var order;

  while (cycle) {
    if (ret) {
      G.Sleep(Interval);
    }
    blockGetInfo(onOrders);
    orders = GlobalInfo.orders;
    order = foundOrder(orders, Id, dirs);

    if (order === null) {
      break;
    } else {
      exchange.CancelOrder(order.Id);
      ret = true;
    }
  }
  return ret;
}

function valuesToString(values, vpos) {
  var result = '';
  var pos;
  var i;

  if (typeof pos === 'undefined') {
    pos = 0;
  } else {
    pos = vpos;
  }
  for (i = pos; i < values.length; i++) {
    if (i > pos) {
      result += ' ';
    }
    if (values[i] === null) {
      result += 'null';
    } else if (typeof values[i] === 'undefined') {
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

function account2balance(account, symbol) {
  var accounts = account.SubAccounts;

  if (accounts.hasOwnProperty(symbol)) {
    return accounts[symbol].AccountRights;
  }
  return null;
}

// getHoldPosition get the position hold
function getHoldPosition(positions, dir) {
  var len = positions.length;
  var i;

  if (len === 0) {
    return null;
  }
  for (i = 0; i < len; i += 1) {
    if (dir === 0 && positions[i].TradeType === TRADE_TYPE_LONG) {
      return positions[i];
    }
    if (dir === 1 && positions[i].TradeType === TRADE_TYPE_SHORT) {
      return positions[i];
    }
  }
  return null;
}

function positionDiff(positions) {
  var buyPos = getHoldPosition(positions, 0);
  var sellPos = getHoldPosition(positions, 1);
  var buyPosAmount = buyPos === null ? 0 : buyPos.Amount;
  var sellPosAmount = sellPos === null ? 0 : sellPos.Amount;
  var diffAmount = buyPosAmount - sellPosAmount;

  return diffAmount;
}

function totProfit(positions) {
  var buyPos = getHoldPosition(positions, 0);
  var sellPos = getHoldPosition(positions, 1);
  var buyProfitRate = 0;
  var sellProfitRate = 0;

  if (!(buyPos === null)) {
    buyProfitRate = position2Val(buyPos);
  }

  if (!(sellPos === null)) {
    sellProfitRate = position2Val(sellPos);
  }
  return buyProfitRate + sellProfitRate;
}

function GridTrader(dir) {
  var traderDir = dir;
  var vId = 0;
  var realId;
  var found;
  var profitPrice = -1;

  var lenObj;
  var waitOpenOrderBooks = {};
  var waitCoverOrderBooks = {};
  var waitCloseOrderBooks = {};
  var closedOrderBooks = {};
  var orderId;

  this.SetProfitPrice = function(price) {
    profitPrice = price;
  };

  this.GetProfitPrice = function() {
    return profitPrice;
  };

  this.BooksLen = function() {
    lenObj.waitOpen = waitOpenOrderBooks.length;
    lenObj.waitCover = waitCoverOrderBooks.length;
    lenObj.waitClose = waitCloseOrderBooks.length;
    lenObj.history = closedOrderBooks.length;
    return lenObj;
  };

  this.Debug = function() {
    G.Log('waitOpenOrderBooks List:');
    for (orderId in waitOpenOrderBooks) {
      G.Log(waitOpenOrderBooks[orderId]);
    }
    return;
  };

  this.Buy = function(price, amount, extra) {
    G.Log(
      'buy price: ' + price + ',' + 'amount: ' + amount + ', extra:' + extra
    );

    if (typeof extra === 'undefined') {
      extra = '';
    } else {
      extra = valuesToString(arguments, 2);
    }
    vId++;
    orderId = 'V' + vId;

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
    G.Log('Sell price: ' + price + ', amount: ' + amount + ', extra:' + extra);
    //return

    if (typeof extra === 'undefined') {
      extra = '';
    } else {
      extra = valuesToString(arguments, 2);
    }
    vId++;
    var orderId = 'V' + vId;

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

  this.CloseOpenAndCover = function() {
    for (orderId in orderBooks) {
      if (
        orderBooks[orderId].Status === STATE_WAIT_OPEN ||
        orderBooks[orderId].Status === STATE_WAIT_COVER
      ) {
        orderBooks[orderId].Status = STATE_WAIT_CLOSE;
      }
    }
  };

  this.GetHistoryOrders = function() {
    return hisBooks;
  };

  this.ConvertOrder2Close = function(orderId) {
    order = orderBooks[orderId];
    if (order != null) {
      order.Status = STATE_END_CLOSE;
    }
  };

  this.GetOrder = function(orderId) {
    if (typeof orderId === 'string') {
      return exchange.GetOrder(orderId);
    }
    if (typeof hisBooks[orderId] !== 'undefined') {
      return hisBooks[orderId];
    }
    if (typeof orderBooks[orderId] !== 'undefined') {
      return orderBooks[orderId];
    }
    return null;
  };

  // 一个格子的状态机转换
  this.PollOne = function(order, ticker, exchangeOrders, ext) {
    var pfn = order.Type == ORDER_TYPE_BUY ? exchange.Buy : exchange.Sell;
    var pfnDir = order.Type == ORDER_TYPE_BUY ? 'buy' : 'sell';
    var coverPfn = order.Type == ORDER_TYPE_BUY ? exchange.Sell : exchange.Buy;
    var coverPfnDir = order.Type == ORDER_TYPE_BUY ? 'closebuy' : 'closesell';

    // 等待开仓的订单
    if (order.Status == STATE_WAIT_OPEN) {
      var diff = floatReverse(
        order.Type == ORDER_TYPE_BUY ?
          ticker.Buy - order.Price :
          order.Price - ticker.Sell
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
          '(距离: ' +
          diff +
          (order.Type == ORDER_TYPE_BUY ?
            ' 买一: ' + ticker.Buy :
            ' 卖一: ' + ticker.Sell) +
          ')' +
          ' VID:' +
          order.VID
      );
      if (realId != null) {
        order.OpenId = realId;
        order.Status = STATE_WAIT_COVER;
        traderMap.set(order.OpenId, true);
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
          order.Extra + ' 平仓 价格:' + order.CoverPrice + ' VID:' + order.VID
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
        G.Log('close order:' + order.CoverId);
        order.Status = STATE_END_CLOSE;
        traderMap.delete(order.Price);
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

// 动态再平衡, 注意不会平仓调hold部分的仓位
function resetAccount(dir) {
  G.Log('平衡账户mode:', dir);
  cancelPending(dir);
  while (true) {
    blockGetInfo(onOrders, onPosition);
    var orders = order2DirOrder(GlobalInfo.orders, BuyFirst ? 0 : 1);
    var positions = GlobalInfo.positions;

    G.Log('positions:', positions);
    G.Log('orders:', orders);
    var pos = BuyFirst ?
      getHoldPosition(positions, 0) :
      getHoldPosition(positions, 1);

    G.Log('pos:', pos);
    var holdAmount = pos == null ? 0 : pos.Amount;

    if (holdAmount == 0 && orders.length == 0) {
      break;
    }
    if (holdAmount > 0) {
      var leftAmount = holdAmount;

      if (BuyFirst) {
        G.Log('平多仓', leftAmount);
        exchange.SetDirection('closebuy');
        var closeId = exchange.Sell(-1, leftAmount, '平多仓');
      } else {
        G.Log('平空仓', leftAmount);
        exchange.SetDirection('closesell');
        var closeId = exchange.Buy(-1, leftAmount, '平空仓');
      }
    }
    G.Sleep(Interval);
    if (orders.length != 0) {
      cancelPending(dir);
    }
  }
  G.Log('平衡完成');
}

function exit() {
  G.Log('策略开始停止');
  //resetAccount(BuyFirst ? 0 : 1);
  G.Log('策略成功停止');
}

function main() {
  exchange.SetStockType(Coin);
  exchange.SetContractType(ContractType); // 设置合约
  exchange.SetMarginLevel(MarginLevel); // 设置杠杆
  blockGetInfo(onOrders, onAccount, onPosition, onTicker);
}
