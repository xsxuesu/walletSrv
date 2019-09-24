#!/bin/bash
echo "setup hot wallet"
## check and make dir
if [ ! -d "/var/wallet/hotservice" ]; then
  mkdir -p /var/wallet/hotservice
  mkdir -p /var/wallet/hotservice/keystore/eth/
fi
if [ ! -a "/var/wallet/hotservice/log" ]; then
  touch /var/wallet/hotservice/log
fi

CurPath=$(pwd)
echo ${CurPath}

rm -rf ${CurPath}/target/hotservice

## copy config.json
cp ${CurPath}/hotWallet/config/config.json  /var/wallet/hotservice/config.json

## cd hotwallet go build
cd ${CurPath}/hotWallet &&  go build -o ${CurPath}/target/hotservice  main.go

## copy hotservice
ps -ef |grep hotservice |awk '{print $2}'|xargs kill -9
cp ${CurPath}/target/hotservice  /var/wallet/hotservice/hotservice

## start hot service
cd /var/wallet/hotservice && nohup /var/wallet/hotservice/hotservice -f /var/wallet/hotservice/config.json >> /var/wallet/hotservice/log 2>&1 &


