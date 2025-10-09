package logic

import (
	code "GoMeeting/pkg/result"
	"context"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMeetingInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMeetingInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeetingInfoLogic {
	return &GetMeetingInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMeetingInfoLogic) GetMeetingInfo(in *meeting.GetMeetingInfoReq) (*meeting.GetMeetingInfoResp, error) {
	//根据会议号查询会议信息主键
	meetingInfo, err := l.svcCtx.MeetingInfoModel.FindOneByMeetingId(l.ctx, in.MeetingId)
	//数据库操作出错
	if err != nil && err != sqlc.ErrNotFound {
		return &meeting.GetMeetingInfoResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	if err == sqlc.ErrNotFound {
		return &meeting.GetMeetingInfoResp{
			Code: code.MeetingNotExistCode,
		}, nil
	}

	var startTimeValue int64
	if meetingInfo.StartTime.Valid {
		startTimeValue = meetingInfo.StartTime.Time.Unix()
	} else {
		startTimeValue = 0 // 或其他默认值
	}
	return &meeting.GetMeetingInfoResp{
		Code:        code.SUCCESSCode,
		MeetingId:   meetingInfo.MeetingId,
		MeetingName: meetingInfo.MeetingName,
		JoinType:    meetingInfo.JoinType,
		StartTime:   startTimeValue,
	}, nil
}
