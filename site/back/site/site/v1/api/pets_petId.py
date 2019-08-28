# -*- coding: utf-8 -*-
from __future__ import absolute_import, print_function

from . import ApiHandler
from .. import schemas


class PetsPetid(ApiHandler):

    def get(self, petId):
        print("get")
        return [], 200, None