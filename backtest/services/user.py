# -*- coding: UTF-8 -*-

'''
    @brief service for user
    @author: xiyan
    @date 2016-12-26
'''

import hashlib
import tornado.gen
import time
import tornado.httpclient
import random
import uuid
import pdb

from models.user.user_model import User, ResetPasswd

from services.base import BaseService 

from util.time_common import DAY_PATTERN, timestamp_to_string, FULL_PATTERN, week_seconds, day_seconds, month_seconds
from util.oauth2 import *

__all__ = ['EditorService', 'AdminService', 'SdkAdminService']

NoUserException = -1
PasswdErrorException = -2
NoAppException = -3
UserExistsException = -4
UserSameNameException = -5

PASSWD_SALT = "pass"
CHECK_TIMEOUT = -12
PASSWD_TIMEOUT = 60
MES_PAGE_COUNT = 20
ACTIVE_USER_TIME = day_seconds * 3


'''普通用户管理接口'''
class UserService(BaseService):
    
    user_m = User()
    reset_passwd_m = ResetPasswd()

    @tornado.gen.engine
    def get(self, user_id=None, name=None, email=None, callback=None):
        '''获取用户信息'''
        query = {}
        if user_id:
            query['user_id'] = user_id
        if name:
            query['name'] = name
        if email:
            query['email'] = email
        user = yield tornado.gen.Task(self.user_m.find_one, query)
        callback(user)

    @tornado.gen.engine
    def login(self, email, password, token, token_from, openid=None, app_id=None, callback=None):

        third_id = None
        
        if token:
            third_user = yield tornado.gen.Task(self._get_data_from_social, token, token_from, openid=openid)
            if third_user:
                third_id = third_user['social_uid']
            else:
                callback((None, None))
                return
        
        if third_id:
            query = {
                "social_uid": third_id,
                "source": token_from,
            }
        else:
            query = {
                'email': email
            }

        user = yield tornado.gen.Task(self.user_m.find_one, query)
        if not user:
            callback((None, NoUserException))
            return

        if not third_id:
            password = hashlib.md5("%s_%s" % (PASSWD_SALT, password)).hexdigest()
            logging.error("password: %s, db_passwd: %s" % (password, user["passwd"]))
            if user["passwd"] != password:
                callback((None, PasswdErrorException))
                return

        # update third user image and recommend follows
        if third_id:
            yield tornado.gen.Task(self._update_logo_if_have_not_custom, user)

        user_info = yield tornado.gen.Task(self.info, user, None, None)

        if third_id:
            user_info['third_name'] = third_user['name']
        else:
            user_info['third_name'] = ""

        callback((user['user_id'], user_info))

    @tornado.gen.engine
    def login_from_wechat_webview(self, code, callback=None):
        '''第三方登录'''
        wc_api = OAuth2.instance(WEBCHAT_T)
        if hasattr(setting, "proxy_host"):
            wc_api.set_proxy(setting.proxy_host, setting.proxy_port)
        token_info = yield tornado.gen.Task(wc_api.get_token_info, code)
        if not "access_token" in token_info:
            callback((None, None))
            return

        user_id, user_info = yield tornado.gen.Task(self.login, None, None, token_info["access_token"], WEBCHAT_T, token_info["openid"])
        if not user_id:
            user_id, user_info = yield tornado.gen.Task(self.register, None, None, token_info["access_token"], WEBCHAT_T, token_info["openid"], allow_random_name=1)

        callback((user_id, user_info))

    @tornado.gen.engine
    def login_from_weibo_webview(self, code, callback=None):
        '''第三方登录'''
        weibo_api = OAuth2.instance(WEIBO_T)
        if hasattr(setting, "proxy_host"):
            weibo_api.set_proxy(setting.proxy_host, setting.proxy_port)
        token_info = yield tornado.gen.Task(weibo_api.get_token_info_from_code, code)
        if not "access_token" in token_info:
            callback((None, None))
            return

        user_id, user_info = yield tornado.gen.Task(self.login, None, None, token_info["access_token"], WEIBO_T)
        if not user_id:
            user_id, user_info = yield tornado.gen.Task(self.register, None, None, token_info["access_token"], WEIBO_T, allow_random_name=1)

        callback((user_id, user_info))


    @tornado.gen.engine
    def register(self, email, passwd, token, token_from, openid=None, extra={}, app_id=None, allow_random_name=None, callback=None):
        third_id = None
        #第三方登录的支持
        if token:
            third_user = yield tornado.gen.Task(self._get_data_from_social, token, token_from, openid=openid)
            if third_user:
                third_id = third_user['social_uid']
            else:
                callback((None, None))
                return

        query = {"email": email}

        #第三方登录的支持
        if token:
            query = {"social_uid": third_id, "source": token_from}

        user = yield tornado.gen.Task(self.user_m.find_one, query)

        #pdb.set_trace()
        
        if user:
            callback((None, UserExistsException))
            return

        #第三方登录的支持
        if token:
            user = third_user.copy()
        else:
            try:
                name = email.split('@')[0]
            except Exception, e:
                name = ""
            user = {'name': name}
            user["email"] = email
            logging.error("passwd: %s" % passwd)
            passwd = hashlib.md5("%s_%s" % (PASSWD_SALT, passwd)).hexdigest()
            logging.error("md5_passwd: %s" % passwd)
            user["passwd"] = passwd

            for k, v in extra.items():
                if v:
                    user[k] = v

        if extra.get('name', None):
            user['name'] = extra['name']

        # check same name
        if user.get('name', ''):
            query = {
                'name': user['name']
            }
            same_name_user = yield tornado.gen.Task(self.user_m.find_one, query)
            if same_name_user:
                if allow_random_name:
                    newname = yield tornado.gen.Task(self._user_name_random_int, user['name'])
                    user['name'] = newname
                else:
                    data = {
                        'name': same_name_user['name'],
                        'photo': self._get_user_photo(user, size=USER_IMAGE_SMALL)
                    }
                    callback((data, UserSameNameException))
                    return

        new_user_id = yield tornado.gen.Task(self.user_m.insert, user, upsert=True, safe=True)
        ext_data = {
            'user_id': new_user_id
        }

        new_user = yield tornado.gen.Task(self.info, None, new_user_id, None)

        if third_id:
            new_user['third_name'] = third_user['name']
        else:
            new_user['third_name'] = ""

        callback((new_user_id, new_user))

    @tornado.gen.engine
    def _user_name_random_int(self, oldname, callback):
        token = random.randint(1, 9999)
        newname = oldname + "#" + str(token)
        checksame = False
        while not checksame:
            query = {
                'name': newname
            }
            same_name_user = yield tornado.gen.Task(self.user_m.find_one, query)
            if not same_name_user:
                checksame = True
            else:
                token = random.randint(1, 9999)
                newname = oldname + "#" + str(token)
        callback(newname)

    @tornado.gen.engine
    def search(self, page, pagesize, token=None, callback=None):
        "分页获取用户的信息, token 为正则匹配(可选)"
        query = {
            'pos': (page - 1) * pagesize,
            'count': pagesize,
            'conditions':{}
        }
        if not token is None:
            conditions = {
                'name': {"$regex": token}
            }
            query['conditions'] = conditions

        user_list = yield tornado.gen.Task(self.user_m.get_list, query)
        callback(user_list)

    @tornado.gen.engine
    def get_all_ids(self, callback=None):
        "获取所有用户的id"
        conditions = {}
        query = {
            'conditions': conditions,
            "fields": {"user_id": 1}
        }
        user_list = yield tornado.gen.Task(self.user_m.get_list, query)

        results = [u['user_id'] for u in user_list]

        callback(results)

    @tornado.gen.engine
    def delete(self, user_id, callback=None):
        "根据用户的名字删除"
        print user_id
        query = {"user_id": user_id}
        yield tornado.gen.Task(self.user_m.delete, query)
        callback(user_id)

    @tornado.gen.engine
    def count(self, token, callback):
        "查询用户总数"
        #@TODO
        query = {
        
        }
        c = yield tornado.gen.Task(self.user_m.count, query)
        callback(c)

    @tornado.gen.engine
    def reset_passwd(self, query, callback):
        data = query.copy()
        old_data = yield tornado.gen.Task(self.user_m.find_one, query)
        if not old_data:
            callback(400)

        yield tornado.gen.Task(self.reset_passwd_m.delete, query)

        reset_passwd_id = yield tornado.gen.Task(self.reset_passwd_m.insert, data)
        callback(reset_passwd_id)

    @tornado.gen.engine
    def get_reset_passwd_by_id(self, reset_passwd_id, callback):
        reset_passwd = yield tornado.gen.Task(self.reset_passwd_m.get_by_id, reset_passwd_id)
        callback(reset_passwd)

    @tornado.gen.engine
    def update_user_passwd(self, reset_passwd_id, email, passwd, callback):

        query = {
            "email": email
        }
        # user_password_md5 = hashlib.md5(user_password).hexdigest()
        # user_password_md5 = hashlib.md5(user_password_md5 + app_secret).hexdigest()
        
        new_passwd = hashlib.md5("%s_%s" % (PASSWD_SALT, passwd)).hexdigest()

        update_set = {
            "$set": {
                "passwd": new_passwd
            }
        }
        yield tornado.gen.Task(self.user_m.update, query, update_set)

        query = {
            'reset_passwd_id': reset_passwd_id
        }
        yield tornado.gen.Task(self.reset_passwd_m.delete, query)
        callback(None)

    def _parse_utf8_string(self, value):
        import chardet
        parse_value = chardet.detect(value)
        logging.error("parse_value:%s" % str(parse_value))
        if parse_value["encoding"] == "utf-8":
            return value
        else:
            return value.decode("gb18030", "ignore").encode("utf-8")

    @tornado.gen.engine
    def modify(self, user_id, update_values, callback=None):
        name = update_values.get('name', '')
        if name:
            query = {
                'name': name
            }
            exists_name_user = yield tornado.gen.Task(self.user_m.find_one, query)
            if exists_name_user:
                callback(UserExistsException)
                return

        query = {'user_id': user_id}
        update_set = {
            "$set": update_values,
        }
        yield tornado.gen.Task(self.user_m.update, query, update_set)

        new_info = yield tornado.gen.Task(self.info, None, user_id, None)
        callback(new_info)

    @tornado.gen.engine
    def info(self, finded_user, user_id, self_user_id, viewMore=False, callback=None):
        if not finded_user:
            user = yield tornado.gen.Task(self.user_m.get_by_id, user_id)
        else:
            user = finded_user
        if not user:
            callback({})
            return
        result = {
            'email': user.get('email', ''),
            'name': user.get('name', ''),
            'user_id': user['user_id'],
            'sex': user.get('sex', 0),
            'age': user.get('age', 0),
            'address': user.get('address', ''),
            'description': user.get('description', '')
        }
        callback(result)

    @tornado.gen.engine
    def _get_data_from_social(self, token, token_from, openid=None, callback=None):
        logging.error('from 3rd %s %s' % (token_from, token))
        if token_from == WEIBO_T:
            weibo_api = OAuth2.instance(WEIBO_T)
            if hasattr(setting, "proxy_host"):
                weibo_api.set_proxy(setting.proxy_host, setting.proxy_port)
            token_info = yield tornado.gen.Task(weibo_api.get_token_info, token)
            if not token_info:
                callback(None)
                return
            uid = token_info.get('uid', '')
            user_info = yield tornado.gen.Task(weibo_api.get_user_info, token, uid)
            if not user_info:
                callback(None)
                return
            if user_info['gender'] == 'm':
                sex = 0
            elif user_info['gender'] == 'f':
                sex = 1
            else:
                sex = -1

            user_data = {
                "social_uid": uid,
                "source": WEIBO_T,
                "name": user_info['name'],
                "social_logo": user_info['avatar_large'],
                "social_big_logo": user_info.get('avatar_hd', ''),
                "sex": sex,
                "address": user_info['location'],
                "description": user_info['description']
            }

            callback(user_data)
        elif token_from == QQ_T:
            qq_api = OAuth2.instance(QQ_T)
            if hasattr(setting, "proxy_host"):
                qq_api.set_proxy(setting.proxy_host, setting.proxy_port)
            token_info = yield tornado.gen.Task(qq_api.get_token_info, token)
            if not token_info:
                callback(None)
                return
            uid = token_info.get('openid', '')
            appkey = token_info.get('client_id', '')
            user_info = yield tornado.gen.Task(qq_api.get_user_info, token, appkey, uid)
            if user_info['ret'] < 0:
                callback(None)
                return
            if user_info['gender'] == u'男':
                sex = 0
            else:
                sex = 1
            user_data = {
                "social_uid": uid,
                "source": QQ_T,
                "name": user_info['nickname'],
                "social_logo": user_info.get('figureurl_qq_2', user_info['figureurl_qq_1']),
                "social_big_logo": self._get_big_qq_img(user_info.get('figureurl_qq_2', user_info['figureurl_qq_1'])),
                "sex": sex,
                "address": user_info['city'],
                "description": ''
            }

            callback(user_data)
        elif token_from == WEBCHAT_T:
            wc_api = OAuth2.instance(WEBCHAT_T)
            if hasattr(setting, "proxy_host"):
                wc_api.set_proxy(setting.proxy_host, setting.proxy_port)
            user_info = yield tornado.gen.Task(wc_api.get_user_info, token, openid)
            if (not user_info) or (user_info.get('errcode', '')):
                callback(None)
                return
            user_data = {
                "social_uid": user_info['unionid'],
                "source": WEBCHAT_T,
                "name": user_info['nickname'],
                "social_logo": self._get_small_wechat_img(user_info.get('headimgurl', '')),
                "social_big_logo": user_info.get('headimgurl', ''),
                "sex": 0 if user_info['sex'] == 1 else 1,
                "address": user_info['city'],
                "description": ''
            }

            callback(user_data)
        else:
            callback(None)

