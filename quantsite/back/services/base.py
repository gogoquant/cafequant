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

_all_ = ['BaseService','FieldException','NoFileException','get_file_data', 'RED_PACKET_ITEMS']

# 红包组件参数定义
RED_PACKET_ITEMS = {
    'red_game_duration': 'red_game_duration',   # 单次游戏时长, unit：s
    'red_total_cash': 'red_total_cash',         # 单次红包总额
    'red_packet_count': 'red_packet_count',     # 单次红包个数
    'red_allot': 'red_allot',                   # 分配方式 > 1 随机分配， 2 平均分配
    'red_ratio': 'red_ratio',                   # 中奖率
    'red_auth': 'red_auth',                     # 资格认证方式 > 1 imei验证， 2 微信验证
    'red_auth_age': 'red_auth_age',             # 资格认证时效 > 0 仅允许验证一次(永久有效)， 1 仅允许验证三次， 2 一天内一次， 3 一天内三次
    'red_shared_can_again': 'red_shared_can_again',  # 分享后是否可增加一次资格
    'auto_play': 'auto_play',                   # 游戏触发方式  0 点击按钮 ; 1 开始时延时触发
    'start_time': 'start_time',                 #延时触发时延时时间(秒数)
    'logo_md5': 'logo_md5',
    'brand_name': 'brand_name',
}

CONVERT_TYPE_MONGO_ITEMS = {
    "red_game_duration": int,
    "red_total_cash": float,
    "red_packet_count": int,
    "red_allot": int,
    "red_ratio": int,
    "red_auth": int,
    "red_auth_age": int,
    "red_shared_can_again": int,
    "red_begin_delay": int,
    "red_custom": int,
    "auto_play": int,
    "start_time": float,
}


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
    
    def parse_page(self,page,page_size):
        query = {}
        page      = page or 1
        page_size = page_size or setting.PAGE_SIZE
        pos       = pos = (int(page)-1) * int(page_size)

        query['pos']   = pos
        query['count'] = page_size
        return query
    
    def _(self, text):
        ''' Localisation shortcut '''
        return text
        #return self.locale.translate(text).encode("utf-8")
    
    def get_fields(self,model_cls):
        fields = {}
        for f in model_cls.__dict__:
            if type(model_cls.__dict__[f]) == bool:
                fields[f] = model_cls.__dict__[f]
        return fields


class NoFileException(Exception):
    pass




