package logic

import (
	"GoMeeting/pkg/ctxdata"
	code "GoMeeting/pkg/result"
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	// 无人时结束会议(待实现)
	// 主持人结束会议后实现成员批量退会(待实现)

	//主持人才可结束会议
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

	//成员批量退会
	err = allLeaveMeeting(l, in.MeetingId)
	if err != nil {
		return &meeting.Result{
			Code: code.MeetingEndFailCode,
			Msg:  "成员批量退会失败",
		}, nil
	}

	return &meeting.Result{
		Code: code.SUCCESSCode,
	}, nil
}

func allLeaveMeeting(l *EndMeetingLogic, meetingId uint64) error {
	// 获取指定会议的所有成员ID
	memberSetKey := fmt.Sprintf(ctxdata.MeetingMemberPrefix, meetingId)
	members, err := l.svcCtx.Redis.SmembersCtx(l.ctx, memberSetKey)
	if err != nil {
		l.Logger.Errorf("从Redis获取会议成员ID失败: %v", err)
		return err
	}
	for _, memberIdStr := range members {
		// 调用同一个服务内的 LeaveMeeting 方法
		// 将字符串转换为 uint64
		memberId, err := strconv.ParseUint(memberIdStr, 10, 64)
		if err != nil {
			l.Logger.Errorf("成员ID转换失败: %v", err)
			return err
		}
		err = leaveMeeting(l, meetingId, memberId)
		if err != nil {
			l.Logger.Errorf("调用RPC本地LeaveMeeting 批量退会失败: %v", err)
			return err
		}
	}
	return nil
}

// 调用同一个服务内的 LeaveMeeting 方法
func leaveMeeting(l *EndMeetingLogic, meetingId uint64, userId uint64) error {
	// 调用同一个服务内的 JoinMeeting 方法
	leaveReq := &meeting.LeaveMeetingReq{
		UserId:    userId,
		MeetingId: meetingId,
	}
	// 创建 JoinMeetingLogic 实例并调用
	leaveLogic := NewLeaveMeetingLogic(l.ctx, l.svcCtx)
	leaveResult, err := leaveLogic.LeaveMeeting(leaveReq)
	if err != nil {
		l.Logger.Errorf("调用RPC本地LeaveMeeting 失败: %v", err)
		return err
	}
	// 检查加入会议的结果
	if leaveResult.Code != code.SUCCESSCode {
		l.Logger.Errorf("调用RPC本地LeaveMeeting离开会议失败: %v", err)
		return errors.New("调用RPC本地JoinMeeting离开会议失败")
	}
	return nil
}
