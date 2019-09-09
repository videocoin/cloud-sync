package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	RPCAddr string `default:"0.0.0.0:5009"`
	Bucket  string `required:"true" envconfig:"BUCKET" default:"testdserkin"`

	Logger *logrus.Entry `envconfig:"error"`
}
