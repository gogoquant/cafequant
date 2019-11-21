############################################
# version : xiyanxiyan10/quant
# quant 编译后台
# dist  编译前台
############################################

# 设置继承自alpine官方镜像
# depend https://github.com/xiyanxiyan10/docker-alpine-glibc
FROM registry.cn-hangzhou.aliyuncs.com/misaka/alpine-go:alpine-go

# 下面是一些创建者的基本信息
MAINTAINER xiyanxiyan10 (xiyanxiyan10@hotmail.com)
# set env
ENV LC_ALL C

# 设置环境变量，所有操作都是非交互式的
ENV DEBIAN_FRONTEND noninteractive

RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf
RUN echo "nameserver 114.114.114.114" >> /etc/resolv.conf

ENV LANG=C.UTF-8

#mkdir
RUN mkdir /opt/app -p 
RUN mkdir /opt/app/web -p 

#add file
COPY ./quant     /opt/app/quant
COPY ./dist      /opt/app/web/dist

VOLUME ["/data"]

EXPOSE 9876

WORKDIR /opt/app
RUN chmod 777 /opt/app/quant

CMD ["/opt/app/quant"]
