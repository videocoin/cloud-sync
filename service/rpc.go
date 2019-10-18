package service

import (
	"bytes"
	"context"
	"net"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-api/rpc"
	v1 "github.com/videocoin/cloud-api/syncer/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RpcServerOptions struct {
	Addr   string
	Writer *Writer
	Logger *logrus.Entry
}

type RpcServer struct {
	addr   string
	writer *Writer
	grpc   *grpc.Server
	listen net.Listener
	logger *logrus.Entry
}

func NewRpcServer(opts *RpcServerOptions) (*RpcServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcOpts = append(grpcOpts, grpc.MaxRecvMsgSize(1024*1024*1024))

	grpcServer := grpc.NewServer(grpcOpts...)

	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}

	rpcServer := &RpcServer{
		addr:   opts.Addr,
		writer: opts.Writer,
		grpc:   grpcServer,
		listen: listen,
		logger: opts.Logger,
	}

	v1.RegisterSyncerServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *RpcServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *RpcServer) Health(ctx context.Context, req *protoempty.Empty) (*rpc.HealthStatus, error) {
	return &rpc.HealthStatus{Status: "OK"}, nil
}

func (s *RpcServer) Sync(ctx context.Context, req *v1.SyncRequest) (*protoempty.Empty, error) {
	s.logger.WithField("path", req.Path).Info("syncing")

	go func(ctx context.Context, req *v1.SyncRequest) {
		logger := s.logger.WithField("path", req.Path)

		data := req.GetData()
		if data == nil {
			logger.Error("empty data")
			return
		}

		ws := s.writer.NewSession(context.Background(), req.Path, req.ContentType)
		defer func() {
			err := ws.Close(true)
			if err != nil {
				logger.Errorf("failed to close: %s", err)
			}
		}()

		err := ws.Write(bytes.NewReader(data))
		if err != nil {
			logger.Errorf("failed to write: %s", err.Error())
			return
		}

		logger.Infof("successfully synced file")
	}(ctx, req)

	return &protoempty.Empty{}, nil
}
