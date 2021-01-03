#!/bin/bash

# apply all of the monitoring
kubectl apply -f prometheus/clusterRole.yml
kubectl apply -f prometheus/config-map.yaml
kubectl apply -f prometheus/prometheus-deployment.yml

# create chatback services
kubectl apply -f chatback.yml