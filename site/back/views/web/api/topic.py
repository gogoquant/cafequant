#-*- coding: UTF-8 -*-
'''
    @brief 文章管理
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

from services.topic import TopicService
from services.pubmsg import PubMsgService
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


@route(r"/api/topic/add", name="api.topic.add")
class TopicAddHandler(WebAsyncAuthHandler):
    topic_s = TopicService()
    pub_s = PubMsgService()

    @tornado.gen.engine
    def _post_(self):

        #import pdb
        #pdb.set_trace()

        body = self.request.body
        js = simplejson.loads(body, encoding='utf-8')

        #获取用户数据
        title = js.get("title", None)
        brief = js.get("brief", None)
        content = js.get("content", None)
        tag = js.get("tag", None)
        #file_format = js.get("file_format", None)

        #从session中获取用户id
        if not hasattr(self, 'editorinfo'):
            user_id = None
        else:
            user_id = self.editorinfo["user_id"]

        if title is None or content is None:
            self.render(msg="参数不足")
            return

        #only 20 char
        if len(brief) >= 20:
            brief = brief[0:20]

        topic_id = yield tornado.gen.Task(
            self.topic_s.add,
            title=title,
            tag=tag,
            content=content,
            user_id=user_id,
            brief=brief)

        self.render_success(msg=topic_id)
        return


'''
    获取用户文章列表
'''


@route(r"/api/topic/list", name="api.topic.list")
class TopicListHandler(WebHandler):
    topic_s = TopicService()
    user_s = UserService()
    tag_s = TagService()

    @tornado.gen.engine
    def _post_(self):

        #pdb.set_trace()

        #获取用户数据
        title = self.get_argument("title", None)
        user_id = self.get_argument("user_id", None)
        tag_id = self.get_argument("tag_id", None)
        file_format = self.get_argument("file_format", None)
        sort_key = self.get_argument("sort_key", None)

        if sort_key is None:
            sort_key = "mtime"

        pos = self.get_argument("pos", None)
        count = self.get_argument("count", None)

        logging.info("pos %s count %s tag_id %s user_id %s" %
                     (pos, count, tag_id, user_id))

        topics = yield tornado.gen.Task(
            self.topic_s.get_list,
            tag_id=tag_id,
            title=title,
            file_format=file_format,
            pos=pos,
            count=count,
            sort_key=sort_key)

        for topic in topics:
            user_info = yield tornado.gen.Task(
                self.user_s.get,
                user_id=topic["user_id"],
                name=None,
                email=None)
            topic["user_info"] = user_info

            tag_idx = topic.get("tag_id", None)

            if not tag_idx is None:
                tag_info = yield tornado.gen.Task(self.tag_s.get,
                                                  topic["tag_id"])
                topic["tag_info"] = tag_info
            else:
                tag_info = {"title": "default"}
                topic["tag_info"] = tag_info

        logging.info("topic count %d" % len(topics))

        self.render_success(msg=topics)
        return


'''
    更新用户文章
'''


@route(r"/api/topic/update", name="api.topic.update")
class TopicUpdateHandler(WebAsyncAuthHandler):
    topic_s = TopicService()

    @tornado.gen.engine
    def _post_(self):

        #pdb.set_trace()

        body = self.request.body
        js = simplejson.loads(body, encoding='utf-8')

        #获取用户数据
        topic_id = js.get("topic_id", None)
        title = js.get("title", None)
        content = js.get("content", None)
        user_id = js.get("user_id", None)
        tag_id = js.get("tag_id", None)
        brief = js.get("brief", None)

        #only 20 char
        if len(brief) >= 20:
            brief = brief[0:20]

        file_format = self.get_argument("file_format", None)

        yield tornado.gen.Task(self.topic_s.update, topic_id, user_id, title,
                               brief, content, tag_id, file_format)

        self.render_success(msg="success")
        return


'''
    获取用户文章详细数据
'''


@route(r"/api/topic/get", name="api.topic.get")
class TopicGetHandler(WebHandler):
    topic_s = TopicService()
    user_s = UserService()
    tag_s = TagService()

    @tornado.gen.engine
    def _post_(self):

        #获取用户数据
        topic_id = self.get_argument("topic_id", None)

        logging.info(topic_id)

        topic = yield tornado.gen.Task(self.topic_s.get, topic_id)

        user_info = yield tornado.gen.Task(
            self.user_s.get, user_id=topic["user_id"], name=None, email=None)
        topic["user_info"] = user_info
        #logging.info(topic)

        tag_id = topic.get("tag_id", None)

        if not tag_id is None:
            tag_info = yield tornado.gen.Task(self.tag_s.get, topic["tag_id"])
            topic["tag_info"] = tag_info
        else:
            tag_info = {"title": "default"}
            topic["tag_info"] = tag_info

        self.render_success(msg=topic)


'''
    获取用户文章总数
'''


@route(r"/api/topic/tot", name="api.topic.tot")
class TopicTotHandler(WebHandler):
    topic_s = TopicService()

    @tornado.gen.engine
    def _post_(self):
        query = {}
        tot = yield tornado.gen.Task(self.topic_s.count, query)
        self.render_success(msg=tot)


'''
    删除文章
'''


@route(r"/api/topic/delete", name="api.topic.delete")
class TopicDeleteHandler(WebAsyncAuthHandler):
    topic_s = TopicService()

    @tornado.gen.engine
    def _post_(self):

        topic_id = self.get_argument("topic_id", "")

        yield tornado.gen.Task(self.topic_s.delete, topic_id)

        self.render_success(msg="success")
