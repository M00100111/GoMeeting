package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	Name     string
	ListenOn string
	Pattern  string

	Redis redis.RedisConf

	Jwt struct {
		AccessSecret string
		AccessExpire int64
	}
	MongoDB struct {
		Url string
		Db  string
	}

	KafkaPusherConf struct {
		Brokers []string
		Topic   string
	}
	KafkaConsumerConf kq.KqConf
}
