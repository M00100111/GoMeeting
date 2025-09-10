package svc

import (
	"GoMeeting/rpcs/meeting/models"
	"GoMeeting/rpcs/meeting/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	models.MeetingInfoModel
	models.MeetingMemberModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:             c,
		Redis:              redis.MustNewRedis(c.Redisx),
		MeetingInfoModel:   models.NewMeetingInfoModel(sqlConn, c.Cache),
		MeetingMemberModel: models.NewMeetingMemberModel(sqlConn, c.Cache),
	}
}
