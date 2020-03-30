package service

import (
	"github.com/go-redis/redis"
	"github.com/videocoin/cloud-sync/eventbus"
)

type Service struct {
	cfg        *Config
	rpc        *RPCServer
	eb         *eventbus.EventBus
	httpServer *HTTPServer
}

func NewService(cfg *Config) (*Service, error) {
	redisOpts, err := redis.ParseURL(cfg.DBURI)
	if err != nil {
		return nil, err
	}

	redisOpts.MaxRetries = 3
	redisOpts.PoolSize = 50

	dbcli := redis.NewClient(redisOpts)

	ds, err := NewDatastore(dbcli)
	if err != nil {
		return nil, err
	}

	ebConfig := &eventbus.Config{
		URI:    cfg.MQURI,
		Name:   cfg.Name,
		Logger: cfg.Logger.WithField("system", "eventbus"),
		Bucket: cfg.Bucket,
	}
	eb, err := eventbus.New(ebConfig)
	if err != nil {
		return nil, err
	}

	rpcOptions := &RPCServerOptions{
		Addr:   cfg.RPCAddr,
		Logger: cfg.Logger,
		DS:     ds,
		EB:     eb,
		Bucket: cfg.Bucket,
	}

	rpc, err := NewRPCServer(rpcOptions)
	if err != nil {
		return nil, err
	}

	httpServerOpts := &HTTPServerOptions{
		Addr:   cfg.HTTPAddr,
		Logger: cfg.Logger.WithField("system", "http-server"),
		Bucket: cfg.Bucket,
		DS:     ds,
		EB:     eb,
	}
	hs, err := NewHTTPServer(httpServerOpts)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg:        cfg,
		rpc:        rpc,
		eb:         eb,
		httpServer: hs,
	}

	return svc, nil
}

func (s *Service) Start(errCh chan error) {
	go func() {
		errCh <- s.rpc.Start()
	}()

	go func() {
		errCh <- s.eb.Start()
	}()

	go func() {
		errCh <- s.httpServer.Start()
	}()
}

func (s *Service) Stop() error {
	return s.eb.Stop()
}
