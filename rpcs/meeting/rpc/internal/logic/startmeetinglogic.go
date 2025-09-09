package logic

import (
	"GoMeeting/pkg/ctxdata"
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"strconv"
	"time"
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
	//会议是否已开始

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
			Msg:  "非会议主持人无法开启会议",
		}, nil
	}

	//当会议已开始但用户为会议主持人则跳转到申请入会请求，上面已鉴权
	if meetingInfo.Status == 1 {
		//申请入会,修改当前人的入会信息
		err = joinMeeting(l, in.MeetingId, in.UserIndex, in.Password)
		if err != nil {
			return &meeting.Result{
				Code: code.SYS_ERRORCode,
				Msg:  "调用本地JoinMeeting加入会议失败",
			}, nil
		}
	}

	//开启会议查询是否有以当前人为房主的进行中的其他会议
	meetingInfo2, err := l.svcCtx.MeetingInfoModel.FindOneByUserId(l.ctx, in.UserIndex)
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
	if meetingInfo.Id != meetingInfo2.Id && meetingInfo2.Status == 1 {
		return &meeting.Result{
			Code: code.MeetingUserInOtherMeetingCode,
		}, nil
	}

	//修改会议信息
	meetingInfo.Status = 1 //开启会议
	//开启会议的时间
	meetingInfo.StartTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	meetingInfo.EndTime = sql.NullTime{}
	// 如果需要更新会议名称、加入类型或密码
	if in.MeetingName != "" {
		meetingInfo.MeetingName = in.MeetingName
	}
	meetingInfo.JoinType = in.JoinType
	// 只有在提供了新密码时才更新
	if in.JoinType == 0 {
		meetingInfo.MeetingPassword = ""
	}
	if in.JoinType == 1 && in.Password != "" {
		meetingInfo.MeetingPassword = in.Password
	}
	err = l.svcCtx.MeetingInfoModel.Update(l.ctx, meetingInfo)
	if err != nil {
		return &meeting.Result{
			Code: code.ErrDbOpCode,
		}, nil
	}

	//记录会议信息到Redis,会议超过一定时间未结束自动结束
	// 将会议ID和过期时间(当前时间+5小时)存入ZSet
	expireTime := meetingInfo.StartTime.Time.Add(5 * time.Hour).Unix()
	_, err = l.svcCtx.Redis.Zadd(ctxdata.MeetingOngoingPrefix, expireTime, strconv.FormatUint(in.MeetingId, 10))
	if err != nil {
		l.Logger.Errorf("将会议信息存入Redis失败: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
			Msg:  "将会议信息存入Redis失败",
		}, nil
	}

	//申请入会,修改当前人的入会信息
	err = joinMeeting(l, in.MeetingId, in.UserIndex, in.Password)
	if err != nil {
		return &meeting.Result{
			Code: code.SYS_ERRORCode,
			Msg:  "调用本地JoinMeeting加入会议失败",
		}, nil
	}
	return &meeting.Result{
		Code: code.SUCCESSCode,
	}, nil
}

// 调用同一个服务内的 JoinMeeting 方法
func joinMeeting(l *StartMeetingLogic, meetingId, userIndex uint64, meetingPassword string) error {
	// 调用同一个服务内的 JoinMeeting 方法
	joinReq := &meeting.JoinMeetingReq{
		UserIndex: userIndex,
		MeetingId: meetingId,
		Password:  meetingPassword,
	}
	// 创建 JoinMeetingLogic 实例并调用
	joinLogic := NewJoinMeetingLogic(l.ctx, l.svcCtx)
	joinResult, err := joinLogic.JoinMeeting(joinReq)
	if err != nil {
		l.Logger.Errorf("调用RPC本地JoinMeeting 失败: %v", err)
		return err
	}
	// 检查加入会议的结果
	if joinResult.Code != code.SUCCESSCode {
		l.Logger.Errorf("调用RPC本地JoinMeeting加入会议失败: %v", err)
		return errors.New("调用RPC本地JoinMeeting加入会议失败")
	}
	return nil
}
