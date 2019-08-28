# -*- coding: utf-8 -*-
from __future__ import absolute_import, print_function

from . import ApiHandler
from .. import schemas


class Pets(ApiHandler):

    def get(self):
        print(self.args)

        return [], 200, {}

    def post(self):

        return None, 201, None