//按键响应集成接口
var MSGLIST_EVENT = 0
var MSGTOT_EVENT = 1
var MSGPREV_EVENT = 2
var MSGNEXT_EVENT = 3
var MSGDELETE_EVENT = 4
var MSGEDIT_EVENT = 5

function ajax_success(response, activeType, param) {
    console.log("opt in ajax success", operatType);
    if (response.status == "success") {
        out = response.msg;
        switch (operatType){
            case MSGLIST_EVENT:                                    //刷新用户列表
                msgs_str = msgMaster.buildMsgList(out);
                $('#msgList').html(msgs_str);

                //挂载删除事件
                deleteRegister();
                //挂载编辑事件
                editRegister();

                break;
            case MSGTOT_EVENT:
                msgMaster.instot = out;                            //设置总用户数

                totMsg = Number(out);                              //重新设置页面数
                pagesize = Number(msgMaster.pagesize);
                pagecurr = Number(msgMaster.pagecurr);
                pagetot = Math.ceil(totMsg/pagesize)
                    msgMaster.pagetot  = String(pagetot);

                console.log("user tot:", out);
                console.log("page tot:", pagetot);

                //验值并恢复
                if(pagecurr < 0){
                    msgMaster.pagecurr = 0;
                }

                if(pagecurr >= pagetot){                            //如果当前页位置超过最大则重置
                    msgMaster.pagecurr = pagetot -1
                }

                msgMaster.getMsgList(null, msgMaster.pagecurr*msgMaster.pagesize, msgMaster.pagesize); //查询总数后将当前页刷新下
                break;
            case MSGDELETE_EVENT:
                msgMaster.getMsgTot();                            //获取总数
                break;
            case MSGEDIT_EVENT:
                break;
            default:
                ret = "默认响应:" + operatType.toString();
                console.log(ret);
                break;
                //$(".load-div-main").hide();
        }
    }
}
function  ajax_error(operatType, param){
    alert("服务器错误");
}

//注册按键的回调事件
$("#MsgNext").on("click", function(){
    //数据获取的地址
    var reqUrl = "";
    //获取事件名称
    var operatType = parseInt($(this).attr("misaka-operat"));
    //获取事件行为
    var activeType = "POST";
    //用于存放根据事件类型提取的数据
    var data = {};
    
    msgMaster.pagecurr  = Number(msgMaster.pagecurr) + 1;
    msgMaster.getMsgTot();
})


//注册按键的回调事件
$("#msg_new").on("click", function(){
    reqUrl = "/adminweb/msg";
    reqUrl = reqUrl;
    //redirect
    window.open(reqUrl);  
})

//注册按键的回调事件
$("#MsgPrev").on("click", function(){
    //数据获取的地址
    var reqUrl = "";
    //获取事件名称
    var operatType = parseInt($(this).attr("misaka-operat"));
    //获取事件行为
    var activeType = "POST";
    //用于存放根据事件类型提取的数据
    var data = {};
    msgMaster.pagecurr  = Number(msgMaster.pagecurr) - 1;
    msgMaster.getMsgTot();
})


//注册按键的删除事件
function deleteRegister(){
    $(".MsgDelete").on("click", function(){
        //数据获取的地址
        var reqUrl = "";
        //获取用户id
        var msg_id = $(this).attr("msg_id");
        //获取事件行为
        var activeType = "POST";
        //用于存放根据事件类型提取的数据
        var data = {};
        msgMaster.deleteMsg(msg_id);
    })
}

function deleteMsg(msg_id){
    reqUrl = "/api/pub/delete";
    data = {};
    data["msg_id"] = msg_id;
    operatType = MSGDELETE_EVENT;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

//注册按键的删除事件
function editRegister(){
    $(".MsgEdit").on("click", function(){
        //数据获取的地址
        var reqUrl = "";
        //获取用户id
        var msg_id = $(this).attr("msg_id");
        //获取事件行为
        var activeType = "GET";
        msgMaster.editMsg(msg_id);
    })
}

function editMsg(msg_id){
    reqUrl = "/adminweb/msg";
    reqUrl = reqUrl + "?msg_id=" + msg_id;
    //redirect
    window.open(reqUrl);  
}

function getMsgList(token, page, pagesize){
    data = {};
    data.pos = page
    data.count = pagesize
    data.sort_key = "mtime";
    
    reqUrl = "/api/pub/search"
    if(token!=null)
        data.token = String(token)

    console.log("search msg %s->%s", String(data.pos), String(data.count));

    operatType = MSGLIST_EVENT;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

function getMsgTot(){
    data = {}
    reqUrl = "/api/pub/count";
    operatType = MSGTOT_EVENT;


    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

function buildMsgList(js){
    msgs_str = ""
    date = new Date()
    msgs = js
    for(var i=0,l=msgs.length; i<l; i++){
            msgs_str += "<tr>"
            //获取用户对象
            msg = msgs[i]
            //拼接显示信息
            //msgs_str += "<td class=\"center\">" + msg["msg_id"]      + "</td>"
            msgs_str += "<td class=\"center\">" + msg["title"]         + "</td>"
            date.setTime(Number(msg["date"]) * 1000);
            msgs_str += "<td class=\"center\">" + date.toDateString()  + "</td>"
            date.setTime(Number(msg["last_modify"]) * 1000);
            msgs_str += "<td class=\"center\">" + date.toDateString()  + "</td>"
            msgs_str +=
            "<td class=\"center\">"+
                "<a class=\"btn btn-info  MsgEdit\"" + " msg_id=" + msg["pubmsg_id"] + ">" +
                    "<i id =\"MsgEdit\" class=\"glyphicon glyphicon-edit icon-white\"></i>" +
                        "Edit"+
                "</a>" +
                "<a class=\"btn btn-danger MsgDelete\"" + " msg_id=" + msg["pubmsg_id"] + ">" +
                    "<i class=\"glyphicon glyphicon-trash icon-white\" ></i>" +
                        "Delete" +
                "</a>" +
            "</td>"
            msgs_str += "</tr>"
    }
    return msgs_str
}


//极简类
var MsgFactor = {
    //类属性
    //构造方法
    name:"MsgFacotr",
    createNew: function(){
        var obj = Factor.createNew();
        obj.ajaxInit(ajax_success, ajax_error);
        
        obj.pagesize ="5"  //每页的数据量
        obj.pagecurr ="0"   //当前选中的页
        obj.pagetot ="0"    //页总数
        obj.instot ="0"     //内容数
        obj.numSelect = "0" 
        obj.buildMsgList = buildMsgList
        obj.getMsgList = getMsgList
        obj.getMsgTot = getMsgTot
        obj.deleteMsg = deleteMsg
        obj.editMsg = editMsg
        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

//实例化
msgMaster = MsgFactor.createNew();

//查询用户总量
msgMaster.getMsgTot();


