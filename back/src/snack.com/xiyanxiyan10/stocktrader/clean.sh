sudo docker images |grep none |awk '{print $3}'|xargs sudo docker rmi
