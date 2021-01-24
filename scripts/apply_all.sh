#!/bin/bash

USE_REMOTE='false'

print_usage() {
  printf "Usage: \n\tr - use remote?"
}

while getopts 'r' flag; do
  case "${flag}" in
    r) USE_REMOTE='true' ;;
    *) print_usage
       exit 1 ;;
  esac
done

# apply all of the monitoring
kubectl apply -f prometheus/clusterRole.yml
kubectl apply -f prometheus/config-map.yaml
kubectl apply -f prometheus/prometheus-deployment.yml

# create chatback services
if [[ $USE_REMOTE == 'true' ]]; then
    kubectl apply -f chatback-remote.yml
else
    kubectl apply -f chatback.yml
fi