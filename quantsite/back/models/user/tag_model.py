#-*- coding: UTF-8 -*-
'''
Created on 2017-01-03

@author: xiyanxiyan10
@module: tag_model
'''
import simplejson
import hashlib
import logging
import traceback

from iseecore.models import AsyncBaseModel

#文章标签
class Tag(AsyncBaseModel):
    #TODO 不能有一样的标签
    title            = True    #tag 简介
    user_id          = False   #创建tag的用户id 
    
    # meta
    table           = "tag"
    key             = "tag_id"

