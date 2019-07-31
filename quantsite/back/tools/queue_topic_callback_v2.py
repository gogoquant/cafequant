from tornado import ioloop
from tornado.options import options
from stormed import Connection
from pymongo import Connection as Conn
import pymongo
import os
import logging, time

import setting

base_file_path = os.path.dirname(__file__)
if base_file_path == '':
    base_file_path = './'
os.chdir(base_file_path)

# usage: {'index': [0 ... 8], 'count': index total count, 'result_map': {'0': [keypoint, rate_value, similarity] ... '8':}}
recs_dic = {}
# rate value for subimage
index_map = {'0':1, '1':1.5, '2':1, '3':1.5, '4':2, '5':1.5, '6':1, '7':1.5, '8':1}

TRAINING_STATE_SUCCESS = 1


class TopicCallback:
    def topic_callback_request(self,msg):
        global ch
        global db_conn
        global librecs
        state = msg.headers["state"]
        logging.error(msg.headers)
        groupname = "group1"
        md5 = None
        if msg.headers.has_key("GroupName"):
            groupname = msg.headers["GroupName"]
        if msg.headers.has_key("ImageName"):
            md5 = msg.headers["ImageName"]

        logging.error("msg : %s" %msg)
        remote_fileid = msg.headers["FileId"]
        #skip blur file
        if remote_fileid.find("_blur") > 0:
            logging.error("skip blur file")
            msg.ack()
            return
        keynumber = msg.headers["KeyPointNumber"]
        if state == "204 FAILED":
            """ logging group failed, training failed """
            db = db_conn.training
            db.groups.update({"groupname":groupname}, {"$inc":{"training_fail_count":1}})
            db.tlog.insert({"time_key":int(time.time()), "groupname":groupname, \
                                "imagename":remote_fileid, "state":state})
            msg.ack()
        else:
            """ logging succeed,update files files """
            """ in version2.0, do not need to compute similarity """
            db = db_conn.files
            update_set = {}
            update_set["$set"] = {"key_point": keynumber}
            db.files.update({"md5": md5}, update_set)

            self.keynumber = keynumber
            self.origin_md5 = md5
            self.update_res_picture(1)
            msg.ack()

    def update_res_picture(self, similarity):
        global db_conn
        if not hasattr(self, 'keynumber'):
            self.keynumber = 0
        db = db_conn.resource
        query = {
            'md5': self.origin_md5,
            'delete_flag': {'$ne': '1'}
        }
        res_picture_query = db.res_picture.find_one(query)
        if not res_picture_query:
            logging.error('has res_picture | md5:%s' % self.origin_md5)
            return
        update_set = {
            '$set': {
                'training_state': TRAINING_STATE_SUCCESS,
                'similarity': similarity,
                'key_point': self.keynumber
            }
        }
        query = {
            'md5': self.origin_md5,
            'delete_flag': {'$ne': '1'}
        }
        db.res_picture.update(query, update_set, multi=True)


def topic_finish_request(msg):
    topic_callback = TopicCallback()
    topic_callback.topic_callback_request(msg)


def on_amqp_connection():
    global ch
    ch = conn.channel()
    # ch.exchange_declare(exchange=options.Search_CallBack_Exchange,
    #                     type="topic", durable=True)
    ch.queue_declare(queue=options.Topic_CallBack_Queue, durable=True)
    ch.queue_bind(queue=options.Topic_CallBack_Queue,
                  exchange=options.Search_CallBack_Exchange,
                  routing_key=options.Topic_CallBack_Queue)
    ch.consume(options.Topic_CallBack_Queue, topic_finish_request, no_ack=False)
    logging.error('consume:[E %s][Q %s][R %s]' %
                  (options.Search_CallBack_Exchange,
                   options.Topic_CallBack_Queue,
                   options.Topic_CallBack_Queue))

    base_file_path = os.path.dirname(__file__)


def on_amqp_disconnect():
    logging.error("[*]on_amqp_disconnect")
    logging.error("[*]will close tornado")
    try:
        ioloop.IOLoop.instance().stop()
    except AttributeError, e:
        logging.error(e.message)
        pass
    finally:
        exit(1)
ch = None
options.log_to_stderr = True
options.parse_command_line()


mongo_host = 'mongodb://%s:%s@%s' % (setting.MONGO_USER,
                                     setting.MONGO_PASS,
                                     setting.MONGO_HOST)
mongo_port = setting.MONGO_PORT
mq_host = setting.MQ_HOST
mq_heartbeat = setting.MQ_HEARTBEAT
mq_username = setting.MQ_USERNAME
mq_password = setting.MQ_PASSWORD

db_conn = Conn(mongo_host, mongo_port)
conn = Connection(host=mq_host, username=mq_username,
                  password=mq_password, heartbeat=mq_heartbeat)

base_file_path = os.path.dirname(__file__)

conn.connect(on_amqp_connection)
conn.on_disconnect = on_amqp_disconnect

try:
    ioloop.IOLoop.instance().start()
except KeyboardInterrupt:
    conn.close()
