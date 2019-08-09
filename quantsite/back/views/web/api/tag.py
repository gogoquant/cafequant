#-*- coding: UTF-8 -*-
'''
    @brief tag管理
    @author: xiyan
    @data 2016-12-26

'''

import logging
import tornado.gen
import tornado.web
import pdb

from tornado.options import define, options
from iseecore.routes import route
from views.web.base import WebHandler, WebAsyncAuthHandler

from services.user import UserService, NoUserException, PasswdErrorException, UserExistsException, UserSameNameException
from services.tag import TagService

from stormed import Connection
from tornadomail.message import EmailMessage, EmailMultiAlternatives
import simplejson
import hashlib

import setting
from pycket.driver import Driver
'''
    添加用户文章
'''


@route(r"/api/tag/add", name="api.tag.add")
class TagAddHandler(WebAsyncAuthHandler):
    tag_s = TagService()

    @tornado.gen.engine
    def _post_(self):
        #pdb.set_trace()

        #获取用户数据
        title = self.get_argument("title", None)

        if not hasattr(self, 'editorinfo'):
            user_id = None
        else:
            user_id = self.editorinfo["user_id"]

        if title is None:
            self.render(msg="参数不足")
            return

        tag_id = yield tornado.gen.Task(
            self.tag_s.add, title=title, user_id=user_id)

        self.render_success(msg=tag_id)
        return


'''
    获取用户文章列表
'''


@route(r"/api/tag/search", name="api.tag.search")
class TagListHandler(WebHandler):
    tag_s = TagService()

    @tornado.gen.engine
    def _post_(self):

        #pdb.set_trace()

        pos = self.get_argument("pos", None)
        count = self.get_argument("count", None)
        title = self.get_argument("title", None)

        logging.info("pos %s count %s" % (pos, count))

        tags = yield tornado.gen.Task(
            self.tag_s.get_list, title=title, pos=pos, count=count)

        logging.info("tag count %d" % len(tags))

        self.render_success(msg=tags)
        return


'''
    更新用户文章
'''


@route(r"/api/tag/update", name="api.tag.update")
class TagUpdateHandler(WebAsyncAuthHandler):
    tag_s = TagService()

    @tornado.gen.engine
    def _post_(self):
        #获取用户数据
        tag_id = self.get_argument("tag_id", None)
        title = self.get_argument("title", None)

        yield tornado.gen.Task(self.tag_s.update, tag_id=tag_id, title=title)

        self.render_success(msg="success")
        return


'''
    获取用户文章详细数据
'''


@route(r"/api/tag/get", name="api.tag.get")
class TagGetHandler(WebHandler):
    tag_s = TagService()

    @tornado.gen.engine
    def _post_(self):
        #获取用户数据
        tag_id = self.get_argument("tag_id", None)

        logging.info(tag_id)

        tag = yield tornado.gen.Task(self.tag_s.get, tag_id)

        logging.info(tag)

        self.render_success(msg=tag)


'''
    获取用户tag总数
'''


@route(r"/api/tag/tot", name="api.tag.tot")
class TagTotHandler(WebHandler):
    tag_s = TagService()

    @tornado.gen.engine
    def _post_(self):
        query = {}
        tot = yield tornado.gen.Task(self.tag_s.count, query)
        self.render_success(msg=tot)


'''
    删除tag
'''


@route(r"/api/tag/delete", name="api.tag.delete")
class TagDeleteHandler(WebAsyncAuthHandler):
    tag_s = TagService()

    @tornado.gen.engine
    def _post_(self):
        #pdb.set_trace()

        tag_id = self.get_argument("tag_id", "")

        yield tornado.gen.Task(self.tag_s.delete, tag_id)

        self.render_success(msg="success")
