#-*- coding: UTF-8 -*-
'''
Xiaolin.Dai
history model
'''
import logging

from models.base_model import *

class History(BaseModel):
    
    def __init__(self):
        BaseModel.__init__(self)
        self.db  = self.connection.user
        self.dao = self.db.history
        self.key = "history_id"
    
    def field(self):
        return {
            "history_id"         : 1,
            "resource_id"        : 1,
            "resource_tag_title" : 1,
            "user_id"            : 1,
            "date"               : 1
        }
