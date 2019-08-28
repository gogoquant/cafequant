# -*- coding: UTF-8 -*-
'''
Created on 2017-3-11
@author: lancelot

'''

# -*- coding: utf-8 -*-
import os.path
from tornado.options import define, options
import tornado.httpserver
import tornado.ioloop
import logging
import tornado.web
import pdb

from ctypes import cdll
from iseecore.routes import route
from iseecore.services import init_service
import setting
import pdb 
#from core import load_tornado_settings

#modules = ['v1']
#config = load_tornado_settings(*modules)

#from tornadomail.backends.smtp import EmailBackend

define("port", default=8080, help="Run server on a specific port", type=int)


# 自定义Application类，继承于Application
class Application(tornado.web.Application):

    def __init__(self):
    #def __init__(self, url_list, **app_settings):
        #tornado.web.Application.__init__(self, url_list, **app_settings)
        #self.config = config

        def init():
            init_service.load_init_services()

            logging.info("register all route")
            import views.web.api.user
            logging.info("register route->%s", views.web.api.user.__package__)
            #import views.web.api.topic
            #import views.web.api.tag
            #import views.web.api.chat
            #import views.web.api.pubmsg
            #import views.web.api.upload

            handlers = route.get_routes()

            # 定义setting，对tornado.web.Application进行设置
            settings = dict(
                blog_title=u"quant",

                #设置模板文件路径
                #template_path=os.path.join(
                #   os.path.dirname(__file__), setting.MISAKA_TEMPLATE_DIR),

                #设置static文件路径, swaggar 使用
                static_path = os.path.join(os.path.dirname(__file__), "site/site/static"),
                
                debug=True,
                cookie_secret="NyM8vWu0Slec0GLonsdI3ebBX0VEqU3Vm8MsHiN5rrc=",
                app_secret="XOJOwYhNTgOTJqnrszG3hWuAsTmVz0GisOCY6R5d1E8=",
                login_url="/",
                autoescape=None,
                gzip=True,
            )

            # redis配置
            redis_setting = {
                'host': setting.REDIS_HOST,
                'port': setting.REDIS_PORT,
                'db_sessions': setting.REDIS_SESSIONS,
            }

            # session配置为redis存储
            settings['session'] = {'engine': 'redis', 'storage': redis_setting}

            url_list = []
            #url_list.extend(config.URIS)
            url_list.extend(handlers)
            
            tornado.web.Application.__init__(self, url_list, **settings)

        # init amqp
        from services.tool.amqp import Rabbitmq

        logging.info("Try to init rabbitmq")
        Rabbitmq(back=init)

    @property
    def mail_connection(self):
        return True

    @property
    def base_file_path(self):
        return os.path.dirname(__file__)


def main():
    options.log_to_stderr = True
    options.logging = 'info'
    tornado.options.parse_command_line()
    # 启动http server监听端口
    http_server = tornado.httpserver.HTTPServer(Application())
    http_server.listen(options.port)
    # 设置本地化
    tornado.locale.set_default_locale("zh_CN")
    tornado.locale.load_translations(
        os.path.join(os.path.dirname(__file__), "translations"))
    # 启动ioloop
    tornado.ioloop.IOLoop.instance().start()


if __name__ == '__main__':
    main()
