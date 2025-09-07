package logic

import (
	"GoMeeting/pkg/result"
	"context"
	"fmt"

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
	fmt.Println(in.Msg)
	return &user.PingResp{
		Code: code.SUCCESSCode,
		Msg:  in.Msg + "!!!",
	}, nil
}
