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

func (l *PinguserLogic) Pinguser(req *types.PingReq) (resp *types.Result, err error) {
	// todo: add your logic here and delete this line
	result, err := l.svcCtx.UserRpc.Ping(l.ctx, &user.PingReq{Msg: req.Msg})
	return types.NewSuccessMessageResult(result.Msg + "!!!"), nil
}
