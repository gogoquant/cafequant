//按键响应集成接口
var LoginBt = 0
var LogoutBt = 1
var RegisterBt = 2


//按钮事件回调
function sendButtonOpeat(reqUrl, operatType, data, activeType){
    //console.log("opt in ajax %s", operatType);
    $.ajax({
        url : reqUrl,
        data : data,
        type : activeType,
        //交互数据全部使用js
        dataType : "json",
        success : function (response) {
            if (response.status == "success") {
                out = response.data.operat_out;
                switch (operatType){
                        case LoginBt:
                            jsonToSlider(out);
                            break;
                        case LogoutBt:
                            break;
                        case RegisterBt:
                            break;
                        default:
                            ret = "默认响应:" + operatType.toString();
                            console.log(ret);
                            break;
                            //$(".load-div-main").hide();
                }
            }
        },
        error: function(){
            alert("服务器错误");
            //$(".load-div-main").hide();
        }
    });
}

//注册按键的回调事件
$(".misaka-button").on("click", function(){
    //数据获取的地址
    var reqUrl = "";
    //获取事件名称
    var operatType = parseInt($(this).attr("operat-type"));
    //获取事件行为
    var activeType = parseInt($(this).attr("active-type"));
    //用于存放根据事件类型提取的数据
    var data = {};
    switch (operatType){
        case LoginBt:
            break;
        case LogoutBt:
            break;
        case RegisterBt:
            break;
        default:
            alert("未知命令");
            break;
    }
    console.log("Url:%s, operatType:%d, activeType:%s", reqUrl, operatType, activeType);
    sendButtonOpeat(reqUrl, operatType, data, activeType)
})

//极简类
var ClassMaster = {
    //类属性
    sound:"哇哇",
    //构造方法
    createNew: function(){
        var obj = {};
        obj.name = "test"
        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

//setInterval(genTimer,1000);
