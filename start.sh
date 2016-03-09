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

nohup $TOP/bin/${APP} \
           -cfKeyID=${cfKeyID} \
           -cfKeyFile=${cfKeyFile} \
           -cdnHost=${cdnHost} \
           -cfExpHours=${cfExpHours} \
           -htmlPath=${htmlPath} \
           -debug=${debug} \
           -httpPort=${httpPort} \
           -rootToken=${rootToken} >> ${LOG}/server.log 2>&1 </dev/null &
