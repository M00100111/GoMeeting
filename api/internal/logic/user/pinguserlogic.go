package user

import (
	"GoMeeting/rpcs/user/rpc/user"
	"context"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PinguserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPinguserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PinguserLogic {
	return &PinguserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PinguserLogic) Pinguser(req *types.PingReq) (resp *types.PingResp, err error) {
	result, err := l.svcCtx.User.Ping(l.ctx, &user.PingReq{Msg: req.Msg})
	resp = &types.PingResp{
		Msg: result.Msg,
	}
	return resp, err
}
