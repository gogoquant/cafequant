#!/usr/bin/env python
import time

from tornado.testing import AsyncTestCase

import tornadoredis


def get_callable(obj):
    return hasattr(obj, '__call__')


def async_test_ex(timeout=5):
    def _inner(func):
        def _runner(self, *args, **kwargs):
            try:
                func(self, *args, **kwargs)
            except:
                self.stop()
                raise
            return self.wait(timeout=timeout)
        return _runner
    return _inner


def async_test(func):
    _inner = async_test_ex()
    return _inner(func)


class TestRedisClient(tornadoredis.Client):

    def __init__(self, *args, **kwargs):
        self._on_destroy = kwargs.get('on_destroy', None)
        if 'on_destroy' in kwargs:
            del kwargs['on_destroy']
        super(TestRedisClient, self).__init__(*args, **kwargs)

    def __del__(self):
        super(TestRedisClient, self).__del__()
        if self._on_destroy:
            self._on_destroy()


class RedisTestCase(AsyncTestCase):
    test_db = 9
    test_port = 6379

    def setUp(self):
        super(RedisTestCase, self).setUp()
        self.client = self._new_client()
        self.client.flushdb()

    def tearDown(self):
        try:
            self.client.connection.disconnect()
            del self.client
        except AttributeError:
            pass
        super(RedisTestCase, self).tearDown()

    def _new_client(self, pool=None, on_destroy=None):
        client = TestRedisClient(io_loop=self.io_loop,
                                 port=self.test_port,
                                 selected_db=self.test_db,
                                 connection_pool=pool,
                                 on_destroy=on_destroy)

        return client

    def delayed(self, timeout, cb):
        self.io_loop.add_timeout(time.time() + timeout, cb)
