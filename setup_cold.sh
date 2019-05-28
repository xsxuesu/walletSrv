#!/bin/bash
## check and make dir
if [ ! -d "/var/wallet/coldservice" ]; then
  mkdir -p /var/wallet/coldservice
fi
if [ ! -a "/var/wallet/coldservice/log" ]; then
  touch /var/wallet/coldservice/log
fi

CurPath=$(pwd)
echo ${CurPath}

rm -rf ${CurPath}/target/coldservice

## copy config.json
cp ${CurPath}/coldWallet/config/config.json  /var/wallet/coldservice/config.json

## cd coldwallet go build
cd ${CurPath}/coldWallet  &&  go build -o ${CurPath}/target/coldservice  main.go

## copy coldservice
ps -ef |grep coldservice |awk '{print $2}'|xargs kill -9
cp ${CurPath}/target/coldservice  /var/wallet/coldservice/coldservice

## start cold service
cd /var/wallet/coldservice && nohup /var/wallet/coldservice/coldservice -f /var/wallet/coldservice/config.json >> /var/wallet/coldservice/log 2>&1 &









