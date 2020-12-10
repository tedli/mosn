ARG REGISTRY=192.168.1.52
ARG REPOSITORY=system_containers
ARG BUILDER_IMAGE_NAME=cicd-golang-build-golang
ARG BUILDER_IMAGE_TAG=1.15.6-alpine3.12
ARG BASE_IMAGE_NAME=sofastack-base
ARG BASE_IMAGE_TAG=0.19.0-202012101557

FROM ${REGISTRY}/${REPOSITORY}/${BUILDER_IMAGE_NAME}:${BUILDER_IMAGE_TAG} AS builder

ARG SKIP_UPX=0

COPY . /go/src/mosn.io/mosn
WORKDIR /go/src/mosn.io/mosn

RUN echo $(git rev-parse HEAD 2> /dev/null || true)$(if ! git diff-index --quiet HEAD; then echo -dirty; fi) \
> /usr/local/bin/version.txt && CGO_ENABLED=0 go build -v -mod=vendor -o /tmp/mosn.elf \
-ldflags="-X main.Version=$(cat /usr/local/bin/version.txt) -w -s" \
mosn.io/mosn/cmd/mosn/main 2>&1 && \
test 0 -eq ${SKIP_UPX} && upx -q --best --ultra-brute -o /usr/local/bin/mosn /tmp/mosn.elf || \
mv /tmp/mosn.elf /usr/local/bin/mosn

FROM ${REGISTRY}/${REPOSITORY}/${BASE_IMAGE_NAME}:${BASE_IMAGE_TAG}

COPY --from=builder /usr/local/bin/mosn /usr/local/bin/mosn
COPY --from=builder /usr/local/bin/version.txt /usr/local/bin/version.txt
