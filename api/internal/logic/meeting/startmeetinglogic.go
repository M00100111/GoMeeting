package meeting

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/meeting/rpc/meeting"
	"GoMeeting/rpcs/user/rpc/user"
	"context"
	"fmt"
	"runtime/debug"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartMeetingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStartMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartMeetingLogic {
	return &StartMeetingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartMeetingLogic) StartMeeting(req *types.StartMeetingReq) (resp *types.Result, err error) {
	//全局鉴权
	//参数校验
	if req.MeetingId == 0 || req.MeetingName == "" || (req.JoinType == 0 && req.Password != "") || (req.JoinType == 1 && len(req.Password) != 5) {
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

	result2, err := l.svcCtx.MeetingRpc.StartMeeting(l.ctx, &meeting.StartMeetingReq{
		UserIndex:   result.Index,
		MeetingId:   req.MeetingId,
		MeetingName: req.MeetingName,
		JoinType:    req.JoinType,
		Password:    req.Password,
	})
	//系统错误
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("MeetingRpc.StartMeeting error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	//业务错误
	fmt.Println(result2.Code)
	if result2.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result2), nil
	}
	//统一错误处理和响应格式转换
	return types.NewSuccessResult(), nil
}
