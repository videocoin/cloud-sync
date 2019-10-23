package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	RPCAddr string `default:"0.0.0.0:5021"`
	Bucket  string `required:"true" default:"testvc01" envconfig:"BUCKET"`
	DBURI   string `required:"true" default:"redis://127.0.0.1:6379/0" envconfig:"DBURI"`
	MQURI   string `default:"amqp://guest:guest@127.0.0.1:5672" envconfig:"MQURI"`

	Logger *logrus.Entry `envconfig:"error"`
}
