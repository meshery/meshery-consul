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

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
ENV DISTRO="debian"
ENV SERVICE_ADDR="meshery-consul"
ENV MESHERY_SERVER="http://meshery:9081"
WORKDIR $HOME/.meshery
COPY templates/ ./templates
COPY --from=builder /build/meshery-consul .
USER nonroot:nonroot

CMD ["./meshery-consul"]

