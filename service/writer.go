package service

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
)

type WriterOptions struct {
	Bucket string
	Logger *logrus.Entry
}

type Writer struct {
	bh *storage.BucketHandle
}

func NewWriter(wo *WriterOptions) (*Writer, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bh := client.Bucket(wo.Bucket)
	if _, err = bh.Attrs(ctx); err != nil {
		return nil, err
	}

	return &Writer{
		bh: bh,
	}, nil
}

func (w *Writer) NewSession(ctx context.Context, name string) *WriteSession {
	return NewWriteSession(ctx, w.bh.Object(name))
}

type WriteSession struct {
	ctx context.Context
	hdl *storage.ObjectHandle
	w   *storage.Writer
}

func NewWriteSession(ctx context.Context, oh *storage.ObjectHandle) *WriteSession {
	w := oh.NewWriter(ctx)
	return &WriteSession{
		ctx: ctx,
		hdl: oh,
		w:   w,
	}
}

func (ws *WriteSession) Write(r io.Reader) error {
	if _, err := io.Copy(ws.w, r); err != nil {
		return err
	}

	return nil
}

func (ws *WriteSession) Close(public bool) error {
	if err := ws.w.Close(); err != nil {
		return err
	}

	if public {
		if err := ws.hdl.ACL().Set(ws.ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return err
		}
	}

	return nil
}
