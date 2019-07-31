#-*- coding: UTF-8 -*-

'''
    @普通用户的
    @author: lancelot
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
    主页
'''
@route(r"/", name="web.web.home")
class IndexHandler(WebHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '用户主页'
        self.template   = 'web/index.html'  
        self.nav_active = 'home'
        
        data = {}
        self.render(**data)


'''
    文章列表
'''
@route(r"/web/topicList", name="web.web.topiclist")
class TopicListHandler(WebHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '文章列表主页'
        self.template   = 'web/topicList.html'  
        self.nav_active = 'topiclist'
        
        data = {}
        self.render(**data)


'''
    单文章详情
'''
@route(r"/web/topic", name="web.web.topic")
class TopicHandler(WebHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '文章列表主页'
        self.template   = 'web/topic.html'  
        self.nav_active = 'topic'
        
        data = {}
        self.render(**data)
