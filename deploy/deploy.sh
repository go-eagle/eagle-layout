#!/bin/bash

source deploy/docker_image.sh

# deploy k8s deployment
echo "> deploy k8s deployment"
kubectl apply -f deploy/k8s/go-deployment.yaml

# deploy k8s service
echo "> deploy k8s service"
kubectl apply -f deploy/k8s/go-service.yaml

echo "Done. deploy success."