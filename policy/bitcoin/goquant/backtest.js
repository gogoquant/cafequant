function BackOrderMap() {
  var i;
  var v;
  var struct = function(key, value) {
    this.key = key;
    this.value = value;
  };
  var put = function(key, value) {
    for (i = 0; i < this.arr.length; i++) {
      if (this.arr[i].key === key) {
        this.arr[i].value = value;
        return;
      }
    }
    this.arr[this.arr.length] = new struct(key, value);
  };
  var get = function(key) {
    for (i = 0; i < this.arr.length; i++) {
      if (this.arr[i].key === key) {
        return this.arr[i].value;
      }
    }
    return null;
  };
  var remove = function(key) {
    for (i = 0; i < this.arr.length; i++) {
      v = this.arr.pop();
      if (v.key === key) {
        continue;
      }
      this.arr.unshift(v);
    }
  };
  var size = function() {
    return this.arr.length;
  };
  var isEmpty = function() {
    return this.arr.length <= 0;
  };

  this.arr = new Array();
  this.get = get;
  this.put = put;
  this.remove = remove;
  this.size = size;
  this.isEmpty = isEmpty;
}

ExchangeBack = function() {
  cnt = 1;
  backorders = new BackOrderMap();
  this.Buy = function(price, amount, extra) {
    G.Log('Buy price: ' + price + ', amount: ' + amount + ', extra:' + extra);
    cnt += 1;
    return cnt;
  };
  this.Sell = function(price, amount, extra) {
    G.Log('Sell price: ' + price + ', amount: ' + amount + ', extra:' + extra);
    cnt++;
    return cnt;
  };
  this.CancelOrder = function(dir) {
    G.Log('cancel order:', dir);
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

exchange = new ExchangeBack();
