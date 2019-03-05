/*
 * @brief 通用js库
 * @author lancelot
 * @date 2017/3/23
 *
 */

function getLocalTime(nS) {     
    return new Date(parseInt(nS) * 1000).toLocaleString().replace(/:\d{1,2}$/,' ');     
}

function getRandomNum(Min,Max)
{   
    var Range = Max - Min;   
    var Rand = Math.random();   
    return(Min + Math.round(Rand * Range));   
}   

//获取对应的图片地址
function qiniuGetImg(reqUrl, width){
    //reqUrl="hq.jpg"
    infoUrl="imageinfo";
    //baseUrl="http://7xki2b.com1.z0.glb.clouddn.com"
    viewUrl="imageView/2/w"

    //var url = baseUrl + "/" +  reqUrl;
    var url = reqUrl;

    if (width != null){
        url = url + "?" + viewUrl + "/" + width;
    }
    console.log("build img req  %s", url);
    return url;
}

//获取参数路径
function getQueryString(name)
{
    var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
    var r = window.location.search.substr(1).match(reg);
    if(r != null && r.toString().length>1) return  unescape(r[2]); 
    return null;
}

function ajaxRegister(reqUrl, operatType, data, activeType, contentType){
    this.reqUrl = reqUrl;
    this.operatType = operatType;
    this.data = data;
    this.activeType = activeType;
    this.contentType = contentType;
}

function ajaxInit(func, func_err){
    this.func = func;
    this.func_err = func_err;
}

//发送ajax封装
function ajaxSend(operatType, activeType, reqUrl, contentType, data, param){
    
    var func = this.func;
    var func_err = this.func_err;

    console.log("opt in ajax %s->%s", operatType, reqUrl);
    
    //待发送的ajax数据
    var ajax_body = {
        url : reqUrl,
        type : activeType,
        dataType : "json",
        
        success : function (response) {
            func(response, operatType, param);
        },
        
        error: function(){
            func_err(operatType, param);
        }
    };
    //提交json
    if(contentType == "json"){
        ajax_body["contentType"] = "application/json; charset=utf-8";
        ajax_body["data"] = JSON.stringify(data);
    }else if(contentType == "form"){
        //提交表单
        ajax_body["async"] = true,
        ajax_body["cache"] = false,
        ajax_body["processData"] = false,
        ajax_body["contentType"] = false, 
        ajax_body["data"] = data;
    }else{
        //提交参数
        ajax_body["data"] = data;
    }
    $.ajax(ajax_body);
}

function is_pc() {
    var userAgentInfo = navigator.userAgent;
    var Agents = ["Android", "iPhone",
        "SymbianOS", "Windows Phone",
        "iPad", "iPod"];
    var flag = true;
    for (var v = 0; v < Agents.length; v++) {
        if (userAgentInfo.indexOf(Agents[v]) > 0) {
            flag = false;
            break;
        }
    }
    return flag;
}

//极简类, 构造方法
var Factor = {
    
    //共享属性
    name:"Facotr",
    
    createNew: function(){
        var obj = {};
        
        obj.func = null;
        obj.func_err = null;

        obj.is_pc = is_pc;
        obj.ajaxInit = ajaxInit;
        obj.ajaxSend = ajaxSend;
        obj.ajaxRegister = ajaxRegister;
        obj.getQueryString = getQueryString;
        obj.getLocalTime = getLocalTime;
        obj.qiniuGetImg = qiniuGetImg;
        obj.getRandomNum = getRandomNum;

        obj.debug = function(){ alert(Factor.name); };
        return obj;
    }
};
