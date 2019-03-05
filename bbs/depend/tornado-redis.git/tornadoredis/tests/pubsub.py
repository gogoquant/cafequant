from functools import partial
from tornado import gen

from redistest import RedisTestCase, async_test


class PubSubTestCase(RedisTestCase):

    def setUp(self):
        super(PubSubTestCase, self).setUp()
        self._message_count = 0
        self.publisher = self._new_client()

    def tearDown(self):
        try:
            self.publisher.connection.disconnect()
            del self.publisher
        except AttributeError:
            pass
        super(PubSubTestCase, self).tearDown()

    def pause(self, timeout, callback=None):
        self.delayed(timeout, callback)

    def _expect_messages(self, messages, expected_number,
                         subscribe_callback=None, callback=None):
        self._expected_messages = messages
        self._expected_number = expected_number
        self._subscribe_callback = subscribe_callback
        self._done_callback = callback

    def _handle_message(self, msg):
        self._message_count += 1
        self.assertIn(msg.kind, self._expected_messages)
        expected = self._expected_messages[msg.kind]
        self.assertEqual(msg.pattern, expected[0])
        self.assertEqual(msg.body, expected[1])
        if msg.kind in ('subscribe', 'psubscribe'):
            if self._subscribe_callback:
                cb = self._subscribe_callback
                self._subscribe_callback = None
                cb(True)
        if self._message_count >= self._expected_number:
            if self._done_callback:
                self._done_callback(True)
            self.stop()

    @async_test
    @gen.engine
    def test_pub_sub(self):
        self._expect_messages({'subscribe': ('foo', 1),
                               'message': ('foo', 'bar'),
                               'unsubscribe': ('foo', 0)},
                              3,
                              subscribe_callback=(yield gen.Callback('sub')),
                              callback=(yield gen.Callback('done')))
        yield gen.Task(self.client.subscribe, 'foo')
        self.client.listen(self._handle_message)
        yield gen.Wait('sub')
        yield gen.Task(self.publisher.publish, 'foo', 'bar')
        yield gen.Task(self.publisher.publish, 'foo', 'bar')
        yield gen.Task(self.publisher.publish, 'foo', 'bar')
        yield gen.Task(self.client.unsubscribe, 'foo')

        yield gen.Wait('done')

        self.assertEqual(self._message_count, 3)
        self.stop()

    @async_test
    @gen.engine
    def test_pub_psub(self):
        self._expect_messages({'psubscribe': ('foo.*', 1),
                               'pmessage': ('foo.*', 'bar'),
                               'punsubscribe': ('foo.*', 0),
                               'unsubscribe': ('foo.*', 1)},
                              2,
                              subscribe_callback=(yield gen.Callback('sub')),
                              callback=(yield gen.Callback('done')))
        yield gen.Task(self.client.psubscribe, 'foo.*')
        self.client.listen(self._handle_message)
        yield gen.Wait('sub')
        yield gen.Task(self.publisher.publish, 'foo.1', 'bar')
        yield gen.Task(self.publisher.publish, 'bar.1', 'zar')
        yield gen.Task(self.client.punsubscribe, 'foo.*')

        yield gen.Wait('done')

        self.assertEqual(self._message_count, 2)
        self.stop()
