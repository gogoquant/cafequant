# -*- coding: UTF-8 -*-
'''
    @brief service for tag
    @author: xiyan
    @date 2017-01-03
'''

import hashlib
import tornado.gen
import time
import tornado.httpclient
import random
import uuid
import pdb

from models.user.user_model import User, ResetPasswd
from models.user.tag_model import Tag

from services.base import BaseService

from util.time_common import DAY_PATTERN, timestamp_to_string, FULL_PATTERN, week_seconds, day_seconds, month_seconds
#from util.oauth2 import *

__all__ = ['TagService']
'''文章管理'''


class TagService(BaseService):

    tag_m = Tag()

    @tornado.gen.engine
    def get(self, tag_id=None, callback=None):
        '''使用id精确检索某个tag'''
        query = {}
        query['tag_id'] = tag_id
        tag = yield tornado.gen.Task(self.tag_m.find_one, query)
        callback(tag)

    @tornado.gen.engine
    def count(self, query, callback):
        "查询tag总数"
        #@TODO
        query = {}
        c = yield tornado.gen.Task(self.tag_m.count, query)
        callback(c)

    @tornado.gen.engine
    def get_list(self,
                 user_id=None,
                 title=None,
                 pos=None,
                 count=None,
                 callback=None):
        '''检索文章分页列表, 根据时间排序'''

        args = {"pos": pos, "count": count, 'sorts': [['mtime', setting.DESC],]}

        #设置过滤条件
        conditions = {}

        #设置过滤
        if user_id:
            conditions['user_id'] = user_id
        if title:
            conditions['title'] = {"$regex": title}

        args["conditions"] = conditions

        tags = yield tornado.gen.Task(self.tag_m.get_list, args)
        callback(tags)

    @tornado.gen.engine
    def add(self, title, user_id, callback=None):
        tag = {'title': title, "user_id": user_id}

        new_tag_id = yield tornado.gen.Task(
            self.tag_m.insert, tag, upsert=True, safe=True)
        callback(new_tag_id)

    @tornado.gen.engine
    def update(self, tag_id, title=None, callback=None):
        query = {"tag_id": tag_id}

        args = {}

        if not title is None:
            args['title'] = title

        update_set = {'$set': args}

        yield tornado.gen.Task(self.tag_m.update, query, update_set)

        callback(None)

    @tornado.gen.engine
    def delete(self, tag_id, callback=None):
        "根据tag的名字删除"
        print tag_id
        query = {"tag_id": tag_id}
        yield tornado.gen.Task(self.tag_m.delete, query)
        callback(tag_id)
