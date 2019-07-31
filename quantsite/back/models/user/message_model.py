#-*- coding: UTF-8 -*-

'''
    message used for chat

'''

from iseecore.models import AsyncBaseModel

class ChatMessage(AsyncBaseModel):
    user_id    = False
    title      = False

    ctime      = False

    table           = "chatmessage"
    key             = "chatmessage_id"

#动态日志
class PubMessage(AsyncBaseModel):
    
    #用户id
    user_id    = False
    
    #消息标题
    title      = False

    #消息简介
    brief      = False

    #消息content
    content    = False

    #图片地址
    img_url    = False


    #链接跳转地址
    redir_url = False


    #阅读标志
    read       = False

    #发送对象
    user_to    = False

    #消息类型
    msg_type   = False


    #消息创建时间
    ctime      = False

    table           = "pubmsg"
    key             = "pubmsg_id"

#系统管理日志
class LogMessage(AsyncBaseModel):
    
    #用户id
    user_id    = False
    
    #日志级别
    level      = False

    #日志类型
    msgtype  = False


    #消息创建时间
    ctime      = False

    table           = "logmsg"
    key             = "logmsg_id"
