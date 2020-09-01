// author mhw
// brief safe box 
// created 201910281001
// update  201912250000

// 平仓
function cover(pos) {
    var closeId = 0
    if (pos.length > 0) {
        var leftAmount = pos[0].Amount
        while(leftAmount > 0){
            //分批平仓
            var subAmount = leftAmount > MaxCover ? MaxCover : leftAmount
            if (pos[0].Type == 0) { //平多仓，采用盘口吃单，会损失手续费，可改为盘口挂单，会增加持仓风险。
                Log("平多仓", subAmount)  
                exchange.SetDirection('sell')
                var closeId = exchange.Sell(-1, subAmount, '平多仓')
            } else {
                Log("平空仓", subAmount)
                exchange.SetDirection('buy')
                var closeId = exchange.Buy(-1, subAmount, '平空仓')
            }
            Sleep(Intervel)
            leftAmount -= subAmount
        }
    }
}

// watchBox 箱体检测
function watchBox(price) {
    if (price > HighBox) {
        return true
    }
    if (price < LowBox) {
        return true
    }
    return false
}


function main() {
    exchange.SetContractType('swap')
    Log(exchange.GetAccount())
    while (true) {

        ticker_curr = exchange.GetTicker()
        Sleep(Intervel)
        last_price = ticker_curr.Last

        // 箱体,安全检测
        if (watchBox(last_price)) {
            Log("异常价格:" + JSON.stringify({
                lastPrice: last_price,
                HighBox: HighBox,
                LowBox: LowBox
            }))
            break
        }else{
            LogStatus(JSON.stringify({
                lastPrice: last_price,
                HighBox: HighBox,
                LowBox: LowBox
            }))
        }
        Sleep(Intervel)
       
    }

    // 平仓清理
    while (true) {
        Log("异常情况尝试平仓")
        // 定时
        var pos = exchange.GetPosition()
        if(typeof pos == "undefined" || pos == null){
            continue
        }

        if (pos.length <= 0) {
            Log("异常情况清理完毕程序退出")
            break
        }
        //尝试平仓
        cover(pos)
        
        Sleep(Intervel)
    }
}

