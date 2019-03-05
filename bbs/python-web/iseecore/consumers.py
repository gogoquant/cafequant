#-*- coding: UTF-8 -*-

import tornado.web
import logging
import traceback
import simplejson

class consumer:
    _method_dict_ = {}
    
    @classmethod
    def add_method(cls,req_id , method):
        cls._method_dict_[req_id] = method
    
    @classmethod
    def consume(cls,msg):
        try:
            cls._consume_(msg)
        except (Exception, IndexError), e:
            exc_msg = traceback.format_exc()
            logging.error(exc_msg)
    
    @classmethod
    def _consume_(cls,msg):
        headers = msg.headers
        req_id = None
        if 'request_id' in headers:
            req_id = headers['request_id']
        elif 'req_id' in headers:
            req_id = headers['req_id']
        else:
            body = simplejson.loads(str(msg.body).replace("\\", ""))
            if 'request_id' in body:
                req_id = body['request_id']
            elif 'req_id' in body:
                req_id = body['req_id']
        if req_id:
            if req_id in cls._method_dict_:
                method = cls._method_dict_.pop(req_id)
                method(msg)
    @classmethod
    def reset(cls):
        cls._method_dict_ = {}

