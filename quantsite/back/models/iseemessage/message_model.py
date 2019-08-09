#!/usr/bin/python
# -*- coding: utf-8 -*-

__author__ = 'anti-moth'
from iseecore.models import MessageAsyncBaseModel

_all_ = [
    'IseeMessage',
    'ExpiredDeliverMessage',
    'DeliverMessage',
    'ExpiredIseeMessage',
]


class DeliverMessage(MessageAsyncBaseModel):
    """用户级消息存储"""
    msg_id = True
    msg_type = True  # 整数。 消息发起者和消息获取者的一种约定方式，方便解析消息(msg_data)
    msg_date = True  # 秒数。 用于消息的排序。因为会插入一对多消息，使用id排序无法完成。
    msg_from = True  # 消息的发起者
    msg_data = {}  # 单一接收者消息存储。json数据，由消息创建者自行决定结构，用msg_type以判别不同的格式
    msg_deliver = {  # 分发条件(到用户,应该从用户信息的角度去设置，至少能够筛选到某一个特定的用户)
        'editor_id': '',
        'app_id': '',
    }
    # 一对多消息的扩散方式，当用户级查询发现有新的一对多消息(n条>1)时则扩散(n条)对应于此用户的记录。
    # 扩散应包含内容：msg_type、 msg_date、 msg_deliver(分发条件应重新组合，以应对用户级查询)、 isee_msg_id
    # 消息本体仍从一对多消息集合里获取，当扩散或获取时可使用redis缓存。
    isee_msg_id = True

    is_read = True  # 阅读标志。0表示未读， 1标识已读

    table = "deliver_message"
    key = "msg_id"


class ExpiredDeliverMessage(MessageAsyncBaseModel):
    """用户级过期消息存储"""
    msg_id = True
    msg_type = True
    msg_date = True
    msg_from = True
    msg_data = {}
    msg_deliver = {
        'editor_id': '',
        'app_id': '',
    }
    isee_msg_id = True
    is_read = True

    table = "expired_deliver_message"
    key = "msg_id"


class IseeMessage(MessageAsyncBaseModel):
    """一对多消息存储，如推送消息等"""
    msg_id = True
    msg_type = True
    msg_data = {}
    msg_from = True
    msg_date = True
    msg_deliver = {  # 多用户筛选条件，得到的是一个用户组。
        'vip_level': '',
        'app_id': '',
    }

    table = 'isee_message'
    key = 'msg_id'


class ExpiredIseeMessage(MessageAsyncBaseModel):
    """一对多过期消息存储"""
    msg_id = True
    msg_type = True
    msg_data = {}
    msg_from = True
    msg_date = True
    msg_deliver = False

    table = 'expired_isee_message'
    key = 'msg_id'
