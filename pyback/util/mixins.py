#!/usr/bin/env python
# coding=utf8
try:
    import psyco
    psyco.full()
except:
    pass
import traceback
import sys
import os
import httplib
import tornado
import logging

from models.ia.analytics_model import Analytics, Imei
import tornado.gen


class LogMixin(object):

    @tornado.gen.engine
    def write_log(self):
        # if not self.get_argument("behaviors", None):
        imei = ""
        uris = self.request.uri.split('?imei=')

        imei = self.get_argument("imei", "")
        if not imei:
            imei = self.request.headers.get("Imei", "")

        if not imei:
            logging.error("no imei")
            return

        analytics_data = {
            "class_name": self.__class__.__name__,
            "uri": uris[0],
            "user_id": self.get_secure_cookie("user_id"),
            "method": self.request.method,
            "remote_ip": self.request.remote_ip,
            "device": imei
        }
        # save_keys = ["X-Real-Ip", "referer", "device"]

        save_keys = ["X-Real-Ip", "referer"]

        for key in save_keys:
            analytics_data[key] = self.request.headers.get(key, "")

        if imei:

            analytics_class = Analytics()
            yield tornado.gen.Task(analytics_class.insert, analytics_data)

            if (self.request.uri.find("/resource/init") >= 0):
                user_data = {}
                save_keys = [
                    "software", "imei", "system", "market", "device", "os",
                    "network", "appkey"
                ]
                for key in save_keys:
                    user_data[key] = self.get_argument(key, "")

                imei_class = Imei()
                old_data = yield tornado.gen.Task(imei_class.find_one,
                                                  {"imei": user_data["imei"]})
                if old_data:
                    system = user_data.get("system", "")
                    if system != "":
                        yield tornado.gen.Task(imei_class.update,
                                               {"imei": user_data["imei"]},
                                               {"$set": user_data})
                    else:
                        yield tornado.gen.Task(imei_class.update,
                                               {"imei": user_data["imei"]},
                                               {"$set": {}})
                else:
                    yield tornado.gen.Task(imei_class.insert, user_data)
