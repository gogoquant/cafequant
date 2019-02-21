# QunatPro

# 1. 用户画像
- 业余散户
- 交易员
- 程序员
- 中小企业

# 2. 盈利模式
- 在线简单服务赚取流量
- 托管交易赚取租金
- 吸引第三方赚取流量
- 第三方合作
- 开设课程
- 参考发明者

# 3. 竞品分析

# 4. 优势
- 简化Ai辅助新手
- 门槛低
- 团队人员少

# 5. 劣势


# 6. 成本评估


# 7. 产品生命周期


# 8. 交易渠道支持

- bitcoin  
- 期货
- 股票
- 小众途径


# 9. 功能
- 历史数据回测
- 在线回测
- 实盘交易
- 历史行情智能分析
- 简易人工智能训练
- 论坛交流

# 10. 伪代码核心流程
```
// 回测: action-->waitCheck-->Backtest-->draw-->action_virtual, data from the history files
// 在线模拟: action-->waitCheck-->Backtest-->draw-->action_virtual, data from online 
// 实盘: action-->waitCheck-->draw-->action_real, data from online 
// 绘图: 默认有绘图操作,用户也可在操作中自己绘图
//
//
//

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
