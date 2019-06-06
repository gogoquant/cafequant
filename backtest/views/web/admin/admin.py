#-*- coding: UTF-8 -*-

'''
    @管理员的view
    @author: xiyan
    @data 2016-12-26

'''

import logging
import tornado.gen
import tornado.web
import pdb

from tornado.options import define, options
from iseecore.routes  import route
from views.web.base import *

from services.user import UserService, NoUserException, PasswdErrorException, UserExistsException, UserSameNameException

from stormed import Connection
from tornadomail.message import EmailMessage, EmailMultiAlternatives
import simplejson
import hashlib

import setting
from pycket.driver import Driver


'''
    管理员用户登录
'''
@route(r"/adminweb/adminLogin", name="web.adminweb.adminLogin")
class AdminLoginHandler(WebHandler):
    user_s = UserService()

    @tornado.gen.engine
    def _get_(self):
        self.title      = 'user login'
        self.template   = 'admin/admin_login.html'  
        self.nav_active = 'admin'
        
        data = {}
        self.render(**data)

    @tornado.gen.engine
    def _post_(self):

        #pdb.set_trace()

        #获取用户基本信息
        email = self.get_argument("email", None)
        passwd = self.get_argument("passwd", "")
        token = self.get_argument("access_token", None)
        token_from = self.get_argument("source", None)
        openid = self.get_argument("openid", None)
        remember        = int(self.get_argument('remember',0))
        
        #必要信息缺失则直接错误
        if (not email or not passwd) and (not token or not token_from):
            #self.render_error(status=400, code=400, msg='no email or token')
            self.redirect('/adminweb/adminError', permanent=True)
            return

        #对比，更新cookie
        if remember:
            Driver.EXPIRE_SECONDS = 30 * 24 * 60 * 60  
            if remember != editor_remember_cookie:
                self.set_cookie(EDITOR_NAME_REMEMBER_COOKIE_KEY, str(remember), expires_days=30)
            if editor_name != editor_name_cookie:
                self.set_cookie(EDITOR_NAME_COOKIE_KEY, editor_name, expires_days=30)
        else:
            Driver.EXPIRE_SECONDS = 1 * 24 * 60 * 60   
            self.clear_cookie(EDITOR_NAME_REMEMBER_COOKIE_KEY)
            self.clear_cookie(EDITOR_NAME_COOKIE_KEY)

        user_id, editorinfo = yield tornado.gen.Task(self.user_s.login, email, passwd, token, token_from, openid=openid, app_id=None)
        
        #pdb.set_trace()

        #根据返回值作处理
        if not editorinfo:
            self.redirect('/adminweb/adminError', permanent=True)
            #self.render_error(status=400, code=400, msg='no social data')
            return

        if editorinfo == NoUserException:
            self.redirect('/adminweb/adminError', permanent=True)
            #self.render_error(status=404, code=404, msg='no user')
            return

        if editorinfo == PasswdErrorException:
            self.redirect('/adminweb/adminError', permanent=True)
            #self.render_error(status=403, code=403, msg='passwd error')
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
    admin home
'''
@route(r"/adminweb/index", name="web.adminweb.index")
class AdminHomeHandler(WebAsyncAuthHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '管理员主页'
        self.template   = 'admin/admin_home.html'  
        self.nav_active = 'admin'
        
        data = {}
        self.render(**data)

'''
    admin user list
'''
@route(r"/adminweb/user_list", name="web.adminweb.user_list")
class AdminHomeHandler(WebAsyncAuthHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '管理员用户列表主页'
        self.template   = 'admin/admin_user.html'  
        self.nav_active = 'admin'
        
        data = {}
        self.render(**data)


'''
    admin topic editor
'''
@route(r"/adminweb/topic_editor", name="web.adminweb.topic_editor")
class AdminHomeHandler(WebAsyncAuthHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '管理员topic编辑页面'
        self.template   = 'admin/topic_editor.html'  
        self.nav_active = 'admin'
        
        data = {}
        topic_id = self.get_argument("topic_id", "")
        data["topic_id"] = topic_id
        self.render(**data)

'''
    admin topic list
'''
@route(r"/adminweb/topic_list", name="web.adminweb.topic_list")
class AdminHomeHandler(WebAsyncAuthHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '管理员用户文章管理主页'
        self.template   = 'admin/admin_topic.html'  
        self.nav_active = 'topic'
        
        data = {}
        self.render(**data)


'''
    admin msg list
'''
@route(r"/adminweb/msg_list", name="web.adminweb.msg_list")
class AdminHomeHandler(WebAsyncAuthHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '消息管理主页'
        self.template   = 'admin/msg_list.html'  
        self.nav_active = 'topic'
        
        data = {}
        self.render(**data)

'''
    admin msg
'''
@route(r"/adminweb/msg", name="web.adminweb.msg")
class AdminHomeHandler(WebAsyncAuthHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '消息编辑主页'
        self.template   = 'admin/msg.html'  
        self.nav_active = 'msg'
        
        data = {}
        self.render(**data)



'''
    admin error
'''
@route(r"/adminweb/adminError", name="web.adminweb.adminError")
class AdminErorHandler(WebHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '管理员error page'
        self.template   = 'admin/admin_error.html'  
        self.nav_active = 'admin'
        
        data = {}
        self.render(**data)

'''
    admin upload
'''
@route(r"/adminweb/upload", name="web.adminweb.upload")
class AdminErorHandler(WebAsyncAuthHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '管理员upload'
        self.template   = 'admin/file_upload.html'  
        self.nav_active = 'admin'
        
        data = {}
        self.render(**data)


