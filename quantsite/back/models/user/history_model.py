#-*- coding: UTF-8 -*-
'''
hongwei.mei
history model
'''
import logging

from iseecore.models import AsyncBaseModel

class History(AsyncBaseModel):
    
    def __init__(self):
        AsyncBaseModel.__init__(self)
        
        table           = "history"
        key             = "history_id"
    
    def field(self):
        return {
            "history_id"         : 1,
            "resource_id"        : 1,
            "resource_tag_title" : 1,
            "user_id"            : 1,
            "date"               : 1
        }
