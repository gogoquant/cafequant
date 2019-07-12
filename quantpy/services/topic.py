# -*- coding: UTF-8 -*-

'''
    @brief service for topic
    @author: xiyan
    @date 2017-01-03
'''

import string
import hashlib
import tornado.gen
import time
import tornado.httpclient
import random
import uuid
import pdb

from models.user.user_model import User, ResetPasswd
from models.user.topic_model import Topic

from services.base import BaseService

from util.time_common import DAY_PATTERN, timestamp_to_string, FULL_PATTERN, week_seconds, day_seconds, month_seconds
from util.oauth2 import *

__all__ = ['EditorService', 'AdminService', 'SdkAdminService']


'''文章管理'''
class TopicService(BaseService):
    
    topic_m = Topic()

    @tornado.gen.engine
    def get(self, topic_id=None, callback=None):
        '''使用id精确检索文章某个'''
        query = {}
        query['topic_id'] = topic_id
        topic = yield tornado.gen.Task(self.topic_m.find_one, query)
        callback(topic)

    @tornado.gen.engine
    def count(self, query, callback):
        "查询用户总数"
        #@TODO
        query = {
        
        }
        c = yield tornado.gen.Task(self.topic_m.count, query)
        callback(c)

    @tornado.gen.engine
    def get_list(self, user_id=None, tag_id=None,title=None, file_format=None, pos=None, count=None, sort_key=None, callback=None):
        '''检索文章分页列表, 根据时间排序'''
        
        #pdb.set_trace()

        args = {
            "pos": string.atoi(pos),
            "count": string.atoi(count),
            'sorts': [
                [sort_key, setting.DESC],
            ]
        }

        #设置过滤条件
        conditions = {}
    
        #设置过滤
        if user_id:
            conditions['user_id'] = user_id
        if title:
            conditions['title'] = {"$regex":title}
        if file_format:
            conditions['file_format'] = file_format
        if tag_id:
            conditions['tag_id'] = tag_id

        args["conditions"] = conditions

        topics = yield tornado.gen.Task(self.topic_m.get_list, args)
        callback(topics)

    @tornado.gen.engine
    def add(self, title, content, brief, tag=None, user_id=None, file_format=None, callback=None):
        topic = {'title':title, "brief":brief, "content":content, "user_id":user_id}
        
        if tag:
            topic['tag'] = tag
        if file_format:
            topic['file_format'] = file_format

        
        topic['ctime'] = time.time()
        topic['mtime'] = time.time()

        new_topic_id = yield tornado.gen.Task(self.topic_m.insert, topic, upsert=True, safe=True)
        callback(new_topic_id)

    @tornado.gen.engine
    def update(self, topic_id, user_id=None, title = None, brief=None, content= None, tag=None, file_format=None, callback=None):
        #pdb.set_trace()
        
        query = {"topic_id":topic_id}
        args = {}

        if not title is None:
            args['title']= title

        if not content is None:
            args['content']= content

        if not tag is None:
            args['tag_id']= tag

        if not file_format is None:
            args['file_format']= file_format

        if not brief is None:
            args['brief']= brief

        args['mtime'] = time.time()

        update_set = {'$set': args}

        yield tornado.gen.Task(self.topic_m.update, query, update_set)
        
        callback(None)

    @tornado.gen.engine
    def delete(self, topic_id, callback=None):
        "根据topic的名字删除"
        print topic_id
        query = {"topic_id": topic_id}
        yield tornado.gen.Task(self.topic_m.delete, query)
        callback(topic_id)

