import random

from tornado import gen

from tornadoredis.exceptions import ResponseError

from redistest import RedisTestCase, async_test


class MiscTestCase(RedisTestCase):
    @async_test
    @gen.engine
    def test_response_error(self):
        res = yield gen.Task(self.client.set, 'foo', 'bar')
        self.assertTrue(res)
        res = yield gen.Task(self.client.llen, 'foo')
        self.assertIsInstance(res, ResponseError)
        self.stop()

    @async_test
    @gen.engine
    def test_for_memory_leaks(self):
        '''
        Tests if a Client instance destroyed properly
        '''
        def some_code(callback=None):
            c = self._new_client(on_destroy=callback)
            c.get('foo')

        for __ in xrange(1, 3):
            yield gen.Task(some_code)

        self.stop()

    @async_test
    @gen.engine
    def test_for_memory_leaks_gen(self):
        '''
        Find a way to destroy client instances created by
        tornado.gen-wrapped functions.
        '''
        @gen.engine
        def some_code(callback=None):
            c = self._new_client(on_destroy=callback)
            n = '%d' % random.randint(1, 1000)
            yield gen.Task(c.set, 'foo', n)
            n2 = yield gen.Task(c.get, 'foo')
            self.assertEqual(n, n2)

        yield gen.Task(some_code)

        self.stop()

    @async_test
    @gen.engine
    def test_with(self):
        '''
        Find a way to destroy client instances created by
        tornado.gen-wrapped functions.
        '''
        @gen.engine
        def some_code(callback=None):
            with self._new_client(on_destroy=callback) as c:
                n = '%d' % random.randint(1, 1000)
                yield gen.Task(c.set, 'foo', n)
                n2 = yield gen.Task(c.get, 'foo')
                self.assertEqual(n, n2)

        yield gen.Task(some_code)

        self.stop()

    @async_test
    @gen.engine
    def test_for(self):
        '''
        Find a way to destroy client instances created by
        tornado.gen-wrapped functions.
        '''
        @gen.engine
        def some_code(callback=None):
            c = self._new_client(on_destroy=callback)
            for n in xrange(1, 5):
                yield gen.Task(c.set, 'foo', n)
                n2 = yield gen.Task(c.get, 'foo')
                self.assertEqual('%d' % n, n2)

        yield gen.Task(some_code)

        self.stop()
