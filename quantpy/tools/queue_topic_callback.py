from tornado import ioloop
from tornado.options import define,options
from stormed import Message, Connection
from pymongo import Connection as Conn
import pymongo
import os
import logging, time
import ConfigParser
import simplejson
import uuid
import os.path
from ctypes import cdll
from ctypes import *

from util.fastdfs_utils import DfsUtils
from util.consumers import consumer
import setting

base_file_path =  os.path.dirname(__file__)
if base_file_path == '':
    base_file_path = './'
os.chdir(base_file_path)

#usage: {'index': [0 ... 8], 'count': index total count, 'result_map': {'0': [keypoint, rate_value, similarity] ... '8':}}
recs_dic = {}
#rate value for subimage
index_map = {'0':1, '1':1.5, '2':1, '3':1.5, '4':2, '5':1.5, '6':1, '7':1.5, '8':1}

TRAINING_STATE_SUCCESS = 1
class TopicCallback:
    def topic_callback_request(self,msg):
        global ch
        global db_conn
        global librecs
        state     = msg.headers["state"]
        logging.error(msg.headers)
        groupname = "group1"
        md5       = None
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
            db = db_conn.files
            update_set = {}
            update_set["$set"] = {"key_point": keynumber}
            db.files.update({"md5": md5}, update_set)
            f = db.files.find_one({'md5':md5})
            file_id = f.get('file_id')
            file_id_splits = file_id.split('.')
            file_suffix = file_id_splits[-1]
            if file_suffix == 'png':
                file_suffix = 'jpg'
            file_id = file_id_splits[0] + '_1024.' + file_suffix
            file_temp = "/tmp/%s."%uuid.uuid4().get_hex() + '.' + file_suffix
            fastdfs_client.download_to_file(file_temp,str(file_id))

            fun_recs = librecs.reload
            fun_recs.argtypes = [c_char_p, c_char_p]

            new_file_temp_path = "/tmp/%s."%uuid.uuid4().get_hex()
            fun_recs(file_temp , new_file_temp_path)

            try:
                os.remove(file_temp)
            except Exception, e:
                pass
            if not os.path.exists(new_file_temp_path + "0.jpg"):
                msg.ack()
                return

            self.keynumber = keynumber
            self.origin_md5 = md5
            req_id = str(id(self))
            recs_dic[req_id] = {"index": [], "count": 0, "result_map": {}}
            for i in xrange(9):
                new_file_temp = "%s%d.jpg" % (new_file_temp_path, i)
                if os.path.exists(new_file_temp):
                    new_req_id = "%s.%s" % (req_id, str(i))
                    recs_dic[req_id]["index"].append(i)
                    recs_dic[req_id]["count"] = recs_dic[req_id]["count"] + 1
                    consumer.add_method(new_req_id, self.search_callback_request)
                    with open(new_file_temp, "rb") as file_data:
                        binary_content = file_data.read()
                    search_msg = Message(binary_content, delivery_mode=2, reply_to=options.Search_CallBack_Queue, headers={"Method":    "Search",
                                                                                                                    "GroupName":str(groupname), "req_id": new_req_id,
                                                                                                                    "hash":str(int(req_id, 16))})
                    ch.publish(search_msg, exchange="amq.direct", routing_key = "sift_extract")
                    try:
                        os.remove(new_file_temp)
                    except Exception, e:
                        pass
                else:
                    continue

            if recs_dic[req_id]["count"] == 0:
                recs_dic.pop(req_id)

            db = db_conn.training
            db.tlog.insert({"time_key":int(time.time()), "groupname":groupname, \
                                "imagename":remote_fileid, "state":state, \
                                "knumber":keynumber, "sdtime":msg.headers["SiftDownloadTime"], \
                                "stime":msg.headers["SiftTime"], "tdtime":msg.headers["IndexDownloadTime"], \
                                "ttime":msg.headers["IndexTime"], "utime":msg.headers["UploadTime"]})
            msg.ack()

    def search_callback_request(self,msg):
        data = str(msg.body).replace("\\", "")
        data = simplejson.loads(data)
        #logging.error(str(data))
        headers = msg.headers
        state = data["state"]
        key_num = 0
        if "KeyPointNumber" in headers:
            key_num = headers["KeyPointNumber"]

        if not hasattr(self,'origin_md5'):
            logging.error('self has no topic callback')
            msg.ack()
            return

        req_id_and_index = data["req_id"].split('.')
        cur_req_id = req_id_and_index[0]
        cur_index = req_id_and_index[1]
        if not recs_dic.has_key(cur_req_id):
            logging.error('recs dict has no cur request callback')
            msg.ack()
            return
        if state == "204 FAILED":
            recs_dic[cur_req_id]["result_map"][cur_index] = [key_num, index_map[cur_index], 0]
        else:
            image_path = data["image_path"]
            if image_path == "204 FAILED":
                recs_dic[cur_req_id]["result_map"][cur_index] = [key_num, index_map[cur_index], 0]
            else:
                image_path = image_path.split('}')[0]
                image_path_list = image_path.split('_')
                try:
                    file_md5 = image_path_list[0] + '_' + image_path_list[1]
                    if file_md5 == self.origin_md5:
                        recs_dic[cur_req_id]["result_map"][cur_index] = [key_num, index_map[cur_index], image_path_list[2]]
                    else:
                        recs_dic[cur_req_id]["result_map"][cur_index] = [key_num, index_map[cur_index], -1]
                except IndexError:
                    logging.error("parse imageData error IndexError")
                    recs_dic[cur_req_id]["result_map"][cur_index] = [key_num, index_map[cur_index], 0]
        if len(recs_dic[cur_req_id]["result_map"].values()) == recs_dic[cur_req_id]["count"]:
            flat_list = zip(*recs_dic[cur_req_id]["result_map"].values())
            total_key_count = float(sum(flat_list[0]))
            logging.error("total_key_count: %s" % total_key_count)
            rate_value = 0.0
            for key_same in recs_dic[cur_req_id]["result_map"].values():
                if key_same[0] == 0 or key_same[2] == 0:
                    key_rate = 0
                    same_rate = 0
                elif key_same[2] == -1:
                    key_rate = 1
                    same_rate = -1
                else:
                    key_rate = min(max(key_same[0] * 9 / total_key_count, 0.6), 1.2)
                    same_rate = min(max(float(key_same[2]) / 30, 0.5), 1.3)
                cur_rate = key_rate * same_rate * key_same[1]
                logging.error("key_rate: %s, same_rate: %s, cur_rate: %s" % (key_rate, same_rate, cur_rate))
                rate_value = rate_value + cur_rate

            rate_value = min(max(float(total_key_count) / 600, 0.6), 1.1) * rate_value
            logging.error("final rate_value: %s" % rate_value)
            self.update_res_picture(rate_value)
            recs_dic.pop(cur_req_id)

        logging.error('finish')
        msg.ack()
    def update_res_picture(self,similarity):
        global db_conn
        if not hasattr(self,'keynumber'):
            self.keynumber = 0
        db = db_conn.resource
        query = {
            'md5' : self.origin_md5,
            'delete_flag' : {'$ne' : '1'}
        }
        res_picture_query = db.res_picture.find_one(query)
        if not res_picture_query:
            logging.error('has res_picture | md5:%s'%self.origin_md5)
            return
        #if not 'similarity' in res_picture_query or res_picture_query.get('similarity',-1) <  similarity :
        update_set = {
            '$set': {
                'training_state' : TRAINING_STATE_SUCCESS,
                'similarity' : similarity,
                'key_point'  : self.keynumber
            }
        }
        query = {
            'md5' : self.origin_md5,
            'delete_flag' : {'$ne' : '1'}
        }
        db.res_picture.update(query,update_set,multi = True)


def topic_finish_request(msg):
    topic_callback = TopicCallback()
    topic_callback.topic_callback_request(msg)

def on_amqp_connection():
    global ch,fastdfs_client
    ch = conn.channel()
    ch.exchange_declare(exchange=options.Search_CallBack_Exchange, type="topic", durable=True)
    ch.queue_declare(queue=options.Topic_CallBack_Queue, durable=True)
    ch.queue_declare(queue = options.Search_CallBack_Queue, durable=True)
    ch.queue_bind(queue=options.Topic_CallBack_Queue, exchange=options.Search_CallBack_Exchange, \
                  routing_key=options.Topic_CallBack_Queue)
    ch.queue_bind(queue=options.Search_CallBack_Queue, exchange=options.Search_CallBack_Exchange,
        routing_key=options.Search_CallBack_Queue)
    ch.consume(options.Topic_CallBack_Queue, topic_finish_request, no_ack=False)
    ch.consume(options.Search_CallBack_Queue, consumer.consume, no_ack=False)
    logging.error('consume:[E %s][Q %s][R %s]'%(options.Search_CallBack_Exchange,options.Topic_CallBack_Queue,options.Topic_CallBack_Queue))
    logging.error('consume:[E %s][Q %s][R %s]'%(options.Search_CallBack_Exchange,options.Search_CallBack_Queue,options.Search_CallBack_Queue))


    base_file_path =  os.path.dirname(__file__)
    fastdfs_client_conf_path = os.path.join(base_file_path, "fastdfs_client.conf")
    fastdfs_client = DfsUtils(fastdfs_client_conf_path)

def on_amqp_disconnect():
    logging.error( "[*]on_amqp_disconnect")
    logging.error( "[*]will close tornado")
    try:
        tornado.ioloop.IOLoop.instance().stop()
    except AttributeError, e:
        pass
    finally:
        exit(1)
ch = None
options.log_to_stderr = True
options.parse_command_line()


mongo_host   = 'mongodb://%s:%s@%s' %(setting.MONGO_USER,setting.MONGO_PASS,setting.MONGO_HOST)
mongo_port   = setting.MONGO_PORT
mq_host      = setting.MQ_HOST
mq_heartbeat = setting.MQ_HEARTBEAT
mq_username  = setting.MQ_USERNAME
mq_password  = setting.MQ_PASSWORD

options['Search_CallBack_Queue'].set(options.Queue_Topic_Callback_Search_Callback_Queue)
db_conn = Conn(mongo_host, mongo_port)
conn = Connection(host=mq_host, username=mq_username, password=mq_password, heartbeat=mq_heartbeat)

base_file_path         = os.path.dirname(__file__)
libimagerecs_path = os.path.join(base_file_path, "./librecs.so")
librecs = cdll.LoadLibrary(libimagerecs_path)

conn.connect(on_amqp_connection)
conn.on_disconnect = on_amqp_disconnect

try:
    ioloop.IOLoop.instance().start()
except KeyboardInterrupt:
    conn.close()
