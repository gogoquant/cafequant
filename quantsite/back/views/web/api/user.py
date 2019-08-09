#-*- coding: UTF-8 -*-

'''
    @brief topic manager 
    @author: snack
    @data 2019-08-09
'''

import logging
import tornado.gen
import tornado.web
import pdb
import setting

from tornado.options import define, options
from iseecore.routes  import route
from views.web.base import WebHandler, WebAsyncAuthHandler
#from views.web.base import *

from services.user import UserService, NoUserException, PasswdErrorException, UserExistsException, UserSameNameException

from stormed import Connection
from tornadomail.message import EmailMessage, EmailMultiAlternatives
import simplejson
import hashlib

import setting
from util.time_common import month_seconds


'''
    用户注册句柄
'''
@route(r"/api/user/register", name="api.user.register")
class RegisterHandler(WebHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _post_(self):

        #获取用户的基本信息
        email = self.get_argument("email", None)
        passwd = self.get_argument("passwd", None)
        token = self.get_argument("access_token", None)
        token_from = self.get_argument("source", None)
        
        #获取用户的第三方认证id
        openid = self.get_argument("openid", None)
        
        #额外信息
        extra = {}
        extra['sex'] = self.get_argument("sex", 0)
        extra['name'] = self.get_argument("name", "")
        
        #是否允许随机名
        allow_random_name = self.get_argument("allow_same", None)

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

@route(r"/api/user/update", name="user.update")
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

        result = yield tornado.gen.Task(self.user_s.modify, user_id, update_values)

        if result == UserExistsException:
            self.render_error(status=400, code=400, msg='exists name')
            return

        self.write(simplejson.dumps(update_values))
        self.finish()

    @tornado.gen.engine
    def _get_(self):
        self._post_()

@route(r"/api/user/logout", name="api.user.logout")
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

@route(r"/api/user/change_passwd/(.*)", name="api.user.change_passwd")
class ForgetPasswdHandler(WebAsyncAuthHandler):
    user_service = UserService()
    template = "web/user/user_reset_passwd.html"

    @tornado.gen.engine
    def _get_(self, reset_passwd_id):

        reset_passwd = yield tornado.gen.Task(self.user_service.get_reset_passwd_by_id, reset_passwd_id)
        # if not reset_passwd:
        #     raise tornado.web.HTTPError(400)
        response = {'reset_passwd_id': reset_passwd_id}
        response["flg"] = self.get_argument('success', 0)
        self.render(**response)

    @tornado.gen.engine
    def _post_(self, reset_passwd_id):

        reset_passwd = yield tornado.gen.Task(self.user_service.get_reset_passwd_by_id, reset_passwd_id)
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

        yield tornado.gen.Task(self.user_service.update_user_passwd, reset_passwd_id, email, user_password)

        self.finish()

@route(r"/api/user/info/(.+)", name="api.user.info")
class UserInfoHandler(WebAsyncAuthHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _get_(self, user_id):
        self_user_id = self.get_secure_cookie("user_id")
        user_info = yield tornado.gen.Task(self.user_s.info, None, user_id, self_user_id, viewMore=True)

        if not user_info:
            logging.error('no user')
            self.render_error(status=400, code=400, msg='no user')
            return

        self.write(simplejson.dumps(user_info))
        self.finish()


@route(r"/api/user/search", name="api.user.search")
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

        users = yield tornado.gen.Task(self.user_s.search, page_num, pagesize_num, token)
        
        #pdb.set_trace()

        self.render_success(msg=users)
        return


@route(r"/api/user/tot", name="api.user.tot")
class UserTotHandler(WebHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _post_(self):
        #pdb.set_trace()

        token = self.get_argument("token", None)

        users = yield tornado.gen.Task(self.user_s.count, token)
        
        self.render_success(msg=users)
        return


@route(r"/api/user/delete", name="api.user.delete")
class UserDeleteHandler(WebAsyncAuthHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _post_(self):
        #pdb.set_trace()

        user_id = self.get_argument("user_id", "")

        yield tornado.gen.Task(self.user_s.delete, user_id)
        
        self.render_success(msg="success")
        return

