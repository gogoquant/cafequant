//按键响应集成接口

var TAGLOAD_EVENT = 0;
var LISTLOAD_EVENT = 1;


function buildTopicList(json)  
{  
    var oProduct, $topicLine, iHeight, iTempHeight;  
    for(var i=0, l=json.length; i<l; i++)  
    {  
        if(json[i]['tag_info'] == null){
            continue;
        }

        oProduct = json[i];  

        $topicLine = null;

        // 找出当前高度最小的topic元素, 新内容添加到该列  
        iHeight = -1;  
        $('#topic_list li').each(function(){  
            iTempHeight = Number( $(this).height() );  
            if(iHeight==-1 || iHeight>iTempHeight)  
            {  
                iHeight = iTempHeight;  
                $topicLine = $(this);  
            }  
        }); 

        var elemStr = "";
        elemStr = elemStr + "<li class = \"topic_elem\">";
        elemStr = elemStr + "<div class=\"topic_elem_title\"><p>" ;
        elemStr = elemStr + "<a href=" + "/web/topic?topic_id="  + json[i]["topic_id"] + ">";
        elemStr = elemStr + json[i]["title"] + "</p></a></div>";
        elemStr = elemStr + "<div class=\"topic_elem_brief\"><p>" + json[i]["brief"] + "</p></div>";
        elemStr = elemStr + "<div class=\"topic_elem_tag\">" + '<span class=\"mini_tag\">from</span>' + '<span class="mini_info"> '   
            + json[i]["tag_info"]["title"] + "</span></div>"
        elemStr = elemStr + "<div class=\"topic_elem_ext\"><p>"  + this.getLocalTime(json[i]["mtime"]) + "</p></div>";
        elemStr = elemStr + "</li>" 

        $item = $(elemStr).hide();  

        if($topicLine == null){
            $('#topic_list').html(elemStr);
        }else{
            $topicLine.append($item); 
        }
        $item.fadeIn();  
    }  
}  

function ajax_success(response, operatType, param) {
    console.log("opt in ajax success with type:", operatType);
    if (response.status == "success") {
        out = response.msg;
        switch (operatType){
            case TAGLOAD_EVENT:    
                tag_str = master.buildTagList(out);
                $('#tag_lst').html(tag_str);
                tagRegister();
                break;
            case LISTLOAD_EVENT:  
                if(out.length < 1){
                    master.pageend = 1;
                }else{
                    buildTopicList(out);
                }
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

//注册按键的删除事件
function tagRegister(){
    $(".tag_button").on("click", function(){
        //数据获取的地址
        var reqUrl = "";

        //获取tag_id,全局查找截获的为空
        //

        $curr = this;

        $('#tag_lst>ul>li').each(function(){  
            if(this != $curr)  
            {  
                $(this).removeClass("active");
            }  
        }); 

        $(this).addClass("active");

        
        var tag_id = $(this).attr("tag_id");
        //获取事件行为
        var activeType = "GET";
        
        $('#topic_list').empty();

        //重置master的游标
        master.pagecurr = 0;
        master.tag_id = tag_id;

        master.topicLoad();  
    })
}


//注册按键的删除事件
function sortRegister(){
    $(".display_list").on("click", function(){
        //数据获取的地址
        var reqUrl = "";

        $curr = this;

        $('#topic_head>ul>li').each(function(){  
            if(this != $curr)  
            {  
                $(this).removeClass("active");
            }  
        }); 
        $(this).addClass("active");


        $('#topic_list').empty();
            
        //set sort policy
        master.pagecurr = 0;
        master.sort_key = $(this).attr("sort");
        
        master.topicLoad();
    })
}

function buildTagList(js){
    tags_str = "";
    tags = js;
    tags_str += "<ul>"
    
    /*全局查找标签*/
    tags_str = tags_str + "<li class=\"tag_button active\">"
    tags_str += "全部";
    tags_str += "</li>";

    for(var i=0,l=tags.length; i<l; i++){
            tag = tags[i];
            tags_str = tags_str + "<li class=\"tag_button\" tag_id= " + tag["tag_id"] +">";
            tags_str += tag["title"];
            tags_str += "</li>";
    }
    tags_str += "</ul>"
    return tags_str;
}

function tagLoad(){
    reqUrl = "/api/tag/search";
    operatType = TAGLOAD_EVENT;
    data = {};
    data["pos"] = 0;
    data["count"] = 8;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

function topicLoad(){
    pos = master.pagecurr*master.pagesize
    count = master.pagesize
    tag_id = master.tag_id;
    sort_key = master.sort_key;

    reqUrl = "/api/topic/search";
    operatType = LISTLOAD_EVENT;
    data = {};

    data["pos"] = pos;
    data["count"] = count;

    if(tag_id != null){
        data["tag_id"] = tag_id;  
    }

    data["sort_key"] = sort_key;


    console.log("topic load pos %d count %d", pos, count);
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

        obj.pagesize = 10
        obj.pagecurr  = 0
        obj.pageend = 0;
        obj.tag = null

        obj.title = null     //标题
        obj.content = null   //内容

        obj.sort_key = "mtime"

        obj.tagLoad = tagLoad
        obj.topicLoad = topicLoad
        obj.buildTagList = buildTagList
        obj.buildTopicList = buildTopicList
        obj.sortRegister = sortRegister

        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

//实例化
master = TopicFactor.createNew();

$(window).scroll(function(){  
    // 当滚动到最底部以上100像素时， 加载新内容  
    if ($(document).height() - $(this).scrollTop() - $(this).height()<100) {
        if(master.pageend != 1){
            master.pagecurr = master.pagecurr + 1;
            master.topicLoad();
        }
    }
});  
  
$(document).ready(function(){  
    //加载标签
    master.tagLoad();
    //加载文章
    master.topicLoad();  

    master.sortRegister();
});   
