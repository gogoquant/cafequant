# -*- coding: utf-8 -*-

###
### DO NOT CHANGE THIS FILE
### 
### The code is auto generated, your change will be overwritten by 
### code generating.
###
from __future__ import absolute_import

from .api.pets import Pets
from .api.pets_petId import PetsPetid
from .api.users import Users


url_prefix = 'v1'

routes = [
    dict(resource=Pets, urls=[r"/pets"], endpoint='pets'),
    dict(resource=PetsPetid, urls=[r"/pets/(?P<petId>[^/]+?)"], endpoint='pets_petId'),
    dict(resource=Users, urls=[r"/users"], endpoint='users'),
]

def load_uris(config):
    try:
        config.update_uri(routes, url_prefix)
    except:
        pass