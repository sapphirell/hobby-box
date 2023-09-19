#!/bin/bash

print() {
    echo -e "\033[33mINFO:${1}\033[0m"
}

cmd() {
    echo -e "\033[36mCOMMAND: ${1} \033[0m"
    $1
}

docker rmi -f fantuanpu:v1
docker build -t fantuanpu:v1 .

p=`docker ps | grep fantuanpu | grep -v grep| awk '{print $1}'`
if [ $p ]; then
    docker stop $p
fi
docker run -p 80:8080 --link mysql8.0.32:localdb --link some-redis:localredis -d fantuanpu:v1