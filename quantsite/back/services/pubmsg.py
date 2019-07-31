# -*- coding: UTF-8 -*-

'''
    @brief service for pubmsg
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
from models.user.message_model import PubMessage

from services.base import BaseService 

from util.time_common import DAY_PATTERN, timestamp_to_string, FULL_PATTERN, week_seconds, day_seconds, month_seconds
from util.oauth2 import *

__all__ = ['EditorService', 'AdminService', 'SdkAdminService']


'''留言板管理'''
class PubMsgService(BaseService):
    _limit_map = {} 

    def __init__(self):
        self.msg_m = PubMessage()

    #设置某类消息的阀值
    @classmethod
    def  setLimit(cls, msg_type, limit):
        cls._limit_map[msg_type] = limit


    #获取某类消息的阀值
    @classmethod
    def  getLimit(cls, msg_type):
        if msg_type in cls._limit_map.keys():
            return cls._limit_map[msg_type]
        else:
            return None

    @tornado.gen.engine
    def count(self, query=None, callback=None):
        "查询msg总数"
        query = {
        
        }
        c = yield tornado.gen.Task(self.msg_m.count, query)
        callback(c)


    @tornado.gen.engine
    def get(self, msg_id=None, callback=None):
        '''使用id精确检索某个msg'''
        query = {}
        query['pubmsg_id'] = msg_id
        msg = yield tornado.gen.Task(self.msg_m.find_one, query)
        callback(msg)


    @tornado.gen.engine
    def get_list_asc(self, pos, count, conditions, callback=None):
        args = {
            "pos": pos,
            "count": count,
            'sorts': [
                ['ctime', setting.ASC],
            ]
        }
        #设置过滤条件
        #conditions = {}
    
        args["conditions"] = conditions

        msgs = yield tornado.gen.Task(self.msg_m.get_list, args)
        callback(msgs)


    @tornado.gen.engine
    def get_list_desc(self, pos, count, conditions, callback=None):
        args = {
            "pos": pos,
            "count": count,
            'sorts': [
                ['ctime', setting.DESC],
            ]
        }

        #设置过滤条件
        #conditions = {}
    
        args["conditions"] = conditions

        msgs = yield tornado.gen.Task(self.msg_m.get_list, args)
        
        for msg in msgs:
            if msg["user_id"] is None:
                msg["name"] = "匿名"
            else:
                #非匿名
                pass

        #import pdb
        #pdb.set_trace()

        callback(msgs)

    @tornado.gen.engine
    def add(self, title, brief=None, content=None, user_id=None, img_url=None, \
            redir_url=None, user_to=None,msg_type=None, callback=None):
        
        read = 0

        msg = {'title':title, "brief":brief, "content":content, "user_id":user_id,\
                "redir_url":redir_url, "img_url":img_url, "user_to":user_to, "msg_type":msg_type, "read":read,
                "ctime":time.time() }

        #pdb.set_trace()
        #填充新消息
        new_msg_id = yield tornado.gen.Task(self.msg_m.insert, msg, upsert=True, safe=True)


        #如果该类消息有阀值，则删除多余的消息
        limit = PubMsgService.getLimit(msg_type)

        if not limit is None:
            #获取该类消息总数
            query = {}
            query['msg_type'] = msg_type
            msg_tot = yield tornado.gen.Task(self.msg_m.count, query)

            logging.info("%d/%d" % (msg_tot, limit))       

            #删除多余的消息
            if msg_tot > limit:
                query = {}
                msg_num = msg_tot - limit
                msgs = yield tornado.gen.Task(self.get_list_asc, 0, msg_num, {"msg_type":msg_type})
                logging.info("delete msg tot %d" % len(msgs))       
                for msg in msgs:
                    query = {"pubmsg_id": msg["pubmsg_id"]}
                    yield tornado.gen.Task(self.msg_m.delete, query)
        
        callback(new_msg_id)

    @tornado.gen.engine
    def delete(self, msg_id, callback=None):

        "根据msg的id删除"
        query = {"pubmsg_id": msg_id}
        yield tornado.gen.Task(self.msg_m.delete, query)
        callback(msg_id)


    @tornado.gen.engine
    def update(self, msg_id=None, user_id=None, user_to=None, title =None, \
            brief=None, content=None, img_url=None, redir_url=None, msg_type=None, callback=None):
        #pdb.set_trace()
        
        query = {"pubmsg_id":msg_id}
        args = {}

        if not title is None:
            args['title']= title

        if not content is None:
            args['content']= content

        if not brief is None:
            args['brief']= brief

        if not user_id is None:
            args['user_id']= user_id

        if not user_to is None:
            args['user_to']= user_to

        if not msg_type is None:
            args['msg_type']= msg_type

        if not img_url is None:
            args['img_url']= img_url

        if not redir_url is None:
            args['redir_url']= redir_url

        args['mtime'] = time.time()

        update_set = {'$set': args}

        yield tornado.gen.Task(self.msg_m.update, query, update_set)
    

#公告消息限制
PubMsgService.setLimit("publish", 1)

#动态消息限制
PubMsgService.setLimit("active", 500)
