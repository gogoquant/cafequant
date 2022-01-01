// This is an example algorithm

// scripts to loader the data from gateway
function main() {
  var Constract = 'quarter';
  var Symbol = 'BTC/USD';
  var Period = 'M5';
  var IO = 'online';
  var IVL = 5;
  var SIZE = 5;
  
  var records;
  var record;
  var i;
  
  E.SetIO(IO);
  E.SetContractType(Constract);                                                                                              
  E.SetStockType(Symbol);
  E.Ready();
  
  while (true) {
    records = E.GetRecords(Period, '', SIZE);
    if(records === null){
      G.Sleep(IVL * 1000 * 60); 
      continue;
    } 
    records = records[0]
    G.Log(records);
    //var records = records.parseJSON(); 
    for (var i in records) { 
        record = records[i] 
        //G.Log(i, "record is :", records[i]);
        //{"Time":1610433000,"Open":37863.48,"High":38000,"Low":37854.61,"Close":37891.36,"Volume":497520}
        E.BackPutOHLC(record.Time, record.Open, record.High, record.Low, record.Close, record.Volume, "unknown", Period)  
    }
    G.Sleep(IVL * 1000 * 60); 
  }
} 
