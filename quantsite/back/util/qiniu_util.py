#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
@author:lancelot
@date:2017/3/24
"""

import os
import logging
import qiniu
import mimetypes
import setting

ACCESS_KEY = "test"
SECRET_KEY = "test"

DEFAULT_APP = setting.DEFAULT_APP_ID

TEST_BUCKET = "testidealsee"

BUCKETS = {
    "test": "testidealsee",
    setting.DEFAULT_APP_ID: "xiyanxiyan10",
    "huanshi": "yixun",
    "yixunfiles": "yixunfiles",
    "testyixun": "testyixun",
}

BUCKETS_URL = {
    "test": "http://7xiou8.com2.z0.glb.qiniucdn.com/",
    setting.DEFAULT_APP_ID: "http://7xofx2.com2.z0.glb.qiniucdn.com/",
    "huanshi": "http://7xofx2.com2.z0.glb.qiniucdn.com/",
    "yixunfiles": "http://7xsom5.com2.z0.glb.qiniucdn.com/",
}


def get_content_type(filename):
    return mimetypes.guess_type(filename)[0] or 'application/octet-stream'


qiniuAuth = qiniu.Auth(ACCESS_KEY, SECRET_KEY)
qiniu_bucket_manager = qiniu.BucketManager(qiniuAuth)


# 从七牛获取上传token
# 用于客户端请求上传,服务器端自主上传到七牛
def get_upload_token(filename, app_id=setting.DEFAULT_APP_ID):
    bucket = BUCKETS[app_id]
    expires = 3600
    policy = {}
    strict_policy = True
    # ext = filename.split(".")[-1]
    key = None

    # if ext in setting.EXT_PIC_LIST:
    #     # policy['persistentOps'] = 'imageView2/0/q/100/w/320'
    #     pass
    # elif ext in setting.EXT_VIDEO_LIST:
    #     policy['persistentOps'] = 'imageView2/0/q/100/w/320'

    token = qiniuAuth.upload_token(bucket, key=key, expires=expires, policy=policy, strict_policy=strict_policy)
    return token


# 使用文件路径上传
# 如果文件大小大于4m，自动断点续传
# 如果文件名已存在将返回错误
def upload_by_path(token, file_path, filename=None, socket_proxy_host=None, socket_proxy_port=None):

    if socket_proxy_host and socket_proxy_port:
        import socks
        import socket
        socks.setdefaultproxy(socks.PROXY_TYPE_SOCKS5, socket_proxy_host, socket_proxy_port)
        # socks.setdefaultproxy(socks.PROXY_TYPE_SOCKS5, "10.0.1.62", 10080)
        # socks.setdefaultproxy(socks.PROXY_TYPE_HTTP, "10.0.1.62", 3128)
        # socks.setdefaultproxy(socks.PROXY_TYPE_HTTPS, "10.0.1.62", 3128)
        socket.socket = socks.socksocket

    content_type = get_content_type(file_path)
    ret, info = qiniu.put_file(token, filename, file_path, mime_type=content_type, check_crc=True)
    if socket_proxy_host and socket_proxy_port:
        socks.setdefaultproxy()
        socket.socket = socks.socksocket

    return ret, info


# 使用文件流的方式上传文件
def upload_by_stream(token, file_name, data, file_size=None, file_key=None):
    mime_type = get_content_type(file_name)
    logging.error("mime_type:%s" % mime_type)
    ret, info = qiniu.put_data(token, file_key, data, mime_type=mime_type, check_crc=True)
    return ret, info


# 使用文件流的方式上传文件
def upload_by_stream(token, file_name, data, file_size=None, file_key=None):
    mime_type = get_content_type(file_name)
    logging.error("mime_type:%s" % mime_type)
    ret, info = qiniu.put_data(token, file_key, data, mime_type=mime_type, check_crc=True)
    return ret, info


if __name__ == '__main__':

    # upload by path
    token = get_upload_token('test.mp4', app_id="test")
    logging.error("token")
    logging.error(token)

    path = '/home/wangande/testupload.mp4'
    filename = "testupload.mp4"

    file_size = os.stat(path).st_size
    logging.error('size')
    logging.error(file_size)
    ret, info = upload_by_path(token, path, file_size, filename="testupload.mp4")
    logging.error("ret:%s" % ret)
    logging.error("info:%s" % info)

    # upload by stream
    with open(path, 'rb') as f:
        ret, info = upload_by_stream(token, filename, f.read(), file_size=file_size)
        logging.error("ret:%s" % ret)
        logging.error("info:%s" % info)

        if info.status_code == 200:
            logging.error("url:%s" % BUCKETS_URL.get("arlaunch") + ret['key'])



