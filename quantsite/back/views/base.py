# -*- coding: utf-8 -*-
'''
    @brief handler base
    @author: xiyanxiyan10
    @note  base class for web handler
'''

import os
import logging
import tornado.web
import simplejson

from tornado.options import define, options
from tornado import locale
from pycket.session import SessionMixin
from iseecore.routes import route

import setting

class RequestHandler(tornado.web.RequestHandler, SessionMixin):
    def __init__(self):
        self.initialize()
        
    def tileSet(self, title):
        self.title = title
        
    def ttemplateSet(self, template):
        self.template = template

    def initialize(self):
        self.title = ""
        self.template = ""
        self.response = {}
        self.ch = options.group_dict('ch')
        self.fastdfs_client = options.group_dict('fastdfs_client')
        self.redis_client = options.group_dict('redis_client')

        if not self.ch:
            logging.error('need channel in tornado.options')
        
        if not self.fastdfs_client:
            logging.error('need fastdfs_client in tornado.options')

        if not self.redis_client:
            logging.error('need redis_client in tornado.options')

    def prepare(self):
        self.view_permission = "00000"
        #self.edit_permission = None
        

    def _get_(self, *args, **kwargs):
        raise tornado.web.HTTPError(405)

    def _post_(self, *args, **kwargs):
        raise tornado.web.HTTPError(405)  
    
    def _put_(self, *args, **kwargs):
        raise tornado.web.HTTPError(405) 
    
    def _delete_(self, *args, **kwargs):
        raise tornado.web.HTTPError(405)          

    def get_user_locale(self):
        if hasattr(self, "lan_form_arg"):
            lan = ''#self.lan_form_arg
        else:
            lan = self.get_secure_cookie("user_locale")
        self.__locale = lan
        
        print lan
        if not lan:
            bl = self.get_browser_locale()
            bl_code = bl.code

            print bl_code
            self.__locale = bl_code
            self.set_secure_cookie("user_locale", bl_code)
            return bl
        return locale.get(lan)
    
    def init_locale(self):
         lan = self.get_secure_cookie("user_locale")
         if not lan:
             lan = 'zh_CN'
             self.set_secure_cookie("user_locale", lan)
         self.__locale = lan

    def get_lan(self):
        return self.__locale

    def _(self, text):
        ''' Localisation shortcut '''
        return self.locale.translate(text).encode("utf-8")

    def _u(self, text):
        return self.locale.translate(text)

    def render(self, **kwargs):
        
        title = self.title
        template = self.template
        
        # add public env 
        kwargs['site_url'] = setting.SITE_URL
        kwargs['isset'] = self.isset
        kwargs['make_url'] = self.make_url
        kwargs['get_comma'] = self.get_comma

        self.namespace = kwargs
        tornado.web.RequestHandler.render(self, template, title=title, lan=self.__locale, **kwargs)

    def isset(self, v):
        return self.namespace.has_key(v)

    def make_url(self, url, **kwargs):
        return url

    def get_comma(self, code, **kwargs):
        lan = self.get_secure_cookie("user_locale")
        if code == 'comma':
            if lan == 'zh_CN':
                return '，'
            elif lan == 'ja_JP':
                return '、'
            elif lan == 'en_US':
                return ', '
        elif code == 'stop':
            if lan == 'zh_CN':
                return '、'
            elif lan == 'ja_JP':
                return '、'
            elif lan == 'en_US':
                return ', '
        return ''

    @property
    def app_name(self):
        return None

    @property
    def base_file_path(self):
        return self.application.base_file_path

    def get_browser_locale(self, default="en_US"):
        """Determines the user's locale from Accept-Language header.

        See http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.4
        """
        if "Accept-Language" in self.request.headers:
            languages = self.request.headers["Accept-Language"].split(",")
            locales = []
            for language in languages:
                parts = language.strip().split(";")
                if len(parts) > 1 and parts[1].startswith("q="):
                    try:
                        score = float(parts[1][2:])
                    except (ValueError, TypeError):
                        score = 0.0
                else:
                    score = 1.0
                locales.append((parts[0], score))
            if locales:
                locales.sort(key=lambda (l, s): s, reverse=True)
                codes = [self._true_code(l[0]) for l in locales]
                return locale.get(*codes)
        return locale.get(default)

    def _true_code(self, code):
        if code == 'zh':
            return 'zh_CN'
        elif code == 'ja':
            return 'ja_JP'
        elif code == 'en':
            return 'en_US'
        else:
            return code

@route(r"/(.*)", name="error")
class ErrorHandler(RequestHandler):
    def prepare(self):
        super().prepare()
        self.set_status(404)
        raise tornado.web.HTTPError(404)
