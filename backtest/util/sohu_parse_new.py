#!/usr/bin/python
# -*- coding: utf-8 -*-
__author__ = 'zhangmin'
'''
搜狐视频首页的链接包含list,如http://tv.sohu.com/s2015/newslist/?vid=2534613，云视链无法解析需要先解到其对应的播放地址
'''
from wantv import get_video_url
import common
import re
import simplejson

def SearchSohuVideoUrl(v_url, useragent=None):
    req = re.compile(r"vid=([^&\?]*)") #匹配出vid
    vids = req.findall(v_url)
    if not vids:
        return v_url
    vid = vids[0]
    html = common.get_html(v_url) #匹配出playlistid，用于发送请求
    req = re.compile(r"var PLAYLIST_ID = ([^;]*)")
    play_list_ids = req.findall(html)
    if not play_list_ids:
        return v_url
    play_list_id = play_list_ids[0]
    request_url = "http://pl.hd.sohu.com/videolist?callback=newest&playlistid=%s"%play_list_id
    response = common.get_html(request_url, encoding='GBK')
    try:
        json_obj = simplejson.loads(response[response.find('{') : response.rfind('}') + 1])
        videos =json_obj.get('videos', [])
        for video in videos:
            cur_vid = str(video.get('vid', ''))
            if vid != cur_vid: continue
            cur_url = video.get('pageUrl', '')
            if not cur_url:
                return v_url
            return cur_url
    except Exception:
        return v_url


def sohu_download(v_url, useragent=None):
    if 'list' in v_url:
        v_url = SearchSohuVideoUrl(v_url, useragent)
    return get_video_url(v_url, useragent)


if __name__ == '__main__':
    import sys
    useragent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.146 Safari/537.36"
    try:
        v_url = sys.argv[1]
    except Exception:
        v_url = "http://tv.sohu.com/s2015/newslist/?vid=2534613"
    print "url=%s, v_url="%v_url,sohu_download(v_url, useragent)


