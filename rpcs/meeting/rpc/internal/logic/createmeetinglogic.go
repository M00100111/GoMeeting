package logic

import (
	"GoMeeting/rpcs/meeting/models"
	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type CreateMeetingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMeetingLogic {
	return &CreateMeetingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMeetingLogic) CreateMeeting(in *meeting.CreateMeetingReq) (*meeting.Result, error) {
	// 创建会议信息
	meetingInfo := &models.MeetingInfo{
		MeetingId:   strconv.Itoa(int(in.UserId)),
		MeetingName: in.Username + "的会议",
		UserId:      in.UserId,
	}
	err := l.svcCtx.MeetingInfoModel.CreateMeeting(l.ctx, meetingInfo)
	if err != nil {
		return nil, err
	}
	return &meeting.Result{}, nil
}
