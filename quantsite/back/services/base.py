#-*- coding: UTF-8 -*-
'''
    @brief handler for service 
    @author: xiyanxiyan10
    @data 2016-12-13
    @note 服务器相关类
'''

import tornado.gen
import setting
from tornado.options import options
import logging

_all_ = [
    'BaseService', 'FieldException', 'NoFileException', 'get_file_data',
    'RED_PACKET_ITEMS'
]


class FieldException(object):

    def __init__(self, arg):
        self.arg = arg


class BaseService(dict):

    def __init__(self):
        self.ch = options.group_dict('ch')
        self.fastdfs_client = options.group_dict('fastdfs_client')
        self.redis_client = options.group_dict('redis_client')

        self.riak_client = options.group_dict('riak_client')

    def init(self):
        pass

    def parse_page(self, page, page_size):
        query = {}
        page = page or 1
        page_size = page_size or setting.PAGE_SIZE
        pos = pos = (int(page) - 1) * int(page_size)

        query['pos'] = pos
        query['count'] = page_size
        return query

    def _(self, text):
        ''' Localisation shortcut '''
        return text
        #return self.locale.translate(text).encode("utf-8")

    def get_fields(self, model_cls):
        fields = {}
        for f in model_cls.__dict__:
            if type(model_cls.__dict__[f]) == bool:
                fields[f] = model_cls.__dict__[f]
        return fields


class NoFileException(Exception):
    pass
