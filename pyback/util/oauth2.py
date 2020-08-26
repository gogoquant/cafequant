#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
Oauth2 协议封装，使用异步httpclient
支持weibo,qq,renren,douban

client_id   true    string  申请应用时分配的AppKey。
client_secret   true    string  申请应用时分配的AppSecret。
grant_type  true    string  请求的类型，填写authorization_code

grant_type为authorization_code时
    必选  类型及范围   说明
code    true    string  调用authorize获得的code值。
redirect_uri    true    string  回调地址，需需与注册应用里的回调地址一致。

curl -v -d "client_id=320064844&client_secret=6ca10df52bdc5de775d7c16dd6593daf&grant_type=authorization_code&code=f8283ac0af2b376dc4ae7736780c4a1c&redirect_uri=http://www.idealsee.com/api/oauth2/callback" https://api.weibo.com/oauth2/access_token
@file:oauth2.py
@modul:oauth2
@author:xiaolin@idealsee.cn
@date:2014-09-10
"""
import tornado.gen
import simplejson
import logging
from functools import wraps
from urllib import urlencode
from new import classobj
from tornado.httpclient import AsyncHTTPClient, HTTPRequest, HTTPError
from tornado.httputil import HTTPHeaders

import sys
sys.path.append("../")
import setting
import mimetypes

reload(sys)
sys.setdefaultencoding("utf-8")

WEIBO_T = "weibo"
QQ_T = "qq"
WEBCHAT_T = "wechat"

# yowo.idealsee.com APP_ID
WEBCHAT_OPEN_APP_ID = "wx79124388947c0897"
WEBCHAT_OPEN_SECRET = "07b12843a0c78c514104754c3d10e4d0"

# yowo.idealsee.com APP_ID
WEIBO_APPKEY = "453970566"
WEIBO_SECRETKEY = "fea32b2ee419b4a7c0d955ab5861aad1"
WEIBO_DEFAULT_REDIRECT_URI = "http://yowo.idealsee.com"

AsyncHTTPClient.configure("tornado.curl_httpclient.CurlAsyncHTTPClient")


def encode_multipart_formdata(fields, files):
    """
    fields is a sequence of (name, value) elements for regular form fields.
    files is a sequence of (name, filename, value) elements for data to be uploaded as files
    Return (content_type, body) ready for httplib.HTTP instance
    """
    BOUNDARY = '----------ThIs_Is_tHe_bouNdaRY_$'
    CRLF = '\r\n'
    L = []
    for (key, value) in fields:
        L.append('--' + BOUNDARY)
        L.append('Content-Disposition: form-data; name="%s"' % key)
        L.append('')
        L.append(value)
    for (key, filename, value) in files:
        L.append('--' + BOUNDARY)
        L.append('Content-Disposition: form-data; name="%s"; filename="%s"' %
                 (key, filename))
        L.append('Content-Type: %s' % get_content_type(filename))
        L.append('')
        L.append(value)
    L.append('--' + BOUNDARY + '--')
    L.append('')
    body = CRLF.join(L)
    content_type = 'multipart/form-data; boundary=%s' % BOUNDARY
    return content_type, body


def get_content_type(filename):
    return mimetypes.guess_type(filename)[0] or 'application/octet-stream'


class Oauth2Error(Exception):
    """Occurred when doing API call"""

    def __init__(self, site_name, url, error_msg, *args):
        self.site_name = site_name
        self.url = url
        self.error_msg = error_msg
        Exception.__init__(self, error_msg, *args)


def _http_error_handler(func):

    @wraps(func)
    def deco(self, *args, **kwargs):
        try:
            res = func(self, *args, **kwargs)
        except HTTPError as e:
            raise Oauth2Error(self.site_name, e.url, e.read())
        except URLError as e:
            raise Oauth2Error(self.site_name, args[0], e.reason)

        error_key = getattr(self, 'RESPONSE_ERROR_KEY', None)
        if error_key is not None and error_key in res:
            raise Oauth2Error(self.site_name, args[0], res)

        return res

    return deco


class OAuth2(object):
    """Base OAuth2 class, Sub class must define the following settings:

    AUTHORIZE_URL    - Asking user to authorize and get token
    ACCESS_TOKEN_URL - Get authorized access token

    And the bellowing should define in settings file

    REDIRECT_URI     - The url after user authorized and redirect to
    CLIENT_ID        - Your client id for the social site
    CLIENT_SECRET    - Your client secret for the social site

    Also, If the Website needs scope parameters, your should add it too.

    SCOPE            - A list type contains some scopes

    Details see: http://tools.ietf.org/html/rfc6749


    SubClass MUST Implement the following three methods:

    build_api_url(self, url)
    build_api_data(self, **kwargs)
    parse_token_response(self, res)
    """

    def __init__(self):
        self.proxy_host = None
        self.proxy_port = None

    @classmethod
    def instance(cls, social_type):
        if social_type == WEIBO_T:
            return Weibo()
        elif social_type == QQ_T:
            return QQ()
        elif social_type == WEBCHAT_T:
            return Webchat()

    def set_proxy(self, host, port):
        self.proxy_host = host
        self.proxy_port = port

    @tornado.gen.engine
    @_http_error_handler
    def http_get(self, url, data="", parse=True, callback=None):
        http_client = AsyncHTTPClient()
        headers = self.get_headers()
        http_request = HTTPRequest(
            '%s?%s' % (url, urlencode(data)),
            headers=headers,
            proxy_host=self.proxy_host,
            proxy_port=self.proxy_port)
        response = yield tornado.gen.Task(http_client.fetch, http_request)
        if response.error:
            logging.error("Error:%s", response.error)
            logging.error("Error_body:%s", response.body)
            callback(None)
        else:
            response = response.body
            if parse:
                try:
                    response = simplejson.loads(response)
                except Exception, e:
                    response = None
            callback(response)

    @tornado.gen.engine
    @_http_error_handler
    def http_post(self, url, data="", parse=True, callback=None):
        http_client = AsyncHTTPClient()
        headers = self.get_headers()
        http_request = HTTPRequest(
            url,
            method="POST",
            headers=headers,
            body=urlencode(data),
            proxy_host=self.proxy_host,
            proxy_port=self.proxy_port)
        response = yield tornado.gen.Task(http_client.fetch, http_request)
        if response.error:
            logging.error("Error:%s", response.error)
            logging.error("Error_body:%s", response.body)
            callback(None)
        else:
            response = response.body
            if parse:
                try:
                    response = simplejson.loads(response)
                except Exception, e:
                    response = None
            callback(response)

    @tornado.gen.engine
    @_http_error_handler
    def http_upload_post(self, url, data="", parse=True, callback=None):
        http_client = AsyncHTTPClient()

        fields = []
        files = []
        for key, value in data.items():
            if isinstance(value, file):
                files.append([key, value.name, value.read()])
            else:
                fields.append([key, value])
        content_type, body = encode_multipart_formdata(fields, files)
        headers = HTTPHeaders({"content-type": content_type})

        http_request = HTTPRequest(
            url,
            method="POST",
            headers=headers,
            body=body,
            proxy_host=self.proxy_host,
            proxy_port=self.proxy_port)
        print http_request
        response = yield tornado.gen.Task(http_client.fetch, http_request)
        if response.error:
            logging.error("Error:%s", response.error)
            logging.error("Error_body:%s", response.body)
            callback(None)
        else:
            response = response.body
            if parse:
                try:
                    response = simplejson.loads(response)
                except Exception, e:
                    response = None
            callback(response)

    def get_headers(self):
        """Sub class rewiter this function If it's necessary to add headers"""
        return None


class Weibo(OAuth2):
    GET_TOKEN_INFO_URL = "https://api.weibo.com/oauth2/get_token_info"
    GET_USER_INFO_URL = "https://api.weibo.com/2/users/show.json"
    GET_USER_FOLLOW_URL = "https://api.weibo.com/2/friendships/friends.json"
    POST_WEIBO_URL = "https://upload.api.weibo.com/2/statuses/upload.json"
    # POST_WEIBO_URL = "https://api.weibo.com/2/statuses/update.json"
    # POST_WEIBO_URL = "https://api.weibo.com/2/statuses/upload_url_text.json"
    GET_TOKEN_FROM_CODE_URL = "https://api.weibo.com/oauth2/access_token"
    """
    http://open.weibo.com/wiki/Oauth2/get_token_info

   oauth2/get_token_info
    查询用户access_token的授权相关信息，包括授权时间，过期时间和scope权限。
    URL
    https://api.weibo.com/oauth2/get_token_info
    HTTP请求方式
    POST
    请求参数
    access_token：用户授权时生成的access_token。
    返回数据
     {
           "uid": 1073880650,
           "appkey": 1352222456,
           "scope": null,
           "create_at": 1352267591,
           "expire_in": 157679471
     }

    返回值字段   字段类型    字段说明
    uid string  授权用户的uid。
    appkey  string  access_token所属的应用appkey。
    scope   string  用户授权的scope权限。
    create_at   string  access_token的创建时间，从1970年到创建时间的秒数。
    expire_in   string  access_token的剩余时间，单位是秒数。
    """

    @tornado.gen.engine
    def get_token_info(self, access_token, callback):
        post_data = {'access_token': access_token}
        token_info = yield tornado.gen.Task(self.http_post,
                                            self.GET_TOKEN_INFO_URL, post_data)
        callback(token_info)

    """
    http://open.weibo.com/wiki/2/users/show

    users/show
    根据用户ID获取用户信息
    URL
    https://api.weibo.com/2/users/show.json
    支持格式
    JSON
    HTTP请求方式
    GET
    是否需要登录
    是
    关于登录授权，参见 如何登录授权
    访问授权限制
    访问级别：普通接口
    频次限制：是
    关于频次限制，参见 接口访问权限说明
    请求参数
        必选  类型及范围   说明
    source  false   string  采用OAuth授权方式不需要此参数，其他授权方式为必填参数，数值为应用的AppKey。
    access_token    false   string  采用OAuth授权方式为必填参数，其他授权方式不需要此参数，OAuth授权后获得。
    uid false   int64   需要查询的用户ID。
    screen_name false   string  需要查询的用户昵称。
    注意事项
    参数uid与screen_name二者必选其一，且只能选其一
    调用样例及调试工具
    API测试工具
    返回结果
    JSON示例
    {
        "id": 1404376560,
        "screen_name": "zaku",
        "name": "zaku",
        "province": "11",
        "city": "5",
        "location": "北京 朝阳区",
        "description": "人生五十年，乃如梦如幻；有生斯有死，壮士复何憾。",
        "url": "http://blog.sina.com.cn/zaku",
        "profile_image_url": "http://tp1.sinaimg.cn/1404376560/50/0/1",
        "domain": "zaku",
        "gender": "m",
        "followers_count": 1204,
        "friends_count": 447,
        "statuses_count": 2908,
        "favourites_count": 0,
        "created_at": "Fri Aug 28 00:00:00 +0800 2009",
        "following": false,
        "allow_all_act_msg": false,
        "geo_enabled": true,
        "verified": false,
        "status": {
            "created_at": "Tue May 24 18:04:53 +0800 2011",
            "id": 11142488790,
            "text": "我的相机到了。",
            "source": "<a href="http://weibo.com" rel="nofollow">新浪微博</a>",
            "favorited": false,
            "truncated": false,
            "in_reply_to_status_id": "",
            "in_reply_to_user_id": "",
            "in_reply_to_screen_name": "",
            "geo": null,
            "mid": "5610221544300749636",
            "annotations": [],
            "reposts_count": 5,
            "comments_count": 8
        },
        "allow_all_comment": true,
        "avatar_large": "http://tp1.sinaimg.cn/1404376560/180/0/1",
        "verified_reason": "",
        "follow_me": false,
        "online_status": 0,
        "bi_followers_count": 215
    }

    关于错误返回值与错误代码，参见 错误代码说明
    返回字段说明
    返回值字段   字段类型    字段说明
    id  int64   用户UID
    idstr   string  字符串型的用户UID
    screen_name string  用户昵称
    name    string  友好显示名称
    province    int 用户所在省级ID
    city    int 用户所在城市ID
    location    string  用户所在地
    description string  用户个人描述
    url string  用户博客地址
    profile_image_url   string  用户头像地址（中图），50×50像素
    profile_url string  用户的微博统一URL地址
    domain  string  用户的个性化域名
    weihao  string  用户的微号
    gender  string  性别，m：男、f：女、n：未知
    followers_count int 粉丝数
    friends_count   int 关注数
    statuses_count  int 微博数
    favourites_count    int 收藏数
    created_at  string  用户创建（注册）时间
    following   boolean 暂未支持
    allow_all_act_msg   boolean 是否允许所有人给我发私信，true：是，false：否
    geo_enabled boolean 是否允许标识用户的地理位置，true：是，false：否
    verified    boolean 是否是微博认证用户，即加V用户，true：是，false：否
    verified_type   int 暂未支持
    remark  string  用户备注信息，只有在查询用户关系时才返回此字段
    status  object  用户的最近一条微博信息字段 详细
    allow_all_comment   boolean 是否允许所有人对我的微博进行评论，true：是，false：否
    avatar_large    string  用户头像地址（大图），180×180像素
    avatar_hd   string  用户头像地址（高清），高清头像原图
    verified_reason string  认证原因
    follow_me   boolean 该用户是否关注当前登录用户，true：是，false：否
    online_status   int 用户的在线状态，0：不在线、1：在线
    bi_followers_count  int 用户的互粉数
    lang    string  用户当前的语言版本，zh-cn：简体中文，zh-tw：繁体中文，en：英语
    """

    @tornado.gen.engine
    def get_user_info(self, access_token, uid, callback):
        post_data = {'access_token': access_token, 'uid': uid}
        userinfo = yield tornado.gen.Task(self.http_get, self.GET_USER_INFO_URL,
                                          post_data)
        callback(userinfo)

    """
    获取用户的关注列表

    请求参数
    必选  类型及范围   说明
    source  false   string  采用OAuth授权方式不需要此参数，其他授权方式为必填参数，数值为应用的AppKey。
    access_token    false   string  采用OAuth授权方式为必填参数，其他授权方式不需要此参数，OAuth授权后获得。
    uid false   int64   需要查询的用户UID。
    screen_name false   string  需要查询的用户昵称。
    count   false   int 单页返回的记录条数，默认为50，最大不超过200。
    cursor  false   int 返回结果的游标，下一页用返回值里的next_cursor，上一页用previous_cursor，默认为0。
    trim_status false   int 返回值中user字段中的status字段开关，0：返回完整status字段、1：status字段仅返回status_id，默认为1。
    
    返回结果
    JSON示例
    {
        "users": [
            {
                "id": 1404376560,
                "screen_name": "zaku",
                "name": "zaku",
                "province": "11",
                "city": "5",
                "location": "北京 朝阳区",
                "description": "人生五十年，乃如梦如幻；有生斯有死，壮士复何憾。",
                "url": "http://blog.sina.com.cn/zaku",
                "profile_image_url": "http://tp1.sinaimg.cn/1404376560/50/0/1",
                "domain": "zaku",
                "gender": "m",
                "followers_count": 1204,
                "friends_count": 447,
                "statuses_count": 2908,
                "favourites_count": 0,
                "created_at": "Fri Aug 28 00:00:00 +0800 2009",
                "following": false,
                "allow_all_act_msg": false,
                "remark": "",
                "geo_enabled": true,
                "verified": false,
                "status": {
                    "created_at": "Tue May 24 18:04:53 +0800 2011",
                    "id": 11142488790,
                    "text": "我的相机到了。",
                    "source": "<a href="http://weibo.com" rel="nofollow">新浪微博</a>",
                    "favorited": false,
                    "truncated": false,
                    "in_reply_to_status_id": "",
                    "in_reply_to_user_id": "",
                    "in_reply_to_screen_name": "",
                    "geo": null,
                    "mid": "5610221544300749636",
                    "annotations": [],
                    "reposts_count": 5,
                    "comments_count": 8
                },
                "allow_all_comment": true,
                "avatar_large": "http://tp1.sinaimg.cn/1404376560/180/0/1",
                "verified_reason": "",
                "follow_me": false,
                "online_status": 0,
                "bi_followers_count": 215
            },
            ...
        ],
        "next_cursor": 5,
        "previous_cursor": 0,
        "total_number": 668
    }

    """

    @tornado.gen.engine
    def get_user_follow(self, access_token, uid, callback):
        post_data = {'access_token': access_token, 'uid': uid}
        userinfo = yield tornado.gen.Task(self.http_get,
                                          self.GET_USER_FOLLOW_URL, post_data)
        callback(userinfo)

    @tornado.gen.engine
    def post_weibo_with_pic(self, access_token, text, file_tmp, callback):
        f = open(file_tmp)
        post_data = {
            'access_token': str(access_token),
            'status': text,
            'pic': f
        }
        userinfo = yield tornado.gen.Task(self.http_upload_post,
                                          self.POST_WEIBO_URL, post_data)
        f.close()
        callback(userinfo)

    """
    oauth2/access_token
    OAuth2的access_token接口
    URL
    https://api.weibo.com/oauth2/access_token
    HTTP请求方式
    POST
    请求参数
        必选  类型及范围   说明
    client_id   true    string  申请应用时分配的AppKey。
    client_secret   true    string  申请应用时分配的AppSecret。
    grant_type  true    string  请求的类型，填写authorization_code

    grant_type为authorization_code时
        必选  类型及范围   说明
    code    true    string  调用authorize获得的code值。
    redirect_uri    true    string  回调地址，需需与注册应用里的回调地址一致。

    返回数据
     {
           "access_token": "ACCESS_TOKEN",
           "expires_in": 1234,
           "remind_in":"798114",
           "uid":"12341234"
     }

    返回值字段   字段类型    字段说明
    access_token    string  用于调用access_token，接口获取授权后的access token。
    expires_in  string  access_token的生命周期，单位是秒数。
    remind_in   string  access_token的生命周期（该参数即将废弃，开发者请使用expires_in）。
    uid string  当前授权用户的UID。
    """

    @tornado.gen.engine
    def get_token_info_from_code(self, code, callback):
        post_data = {
            "client_id": WEIBO_APPKEY,
            "client_secret": WEIBO_SECRETKEY,
            "code": code,
            "grant_type": "authorization_code",
            "redirect_uri": WEIBO_DEFAULT_REDIRECT_URI
        }
        token_info = yield tornado.gen.Task(self.http_post,
                                            self.GET_TOKEN_FROM_CODE_URL,
                                            post_data)
        callback(token_info)


class QQ(OAuth2):
    GET_TOKEN_INFO_URL = "https://graph.qq.com/oauth2.0/me"
    GET_USER_INFO_URL = "https://graph.qq.com/user/get_user_info"
    """
    http://wiki.connect.qq.com/%E8%8E%B7%E5%8F%96%E7%94%A8%E6%88%B7openid_oauth2-0
    1 请求地址

    PC网站：https://graph.qq.com/oauth2.0/me
    WAP网站：https://graph.z.qq.com/moc2/me
    2 请求方法

    GET
    3 请求参数

    请求参数请包含如下内容：
    参数  是否必须    含义
    access_token    必须  在Step1中获取到的access token。

    4 返回说明

    PC网站接入时，获取到用户OpenID，返回包如下：
    1
    callback( {"client_id":"YOUR_APPID","openid":"YOUR_OPENID"} );
    """

    @tornado.gen.engine
    def get_token_info(self, access_token, callback):
        post_data = {'access_token': access_token}
        token_info = yield tornado.gen.Task(
            self.http_get, self.GET_TOKEN_INFO_URL, post_data, parse=False)
        if 'callback(' in token_info:
            token_info = token_info[token_info.index('(') +
                                    1:token_info.rindex(')')]
            token_info = simplejson.loads(token_info)
        else:
            token_info = token_info.split('&')
            token_info = [_r.split('=') for _r in token_info]
            token_info = dict(token_info)
        callback(token_info)

    """
    http://wiki.connect.qq.com/openapi%E8%B0%83%E7%94%A8%E8%AF%B4%E6%98%8E_oauth2-0
    参数  含义
    access_token    可通过使用Authorization_Code获取Access_Token 或来获取。
    access_token有3个月有效期。
    oauth_consumer_key  申请QQ登录成功后，分配给应用的appid
    openid  用户的ID，与QQ号码一一对应。
    可通过调用https://graph.qq.com/oauth2.0/me?access_token=YOUR_ACCESS_TOKEN 来获取。

    http://wiki.connect.qq.com/get_user_info
    3.4返回参数说明
        参数说明    描述
        ret 返回码
        msg 如果ret<0，会有相应的错误信息提示，返回数据全部用UTF-8编码。
        nickname    用户在QQ空间的昵称。
        figureurl   大小为30×30像素的QQ空间头像URL。
        figureurl_1 大小为50×50像素的QQ空间头像URL。
        figureurl_2 大小为100×100像素的QQ空间头像URL。
        figureurl_qq_1  大小为40×40像素的QQ头像URL。
        figureurl_qq_2  大小为100×100像素的QQ头像URL。需要注意，不是所有的用户都拥有QQ的100x100的头像，但40x40像素则是一定会有。
        gender  性别。 如果获取不到则默认返回"男"
        is_yellow_vip   标识用户是否为黄钻用户（0：不是；1：是）。
        vip 标识用户是否为黄钻用户（0：不是；1：是）
        yellow_vip_level    黄钻等级
        level   黄钻等级
        is_yellow_year_vip  标识是否为年费黄钻用户（0：不是； 1：是）
    """

    @tornado.gen.engine
    def get_user_info(self, access_token, oauth_consumer_key, openid, callback):
        post_data = {
            'access_token': access_token,
            'oauth_consumer_key': oauth_consumer_key,
            'openid': openid
        }
        userinfo = yield tornado.gen.Task(self.http_get, self.GET_USER_INFO_URL,
                                          post_data)
        callback(userinfo)


class Webchat(OAuth2):
    # GET_TOKEN_INFO_URL = "https://graph.qq.com/oauth2.0/me"
    GET_USER_INFO_URL = "https://api.weixin.qq.com/sns/userinfo"
    GET_TOKEN_FROM_CODE_URL = "https://api.weixin.qq.com/sns/oauth2/access_token"
    """
    https://open.weixin.qq.com/cgi-bin/frame?t=resource/res_main_tmpl
    参数  含义
    access_token    通过code获取access_token 刷新或续期access_token使用 
    openid  普通用户的标识，对当前开发者帐号唯一，获取该openid的好友列表
    通过code获取access_token 刷新或续期access_token使用

    返回参数说明
        openid  普通用户的标识，对当前开发者帐号唯一
        nickname    普通用户昵称
        sex 普通用户性别，1为男性，2为女性
        province    普通用户个人资料填写的省份
        city    普通用户个人资料填写的城市
        country 国家，如中国为CN
        headimgurl  用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空
        privilege   用户特权信息，json数组，如微信沃卡用户为（chinaunicom）
        unionid 用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的unionid是唯一的。

        {
        'province': 'Auckland', 
        'openid': 'oIK-pjmHPT4Jd-1p0V6lTqTfoHQY', 
        'headimgurl': 'http://wx.qlogo.cn/mmopen/GSoEvnap96LS4fwPiayyDgmUrkdXNoSpKib4xCwRCns23jc8XSbCQbBEktpbiaLQTPvPCSJvWN55ChcTHqulRrBxxuUoDACGMnD/0', 
        'language': 'en', 
        'city': '', 
        'country': 'NZ', 
        'sex': 1, 
        'unionid': 'ozdGCuE7TsWUrRHQSpxT9m7ucK14', 
        'privilege': [], 
        'nickname': u'\u7b28\u9ed1\u5154'
        }

    """

    @tornado.gen.engine
    def get_user_info(self, access_token, openid, callback):
        post_data = {'access_token': access_token, 'openid': openid}
        userinfo = yield tornado.gen.Task(self.http_get, self.GET_USER_INFO_URL,
                                          post_data)
        callback(userinfo)

    """
    通过code获取access_token

    请求说明
    http请求方式: GET
    https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
    参数说明
    参数  是否必须    说明
    appid   是   应用唯一标识，在微信开放平台提交应用审核通过后获得
    secret  是   应用密钥AppSecret，在微信开放平台提交应用审核通过后获得
    code    是   填写第一步获取的code参数
    grant_type  是   填authorization_code

    返回说明
    正确的返回：
    { 
    "access_token":"ACCESS_TOKEN", 
    "expires_in":7200, 
    "refresh_token":"REFRESH_TOKEN",
    "openid":"OPENID", 
    "scope":"SCOPE" 
    }
    参数  说明
    access_token    接口调用凭证
    expires_in  access_token接口调用凭证超时时间，单位（秒）
    refresh_token   用户刷新access_token
    openid  授权用户唯一标识
    scope   用户授权的作用域，使用逗号（,）分隔
    错误返回样例：
    {
    "errcode":40029,"errmsg":"invalid code"
    }
    """

    @tornado.gen.engine
    def get_token_info(self, code, callback):
        post_data = {
            "appid": WEBCHAT_OPEN_APP_ID,
            "secret": WEBCHAT_OPEN_SECRET,
            "code": code,
            "grant_type": "authorization_code"
        }
        token_info = yield tornado.gen.Task(self.http_get,
                                            self.GET_TOKEN_FROM_CODE_URL,
                                            post_data)
        callback(token_info)


if __name__ == '__main__':
    logging.basicConfig(
        format='%(asctime)s %(filename)s:%(lineno)d %(levelname)s %(message)s',
        level=logging.INFO)
    import tornado.ioloop

    # oauth2 = OAuth2.instance('qq')
    # if hasattr(setting,"proxy_host"):
    #     oauth2.set_proxy("10.0.1.32",31287)
    # def get_user_info_response(response):
    #     logging.info(response)
    # def callback(response):
    #     logging.info(response)
    #     # get_user info
    #     uid = response.get('openid', '')
    #     appkey = response.get('client_id','')
    #     oauth2.get_user_info(access_token,appkey, uid,get_user_info_response)
    # access_token = '3DC1C24331F6468C20D0F808AAB5A652'

    oauth2 = OAuth2.instance('weibo')
    if hasattr(setting, "proxy_host"):
        oauth2.set_proxy("10.0.1.32", 31287)

    def get_user_info_response(response):
        logging.info(response)

    def callback(response):
        logging.info(response)
        # uid = response.get('uid', '')
        # oauth2.get_user_info(access_token, uid,get_user_info_response)

    access_token = "2.003Xt_pCanOfgCd6dcf0a39efEYWCE"
    # access_token = "OezXcEiiBSKSxW0eoylIeAtwmYQypM5splLYSpYTUbOtIR7ftwsTo5O2axcSedADT7igkMZf8UicDKQirp4TIgaPUYopGN267v3YY0KUdvyk4Ri31YIeV285rtuBxRbDNKlkV2ENzbS5i_sF7c19GA"
    # openid = "oIK-pjmHPT4Jd-1p0V6lTqTfoHQY"

    # oauth2.get_token_info(access_token,callback)
    # oauth2.get_user_info(access_token, openid, callback)
    oauth2.post_weibo_with_pic(access_token, 'test', '../3.jpg', callback)
    tornado.ioloop.IOLoop.instance().start()
