#!/bin/bash
docker run --name mongo-db -d mongo
docker build -t things-api .
docker run --name iothings -p 127.0.0.1:4000:4000 --link mongo-db:mongo -d things-api
