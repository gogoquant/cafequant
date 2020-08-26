# -*- coding: UTF-8 -*-

import socket
import fcntl
import struct
import logging
import pdb
import tornado.ioloop
from tornado.options import options
#from stormed import Connection as StormedConnection

# from iseecore.services import init_service
from iseecore.consumers import consumer
from iseecore.models import SyncData
import setting


class Rabbitmq(object):
    """docstring for Rabbitmq"""
    conn = None
    ch = None

    # 是否准备好（是否已经初始化）
    is_ready = False

    def __init__(self, back=None):
        ''' 初始化ampq '''

        if not self.is_ready:
            self.back = back
            #self._replace_callback_queue_()
            #self.conn = StormedConnection(
            #    host=setting.MQ_HOST,
            #    username=setting.MQ_USERNAME,
            #    password=setting.MQ_PASSWORD,
            #    heartbeat=setting.MQ_HEARTBEAT
            #)

            #self.conn.connect(self._on_amqp_connection_)
            #self.conn.on_disconnect = self._on_amqp_disconnect_
            #self.is_ready = True

            # @TODO disable rabbitmq
            self._on_amqp_connection_()

    def add_consume(self, queue):
        '''追加回掉消费者'''
        self.ch.consume(queue, consumer.consume, no_ack=True)

    def _replace_callback_queue_(self):
        '''向配置句柄中再追加些queue相关的宏'''
        #import pdb
        #pdb.set_trace()
        cur_ip = self._get_ip_address_('eth0')
        logging.error('cur_ip:%s' % cur_ip)
        ip_array = cur_ip.split('.')

    def _on_amqp_connection_(self):
        '''初始化本地的rabbitmq'''

        def build():
            self.ch = self.conn.channel()

            #通知消息同步用的exchange
            self.ch.exchange_declare(
                exchange=options.Sync_Send_Data_Exchange,
                type="topic",
                durable=True)

            options["ch"] = self.ch
            SyncData.configure(self.ch)
            self._consume_queues_()
            logging.error('[init]Rabbitmq init success')

        #初始化rabbitmq queue
        #disable rabbitmq
        #build()

        #self.ch.on_error = build

        if self.back:
            self.back()

    def _consume_queues_(self):
        consume_queues = [
            #注册回掉队列
        ]
        for consume_queue in consume_queues:
            self.ch.consume(consume_queue, consumer.consume, no_ack=True)

    def _on_amqp_disconnect_(self):
        logging.error("[*]on_amqp_disconnect")
        logging.error("[*]will close tornado")
        try:
            tornado.ioloop.IOLoop.instance().stop()
        except AttributeError as e:
            logging.error("ioloop stop error! %s" % e.message)
            pass
        finally:
            exit(1)

    def _get_ip_address_(self, ifname):
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        return socket.inet_ntoa(
            fcntl.ioctl(s.fileno(), 0x8915, struct.pack('256s',
                                                        ifname[:15]))[20:24])
