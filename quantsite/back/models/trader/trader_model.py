#-*- coding: UTF-8 -*-

import simplejson
import hashlib
import logging
import traceback

from iseecore.models import AsyncBaseModel

#交易所
class Exchange(AsyncBaseModel):

    exchange_id    = False     #ID
    
    name            = False    #交易所名
    
    brief            = False    #交易所简介

    content          = False    #文章内容

    accessKey        = False    #通行key

    secretKey        = False    #secret key

    ctime            = False    #创建时间

    mtime            = False    #最后编辑时间

    # meta
    table           = "exchange"
    key             = "exchange_id"    

#策略
class Stragey (AsyncBaseModel):

    stragey_id       = False    #ID

    text             = False    #策略

    ctime            = False    #创建时间

    mtime            = False    #最后编辑时间
    
    # meta
    table           = "stragey"
    key             = "stragey_id"

#交易
class Trader(AsyncBaseModel):

    trader_id        = False    #ID
    
    name             = False    #标题
    

    ctime            = False    #创建时间

    mtime            = False    #最后编辑时间

    # meta
    table           = "trader"
    key             = "trader_id"    
