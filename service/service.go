package service

type Service struct {
	cfg *Config
	rpc *RpcServer
}

func NewService(cfg *Config) (*Service, error) {
	writerOptions := &WriterOptions{
		Bucket: cfg.Bucket,
		Logger: cfg.Logger,
	}

	writer, err := NewWriter(writerOptions)
	if err != nil {
		return nil, err
	}

	rpcOptions := &RpcServerOptions{
		Addr:   cfg.RPCAddr,
		Writer: writer,
		Logger: cfg.Logger,
	}

	rpc, err := NewRpcServer(rpcOptions)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg: cfg,
		rpc: rpc,
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.rpc.Start()
	return nil
}

func (s *Service) Stop() error {
	return nil
}
