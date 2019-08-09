#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
权限定义
每一个元素用一个二进制位控制。
"""

import logging
import simplejson
from services.defines import PRIV_DEFAULT, PRIV_MAX


class Permission(object):

    def _init_(self, pstr=PRIV_DEFAULT):
        self.msize = PRIV_MAX
        self.permissionInit(pstr)

    def permissionInit(self, pstr):
        '''设置权限'''
        self.size = len(pstr)
        '''过长的串重设长度'''
        if self.size > self.msize:
            self.size = self.msize

        self.pstr = pstr

    def permissionUnset(self, i):
        '''清空权限'''
        if i >= self.size:
            return False
        else:
            self.pstr[i] = '0'

    def permissionSet(self, i):
        '''设置权限'''
        if i >= self.size:
            return False
        else:
            self.pstr[i] = '1'

    def permissionCheckOne(self, i):
        '''检查权限'''
        if i >= self.msize:
            return False
        if self.pstr[i] == '0':
            return False
        return True

    def permissionSize(self):
        '''获取权限串长度'''
        return self.size

    def permissionMSize(self):
        '''获取最大可能权限长度'''
        return self.msize

    def permissionGet(self):
        '''获取权限串'''
        return self.pstr, self.size

    def permissionCheck(self, p):
        lftSize = self.permissionSize()
        rhtSize = p.permissionSize()

        if lftSize != rhtSize:
            return None

        for i in range(rhtSize):
            if p.permissionCheckOne(i):
                if not self.permissionCheckOne(i):
                    return False
        return True
