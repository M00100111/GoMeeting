package svc

import (
	"GoMeeting/rpcs/ws/internal/config"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	KafkaPusher *kq.Pusher // 生产者
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		Redis:       redis.MustNewRedis(c.Redis),
		KafkaPusher: kq.NewPusher(c.KafkaPusherConf.Brokers, c.KafkaPusherConf.Topic),
	}
}
