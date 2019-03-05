#!/usr/bin/python
# -*- coding: utf-8 -*-

__author__ = 'anti-moth'

import inspect

ERROR_JSON_STATUS = -1
ERROR_HTTP_STATUS = 400

class YixunErrrorCodeDefine(object):
    YIXUN_STANDARD_ERROR = 0x0001
    WEATHER_INVALID_TIME_STAMP = 0x0002

code_map_msg = {
    YixunErrrorCodeDefine.YIXUN_STANDARD_ERROR: "standard error",
    YixunErrrorCodeDefine.WEATHER_INVALID_TIME_STAMP: "Invalid tt. can't conver to int",
}


class YixunStandardError(Exception):

    def __init__(self, error_code):
        self.error_code = error_code
        try:
            current_call = inspect.stack()[1]
            _iframe = current_call[0]
            self.line_no = _iframe.f_lineno
            self.module_name = _iframe.f_globals.get("__name__", "")
            self.method_name = current_call[3]
            self.class_name = _iframe.f_locals.get('self', None).__class__.__name__

        except (IndexError, AttributeError):
            self.line_no = ''
            self.module_name = ''
            self.method_name = ''
            self.class_name = ''

    def __repr__(self):
        msg = code_map_msg.get(self.error_code, '')
        return "[*] YixunError:%s > %s. module: %s, class: %s, method: %s, line: %s " % (self.error_code, msg,
                                                                                         self.module_name,
                                                                                         self.class_name,
                                                                                         self.method_name,
                                                                                         self.line_no)

    def __str__(self):
        return code_map_msg.get(self.error_code, 'not find any match msg for code:%s' % self.error_code)


class PlatformDefine(object):
    ANDROID = "android"
    IPHONE = "iphone"
    IOS = "iphone"
    BROWSER = "browser"
    UNKNOWN = "unknown_platform"

    UNKNOWN_ID = -1
    ANDROID_ID = 0
    IOS_ID = 1
    IPHONE_ID = 2
    BROWSER_ID = 3

class MisakaPriv(object):
    NONE = -1
    BROWSER = 0  #浏览权限
    TUCAO = 1    #吐槽权限
    BOOK  = 2    #书籍下载权限
    APP  = 3     #应用
    MAX   = 4
    
    def OrgPermisson(self):
        return "01234"
