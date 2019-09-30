package main

import (
	"context"
	"crypto/rand"

	v1 "github.com/videocoin/cloud-api/syncer/v1"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("127.0.0.1:5010", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := v1.NewSyncerServiceClient(conn)
	nByte := 5
	data := make([]byte, nByte*1024*1024)

	rand.Read(data)

	ctx := context.Background()

	for i := 0; i < 1; i++ {
		stream, err := client.Sync(ctx)
		if err != nil {
			panic(err)
		}

		defer stream.CloseSend()

		err = stream.Send(&v1.SyncRequest{
			SyncOneof: &v1.SyncRequest_Meta{
				Meta: &v1.Metadata{
					Path: "data",
				},
			},
		})

		if err != nil {
			panic(err)
		}

		err = stream.Send(&v1.SyncRequest{
			SyncOneof: &v1.SyncRequest_Data{
				Data: data,
			},
		})

		if err != nil {
			panic(err)
		}

		_, err = stream.CloseAndRecv()
		if err != nil {
			panic(err)
		}
	}
}
