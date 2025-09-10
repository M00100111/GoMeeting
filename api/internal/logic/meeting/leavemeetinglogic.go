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

type LeaveMeetingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLeaveMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveMeetingLogic {
	return &LeaveMeetingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LeaveMeetingLogic) LeaveMeeting(req *types.LeaveMeetingReq) (resp *types.Result, err error) {
	if req.MeetingId == 0 || req.UserId == 0 {
		return types.NewErrorResultWithCode(code.ParamErrorCode), nil
	}
	result, err := l.svcCtx.MeetingRpc.LeaveMeeting(l.ctx, &meeting.LeaveMeetingReq{
		UserId:    req.UserId,
		MeetingId: req.MeetingId,
	})
	if err != nil {
		// 错误日志
		l.Logger.Errorf("MeetingRpc.LeaveMeeting error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	if result.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result), nil
	}
	return types.NewSuccessResult(), nil
}
