package logic

import (
	"context"

	"GoMeeting/rpcs/user/rpc/internal/svc"
	"GoMeeting/rpcs/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 定义功能及请求与响应结构体
func (l *PingLogic) Ping(in *user.PingReq) (*user.PingResp, error) {
	// todo: add your logic here and delete this line
	return &user.PingResp{
		Msg: in.Msg + "!!!",
	}, nil
}
