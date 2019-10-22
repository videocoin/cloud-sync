FROM golang:1.12.4 as builder
WORKDIR /go/src/github.com/videocoin/cloud-sync
COPY . .
RUN make build

FROM blitznote/debase:18.04

RUN apt-get update
RUN apt-get install -y ca-certificates
COPY --from=builder /go/src/github.com/videocoin/cloud-sync/bin/syncer /opt/videocoin/bin/syncer
CMD ["/opt/videocoin/bin/syncer"]

