FROM golang:1.14 as builder

WORKDIR /go/src/github.com/videocoin/cloud-sync
COPY . .

RUN make build


FROM bitnami/minideb:stretch

RUN apt-get update
RUN apt-get install -y ca-certificates curl

COPY --from=builder /go/src/github.com/videocoin/cloud-sync/bin/syncer /syncer

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 && \
   curl -L -k https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 --output /bin/grpc_health_probe && chmod +x /bin/grpc_health_probe

CMD ["syncer"]
