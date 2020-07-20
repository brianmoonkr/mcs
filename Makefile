all: usage

BINARY=mcs
ORG=cojam
NAME=mcs-dev
VERSION=0.1.1
BUILD=${VERSION}.2
DIST=ubuntu-18.04
IMAGE=${ORG}/${NAME}:${BUILD}-${DIST}

usage: 
	@echo ""
	@echo "usage: make [edit|build|push]"
	@echo ""

build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BINARY}
	docker build -t ${IMAGE} -f Dockerfile.dev .

push:
	docker push ${IMAGE}
