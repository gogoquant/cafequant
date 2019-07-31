#-*- coding: UTF-8 -*-
'''
Created on 2017-01-03

@author: xiyanxiyan10
@module: topic_model

'''
import simplejson
import hashlib
import logging
import traceback

from iseecore.models import AsyncBaseModel

#文章
class Topic(AsyncBaseModel):
    
    title            = False    #文章标题
    
    brief            = False    #文章简介

    content          = False    #文章内容

    tag_id           = False    #文章类型
    
    file_type        = False    #文章格式类型
    
    user_id          = False    #用户id

    
    published        = False    #是否发表

    reader           = False    #阅读次数

    star             = False    #关注次数 

    mark             = False    #特殊标记,预留


    ctime            = False    #创建时间

    mtime            = False    #最后编辑时间

    # meta
    table           = "topic"
    key             = "topic_id"


#文章评论
class TopicDis(AsyncBaseModel):

    topic_id        = False     #文章ID
    
    user_id         =  False    #用户ID

    content         = False     #评论内容

    father_id       = False     #被引用的评论

    ctime            = False    #创建时间
    
    # meta
    table           = "topicdis"
    key             = "topicdis_id"
