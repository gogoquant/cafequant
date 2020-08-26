#-*- coding: UTF-8 -*-
'''
    @brief start vpn
    @author: xiyanxiyan10
    @data 2020-01-15
'''
from services.vpn import VPNService
from iseecore.routes import route
from views.web.base import AdminHandler
import pdb


@route(r"/v1/utils/vpn", name="v1.utils.vpn")
class VPNHandler(AdminHandler):

    def post(self):
        vpn_s = VPNService()
        vpn_s.reverse()
        self.write("success")
