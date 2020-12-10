REGISTRY ?= 192.168.1.52
REPOSITORY ?= system_containers
BUILDER_IMAGE_NAME ?= cicd-golang-build-golang
BUILDER_IMAGE_TAG ?= 1.15.6-alpine3.12
BASE_IMAGE_NAME ?= sofastack-base
BASE_IMAGE_TAG ?= 0.19.0-202012101557
SKIP_UPX ?= 0
IMAGE_NAME ?= proxyv2
IMAGE_TAG ?= 0.19.0-mosn-$(shell date '+%Y%m%d%H%M')

.PHONY: build

build:
	docker build --no-cache --build-arg REGISTRY=${REGISTRY} \
--build-arg REPOSITORY=${REPOSITORY} \
--build-arg BUILDER_IMAGE_NAME=${BUILDER_IMAGE_NAME} \
--build-arg BUILDER_IMAGE_TAG=${BUILDER_IMAGE_TAG} \
--build-arg BASE_IMAGE_NAME=${BASE_IMAGE_NAME} \
--build-arg BASE_IMAGE_TAG=${BASE_IMAGE_TAG} \
--build-arg SKIP_UPX=${SKIP_UPX} \
-t ${REGISTRY}/${REPOSITORY}/${IMAGE_NAME}:${IMAGE_TAG} -f ci.dockerfile .
