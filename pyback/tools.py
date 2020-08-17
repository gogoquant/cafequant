#!/usr/bin/python
# -*- coding: utf-8 -*-
"""
所有工具实现在此

@file:tools.py
@modul:tools
@author:xiaolin@idealsee.cn
@date:2015-03-21
"""
import os
import re
import os.path

from optparse import OptionParser

# 扫描模板，生成中文


class GenZhCNFile():

    def __init__(self, template_dir, out_file):
        self.template_dir = template_dir
        self.out_file = out_file

    def run(self):
        zhcndict = self.get_template_zhcndict()
        self.write_zhchfile(zhcndict)

    def get_template_zhcndict(self):
        template_zhcndict = {}
        for f in os.listdir(self.template_dir):
            template_zhcndict[f] = []

            file_path = os.path.join(self.template_dir, f)
            if not os.path.exists(file_path):
                continue

            line_list = []
            with open(file_path, "r") as file_data:
                line_list = file_data.readlines()
            if not line_list:
                continue

            for line in line_list:
                if line.startswith("<!--"):
                    continue
                line = line.decode("utf-8")
                zhcnlist = re.findall(
                    u'[\u4e00-\u9fa5|\d][\u4e00-\u9fa5|\w|\°|\&\，|\；|\。|\、|\/]+',
                    line)
                template_zhcndict[f] += zhcnlist
        return template_zhcndict

    def write_zhchfile(self, zhcndict):
        out_fd = open(self.out_file, "w")
        lines = []
        for f, ls in zhcndict.items():
            lines.append(u"#%s" % f)
            ls = list(set(ls))
            for zhcn in ls:
                has_zh = re.findall(u'[\u4e00-\u9fa5]+', zhcn)
                if not has_zh:
                    ls.remove(zhcn)
            lines += ls
        l2 = list(set(lines))
        l2.sort(key=lines.index)
        ws = "\n".join(l2)
        out_fd.write(ws.encode("utf-8"))
        out_fd.close()


if __name__ == '__main__':
    usage = "tools.py method [options]"
    methods = {"gen": u"扫描模板，生成中文 -d 输入模板目录 -o 输出文件"}
    method_helps = []
    for method, help in methods.items():
        method_helps.append("%s:%s" % (method, help))
    usage += "\n" + ("\n").join(method_helps)

    parser = OptionParser(usage=usage)
    parser.add_option(
        "-f",
        "--file",
        dest="filename",
        help="write report to FILE",
        metavar="FILE")
    parser.add_option("-d", "--dir", dest="dir", help="open dir")
    parser.add_option("-o", "--out", dest="out", help="out file")
    parser.add_option(
        "-q",
        "--quiet",
        action="store_false",
        dest="verbose",
        default=True,
        help="don't print status messages to stdout")
    (options, args) = parser.parse_args()
    assert len(args)

    arg = args[0]

    if args[0] == "gen":
        if not options.dir or not options.out:
            parser.print_help()
            exit(1)
        g = GenZhCNFile(options.dir, options.out)
        g.run()
