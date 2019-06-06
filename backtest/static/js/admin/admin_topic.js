//按键响应集成接口
var TOPICLIST_EVENT = 0
var TOPICTOT_EVENT = 1
var TOPICPREV_EVENT = 2
var TOPICNEXT_EVENT = 3
var TOPICDELETE_EVENT = 4
var TOPICEDIT_EVENT = 5

function ajax_success(response, activeType, param) {
    console.log("opt in ajax success", operatType);
    if (response.status == "success") {
        out = response.msg;
        switch (operatType){
            case TOPICLIST_EVENT:                                    //刷新用户列表
                topics_str = topicMaster.buildTopicList(out);
                $('#topicList').html(topics_str);

                //挂载删除事件
                deleteRegister();
                //挂载编辑事件
                editRegister();

                break;
            case TOPICTOT_EVENT:
                topicMaster.instot = out;                            //设置总用户数

                totTopic = Number(out);                              //重新设置页面数
                pagesize = Number(topicMaster.pagesize);
                pagecurr = Number(topicMaster.pagecurr);
                pagetot = Math.ceil(totTopic/pagesize)
                    topicMaster.pagetot  = String(pagetot);

                console.log("user tot:", out);
                console.log("page tot:", pagetot);

                //验值并恢复
                if(pagecurr < 0){
                    topicMaster.pagecurr = 0;
                }

                if(pagecurr >= pagetot){                            //如果当前页位置超过最大则重置
                    topicMaster.pagecurr = pagetot -1
                }

                topicMaster.getTopicList(null, topicMaster.pagecurr*topicMaster.pagesize, topicMaster.pagesize); //查询总数后将当前页刷新下
                break;
            case TOPICDELETE_EVENT:
                topicMaster.getTopicTot();                            //获取总数
                break;
            case TOPICEDIT_EVENT:
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
$("#TopicNext").on("click", function(){
    //数据获取的地址
    var reqUrl = "";
    //获取事件名称
    var operatType = parseInt($(this).attr("misaka-operat"));
    //获取事件行为
    var activeType = "POST";
    //用于存放根据事件类型提取的数据
    var data = {};
    
    topicMaster.pagecurr  = Number(topicMaster.pagecurr) + 1;
    topicMaster.getTopicTot();
})


//注册按键的回调事件
$("#topic_new").on("click", function(){
    reqUrl = "/adminweb/topic_editor";
    reqUrl = reqUrl;
    //redirect
    window.open(reqUrl);  
})

//注册按键的回调事件
$("#TopicPrev").on("click", function(){
    //数据获取的地址
    var reqUrl = "";
    //获取事件名称
    var operatType = parseInt($(this).attr("misaka-operat"));
    //获取事件行为
    var activeType = "POST";
    //用于存放根据事件类型提取的数据
    var data = {};
    topicMaster.pagecurr  = Number(topicMaster.pagecurr) - 1;
    topicMaster.getTopicTot();
})


//注册按键的删除事件
function deleteRegister(){
    $(".TopicDelete").on("click", function(){
        //数据获取的地址
        var reqUrl = "";
        //获取用户id
        var topic_id = $(this).attr("topic_id");
        //获取事件行为
        var activeType = "POST";
        //用于存放根据事件类型提取的数据
        var data = {};
        topicMaster.deleteTopic(topic_id);
    })
}

function deleteTopic(topic_id){
    reqUrl = "/api/topic/delete";
    data = {};
    data["topic_id"] = topic_id;
    operatType = TOPICDELETE_EVENT;


    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

//注册按键的删除事件
function editRegister(){
    $(".TopicEdit").on("click", function(){
        //数据获取的地址
        var reqUrl = "";
        //获取用户id
        var topic_id = $(this).attr("topic_id");
        //获取事件行为
        var activeType = "GET";
        topicMaster.editTopic(topic_id);
    })
}

function editTopic(topic_id){
    reqUrl = "/adminweb/topic_editor";
    reqUrl = reqUrl + "?topic_id=" + topic_id;
    //redirect
    window.open(reqUrl);  
}

function getTopicList(token, page, pagesize){
    data = {};
    data.pos = page
    data.count = pagesize
    data.sort_key = "mtime";
    
    reqUrl = "/api/topic/search"
    if(token!=null)
        data.token = String(token)

    console.log("search topic %s->%s", String(data.pos), String(data.count));

    operatType = TOPICLIST_EVENT;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

function getTopicTot(){
    data = {}
    reqUrl = "/api/topic/tot";
    operatType = TOPICTOT_EVENT;


    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

function buildTopicList(js){
    topics_str = ""
    date = new Date()
    topics = js
    for(var i=0,l=topics.length; i<l; i++){
            topics_str += "<tr>"
            //获取用户对象
            topic = topics[i]
            //拼接显示信息
            //topics_str += "<td class=\"center\">" + topic["topic_id"]      + "</td>"
            topics_str += "<td class=\"center\">" + topic["title"]         + "</td>"
            date.setTime(Number(topic["date"]) * 1000);
            topics_str += "<td class=\"center\">" + date.toDateString()  + "</td>"
            date.setTime(Number(topic["last_modify"]) * 1000);
            topics_str += "<td class=\"center\">" + date.toDateString()  + "</td>"
            topics_str +=
            "<td class=\"center\">"+
                "<a class=\"btn btn-info  TopicEdit\"" + " topic_id=" + topic["topic_id"] + ">" +
                    "<i id =\"TopicEdit\" class=\"glyphicon glyphicon-edit icon-white\"></i>" +
                        "Edit"+
                "</a>" +
                "<a class=\"btn btn-danger TopicDelete\"" + " topic_id=" + topic["topic_id"] + ">" +
                    "<i class=\"glyphicon glyphicon-trash icon-white\" ></i>" +
                        "Delete" +
                "</a>" +
            "</td>"
            topics_str += "</tr>"
    }
    return topics_str
}


//极简类
var TopicFactor = {
    //类属性
    //构造方法
    name:"TopicFacotr",
    createNew: function(){
        var obj = Factor.createNew();
        obj.ajaxInit(ajax_success, ajax_error);
        
        obj.pagesize ="5"  //每页的数据量
        obj.pagecurr ="0"   //当前选中的页
        obj.pagetot ="0"    //页总数
        obj.instot ="0"     //内容数
        obj.numSelect = "0" 
        obj.buildTopicList = buildTopicList
        obj.getTopicList = getTopicList
        obj.getTopicTot = getTopicTot
        obj.deleteTopic = deleteTopic
        obj.editTopic = editTopic
        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

//实例化
topicMaster = TopicFactor.createNew();

//查询用户总量
topicMaster.getTopicTot();


