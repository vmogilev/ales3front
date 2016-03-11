#!/bin/bash

cd /data/app/dev.alcalcs.com/root/
./stop.sh
sudo setcap -r ./bin/ales3front

docker build -t alesi2/ales3front .
docker tag alesi2/ales3front:latest 639518868771.dkr.ecr.us-east-1.amazonaws.com/alesi2/ales3front:latest

aws ecr get-login --region us-east-1 > /tmp/dockerlogin.sh
sh /tmp/dockerlogin.sh
rm -f /tmp/dockerlogin.sh

docker push 639518868771.dkr.ecr.us-east-1.amazonaws.com/alesi2/ales3front:latest

sudo setcap 'cap_net_bind_service=+ep' ./bin/ales3front
. ./conf/ales3front.env
./start.sh >> ./log/start.log 2>&1 </dev/null
