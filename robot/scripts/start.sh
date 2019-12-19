#!/bin/sh
echo "10.66.2.81 registry.changhong.com" >> /etc/hosts
#/deploy/apigateway
nohup ping 127.0.0.1 &
