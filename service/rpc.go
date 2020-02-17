package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	protoempty "github.com/gogo/protobuf/types"
	"github.com/grafov/m3u8"
	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/syncer/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-sync/eventbus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type RpcServerOptions struct {
	Addr   string
	Bucket string
	Logger *logrus.Entry
	DS     *Datastore
	EB     *eventbus.EventBus
}

type RpcServer struct {
	addr   string
	bucket string
	grpc   *grpc.Server
	listen net.Listener
	logger *logrus.Entry
	ds     *Datastore
	eb     *eventbus.EventBus
	gscli  *storage.Client
	bh     *storage.BucketHandle
}

func NewRpcServer(opts *RpcServerOptions) (*RpcServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcOpts = append(grpcOpts, grpc.MaxRecvMsgSize(1024*1024*1024))

	grpcServer := grpc.NewServer(grpcOpts...)
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	gscli, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	bh := gscli.Bucket(opts.Bucket)

	rpcServer := &RpcServer{
		addr:   opts.Addr,
		bucket: opts.Bucket,
		grpc:   grpcServer,
		listen: listen,
		logger: opts.Logger,
		ds:     opts.DS,
		eb:     opts.EB,
		bh:     bh,
		gscli:  gscli,
	}

	v1.RegisterSyncerServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *RpcServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *RpcServer) Sync(ctx context.Context, req *v1.SyncRequest) (*protoempty.Empty, error) {
	return &protoempty.Empty{}, nil
}

func (s *RpcServer) uploadSegment(ctx context.Context, streamID string, segmentNum int, ct string, data []byte) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	objectName := fmt.Sprintf("%s/%d.ts", streamID, segmentNum)

	logger := s.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"segment_num": segmentNum,
		"bucket":      s.bucket,
		"object_name": objectName,
	})

	logger.Info("uploading segment")

	gctx := context.Background()
	cli, err := storage.NewClient(gctx)
	if err != nil {
		return nil, nil, err
	}
	defer cli.Close()

	bh := cli.Bucket(s.bucket)
	if _, err := bh.Attrs(ctx); err != nil {
		return nil, nil, err
	}

	obj := bh.Object(objectName)
	w := obj.NewWriter(ctx)
	w.CacheControl = "no-cache"

	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		return nil, nil, err
	}

	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(ctx)

	logger.Info("segment has been uploaded successfully")

	return obj, attrs, err
}

func (s *RpcServer) generateAndUploadLiveMasterPlaylist(ctx context.Context, streamID string, segmentNum int) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	objectName := fmt.Sprintf("%s/index.m3u8", streamID)

	logger := s.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"segment_num": segmentNum,
		"bucket":      s.bucket,
		"object_name": objectName,
	})

	logger.Info("generating live master playlist")

	p, err := m3u8.NewMediaPlaylist(uint(segmentNum), uint(segmentNum))
	if err != nil {
		return nil, nil, err
	}

	for num := 1; num <= segmentNum; num++ {
		err := p.Append(fmt.Sprintf("%d.ts", num), 10, "")
		if err != nil {
			return nil, nil, err
		}
	}

	data := p.Encode().Bytes()

	logger.Info("uploading live master playlist")

	gctx := context.Background()
	cli, err := storage.NewClient(gctx)
	if err != nil {
		return nil, nil, err
	}
	defer cli.Close()

	bh := cli.Bucket(s.bucket)
	if _, err := bh.Attrs(ctx); err != nil {
		return nil, nil, err
	}

	obj := bh.Object(objectName)
	w := obj.NewWriter(ctx)
	w.CacheControl = "no-cache"

	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		return nil, nil, err
	}

	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(ctx)

	logger.Info("live master playlist has been uploaded successfully")

	return obj, attrs, err
}

func parseReqPath(path string) (string, int, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", 0, errors.New("wrong request path format")
	}

	streamID := parts[0]

	sparts := strings.Split(parts[1], ".")
	if len(sparts) != 2 {
		return "", 0, errors.New("wrong segment format")
	}

	segmentNum, err := strconv.Atoi(sparts[0])
	if err != nil {
		return "", 0, err
	}

	return streamID, segmentNum, nil
}
