#!/bin/bash

IMAGE_PREFIX="chatback"

docker build -t $IMAGE_PREFIX-db ./db
docker build -t $IMAGE_PREFIX-client ./client
docker build -t $IMAGE_PREFIX-server ./server
docker build -t $IMAGE_PREFIX-broker ./broker
