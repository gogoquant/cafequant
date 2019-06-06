//按键响应集成接口
var LOAD_EVENT = 0
var CHATADD_EVENT   = 1

function ajax_success(response, operatType, param) {
    console.log("opt in ajax success", operatType);
    if (response.status == "success") {
        out = response.msg;
        switch (operatType){
            case LOAD_EVENT:    
                chat_str = master.buildChatList(out);
                $('#board_message').html(chat_str);
                break;
            case CHATADD_EVENT:
                master.chatload();
                $("#board_blank").val('');
                break;
            default:
                ret = "默认响应:" + operatType.toString();
                console.log(ret);
                break;
                //$(".load-div-main").hide();
        }
    }
}

function ajax_error(operatType, param){
    alert("服务器错误");
}

function buildChatList(js){
    chats_str = "";
    chats = js;
    
    for(var i=0,l=chats.length; i<l; i++){
            chat = chats[i];
            chats_str += "<p>";
            chats_str += chat["name"];
            chats_str += ":";
            chats_str += chat["title"];
            chats_str += "</p>";
    }
    return chats_str;
}

function chatload(){
    console.log("load trigger!");
    reqUrl = "/api/chat/search";
    operatType = LOAD_EVENT;
    data = {};
    data["pos"] = 0;
    data["count"] = 10;

    master.ajaxInit(ajax_success, ajax_error);
    master.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

function add(){
    reqUrl = "/api/chat/add";
    operatType = CHATADD_EVENT;
    data = {};
    data["title"] = this.title;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}


//注册按键的回调事件
$("#board_button").on("click", function(){
    //数据获取的地址
    var reqUrl = "";
    //获取事件行为
    var activeType = "POST";
    //提取用户信息
    var title = $("#board_blank").val();
    master.title = title;
    master.add();
})

//极简类
var ChatFactor = {
    //类属性
    //构造方法
    name:"ChatFacotr",
    createNew: function(){

        var obj = Factor.createNew();
        obj.ajaxInit(ajax_success, ajax_error);
        
        obj.title = null     //标题
        obj.content = null   //内容
        obj.chatload = chatload
        obj.add = add
        obj.buildChatList = buildChatList
        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

$(document).ready(function(){  
    master = ChatFactor.createNew();
    if (master.is_pc()){
    }else{
        //window.location.href='http://www.lancelot.top:8100';
    }
});   
//master.chatload()
//chatTimer = window.setInterval(master.chatload,60000); 
