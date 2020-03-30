package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/cloud-pkg/mqmux"
	tracerext "github.com/videocoin/cloud-pkg/tracer"
	"google.golang.org/api/iterator"
)

type Config struct {
	Logger *logrus.Entry
	URI    string
	Name   string
	Bucket string
}

type EventBus struct {
	logger *logrus.Entry
	mq     *mqmux.WorkerMux
	bucket string
}

func New(c *Config) (*EventBus, error) {
	mq, err := mqmux.NewWorkerMux(c.URI, c.Name)
	if err != nil {
		return nil, err
	}
	if c.Logger != nil {
		mq.Logger = c.Logger
	}

	return &EventBus{
		logger: c.Logger,
		mq:     mq,
		bucket: c.Bucket,
	}, nil
}

func (e *EventBus) Start() error {
	err := e.mq.Publisher("streams.status")
	if err != nil {
		return err
	}

	err = e.mq.Consumer("streams.delete", 1, false, e.handleStreamDelete)
	if err != nil {
		return err
	}

	return e.mq.Run()
}

func (e *EventBus) Stop() error {
	return e.mq.Close()
}

func (e *EventBus) EmitUpdateStreamStatus(ctx context.Context, id string, status streamsv1.StreamStatus) error {
	headers := make(amqp.Table)

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		ext.SpanKindRPCServer.Set(span)
		ext.Component.Set(span, "syncer")
		err := span.Tracer().Inject(
			span.Context(),
			opentracing.TextMap,
			mqmux.RMQHeaderCarrier(headers),
		)
		if err != nil {
			e.logger.Errorf("failed to span inject: %s", err)
		}
	}

	event := &privatev1.Event{
		Type:     privatev1.EventTypeUpdateStatus,
		StreamID: id,
		Status:   status,
	}
	err := e.mq.PublishX("streams.status", event, headers)
	if err != nil {
		return err
	}
	return nil
}

func (e *EventBus) handleStreamDelete(d amqp.Delivery) error {
	var span opentracing.Span
	tracer := opentracing.GlobalTracer()
	spanCtx, err := tracer.Extract(opentracing.TextMap, mqmux.RMQHeaderCarrier(d.Headers))

	e.logger.Debugf("handling body: %+v", string(d.Body))

	if err != nil {
		span = tracer.StartSpan("eventbus.handleStreamDelete")
	} else {
		span = tracer.StartSpan("eventbus.handleStreamDelete", ext.RPCServerOption(spanCtx))
	}

	defer span.Finish()

	req := new(privatev1.Event)
	err = json.Unmarshal(d.Body, req)
	if err != nil {
		tracerext.SpanLogError(span, err)
		return err
	}

	span.SetTag("stream_id", req.StreamID)
	span.SetTag("event_type", req.Type.String())

	e.logger.Infof("handling request %+v", req)

	path := fmt.Sprintf("%s/", req.StreamID)

	e.logger.Infof("removing %s", path)

	gctx := context.Background()
	cli, err := storage.NewClient(gctx)
	if err != nil {
		return err
	}
	defer cli.Close()

	bh := cli.Bucket(e.bucket)
	it := bh.Objects(gctx, &storage.Query{
		Prefix:    path,
		Delimiter: "",
	})

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		e.logger.Infof("removing %s", attrs.Name)

		if err := bh.Object(attrs.Name).Delete(gctx); err != nil {
			e.logger.Infof("failed to delete object %s: %s", attrs.Name, err)
			continue
		}
	}

	return nil
}
