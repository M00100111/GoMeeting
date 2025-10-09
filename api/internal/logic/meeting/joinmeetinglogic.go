package meeting

import (
	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/meeting/rpc/meeting"
	"GoMeeting/rpcs/user/rpc/user"
	"context"
	"fmt"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinMeetingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinMeetingLogic {
	return &JoinMeetingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinMeetingLogic) JoinMeeting(req *types.JoinMeetingReq) (resp *types.Result, err error) {
	//鉴权
	if req.UserId == 0 || req.MeetingId == 0 {
		return types.NewErrorResultWithCode(code.ParamErrorCode), nil
	}
	//调用RPC服务
	result, err := l.svcCtx.UserRpc.GetUserInfoByUserId(l.ctx, &user.GetUserInfoByUserIdReq{
		UserId: req.UserId,
	})
	if err != nil {
		l.Logger.Errorf("UserRpc.GetIndexByUserId error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	if result.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result), nil
	}

	result2, err := l.svcCtx.MeetingRpc.JoinMeeting(l.ctx, &meeting.JoinMeetingReq{
		UserIndex:    result.UserInfo.Index,
		UserId:       req.UserId,
		MeetingId:    req.MeetingId,
		Password:     req.Password,
		Username:     result.UserInfo.Username,
		Sex:          result.UserInfo.Sex,
		Email:        result.UserInfo.Email,
		MicStatus:    req.MicStatus,
		CameraStatus: req.CameraStatus,
		ScreenStatus: req.ScreenStatus,
	})
	if err != nil {
		l.Logger.Errorf("UserRpc.JoinMeeting error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	fmt.Println(result2)
	if result2.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result2), nil
	}

	return types.NewSuccessResult(), nil
}
