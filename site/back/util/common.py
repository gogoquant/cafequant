# -*- coding: utf-8 -*-
import urllib2
import os.path
import sys
import re

default_encoding = sys.getfilesystemencoding()
if default_encoding.lower() == 'ascii':
    default_encoding = 'utf-8'


def to_native_string(s):
    if type(s) == unicode:
        return s.encode(default_encoding)
    else:
        return s


def tr(s):
    if type(s) == unicode:
        return s.encode(default_encoding)
    else:
        return s


def r1(pattern, text):
    m = re.search(pattern, text)
    if m:
        return m.group(1)


def match1(text, *patterns):
    """Scans through a string for substrings matched some patterns (first-subgroups only).

    Args:
        text: A string to be scanned.
        patterns: Arbitrary number of regex patterns.

    Returns:
        When only one pattern is given, returns a string (None if no match found).
        When more than one pattern are given, returns a list of strings ([] if no match found).
    """

    if len(patterns) == 1:
        pattern = patterns[0]
        match = re.search(pattern, text)
        if match:
            return match.group(1)
        else:
            return None
    else:
        ret = []
        for pattern in patterns:
            match = re.search(pattern, text)
            if match:
                ret.append(match.group(1))
        return ret


def parse_query_param(url, param):
    """Parses the query string of a URL and returns the value of a parameter.

    Args:
        url: A URL.
        param: A string representing the name of the parameter.

    Returns:
        The value of the parameter.
    """

    try:
        return parse.parse_qs(parse.urlparse(url).query)[param][0]
    except:
        return None


def unicodize(text):
    return re.sub(r'\\u([0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f][0-9A-Fa-f])',
                  lambda x: chr(int(x.group(0)[2:], 16)), text)


def r1_of(patterns, text):
    for p in patterns:
        x = r1(p, text)
        if x:
            return x


def unescape_html(html):
    import xml.sax.saxutils
    html = xml.sax.saxutils.unescape(html)
    html = re.sub(r'&#(\d+);', lambda x: unichr(int(x.group(1))), html)
    return html


def ungzip(s):
    from StringIO import StringIO
    import gzip
    buffer = StringIO(s)
    f = gzip.GzipFile(fileobj=buffer)
    return f.read()


def undeflate(s):
    import zlib
    return zlib.decompress(s, -zlib.MAX_WBITS)


def get_response(url, useragent=None, Referer=None, **kwargs):
    headers = {}
    if useragent:
        headers["User-Agent"] = useragent
    if Referer:
        headers["Referer"] = Referer
    headers.update(kwargs)
    request = urllib2.Request(url, headers=headers)
    response = urllib2.urlopen(request, timeout=10)
    data = response.read()
    if response.info().get('Content-Encoding') == 'gzip':
        data = ungzip(data)
    elif response.info().get('Content-Encoding') == 'deflate':
        data = undeflate(data)
    response.data = data
    return response


def get_html(url, encoding=None, useragent=None, Referer=None, **kwargs):
    content = get_response(
        url, useragent=useragent, Referer=Referer, **kwargs).data
    if encoding:
        content = content.decode(encoding)
    return content


def get_videojj_response(url, appkey_b64):

    headers = {
        "Host":
            "videojj.com",
        "User-Agent":
            "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:45.0) Gecko/20100101 Firefox/45.0",
        "Accept":
            "application/json, text/javascript, */*; q=0.01",
        "Accept-Language":
            "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
        "Accept-Encoding":
            "gzip, deflate",
        "Appkey":
            appkey_b64,
        "Referer":
            "http://www.baidu.com",
        "If-None-Match":
            'W/"207e-112555361"',
        "Connection":
            "keep-alive"
    }

    request = urllib2.Request(url, headers=headers)
    response = urllib2.urlopen(request)
    data = response.read()

    import json
    import StringIO
    import gzip
    try:
        data = StringIO.StringIO(data)
        gzipper = gzip.GzipFile(fileobj=data)
        html = gzipper.read()
        djson = json.loads(html)
        # for index, item in enumerate(djson["msg"]["segs"]["MP4-SD"]):
        #     print "MP4 - ", index, " : ",
        #     print item["url"]
    except:
        djson = json.loads(data)
        print djson

    return djson
    # import pdb
    # pdb.set_trace()
    # if response.info().get('Content-Encoding') == 'gzip':
    #     data = ungzip(data)
    # elif response.info().get('Content-Encoding') == 'deflate':
    #     data = undeflate(data)
    # response.data = data
    # return response


def get_decoded_html(url, useragent=None):
    response = get_response(url, useragent=useragent)
    data = response.data
    charset = r1(r'charset=([\w-]+)', response.headers['content-type'])
    if charset:
        return data.decode(charset)
    else:
        return data


def print_info(site_info, title, type, size):
    if type:
        type = type.lower()
    if type in ['3gp']:
        type = 'video/3gpp'
    elif type in ['asf', 'wmv']:
        type = 'video/x-ms-asf'
    elif type in ['flv', 'f4v']:
        type = 'video/x-flv'
    elif type in ['mkv']:
        type = 'video/x-matroska'
    elif type in ['mp3']:
        type = 'audio/mpeg'
    elif type in ['mp4']:
        type = 'video/mp4'
    elif type in ['mov']:
        type = 'video/quicktime'
    elif type in ['ts']:
        type = 'video/MP2T'
    elif type in ['webm']:
        type = 'video/webm'

    if type in ['video/3gpp']:
        type_info = "3GPP multimedia file (%s)" % type
    elif type in ['video/x-flv', 'video/f4v']:
        type_info = "Flash video (%s)" % type
    elif type in ['video/mp4', 'video/x-m4v']:
        type_info = "MPEG-4 video (%s)" % type
    elif type in ['video/MP2T']:
        type_info = "MPEG-2 transport stream (%s)" % type
    elif type in ['video/webm']:
        type_info = "WebM video (%s)" % type
    # elif type in ['video/ogg']:
    #    type_info = "Ogg video (%s)" % type
    elif type in ['video/quicktime']:
        type_info = "QuickTime video (%s)" % type
    elif type in ['video/x-matroska']:
        type_info = "Matroska video (%s)" % type
    # elif type in ['video/x-ms-wmv']:
    #    type_info = "Windows Media video (%s)" % type
    elif type in ['video/x-ms-asf']:
        type_info = "Advanced Systems Format (%s)" % type
    # elif type in ['video/mpeg']:
    #    type_info = "MPEG video (%s)" % type
    elif type in ['audio/mpeg']:
        type_info = "MP3 (%s)" % type
    else:
        type_info = "Unknown type (%s)" % type

    print("Video Site:", site_info)
    print("Title:     ", tr(title))
    print("Type:      ", type_info)
    print("Size:      ", round(size / 1048576,
                               2), "MiB (" + str(size) + " Bytes)")
    print()


def url_save(url, filepath, bar, refer=None):
    print "url: %s" % url
    print "refer: %s" % refer
    headers = {}
    if refer:
        headers['Referer'] = refer
    request = urllib2.Request(url, headers=headers)
    response = urllib2.urlopen(request)
    file_size = int(response.headers['content-length'])
    assert file_size
    if os.path.exists(filepath):
        if file_size == os.path.getsize(filepath):
            if bar:
                bar.done()
            print 'Skip %s: file already exists' % os.path.basename(filepath)
            return
        else:
            if bar:
                bar.done()
            print 'Overwriting', os.path.basename(filepath), '...'
    with open(filepath, 'wb') as output:
        received = 0
        while True:
            buffer = response.read(1024 * 256)
            if not buffer:
                break
            received += len(buffer)
            output.write(buffer)
            if bar:
                bar.update_received(len(buffer))
    assert received == file_size == os.path.getsize(
        filepath), '%s == %s == %s' % (received, file_size,
                                       os.path.getsize(filepath))


def url_size(url):
    request = urllib2.Request(url)
    request.get_method = lambda: 'HEAD'
    response = urllib2.urlopen(request)
    size = int(response.headers['content-length'])
    return size


def url_size(url):
    size = int(urllib2.urlopen(url).headers['content-length'])
    return size


def urls_size(urls):
    return sum(map(url_size, urls))


def url_info(url, faker=False):
    if faker:
        response = urllib2.urlopen(
            urllib2.Request(url, headers=fake_headers), None)
    else:
        response = urllib2.urlopen(urllib2.Request(url))

    headers = response.headers

    type = headers['content-type']
    mapping = {
        'video/3gpp': '3gp',
        'video/f4v': 'flv',
        'video/mp4': 'mp4',
        'video/MP2T': 'ts',
        'video/quicktime': 'mov',
        'video/webm': 'webm',
        'video/x-flv': 'flv',
        'video/x-ms-asf': 'asf',
        'audio/mpeg': 'mp3'
    }
    if type in mapping:
        ext = mapping[type]
    else:
        type = None
        if headers['content-disposition']:
            try:
                filename = parse.unquote(
                    r1(r'filename="?([^"]+)"?', headers['content-disposition']))
                if len(filename.split('.')) > 1:
                    ext = filename.split('.')[-1]
                else:
                    ext = None
            except:
                ext = None
        else:
            ext = None

    if headers.get('transfer-encoding') != 'chunked':
        size = int(headers['content-length'])
    else:
        size = None

    return type, ext, size


def parse_host(host):
    """Parses host name and port number from a string.
    """
    if re.match(r'^(\d+)$', host) is not None:
        return ("0.0.0.0", int(host))
    if re.match(r'^(\w+)://', host) is None:
        host = "//" + host
    o = parse.urlparse(host)
    hostname = o.hostname or "0.0.0.0"
    port = o.port or 0
    return (hostname, port)


def get_sogou_proxy():
    # sogou_proxy = parse_host(url)
    return None


def set_proxy(proxy):
    proxy_handler = request.ProxyHandler({
        'http': '%s:%s' % proxy,
        'https': '%s:%s' % proxy,
    })
    opener = request.build_opener(proxy_handler)
    request.install_opener(opener)


class SimpleProgressBar:

    def __init__(self, total_size, total_pieces=1):
        self.displayed = False
        self.total_size = total_size
        self.total_pieces = total_pieces
        self.current_piece = 1
        self.received = 0

    def update(self):
        self.displayed = True
        bar_size = 40
        percent = self.received * 100.0 / self.total_size
        if percent > 100:
            percent = 100.0
        bar_rate = 100.0 / bar_size
        dots = percent / bar_rate
        dots = int(dots)
        plus = percent / bar_rate - dots
        if plus > 0.8:
            plus = '='
        elif plus > 0.4:
            plus = '-'
        else:
            plus = ''
        bar = '=' * dots + plus
        bar = '{0:>3.0f}% [{1:<40}] {2}/{3}'.format(percent, bar,
                                                    self.current_piece,
                                                    self.total_pieces)
        sys.stdout.write('\r' + bar)
        sys.stdout.flush()

    def update_received(self, n):
        self.received += n
        self.update()

    def update_piece(self, n):
        self.current_piece = n

    def done(self):
        if self.displayed:
            print
            self.displayed = False


class PiecesProgressBar:

    def __init__(self, total_size, total_pieces=1):
        self.displayed = False
        self.total_size = total_size
        self.total_pieces = total_pieces
        self.current_piece = 1
        self.received = 0

    def update(self):
        self.displayed = True
        bar = '{0:>3}%[{1:<40}] {2}/{3}'.format('?', '?' * 40,
                                                self.current_piece,
                                                self.total_pieces)
        sys.stdout.write('\r' + bar)
        sys.stdout.flush()

    def update_received(self, n):
        self.received += n
        self.update()

    def update_piece(self, n):
        self.current_piece = n

    def done(self):
        if self.displayed:
            print
            self.displayed = False


class DummyProgressBar:

    def __init__(self, *args):
        pass

    def update_received(self, n):
        pass

    def update_piece(self, n):
        pass

    def done(self):
        pass


def escape_file_path(path):
    path = path.replace('/', '-')
    path = path.replace('\\', '-')
    path = path.replace('*', '-')
    path = path.replace('?', '-')
    return path


def download_urls(urls,
                  title,
                  ext,
                  total_size,
                  output_dir='.',
                  refer=None,
                  merge=True):
    assert urls
    assert ext in ('flv', 'mp4')
    if not total_size:
        try:
            total_size = urls_size(urls)
        except:
            import traceback
            import sys
            traceback.print_exc(file=sys.stdout)
            pass
    #title = to_native_string(title)
    #title = escape_file_path(title)
    #filename = '%s.%s' % (title, ext)
    filename = "test.mp4"
    filepath = os.path.join(output_dir, filename)
    if total_size:
        if os.path.exists(
                filepath) and os.path.getsize(filepath) >= total_size * 0.9:
            print 'Skip %s: file already exists' % filepath
            return
        bar = SimpleProgressBar(total_size, len(urls))
    else:
        bar = PiecesProgressBar(total_size, len(urls))
    if len(urls) == 1:
        url = urls[0]
        print 'Downloading %s ...' % filename
        url_save(url, filepath, bar, refer=refer)
        bar.done()
    else:
        flvs = []
        print 'Downloading %s.%s ...' % (title, ext)
        for i, url in enumerate(urls):
            filename = '%s[%02d].%s' % (title, i, ext)
            filepath = os.path.join(output_dir, filename)
            flvs.append(filepath)
            # print 'Downloading %s [%s/%s]...' % (filename, i+1, len(urls))
            bar.update_piece(i + 1)
            url_save(url, filepath, bar, refer=refer)
        bar.done()
        if not merge:
            return
        if ext == 'flv':
            from flv_join import concat_flvs
            concat_flvs(flvs, os.path.join(output_dir, title + '.flv'))
            for flv in flvs:
                os.remove(flv)
        elif ext == 'mp4':
            from mp4_join import concat_mp4s
            concat_mp4s(flvs, os.path.join(output_dir, title + '.mp4'))
            for flv in flvs:
                os.remove(flv)
        else:
            print "Can't join %s files" % ext


def playlist_not_supported(name):

    def f(*args, **kwargs):
        raise NotImplementedError('Play list is not supported for ' + name)

    return f


def script_main(script_name, download, download_playlist=None):
    if download_playlist:
        help = 'python %s.py [--playlist] [-c|--create-dir] [--no-merge] url ...' % script_name
        short_opts = 'hc'
        opts = ['help', 'playlist', 'create-dir', 'no-merge']
    else:
        help = 'python [--no-merge] %s.py url ...' % script_name
        short_opts = 'h'
        opts = ['help', 'no-merge']
    import sys
    import getopt
    try:
        opts, args = getopt.getopt(sys.argv[1:], short_opts, opts)
    except getopt.GetoptError, err:
        print help
        sys.exit(1)
    playlist = False
    create_dir = False
    merge = True
    for o, a in opts:
        if o in ('-h', '--help'):
            print help
            sys.exit()
        elif o in ('--playlist',):
            playlist = True
        elif o in ('-c', '--create-dir'):
            create_dir = True
        elif o in ('--no-merge'):
            merge = False
        else:
            print help
            sys.exit(1)
    if not args:
        print help
        sys.exit(1)

    for url in args:
        if playlist:
            download_playlist(url, create_dir=create_dir, merge=merge)
        else:
            download(url, merge=merge)


def union_dict(*objs):
    _keys = set(sum([obj.keys() for obj in objs], []))
    _total = {}
    for _key in _keys:
        _total[_key] = sum([obj.get(_key, 0) for obj in objs])
    return _total
