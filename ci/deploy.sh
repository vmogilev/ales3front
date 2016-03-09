#!/bin/bash

if [ $1"x" != "YESx" ]; then
    echo "this script is designed to run from hooks/post-receive in staging bare git repo"
    echo "by placing post-receive in hooks/post-receive"
    echo "see: http://www.dbatoolz.com/t/continuous-integration-golang-github.html"
    exit
fi

# ----------- EDIT BEGIN ----------- #

BASE_TOP=/data; export BASE_TOP
APP_TOP=${BASE_TOP}/app; export APP_TOP
GO=/usr/local/go/bin/go; export GO

## go [get|install] ${SRC_PATH}/${APP_NAME}
SRC_PATH=bitbucket.org/alesi2; export SRC_PATH
APP_NAME=ales3front; export APP_NAME

## local http root directory served by go http - ${APP_TOP}/${WWW_PATH}
## for / directory use /root:
##      site.com        -> site.com/root
##      blog.site.com   -> blog.site.com/root
##      site.com/mnt    -> site.com/mnt
WWW_PATH=dev.alcalcs.com/root; export WWW_PATH

## local bare git repo path - ${SRC_NAME}/${APP_NAME}.git
SRC_NAME=alesi2; export SRC_NAME

# ----------- EDIT END ----------- #


GOPATH=${BASE_TOP}/dev/golang; export GOPATH
SOURCE=${GOPATH}/src/${SRC_PATH}/${APP_NAME}; export SOURCE
TARGET=${APP_TOP}/${WWW_PATH}; export TARGET
GIT_DIR=${BASE_TOP}/stage/git/${SRC_NAME}/${APP_NAME}.git; export GIT_DIR

## pre-creating SOURCE DIR solves the issue with:
##  "remote: fatal: This operation must be run in a work tree"
mkdir -p ${SOURCE}
mkdir -p ${TARGET}

## GIT_WORK_TREE=${SOURCE} git checkout -f


## do not prefix go get with GIT_WORK_TREE - it causes the following errors:
##  remote: # cd .; git clone https://github.com/rigingo/dlog /data/app/dev/golang/src/github.com/rigingo/dlog
##  remote: fatal: working tree '/data/app/dev/golang/src/github.com/rigingo/dop' already exists.
##
##GIT_WORK_TREE=${SOURCE} $GO get github.com/rigingo/dop
##GIT_WORK_TREE=${SOURCE} $GO install github.com/rigingo/dop

unset GOBIN
unset GIT_DIR
$GO get ${SRC_PATH}/${APP_NAME}
$GO install ${SRC_PATH}/${APP_NAME}

if [ $? -gt 0 ]; then
    echo "ERROR: compiling ${APP_NAME} - exiting!"
    exit 1
fi

sudo setcap 'cap_net_bind_service=+ep' $GOPATH/bin/${APP_NAME}


# ----------- DEPLOY BEGIN ----------- #

cp -pr ${SOURCE}/html     ${TARGET}/
cp -pr ${SOURCE}/conf     ${TARGET}/
cp -p ${SOURCE}/*.sh        ${TARGET}/
chmod +x ${TARGET}/*.sh

mkdir -p ${TARGET}/bin
## this fails with
## cp: setting attribute ‘security.capability’ for ‘security.capability’: Operation not permitted
## cp --preserve=all $GOPATH/bin/${APP_NAME} ${TARGET}/bin/
## so we have to use setcap again
cp -p $GOPATH/bin/${APP_NAME} ${TARGET}/bin/
sudo setcap 'cap_net_bind_service=+ep' ${TARGET}/bin/${APP_NAME}

. ${TARGET}/conf/${APP_NAME}.env
${TARGET}/stop.sh >> ${TARGET}/log/stop.log 2>&1 </dev/null
${TARGET}/start.sh >> ${TARGET}/log/start.log 2>&1 </dev/null

# ----------- DEPLOY END ----------- #
