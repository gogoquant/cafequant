//
// 回测: action-->waitCheck-->Backtest-->draw-->action_virtual, data from the history files
// 在线模拟: action-->waitCheck-->Backtest-->draw-->action_virtual, data from online 
// 实盘: action-->waitCheck-->draw-->action_real, data from online 
// 绘图: 默认有绘图操作,用户也可在操作中自己绘图
//
// 交易类型支持:
//      bitcoin  
//      期货
//      股票
//      小众韭菜交易途径
//
//
//更多支持:
//      历史行情智能分析
//      简易人工智能训练
//
//
//受众群体:
//      散户
//
//

```
//Entry for this process
function Main(){
    Init()
    main()
}

// wait by time, in case of run too fast
function waitCheck(){

}

// all function try to call backtest
function Backtest(){
    
}

// init the trader
function Init(){
    initposition
    initbacktest
}


//Entry of this policy
function main(){

    // triggered when in backtest mode 
    setcash(cash)
    
    while true{
        ticker  = getticker()

        order()
        sell()
        close()
    }
}

// action_real
def order_real(){} -->order
def sell_real(){}  -->sell
def close_real(){} -->close

// action_virtual
def order_virtual(){} -->order
def sell_virtual(){}  -->sell
def close_virtual(){} -->close
```
