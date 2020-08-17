#! /usr/bin/env python
# coding=utf-8

from email.mime.text import MIMEText
from email.header import Header
from smtplib import SMTP
from services.base import BaseService, Singleton
import config.config as Config
import logging

# 邮件的正文内容
mail_content = '你好，这是Quant 通知邮件'
# 邮件标题
mail_title = 'Quant通知'


class Queue(object):
    """队列"""

    def __init__(self):
        self.items = []

    def empty(self):
        return self.items == []

    def put(self, item):
        """进队列"""
        self.items.insert(0, item)

    def get(self):
        """出队列"""
        return self.items.pop()

    def size(self):
        """返回大小"""
        return len(self.items)


__all__ = ['MailService']


@Singleton
class MailService(BaseService):

    def init(self, host, port, user, pwd, receivers):
        self.user = user
        self.pwd = pwd
        self.host = host
        self.port = port
        self.receivers = []
        for receiver in receivers:
            self.receivers.append(receiver)

        self.msgs = Queue()

    def add_msg(self, msg):
        logging.info("add msg: %s" % msg)
        config = Config.GetConfig()
        cache_size = int(config.get('smtp_cachesize'))
        if self.msgs.size() < cache_size:
            self.msgs.put(msg)

    def send_one(self):
        if self.msgs.empty() is True:
            return True

        try:
            smtp = SMTP('%s:%d' % (self.host, self.port), timeout=2)
            smtp.set_debuglevel(1)
            # smtp.connect(self.host, self.port)
            # smtp.ehlo()
            logging.info("connect server success")
            smtp.login(self.user, self.pwd)
            logging.info("login server success")

            data = self.msgs.get()
            msg = data
            logging.info("publish msg start: %s" % msg)

            msgInfo = MIMEText(msg, "plain", 'utf-8')
            msgInfo["Subject"] = Header(mail_content, 'utf-8')
            msgInfo["From"] = self.user
            msgInfo["To"] = ";".join(self.receivers)
            smtp.sendmail(self.user, self.receivers, msgInfo.as_string())
            logging.info("publish msg end")
            return True
        except Exception as e:
            # Put msg into data again if fails
            logging.info('smtp fail:' + self.host + ":" + str(self.port))
            logging.info(e)
            return False


if __name__ == '__main__':
    pass
