package service

import (
	"context"
	"errors"
	"net"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	protoempty "github.com/gogo/protobuf/types"
	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/syncer/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-sync/eventbus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type RPCServerOptions struct {
	Addr   string
	Bucket string
	Logger *logrus.Entry
	DS     *Datastore
	EB     *eventbus.EventBus
}

type RPCServer struct {
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

func NewRPCServer(opts *RPCServerOptions) (*RPCServer, error) {
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

	rpcServer := &RPCServer{
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

func (s *RPCServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *RPCServer) Sync(ctx context.Context, req *v1.SyncRequest) (*protoempty.Empty, error) {
	return &protoempty.Empty{}, nil
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
