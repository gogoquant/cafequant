#!/usr/bin/env python

__all__ = ['iqiyi_download']
import os
import simplejson
from common import *


def iqiyi_download(url, output_dir='.', merge=True, info_only=False):
    html = get_html(url)

    videoId = r1(r'data-player-videoid="([^"]+)"', html)
    assert videoId

    tvId = r1(r'data-player-tvid="([^"]+)"', html)
    assert tvId
    req = re.compile(r'<script(.*)>([^<]*)</script>')
    scripts = req.findall(html)
    script_attrs, script_inner = zip(*scripts)
    sea_nums = sum(1 for item in script_attrs if -1 != item.find("sea1.2."))
    sbase_nums = sum(1 for item in script_inner if -1 != item.find('s/base'))
    cmd_js = "nodejs ./iqiyi_parse.js %s %s %s %s" % (url, tvId, sea_nums,
                                                      sbase_nums)

    jsStr = os.popen(cmd_js).read().strip()
    jsObj = jsStr.split(' ')

    sc = jsObj[0]
    src = jsObj[1]
    t = jsObj[2]
    jsT = jsObj[3]
    assert sc
    assert src
    metal_url = 'http://cache.m.iqiyi.com/jp/tmts/%s/%s/?type=mp4&sc=%s&src=%s&t=%s&__jsT=%s' % (
        tvId, videoId, sc, src, t, jsT)
    metal_response = get_html(metal_url)
    metal_response = metal_response.replace('var tvInfoJs=', '')
    metal_response = metal_response.replace('\'', '\"')

    metal_response = simplejson.loads(metal_response)
    mp4_url = metal_response['data'].get('mp4Url', None)
    if not mp4_url:
        m3u_url = metal_response['data'].get('m3u', None)
        assert m3u_url
        return [m3u_url]

    assert mp4_url

    temp_url = "http://data.video.qiyi.com/%s" % mp4_url.split("/")[-1].split(
        ".")[0] + ".ts"
    try:
        urllib2.urlopen(temp_url)
    except urllib2.HTTPError as e:
        key = r1(r'key=(.*)', e.geturl())
    assert key
    mp4_url += "?key=%s" % key

    return [mp4_url]


site_info = "iQIYI.com"
download = iqiyi_download
download_playlist = playlist_not_supported('iqiyi')

if __name__ == '__main__':
    html = "http://www.iqiyi.com/w_19rskmbsyx.html?vfrm=3-1-0-1"  #"http://www.iqiyi.com/v_19rro52aao.html"
    print html, '=>', iqiyi_download(html)
