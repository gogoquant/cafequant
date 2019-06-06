//按键响应集成接口
var TOPICSUBMIT_EVENT = 0
var TOPICLOAD_EVENT   = 1
var TAGLOAD_EVENT = 2

function ajax_success(response, operatType, param) {
    console.log("opt in ajax success", operatType);
    if (response.status == "success") {
        out = response.msg;
        switch (operatType){
            case TOPICSUBMIT_EVENT: 
                alert("更新文章成功")
                window.location.href="/adminweb/topic_list";  
                break;
            case TAGLOAD_EVENT:    
                tag_str = topicMaster.buildTagList(out);
                $('#tag_lst').html(tag_str);
                //tag load 完成之后，load topic 来设置tag
                topicMaster.topicUpdate();
                break;
            case TOPICLOAD_EVENT:       //更新页面
                topic = out;
                $('#topic_title').val(out["title"]);
                $('#topic_content').val(out["content"]);
                $('#topic_brief').val(out["brief"]);
                
                //不是org才设置选中的tag，默认是未发布
                if(out['tag_id'] != 'org'){
                    tagSet(out["tag_id"]);
                }
                
                //data["tag_id"] = $('#tag_vec option:selected').attr("tag_id");
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
$("#topic_submit").on("click", function(e){
    e.preventDefault();
    //数据获取的地址
    var reqUrl = "";
    //获取事件行为
    var activeType = "POST";
    //提取用户信息
    var title = $('#topic_title').val();
    var content = $('#topic_content').val();
    var brief = $('#topic_brief').val();
    topicMaster.title = title
    topicMaster.content = content
    topicMaster.brief = brief
    topicMaster.topicSubmit()
})


function buildTagList(js){
    tags_str = "";
    tags = js;
    
    tags_str += "<div class=\"form-group col-md-4\">"
    tags_str += "<select id=\"tag_vec\">"

    tags_str += "<option tag_id=" + "org" + ">";
    tags_str += "未发布";
    tags_str += "</option>";

    for(var i=0,l=tags.length; i<l; i++){
            tag = tags[i];
            tags_str += "<option tag_id=" + tag["tag_id"]+ ">";
            tags_str += tag["title"];
            tags_str += "</option>";
    }

    tags_str+="</select></div>"
    return tags_str;
}

//选定该文章的tag标签
function tagSet(tag_id){
    $("#tag_vec option").each(function(){
        if( $(this).attr("tag_id") == tag_id){
            $(this).attr('selected', 'selected');   
        }
    })
}

function topicSubmit(){
    var operatType = TOPICSUBMIT_EVENT;

    //构建将发送的body
    var data = {};
    
    data["title"] = this.title;
    data["content"] = this.content;
    data["brief"] = this.brief;
    var  tag_id  = $('#tag_vec option:selected').attr("tag_id");
    //org表示未发布的tag
    data["tag_id"] = tag_id
    
    var topic_id = this.getQueryString("topic_id");
    
    if(topic_id == null){
        reqUrl = "/api/topic/add";
    }else{
        reqUrl = "/api/topic/update";
        data["topic_id"] = topic_id;
    }

    this.ajaxSend(operatType, "POST", reqUrl, "json", data, null);
}

function topicUpdate(){
    reqUrl = "/api/topic/get";
    operatType = TOPICLOAD_EVENT;

    topic_id = this.getQueryString("topic_id");
    if (topic_id == null){
        return;
    }

    //构建将发送的body
    data = {};
    data["topic_id"] = topic_id;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}


function tagUpdate(){
    reqUrl = "/api/tag/search";
    operatType = TAGLOAD_EVENT;
    data = {};
    data["pos"] = 0;
    data["count"] = 30;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

//极简类
var TopicFactor = {
    //类属性
    //构造方法
    name:"TopicFacotr",
    createNew: function(){
        var obj = Factor.createNew();

        obj.title = null     //标题
        obj.content = null   //内容

        obj.ajaxInit(ajax_success, ajax_error);

        obj.topicSubmit = topicSubmit
        obj.topicUpdate = topicUpdate
        obj.tagUpdate = tagUpdate
        obj.buildTagList = buildTagList
        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

//实例化
topicMaster = TopicFactor.createNew();
topicMaster.tagUpdate();
