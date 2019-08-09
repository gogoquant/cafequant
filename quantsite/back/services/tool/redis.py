#-*- coding: UTF-8 -*-
import logging

from tornado.options import options
from tornado.options import define

from iseecore.services import init_service
import setting

import tornadoredis

CONNECTION_POOL = tornadoredis.ConnectionPool(
    max_connections=100,
    wait_for_available=True,
    host=setting.REDIS_HOST,
    port=setting.REDIS_PORT)


@init_service()
class Redis():

    def __init__(self):
        logging.error("redis_client init")
        define(
            "redis_client",
            tornadoredis.Client(connection_pool=CONNECTION_POOL),
            help="redis handler")
        #options['redis_client'] = tornadoredis.Client(connection_pool=CONNECTION_POOL)
        logging.error('[init]redis_client init success')
