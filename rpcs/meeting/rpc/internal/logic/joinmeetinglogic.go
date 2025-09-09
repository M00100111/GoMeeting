package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/meeting/models"
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"time"

	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinMeetingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinMeetingLogic {
	return &JoinMeetingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinMeetingLogic) JoinMeeting(in *meeting.JoinMeetingReq) (*meeting.Result, error) {
	//获取会议信息
	meetingInfo, err := l.svcCtx.MeetingInfoModel.FindOneByMeetingId(l.ctx, in.MeetingId)
	//数据库操作出错
	if err != nil && err != sqlc.ErrNotFound {
		return &meeting.Result{
			Code: code.ErrDbOpCode,
		}, nil
	}
	if err == sqlc.ErrNotFound {
		return &meeting.Result{
			Code: code.MeetingNotExistCode,
		}, nil
	}
	//判断会议状态是否已结束
	if meetingInfo.Status == 0 {
		//未到达会议时间
		if meetingInfo.StartTime.Valid && meetingInfo.StartTime.Time.After(time.Now()) {
			return &meeting.Result{
				Code: code.MeetingNotStartedCode,
			}, nil
		}
		//已过会议时间且已更新会议结束时间
		if meetingInfo.EndTime.Valid && meetingInfo.EndTime.Time.Before(time.Now()) {
			return &meeting.Result{
				Code: code.MeetingAlreadyEndedCode,
			}, nil
		}
	}
	//判断加入方式
	//验证密码
	if meetingInfo.JoinType == 1 && meetingInfo.MeetingPassword != in.Password {
		return &meeting.Result{
			Code: code.MeetingJoinCodeErrorCode,
		}, nil
	}

	//查询会议成员记录
	meetingMember, err := l.svcCtx.MeetingMemberModel.FindOneByMeetingIdUserId(l.ctx, meetingInfo.Id, in.UserIndex)
	//数据库操作出错
	if err != nil && err != sqlc.ErrNotFound {
		return &meeting.Result{
			Code: code.ErrDbOpCode,
		}, nil
	}
	if err == sqlc.ErrNotFound {
		//第一次加入新增信息
		_, err := l.svcCtx.MeetingMemberModel.Insert(l.ctx, &models.MeetingMember{
			MeetingId: meetingInfo.Id,
			UserId:    in.UserIndex,
			LastJoinTime: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			return &meeting.Result{
				Code: code.ErrDbOpCode,
				Msg:  "新增会议成员信息失败",
			}, nil
		}
	} else {
		//判断是否被拉黑
		if meetingMember.UserStatus != 0 {
			return &meeting.Result{
				Code: code.MeetingUserNotAllowedCode,
			}, nil
		}
		//非第一次加入修改入会时间
		meetingMember.LastJoinTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		err = l.svcCtx.MeetingMemberModel.Update(l.ctx, meetingMember)
		if err != nil {
			return &meeting.Result{
				Code: code.ErrDbOpCode,
				Msg:  "修改会议成员信息失败",
			}, nil
		}
	}
	//入会成功
	return &meeting.Result{
		Code: code.SUCCESSCode,
	}, nil
}
