# -*- coding: utf-8 -*-
import os
import pymongo
from tornado.options import define

DEFAULT_APP_ID = 'misaka'
DEFAULT_GROUP_NAME = 'xiyanxiyan10'
APPS_BASE_PATH = 'apps'
DEFAULT_APP_ID = '5268b93b602ede779f92501e'
CORE_APP_NAME = 'misaka'
NEED_SYNC = True
SITE_URL = "10.0.1.162"
MEDIA_BASE_URL = "http://10.0.1.125:9000"
ELA_LOGGING_URL = 'http://10.0.1.120:8080/api/v1/proxy/namespaces/kube-system/services/elasticsearch-logging/'
CDN_URL = "http://7xi53e.com1.z0.glb.clouddn.com"
LOCAL_SITE_URL = "127.0.0.1"

#该网站跳转的域名
MISAKA_DNS = "www.lancelot.top"

LOGIN_URL = '/admin/adminLogin'
ADMINLOGIN_URL = '/user/Login'

#@TODO
YIXUN_DEFAULT_FILE_BACKET = ""
YIXUN_DEFAULT_FILE_BACKET = ""

#降序mongo sort
DESC = pymongo.DESCENDING
#升序mongo sort
ASC = pymongo.ASCENDING

PAGE_SIZE = 20

VUFORIA = [{
    "user_name": "xiyanxiyan10@hotmail.com",
    "user_pass": "meihongwei",
    "user_id": "3363653"
}]

# admin role
ROLE_ADMIN = "admin"
ROLE_EDITOR = "editor"
ROLE_TEST = "test"
ROLE_GUEST = "guest"

# support file ext
EXT_PIC_LIST = ["jpg", "png", "jpeg", "gif"]
EXT_VIDEO_LIST = ["mp4", "flv", "wmv", "m4v", "f4v"]
EXT_AUDIO_LIST = ["mp3", "ogg"]
EXT_3D_LIST = ["unity3d", "obj", "dae"]
FILE_TYPE_EQ_EXTS = [
    'unity3d', 'tracker_data', 'data', 'dat', 'xml', 'mp3', 'ipa', 'zip', 'fbx'
]

#redis conf file
REDIS_HOST = os.getenv("REDIS_MASTER_PORT_6379_TCP_ADDR", '127.0.0.1')
REDIS_PORT = int(os.getenv("REDIS_MASTER_PORT_6379_TCP_PORT", 6379))
REDIS_SESSIONS = 10

# api use
API_STATE = True
VERSION_STATE = ''

# amqp
MQ_HOST = os.getenv("MQ_PORT_5672_TCP_ADDR", "10.0.1.37")
MQ_HEARTBEAT = 30
MQ_USERNAME = os.getenv("MQ_USERNAME", "lancelot")
MQ_PASSWORD = os.getenv("MQ_PASSWORD", "lancelot")

# mongo
MONGO_HOST = "127.0.0.1"
MONGO_PORT = int(27017)
MONGO_USER = "mysite"
MONGO_PASS = "password"
MONGO_DB = "mysite"
MONGO_ID = "my_id"

#CDN
IS_CDN = False
STATIC_CDN_URL = "http://www.baidu.com"
STATIC_URL = "http://www.lancelot.top"

#静态文件路径
MISAKA_STATIC_DIR = "/home/ubuntu/project/misakaWeb/python-web/static"
MISAKA_TEMPLATE_DIR = "/home/ubuntu/project/misakaWeb/python-web/template"
MISAKA_DEFAULT_FILE_BACKET = "test"

#检索回调rpc
define(
    "Search_CallBack_Exchange",
    default="search_callback",
    help="misaka search callback exchange")

#消息同步
# send message by search_sync exchange
define(
    "Sync_Send_Data_Exchange",
    default="sync_send",
    help="Sync_Send_Data_Exchange")
define(
    "Sync_Send_Data_Routing_Key",
    default="misaka_sync_db.data",
    help="Sync_Send_Data_Routing_Key")

define(
    "Sync_Receive_Data_Exchange",
    default="sync_recive",
    help="Sync_Receive_Data_Exchange")
define(
    "Sync_Receive_Data_Queue",
    default="misaka_sync_db_queue",
    help="Sync_Receive_Data_Queue")
define(
    "Sync_Receive_Data_Routing_Key",
    default="misaka_sync_db.*",
    help="Sync_Receive_Data_Routing_Key")

#用户信息变更queue
define(
    "User_Status_Exchange", default="user_status", help="user status exchange")
define("User_Status_Queue", default="user_status", help="user status queue")
