############################################
# version : xiyanxiyan10/quant
# quant 编译后台
# dist  编译前台
############################################

# 设置继承自alpine官方镜像
FROM alpine

# 下面是一些创建者的基本信息
MAINTAINER xiyanxiyan10 (xiyanxiyan10@hotmail.com)

# set env
ENV LC_ALL C

# 设置环境变量，所有操作都是非交互式的
ENV DEBIAN_FRONTEND noninteractive

RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf
RUN echo "nameserver 114.114.114.114" >> /etc/resolv.conf
RUN yum install -y sqlite 

#mkdir
RUN mkdir /opt/app -p 

#add file
COPY /tmp/quant     /opt/app/quant
COPY /tmp/dist      /opt/app/dist

VOLUME ["/data"]

EXPOSE 9876

WORKDIR /opt/app

CMD ["/opt/app/quant"]
