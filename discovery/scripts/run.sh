#!/bin/bash

kubectl delete pod demo
kubectl run --rm -i demo --image=discovery --image-pull-policy='Never'