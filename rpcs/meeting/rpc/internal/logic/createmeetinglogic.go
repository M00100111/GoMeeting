package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/meeting/models"
	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"runtime/debug"
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
		MeetingId:   in.MeetingId,
		MeetingName: in.Username + "的会议",
		UserId:      in.UserId,
	}
	err := l.svcCtx.MeetingInfoModel.CreateMeeting(l.ctx, meetingInfo)
	if err != nil {
		l.Logger.Errorf("注册用户会议号和会议成员信息失败: %v, stack: %s", err, debug.Stack())
		return &meeting.Result{
			Code: code.ErrDbOpCode,
		}, nil
	}

	return &meeting.Result{
		Code: code.SUCCESSCode,
	}, nil
}
