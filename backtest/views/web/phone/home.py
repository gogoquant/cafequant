#-*- coding: UTF-8 -*-

'''
    @手机用户的
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
    手机主页
'''
@route(r"/phone/index", name="web.phone.home")
class IndexHandler(WebHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = '用户主页'
        self.template   = 'phone/demo.html'  
        self.nav_active = 'home'
        
        data = {}
        self.render(**data)


