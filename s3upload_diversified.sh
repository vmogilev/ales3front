#!/bin/bash

FILE=$1
SDIR=$2
ROOT=$3
BUCK=$4

usage() {

    name=`basename $0`
    echo "${name} localfile dir_under_release rootname bucket "
    exit 1
}

if [ -z "$FILE" ] || [ -z "$SDIR" ] || [ -z "$ROOT" ] || [ -z "$BUCK" ]; then
    usage;
fi

if [ ! -f $FILE ]; then
    usage;
fi

key=`basename $FILE`

echo "
    FILE=${FILE}
    SDIR=${SDIR}
    ROOT=${ROOT}
    BUCK=${BUCK}
    key=${key}
"

aws s3api put-object --bucket ${BUCK} --key uploads/release/${SDIR}/${key} \
                     --body ${FILE} \
                     --content-disposition filename\=\"${ROOT}\"
