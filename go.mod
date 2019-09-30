module github.com/videocoin/cloud-sync

go 1.12

replace github.com/videocoin/cloud-api => ../cloud-api

require (
	cloud.google.com/go v0.37.4
	github.com/gogo/protobuf v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.4.2
	github.com/videocoin/cloud-api v0.2.13
	github.com/videocoin/cloud-pkg v0.0.5
	google.golang.org/grpc v1.23.0
)
