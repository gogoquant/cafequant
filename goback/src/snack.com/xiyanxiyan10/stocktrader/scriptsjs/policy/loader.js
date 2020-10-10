 
function main() {
  var Constract = 'quarter';
  var Symbol = 'BTC/USD';
  var Period = 'M30';
  var IO = 'online';
  var records;
  var i;
  var elem;
  var timestamp;
  var vec = [];
 
  E.SetIO(IO);
  E.SetContractType(Constract);
  E.SetStockType(Symbol);
  E.Ready();
 
  while (true) {
    records = E.GetRecords(Period, '');
    if(records === null){
      G.Sleep(5 * 1000); 
        continue;
    }
    //records = JSON.parse(records)
    for (i in records) {
      var v = {}
      elem = records[i]
      timestamp = new Date(elem.Time*1000).toLocaleString();
      v.Timestamp = timestamp
      vec.push(v)
    }
    G.Log(vec);      
    G.Sleep(5 * 1000);                                                                                                                                                       
  }
}
