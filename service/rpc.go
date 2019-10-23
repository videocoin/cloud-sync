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
	"github.com/videocoin/cloud-api/rpc"
	streamsv1 "github.com/videocoin/cloud-api/streams/v1"
	v1 "github.com/videocoin/cloud-api/syncer/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-sync/eventbus"
	"google.golang.org/grpc"
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

func (s *RpcServer) Health(ctx context.Context, req *protoempty.Empty) (*rpc.HealthStatus, error) {
	return &rpc.HealthStatus{Status: "OK"}, nil
}

func (s *RpcServer) Sync(ctx context.Context, req *v1.SyncRequest) (*protoempty.Empty, error) {
	s.logger.WithField("path", req.Path).Info("syncing")

	go func(ctx context.Context, req *v1.SyncRequest) {
		logger := s.logger.WithField("path", req.Path)

		streamID, segmentNum, err := parseReqPath(req.Path)
		if err != nil {
			logger.Errorf("failed to parse request path: %s", err)
			return
		}

		data := req.GetData()
		if data == nil {
			logger.Error("empty data")
			return
		}

		emptyCtx := context.Background()

		err = s.uploadSegment(emptyCtx, streamID, segmentNum, req.ContentType, data)
		if err != nil {
			logger.Errorf("failed to upload segment: %s", err.Error())
			return
		}

		err = s.ds.AddSegment(streamID, segmentNum)
		if err != nil {
			logger.Errorf("failed to add segment: %s", err.Error())
			return
		}

		if segmentNum == 1 {
			logger.Info("updating stream status as ready")
			err = s.eb.EmitUpdateStreamStatus(ctx, streamID, streamsv1.StreamStatusReady)
			if err != nil {
				logger.Errorf("failed to update stream status: %s", err)
			}
		}

		logger.Info("generating and uploading live master playlist")
		err = s.generateAndUploadLiveMasterPlaylist(emptyCtx, streamID, segmentNum)
		if err != nil {
			logger.Errorf("failed to generate live master playlist: %s", err.Error())
			return
		}

	}(ctx, req)

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) uploadSegment(ctx context.Context, streamID string, segmentNum int, ct string, data []byte) error {
	logger := s.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"segment_num": segmentNum,
	})

	objectName := fmt.Sprintf("%s/%d.ts", streamID, segmentNum)

	logger = logger.WithField("object_name", objectName)

	w := s.bh.Object(objectName).NewWriter(ctx)
	w.ACL = []storage.ACLRule{
		storage.ACLRule{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	w.ContentType = ct
	defer func() {
		err := w.Close()
		if err != nil {
			logger.Errorf("failed to close: %s", err)
		}
	}()

	logger.Info("uploading segment")

	if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
		return err
	}

	logger.Infof("successfully synced segment")

	objAttrs, err := s.bh.Object(objectName).Attrs(ctx)
	if err != nil {
		return err
	}

	logger.Infof("segment object %+v\n", objAttrs)

	return nil
}

func (s *RpcServer) generateAndUploadLiveMasterPlaylist(ctx context.Context, streamID string, segmentNum int) error {
	logger := s.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"segment_num": segmentNum,
	})

	objectName := fmt.Sprintf("%s/index.m3u8", streamID)

	logger = logger.WithField("object_name", objectName)

	p, err := m3u8.NewMediaPlaylist(uint(segmentNum), uint(segmentNum))
	if err != nil {
		return err
	}

	for num := 1; num < segmentNum; num++ {
		err := p.Append(fmt.Sprintf("%d.ts", num), 10, "")
		if err != nil {
			return err
		}
	}

	data := p.Encode().Bytes()

	w := s.bh.Object(objectName).NewWriter(ctx)
	w.ACL = []storage.ACLRule{
		storage.ACLRule{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}
	w.ContentType = "application/x-mpegURL"
	defer func() {
		err := w.Close()
		if err != nil {
			logger.Errorf("failed to close: %s", err)
		}
	}()

	logger.Info("uploading live master playlist")

	if _, err = io.Copy(w, bytes.NewReader(data)); err != nil {
		return err
	}

	logger.Infof("successfully synced live master playlist")

	objAttrs, err := s.bh.Object(objectName).Attrs(ctx)
	if err != nil {
		return err
	}

	logger.Infof("playlist object %+v\n", objAttrs)

	return nil
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
