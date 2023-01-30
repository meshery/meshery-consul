FROM golang:1.17 as bd
ARG VERSION
ARG GIT_COMMITSHA
WORKDIR /github.com/layer5io/meshery-consul
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -ldflags="-w -s -X main.version=$VERSION -X main.gitsha=$GIT_COMMITSHA" -a -o /meshery-consul .
RUN find . -name "*.go" -type f -delete; mv consul /

FROM alpine:3.15
USER root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk --update add ca-certificates && \
    mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN chmod 775 /home/appuser/templates/oam/workloads/ \
 && chown -R appuser:appuser /home/appuser/templates/oam/workloads/

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
