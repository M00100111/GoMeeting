package svc

import (
	"GoMeeting/rpcs/user/models"
	"GoMeeting/rpcs/user/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	models.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		Redis:     redis.MustNewRedis(c.Redisx),
		UserModel: models.NewUserModel(sqlConn, c.Cache),
	}
}
