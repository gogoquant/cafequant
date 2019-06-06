#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
通过云视链进行视频解析

支持 youku tudou 腾讯 乐视
不支持 iqiyi

@file:wantv.py
@modul:wantv
@date:2015-03-25
"""

import os
import sys
import simplejson
import logging
import base64
from common import *
import time

import redis

sys.path.append("../")
import setting

HTML_BASE_URL = "http://videojj.com/api/videos/parse?url="
REFERER_BASE_URL = "http://videojj.com/"
USER_AGENT = "Mozilla/5.0 (Linux; U; Android 4.3; en-us; SM-N900T Build/JSS15J) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30"
# qq    "HD", "SD",
# iqiyi "720P", "1080P", "Fluency", "HD", "Speed",
# youku "3GP-HD", "1080P", "MP4-HD", "Mobile-m3u8-1080P", "Mobile-m3u8-MP4-HD", "Mobile-m3u8-SD", "Mobile-m3u8-SuperHD", "SD", "SuperHD"
# tudou "Mobile-MP4-SD" "Mobile-m3u8-SD" "SD"
# sohu  "HD" "SD" "SuperHD"
MAP_SITES_TO_VIDEO_TYPES = {
    'youku': ('MP4-HD', '3GP-HD', 'Mobile-m3u8-MP4-HD', ),
    'qq': ('SD', 'HD', '720P', ),
    'tudou': ('Mobile-MP4-SD', 'HD',),
    'sohu': ('SD', 'HD', 'SuperHD', 'Orignal'),
}

URL_TIMEOUT = 1 * 60 * 60

pool = redis.ConnectionPool(host=setting.REDIS_HOST, port=setting.REDIS_PORT,db=0)
r = redis.Redis(connection_pool=pool)


def get_video_url_by_html(base64_url, useragent=None):
    logging.error("get video url by videojj web")
    html_url = HTML_BASE_URL + base64_url
    referer_url = REFERER_BASE_URL  # + base64_url
    if useragent is None:
        useragent = USER_AGENT

    appkey = "VJjJcUeTl"  # 账号
    appkey_b64 = base64.b64encode(appkey)

    # html_response = get_html(html_url,useragent=useragent,Referer=referer_url)
    response = get_videojj_response(html_url, appkey_b64)
    # response = simplejson.loads(html_response)
    assert response.get("status",-1) == 0,"error status"
    segs = response.get("msg",{}).get("segs",{})
    assert segs,"none video"
    site = response.get("msg", {}).get("site", "")
    for v_type in MAP_SITES_TO_VIDEO_TYPES.get(site, []):
        videos = segs.get(v_type,[])
        if videos:
            video_url = videos[0].get('url')
            break
    else:
        keys = segs.keys()
        video_url = segs[keys[0]][0].get('url')
    return video_url


def get_video_url_by_cache(url):
    return r.get(url)
"""
删除cache 中的video_url
"""


def del_video_url_in_cache(url):
    base64_url = encode_url(url)
    r.delete(base64_url)


def save_video_url(url, video_url):
    r.set(url,video_url)
    r.expire(url,URL_TIMEOUT)


def get_video_url(url, useragent=None):
    base64_url = encode_url(url)
    cache_url = get_video_url_by_cache(base64_url)

    if cache_url:
        return [cache_url]
    try:
        video_url = get_video_url_by_html(base64_url, useragent)
    except (IndexError, KeyError):
        logging.error("error occurred when getting video url by videojj web")
        return []
    assert video_url
    save_video_url(base64_url, video_url)
    return [video_url]


def encode_url(url):
    assert url
    base64_url = base64.encodestring(url).replace("\n", "")
    return base64_url

if __name__ == '__main__':
    try:
        html_url = sys.argv[1]
    except IndexError:
        # html_url = "http://www.tudou.com/listplay/Q0MiBE2DPCs/zAV9Di5yzu0.html"
        html_url = "http://my.tv.sohu.com/pl/9029407/81371700.shtml"
    assert html_url

    print html_url, "==>", get_video_url(html_url)

