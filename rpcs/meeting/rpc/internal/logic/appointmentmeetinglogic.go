package logic

import (
	"context"

	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppointmentMeetingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppointmentMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppointmentMeetingLogic {
	return &AppointmentMeetingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppointmentMeetingLogic) AppointmentMeeting(in *meeting.AppointmentMeetingReq) (*meeting.Result, error) {
	// todo: add your logic here and delete this line
	//预定开启时间
	//通知会议成员
	return &meeting.Result{}, nil
}
