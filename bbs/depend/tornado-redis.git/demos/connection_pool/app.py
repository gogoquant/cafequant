import logging

import tornado.httpserver
import tornado.web
import tornado.ioloop
import tornado.gen

import tornadoredis


logging.basicConfig(level=logging.DEBUG)

log = logging.getLogger('app')

# Create a global redis connection pool of
# ten Redis server connections to use with this demo.
# Let redis clients wait for available connection instead of
# raising an exception if none is available.
CONNECTION_POOL = tornadoredis.ConnectionPool(max_connections=10,
                                              wait_for_available=True)

NUMBER_OF_CLIENTS = 5


class MainHandler(tornado.web.RequestHandler):

    @tornado.gen.engine
    def incr_counter(self, client, multiplier, callback=None):
        k = 'counter%d' % multiplier
        res = yield tornado.gen.Task(client.incrby, k, multiplier)
        # Ensure that demo keys will expire in 2 minutes.
        # Note that to execute this command a client has to wait
        # for other clients to complete their INCRBY calls.
        # You may check it using the redis-cli command-line utility
        # and the MONITOR command.
        yield tornado.gen.Task(client.expire, k, 120)
        # Return the number of visits multiplied by specified value
        callback(res)

    @tornado.web.asynchronous
    @tornado.gen.engine
    def get(self):
        '''
        Register the number of page views and return it on the page.

        Create NUMBER_OF_CLIENTS redis clients and connect them
        to the global connection pool.
        '''
        indexes = range(1, NUMBER_OF_CLIENTS + 1)
        clients = (tornadoredis.Client(connection_pool=CONNECTION_POOL)
                   for __ in indexes)

        # Register page views using %NUMBER_OF_CLIENTS% redis clients.
        # Run redis commands "simultaneously" to simulate a complex
        # application creating a number of Client instances.
        values = yield map(lambda c, n: tornado.gen.Task(self.incr_counter,
                                                         c, n),
                           clients,
                           indexes)

        # Create a new client and get the redis server information
        # to display it on a demo page.
        c = tornadoredis.Client(connection_pool=CONNECTION_POOL)
        info = yield tornado.gen.Task(c.info)
        values = map(lambda n, v: (n, v), indexes, values)
        self.render("template.html",
                    title="Connection Pool Demo",
                    values=values,
                    info=info)


application = tornado.web.Application([
    (r'/', MainHandler),
])


if __name__ == '__main__':
    http_server = tornado.httpserver.HTTPServer(application)
    http_server.listen(8888)
    print 'Demo is runing at 0.0.0.0:8888\nQuit the demo with CONTROL-C'
    tornado.ioloop.IOLoop.instance().start()
