//按键响应集成接口
var MSGSUBMIT_EVENT = 0
var MSGLOAD_EVENT   = 1

function ajax_success(response, operatType, param) {
    console.log("opt in ajax success", operatType);
    if (response.status == "success") {
        out = response.msg;
        switch (operatType){
            case MSGSUBMIT_EVENT: 
                alert("更新文章成功")
                window.location.href="/admin/msg_list";  
                break;
            case MSGLOAD_EVENT:       //更新页面
                msg = out;
                $('#msg_title').val(out["title"]);
                $('#msg_content').val(out["content"]);
                $('#msg_brief').val(out["brief"]);
                $('#msg_type').val(out["msg_type"]);
                $('#img_url').val(out["img_url"]);
                $('#redir_url').val(out["redir_url"]);
                $('#user_id').val(out["user_id"]);
                $('#user_to').val(out["user_to"]);
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

//注册按键的回调事件
$("#msg_submit").on("click", function(e){
    e.preventDefault();
    //数据获取的地址
    var reqUrl = "";
    //获取事件行为
    var activeType = "POST";
    //提取用户信息
    msgMaster.msgSubmit()
})


function msgSubmit(){
    operatType = MSGSUBMIT_EVENT;

    //构建将发送的body
    data = {};

    var title = $('#msg_title').val();
    var content = $('#msg_content').val();
    var brief = $('#msg_brief').val();
    var msg_type = $('#msg_type').val();
    var img_url = $('#img_url').val();
    var redir_url = $('#redir_url').val();
    var img_url = $('#img_url').val();
    var user_id = $('#user_id').val();
    var user_to = $('#user_to').val();

    
    data["title"] = title;
    data["content"] = content;
    data["brief"] = brief;
    data["msg_type"] = msg_type;
    data["img_url"] = img_url;
    data["redir_url"] = redir_url;
    data["user_id"] = user_id;
    data["user_to"] = user_to;

    msg_id = this.getQueryString("msg_id");
    
    if(msg_id == null){
        reqUrl = "/api/pub/add";
    }else{
        reqUrl = "/api/pub/update";
        data["msg_id"] = msg_id;
    }

    this.ajaxSend(operatType, "POST", reqUrl, "json", data, null);
}

function msgUpdate(){
    reqUrl = "/api/pub/get";
    operatType = MSGLOAD_EVENT;

    msg_id = this.getQueryString("msg_id");
    if (msg_id == null){
        return;
    }

    //构建将发送的body
    data = {};
    data["msg_id"] = msg_id;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}


//极简类
var MsgFactor = {
    //类属性
    //构造方法
    name:"MsgFacotr",
    createNew: function(){
        var obj = Factor.createNew();

        obj.title = null     //标题
        obj.content = null   //内容

        obj.ajaxInit(ajax_success, ajax_error);

        obj.msgSubmit = msgSubmit
        obj.msgUpdate = msgUpdate
        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

//实例化
msgMaster = MsgFactor.createNew();
msgMaster.msgUpdate();
