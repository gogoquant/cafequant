#-*- coding: UTF-8 -*-
'''
    @brief chat管理
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
from services.chat import ChatMsgService

from stormed import Connection
from tornadomail.message import EmailMessage, EmailMultiAlternatives
import simplejson
import hashlib

import setting
from pycket.driver import Driver


'''
    添加留言
'''
@route(r"/api/chat/add", name="api.chat.add")
class TagAddHandler(WebAsyncAuthHandler):
    chat_s = ChatMsgService()

    @tornado.gen.engine
    def _post_(self):
        #pdb.set_trace()

        #获取用户数据
        title = self.get_argument("title", None)

        if len(title) > 9:
            title = title[0:9]
        
        if not hasattr(self, 'editorinfo'):
            user_id = None
        else:
            user_id = self.editorinfo["user_id"]
        
        if title is None:
            self.render(msg="参数不足")
            return
        
        chat_id = yield tornado.gen.Task(self.chat_s.add, title=title, user_id=user_id)
        
        self.render_success(msg=chat_id)
        return

'''
    获取留言列表
'''
@route(r"/api/chat/search", name="api.chat.search")
class TagListHandler(WebHandler):
    chat_s = ChatMsgService()

    @tornado.gen.engine
    def _post_(self):

        #pdb.set_trace()
        pos = 0
        count = self.get_argument("count", 0)

        chats = yield tornado.gen.Task(self.chat_s.get_list_desc, pos=pos, count=count)
        
        logging.info("chat count %d" % len(chats) )
        
        self.render_success(msg=chats)
        return

'''
    获取留言数量
'''
@route(r"/api/chat/count", name="api.chat.count")
class TagListHandler(WebHandler):
    chat_s = ChatMsgService()

    @tornado.gen.engine
    def _get_(self):

        tot = yield tornado.gen.Task(self.chat_s.count)
        self.render_success(msg=tot)
        return
