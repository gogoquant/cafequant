#-*- coding: UTF-8 -*-c
"""

"""
import logging
import asyncmongo
from pymongo import Connection
from tornado.options import options
from tornado.options import define

from iseecore.services import init_service
from iseecore.models import AsyncBaseModel
import setting


@init_service()
class Mongodb():

    """docstring for Mongodb"""

    def __init__(self):
        asyn_client = asyncmongo.Client(
            pool_id='mydb',
            host=setting.MONGO_HOST,
            port=setting.MONGO_PORT,
            #dbuser=setting.MONGO_USER,
            #dbpass=setting.MONGO_PASS,
            dbname='misakaweb',
            maxcached=150,
            maxconnections=150,
        )
        connection = Connection(setting.MONGO_HOST, setting.MONGO_PORT)
        
        define("async_client", default=asyn_client, help="async connection")
        #options["asyn_client"] = asyn_client
        AsyncBaseModel.configure(asyn_client)

        define("mono_conn", default=connection, help="mongo connection")
        #options["mono_conn"] = connection
        logging.error('[init]Mongodb init success')
