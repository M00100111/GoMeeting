package svc

import (
	"GoMeeting/api/internal/config"
	"GoMeeting/rpcs/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis

	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.Redis),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
