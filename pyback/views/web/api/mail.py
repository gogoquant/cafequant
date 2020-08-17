# -*- coding: UTF-8 -*-
'''
    @brief publish mail to user
    @author: xiyanxiyan10
    @data 2019-01-15
'''

import logging
from services.mail import MailService
from views.web.base import AdminHandler
from iseecore.routes import route


@route(r"/v1/utils/mail", name="v1.utils.mail")
class MailHandler(AdminHandler):

    def _post_(self, *args, **kwargs):
        logging.info('call mail send')
        msg = self.get_argument('msg')
        mail_s = MailService()
        mail_s.add_msg(msg)
        self.write("success")
