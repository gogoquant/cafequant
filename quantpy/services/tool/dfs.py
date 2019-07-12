# -*- coding: UTF-8 -*-
import logging
from tornado.options import options
from iseecore.services import init_service
#from iseecore.fastdfs_utils import DfsUtils


@init_service()
class DFS():
    def __init__(self):
        #TODO ignore dfs
        #options["fastdfs_client"] = DfsUtils()
        logging.error('[init]FastDFS init success')
