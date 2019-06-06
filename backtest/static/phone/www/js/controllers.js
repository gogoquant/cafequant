angular.module('starter.controllers', [])

.controller('DashCtrl', function($scope) {})

.controller('ChatsCtrl', function($scope, $stateParams, $http, $state, $ionicPopup) {
    
    $scope.hasmore=true;
    
    var run = false;    //模拟线程锁机制  防止多次请求 含义：是否正在请求。请注意，此处并非加入到了就绪队列，而是直接跳过不执行
    var obj = {pos:0,count:8};
    var result = infoPull(obj,1);  

    console.log($scope.hasmore+"是否加载更多");
    
    //刷新
    $scope.doRefresh = function(){
         var obj_data = {pos:0,count:8};
         var result = infoPull(obj_data,2);
         $scope.$broadcast('scroll.refreshComplete');
    };
    
    //滚动加载
    $scope.loadMore = function(){   
        var old = $scope.items;
        if(old!=undefined){
          var result = infoPull(obj,3); 
        }   
        $scope.$broadcast('scroll.infiniteScrollComplete');
    };

    /* state:1初始化，2刷新，3加载更多 */
    function infoPull(obj_data,state){
        if(!run){
            run = true;
            $http({
               method:"POST",
               url:'http://www.lancelot.top:80/api/topic/search',
               params: obj_data,
            }).success(function(data, status) {
                  run = false;
                  //加载更多，则追加内容
                  if (state ==3 ) {
                      $scope.items = $scope.items.concat(data.msg);
                      if (data.msg==null||data.msg.length==0) {
                          console.log("结束");
                          $scope.hasmore=false;
                      }else{
                          obj.pos += obj.count; 
                      }
                  }else{
                      $scope.items = data.msg;
                  } 
              }).error(function(data, status) {
              
              });
        }
    }
})

.controller('ChatDetailCtrl', function($scope, $stateParams, $http, Chats) {
    var run = false;//模拟线程锁机制  防止多次请求 含义：是否正在请求。请注意，此处并非加入到了就绪队列，而是直接跳过不执行
    obj_data = { "topic_id" : $stateParams.chatId };

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
    infoPull(obj_data);
    
    function infoPull(obj_data){
        if(!run){
            run = true;
            $http({
               method:"POST",
               url:'http://www.lancelot.top:80/api/topic/get',
               params: obj_data,
            }).success(function(data, status) {
                $scope.topic = { "title":data.msg.title, "content" : marked(data.msg.content) }
            }).error(function(data, status) {
              
            });
        }
    }
})

.controller('AccountCtrl', function($scope) {
  $scope.settings = {
    enableFriends: true
  };
});
