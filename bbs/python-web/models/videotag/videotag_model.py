#!/usr/bin/python
# -*- coding: utf-8 -*-

from iseecore.models import AsyncBaseModel

class VideoTag(AsyncBaseModel):
    tag_id = True
    thumb_md5 = True  # 缩列图md
    time = True  # tag创建时间 ms，相当视频开始时间
    x = True  # tag创建矩形区域
    y = True  # tag创建矩形区域
    width = True  # tag创建矩形区域
    height = True  # tag创建矩形区域
    tracklines = True  # 跟踪列表 包含矩形区域rect(x,y,width,height,curr_time)

    # meta
    table = 'videotag'
    key = 'tag_id'


class VidetoTagCache(AsyncBaseModel):
    cache_id = True
    md5 = True
    tag_list = True

    # meta
    table = "videotag_cache"
    key = "cache_id"
