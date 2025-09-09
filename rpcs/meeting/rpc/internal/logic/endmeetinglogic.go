package logic

import (
	"GoMeeting/pkg/ctxdata"
	code "GoMeeting/pkg/result"
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"strconv"
	"time"

	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"

	"github.com/zeromicro/go-zero/core/logx"
)

type EndMeetingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEndMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EndMeetingLogic {
	return &EndMeetingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EndMeetingLogic) EndMeeting(in *meeting.EndMeetingReq) (*meeting.Result, error) {
	//获取会议信息
	//根据会议号查询会议信息主键
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
	//判断会议状态

	//获取成员信息
	//根据会议信息主键与用户主键查询会议信息
	member, err := l.svcCtx.MeetingMemberModel.FindOneByMeetingIdUserId(l.ctx, meetingInfo.Id, in.UserIndex)
	//数据库操作出错
	if err != nil && err != sqlc.ErrNotFound {
		return &meeting.Result{
			Code: code.ErrDbOpCode,
		}, nil
	}
	if err == sqlc.ErrNotFound {
		return &meeting.Result{
			Code: code.ErrRecordNotExistCode,
		}, nil
	}
	//鉴权
	if member.UserType != 1 {
		return &meeting.Result{
			Code: code.UserNoPermissionCode,
			Msg:  "非会议主持人无法关闭会议",
		}, nil
	}

	//会议已结束
	if meetingInfo.Status == 0 && meetingInfo.EndTime.Valid && meetingInfo.EndTime.Time.Before(time.Now()) {
		return &meeting.Result{
			Code: code.MeetingAlreadyEndedCode,
		}, nil
	}

	//关闭会议
	meetingInfo.Status = 0
	meetingInfo.EndTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	err = l.svcCtx.MeetingInfoModel.Update(l.ctx, meetingInfo)
	if err != nil {
		return &meeting.Result{
			Code: code.MeetingEndFailCode,
		}, nil
	}

	// 删除redis缓存
	_, err = l.svcCtx.Redis.Zrem(ctxdata.MeetingOngoingPrefix, strconv.FormatUint(in.MeetingId, 10))
	if err != nil {
		l.Logger.Errorf("从Redis删除会议信息失败: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
			Msg:  "从Redis删除会议信息失败",
		}, nil
	}

	return &meeting.Result{
		Code: code.SUCCESSCode,
	}, nil
}
