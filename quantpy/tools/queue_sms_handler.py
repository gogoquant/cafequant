#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
发送短信

"""
import logging

import tornado.gen
from tornado.options import options

from tornado import ioloop
import simplejson

from stormed import Message, Connection
import setting
from services.tool.db import Mongodb
from iseecore.models import SyncData


# 一些变量
Sms_Exchange = "_yixin_sms_send"
Sms_Queue = "_yixin_sms_send"


class SMSHandler(object):
    def __init__(self):
        from services.sms import SmsService
        self.sms_s = SmsService()

    @tornado.gen.engine
    def init(self, callback):
        callback(None)

    @tornado.gen.engine
    def run(self, msg):
        method = msg.headers["Method"]
        content = simplejson.loads(msg.body)
        # try:
        logging.error("receive %s" % method)
        if method == "pull_status":
            yield tornado.gen.Task(self.sms_s.pull_sms_status)
        else:
            yield tornado.gen.Task(self.sms_s.send_veri_code, content, method)
        # except Exception, e:
        #     logging.error(e.message)
        #     raise e
        # finally:
        msg.ack()


def on_amqp_connection():
    global ch
    ch = conn.channel()
    options["ch"] = ch
    SyncData.configure(ch)
    ch.exchange_declare(exchange=Sms_Exchange, type="direct", durable=True)
    ch.queue_declare(queue=Sms_Queue, durable=True)
    ch.queue_bind(
        queue=Sms_Queue,
        exchange=Sms_Exchange,
        routing_key=Sms_Queue
    )
    Mongodb()
    vh = SMSHandler()

    def callback(p):
        logging.info('bind success!')

    vh.init(callback)
    ch.consume(Sms_Queue, vh.run, no_ack=False)

    # TEST send sms
    request = Message(
        simplejson.dumps({
            "code": "9999",
            "mobile": "18980098109",
            "app_id": "5268b93b602ede779f92501e"
        }),
        delivery_mode=2,
        headers={
            "Method": "veri_code"
        }
    )
    # ch.publish(
    #     request, 
    #     exchange    = Sms_Exchange, 
    #     routing_key = Sms_Queue
    # )


def on_amqp_disconnect():
    logging.error("[*]on_amqp_disconnect")
    logging.error("[*]will close tornado")
    try:
        tornado.ioloop.IOLoop.instance.stop()
    except AttributeError, e:
        pass
    finally:
        exit(1)


if __name__ == '__main__':
    logging.basicConfig(format='[%(asctime)s %(filename)s:%(lineno)d %(levelname)s] %(message)s', level=logging.DEBUG)

    mq_host = setting.MQ_HOST
    mq_heartbeat = setting.MQ_HEARTBEAT
    mq_username = setting.MQ_USERNAME
    mq_password = setting.MQ_PASSWORD

    conn = Connection(host=mq_host, username=mq_username, password=mq_password, heartbeat=mq_heartbeat)
    conn.connect(on_amqp_connection)
    conn.on_disconnect = on_amqp_disconnect

    try:
        ioloop.IOLoop.instance().start()
    except KeyboardInterrupt:
        conn.close()
