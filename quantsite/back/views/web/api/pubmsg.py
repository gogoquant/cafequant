#-*- coding: UTF-8 -*-
'''
    @brief 消息管理
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
from services.pubmsg import PubMsgService

from stormed import Connection
from tornadomail.message import EmailMessage, EmailMultiAlternatives
import simplejson
import hashlib

import setting
from pycket.driver import Driver


'''
    发送消息
'''
@route(r"/api/pub/add", name="api.pub.add")
class PubAddHandler(WebAsyncAuthHandler):

    @tornado.gen.engine
    def _post_(self):
        self.pub_s = PubMsgService()
        #pdb.set_trace()

        body = self.request.body
        js = simplejson.loads(body, encoding='utf-8')


        title =     js.get("title", None)
        brief =     js.get("brief", None)
        content =   js.get("content", None)
        user_to =   js.get("user_to", None)
        msg_type =  js.get("msg_type", None)
        img_url =   js.get("img_url", None)
        redir_url = js.get("redir_url", None)
        
        #pdb.set_trace()

        if not hasattr(self, 'editorinfo'):
            user_id = None
        else:
            user_id = self.editorinfo["user_id"]
        
        if title is None:
            self.render(msg="参数不足")
            return
        
        pub_id = yield tornado.gen.Task(self.pub_s.add, title=title, brief=brief, content=content, \
                user_id=user_id, img_url=img_url, redir_url=redir_url, \
                user_to=user_to,msg_type=msg_type)

        self.render_success(msg=pub_id)
        return

'''
    获取消息集合
'''
@route(r"/api/pub/search", name="api.pub.search")
class PubSearchHandler(WebHandler):

    @tornado.gen.engine
    def _post_(self):
        self.pub_s = PubMsgService()

        #pdb.set_trace()
        
        pos = self.get_argument("pos", 0)
        count = self.get_argument("count", 0)

        #查询对应类型的消息集合
        msg_type = self.get_argument("msg_type", None)

        conditions = {}

        if not msg_type is None:
            conditions['msg_type'] = msg_type

        pubs = yield tornado.gen.Task(self.pub_s.get_list_desc, pos=pos, count=count, conditions=conditions)
        
        logging.info("pub count %d" % len(pubs) )
        
        self.render_success(msg=pubs)
        return

'''
    获取留言数量
'''
@route(r"/api/pub/count", name="api.pub.count")
class PubCountHandler(WebHandler):
    pub_s = PubMsgService()

    @tornado.gen.engine
    def _get_(self):
        self.pub_s = PubMsgService()

        tot = yield tornado.gen.Task(self.pub_s.count)
        self.render_success(msg=tot)
        return

    @tornado.gen.engine
    def _post_(self):
        self._get_()
        return

'''
    删除动态
'''
@route(r"/api/pub/delete", name="api.pub.delete")
class PubDeleteHandler(WebAsyncAuthHandler):
    @tornado.gen.engine
    def _post_(self):
        self.pub_s = PubMsgService()
        msg_id = self.get_argument("msg_id", 0)

        yield tornado.gen.Task(self.pub_s.delete, msg_id)
        self.render_success(msg=msg_id)
        return

'''
    获取某个
'''
@route(r"/api/pub/get", name="api.pub.get")
class PubGetHandler(WebHandler):
    
    @tornado.gen.engine
    def _post_(self):
        self.pub_s = PubMsgService()
        msg_id = self.get_argument("msg_id", 0)

        msg = yield tornado.gen.Task(self.pub_s.get, msg_id)
        self.render_success(msg=msg)
        return

'''
    更新动态
'''
@route(r"/api/pub/update", name="api.pub.update")
class PubUpdateHandler(WebAsyncAuthHandler):
    
    @tornado.gen.engine
    def _post_(self):
        self.pub_s = PubMsgService()

        body = self.request.body
        js = simplejson.loads(body, encoding='utf-8')
        
        msg_id =    js.get("msg_id", None)
        title =     js.get("title", None)
        brief =     js.get("brief", None)
        content =   js.get("content", None)
        user_to =   js.get("user_to", None)
        msg_type =  js.get("msg_type", None)
        img_url =   js.get("img_url", None)
        redir_url = js.get("redir_url", None)
        
        #pdb.set_trace()

        if not hasattr(self, 'editorinfo'):
            user_id = None
        else:
            user_id = self.editorinfo["user_id"]

        yield tornado.gen.Task(self.pub_s.update, msg_id=msg_id, user_id=user_id, user_to=user_to, title =title, \
            brief=brief, content=content, img_url=img_url, redir_url=redir_url, msg_type=msg_type)
        
        self.render_success(msg=msg_id)
        return
