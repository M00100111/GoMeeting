package svc

import (
	"GoMeeting/api/internal/config"
	"GoMeeting/api/internal/middleware"
	"GoMeeting/rpcs/meeting/rpc/meetingclient"
	"GoMeeting/rpcs/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	JwtAuth    rest.Middleware
	UserRpc    userclient.User
	MeetingRpc meetingclient.Meeting
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redis),
		JwtAuth:    middleware.NewJwtAuthMiddleware(c.JwtAuth.AccessSecret).Handle,
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		MeetingRpc: meetingclient.NewMeeting(zrpc.MustNewClient(c.MeetingRpc)),
	}
}

//注册中间件函数
