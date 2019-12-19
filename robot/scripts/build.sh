#!/usr/bin/env bash
version=$1
if [ ! -n "$version" ] ;then
  version=0.0.1
fi
command -v go >/dev/null 2>&1
existsGO=$?
command -v docker >/dev/null 2>&1
existsDocker=$?
if [ $existsGO -ne 0 ] ;then
   echo "please install golang"
   exit 0
fi
env GOOS=linux GOARCH=amd64 go build -o ./build/apigateway -v ./cmd/example/...
if [ $existsDocker -ne 0 ];then
    echo "please install docker"
   exit 0
fi
docker build -f ./build/package/Dockerfile   -t reg.changhong.io/matrix-mce/apigateway:${version} .
rm -rf ./build/apigateway
docker push reg.changhong.io/matrix-mce/apigateway:${version}
exit 0
