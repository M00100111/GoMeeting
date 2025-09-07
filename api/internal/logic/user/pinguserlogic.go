package user

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/user/rpc/user"
	"context"
	"runtime/debug"

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
	result, err := l.svcCtx.UserRpc.Ping(l.ctx, &user.PingReq{Msg: req.Msg})
	//系统错误
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("UserRpc.Ping error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	//业务错误
	if result.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result), nil
	}
	return types.NewSuccessMessageResult(result.Msg), nil
}
