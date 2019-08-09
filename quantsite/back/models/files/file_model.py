#-*- coding: UTF-8 -*-
'''
    @brief 文件存储相关模板
    @author lancelot
    @date 2017/1/5
'''

import logging
import simplejson

from iseecore.models import AsyncBaseModel


class File(AsyncBaseModel):

    md5 = True  #文件MD5
    meta = False  #文件属性
    store = False  #存储容器

    width = False  #图片宽度，可选
    height = False  #图片高度, 可选

    file_name = False  #文件名字
    file_size = False  #文件大小

    qn_md5 = False  #七牛md5
    qn_key = False  #七牛key

    # meta
    table = 'file'
    key = 'file_id'
