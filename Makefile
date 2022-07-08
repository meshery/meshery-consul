GOPATH = $(shell go env GOPATH)
BUILDER=buildx-multi-arch

GIT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_STRIPPED_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1` | cut -c 2-)
v ?= 1.17.8 # Default go version to be used

.PHONY: check
check: error
	golangci-lint run

.PHONY: check-clean-cache
check-clean-cache:
	golangci-lint cache clean

.PHONY: protoc-setup
protoc-setup:
	cd meshes
	wget https://raw.githubusercontent.com/layer5io/meshery/master/meshes/meshops.proto

.PHONY: proto
proto:	
	protoc -I meshes/ meshes/meshops.proto --go_out=plugins=grpc:./meshes/

.PHONY: docker
docker:
	docker build -t layer5/meshery-consul .

.PHONY: docker-run
docker-run:
	(docker rm -f meshery-consul) || true
	docker run --name meshery-consul -d \
	-p 10002:10002 \
	-e DEBUG=true \
	layer5/meshery-consul

.PHONY: run
run:
	go$(v) mod tidy -compat=1.17; \
	DEBUG=true go run main.go

run-force-dynamic-reg:
	FORCE_DYNAMIC_REG=true DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go

.PHONY: error
error:
	go run github.com/layer5io/meshkit/cmd/errorutil -d . analyze -i ./helpers -o ./helpers
