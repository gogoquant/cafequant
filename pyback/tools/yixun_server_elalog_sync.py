#!/usr/bin/env python
# coding=utf-8

__author__ = 'wangande'

import sys
import logging
import hashlib
import urllib2
import getopt
import time
import re
import tornado.gen
import simplejson
import setting
import util.time_common as time_common

from tornado import ioloop
from tornado.options import options
from stormed import Connection as StormedConnection
from services.tool.db import Mongodb
from models.resource.comment_model import ResourceComment
from models.user.collection_model import Collection
from models.ia.analytics_model import Imei
from iseecore.models import SyncData

YIXUN_RES_LOG_NODE = setting.YIXUN_RES_LOG_NODE


class YiXinServerElalogSync(object):

    def __init__(self, day_start, day_end):
        self.day_s = day_start
        self.day_e = day_end
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
        logging.error("start")
        ch = mq_conn.channel()
        options["ch"] = ch
        SyncData.configure(ch)
        # self.fastdfs_client = DfsUtils()
        from services.resource_statistics import ResourceStatisticsService
        from services.resource import ResourcePictureService
        from services.ela_logging import ElasticLoggingService
        from services.ip_db import IpService
        from services.video import VideoService
        from services.browser import BrowserService
        self.resource_picture_service = ResourcePictureService()
        self.res_statistics_service = ResourceStatisticsService()
        self.ela_logging_service = ElasticLoggingService()
        self.browser_s = BrowserService()
        self.collection_m = Collection()
        self.comment_m = ResourceComment()
        self.imei_m = Imei()
        self.ip_s = IpService()
        self.video_s = VideoService()

        # todo @timestamp except
        batch_start_time, batch_end_time = self.get_ela_query_time(
            self.day_s, self.day_e)
        scan_info_pattern = ".*scan_info.*"
        h5_share_read_pattern = ".*h5_share_read_info.*"
        video_share_read_pattern = ".*video_share_read_info.*"
        ar_create_browser_pattern = ".*ar_create_browser.*"
        logging.error('%s-%s' % (batch_start_time, batch_end_time))
        # Ela logging 全文查询， 查询从 batch_start_time 到 batch_end_time 符合条件 pattern 的log
        query = {
            "query": {
                "bool": {
                    "must": [{
                        "bool": {
                            "should": [
                                {
                                    "regexp": {
                                        "log": scan_info_pattern
                                    }
                                },
                                {
                                    "regexp": {
                                        "log": h5_share_read_pattern
                                    }
                                },
                                {
                                    "regexp": {
                                        "log": video_share_read_pattern
                                    }
                                },
                                {
                                    "regexp": {
                                        "log": ar_create_browser_pattern
                                    }
                                },
                            ]
                        }
                    }, {
                        "range": {
                            "@timestamp": {
                                "gte": batch_start_time,
                                "lte": batch_end_time,
                            }
                        }
                    }]
                }
            },
            # 'timeout': '120ms'
        }

        # 获取符合的log 记录
        log_data = yield tornado.gen.Task(self.get_logging_data, query)
        logging.error("logo_data is (%s) ." % len(log_data))
        # 统计log记录中的识别总数
        res_stats_records, imei_records, browser_stats_records = yield tornado.gen.Task(
            self._make_data_records, log_data)
        res_stats_list = yield tornado.gen.Task(self.make_records_to_list,
                                                res_stats_records)

        yield tornado.gen.Task(
            self.res_statistics_service.server_insert_or_update_list,
            res_stats_list)
        yield tornado.gen.Task(self.update_imei_list, imei_records)
        yield tornado.gen.Task(self.update_browser_stats, browser_stats_records)

        logging.error("[*] sync_info send doc:%s" %
                      simplejson.dumps(res_stats_list))
        logging.error("sync_info send imei list doc:%s" %
                      simplejson.dumps(imei_records))
        logging.error("sync_info ar create browser list doc:%s" %
                      simplejson.dumps(browser_stats_records))

        mq_conn.close(self.done)

    def get_ela_query_time(self, start_time, end_time):

        date_pattern = '%Y-%m-%dT%H:%M:%S+08:00'

        start_time_timestamp = time_common.string_to_timestamp(
            start_time, time_common.FULL_PATTERN) - time_common.hour_seconds * 8
        end_time_timestamp = time_common.string_to_timestamp(
            end_time, time_common.FULL_PATTERN) - time_common.hour_seconds * 8

        batch_start_time = time_common.timestamp_to_string(
            start_time_timestamp, date_pattern)
        batch_end_time = time_common.timestamp_to_string(
            end_time_timestamp, date_pattern)

        return batch_start_time, batch_end_time

    @tornado.gen.engine
    def get_logging_data(self, query, callback=None):
        log_list = []
        log_total = yield tornado.gen.Task(
            self.ela_logging_service.get_logging_total, query)
        if not log_total:
            logging.info('Not searching for the record')
        else:
            query['size'] = log_total
            log_list = yield tornado.gen.Task(
                self.ela_logging_service.get_logging_list, query)

        callback(log_list)

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

    @tornado.gen.engine
    def update_imei_list(self, imei_records, callback=None):

        for imei, imei_data in imei_records.items():

            old_imei = yield tornado.gen.Task(self.imei_m.find_one,
                                              {"imei": imei})
            if old_imei:
                network = old_imei.get('network', {})
                area = old_imei.get('area', {})

                for net_name, net_count in imei_data.get('network', {}).items():
                    if net_name in network:
                        network[net_name] += net_count
                    else:
                        network[net_name] = net_count

                for area_name, area_count in imei_data.get('area', {}).items():
                    if area_name in area:
                        area[area_name] += area_count
                    else:
                        area[area_name] = area_count

                update_data = {'network': network, 'area': area}
                yield tornado.gen.Task(self.imei_m.update, {"imei": imei},
                                       {"$set": update_data})

            else:
                imei_data["imei"] = imei
                imei_data["app_id"] = setting.DEFAULT_APP_ID
                yield tornado.gen.Task(self.imei_m.insert, imei_data)

        callback(None)

    @tornado.gen.engine
    def update_browser_stats(self, browser_stats_records, callback=None):
        browser_stats_list = []

        for k, v in browser_stats_records.iteritems():
            browser_stats_list.append(v)
        yield tornado.gen.Task(self.browser_s.insert_or_update_browser,
                               browser_stats_list)
        callback(None)

    @tornado.gen.engine
    def _get_network_area(self, client_ip, callback=None):
        try:
            network = ''
            area = ''
            ip_info = yield tornado.gen.Task(self.ip_s.select_ip_info,
                                             client_ip)
            if ip_info:
                network = ip_info.get('isp', '')
                area = ip_info.get('province', '')
            callback((network, area))
        except Exception, e:
            callback((None, None))

    @tornado.gen.engine
    def match_scan_record(self, match_info, res_stats_records, callback=None):

        md5 = match_info.get("md5", "")
        client_ip = match_info.get("client_ip")
        imei = match_info.get("imei", "")
        system = match_info.get("system", "")
        record_date = match_info.get("date", self.day_s)
        stats_date = time_common.string_to_timestamp(record_date,
                                                     time_common.DAY_PATTERN)

        if md5 and imei:
            m = hashlib.md5()
            m.update(md5 + record_date + YIXUN_RES_LOG_NODE)
            stats_md5 = m.hexdigest()
            network = ''
            area = ''
            if client_ip:
                network, area = yield tornado.gen.Task(self._get_network_area,
                                                       client_ip)

            init_record = {
                'md5': md5,  # 主题的MD5
                'stats_date': stats_date,  # 统计的时间
                'stats_md5': stats_md5,  # 记录标识确定唯一性
                'uv_scan': 0,  # uv扫描数
                'pv_scan_map': {},  # pv扫描
                'uv_share_read': 0,  # uv分享阅读数
                'uv_h5_share_read': 0,  # uvh5分享阅读
                'uv_video_share_read': 0,
                'pv_scan_network_map': {},  # 扫描网络
                'pv_scan_area_map': {},  # 扫描地区
                'pv_scan_dev_map': {},  # 扫描设备
                'uv_h5_share_read_network': {},  # h5分享阅读网络
                'uv_h5_share_read_area': {},  # h5分享阅读地区
                'uv_h5_share_read_dev': {},  # h5分析阅读设备
                'uv_video_share_network': {},  # 视频分享阅读网络
                'uv_video_share_area': {},  # 视频分享阅读地区
                'uv_video_share_dev': {},  # 视频分享阅读设备
                'imei_record_map': {}  # 主题浏览客户
            }

            record = res_stats_records.get(md5, init_record)

            record['uv_scan'] += 1
            record['pv_scan_map'][imei] = 1

            imei_record = record['imei_record_map'].get(imei, {
                "system": "",
                "network": {},
                "area": {}
            })

            if network:
                net_imei_map = record['pv_scan_network_map'].get(network, {})
                net_imei_map[imei] = 1
                record['pv_scan_network_map'][network] = net_imei_map
                imei_record["network"][network] = 1

            if area:
                area_imei_map = record['pv_scan_area_map'].get(area, {})
                area_imei_map[imei] = 1
                record['pv_scan_area_map'][area] = area_imei_map
                imei_record["area"][area] = 1

            if system:
                system_imei_map = record['pv_scan_dev_map'].get(system, {})
                system_imei_map[imei] = 1
                record['pv_scan_dev_map'][system] = system_imei_map
                imei_record["system"] = system

            record['imei_record_map'][imei] = imei_record

            res_stats_records[md5] = record
        callback(res_stats_records)

    @tornado.gen.engine
    def match_h5_share_read(self, match_info, res_stats_records, callback=None):
        resource_id = match_info.get('resource_id', '')
        client_ip = match_info.get('client_ip', '')
        system = match_info.get("system", "")
        record_date = match_info.get("date", self.day_s)
        stats_date = time_common.string_to_timestamp(record_date,
                                                     time_common.DAY_PATTERN)
        md5 = ""
        if resource_id:
            res_pic = yield tornado.gen.Task(
                self.resource_picture_service.one, resource_id=resource_id)
            if res_pic:
                md5 = res_pic.get("md5", "")

        if md5:
            m = hashlib.md5()
            m.update(md5 + record_date + YIXUN_RES_LOG_NODE)
            stats_md5 = m.hexdigest()
            network = ''
            area = ''
            if client_ip:
                network, area = yield tornado.gen.Task(self._get_network_area,
                                                       client_ip)

            init_record = {
                'md5': md5,  # 主题的MD5
                'stats_date': stats_date,  # 统计的时间
                'stats_md5': stats_md5,  # 记录标识确定唯一性
                'uv_scan': 0,  # uv扫描数
                'pv_scan_map': {},  # pv扫描
                'uv_share_read': 0,  # uv分享阅读数
                'uv_h5_share_read': 0,  # uvh5分享阅读
                'uv_video_share_read': 0,
                'pv_scan_network_map': {},  # 扫描网络
                'pv_scan_area_map': {},  # 扫描地区
                'pv_scan_dev_map': {},  # 扫描设备
                'uv_h5_share_read_network': {},  # h5分享阅读网络
                'uv_h5_share_read_area': {},  # h5分享阅读地区
                'uv_h5_share_read_dev': {},  # h5分析阅读设备
                'uv_video_share_network': {},  # 视频分享阅读网络
                'uv_video_share_area': {},  # 视频分享阅读地区
                'uv_video_share_dev': {},  # 视频分享阅读设备
                'imei_record_map': {}  # 主题浏览客户
            }

            record = res_stats_records.get(md5, init_record)
            record['uv_h5_share_read'] += 1

            if network:
                if network in record['uv_h5_share_read_network']:
                    record['uv_h5_share_read_network'][network] += 1
                else:
                    record['uv_h5_share_read_network'][network] = 1

            if area:
                if area in record['uv_h5_share_read_area']:
                    record['uv_h5_share_read_area'][area] += 1
                else:
                    record['uv_h5_share_read_area'][area] = 1

            if not system:
                system = "web"

            if system in record['uv_h5_share_read_dev']:
                record['uv_h5_share_read_dev'][system] += 1
            else:
                record['uv_h5_share_read_dev'][system] = 1

            res_stats_records[md5] = record

        callback(res_stats_records)

    @tornado.gen.engine
    def match_video_share_read(self,
                               match_info,
                               res_stats_records,
                               callback=None):

        video_id = match_info.get('video_id', '')
        client_ip = match_info.get('client_ip', '')
        system = match_info.get("system", "")
        record_date = match_info.get("date", self.day_s)
        stats_date = time_common.string_to_timestamp(record_date,
                                                     time_common.DAY_PATTERN)

        md5 = ""
        video = yield tornado.gen.Task(self.video_s.one, video_id)
        if video:
            resource_id = video.get("resource_id", "")
            if resource_id:
                resource_pic = yield tornado.gen.Task(
                    self.resource_picture_service.one, resource_id=resource_id)
                if resource_pic:
                    md5 = resource_pic.get("md5", "")

        if md5:
            m = hashlib.md5()
            m.update(md5 + record_date + YIXUN_RES_LOG_NODE)
            stats_md5 = m.hexdigest()
            network = ''
            area = ''
            if client_ip:
                network, area = yield tornado.gen.Task(self._get_network_area,
                                                       client_ip)

            init_record = {
                'md5': md5,  # 主题的MD5
                'stats_date': stats_date,  # 统计的时间
                'stats_md5': stats_md5,  # 记录标识确定唯一性
                'uv_scan': 0,  # uv扫描数
                'pv_scan_map': {},  # pv扫描
                'uv_share_read': 0,  # uv分享阅读数
                'uv_h5_share_read': 0,  # uvh5分享阅读
                'uv_video_share_read': 0,
                'pv_scan_network_map': {},  # 扫描网络
                'pv_scan_area_map': {},  # 扫描地区
                'pv_scan_dev_map': {},  # 扫描设备
                'uv_h5_share_read_network': {},  # h5分享阅读网络
                'uv_h5_share_read_area': {},  # h5分享阅读地区
                'uv_h5_share_read_dev': {},  # h5分析阅读设备
                'uv_video_share_network': {},  # 视频分享阅读网络
                'uv_video_share_area': {},  # 视频分享阅读地区
                'uv_video_share_dev': {},  # 视频分享阅读设备
                'imei_record_map': {}  # 主题浏览客户
            }

            record = res_stats_records.get(md5, init_record)
            record['uv_video_share_read'] += 1

            if network:
                if network in record["uv_video_share_network"]:
                    record["uv_video_share_network"][network] += 1
                else:
                    record["uv_video_share_network"][network] = 1

            if area:
                if area in record["uv_video_share_area"]:
                    record["uv_video_share_area"][area] += 1
                else:
                    record["uv_video_share_area"][area] = 1

            if not system:
                system = "web"

            if system in record['uv_video_share_dev']:
                record['uv_video_share_dev'][system] += 1
            else:
                record['uv_video_share_dev'][system] = 1

            res_stats_records[md5] = record

        callback(res_stats_records)

    @tornado.gen.engine
    def match_ar_create_browser(self,
                                match_info,
                                browser_stats_records,
                                callback=None):

        system = match_info.get("system", "")
        browser = match_info.get("browser", "")

        if system and browser:
            m = hashlib.md5()
            m.update(system + browser)
            stats_md5 = m.hexdigest()

            if stats_md5 not in browser_stats_records:
                browser_stats_records[stats_md5] = {
                    "browser": browser,
                    "system": system,
                    "uv_visit": 1,
                }
            else:
                browser_stats_records[stats_md5]["uv_visit"] += 1

        callback(browser_stats_records)

    @tornado.gen.engine
    def _make_data_records(self, logo_data, callback=None):

        scan_info_pattern = re.compile(r".*scan_info:(.*)")
        h5_share_read_pattern = re.compile(r".*h5_share_read_info:(.*)")
        video_share_read_pattern = re.compile(r".*video_share_read_info:(.*)")
        ar_create_browser_pattern = re.compile(r".*ar_create_browser:(.*)")

        res_stats_records = dict()
        imei_records = dict()
        browser_stats_records = dict()

        for log in logo_data:
            if scan_info_pattern.match(log):
                info = scan_info_pattern.match(log)
                match_info = eval(info.group(1))
                res_stats_records = yield tornado.gen.Task(
                    self.match_scan_record, match_info, res_stats_records)
            if h5_share_read_pattern.match(log):
                info = h5_share_read_pattern.match(log)
                match_info = eval(info.group(1))
                res_stats_records = yield tornado.gen.Task(
                    self.match_h5_share_read, match_info, res_stats_records)

            if video_share_read_pattern.match(log):
                info = video_share_read_pattern.match(log)
                match_info = eval(info.group(1))
                res_stats_records = yield tornado.gen.Task(
                    self.match_video_share_read, match_info, res_stats_records)
            if ar_create_browser_pattern.match(log):
                info = ar_create_browser_pattern.match(log)
                match_info = eval(info.group(1))
                browser_stats_records = yield tornado.gen.Task(
                    self.match_ar_create_browser, match_info,
                    browser_stats_records)

        # 整理 pv统计的数据
        for k, v in res_stats_records.items():
            v['pv_scan'] = len(v['pv_scan_map'])
            del v['pv_scan_map']

            v['pv_scan_network'] = {}
            for network, imei_num in v['pv_scan_network_map'].items():
                v['pv_scan_network'][network] = len(imei_num)
            del v['pv_scan_network_map']

            v['pv_scan_area'] = {}
            for area, imei_num in v['pv_scan_area_map'].items():
                v['pv_scan_area'][area] = len(imei_num)
            del v['pv_scan_area_map']

            v['pv_scan_dev'] = {}
            for area, imei_num in v['pv_scan_dev_map'].items():
                v['pv_scan_dev'][area] = len(imei_num)
            del v['pv_scan_dev_map']

            imei_user = []
            for imei, imei_value in v['imei_record_map'].items():
                imei_info = imei_records.get(imei, {
                    'system': "",
                    'network': {},
                    'area': {}
                })

                if imei_value["system"]:
                    imei_info["system"] = imei_value["system"]

                if imei_value["network"]:
                    imei_info["network"].update(imei_value["network"])

                if imei_value["area"]:
                    imei_info["area"].update(imei_value["area"])

                imei_records[imei] = imei_info

                if imei not in imei_user:
                    imei_user.append(imei)

            v["imei_user"] = imei_user

            del v['imei_record_map']

        callback((res_stats_records, imei_records, browser_stats_records))

    @tornado.gen.engine
    def make_records_to_list(self, res_stats_records, callback=None):
        res_stats_list = []
        for k, v in res_stats_records.items():
            res_pic = yield tornado.gen.Task(
                self.resource_picture_service.one, md5=v.get('md5'))
            if res_pic:
                v['resource_id'] = res_pic['resource_id']
                v['resource_picture_id'] = res_pic['resource_picture_id']
                query = {"resource_id": res_pic['resource_id']}

                # start_time = v['stats_date'] - time_common.day_seconds
                # end_time = v['stats_date']

                start_time = time_common.string_to_timestamp(
                    self.day_s, time_common.FULL_PATTERN)
                end_time = time_common.string_to_timestamp(
                    self.day_e, time_common.FULL_PATTERN)

                query["last_modify"] = {"$gte": start_time, "$lt": end_time}
                logging.error(query)
                # 获取统计时间的收藏数和评论
                v['collection'] = yield tornado.gen.Task(
                    self.collection_m.count, query.copy())
                v['comment'] = yield tornado.gen.Task(self.comment_m.count,
                                                      query.copy())
                res_stats_list.append(v)

        callback(res_stats_list)


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

    if re.match(p, day_start) and re.match(p, day_end):
        logging.error(day_start)
        logging.error(day_end)
        YiXinServerElalogSync(day_start, day_end)
        try:
            ioloop.IOLoop.instance().start()
        except KeyboardInterrupt:
            ioloop.IOLoop.instance().stop()
            mq_conn.close()
    else:
        pass
