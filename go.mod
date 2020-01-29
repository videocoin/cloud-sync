module github.com/videocoin/cloud-sync

go 1.12

require (
	cloud.google.com/go v0.37.4
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gogo/protobuf v1.3.1
	github.com/grafov/m3u8 v0.11.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/opentracing/opentracing-go v1.1.0
	github.com/patrickmn/sortutil v0.0.0-20120526081524-abeda66eb583
	github.com/sirupsen/logrus v1.4.2
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/videocoin/cloud-api v0.2.14
	github.com/videocoin/cloud-pkg v0.0.6
	google.golang.org/grpc v1.23.0
)

replace github.com/videocoin/cloud-api => ../cloud-api

replace github.com/videocoin/cloud-pkg => ../cloud-pkg
