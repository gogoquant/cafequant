//按键响应集成接口
var TOPICLOAD_EVENT   = 1

function ajax_success(response, operatType) {
    console.log("opt in ajax success", operatType);
    if (response.status == "success") {
        out = response.msg;
        switch (operatType){
            case TOPICLOAD_EVENT:       //更新页面
                msg = out;
                $('#msg_title').html("<p>" + out["title"] + "</p>");
                $('#msg_info').html("<p>作者: " + out["user_info"]["name"] + "</p>");
                $('#tail_info').html("<p>" + master.getLocalTime(out["mtime"]) + "</p>");
                document.getElementById('msg_body').innerHTML = marked(out["content"]);
                break;
            default:
                ret = "默认响应:" + operatType.toString();
                console.log(ret);
                break;
                //$(".load-div-main").hide();
        }
    }
}

function ajax_error(operatType){
    alert("服务器错误");
}
    
function msgLoad(){
    reqUrl = "/api/msg/get";
    operatType = TOPICLOAD_EVENT;

    /*获取需要查询的文章id*/
    msg_id = getQueryString("msg_id");
    if (msg_id == null){
        return;
    }

    //构建将发送的body
    data = {};
    data["msg_id"] = msg_id;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

//极简类
var TopicFactor = {
    //类属性
    //构造方法
    name:"TopicFacotr",
    

    createNew: function(){
        var obj = Factor.createNew();
        obj.ajaxInit(ajax_success, ajax_error);
        
        obj.title = null     //标题
        obj.content = null   //内容
        obj.msgLoad = msgLoad
        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

marked.setOptions({
    renderer: new marked.Renderer(),
    gfm: true,
    tables: true,
    breaks: false,
    pedantic: false,
    sanitize: true,
    smartLists: true,
    smartypants: false
});

$(document).ready(function(){  
    console.log(marked('I am using __markdown__.'));
    master = TopicFactor.createNew();
    master.msgLoad();
});   
