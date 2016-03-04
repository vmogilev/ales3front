#!/bin/bash

err() {
    echo "ERROR: ${1} - exiting"
    exit 1
}

APP=ales3front

if [ -f "$GOPATH/bin/${APP}" ]; then
    echo "OK: found $GOPATH/bin/${APP}"
else
    err "application binary is missing in $GOBIN/${APP}"
fi

nohup $GOPATH/bin/${APP} \
           -cfKeyID=${cfKeyID} \
           -cfKeyFile=${cfKeyFile} \
           -cdnHost=${cdnHost} \
           -cfExpHours=${cfExpHours} \
           -htmlPath=${htmlPath} \
           -debug=${debug} \
           -httpPort=${httpPort} \
           -rootToken=${rootToken} >> ${DOPROOT}/server.log 2>&1 </dev/null &

