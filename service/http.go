package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/grafov/m3u8"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
	streamsv1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-pkg/logrusext"
	"github.com/videocoin/cloud-sync/eventbus"
)

type HttpServerOptions struct {
	Addr   string
	Bucket string
	Logger *logrus.Entry
	DS     *Datastore
	EB     *eventbus.EventBus
}

type HttpServer struct {
	logger *logrus.Entry
	e      *echo.Echo
	addr   string
	bucket string
	ds     *Datastore
	eb     *eventbus.EventBus
	gscli  *storage.Client
	bh     *storage.BucketHandle
}

func NewHttpServer(opts *HttpServerOptions) (*HttpServer, error) {
	gscli, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}

	bh := gscli.Bucket(opts.Bucket)
	_, err = bh.Attrs(context.Background())
	if err != nil {
		return nil, err
	}

	return &HttpServer{
		logger: opts.Logger,
		e:      echo.New(),
		gscli:  gscli,
		bh:     bh,
		addr:   opts.Addr,
		ds:     opts.DS,
		eb:     opts.EB,
		bucket: opts.Bucket,
	}, nil
}

func (hs *HttpServer) Start() error {
	hs.e.Use(logrusext.Hook())
	hs.e.Use(middleware.Recover())
	hs.e.HideBanner = true
	hs.e.HidePort = true
	hs.e.DisableHTTP2 = true

	hs.e.Logger = logrusext.MWLogger{Entry: hs.logger}

	hs.e.POST("/api/v1/sync", hs.upload)

	hs.logger.Infof("starting http server on %s", hs.addr)

	return hs.e.Start(hs.addr)
}

func (hs *HttpServer) upload(c echo.Context) error {
	path := c.FormValue("path")
	ct := c.FormValue("ct")

	file, err := c.FormFile("file")
	if err != nil {
		hs.logger.Errorf("failed to get file: %s", err)
		return err
	}

	src, err := file.Open()
	if err != nil {
		hs.logger.Errorf("failed to open file: %s", err)
		return err
	}
	defer src.Close()

	logger := hs.logger.WithField("path", path)

	streamID, segmentNum, err := parseReqPath(path)
	if err != nil {
		e := fmt.Errorf("failed to parse request path: %s", err)
		hs.logger.Error(e)
		return e
	}

	emptyCtx := context.Background()
	_, _, err = hs.uploadSegment(emptyCtx, streamID, segmentNum, ct, src)
	if err != nil {
		e := fmt.Errorf("failed to upload segment: %s", err)
		hs.logger.Error(e)
		return e
	}

	logger.Info("generating and uploading live master playlist")
	_, _, err = hs.generateAndUploadLiveMasterPlaylist(emptyCtx, streamID, segmentNum)
	if err != nil {
		e := fmt.Errorf("failed to generate live master playlist: %s", err.Error())
		hs.logger.Error(e)
		return e
	}

	err = hs.ds.AddSegment(streamID, segmentNum)
	if err != nil {
		e := fmt.Errorf("failed to add segment: %s", err.Error())
		hs.logger.Error(e)
		return e
	}

	if segmentNum == 1 {
		logger.Info("updating stream status as ready")
		err = hs.eb.EmitUpdateStreamStatus(emptyCtx, streamID, streamsv1.StreamStatusReady)
		if err != nil {
			logger.Errorf("failed to update stream status: %s", err)
		}
	}

	return c.NoContent(http.StatusAccepted)
}

func (hs *HttpServer) uploadSegment(ctx context.Context, streamID string, segmentNum int, ct string, src multipart.File) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	objectName := fmt.Sprintf("%s/%d.ts", streamID, segmentNum)

	logger := hs.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"segment_num": segmentNum,
		"bucket":      hs.bucket,
		"object_name": objectName,
	})

	logger.Info("uploading segment")

	emptyCtx := context.Background()
	obj := hs.bh.Object(objectName)
	w := obj.NewWriter(emptyCtx)

	if _, err := io.Copy(w, src); err != nil {
		return nil, nil, err
	}

	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	if err := obj.ACL().Set(emptyCtx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(emptyCtx)
	if err != nil {
		return obj, attrs, err
	}

	logger.Info("segment has been uploaded successfully")

	return obj, attrs, err
}

func (hs *HttpServer) generateAndUploadLiveMasterPlaylist(ctx context.Context, streamID string, segmentNum int) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	objectName := fmt.Sprintf("%s/index.m3u8", streamID)

	logger := hs.logger.WithFields(logrus.Fields{
		"stream_id":   streamID,
		"segment_num": segmentNum,
		"bucket":      hs.bucket,
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

	obj := hs.bh.Object(objectName)
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
	if err != nil {
		return obj, attrs, err
	}

	logger.Info("live master playlist has been uploaded successfully")

	return obj, attrs, err
}
