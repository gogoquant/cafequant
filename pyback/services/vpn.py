import time
import pdb
import os

from services.base import BaseService

__all__ = ['VPNService']


class VPNService(BaseService):

    def start(self):
        res = self.get_vpn()
        if res == "close":
            print("close->open")
            os.system(
                "/usr/local/bin/ssserver -c /etc/shadowsocks-python/config.json -d start --pid-file=/root/shadow.pid"
            )

    def reverse(self):
        res = self.get_vpn()
        if res == "open":
            self.stop()
        else:
            self.start()

    def stop(self):
        res = self.get_vpn()
        if res == "open":
            print("open->close")
            pid = self.get_vid()
            print(pid)
            os.system("kill " + pid)

    def get_vid(self):
        fp = open("/root/shadow.pid")
        pid = fp.read()
        return pid

    def get_vpn(self):
        res = os.popen(
            "ps -aux|grep /etc/shadowsocks-python/config.json|wc -l").read()
        if (res == '3\n'):
            res = "open"
        else:
            res = "close"
        return res
