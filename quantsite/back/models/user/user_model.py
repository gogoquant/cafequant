#-*- coding: UTF-8 -*-
'''
Created on 2016-12-08

@author: xiyanxiyan10
@module: user_model
'''

import simplejson
import hashlib
import logging
import traceback

from iseecore.models import AsyncBaseModel


#用户数据
class User(AsyncBaseModel):
    user_id = True  #用户id
    group_id = True  #用户组id

    permission = True  #用户权限
    email = True  #用户邮箱
    passwd = True  #用户密码
    name = True  #用户昵称
    sex = True  #性别 0 male, 1 female , -1 unknown
    age = True  #用户年龄
    logo_md5 = True  #头像md5
    address = False  #家庭地址
    description = False  #用户简述

    social_uid = False  #第三方唯一标识
    source = False  #第三方登录来源

    social_logo = False  #小logo
    social_big_logo = False  #大logo
    social_full_logo = False  #全屏幕logo
    permission = False  #用户权限

    # meta
    table = 'user'
    key = 'user_id'


class ResetPasswd(AsyncBaseModel):
    reset_passwd_id = True
    # user_id         = True
    email = True

    #meta
    table = 'reset_passwd'


class Follow(AsyncBaseModel):
    follow_id = True
    user_id = True
    target_id = True

    table = "follow"
    key = "follow_id"


class RecommendFollow(AsyncBaseModel):
    recommend_follow_id = True
    user_id = True
    recommend_user_id = True
    recommend_type = False
    recommend_type_count = False
    followed = False

    table = "recommend_follow"
    key = "recommend_follow_id"
