#!/bin/bash

docker stop dd-agent
docker stop ales3front-dev
docker rm $(docker ps -a -q)

docker run -d --name dd-agent -h `hostname` -p 8125:8125/udp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v /proc/:/host/proc/:ro \
        -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
        -e API_KEY=5db2d30a51b1a50b77fa3741cc9bc3a2 \
        datadog/docker-dd-agent:latest

docker run -d --publish 80:8080 \
           --name ales3front-dev \
           --link dd-agent:dd-agent \
           -v $HOME/.aws/:/root/.aws/ \
           -v $HOME/.ssh/:/root/.ssh/ \
           -v '/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt' \
           639518868771.dkr.ecr.us-east-1.amazonaws.com/alesi2/ales3front:dev \
           -cfKeyFile=/root/.ssh/pk-APKAIOYPWYODYKDZBPQA.pem \
           -cfKeyID=APKAIOYPWYODYKDZBPQA \
           -debug=true \
           -htmlPath=/app/html \
           -ddAgent=dd-agent:8125 \
           -ddPrefix=devdoc

docker ps -a
