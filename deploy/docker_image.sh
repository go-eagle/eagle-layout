#!/bin/bash

NAMESPACE="go-microservices"
REGISTRY="registry.cn-hangzhou.aliyuncs.com"
IMAGE_NAME=$1
TAG=$2

if [[ -z "$IMAGE_NAME" ]] || [[ -z "$TAG" ]]; then
  echo "tag is empty"
  echo "usage: sh $0 <SERVICE_NAME> <TAG>"
  exit 1
fi

# Replace all occurrences of eagle-layout with the new name
DOCKERFILE_PATH="deploy/docker/Dockerfile"
sed -i "" "s/eagle-layout/${IMAGE_NAME}/g" ${DOCKERFILE_PATH}
echo "Successfully replaced 'eagle-layout' with '${IMAGE_NAME}' in ${DOCKERFILE_PATH}"

# build image
echo "> build docker image"
docker build -t ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${TAG} -t ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:latest -f deploy/docker/Dockerfile .

# docker push new-repo:tagname
echo "> push docker image to remote hub"
docker push ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${TAG}

echo "Done. push docker image success."