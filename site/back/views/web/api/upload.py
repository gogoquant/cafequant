#-*- coding: UTF-8 -*-
'''
    @brief 文件上传
    @author: lancelot
    @data 2017-3-23

'''

import logging
import tornado.gen
import tornado.web
import pdb
import traceback

from tornado.options import define, options
from iseecore.routes import route

from views.web.base import *
from util.qiniu_util import *

from services.file import FileService

from stormed import Connection
import simplejson
import hashlib

import setting
from pycket.driver import Driver
'''
    上传文件
'''


@route(r"/api/file/add", name="api.file.add")
class FileAddHandler(WebAsyncAuthHandler):
    file_s = FileService()

    #获取七牛使用的参数
    def parse_qiniu_argument(self):
        self.file_name = self.get_argument('file.name', None)
        self.file_path = self.get_argument('file.path', None)
        self.file_md5 = self.get_argument('file.md5', None)
        self.file_type = self.get_argument('file.content_type', None)
        self.file_store = self.get_argument('file.file_store', None)
        self.file_meta = self.file_type

        self.file_width = None
        self.file_height = None
        self.file_size = None

    @tornado.gen.engine
    def _post_(self):
        #pdb.set_trace()

        #获取参数
        self.parse_qiniu_argument()

        query = {
            "md5": self.file_md5,
        }

        logging.info('file from nginx name %s, path %s, md5 %s' %
                     (self.file_name, self.file_path, self.file_md5))

        #查询是否已经有该文件
        f = yield tornado.gen.Task(self.file_s.get_by_m5, self.file_md5)

        #有相同文件则返回该文件id
        if f and len(f) != 0:
            logging.info('file found success')
            self.render_success(msg=f)
        else:
            logging.info('file found fail')

        #@TODO 上传到七牛
        self.upload_to_qiniu(self.file_path, self.file_md5, self.file_md5)

        #记录该文件到数据库
        file_id = yield tornado.gen.Task(self.file_s.add, md5=self.file_md5, meta=self.file_meta, store=self.file_store, \
                width=self.file_width, height=self.file_height, file_name=self.file_name, \
                file_size=self.file_size)

        #返回文件在数据库中的id以及其md5
        self.render_success(msg={"file_id": file_id, "file_md5": self.file_md5})
        return

    #文件上传到七牛
    #@TODO 放入rabbitmq做下半段操作，防止影响服务器效率
    def upload_to_qiniu(self, file_path, file_name, file_md5):
        try:
            #获取上传token
            qn_token = get_upload_token(file_name, setting.DEFAULT_APP_ID)
            logging.error('upload md5 %s, token %s' % (file_md5, qn_token))

            #开始上传
            ret, info = upload_by_path(qn_token, file_path, file_name)
            if ret.get("key", None) != file_name:
                logging.error(
                    '[*] train_template_local upload to qiniu error. file_name = %s '
                    % file_name)
                return

        except Exception, e:
            traceback.print_exc()
        pass


'''
    获取文件
'''


@route(r"/api/file/get", name="api.file.get")
class FileGetHandler(WebHandler):
    file_s = FileService()

    @tornado.gen.engine
    def _post_(self):

        #获取用户数据
        md5 = self.get_argument("md5", None)

        logging.info(md5)

        vfile = yield tornado.gen.Task(self.file_s.get_by_m5, md5)

        self.render_success(msg=vfile)


'''
    获取文件总数
'''


@route(r"/api/file/tot", name="api.file.tot")
class FileGetHandler(WebHandler):
    file_s = FileService()

    @tornado.gen.engine
    def _post_(self):
        query = {}
        tot = yield tornado.gen.Task(self.file_s.count, query)
        self.render_success(msg=tot)


'''
    删除文件
'''


@route(r"/api/file/delete", name="api.file.delete")
class FileGetHandler(WebAsyncAuthHandler):
    file_s = FileService()

    @tornado.gen.engine
    def _post_(self):

        md5 = self.get_argument("md5", "")

        yield tornado.gen.Task(self.file_s.delete, md5)

        self.render_success(msg="success")
