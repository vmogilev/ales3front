#!/bin/bash

docker stop dev01_docker
docker rm $(docker ps -a -q)
docker run --publish 8080:8080 \
           --name ales3front_dev \
           --rm \
           -v $HOME/.aws/:/.aws/ \
           -v $HOME/.ssh/:/.ssh/ \
           -v '/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt' \
           639518868771.dkr.ecr.us-east-1.amazonaws.com/alesi2/ales3front:latest \
           -cfKeyFile=/.ssh/pk-APKAIOYPWYODYKDZBPQA.pem \
           -cfKeyID=APKAIOYPWYODYKDZBPQA \
           -debug=true \
           -htmlPath=/app/html \
           -ddAgent=localhost:8125 \
           -ddPrefix=dev01_docker

