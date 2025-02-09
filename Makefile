# Copyright Meshery Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

include build/Makefile.core.mk
include build/Makefile.show-help.mk

#-----------------------------------------------------------------------------
# Environment Setup
#-----------------------------------------------------------------------------
BUILDER=buildx-multi-arch
ADAPTER=consul

#-----------------------------------------------------------------------------
# Docker-based Builds
#-----------------------------------------------------------------------------
.PHONY: docker docker-run lint error test run run-force-dynamic-reg


## Lint check Golang
lint:
	golangci-lint run -c .golangci.yml -v ./...

tidy:
	go mod tidy

verify:
	go mod verify

gobuild:
	go build -o bin/$(ADAPTER) main.go


## Build Adapter container image with "edge-latest" tag
docker:
	DOCKER_BUILDKIT=1 docker build -t meshery/meshery-$(ADAPTER):$(RELEASE_CHANNEL)-latest .

## Run Adapter container with "edge-latest" tag
docker-run:
	(docker rm -f meshery-$(ADAPTER)) || true
	docker run --name meshery-$(ADAPTER) -d \
	-p 10002:10002 \
	-e DEBUG=true \
	meshery/meshery-$(ADAPTER):$(RELEASE_CHANNEL)-latest

## Build and run Adapter locally
run:
	go mod tidy; \
	DEBUG=true GOPROXY=direct GOSUMDB=off go run main.go
