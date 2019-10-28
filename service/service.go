package service

import (
	"github.com/go-redis/redis"
	"github.com/videocoin/cloud-sync/eventbus"
)

type Service struct {
	cfg        *Config
	rpc        *RpcServer
	eb         *eventbus.EventBus
	httpServer *HttpServer
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
	}
	eb, err := eventbus.New(ebConfig)
	if err != nil {
		return nil, err
	}

	rpcOptions := &RpcServerOptions{
		Addr:   cfg.RPCAddr,
		Logger: cfg.Logger,
		DS:     ds,
		EB:     eb,
		Bucket: cfg.Bucket,
	}

	rpc, err := NewRpcServer(rpcOptions)
	if err != nil {
		return nil, err
	}

	httpServerOpts := &HttpServerOptions{
		Addr:   cfg.HTTPAddr,
		Logger: cfg.Logger.WithField("system", "http-server"),
		Bucket: cfg.Bucket,
		DS:     ds,
		EB:     eb,
	}
	hs, err := NewHttpServer(httpServerOpts)
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

func (s *Service) Start() error {
	go s.rpc.Start()
	go s.eb.Start()
	go s.httpServer.Start()

	return nil
}

func (s *Service) Stop() error {
	s.eb.Stop()
	return nil
}
