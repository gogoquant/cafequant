#-*- coding: UTF-8 -*-

import simplejson
import sys
from common import *


def tudou_download_by_iid(iid,
                          title,
                          output_dir='.',
                          merge=True,
                          info_only=False,
                          useragent=None):
    data_string = get_decoded_html(
        'http://www.tudou.com/outplay/goto/getItemSegs.action?iid=%s' % iid,
        useragent=useragent)
    #print "data_string: %s" % str(data_string)
    data = simplejson.loads(data_string)
    #print "iid: %s" % iid
    #print "data: %s" % str(data)
    #print "title: %s" % title
    vids = []
    for k in data:
        if len(data[k]) > 0:
            if data[k][0].has_key("k") and data[k][0].has_key("size"):
                vids.append({"k": data[k][0]["k"], "size": data[k][0]["size"]})

    #temp = max(vids, key=lambda x:x["size"])
    #vid, size = temp["k"], temp["size"]
    from xml.dom.minidom import parseString
    url = None
    ext = None
    for vid_info in vids:
        vid, size = vid_info["k"], vid_info["size"]
        tmp_html = 'http://v2.tudou.com/f?id=%s' % vid
        xml = get_html(tmp_html, useragent=useragent)
        doc = parseString(xml)
        url = [
            n.firstChild.nodeValue.strip()
            for n in doc.getElementsByTagName('f')
        ][0]
        ext = r1(r'http://[\w.]*/(\w+)/[\w.]*', url)
        if ext == "mp4":
            break
    #tmp_html = 'http://ct.v2.tudou.com/f?id=%s' % vid
    #tmp_html = 'http://v2.tudou.com/f?id=%s' % vid
    #print "tmp_html: %s" % tmp_html
    #xml = get_html(tmp_html, useragent=useragent)
    #from xml.dom.minidom import parseString
    #doc = parseString(xml)
    #url = [n.firstChild.nodeValue.strip() for n in doc.getElementsByTagName('f')][0]

    #ext = r1(r'http://[\w.]*/(\w+)/[\w.]*', url)
    """print_info(site_info, title, ext, size)
    if not info_only:
        print url
        #download_urls([url], title, ext, size, output_dir = output_dir, merge = merge)"""
    return url


def tudou_download_by_id(id,
                         title,
                         output_dir='.',
                         merge=True,
                         info_only=False,
                         useragent=None):
    html = get_html(
        'http://www.tudou.com/programs/view/%s/' % id, useragent=useragent)

    iid = r1(r'iid\s*[:=]\s*(\S+)', html)
    title = r1(r'kw\s*[:=]\s*[\'\"]([^\']+?)[\'\"]', html)
    tudou_download_by_iid(
        iid, title, output_dir=output_dir, merge=merge, info_only=info_only)


def tudou_download(url,
                   output_dir='.',
                   merge=True,
                   info_only=False,
                   useragent=None):
    # Embedded player
    id = r1(r'http://www.tudou.com/v/([^/]+)/', url)
    if id:
        return tudou_download_by_id(
            id, title="", info_only=info_only, useragent=useragent)

    html = get_decoded_html(url, useragent=useragent)

    title = r1(r'kw\s*[:=]\s*[\'\"]([^\']+?)[\'\"]', html)
    assert title
    title = unescape_html(title)

    vcode = r1(r'vcode\s*[:=]\s*\'([^\']+)\'', html)
    if vcode:
        print "youku: %s" % vcode
        youku_url = 'http://v.youku.com/v_show/id_%s.html' % vcode
        from youku_parse import youku_download
        return youku_download(youku_url)

    iid = r1(r'iid\s*[:=]\s*(\d+)', html)
    if not iid:
        return tudou_download_playlist(
            url, output_dir, merge, info_only, useragent=useragent)

    result_urls = []
    result_urls.append(
        tudou_download_by_iid(
            iid,
            title,
            output_dir=output_dir,
            merge=merge,
            info_only=info_only,
            useragent=useragent))
    return result_urls


def parse_playlist(url, useragent=None):
    aid = r1('http://www.tudou.com/playlist/p/a(\d+)(?:i\d+)?\.html', url)
    html = get_decoded_html(url, useragent=useragent)
    if not aid:
        aid = r1(r"aid\s*[:=]\s*'(\d+)'", html)
    if re.match(r'http://www.tudou.com/albumcover/', url):
        atitle = r1(r"title\s*:\s*'([^']+)'", html)
    elif re.match(r'http://www.tudou.com/playlist/p/', url):
        atitle = r1(r'atitle\s*=\s*"([^"]+)"', html)
    else:
        raise NotImplementedError(url)
    assert aid
    assert atitle
    import json
    #url = 'http://www.tudou.com/playlist/service/getZyAlbumItems.html?aid='+aid
    url = 'http://www.tudou.com/playlist/service/getAlbumItems.html?aid=' + aid
    return [(atitle + '-' + x['title'], str(x['itemId']))
            for x in json.loads(get_html(url, useragent=useragent))['message']]


def tudou_download_playlist(url,
                            output_dir='.',
                            merge=True,
                            info_only=False,
                            useragent=None):
    videos = parse_playlist(url, useragent=useragent)
    url_list = []
    for i, (title, id) in enumerate(videos):
        print('Processing %s of %s videos...' % (i + 1, len(videos)))
        url_list.append(
            tudou_download_by_iid(
                id,
                title,
                output_dir=output_dir,
                merge=merge,
                info_only=info_only,
                useragent=useragent))
    return url_list


site_info = "Tudou.com"


def main():
    if (len(sys.argv) > 1):
        useragent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.146 Safari/537.36"
        print tudou_download(sys.argv[1], useragent=useragent)
    else:
        print "usage: python tudou_parse.py url"


if __name__ == '__main__':
    main()
