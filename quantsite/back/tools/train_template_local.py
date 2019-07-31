#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import os.path
from tornado import ioloop
from tornado.options import options
import tornado.web
import tornado.gen

import logging

from stormed import Connection as StormedConnection
from stormed import Message
from util.md5 import get_file_md5
import uuid
import zipfile

import setting

try:
    from fdfs_client.client import *
    from fdfs_client.exceptions import *
except ImportError:
    import_path = os.path.abspath('../')
    sys.path.append(import_path)
    from fdfs_client.client import *
    from fdfs_client.exceptions import *

from services.tool.db import Mongodb
from iseecore.models import SyncData

from models.files.file_model import File, QiNiuFilesDefine, QiNiuFiles
from models.resource.template_model import LocalTemplate, \
    LocalGroupTemplate, Template
# from models.resource.resource_model import ResourcePicture
from util import qiniu_util

PID_FILE = '/tmp/traing_template.pid'


class TrainTemplate(object):
    base_file_path = os.path.dirname(__file__)
    if not base_file_path == "":
        os.chdir(base_file_path)

    PROJECT_TARGET_COUNT = 12
    NEED_GROUP_TEMPLATE_COUNT = 2

    def __init__(self):
        options.parse_command_line()
        base_file_path = os.path.dirname(__file__)

        fastdfs_conf_path = os.path.join(base_file_path, "fastdfs_client.conf")
        self.fastdfs_client = Fdfs_client(fastdfs_conf_path)

        mongodb = Mongodb()

        mq_host = setting.MQ_HOST
        mq_heartbeat = setting.MQ_HEARTBEAT
        mq_username = setting.MQ_USERNAME
        mq_password = setting.MQ_PASSWORD

        global conn
        conn = StormedConnection(host=mq_host, username=mq_username,
                                 password=mq_password, heartbeat=mq_heartbeat)
        conn.connect(self.on_ampq_connection)
        conn.on_disconnect = self.on_amqp_disconnect

    def on_ampq_connection(self):
        global ch
        ch = conn.channel()
        options["ch"] = ch
        SyncData.configure(ch)

        self.file_class = File()
        self.file_class.need_sync = False
        self.single_template_class = LocalTemplate()
        # self.single_template_class.need_sync = False
        self.group_template_class = LocalGroupTemplate()
        self.group_template_class.need_sync = False
        # res_pic rate need to sync remote
        # to use vuforia res_pic vu_state need not to set
        # self.res_pic_class = ResourcePicture()
        self.template_class = Template()
        self.db_qiniu_files = QiNiuFiles()

        ch.qos(prefetch_count=1)

        ch.consume("yixun_template_queue", self.on_train_finish, no_ack=False)

        print 'consume ok'

        # conn.close(self.done)

    def on_amqp_disconnect(self):
        print "[*]on_amqp_disconnect"
        print "[*]will close tornado"
        try:
            ioloop.IOLoop.instance().stop()
        except AttributeError, e:
            logging.error(e.message)
        finally:
            exit(1)

    @tornado.gen.engine
    def on_train_finish(self, msg):
        if not msg:
            logging.error("train msg is None!")
        logging.error("train_finish headers: %s" % str(msg.headers))
        data = msg.body
        logging.error("similarity: %s" % float(msg.headers["Similarity"]))
        if not data or not len(data):  # check msg.body. anti-moth 20160127
            logging.error("TrainTemplate.on_train_finish: msg.body is None!!")
            msg.ack()
            return

        targetId = msg.headers.get("TargetId", None)
        similarity = msg.headers.get("Similarity", 0)

        logging.error(similarity)

        is_single = True
        if float(similarity) == 0:
            is_single = False

        # 生成文件的md5
        temp_zip_md5 = uuid.uuid4().get_hex() + '_zip'
        temp_file_path = "/tmp/%s.zip" % temp_zip_md5
        with open(temp_file_path, 'wb') as _fo:
            _fo.write(data)

        zip_file_md5 = get_file_md5(temp_file_path)
        if zip_file_md5:
            temp_zip_md5 = zip_file_md5 + '_zip'

        qiniu_md5 = '%s.zip' % temp_zip_md5

        # 查询文件是否已经存在
        old_file = yield tornado.gen.Task(self.find_file_one, temp_zip_md5)

        # 文件不存在，上传到fdfs
        if not old_file:
            # ios 2.2 template modify, only one zip file. anti-moth 20160126
            zip_file_dic = self.fastdfs_client.upload_by_buffer(data, 'zip')
            if not zip_file_dic:
                logging.error("TrainTemplate.on_train_finish: upload zip file error!!")
                msg.ack()
                return

            zip_file_id = zip_file_dic.get('Remote file_id', None)
            yield tornado.gen.Task(self.insert_file, temp_zip_md5, zip_file_id)
            self.sync_ar_template_file(temp_zip_md5, data)

        # 文件存在且已经上传到七牛
        if old_file and old_file.get("qn_key", "") == qiniu_md5:
            pass
        else:
            yield tornado.gen.Task(self.upload_to_qiniu, temp_file_path, qiniu_md5, temp_zip_md5)

        os.remove(temp_file_path)

        template = {
            'temp_zip_md5': temp_zip_md5,
        }

        # unzip and save file
        # targets_zip_file_path = "/tmp/test_result.zip"
        # zip_name = "test"
        # file_base_path = '/tmp/'
        #
        # with open(targets_zip_file_path, "wb") as file_data:
        #     file_data.write(data)
        #
        # zip_file = zipfile.ZipFile(targets_zip_file_path, 'r',
        #                            zipfile.ZIP_DEFLATED)
        # zip_file.extractall(file_base_path)
        #
        # dat_md5 = uuid.uuid4().get_hex() + '_dat'
        # xml_md5 = uuid.uuid4().get_hex() + '_xml'
        # dat_file = '%s%s.dat' % (file_base_path, zip_name)
        # xml_file = '%s%s.xml' % (file_base_path, zip_name)
        # dat_file_id = self.fastdfs_client.upload_by_filename(dat_file)
        # xml_file_id = self.fastdfs_client.upload_by_filename(xml_file)
        #
        # yield tornado.gen.Task(self.insert_file, dat_md5,
        #                        dat_file_id['Remote file_id'])
        # yield tornado.gen.Task(self.insert_file, xml_md5,
        #                        xml_file_id['Remote file_id'])
        #
        # os.remove(dat_file)
        # os.remove(xml_file)
        # os.remove(targets_zip_file_path)

        # template = {
        #     'dat_md5': dat_md5,
        #     'xml_md5': xml_md5,
        # }

        if is_single:  # query exists template first. anti-moth 20160127
            query_tem = {"md5": targetId}
            cur_tem_data = yield tornado.gen.Task(self.single_template_class.get_list, query_tem)
            has_tem = False
            for tem_data in cur_tem_data:
                tem_id = tem_data["template_id"]
                yield tornado.gen.Task(self.template_class.update, {"template_id": tem_id}, {"$set": template})
                has_tem = True
                logging.error("update new template, template_id = %s, temp_zip_md5 = %s" % (tem_id, temp_zip_md5))

            if not has_tem:
                template_id = yield tornado.gen.Task(self.template_class.insert, template)
                single_template = {
                    'md5': targetId,
                    'template_id': template_id,
                    'similarity': similarity
                }
                yield tornado.gen.Task(self.single_template_class.insert,
                                       single_template)

                logging.error("new template %s" % single_template)
            """
            query = {
                'md5': targetId
            }
            update_set = {
                '$set': {
                    'tracker_state': 1,
                    'vu_rate': int(float(similarity) * 5)
                }
            }
            yield tornado.gen.Task(self.res_pic_class.update,
                                   query, update_set, multi=True)
            """
        else:
            template_id = yield tornado.gen.Task(self.template_class.insert, template)
            query = {'local_group_template_id': targetId}
            update_set = {
                '$set': {
                    'template_id': template_id,
                }
            }
            yield tornado.gen.Task(self.group_template_class.update,
                                   query, update_set)

        msg.ack()

    def sync_ar_template_file(self, file_md5, file_content):
        request = Message(
            file_content,
            delivery_mode=2,
            reply_to=options.Sync_Receive_Data_Queue,
            headers={
                "module": "upload_file",
                "md5": file_md5,
                "file_type": "zip",
                "need_narrow": "0",
                "ext": "zip",
            }
        )
        global ch
        ch.publish(
            request,
            exchange=options.Sync_Send_Data_Exchange,
            routing_key=options.Sync_Send_Data_Routing_Key
        )
        logging.info("[*] sync_data upload_ar_template_file md5: %s " % file_md5)

    @tornado.gen.engine
    def upload_to_qiniu(self, file_path, file_name, temp_zip_md5, callback=None):

        qn_token = qiniu_util.get_upload_token(file_name, setting.YIXUN_DEFAULT_FILE_BACKET)
        ret, info = qiniu_util.upload_by_path(qn_token, file_path, file_name)

        if ret.get("key", None) != file_name:
            logging.error('[*] train_template_local upload to qiniu error. file_name = %s ' % file_name)
            return

        # _db_df = QiNiuFilesDefine
        # _document = {
        #     _db_df.QINIU_MD5: file_name,
        #     _db_df.QINIU_KEY: ret["hash"],
        # }
        # yield tornado.gen.Task(self.db_qiniu_files.update, _document, {"$set": _document}, upsert=True)

        query = {
            "md5": temp_zip_md5
        }

        f = {
            "qn_key": file_name
        }

        update_set = {
            "$set": f,
        }

        self.file_class.need_sync = False
        yield tornado.gen.Task(self.file_class.update, query, update_set, upsert=True, safe=True)

        callback((ret, info))

    @tornado.gen.engine
    def insert_file(self, file_md5, file_id, callback=None):
        self.file_class.need_sync = False
        query = {
            "md5": file_md5,
        }
        f = {
            "file_id": file_id,
        }
        update_set = {
            "$set": f,
        }
        yield tornado.gen.Task(self.file_class.update, query,
                               update_set, upsert=True, safe=True)
        logging.error("update file data")
        logging.error('md5: %s, file_id: %s ' % (file_md5, file_id))
        callback(None)

    @tornado.gen.engine
    def find_file_one(self, file_md5, callback=None):
        query = {
            "md5": file_md5
        }

        file_obj = yield tornado.gen.Task(self.file_class.find_one, query)
        callback(file_obj)

    def done(self):
        logging.error("done")
        ioloop.IOLoop.instance().stop()
        conn.close()

if __name__ == '__main__':
    options.log_to_stderr = True
    options.logging = 'debug'
    try:
        TrainTemplate()
        ioloop.IOLoop.instance().start()
    except KeyboardInterrupt:
        ioloop.IOLoop.instance().stop()
        conn.close()
