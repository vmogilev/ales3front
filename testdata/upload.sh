#!/bin/bash

if [ -z "$1" ]; then
    bucket=support-pub-dev
else
    bucket=$1
fi


../s3upload_diversified.sh ./CODEGUARDIAN_6-4-6_R11.txt.gz testdata-zzz CODEGUARDIAN_6-4-6.txt.gz ${bucket}
../s3upload_diversified.sh ./CODEGUARDIAN_6-4-6_R21.txt.gz testdata-zzz CODEGUARDIAN_6-4-6.txt.gz ${bucket}
../s3upload_diversified.sh ./CODEGUARDIAN_6-4-6_R31.txt.gz testdata-zzz CODEGUARDIAN_6-4-6.txt.gz ${bucket}
../s3upload_diversified.sh ./CODEGUARDIAN_6-4-6_R41.txt.gz testdata-zzz CODEGUARDIAN_6-4-6.txt.gz ${bucket}
../s3upload_diversified.sh ./CODEGUARDIAN_6-4-6_R51.txt.gz testdata-zzz CODEGUARDIAN_6-4-6.txt.gz ${bucket}
../s3upload_diversified.sh ./CODEGUARDIAN_6-7-1_R11.txt.gz testdata-zzz CODEGUARDIAN_6-7-1.txt.gz ${bucket}
../s3upload_diversified.sh ./CODEGUARDIAN_6-7-1_R21.txt.gz testdata-zzz CODEGUARDIAN_6-7-1.txt.gz ${bucket}
../s3upload_diversified.sh ./CODEGUARDIAN_6-7-1_R31.txt.gz testdata-zzz CODEGUARDIAN_6-7-1.txt.gz ${bucket}
../s3upload_diversified.sh ./CODEGUARDIAN_6-7-1_R41.txt.gz testdata-zzz CODEGUARDIAN_6-7-1.txt.gz ${bucket}
../s3upload_diversified.sh ./CODEGUARDIAN_6-7-1_R51.txt.gz testdata-zzz CODEGUARDIAN_6-7-1.txt.gz ${bucket}

