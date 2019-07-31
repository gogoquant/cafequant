# -*- coding: UTF-8 -*-

'''
    @brief service for msg
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
from models.user.message_model import ChatMessage

from services.base import BaseService 

from util.time_common import DAY_PATTERN, timestamp_to_string, FULL_PATTERN, week_seconds, day_seconds, month_seconds
from util.oauth2 import *

__all__ = ['EditorService', 'AdminService', 'SdkAdminService']


'''留言板管理'''
class ChatMsgService(BaseService):
    
    msg_m = ChatMessage()
    limit_m = 8

    @tornado.gen.engine
    def count(self, query=None, callback=None):
        "查询msg总数"
        query = {
        
        }
        c = yield tornado.gen.Task(self.msg_m.count, query)
        callback(c)

    @tornado.gen.engine
    def get_list_asc(self, pos, count, callback=None):
        args = {
            "pos": pos,
            "count": count,
            'sorts': [
                ['ctime', setting.ASC],
            ]
        }
        #设置过滤条件
        conditions = {}
    
        args["conditions"] = conditions

        msgs = yield tornado.gen.Task(self.msg_m.get_list, args)
        callback(msgs)


    @tornado.gen.engine
    def get_list_desc(self, pos, count, callback=None):
        args = {
            "pos": pos,
            "count": count,
            'sorts': [
                ['ctime', setting.DESC],
            ]
        }

        #设置过滤条件
        conditions = {}
    
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
    def add(self, title, user_id, callback=None):
        msg = {'title':title, "user_id":user_id, "ctime":time.time() }

        #填充新消息
        new_msg_id = yield tornado.gen.Task(self.msg_m.insert, msg, upsert=True, safe=True)

        #获取消息总数
        msg_tot = yield tornado.gen.Task(self.msg_m.count)

        #logging.info("%d/%d" % (msg_tot, self.limit_m))       

        #删除多余的消息
        if msg_tot >  self.limit_m:
            msg_num = msg_tot - self.limit_m
            msgs = yield tornado.gen.Task(self.get_list_asc, 0, msg_num)
            
            #logging.info("delete msg tot %d" % len(msgs))       

            for msg in msgs:
                query = {"chatmessage_id": msg["chatmessage_id"]}
                yield tornado.gen.Task(self.msg_m.delete, query)

        callback(new_msg_id)

    @tornado.gen.engine
    def delete(self, msg_id, callback=None):

        "根据msg的id删除"
        query = {"chatmsg_id": msg_id}
        yield tornado.gen.Task(self.msg_m.delete, query)
        callback(msg_id)
