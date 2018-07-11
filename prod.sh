#!/bin/bash
openssl genrsa -out key.rsa 2048
openssl base64 -in key.rsa -out key64.rsa
docker run --name mongo-db -d mongo
docker build -t iothings-api .
docker run --name iothings -p 127.0.0.1:4000:4000 --link mongo-db:mongo -d iothings-api
