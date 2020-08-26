#-*- coding: UTF-8 -*-c
"""

"""
import logging
from pymongo import MongoClient
from motor import motor_tornado
from tornado.options import options
from tornado.options import define

from iseecore.services import init_service
from iseecore.models import AsyncBaseModel
import setting


@init_service()
class Mongodb():
    """docstring for Mongodb"""

    def __init__(self):

        asyn_client = motor_tornado.MotorClient(setting.MONGO_HOST,
                                                setting.MONGO_PORT)
        AsyncBaseModel.configure(asyn_client)
        logging.error('[init]Mongodb init success')
