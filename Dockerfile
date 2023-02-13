FROM golang:1.19 as builder

ARG VERSION
ARG GIT_COMMITSHA
WORKDIR /build
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
# Copy the go source
COPY main.go main.go
COPY internal/ internal/
COPY consul/ consul/
COPY build/ build/

RUN GOPROXY=direct,https://proxy.golang.org CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -a -o meshery-consul main.go

FROM gcr.io/distroless/nodejs:16
ENV DISTRO="debian"
ENV SERVICE_ADDR="meshery-consul"
ENV MESHERY_SERVER="http://meshery:9081"
COPY templates/ ./templates
WORKDIR /
COPY --from=builder /build/meshery-consul .
ENTRYPOINT ["/meshery-consul"]
