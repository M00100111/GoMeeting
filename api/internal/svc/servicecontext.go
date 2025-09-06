package svc

import (
	"GoMeeting/api/internal/config"
	"GoMeeting/rpcs/meeting/rpc/meetingclient"
	"GoMeeting/rpcs/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis

	UserRpc    userclient.User
	MeetingRpc meetingclient.Meeting
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		Redis:      redis.MustNewRedis(c.Redis),
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		MeetingRpc: meetingclient.NewMeeting(zrpc.MustNewClient(c.MeetingRpc)),
	}
}
