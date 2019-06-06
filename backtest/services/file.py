# -*- coding: UTF-8 -*-

'''
    @brief service for file
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

from models.files.file_model import File

from services.base import BaseService 

from util.time_common import DAY_PATTERN, timestamp_to_string, FULL_PATTERN, week_seconds, day_seconds, month_seconds
from util.oauth2 import *
from util.qiniu_util import *

__all__ = ['EditorService', 'AdminService', 'SdkAdminService']


'''文件管理'''
class FileService(BaseService):
    file_m = File()

    #统计文件总数
    @tornado.gen.engine
    def count(self, query=None, callback=None):
        if query is None:
            query = {
        
            }
        c = yield tornado.gen.Task(self.file_m.count, query)
        callback(c)

    #通过文件id获取文件句柄
    @tornado.gen.engine
    def get(self, file_id=None, callback=None):
        '''使用id精确检索文章某个'''
        query = {}
        query['file_id'] = file_id
        vfile = yield tornado.gen.Task(self.file_m.find_one, query)
        callback(vfile)


    #通过文件md5来获取文件句柄
    @tornado.gen.engine
    def get_by_m5(self, md5=None, callback=None):
        '''使用id精确检索文章某个'''
        query = {}
        query['md5'] = md5
        vfile = yield tornado.gen.Task(self.file_m.find_one, query)
        callback(vfile)

    #添加新文件,如果已经有相同的文件，则直接返回该文件
    @tornado.gen.engine
    def add(self, md5, meta, store, file_name = None, file_size = None, width = None, height = None, callback=None):

        vfile = {'md5': md5, "meta":meta, "store":store}
        
        if file_name:
            vfile['file_name'] = file_name

        if file_size:
            vfile['file_size'] = file_size

        file_id = yield tornado.gen.Task(self.file_m.insert, vfile, upsert=True, safe=True)

        callback(file_id)

    #删除文件
    @tornado.gen.engine
    def delete(self, md5, callback=None):

        "根据msg的id删除"
        query = {"md5": md5}
        yield tornado.gen.Task(self.file_m.delete, query)
        callback(md5)

