//按键响应集成接口
var FILESUBMIT_EVENT = 0

function ajax_success(response, operatType, param) {
    console.log("opt in ajax success", operatType);
    if (response.status == "success") {
        out = response.msg;
        switch (operatType){
            case FILESUBMIT_EVENT:
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
$("#file_submit").on("click", function(e){
    e.preventDefault();
    
    //数据获取的地址
    var reqUrl = "/api/file/add";
    //获取事件行为
    var activeType = "POST";

    var operatType = FILESUBMIT_EVENT;

    var form = new FormData(document.getElementById("file_form"));
    //提取用户信息
    //var formData = new FormData();  // 要求使用的html对象
    //formData.appen("file", $( "#file_data" ).eq(0).files[0])

    uploadMaster.ajaxSend(operatType, "POST", reqUrl, "form", form, null);
})


//极简类
var UploadFactor = {
    //类属性
    //构造方法
    name:"UploadFacotr",
    createNew: function(){
        var obj = Factor.createNew();
        obj.ajaxInit(ajax_success, ajax_error);
        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

uploadMaster = UploadFactor.createNew();


