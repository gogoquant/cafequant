#-*- coding: UTF-8 -*-
'''
    @brief topic manager 
    @author: snack
    @data 2019-08-09
'''

import logging, time, datetime, pdb, simplejson, hashlib

import setting

from tornado.options import define, options
from pycket.driver import Driver

import tornado.gen
import tornado.web
from stormed import Connection
from tornadomail.message import EmailMessage, EmailMultiAlternatives

from util.time_common import month_seconds
from iseecore.routes import route
from services.user import UserService, NoUserException, PasswdErrorException, UserExistsException, UserSameNameException
from views.web.base import WebHandler, WebAsyncAuthHandler, EDITOR_NAME_REMEMBER_COOKIE_KEY, \
    EDITOR_NAME_COOKIE_KEY, EDITOR_LOGIN_TIME_COOKIE_KEY, EDITOR_SESSION_COOKIE_KEY, KEEP_LOGIN_SECONDS
    
#from views.web.base import *

'''
    用户登录
'''
@route(r"/v1/users/login", name="v1.users.login")
class AdminLoginHandler(WebHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _post_(self):

        #获取用户基本信息
        email = self.get_argument("email", None)
        passwd = self.get_argument("passwd", "")
        token = self.get_argument("access_token", None)
        token_from = self.get_argument("source", None)
        openid = self.get_argument("openid", None)
        #remember        = int(self.get_argument('remember',0))
        
        #必要信息缺失则直接错误
        if (not email or not passwd) and (not token or not token_from):
            #self.render_error(status=400, code=400, msg='no email or token')
            self.redirect('/adminweb/adminError', permanent=True)
            return

        
        Driver.EXPIRE_SECONDS = 1 * 24 * 60 * 60   
        self.clear_cookie(EDITOR_NAME_REMEMBER_COOKIE_KEY)
        self.clear_cookie(EDITOR_NAME_COOKIE_KEY)

        user_id, editorinfo = yield tornado.gen.Task(self.user_s.login, email, passwd, token, token_from, openid=openid, app_id=None)
        
        #pdb.set_trace()

        #根据返回值作处理
        if not editorinfo:
            self.render_error(status=400, code=400, msg='no social data')
            return

        if editorinfo == NoUserException:
            self.render_error(status=404, code=404, msg='no user')
            return

        if editorinfo == PasswdErrorException:
            self.render_error(status=403, code=403, msg='passwd error')
            return

        else:
            #self.write(simplejson.dumps(editorinfo))
            #self.set_secure_cookie("user_name", editorinfo.get("name", ""))
        
            #record session
            yield tornado.gen.Task( self.session.set,'editorinfo', editorinfo )
        

            #pdb.set_trace()
            
            #set expires cookie
            login_time = time.time()
            login_expires = datetime.datetime.utcnow() + datetime.timedelta(seconds=KEEP_LOGIN_SECONDS)
            
            #set time cookie
            self.set_cookie(EDITOR_LOGIN_TIME_COOKIE_KEY, "%d" % login_time, expires=login_expires)
            
            #set auth cookie and session rebuilder
            login_session = "%d|%s|%s " % (login_time, email, editorinfo.get('passwd', ''))
            self.set_secure_cookie(EDITOR_SESSION_COOKIE_KEY, login_session, expires_days=None)

            #user id cookie
            self.set_secure_cookie("user_id", user_id)
        
        #self.finish()
        self.redirect('/adminweb/index', permanent=True)


'''
    用户注册
'''
@route(r"/v1/users/register", name="v1.users.register")
class RegisterHandler(WebHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _post_(self):
        
        body = self.request.body
        js = simplejson.loads(body, encoding='utf-8')

        #获取用户数据
        email = js.get("email", None)
        passwd = js.get("passwd", None)
        token = js.get("access_token", None)
        token_from = js.get("source", None)

        #获取用户的第三方认证id
        openid = js.get("openid", None)

        #额外信息
        extra = {}
        extra['sex'] = js.get("sex", 0)
        extra['name'] = js.get("name", "")

        #是否允许随机名
        allow_random_name = js.get("allow_same", None)

        #没有必要信息直接保存
        if (not email or not passwd) and (not token or not token_from):
            logging.error('no data')
            logging.error(email)
            logging.error(token)
            logging.error(token_from)
            self.render_error(status=400, code=400, msg='no data')
            return

        #注册用户
        user_id, return_status = yield tornado.gen.Task(self.user_s.register,\
                email, passwd, token, token_from, extra=extra, app_id=None, \
                openid=openid, allow_random_name=allow_random_name)

        #检查返回值，失败则作错误告警
        if not return_status:
            logging.error('no social data')
            logging.error(token)
            logging.error(token_from)
            self.render_error(status=400, code=400, msg='no social data')
            return

        if return_status == UserExistsException:
            self.render_error(status=409, code=409, msg='exists user')
            return

        if return_status == UserSameNameException:
            self.render_error(status=406, code=406, msg=user_id)
            return

        #成功
        self.write(simplejson.dumps(return_status))

        #record user_id and name
        self.set_secure_cookie("user_id", user_id)
        self.set_secure_cookie("user_name", return_status.get("name", ""))
        self.finish()

'''
    用户更新
'''
@route(r"/v1/users/update", name="v1.users.update")
class UserModifyHandler(WebAsyncAuthHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _post_(self):
        name = self.get_argument("name", None)
        sex = self.get_argument("sex", None)
        age = self.get_argument("age", None)
        address = self.get_argument("address", None)
        description = self.get_argument("description", None)
        #user_id = self.api_user["user_id"]
        user_id = ''

        update_values = {}

        if name:
            update_values['name'] = name
        if sex:
            update_values['sex'] = int(sex)
        if age:
            update_values['age'] = int(age)
        if address:
            update_values['address'] = address
        if description:
            update_values['description'] = description

        if not update_values:
            self.render_error(status=400, code=400, msg='update none')
            return

        result = yield tornado.gen.Task(self.user_s.modify, user_id,
                                        update_values)

        if result == UserExistsException:
            self.render_error(status=400, code=400, msg='exists name')
            return

        self.write(simplejson.dumps(update_values))
        self.finish()

    @tornado.gen.engine
    def _get_(self):
        self._post_()

'''
    用户登出
'''
@route(r"/v1/users/logout", name="v1.users.logout")
class LogoutHandler(WebAsyncAuthHandler):

    @tornado.gen.engine
    def _get_(self):
        user_id = self.get_secure_cookie("user_id")
        user_name = self.get_secure_cookie("user_name")
        logging.error("%s %s logout!!!" % (user_name, user_id))

        #直接通过清理cookie来登出
        self.clear_cookie("user_id")
        self.clear_cookie("user_name")
        self.finish()

    @tornado.gen.engine
    def _post_(self):
        self._get_()


'''
@route(r"/api/user/check_login", name="api.user.check_login")
class CheckLoginHandler(WebAsyncAuthHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _get_(self):
        user_id = self.get_secure_cookie("user_id")
        if not user_id:
            raise tornado.web.HTTPError(404)

        user_name = self.get_secure_cookie("user_name")

        self.set_secure_cookie("user_id", user_id)
        self.set_secure_cookie("user_name", user_name)

        yield tornado.gen.Task(self.user_s.refresh_last_online, user_id)
        self.finish()
'''
'''
@route(r"/api/user/forget_passwd", name="api.user.forget_passwd")
class ForgetPasswdHandler(WebAsyncAuthHandler):
    user_service = UserService()

    @tornado.gen.engine
    def _get_(self):
        email = self.get_argument("email", None)
        if not email:
            raise tornado.web.HTTPError(400)

        query = {
            'email': email
        }
        reset_passwd_id = yield tornado.gen.Task(self.user_service.reset_passwd, query)

        logging.error(reset_passwd_id)

        if reset_passwd_id == 400:
            raise tornado.web.HTTPError(400)

        change_passwd_url = setting.SITE_URL + "/user/change_passwd/%s" % (reset_passwd_id)

        import os.path
        import sys
        find_passwd_path = os.path.join(os.path.dirname(__file__), "../../find_passwd.tpl")
        find_passwd_file = open(find_passwd_path, "rb")
        find_passwd_tpl = find_passwd_file.read()

        send_message = find_passwd_tpl % (
            str("亲爱的用户"),
            app_name,
            str(change_passwd_url),
            app_name
        )

        message = EmailMessage(
            "忘记密码",
            send_message,
            setting.mail_mail,
            [email],
            connection=self.application.mail_connection
        )
        logging.error(message)
        message.send(callback=self.send_finish)

    def send_finish(self, num):
        self.finish()

    @tornado.gen.engine
    def _post_(self):
        self._get_()
'''


@route(r"/v1/users/change_passwd/(.*)", name="v1.users.change_passwd")
class ForgetPasswdHandler(WebAsyncAuthHandler):
    user_service = UserService()
    template = "web/user/user_reset_passwd.html"

    @tornado.gen.engine
    def _get_(self, reset_passwd_id):

        reset_passwd = yield tornado.gen.Task(
            self.user_service.get_reset_passwd_by_id, reset_passwd_id)
        if not reset_passwd:
            raise tornado.web.HTTPError(400)
        response = {'reset_passwd_id': reset_passwd_id}
        response["flg"] = self.get_argument('success', 0)
        self.render(**response)

    @tornado.gen.engine
    def _post_(self, reset_passwd_id):

        reset_passwd = yield tornado.gen.Task(
            self.user_service.get_reset_passwd_by_id, reset_passwd_id)
        if not reset_passwd:
            self.render_error(msg=self._('No reset passwd session'), status=203)
            return

        email = reset_passwd["email"]

        user_password = self.get_argument("password", None)
        comfirm_password = self.get_argument('password2', None)

        if not user_password or not comfirm_password:
            self.render_error(msg=self._('No new password!'), status=203)
            return

        if user_password != comfirm_password:
            self.render_error(
                msg=self._('Confirm password is different from new password!'), status=203)
            return

        yield tornado.gen.Task(self.user_service.update_user_passwd,
                               reset_passwd_id, email, user_password)

        self.finish()

'''
    用户信息
'''
@route(r"/v1/users/info/(.+)", name="v1.users.info")
class UserInfoHandler(WebAsyncAuthHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _get_(self, user_id):
        self_user_id = self.get_secure_cookie("user_id")
        user_info = yield tornado.gen.Task(
            self.user_s.info, None, user_id, self_user_id, viewMore=True)

        if not user_info:
            logging.error('no user')
            self.render_error(status=400, code=400, msg='no user')
            return

        self.write(simplejson.dumps(user_info))
        self.finish()


'''
    用户列表
'''
@route(r"/v1/users/list", name="v1.users.list")
class UserSearchHandler(WebHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _post_(self):
        #pdb.set_trace()

        token = self.get_argument("token", None)
        page = self.get_argument("page", None)
        pagesize = self.get_argument("pagesize", None)

        if page is None:
            page_num = int(0)
        else:
            page_num = int(page)

        if pagesize is None:
            pagesize_num = int(0)
        else:
            pagesize_num = int(pagesize)

        users = yield tornado.gen.Task(self.user_s.search, page_num,
                                       pagesize_num, token)

        #pdb.set_trace()

        self.render_success(msg=users)
        return


@route(r"/v1/users/tot", name="v1.users.tot")
class UserTotHandler(WebHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _post_(self):
        #pdb.set_trace()

        token = self.get_argument("token", None)

        users = yield tornado.gen.Task(self.user_s.count, token)

        self.render_success(msg=users)
        return


@route(r"/v1/users/delete", name="v1.users.delete")
class UserDeleteHandler(WebAsyncAuthHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _post_(self):
        #pdb.set_trace()

        user_id = self.get_argument("user_id", "")

        yield tornado.gen.Task(self.user_s.delete, user_id)

        self.render_success(msg="success")
        return
