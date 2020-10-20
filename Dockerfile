FROM golang:1.13 as bd
WORKDIR /github.com/layer5io/meshery-consul
ADD . .
RUN GOPROXY=direct GOSUMDB=off go build -ldflags="-w -s" -a -o /meshery-consul .
RUN find . -name "*.go" -type f -delete; mv consul /

FROM alpine
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN apk --update add ca-certificates && \
    mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

USER appuser
RUN mkdir -p /home/appuser/.kube
WORKDIR /home/appuser
COPY --from=bd /meshery-consul /home/appuser
COPY --from=bd /consul /home/appuser/consul
CMD ./meshery-consul
