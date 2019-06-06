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
    biangua home
'''
@route(r"/app/biangua", name="app.biangua")
class BianGuaHandler(WebHandler):

    @tornado.gen.engine
    def _get_(self):
        self.title      = 'biangua'
        self.template   = 'app/biangua.html'  
        self.nav_active = 'biangua'
        
        data = {}
        self.render(**data)
