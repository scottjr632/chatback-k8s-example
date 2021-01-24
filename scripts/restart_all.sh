#! /bin/bash

./scripts/cleanup_k8s.sh
./scripts/build_docker.sh
./scripts/apply_all.sh
