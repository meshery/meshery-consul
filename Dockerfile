FROM golang:1.19 as builder

ARG VERSION
ARG GIT_COMMITSHA
WORKDIR /build
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN GOPROXY=https://proxy.golang.org,direct go mod download
# Copy the go source
COPY main.go main.go
COPY internal/ internal/
COPY consul/ consul/
COPY build/ build/

RUN GOPROXY=https://proxy.golang.org,direct CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -a -o meshery-consul main.go

FROM alpine:3.15
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk --update add ca-certificates && \
    mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

USER appuser
ENV SERVICE_ADDR="meshery-consul"
ENV MESHERY_SERVER="http://meshery:9081"
RUN mkdir -p /home/appuser/.kube
RUN mkdir -p /home/appuser/.meshery
WORKDIR /home/appuser
COPY templates/ ./templates
COPY --from=builder /build/meshery-consul /home/appuser
COPY consul /home/appuser/consul
CMD ./meshery-consul
