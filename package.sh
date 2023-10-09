#!/bin/bash

print() {
    echo -e "\033[33mINFO:${1}\033[0m"
}

cmd() {
    echo -e "\033[36mCOMMAND: ${1} \033[0m"
    $1
}
git pull
#清理无用镜像
#停止
docker stop $(docker ps -a | grep "Exited" | awk '{print $1 }')
#删除容器
docker rm $(docker ps -a | grep "Exited" | awk '{print $1 }')
#删除镜像
docker rmi $(docker images | grep "none" | awk '{print $3}')


#打包
docker rmi -f fantuanpu:v1
docker build -t fantuanpu:v1 .

p=`docker ps -a | grep fantuanpu | grep -v grep| awk '{print $1}'`
if [[ $p ]]; then
    docker stop $p
fi

p=`docker ps -a | grep main | grep -v grep| awk '{print $1}'`
if [[ $p ]]; then
    docker stop $p
fi
docker run -p 80:8080 -p 443:443 --link mysql8.0.32:localdb --link some-redis:localredis -d fantuanpu:v1