
#varnish to nginx
backend cqserver {
    .host = "127.0.0.1";
    .port = "8080";
    .connect_timeout = 20s;
}

#purge allow
acl purge {
    "localhost";
    "127.0.0.1";
}

#acl
sub vcl_recv {
    if (req.request == "PURGE") {
        if (!client.ip ~ purge) {
            error 405 "Not allowed.";
        }
        return (lookup);
    }
    if (req.http.host ~ "^www.lancelot.top") {
        set req.backend = cqserver;
        if (req.request != "GET" && req.request != "HEAD") {
            return (pipe);
        }else{
            return (lookup);
        }
    }else {
        error 404 "caoqing Cache Server";
        return (lookup);
    }
}

#hit trigger
sub vcl_hit {
    if (req.request == "PURGE") {
        set obj.ttl = 0s;
        error 200 "Purged.";
    }
}

#miss trigger
sub vcl_miss {
    if (req.request == "PURGE") {
        error 404 "Not in cache.";
    }
}

#(1)Varnish通过反向代理请求后端IP为127.0.0.1，端口为8087的web服务器,即nginx服务器监听端口；
#(2)Varnish允许localhost、127.0.0.1、192.168.1.*三个来源IP通过PURGE方法清除缓存；
#(3)Varnish对域名为www.baidu.com的请求进行处理，非www.baidu.com域名的请求则返回"caoqing Cache Server"；
#(4)Varnish对HTTP协议中的GET、HEAD请求进行缓存，对POST请求透过，让其直接访问后端Web服务器。
