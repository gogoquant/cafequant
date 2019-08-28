# -*- coding: UTF-8 -*-

__author__ = 'wangande'

import sys
import logging
import hashlib
import time
import re
import getopt

import tornado.gen
import simplejson
from tornado import ioloop
from tornado.options import options

import setting
from stormed import Connection as StormedConnection
from models.ia.analytics_model import BehaviorLog
import util.time_common as time_common
from services.tool.db import Mongodb
from iseecore.fastdfs_utils import DfsUtils
from iseecore.models import SyncData

YIXUN_RES_LOG_NODE = setting.YIXUN_RES_LOG_NODE

ACTION_CODE = {
    'search': '0',
    'get_animation': '1',
    'download_template': '2',
    'clear': '3',
    'appear': '4',
    'touch_event': '5',
    'discoveryTheme': '6',
    'share': '7',
    'screen_record': '8',
    'H5': '9',
    'page': '10'
}

SHARE_TYPE = {'pic_share': 0, 'h5_share': 1, 'video_share': 2}

TIME_ZONE = {
    0: "0-9",
    1: "10-19",
    2: "20-29",
    3: "30-39",
    4: "40-49",
    5: "50-59",
    6: "60-69",
    7: "70-79",
    8: "80-89",
    9: "90-",
}

NEW_TIME_ZONE = {
    0: "0-5",
    1: "5-10",
    2: "10-15",
    3: "15-20",
    4: "20-25",
    5: "25-30",
    6: "30-35",
    7: "35-40",
    8: "40-45",
    9: "45-50",
    10: "50-55",
    11: "55-60",
    12: "60-65",
    13: "65-70",
    14: "70-75",
    15: "75-80",
    16: "80-85",
    17: "85-90",
    18: "90-"
}


class YiXinClientlogSync(object):

    def __init__(self, day_start, day_end, app_id):
        self.day_s = day_start
        self.day_e = day_end
        self.app_id = app_id
        mongodb = Mongodb()

        mq_host = setting.MQ_HOST
        mq_heartbeat = 3000
        mq_username = setting.MQ_USERNAME
        mq_password = setting.MQ_PASSWORD
        global mq_conn
        mq_conn = StormedConnection(
            host=mq_host,
            username=mq_username,
            password=mq_password,
            heartbeat=mq_heartbeat)

        mq_conn.connect(self.run)
        mq_conn.on_disconnect = self.on_amqp_disconnect

    @tornado.gen.engine
    def run(self):
        logging.error("yixin client log sync start....")
        ch = mq_conn.channel()
        options["ch"] = ch
        SyncData.configure(ch)
        self.fastdfs_client = DfsUtils()
        self.behavior_log_class = BehaviorLog()
        from services.resource_statistics import ResourceStatisticsService
        from services.resource import ResourcePictureService
        self.res_statistics_service = ResourceStatisticsService()
        self.res_pic_service = ResourcePictureService()
        self.systems = ['ios', 'android']

        # todo timezone except record
        batch_start_time = time_common.string_to_timestamp(
            self.day_s, time_common.FULL_PATTERN)
        batch_end_time = time_common.string_to_timestamp(
            self.day_e, time_common.FULL_PATTERN)

        # 原始数据
        query = {
            'date': {
                "$gte": batch_start_time,
                "$lt": batch_end_time
            },
            'app_id': self.app_id
        }

        # 获取符合条件的客户端日志
        behaviors = yield tornado.gen.Task(self.behavior_log_class.get_list,
                                           query)
        logging.error("get client record is %d" % len(behaviors))

        # 把原始数据读取成用户数据
        res_action_data, imei_action_data = yield tornado.gen.Task(
            self.make_data_from_log, behaviors)

        # 分析主题统计相关的用户数据
        res_stats_records = yield tornado.gen.Task(
            self.make_res_action_data_to_record, res_action_data)

        # 将统计数据整理成记录列表
        res_stats_lists = yield tornado.gen.Task(self.make_record_to_list,
                                                 res_stats_records)

        yield tornado.gen.Task(
            self.res_statistics_service.client_insert_or_update_list,
            res_stats_lists)
        doc = simplejson.dumps(res_stats_lists)

        logging.error("[*] sync_info send doc:%s" % (doc))
        # with open("/home/wangande/MyProject/client_log.txt", "a") as f:
        #     f.write(doc)
        #     f.write("\n")

        mq_conn.close(self.done)

    @tornado.gen.engine
    def make_data_from_log(self, behavior_files, callback=None):

        # res_action_data = {
        #     "activetime": {     # 按时间，正常情况下，统计时间应该相同，都为同一天，之所以再根据时间为key
        #                         # 是防止有些日志没有即使上传，日志中存在其他其它天数的日志记录
        #         "imei": {       # 按客户端统计
        #             "md5": {    # 每个客户端当天访问的主题统计
        #                 'read_num': 0,      # 阅读总数
        #                 'read_time': {},    # 阅读时间段的统计
        #                 'share_num': 0,     # 分享次数统计
        #                 'H5_share_num': 0,  # H5分享次数统计
        #             },
        #         },
        #     },
        # }

        # imei_action_data = {
        #   "activetime": {
        #       "imei": {
        #           "page": {
        #               "page_index":1
        #           }
        #       }
        #   }
        # }
        #
        #

        res_action_data = {}
        imei_action_data = {}

        # 获取所有客户端的日志
        # i = 0
        for behavior_file in behavior_files:
            file_id = behavior_file.get('file_id', None)
            file_md5 = behavior_file.get('md5', None)
            if not file_id or not file_md5:
                logging.error("old log file")
                continue

            # 获取客户端内容
            content = self.fastdfs_client.download_to_buffer(str(file_id))
            # file = "./" + str(i) + ".txt"
            #
            # with open(file, "wb+") as f:
            #     f.write(content)
            # i += 1
            # 获取日志的每行的内容
            behaviors = content.split('\n')

            # 获取客户端的信息
            # imei:866299020614807,device:vivo vivo Y13L,system:android,os:4.4.4(19),App:幻视,software:1.4,app-from:idealsee;
            imei_info = behaviors[0].split(';')[0].split(',')
            imei_data = {}
            for info in imei_info:
                key_value = info.split(':')
                imei_data[key_value[0]] = key_value[1]
            # system = imei_data.get('system', '')
            # version = imei_data.get('software', '')
            imei = imei_data.get('imei', '')
            if not imei:
                logging.error('log error not imei sn')
                continue

            app_launch_time = time_common.string_to_timestamp(
                self.day_s, time_common.FULL_PATTERN) * 1000

            timezone_count_flag = 0
            timezone_count_time = 0
            timezone_count_md5 = ''

            for lines in behaviors[1:]:
                if not lines or not lines.strip():
                    continue
                for record in lines.split(';'):
                    # 更新启动时间
                    if 'AppLaunchx' in record and ':' in record:
                        AppLaunchx_l = record.split(':')
                        if AppLaunchx_l[1].isdigit():
                            app_launch_time = int(AppLaunchx_l[1])
                        continue

                    # 退出的日志
                    if 'AppExit' in record and ':' in record:
                        continue

                    if ',' not in record:
                        continue

                    record_list = record.split(',')

                    # 不是统计log记录，继续
                    if (not record_list[0].isdigit()) or (
                            not record_list[1].isdigit()) or (
                                not record_list[2].isdigit()):
                        continue

                    record_time = time_common.timestamp_to_string(
                        app_launch_time / 1000, time_common.DAY_PATTERN)

                    res_record_map = res_action_data.get(record_time, {})
                    res_imei_map = res_record_map.get(imei, {})

                    imei_record_map = imei_action_data.get(record_time, {})
                    imei_imei_map = imei_record_map.get(imei, {})

                    init_map = {
                        'read_num': 0,
                        'read_time': {},
                        'share_num': 0,
                        'h5_share_num': 0,
                        'video_share_num': 0
                    }
                    if record_list[1].strip(
                    ) == ACTION_CODE['get_animation'] and record_list[2].strip(
                    ) == '1' and len(record_list) > 3:
                        md5 = record_list[3]
                        md5_map = res_imei_map.get(md5, init_map)
                        md5_map['read_num'] += 1
                        res_imei_map[md5] = md5_map

                        # 记录
                        timezone_count_flag = 1
                        timezone_count_time = int(record_list[0])
                        timezone_count_md5 = md5

                    elif record_list[1].strip(
                    ) == ACTION_CODE['share'] and len(record_list) > 4:

                        share_type = 0
                        if len(record_list) > 5:
                            share_type == int(record_list[5])

                        md5 = record_list[4]
                        md5_map = res_imei_map.get(md5, init_map)
                        if share_type == SHARE_TYPE['video_share']:
                            md5_map['video_share_num'] += 1
                        elif share_type == SHARE_TYPE['h5_share']:
                            md5_map['h5_share_num'] += 1
                        else:
                            md5_map['share_num'] += 1

                        res_imei_map[md5] = md5_map

                    elif record_list[1].strip() == ACTION_CODE['clear']:
                        if timezone_count_flag and timezone_count_md5 and res_imei_map.get(
                                timezone_count_md5, None) and (int(
                                    record_list[0]) > timezone_count_time):
                            md5_map = res_imei_map.get(timezone_count_md5)
                            read_time = (int(record_list[0]) -
                                         timezone_count_time) / 1000
                            timezone = self.change_time_to_timezone(read_time)
                            md5_map['read_time'][timezone] = md5_map[
                                'read_time'][timezone] + 1 if md5_map[
                                    'read_time'].get(timezone) else 1
                            res_imei_map[timezone_count_md5] = md5_map

                            timezone_count_flag = 0
                            timezone_count_time = 0
                            timezone_count_md5 = ''

                    # elif record_list[1].strip() == ACTION_CODE['page']:
                    #     page_map = imei_imei_map.get("page", {})
                    #     page_index = record_list[3]
                    #     if page_index in page_map:
                    #         page_map[page_index] += 1
                    #     else:
                    #         page_map[page_index] = 1
                    #     imei_imei_map["page"] = page_map

                    res_record_map[imei] = res_imei_map
                    res_action_data[record_time] = res_record_map

                    imei_record_map[imei] = imei_imei_map
                    imei_action_data[record_time] = imei_record_map

        callback((res_action_data, imei_action_data))

    @tornado.gen.engine
    def make_res_action_data_to_record(self, log_data, callback=None):
        # stats_record = {
        #     "stats_md5": {                            # 当天主题的访问情况， 以主题的MD5+时间生成statistics_md5作为唯一标识
        #         'md5': '',                            # 主题md5
        #         'stats_date': '',                     # 日志的时间
        #         'stats_md5': '',                      # 统计的md5 唯一标识
        #         'uv_read': 0,                         # uv阅读数
        #         'pv_read': 0,                         # pv阅读数
        #         'uv_share': 0,                        # uv分享数
        #         'pv_share': 0,                        # pv分享数
        #         'uv_h5_share': 0,                     # uvH5分享数
        #         'pv_h5_share': 0,                     # pvH5分享数
        #         'uv_video_share': 0,                  # uv视频分享数
        #         'pv_video_share': 0,                  # pv视频分享数
        #         'imei_user': [],                      # 当天访问该主题的用户列表
        #         'read_time': {},                      # 当天访问该主题的时间段统计
        #     }
        # }

        res_stats_records = {}
        for activetime, activetime_map in log_data.items():
            if not activetime_map:
                continue
            stats_date = time_common.string_to_timestamp(
                activetime, time_common.DAY_PATTERN)

            for imei, imei_map in activetime_map.items():
                if not imei_map:
                    continue

                for md5, md5_map in imei_map.items():
                    m = hashlib.md5()
                    m.update(md5 + activetime + YIXUN_RES_LOG_NODE)
                    stats_md5 = m.hexdigest()  # 某天某主题的统计，主题和时间确定唯一

                    init_map = {
                        'md5': md5,
                        'stats_date': stats_date,
                        'stats_md5': stats_md5,
                        'uv_read': 0,
                        'pv_read': 0,
                        'uv_share': 0,
                        'pv_share': 0,
                        'uv_h5_share': 0,
                        'pv_h5_share': 0,
                        'uv_video_share': 0,
                        'pv_video_share': 0,
                        'imei_user': [],
                        'read_time': {},
                    }

                    count_map = res_stats_records.get(stats_md5, init_map)

                    count_map['uv_read'] += md5_map['read_num']
                    count_map['pv_read'] += 1 if md5_map.get(
                        'read_num') else 0  # 同个用户当天访问同个主题多次，pv都记为1次

                    count_map['uv_share'] += md5_map['share_num']
                    count_map['pv_share'] += 1 if md5_map.get(
                        'share_num') else 0  # 同个用户当天分享同个主题多次，pv都记为1次

                    count_map['uv_h5_share'] += md5_map['h5_share_num']
                    count_map['pv_h5_share'] += 1 if md5_map.get(
                        'h5_share_num') else 0  # 同个用户当天分享同个主题多次，pv都记为1次

                    count_map['uv_video_share'] += md5_map['video_share_num']
                    count_map['pv_video_share'] += 1 if md5_map.get(
                        'video_share_num') else 0  # 同个用户当天分享同个主题多次，pv都记为1次

                    count_map['imei_user'].append(imei)

                    for timezone, time_count in md5_map.get('read_time',
                                                            {}).items():
                        if timezone in count_map['read_time']:
                            count_map['read_time'][timezone] += time_count
                        else:
                            count_map['read_time'][timezone] = time_count

                    res_stats_records[stats_md5] = count_map

        callback(res_stats_records)

    @tornado.gen.engine
    def make_imei_action_data_to_record(self, log_data, callback=None):
        # stats_record = {
        #     "stats_md5": {                          # 当天主题的访问情况， 以主题的MD5+时间生成stats_md5作为唯一标识
        #         'imei': '',                         # 客户端imei
        #         'stats_date': '',                   # 日志的时间
        #         'stats_md5': '',                    # 统计的md5 唯一标识
        #         'page': {},                         # 当天访问该主题的时间段统计
        #     }
        # }

        statistics_record = {}
        for activetime, activetime_map in log_data.items():
            if not activetime_map:
                continue
            statistics_datetime = time_common.string_to_timestamp(
                activetime, time_common.DAY_PATTERN)

            for imei, imei_map in activetime_map.items():
                if not imei_map:
                    continue

                for md5, md5_map in imei_map.items():
                    m = hashlib.md5()
                    m.update(md5 + activetime)
                    statistics_md5 = m.hexdigest()  # 某天某主题的统计，主题和时间确定唯一

                    init_map = {
                        'md5': md5,
                        'statistics_datetime': statistics_datetime,
                        'statistics_md5': statistics_md5,
                        'uv_read_num': 0,
                        'pv_read_num': 0,
                        'uv_share_num': 0,
                        'pv_share_num': 0,
                        'uv_h5_share_num': 0,
                        'pv_h5_share_num': 0,
                        'imei_user': [],
                        'read_time': {},
                    }

                    count_map = statistics_record.get(statistics_md5, init_map)

                    count_map['uv_read_num'] += md5_map['read_num']
                    count_map['pv_read_num'] += 1 if md5_map.get(
                        'read_num') else 0  # 同个用户当天访问同个主题多次，pv都记为1次

                    count_map['uv_share_num'] += md5_map['share_num']
                    count_map['pv_share_num'] += 1 if md5_map.get(
                        'share_num') else 0  # 同个用户当天分享同个主题多次，pv都记为1次

                    count_map['uv_h5_share_num'] += md5_map['h5_share_num']
                    count_map['pv_h5_share_num'] += 1 if md5_map.get(
                        'h5_share_num') else 0  # 同个用户当天分享同个主题多次，pv都记为1次

                    count_map['imei_user'].append(imei)

                    for timezone, time_count in md5_map.get('read_time',
                                                            {}).items():
                        if timezone in count_map['read_time']:
                            count_map['read_time'][timezone] += time_count
                        else:
                            count_map['read_time'][timezone] = time_count

                    statistics_record[statistics_md5] = count_map

        callback(statistics_record)

    # def change_time_to_timezone(self, read_time):
    #
    #     read_time = int(read_time)
    #     time_key = read_time/10
    #     timezone = TIME_ZONE.get(time_key, "0-9")
    #
    #     return timezone

    def change_time_to_timezone(self, read_time):

        read_time = int(read_time)
        time_key = read_time / 5
        timezone = NEW_TIME_ZONE.get(time_key, "90-")

        return timezone

    @tornado.gen.engine
    def make_record_to_list(self, records, callback=None):
        res_statistics_record_list = []

        for k, v in records.items():
            res_pic = yield tornado.gen.Task(
                self.res_pic_service.one, md5=v.get('md5'))
            if res_pic:
                v['resource_id'] = res_pic['resource_id']
                v['resource_picture_id'] = res_pic['resource_picture_id']
                del v['md5']
                res_statistics_record_list.append(v)
        callback(res_statistics_record_list)

    def done(self):
        ioloop.IOLoop.instance().stop()

    def on_amqp_disconnect(self):
        print "[*]on_amqp_disconnect"
        print "[*]will close tornado"
        try:
            tornado.ioloop.IOLoop.instance.stop()
        except AttributeError, e:
            pass
        finally:
            exit(1)


def usage():
    print "-h get help info"
    print "-t '2015-01-02 08:00' '2015-01-01 08:00'"
    print "-z '2'"


if __name__ == '__main__':
    logging.basicConfig(
        format='[%(asctime)s %(filename)s:%(lineno)d %(levelname)s] %(message)s',
        level=logging.ERROR)
    options.logging = 'debug'
    opts, args = getopt.getopt(sys.argv[1:], "htz")

    day_start = ""
    day_end = ""
    if len(sys.argv) > 1:
        for op, value in opts:
            if op == "-t":
                if len(sys.argv) != 4:
                    usage()
                else:
                    day_start = sys.argv[-1]
                    day_end = sys.argv[-2]

            elif op == "-z":
                if len(sys.argv) != 3:
                    day_end = time_common.timestamp_to_string(
                        int(time.time()), time_common.FULL_PATTERN)
                    day_start = time_common.timestamp_to_string(
                        int(time.time()) - time_common.hour_seconds * 2,
                        time_common.FULL_PATTERN)
                else:
                    day_end = time_common.timestamp_to_string(
                        int(time.time()), time_common.FULL_PATTERN)
                    day_start = time_common.timestamp_to_string(
                        int(time.time()) -
                        time_common.hour_seconds * int(sys.argv[-1]),
                        time_common.FULL_PATTERN)

            elif op == "-h":
                usage()
            else:
                usage()
    else:
        # 昨天
        start_time = int(time.time()) - time_common.day_seconds
        end_time = int(time.time())

        day_start = time_common.timestamp_to_string(start_time,
                                                    time_common.FULL_PATTERN)
        day_end = time_common.timestamp_to_string(end_time,
                                                  time_common.FULL_PATTERN)
    p = r"\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}"

    app_id = setting.DEFAULT_APP_ID
    if re.match(p, day_start) and re.match(p, day_end):
        logging.error(day_start)
        logging.error(day_end)
        YiXinClientlogSync(day_start, day_end, app_id)
        try:
            ioloop.IOLoop.instance().start()
        except KeyboardInterrupt:
            ioloop.IOLoop.instance().stop()
            mq_conn.close()
    else:
        pass
