FROM golang:1.12.4 as builder
WORKDIR /go/src/github.com/videocoin/cloud-sync
COPY . .
RUN make build

FROM bitnami/minideb:jessie
RUN apt-get update && apt-get -y install ca-certificates
COPY --from=builder /go/src/github.com/videocoin/cloud-sync/bin/syncer /opt/videocoin/bin/syncer
CMD ["/opt/videocoin/bin/syncer"]

