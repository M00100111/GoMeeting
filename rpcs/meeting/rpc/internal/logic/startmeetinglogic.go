package logic

import (
	"context"

	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"

	"github.com/zeromicro/go-zero/core/logx"
)

type StartMeetingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStartMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartMeetingLogic {
	return &StartMeetingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StartMeetingLogic) StartMeeting(in *meeting.StartMeetingReq) (*meeting.Result, error) {
	// todo: add your logic here and delete this line

	return &meeting.Result{}, nil
}
