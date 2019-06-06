//按键响应集成接口
var USERLIST_EVENT = 0
var USERTOT_EVENT = 1
var USERPREV_EVENT = 2
var USERNEXT_EVENT = 3
var USERDELETE_EVENT = 4

function ajax_success(response, activeType, param) {
    console.log("opt in ajax success", operatType);
    if (response.status == "success") {
        out = response.msg;
        switch (operatType){
            case USERLIST_EVENT:                                    //刷新用户列表
                users_str = userMaster.buildUserList(out);
                //当用户动态变更则重置上限
                if(out.length < userMaster.pagesize){
                    userMaster.pagetot = userMaster.pagecurr;   //将上限缩小
                }
                $('#userList').html(users_str);

                //挂载删除事件
                deleteRegister();

                break;
            case USERTOT_EVENT:
                userMaster.instot = out;                            //设置总用户数

                totUser = Number(out);                              //重新设置页面数
                pagesize = Number(userMaster.pagesize);
                pagecurr = Number(userMaster.pagecurr);
                pagetot = Math.ceil(totUser/pagesize)
                    userMaster.pagetot  = String(pagetot);
                console.log("page tot:", pagetot);

                if (pagecurr < 1){
                    userMaster.pagecurr = 1;
                }

                if(pagecurr >= pagetot){                            //如果当前页位置超过最大则重置
                    userMaster.pagecurr = String(pagetot);
                }
                userMaster.getUserList(null, userMaster.pagecurr, userMaster.pagesize); //查询总数后将当前页刷新下
                break;
            case USERDELETE_EVENT:
                userMaster.getUserTot();                            //获取总数
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
$("#UserNext").on("click", function(){
    //数据获取的地址
    var reqUrl = "";
    //获取事件名称
    var operatType = parseInt($(this).attr("misaka-operat"));
    //获取事件行为
    var activeType = "POST";
    //用于存放根据事件类型提取的数据
    var data = {};
    userMaster.pagecurr  = Number(userMaster.pagecurr) + 1;
    userMaster.getUserTot();
})


//注册按键的回调事件
$("#UserPrev").on("click", function(){
    //数据获取的地址
    var reqUrl = "";
    //获取事件名称
    var operatType = parseInt($(this).attr("misaka-operat"));
    //获取事件行为
    var activeType = "POST";
    //用于存放根据事件类型提取的数据
    var data = {};
    
    userMaster.pagecurr  = Number(userMaster.pagecurr) - 1;
    userMaster.getUserTot();
})


//注册按键的删除事件
function deleteRegister(){
    $(".UserDelete").click( function(){
        //数据获取的地址
        var reqUrl = "";
        //获取用户id
        var user_id = $(this).attr("user_id");
        //获取事件行为
        var activeType = "POST";
        //用于存放根据事件类型提取的数据
        var data = {};
        userMaster.deleteUser(user_id);
    })
}

function deleteUser(user_id){
    data = {};
    reqUrl = "/api/user/delete?" + "user_id="  + user_id;
    operatType = USERDELETE_EVENT;


    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

function getUserList(token, page, pagesize){
    data = {};
    data.page=page;
    data.pagesize=pagesize;
    
    reqUrl = "/api/user/search";
    if(token!=null)
        data.token = String(token);

    operatType = USERLIST_EVENT;


    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

function getUserTot(){
    data = {}
    reqUrl = "/api/user/tot";
    operatType = USERTOT_EVENT;

 
    this.ajaxSend(operatType, "POST", reqUrl, null, data, null);
}

function buildUserList(js){
    users_str = ""
    date = new Date()
    users = js
    for(var i=0,l=users.length; i<l; i++){
            users_str += "<tr>"
            //获取用户对象
            user = users[i]
            //拼接显示信息
            //users_str += "<td class=\"center\">" + user["user_id"]      + "</td>"
            users_str += "<td class=\"center\">" + user["name"]         + "</td>"
            date.setTime(Number(user["date"]) * 1000);
            users_str += "<td class=\"center\">" + date.toDateString()  + "</td>"
            date.setTime(Number(user["last_modify"]) * 1000);
            users_str += "<td class=\"center\">" + date.toDateString()  + "</td>"
            users_str +=
            "<td class=\"center\">"+
                "<a class=\"btn btn-info\" href=\"#\">" +
                    "<i id =\"UserEdit\" class=\"glyphicon glyphicon-edit icon-white\"></i>" +
                        "Edit"+
                "</a>" +
                "<a class=\"btn btn-danger UserDelete \"" + " user_id=" + user["user_id"] + " >" +
                    "<i class=\"glyphicon glyphicon-trash icon-white \"></i>" +
                        "Delete" +
                "</a>" +
            "</td>"
            users_str += "</tr>"
    }
    return users_str
}


//极简类
var UserFactor = {
    //类属性
    //构造方法
    name:"UserFacotr",
    createNew: function(){

        var obj = Factor.createNew();
        obj.ajaxInit(ajax_success, ajax_error);
        
        obj.pagesize ="5"  //每页的数据量
        obj.pagecurr ="1"   //当前选中的页
        obj.pagetot ="0"    //页总数
        obj.instot ="0"     //内容数
        obj.numSelect = "0" 
        obj.buildUserList = buildUserList
        obj.getUserList = getUserList
        obj.getUserTot = getUserTot
        obj.deleteUser = deleteUser
        obj.print = function(){ alert("喵喵喵"); };
        return obj
    }
};

//实例化
userMaster = UserFactor.createNew();
//查询用户总量
userMaster.getUserTot();


