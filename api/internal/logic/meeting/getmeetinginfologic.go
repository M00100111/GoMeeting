package meeting

import (
	"GoMeeting/pkg/ctxdata"
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/meeting/rpc/meeting"
	"context"
	"runtime/debug"
	"strconv"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMeetingInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMeetingInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeetingInfoLogic {
	return &GetMeetingInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeetingInfoLogic) GetMeetingInfo(req *types.GetMeetingInfoReq) (resp *types.Result, err error) {
	// 根据用户ID获取对应的会议ID
	meetingIdStr, err := l.svcCtx.Redis.HgetCtx(l.ctx, ctxdata.OnMeetingUserPrefix, strconv.FormatUint(req.UserId, 10))
	if err != nil {
		l.Logger.Errorf("Failed to get meeting ID from RedisHash: %v", err)
		return &types.Result{
			Code: code.ErrRedisOpCode,
		}, nil
	}
	// 将获取到的字符串转换为 uint64
	meetingId, err := strconv.ParseUint(meetingIdStr, 10, 64)
	if err != nil {
		l.Logger.Errorf("Failed to parse meeting ID: %v", err)
		return &types.Result{
			Code: code.ErrParamParseCode,
		}, nil
	}

	meetingInfo, err := l.svcCtx.MeetingRpc.GetMeetingInfo(l.ctx, &meeting.GetMeetingInfoReq{
		MeetingId: meetingId,
	})
	//系统错误
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("UserRpc.GetIndexByUserId error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	//业务错误
	if meetingInfo.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(meetingInfo), nil
	}
	return types.NewSuccessDataResult(&types.GetMeetingInfoResp{
		MeetingId:   meetingInfo.MeetingId,
		MeetingName: meetingInfo.MeetingName,
		JoinType:    meetingInfo.JoinType,
		StartTime:   meetingInfo.StartTime,
	}), nil
}
