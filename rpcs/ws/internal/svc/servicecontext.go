package svc

import (
	"GoMeeting/rpcs/ws/internal/config"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	KafkaMeetingPusher *kq.Pusher // 生产者
	KafkaSocialPusher  *kq.Pusher // 生产者
	KafkaWsPusher      *kq.Pusher // 生产者
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		Redis:              redis.MustNewRedis(c.Redis),
		KafkaMeetingPusher: kq.NewPusher(c.KafkaPusherConf.Brokers, c.KafkaPusherConf.Topic[0]),
		KafkaSocialPusher:  kq.NewPusher(c.KafkaPusherConf.Brokers, c.KafkaPusherConf.Topic[1]),
		KafkaWsPusher:      kq.NewPusher(c.KafkaPusherConf.Brokers, c.KafkaPusherConf.Topic[2]),
	}
}

// 创建消费者构建器
func (s *ServiceContext) CreateKafkaConsumer() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: s.Config.KafkaConsumerConf.Brokers,
		GroupID: s.Config.KafkaConsumerConf.Group,
		Topic:   s.Config.KafkaConsumerConf.Topic,
		//MinBytes: 10e3, // 10KB
		MaxBytes: 10e6,                  // 10MB
		MinBytes: 1,                     // 降低到1字节，确保每条消息都能及时消费
		MaxWait:  10 * time.Millisecond, // 设置最大等待时间为100毫秒
	})
}
