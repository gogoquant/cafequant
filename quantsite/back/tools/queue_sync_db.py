#!/usr/bin/python
# -*- coding: utf-8 -*-

import os
from tornado import ioloop
from tornado.options import options
import tornado.gen
import cPickle
import Image

import logging
import simplejson

from stormed import Connection as StormedConnection
import traceback
import setting

try:
    from fdfs_client.client import Fdfs_client
    from fdfs_client.exceptions import DataError
except ImportError:
    import sys
    import_path = os.path.abspath('../')
    sys.path.append(import_path)
    from fdfs_client.client import Fdfs_client
    from fdfs_client.exceptions import DataError

from iseecore.models import SyncData
from iseecore.models import METHOD_INSERT, METHOD_UPDATE, METHOD_DELETE
from models.files.file_model import File

from services.tool.db import Mongodb


import tornadoredis
CONNECTION_POOL = tornadoredis.ConnectionPool(max_connections=100,
                                              wait_for_available=True,
                                              host=setting.REDIS_HOST,
                                              port=setting.REDIS_PORT)
logging.error("redis_client init")
redis_client = tornadoredis.Client(connection_pool=CONNECTION_POOL)
logging.error('[init]redis_client init success')


base_file_path = os.path.dirname(__file__)

if not base_file_path == "":
    os.chdir(base_file_path)


def dict_to_object(d):  # for simplejson loads user-defined objects. anti-moth 20160226
    if '__class__' in d:
        class_name = d.pop('__class__')
        module_name = d.pop('__module__')
        module = __import__(module_name)
        try:
            getattr(module, class_name)
            inst = cPickle.loads(str(d.pop('__data__')))

        except AttributeError:
            logging.error('[*] dict_to_object error. class_name = %s, module_name = %s ' % (class_name, module_name))
            inst = d

    else:
        inst = d

    return inst


@tornado.gen.engine
def sync_data_receive(response, callback=None):
    global db_conn
    if not response:
        logging.error("callback response is None!")
        response.ack()
    else:
        module = response.headers["module"]
        """ if module is upload_file"""
        if module == "upload_file":
            yield tornado.gen.Task(upload_file, response)
            response.ack()
            if callback:
                callback(None)
            return

        """ if module is clear_cache"""
        if module == "clear_cache":
            yield tornado.gen.Task(clear_cache, response)
            response.ack()
            if callback:
                callback(None)
            return

        className = response.headers["className"]
        method = response.headers["method"]

        args = response.body
        args = simplejson.loads(args, object_hook=dict_to_object)

        import sys
        sys.path.append("models")

        if method == METHOD_INSERT:
            yield tornado.gen.Task(sync_insert, module, className, args)
        
        if method == METHOD_UPDATE:
            yield tornado.gen.Task(sync_update, response, module,
                                   className, args)
        if method == METHOD_DELETE:
            yield tornado.gen.Task(sync_delete, module, className, args)

        if className == "Instance":
            yield tornado.gen.Task(operate_cate_group, method, args)

    response.ack()
    if callback:
        callback(None)

FILE_IMAGE = "image"
FILE_VIDEO = "video"
upload_file_tmp = "upload_file.tmp"


@tornado.gen.engine
def upload_file(response, callback=None):
    md5 = response.headers["md5"]
    file_type = response.headers["file_type"]
    binary_content = response.body
    # need_narrow = response.headers["need_narrow"]
    ext = response.headers["ext"]

    logging.error("upload_file response.headers:%s" % str(response.headers))
    width = None
    height = None

    # in the process, do not need to do any image_resize, video_compress
    # and unity generate etc, only need to save file to local module.
    # in addition, this only need to get width and height for image search
    if file_type == FILE_IMAGE:
        with open("upload_file.tmp", "w+") as file_data:
            file_data.write(binary_content)
        try:
            img = Image.open("upload_file.tmp")
            width = img.size[0]
            height = img.size[1]
            os.remove("upload_file.tmp")
        except IOError, e:
            logging.error(e.message)
            pass

    """ upload to local dfs """
    file_id = None
    try:
        file_temp = fastdfs_client.upload_by_buffer(binary_content, ext)
        file_id = file_temp["Remote file_id"]
    except DataError, e:
        logging.error(e.message)
    if not file_id:
        return

    ''' Insert file to local file'model '''
    file_model = File()
    file_model.need_sync = False
    query = {
        "md5": md5,
    }
    f = {
        "file_id": file_id,
    }
    if width:
        f["width"] = width
    if height:
        f["height"] = height
    update_set = {
        "$set": f,
    }
    try:
        yield tornado.gen.Task(file_model.update, query, update_set,
                               upsert=True, safe=True)
    except Exception, e:
        logging.error(e.message)
        exc_msg = traceback.format_exc()
        logging.error(exc_msg)

    """ image resize
    msg = Message(
        "resize",
        delivery_mode=2,
        headers={
            "image_name": md5,
            "FileId": file_id,
            "file_type": file_type,
            "need_narrow": need_narrow
        }
    )
    ch.publish(
        msg,
        exchange=options.Image_Resize_Exchange,
        routing_key=options.Image_Resize_Queue
    )
    """
    if callback:
        callback(None)


@tornado.gen.engine
def clear_cache(response, callback=None):

    # clear resource pages cache
    if response.headers['clear_type'] == 'resource_pages':
        clear_data = simplejson.loads(response.headers['clear_data'])
        md5 = clear_data["md5"]
        h5_md5 = clear_data["h5_md5"]

        logging.error("clear_resource_pages_cache clear_data:%s" % str(clear_data))

        result = yield tornado.gen.Task(redis_client.exists, md5)
        if result:
            logging.info("clear_resource_pages_cache basic info (md5:%s) success!" % md5)
            yield tornado.gen.Task(redis_client.delete, md5)

        h5_result = yield tornado.gen.Task(redis_client.exists, h5_md5)
        if h5_result:
            logging.info("clear_resource_pages_cache h5 info (h5_md5:%s) success!" % h5_md5)
            yield tornado.gen.Task(redis_client.delete, h5_md5)

    if callback:
        callback(None)


@tornado.gen.engine
def sync_insert(module, className, args, callback=None):
    if args:
        method = METHOD_INSERT

        logging.error("sync_insert module:%s,className:%s,args:%s" % (
            module, className, str(args)))
        try:
            import sys
            sys.path.append("models")
            model = __import__(module, fromlist=[className])
        except (ImportError, Exception), e:
            logging.error('[*] sync_insert import. e = %s, class_name = %s, module_name = %s ' % (e, className, module))
            model = __import__(module.split(".")[-1], fromlist=[className])

        model_class = getattr(model, className)()
        model_class.need_sync = False

        method = getattr(model_class, method)

        try:
            yield tornado.gen.Task(method, args)
        except Exception, e:
            logging.error(e.message)
            exc_msg = traceback.format_exc()
            logging.error(exc_msg)
    if callback:
        callback(None)


@tornado.gen.engine
def sync_update(response, module, className, args, callback=None):
    method = METHOD_UPDATE

    args = tuple(args)
    logging.error("sync_update module:%s,className:%s,args:%s" % (
        module, className, str(args)))

    upsert = False
    safe = True
    _upsert_ = response.headers.get("upsert", 0)
    _safe_ = response.headers.get("safe", 0)
    if _upsert_ == "1":
        upsert = True
    if _safe_ == "0":
        safe = True

    try:
        import sys
        sys.path.append("models")
        model = __import__(module, fromlist=[className])
    except (ImportError, Exception), e:
        logging.error('[*] sync_update import. e = %s, class_name = %s, module_name = %s ' % (e, className, module))
        model = __import__(module.split(".")[-1], fromlist=[className])

    model_class = getattr(model, className)()
    model_class.need_sync = False

    method = getattr(model_class, method)

    try:
        yield tornado.gen.Task(method, args[0], args[1],
                               upsert=upsert, safe=safe)
    except Exception, e:
        logging.error(e.message)
        exc_msg = traceback.format_exc()
        logging.error(exc_msg)
    if callback:
        callback(None)


@tornado.gen.engine
def sync_delete(module, className, args, callback=None):
    method = METHOD_DELETE

    logging.error("sync_delete module:%s,className:%s,args:%s" % (
        module, className, str(args)))
    try:
        import sys
        sys.path.append("models")
        model = __import__(module, fromlist=[className])
    except (ImportError, Exception), e:
        logging.error('[*] sync_delete import. e = %s, class_name = %s, module_name = %s ' % (e, className, module))
        model = __import__(module.split(".")[-1], fromlist=[className])

    model_class = getattr(model, className)()
    model_class.need_sync = False

    method = getattr(model_class, method)

    try:
        yield tornado.gen.Task(method, args)
    except Exception, e:
        logging.error(e.message)
        exc_msg = traceback.format_exc()
        logging.error(exc_msg)
    if callback:
        callback(None)


@tornado.gen.engine
def operate_cate_group(method, args, callback=None):
    logging.error("instance sync ok")
    from services.instance import InstanceService
    instance_class = InstanceService()
    if method == METHOD_INSERT:
        groupname = args["groupname"]
        search_instance_count = args["search_instance_count"]
        training_instance_count = args["training_instance_count"]
        yield tornado.gen.Task(
            instance_class.handle_instance,
            ch,
            "new",
            groupname,
            search_instance_count,
            training_instance_count
        )
    elif method == METHOD_UPDATE:
        groupname = args[0]["groupname"]
        update_set = args[1].get('$set')
        search_instance_count = update_set["search_instance_count"]
        training_instance_count = update_set["training_instance_count"]
        yield tornado.gen.Task(
            instance_class.handle_instance,
            ch,
            "modify",
            groupname,
            search_instance_count,
            training_instance_count
        )
    else:
        groupname = args["groupname"]
        yield tornado.gen.Task(
            instance_class.handle_instance,
            ch,
            "del",
            groupname,
            0,
            0
        )

    if callback:
        callback(None)


def on_amqp_connection():
    global ch
    ch = conn.channel()

    """" declare sync receive data queue and exchange """

    ch.queue_declare(queue=options.Sync_Receive_Data_Queue, durable=True)

    ch.queue_bind(
        exchange=options.Sync_Receive_Data_Exchange,
        queue=options.Sync_Receive_Data_Queue,
        routing_key=options.Sync_Receive_Data_Routing_Key)

    ch.queue_bind(
        exchange=options.BGP_CallBack_Receive_Exchange,
        queue=options.Sync_Receive_Data_Queue,
        routing_key=options.Sync_Info_Routing_Key)

    ch.consume(options.Sync_Receive_Data_Queue, sync_data_receive, no_ack=False)
    SyncData.configure(ch)


def on_amqp_disconnect():
    logging.error("[*]on_amqp_disconnect")
    logging.error("[*]will close tornado")
    try:
        tornado.ioloop.IOLoop.instance().stop()
    except AttributeError, e:
        logging.error(e.message)
    finally:
        exit(1)

ch = None
options.log_to_stderr = True
options.parse_command_line()
logging.basicConfig()

base_file_path = os.path.dirname(__file__)

mq_host = setting.MQ_HOST
mq_username = setting.MQ_USERNAME
mq_password = setting.MQ_PASSWORD
mq_heartbeat = setting.MQ_HEARTBEAT

global conn
conn = StormedConnection(host=mq_host, username=mq_username,
                         password=mq_password, heartbeat=mq_heartbeat)
conn.connect(on_amqp_connection)
conn.on_disconnect = on_amqp_disconnect

fastdfs_conf_path = os.path.join(base_file_path, "fastdfs_client.conf")
fastdfs_client = Fdfs_client(fastdfs_conf_path)

Mongodb()

try:
    ioloop.IOLoop.instance().start()
except KeyboardInterrupt:
    ch.close()
    conn.close()
