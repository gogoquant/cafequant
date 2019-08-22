# -*- coding: utf-8 -*-
'''
    @brief handler base
    @author: xiyanxiyan10
    @data 2016-12-13
    @note http 句柄的相关接口封装，主要是集成鉴权机制，以及用户数据重拉
'''

import time
from permission import Permission
from pycket.session import SessionManager
import simplejson
import urlparse
import tornado.web
import tornado.gen
import logging
import datetime

from iseecore.mixins import ExceptionMixin
from views.base import RequestHandler
from services.user import UserService

# from services.nav import Nav
import setting

#登录时间cookie记录
EDITOR_LOGIN_TIME_COOKIE_KEY = 'login_time_cookie'
#session 的cookie记录
EDITOR_SESSION_COOKIE_KEY = 'editor_session_cookie'
#login record
EDITOR_NAME_REMEMBER_COOKIE_KEY = "login_remember"
EDITOR_NAME_COOKIE_KEY = "login_name"

#cookie 超时时间
KEEP_LOGIN_SECONDS = 60 * 60


"""
    WebHandler
    with cdn support 
    with get post pkg 
"""
class WebHandler(RequestHandler):

    def static_url(self, path, use_cdn=None):
        '''get real address, Cdn maybe used'''

        self.require_setting("static_path", "static_url")

        static_handler_class = self.settings.get("static_handler_class",
                                                 tornado.web.StaticFileHandler)

        if use_cdn is None:
            use_cdn = setting.IS_CDN

        if use_cdn:
            base = setting.STATIC_CDN_URL
        else:
            base = ""

        return base + static_handler_class.make_static_url(self.settings, path)

    def set_default_headers(self):
        '''set http header'''
        self.set_header('Access-Control-Allow-Origin', '*')
        self.set_header('Access-Control-Allow-Methods', 'POST, GET, OPTIONS')
        self.set_header('Access-Control-Max-Age', 1000)
        self.set_header('Access-Control-Allow-Headers', '*')
        #self.set_header('Content-type', 'application/json')

    def render(self, **kwargs):
        '''响应封装'''

        #设置导航栏信息
        #nav_active = self.__class__.nav_active if hasattr(self.__class__, 'nav_active') else None
        #kwargs['navbar_list'] = self.navbar_list
        #kwargs['nav_active'] = nav_active
        #self.response.update(kwargs)

        RequestHandler.render(self, **self.response)

    def render_error(self, status=400, code=400, msg=''):
        '''返回错误提示'''
        error = {'status': 'error', 'code': code, 'msg': msg}
        self.set_header("Content-Type", 'application/json')
        self.set_status(status)
        self.write(simplejson.dumps(error))
        self.finish()

    def render_success(self, status=200, code=200, msg=''):
        '''返回api'''
        result = {'status': 'success', 'code': code, 'msg': msg}
        self.set_header("Content-Type", 'application/json')
        self.set_status(status)
        self.write(simplejson.dumps(result))
        self.finish()

    @tornado.web.asynchronous
    @tornado.gen.engine
    def post(self, *args, **kwargs):
        '''post封装'''
        if not hasattr(self, '_post_'):
            raise tornado.web.HTTPError(405)
        self._post_(*args, **kwargs)

    @tornado.web.asynchronous
    @tornado.gen.engine
    def get(self, *args, **kwargs):
        '''get封装'''
        if not hasattr(self, '_get_'):
            raise tornado.web.HTTPError(405)
        self._get_(*args, **kwargs)

    @property
    def navbar_list(self):
        '''头部导航列表'''
        nav_list = [
            {"name": "index", "title": self._("topic list"), "href": "/web/resource"},
            {"name": "infoview", "title": self._("Statics"), "href": "/web/infoview/group"},
        ]
        return nav_list


class WebAsyncAuthHandler(WebHandler):
    """
        该类全部需要用户权限，并且集成session功能，
        所以子类不能使用@tornado.web.authenticated,
        需要非用户权限请使用RequestHandler
    """
    """
        覆盖get,post,put,delete方法，全部使用异步阻塞，
        子类使用_{method}_方法，并且在执行完时加上self.finish()
        需要加上的原因是子类可能也需要执行异步阻塞，
        如果父类调用sele.finish()则子类将不能使用connection
    """

    def permissionSet(self, permission):
        self.view_permission = permission

    def permissionCheck(self, permission):
        if self.view_permission is None:
            permission_pass = True
        else:
            view_permission = Permission(self.view_permission)
            user_permission = Permission(permission)
            permission_pass = user_permission.permissionCheck(view_permission)
        return permission_pass

    @tornado.web.asynchronous
    @tornado.gen.engine
    def get(self, *args, **kwargs):
        if not hasattr(self, '_get_'):
            raise tornado.web.HTTPError(405)

        #超时检测
        session_expired = yield tornado.gen.Task(self.is_session_expired)
        if session_expired:
            return

        #获取session中的用户属性
        self.editorinfo = yield tornado.gen.Task(self.session.get, 'editorinfo')

        #鉴权
        has_priv = yield tornado.gen.Task(self.check_priv, self.editorinfo)

        if not has_priv:
            if self.request.method in ('GET', 'HEAD'):
                url = setting.LOGIN_URL
                self.redirect(url)
                return

        #权限不够，重定向
        if not self.permissionCheck(self.editorinfo['permissions']):
            info = {}
            info["msg"] = self._('Not permitted to visit')
            info["url"] = "/"
            self.send_error(403, info=info)
            return

        self._get_(*args, **kwargs)

    @tornado.web.asynchronous
    @tornado.gen.engine
    def post(self, *args, **kwargs):
        if not hasattr(self, '_post_'):
            raise tornado.web.HTTPError(405)

        session_expired = yield tornado.gen.Task(self.is_session_expired)
        if session_expired:
            return

        self.editorinfo = yield tornado.gen.Task(self.session.get, 'editorinfo')
        has_priv = yield tornado.gen.Task(self.check_priv, self.editorinfo)

        if not has_priv:
            raise tornado.web.HTTPError(403)

        if not self.permissionCheck(self.editorinfo['permissions']):
            raise tornado.web.HTTPError(403)

        self._post_(*args, **kwargs)

    @tornado.web.asynchronous
    @tornado.gen.engine
    def put(self, *args, **kwargs):
        if not hasattr(self, '_put_'):
            raise tornado.web.HTTPError(405)

        session_expired = yield tornado.gen.Task(self.is_session_expired)
        if session_expired:  # for session_expired anti-moth 20160106
            return

        session = SessionManager(self)
        self.editorinfo = yield tornado.gen.Task(session.get, 'editorinfo')
        has_priv = yield tornado.gen.Task(self.check_priv, self.editorinfo)
        if not has_priv:
            raise tornado.web.HTTPError(403)

        if not self.permissionCheck(self.editorinfo['permissions']):
            raise tornado.web.HTTPError(403)

        self._put_(*args, **kwargs)

    @tornado.web.asynchronous
    @tornado.gen.engine
    def delete(self, *args, **kwargs):
        if not hasattr(self, '_delete_'):
            raise tornado.web.HTTPError(405)

        session_expired = yield tornado.gen.Task(self.is_session_expired)

        if session_expired:
            return

        session = SessionManager(self)
        self.editorinfo = yield tornado.gen.Task(session.get, 'editorinfo')
        has_priv = yield tornado.gen.Task(self.check_priv, self.editorinfo)

        if not has_priv:
            raise tornado.web.HTTPError(403)

        if not self.permissionCheck(self.editorinfo['permissions']):
            raise tornado.web.HTTPError(403)

        self._delete_(*args, **kwargs)

    def render(self, **kwargs):
        kwargs['is_login'] = True
        kwargs['editorinfo'] = self.editorinfo
        kwargs['canEdit'] = kwargs.get('canEdit', True)
        WebHandler.render(self, **kwargs)

    @tornado.gen.engine
    def make_editor_login_info(self, editorinfo, callback=None):

        editorinfo = simplejson.loads(editorinfo)
        if "editor_name" not in editorinfo or "editor_password" not in editorinfo:
            callback(True)
            return

        login_time = time.time()
        login_expires = datetime.datetime.utcnow() + datetime.timedelta(
            seconds=KEEP_LOGIN_SECONDS)
        self.set_cookie(
            EDITOR_LOGIN_TIME_COOKIE_KEY,
            "%d" % login_time,
            expires=login_expires)
        login_session = "%d|%s|%s " % (login_time, editorinfo["editor_name"],
                                       editorinfo['editor_password'])
        self.set_secure_cookie(
            EDITOR_SESSION_COOKIE_KEY, login_session, expires_days=None)
        callback(False)

    @tornado.gen.engine
    def is_session_expired(self, callback=None):
        '''后台审计'''
        login_time = self.get_cookie(EDITOR_LOGIN_TIME_COOKIE_KEY)

        user_id = self.get_secure_cookie("user_id")
        session = SessionManager(self)
        user_session = yield tornado.gen.Task(session.get, 'editorinfo')

        #通过审计则更新cookie生存时间
        if user_id and user_session and user_id == user_session.get(
                "user_id", ''):
            #update old time
            login_expires = datetime.datetime.utcnow() + datetime.timedelta(
                seconds=KEEP_LOGIN_SECONDS)

            if login_time is None:
                login_time = time.time()
                self.set_cookie(
                    EDITOR_LOGIN_TIME_COOKIE_KEY,
                    "%d" % login_time,
                    expires=login_expires)
            else:
                self.set_cookie(
                    EDITOR_LOGIN_TIME_COOKIE_KEY,
                    login_time,
                    expires=login_expires)

            callback(False)
            return
        else:
            #后台审计不通过则重定向
            redirect_url = setting.LOGIN_URL
            self.redirect(redirect_url)
            callback(True)
            return

        #cookie生存期耗尽将重定向
        if login_time is None:
            self.clear_cookie(EDITOR_SESSION_COOKIE_KEY)
            redirect_url = setting.LOGIN_URL
            callback(True)
            return

    @tornado.gen.engine
    def check_priv(self, editorinfo, callback=None):
        has_priv = True

        # if not editorinfo:
        if not editorinfo:
            try:
                # 尝试重建session回话
                login_time_lst = self.get_cookie(
                    EDITOR_LOGIN_TIME_COOKIE_KEY).split('|')
                login_time = int(login_time_lst[0])

                session_lst = self.get_secure_cookie(
                    EDITOR_SESSION_COOKIE_KEY).split('|')
                editor_name_cookie = session_lst[1].strip()
                editor_password_md5 = session_lst[2].strip()

                login_time_cookie = int(session_lst[0].strip())

                #校验time cookie
                if login_time != login_time_cookie:
                    has_priv = False

                elif not (editor_name_cookie and editor_password_md5):
                    has_priv = False
                else:
                    #获取用户服务
                    #from services.user import UserService, NoUserException, PasswdErrorException, NoAppException

                    #尝试获取用户信息
                    editor_service = UserService()
                    editor_id, editor = yield tornado.gen.Task(
                        editor_service.login, editor_name_cookie,
                        editor_password_md5)

                    #能检索到用户的话，则将用户部分数据拉取至内存
                    if editor_id is None:
                        has_priv = False
                    else:
                        editorinfo = {
                            'user_id': editor['user_id'],
                            'group_id': editor['group_id'],
                            'email': editor['email'],
                            'sex': editor['sex'],
                            'name': editor['name'],
                            'age': editor['age'],
                            'address': editor['address'],
                            'description': editor['description'],
                            'permissions': editor['permission'],
                            'passwd': editor['passwd'],
                        }

                        #保存会话
                        yield tornado.gen.Task(self.session.set, 'editorinfo',
                                               editorinfo)
                        logging.info(
                            '--INFO-- WebAsyncAuthHandler.check_priv recreate session'
                        )

            except Exception, e:
                #异常处理
                import traceback
                traceback.print_exc()
                logging.error(
                    '--ERROR-- WebAsyncAuthHandler.check_priv error, e=%s' %
                    str(e))
                has_priv = False

        callback(has_priv)
