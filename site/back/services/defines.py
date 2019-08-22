#!/usr/bin/python
# -*- coding: utf-8 -*-

import inspect

ERROR_JSON_STATUS = -1
ERROR_HTTP_STATUS = 400


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


PRIV_DEFAULT = '10000'

PRIV_GET = 0
PRIV_UPDATE = 1
PRIV_CREATE = 2
PRIV_DELETE = 3
PRIV_ADMIN = 4
PRIV_MAX = 5
