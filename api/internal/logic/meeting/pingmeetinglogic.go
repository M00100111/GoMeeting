package meeting

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/meeting/rpc/meeting"
	"context"
	"runtime/debug"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingMeetingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingMeetingLogic {
	return &PingMeetingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingMeetingLogic) PingMeeting(req *types.PingReq) (resp *types.Result, err error) {
	result, err := l.svcCtx.MeetingRpc.Ping(l.ctx, &meeting.PingReq{Msg: req.Msg})
	//系统错误
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("MeetingRpc.Ping error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	//业务错误
	if result.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result), nil
	}
	return types.NewSuccessMessageResult(result.Msg), nil
}
