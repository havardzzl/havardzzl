BINARY="havard"
TIME=$(shell date +"%m-%d-%H")
LONGVER=$(shell git log -1 --pretty=format:"%H")
VERSION=$(shell echo ${LONGVER}|awk '{print substr($$0,1,5)}')
IMAGE=harbor.weizhipin.com\/tcloud\/havard:${TIME}-${VERSION}

# Build the docker image
docker-build:
	@echo ${IMAGE}
	docker build -t ${IMAGE} ./build

init:
	set -e
	rm -rf ./build
	mkdir ./build
	cp ./Dockerfile ./build/Dockerfile
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/${BINARY}

build: init docker-build