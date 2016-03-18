#!/bin/bash

err() {
    echo "ERROR: ${1} - exiting"
    exit 1
}

APP=ales3front

if [ -f "$TOP/bin/${APP}" ]; then
    echo "OK: found $TOP/bin/${APP}"
else
    err "application binary is missing in $TOP/bin/${APP}"
fi

## To start Data Dog Agent Docker Container:
## -----------------------------------------
## docker stop dd-agent
## docker rm $(docker ps -a -q)
## docker run -d --name dd-agent -h `hostname` -p 8125:8125/udp -v /var/run/docker.sock:/var/run/docker.sock -v /proc/:/host/proc/:ro -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro -e API_KEY=5db2d30a51b1a50b77fa3741cc9bc3a2 datadog/docker-dd-agent:latest
## -----------------------------------------

nohup $TOP/bin/${APP} \
           -cfKeyID=${cfKeyID} \
           -cfKeyFile=${cfKeyFile} \
           -cdnHost=${cdnHost} \
           -cfExpHours=${cfExpHours} \
           -htmlPath=${htmlPath} \
           -debug=${debug} \
           -httpPort=${httpPort} \
           -rootToken=${rootToken} \
           -ddAgent=localhost:8125 \
           -ddPrefix=dev01 >> ${LOG}/server.log 2>&1 </dev/null &
