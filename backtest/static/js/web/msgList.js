//按键响应集成接口
var ACTIVELISTLOAD_EVENT = 1;
var PUBLISHLISTLOAD_EVENT = 2;

//追加瀑布
function buildMsgList(json)  
{
    //遍历所有元素
    for(var i=0; i < json.length; i++)  
    {
        buildMsg(json[i])
    }  
}

function buildMsg(json){
        var $grid = msgMaster.waterfull;
        
        var str = ''
        str = '<div class="item" >'
        if(json.img_url != null && json.img_url.length != ""){
            if (json.redir_url != null){
                str = str += '<a href="' + json.redir_url +'">'
            }
            str = str + '<img src="' + msgMaster.qiniuGetImg(json.img_url, 160) + '" />';
            if (json.redir_url != null){
                str = str + '</a>'
            }
        }

        str = str + '<div class="title">'
        if (json.redir_url != null){
                str = str += '<a href="' + json.redir_url +'">'
        }
        
        //太长的标题截取
        var len = json.title.length;
        if (len > 15){
            json.title = json.title.substring(0,15);
        }

        //太长的简介截取
        var len = json.brief.length;
        if (len > 25){
            json.title = json.brief.substring(0,25);
        }

        str = str + json.title 
        if (json.redir_url != null){
                str = str += '<a href="' + json.redir_url +'">'
        }
        str = str + '</div>'
        str = str + '<div class="brief">' + json.brief + '</div>'
        str = str + '</div>'

        loadMsg(str);
}

function loadMsg(str){
        //延迟加载
        setTimeout(function(){ 
            var $grid = msgMaster.waterfull;
            var $elem = $(str); 
            $grid.append($elem);
            $elem.imagesLoaded( function() {                             
                $grid.masonry('appended', $elem );               
            });
        },500 + msgMaster.getRandomNum(200,400));
}

//追加瀑布
function buildPublishList(json)  
{
    buildPublish(json[0]);
}

function buildPublish(json){
    var str =""
    if(json.img_url != null && json.img_url != ""){
        if (json.redir_url != null){
            str = str += '<a href="' + json.redir_url +'">'
        }
        str = str + '<img src="' + msgMaster.qiniuGetImg(json.img_url, 800) + '" />';
        if (json.redir_url != null){
            str = str + '</a>'
        }
    }
    $('#header_publish').html(str);
}

//ajax回调封装
function ajax_success(response, operatType, param) {
    console.log("opt in ajax success", operatType);
    if (response.status == "success") {
        var out = response.msg;
        switch (operatType){
            case ACTIVELISTLOAD_EVENT:   
                {
                    if (out.length < 1){
                        msgMaster.pageend = 1;                   
                    }else{
                        buildMsgList(out);
                    }
                }

                break;
            case PUBLISHLISTLOAD_EVENT:    
                buildPublishList(out);
                break;
            default:
                var ret = "默认响应:" + operatType.toString();
                console.log(ret);
                break;
                //$(".load-div-main").hide();
        }
    }
}

function ajax_error(operatTypem, param){
    alert("服务器错误");
}

function msgLoad(){
    var pos = msgMaster.pagecurr*msgMaster.pagesize
    var count = msgMaster.pagesize
    var tag_id = msgMaster.tag_id;
    var sort_key = msgMaster.sort_key;

    var reqUrl = "/api/pub/search";
    var operatType = ACTIVELISTLOAD_EVENT;
    var data = {};

    data["pos"] = pos;
    data["count"] = count;
    data["msg_type"] = "active";

    console.log("load pos %d, count %d", pos, count);

    if(tag_id != null){
        data["tag_id"] = tag_id;  
    }

    data["sort_key"] = sort_key;

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

//拉取公告图片
function publishLoad(){
    var pos = 0
    var count = 10
    var sort_key = msgMaster.sort_key;

    var reqUrl = "/api/pub/search";
    var operatType = PUBLISHLISTLOAD_EVENT;
    
    var data = {};
    data["pos"] = pos;
    data["count"] = count;
    data["msg_type"] = "publish";
    data["sort_key"] = sort_key;

    console.log("load pos %d, count %d", pos, count);

    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

//极简类
var MsgFactor = {
    //类属性
    //构造方法
    name:"MsgFacotr",

    createNew: function(){
        var obj = Factor.createNew();
        obj.ajaxInit(ajax_success, ajax_error);

        obj.pagesize = 20
        obj.pagecurr  = 0
        obj.pageend = 0;
        obj.tag = null

        obj.title = null     //标题
        obj.content = null   //内容

        //初始化并追加瀑布句柄
        var $msnry = $('#home_pg').masonry({
            columnWidth: 200,//一列的宽度
            gutterWidth:1,//列的间隙
            isRTL:false,//从右到左布局
            
            //isFitWidth:true,// 适应宽度   Boolean
            isResizableL:true,// 是否可调整大小 Boolean
            
            isAnimated:true,//使用jquery的布局变化  Boolean
            animationOptions:{  
                //jquery animate属性 渐变效果  Object { queue: false, duration: 500 }  
            }, 
        });
        
        obj.waterfull = $msnry;
        obj.sort_key = "ctime"
        obj.msgLoad = msgLoad
        obj.buildMsgList = buildMsgList
        obj.buildMsg = buildMsg
        obj.publishLoad = publishLoad
        obj.buildPublishList = buildPublishList
        obj.buildPublish = buildPublish

        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

//注册滚动加载
$(window).scroll(function(){  
    // 当滚动到最底部以上100像素时， 加载新内容  
    if ($(document).height() - $(this).scrollTop() - $(this).height()<100) {
        if(msgMaster.pageend != 1){
            msgMaster.pagecurr = msgMaster.pagecurr + 1;
            msgMaster.msgLoad();
        }
    }
}); 

//实例化
msgMaster = MsgFactor.createNew();
  
//加载动态消息
msgMaster.msgLoad();  

//加载公告消息
msgMaster.publishLoad();  
