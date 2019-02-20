//Entry for this process
function Main(){
    Init()
    main()
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
    setcash(cash)
    
    while true{
        ticker  = getticker()

        order()
        sell()
        close()
    }
}

def order_real(){} -->order
def sell_real(){}  -->sell
def close_real(){} -->close

def order_virtual(){} -->order
def sell_virtual(){}  -->sell
def close_virtual(){} -->close
