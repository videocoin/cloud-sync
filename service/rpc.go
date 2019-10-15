package service

import (
	"bytes"
	"context"
	"fmt"
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
	ws := s.writer.NewSession(context.Background(), req.Path, req.ContentType)

	if data := req.GetData(); data != nil {
		if ws == nil {
			err := fmt.Errorf("failed to start write session")
			s.logger.Errorf(err.Error())

			return nil, err
		}

		go func() {
			defer func() {
				if err := ws.Close(true); err != nil {
					s.logger.Errorf("failed to close: %s", err)
				}
			}()

			if err := ws.Write(bytes.NewReader(data)); err != nil {
				err := fmt.Errorf("failed to write: %s", err.Error())
				s.logger.Errorf(err.Error())
			}

			s.logger.WithField("path", req.Path).Infof("successfully synced file")
		}()
	}

	return &protoempty.Empty{}, nil
}
