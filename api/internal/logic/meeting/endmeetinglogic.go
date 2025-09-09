package meeting

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/meeting/rpc/meeting"
	"GoMeeting/rpcs/user/rpc/user"
	"context"
	"runtime/debug"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EndMeetingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEndMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EndMeetingLogic {
	return &EndMeetingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EndMeetingLogic) EndMeeting(req *types.EndMeetingReq) (resp *types.Result, err error) {
	//参数校验
	if req.MeetingId == 0 {
		return types.NewErrorResultWithCode(code.ParamErrorCode), nil
	}

	//调用RPC服务
	result, err := l.svcCtx.UserRpc.GetIndexByUserId(l.ctx, &user.GetIndexByUserIdReq{
		UserId: req.UserId,
	})
	//系统错误
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("UserRpc.GetIndexByUserId error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	//业务错误
	if result.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result), nil
	}

	//调用RPC服务
	result2, err := l.svcCtx.MeetingRpc.EndMeeting(l.ctx, &meeting.EndMeetingReq{
		UserIndex: result.Index,
		MeetingId: req.MeetingId,
	})
	//系统错误
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("MeetingRpc.EndMeeting error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	//业务错误
	if result2.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result2), nil
	}

	return types.NewSuccessResult(), nil
}
