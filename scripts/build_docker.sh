#!/bin/bash

USER='scottjr632'
USE_REMOTE='false'
IMAGE_PREFIX="chatback"

print_usage() {
  printf "Usage: \n\tr - use remote?\n\tu - user for docker hub\n\tp - image prefix\n"
}

while getopts 'ru:p:' flag; do
  case "${flag}" in
    u) USER='true' ;;
    r) USE_REMOTE='true' ;;
    p) IMAGE_PREFIX="${OPTARG}" ;;
    *) print_usage
       exit 1 ;;
  esac
done

if [[ $USE_REMOTE == 'true' ]]; then
    IMAGE_PREFIX="$USER/$IMAGE_PREFIX"
fi

docker build -t $IMAGE_PREFIX-db ./db
docker build -t $IMAGE_PREFIX-client ./client
docker build -t $IMAGE_PREFIX-server ./server
docker build -t $IMAGE_PREFIX-broker ./broker

if [[ $USE_REMOTE == 'true' ]]; then
    printf "pushing images to remote..."
    docker push $IMAGE_PREFIX-db
    docker push $IMAGE_PREFIX-client
    docker push $IMAGE_PREFIX-server
    docker push $IMAGE_PREFIX-broker
fi
