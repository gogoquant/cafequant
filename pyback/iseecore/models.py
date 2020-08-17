#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
baseModel
"""
import logging
import time
import _pickle as cPickle
import simplejson
import pymongo
import sys
import os
import tornado.gen
import pdb

from bson.objectid import ObjectId
from tornado.options import options
#from stormed import Message

import traceback

__all__ = ['AsyncBaseModel', 'ModelException']

MONGODB_ID = "_id"
DELETE_FLAG = 'delete_flag'  # '1' 已经删除，'0'或者没有该字段表示没有删除


class ModelException(Exception):
    pass


class AsyncBaseModel(object):
    need_sync = False

    @classmethod
    def configure(cls, client):
        cls.async_client = client

    def __init__(self):
        if not hasattr(self, "async_client"):
            raise NotImplementedError("Must configure an async_client.")
        self.className = self.__class__.__name__
        self.module = self.__module__
        if hasattr(self, 'db'):
            db = getattr(self, 'db')
        else:
            db = self.module.split('.')[1]

        if hasattr(self, 'table'):
            table = getattr(self, 'table')
        else:
            table = self.className.lower()

        if hasattr(self, 'key'):
            key = getattr(self, 'key')
        else:
            key = table + '_id'

        self._db_ = db
        self._table_ = table

        self.dao = self.async_client[table][db]
        self._key_ = key

        #self.sync_class = SyncData()

    def exist(self, query, callback):
        one = self.find_one(query)
        if one:
            callback(True)
        else:
            callback(False)

    def unique(self, query, callback=None):
        one = self.find_one(query)
        if one:
            callback(False)
        else:
            callback(True)

    def get(self, value, key=None, callback=None):
        if not key:
            key = self._key_
        callback(self.find_one({key: value}))

    #@tornado.gen.engine
    def get_all(self, query={}, callback=None):
        response_list = []
        conditions = {}
        fields = {}
        sorts = []
        count = 9999
        if "conditions" in query:
            conditions = query["conditions"]
        if "fields" in query:
            fields = query["fields"]
        if "sorts" in query:
            sorts = query["sorts"]
        if 'conditions' not in query:
            conditions = query
        # date_sort_flg = False
        # add conditon delete_flag is not 1
        conditions[DELETE_FLAG] = {'$ne': '1'}
        sorts = [
            [x if x != self._key_ else MONGODB_ID for x in y] for y in sorts
        ]
        if not sorts:
            sorts.append([MONGODB_ID, pymongo.DESCENDING])
        if self._key_ in conditions:
            conditions[MONGODB_ID] = conditions[self._key_]
            del conditions[self._key_]
        if self._key_ in fields:
            fields[MONGODB_ID] = fields[self._key_]
            del fields[self._key_]
        if not fields.keys():
            fields = None
        result_list, error = yield tornado.gen.Task(
            self.dao.find,
            spec=conditions,
            fields=fields,
            sort=sorts,
            limit=count)
        if error is None:
            logging.error(error)

        result_list = result_list[0]
        for result in result_list:
            if MONGODB_ID in result:
                result[self._key_] = str(result[MONGODB_ID])
                del result[MONGODB_ID]
            response_list.append(result)
        del conditions[DELETE_FLAG]
        callback(response_list)

    #@tornado.gen.engine
    def get_list_all(self, query={}, callback=None):
        response_list = []
        conditions = {}
        fields = {}
        sorts = []
        pos = 0
        count = 9999
        if "conditions" in query:
            conditions = query["conditions"]
        if "pos" in query:
            pos = query["pos"]
        if "count" in query:
            count = query["count"]
        if "fields" in query:
            fields = query["fields"]
        if "sorts" in query:
            sorts = query["sorts"]
        if 'conditions' not in query:
            conditions = query

        sorts = [
            [x if x != self._key_ else MONGODB_ID for x in y] for y in sorts
        ]
        if not sorts:
            sorts.append([MONGODB_ID, pymongo.DESCENDING])

        pos = int(pos)
        count = int(count)
        if self._key_ in conditions:
            conditions[MONGODB_ID] = conditions[self._key_]
            del conditions[self._key_]
        if self._key_ in fields:
            fields[MONGODB_ID] = fields[self._key_]
            del fields[self._key_]
        if not fields.keys():
            fields = None
        result_list, _ = yield tornado.gen.Task(
            self.dao.find,
            spec=conditions,
            fields=fields,
            limit=count,
            skip=pos,
            sort=sorts)
        result_list = result_list[0]
        for result in result_list:
            try:
                if MONGODB_ID in result:
                    result[self._key_] = str(result[MONGODB_ID])
                    del result[MONGODB_ID]
                response_list.append(result)
            except Exception as e:
                logging.error(e.message)
                logging.error("result: %s, self._key_: %s " %
                              (str(result), self._key_))

        callback(response_list)

    #@tornado.gen.engine
    def get_list(self, query={}, callback=None):
        response_list = []
        conditions = {}
        fields = {}
        sorts = []
        pos = 0
        count = 9999
        if "conditions" in query:
            conditions = query["conditions"]
        if "pos" in query:
            pos = query["pos"]
        if "count" in query:
            count = query["count"]
        if "fields" in query:
            fields = query["fields"]
        if "sorts" in query:
            sorts = query["sorts"]
        if 'conditions' not in query:
            conditions = query
        # add conditon delete_flag is not 1
        conditions[DELETE_FLAG] = {'$ne': '1'}

        sorts = [
            [x if x != self._key_ else MONGODB_ID for x in y] for y in sorts
        ]
        if not sorts:
            sorts.append([MONGODB_ID, pymongo.DESCENDING])

        pos = int(pos)
        count = int(count)
        if self._key_ in conditions:
            conditions[MONGODB_ID] = conditions[self._key_]
            del conditions[self._key_]
        if self._key_ in fields:
            fields[MONGODB_ID] = fields[self._key_]
            del fields[self._key_]
        if not fields.keys():
            fields = None
        result_list, error = yield tornado.gen.Task(
            self.dao.find,
            spec=conditions,
            fields=fields,
            limit=count,
            skip=pos,
            sort=sorts)
        if error is None:
            logging.error(error)

        result_list = result_list[0]
        for result in result_list:
            try:
                if MONGODB_ID in result:
                    result[self._key_] = str(result[MONGODB_ID])
                    del result[MONGODB_ID]
                response_list.append(result)
            except Exception as e:
                logging.error(e.message)
                logging.error("result: %s, self._key_: %s " %
                              (str(result), self._key_))

        del conditions[DELETE_FLAG]
        callback(response_list)

    #@tornado.gen.engine
    def find_one(self, spec, fields=None, callback=None):
        if self._key_ in spec:
            spec[MONGODB_ID] = spec[self._key_]
            del spec[self._key_]
        # add conditon delete_flag is not 1
        spec[DELETE_FLAG] = {'$ne': '1'}

        model = yield self.dao.find_one(spec)

        if model is None:
            callback(None)
            return

        del spec[DELETE_FLAG]
        callback(model)

    #@tornado.gen.engine
    def get_by_id(self, key_value, callback=None):
        result = yield tornado.gen.Task(self.find_one, {self._key_: key_value})
        callback(result)

    #@tornado.gen.engine
    def count(self, query={}, delete_flag=True, callback=None):
        try:
            if self._key_ in query:
                query[MONGODB_ID] = query[self._key_]
                del query[self._key_]
            # add conditon delete_flag is not 1
            if delete_flag:
                query[DELETE_FLAG] = {'$ne': '1'}

            responses, error = yield tornado.gen.Task(
                self.dao.find,
                spec=query,
                fields={self._key_: 1},
                sort=[[self._key_, 1]])
            if error is None:
                logging.error(error)

            del query[DELETE_FLAG]
            callback(len(responses[0]))
        except Exception as e:
            exc_msg = traceback.format_exc()
            logging.error(exc_msg)
            raise ModelException("count error:%s,e:%s" % (str(query), str(e)))

    #@tornado.gen.engine
    def insert(self,
               doc,
               manipulate=True,
               safe=True,
               check_keys=True,
               callback=None,
               **kwargs):

        send_time = int(time.time())
        doc["date"] = send_time
        doc["last_modify"] = send_time

        # 如果doc中有数组，并且数组中是dict,则向其中加入date和id
        for (_, value) in doc.iteritems():
            if not type(value) == list:
                continue
            if not len(value):
                continue
            if type(value[0]) == dict:
                # if id not exists, add id, date etc for sub doc
                if "id" not in value:
                    for d in value:
                        d["id"] = self.get_id()
                        d["date"] = send_time
                        d["last_modify"] = send_time
        # for sync_data, when _id exists, do not need to get_id()
        if MONGODB_ID not in doc:
            if self._key_ not in doc:
                key_value = self.get_id()
            else:
                key_value = doc[self._key_]
                del doc[self._key_]
        else:
            key_value = doc[MONGODB_ID]
            if self._key_ in doc:
                del doc[self._key_]
        doc[MONGODB_ID] = key_value

        res = yield self.dao.insert_one(doc)
        if res is None:
            callback(None)

        # sync data
        if self.need_sync:
            self.sync_insert_data(doc)

        if callback:
            callback(key_value)

    #@tornado.gen.engine
    def update(self,
               spec,
               document,
               upsert=False,
               manipulate=False,
               safe=True,
               multi=False,
               callback=None,
               **kwargs):

        if self._key_ in spec:
            spec[MONGODB_ID] = spec[self._key_]
            del spec[self._key_]
        # add conditon delete_flag is not 1
        spec[DELETE_FLAG] = {'$ne': '1'}
        old_update_data = document.get("$set", None)
        if old_update_data:
            old_update_data["last_modify"] = int(time.time())
            if DELETE_FLAG in old_update_data:
                del spec[DELETE_FLAG]
            document['$set'] = old_update_data

        unset_data = document.get("$unset", None)
        if unset_data and DELETE_FLAG in unset_data:
            del spec[DELETE_FLAG]

        result = yield tornado.gen.Task(
            self.dao.update,
            spec,
            document,
            upsert=upsert,
            manipulate=manipulate,
            safe=safe,
            multi=multi,
            **kwargs)

        if self.need_sync:
            self.sync_update_data(spec, document, upsert, safe)

        if callback:
            callback(result)

    #@tornado.gen.engine
    def delete(self, spec, safe=True, callback=None, **kwargs):

        if self._key_ in spec:
            spec[MONGODB_ID] = spec[self._key_]
            del spec[self._key_]

        # add conditon delete_flag is not 1
        spec[DELETE_FLAG] = {'$ne': '1'}
        update_set = {'$set': {DELETE_FLAG: '1'}}

        result = yield tornado.gen.Task(
            self.dao.update, spec, update_set, multi=True)

        if self.need_sync:
            self.sync_delete_data(spec)

        if callback:
            callback(result)

    #@tornado.gen.engine
    def find_and_modify(self,
                        spec,
                        document,
                        upsert=False,
                        manipulate=False,
                        safe=True,
                        multi=False,
                        callback=None,
                        **kwargs):

        if self._key_ in spec:
            spec[MONGODB_ID] = spec[self._key_]
            del spec[self._key_]
        # add conditon delete_flag is not 1
        spec[DELETE_FLAG] = {'$ne': '1'}
        old_update_data = document.get("$set", None)
        if old_update_data:
            old_update_data["last_modify"] = int(time.time())
            document['$set'] = old_update_data

        from bson.son import SON

        command = SON([('findAndModify', self._table_), ('query', spec),
                       ('update', document), ('upsert', False), ('new', True)])

        command.update(kwargs)

        result = yield tornado.gen.Task(
            self.async_client.connection("$cmd", self._db_).find_one,
            command,
            _must_use_master=True,
            _is_command=True)

        flag = result[0][0]['value']
        if flag and self.need_sync:
            self.sync_update_data(spec, document)

        callback(flag)

    def get_id(self):
        return str(ObjectId())

    def sync_insert_data(self,
                         doc):  # for mongodb $ operator anti-moth 20160226
        if MONGODB_ID in doc and isinstance(doc[MONGODB_ID], ObjectId):
            doc[MONGODB_ID] = str(doc[MONGODB_ID])
        #self.sync_class.send_insert(self.module, self.className, doc)
        pass

    def sync_update_data(self, spec, doc, upsert=False, safe=False):
        if MONGODB_ID in spec and isinstance(spec[MONGODB_ID], ObjectId):
            spec[MONGODB_ID] = str(spec[MONGODB_ID])
        #self.sync_class.send_update(self.module, self.className, spec, doc, upsert, safe)
        pass

    def sync_delete_data(self, spec):
        if MONGODB_ID in spec and isinstance(spec[MONGODB_ID], ObjectId):
            spec[MONGODB_ID] = str(spec[MONGODB_ID])
        #self.sync_class.send_delete(self.module, self.className, spec)
        pass


METHOD_INSERT = "insert"
METHOD_UPDATE = "update"
METHOD_DELETE = "delete"


def convert_to_builtin_type(
        obj):  # for simplejson dumps user-defined objects. anti-moth 20160226
    # Convert objects to a dictionary of their representation
    class_name = obj.__class__.__name__
    module_name = obj.__module__
    module = __import__(module_name)
    try:
        getattr(module, class_name)
        return {
            '__class__': class_name,
            '__module__': module_name,
            '__data__': cPickle.dumps(obj)
        }
    except AttributeError:
        return obj


# sync data by message queue
class SyncData(object):

    def __init__(self):
        if not hasattr(self, "channel"):
            raise NotImplementedError("Must configure an channel.")

    @classmethod
    def configure(cls, channel):
        cls.channel = channel

    def send_insert(self, module, className, doc):
        pass
        '''
        method = METHOD_INSERT

        doc = simplejson.dumps(doc, default=convert_to_builtin_type)
        request = Message(
            doc,
            delivery_mode=2,
            reply_to=options.Sync_Receive_Data_Queue,
            headers={
                "module": module,
                "className": className,
                "method": method,
            })
        self.channel.publish(
            request,
            exchange=options.Sync_Send_Data_Exchange,
            routing_key=options.Sync_Send_Data_Routing_Key)
        logging.info("[*] sync_data insert module:%s,className:%s,doc:%s" %
                     (module, className, doc))
        '''

    def send_update(self,
                    module,
                    className,
                    spec,
                    doc,
                    upsert=False,
                    safe=False):
        method = METHOD_UPDATE
        pass
        '''
        args = simplejson.dumps((spec, doc), default=convert_to_builtin_type)
        _upsert_ = "0"
        _safe_ = "0"
        if upsert:
            _upsert_ = "1"
        if safe:
            _safe_ = "1"

        request = Message(
            args,
            delivery_mode=2,
            reply_to=options.Sync_Receive_Data_Queue,
            headers={
                "module": module,
                "className": className,
                "method": method,
                "upsert": _upsert_,
                "safe": _safe_
            })
        self.channel.publish(
            request,
            exchange=options.Sync_Send_Data_Exchange,
            routing_key=options.Sync_Send_Data_Routing_Key)
        logging.info(
            "[*] sync_data update module:%s,className:%s,spec:%s,doc:%s,upsert:%s,safe:%s"
            % (module, className, spec, doc, upsert, safe))
        '''

    def send_delete(self, module, className, spec):
        pass
        '''
        method = METHOD_DELETE
        spec = simplejson.dumps(spec, default=convert_to_builtin_type)
        request = Message(
            spec,
            delivery_mode=2,
            reply_to=options.Sync_Receive_Data_Queue,
            headers={
                "module": module,
                "className": className,
                "method": method,
            })
        self.channel.publish(
            request,
            exchange=options.Sync_Send_Data_Exchange,
            routing_key=options.Sync_Send_Data_Routing_Key)
        logging.info("[*] sync_data delete module:%s,className:%s,spec:%s" %
                     (module, className, spec))
        '''


class MessageAsyncBaseModel(AsyncBaseModel):  # anti-moth 20160223
    """统一消息系统中由upsert和insert中生成的objectId"""

    def get_id(self):
        return ObjectId()

    def sync_insert_data(self, doc):
        #self.sync_class.send_insert(self.module, self.className, doc)
        pass

    def sync_update_data(self, spec, doc, upsert=False, safe=False):
        #self.sync_class.send_update(self.module, self.className, spec, doc, upsert, safe)
        pass

    def sync_delete_data(self, spec):
        #self.sync_class.send_delete(self.module, self.className, spec)
        pass
