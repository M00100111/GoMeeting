package logic

import (
	"GoMeeting/pkg/ctxdata"
	code "GoMeeting/pkg/result"
	"context"
	"fmt"
	"strconv"

	"GoMeeting/rpcs/meeting/rpc/internal/svc"
	"GoMeeting/rpcs/meeting/rpc/meeting"

	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveMeetingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveMeetingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveMeetingLogic {
	return &LeaveMeetingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LeaveMeetingLogic) LeaveMeeting(in *meeting.LeaveMeetingReq) (*meeting.Result, error) {
	// 支持关闭会议调用批量删除(待实现)

	// 支持个人离会正常调用
	// ws心跳离线自动调用(待实现)

	// 从Redis Hash中删除成员的会议信息 (成员ID -> 会议ID)
	_, err := l.svcCtx.Redis.HdelCtx(l.ctx, ctxdata.OnMeetingUserPrefix, strconv.FormatUint(in.UserId, 10))
	if err != nil {
		l.Logger.Errorf("Failed to remove member ID from RedisHash: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
		}, nil
	}

	// 将会议室成员id和状态信息从Redis移除
	// 查询成员是否在会议中并删除
	memberSetKey := fmt.Sprintf(ctxdata.MeetingMemberPrefix, in.MeetingId)
	exists, err := l.svcCtx.Redis.SismemberCtx(l.ctx, memberSetKey, in.UserId)
	if err != nil {
		l.Logger.Errorf("Failed to check member ID in RedisSet: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
		}, nil
	}
	if !exists {
		// 成员不存在于会议中
		return &meeting.Result{
			Code: code.MeetingUserNotInMeetingCode,
		}, nil
	}
	// 成员存在，从集合中删除
	_, err = l.svcCtx.Redis.SremCtx(l.ctx, memberSetKey, in.UserId)
	if err != nil {
		l.Logger.Errorf("Failed to remove member ID from RedisSet: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
		}, nil
	}

	// 查询成员状态是否存在，如果存在则删除
	memberStatusKey := fmt.Sprintf(ctxdata.MeetingMemberDetailPrefix, in.MeetingId)
	exists, err = l.svcCtx.Redis.HexistsCtx(l.ctx, memberStatusKey, strconv.FormatUint(in.UserId, 10))
	if err != nil {
		l.Logger.Errorf("Failed to check member status in RedisHash: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
		}, nil
	}
	if !exists {
		// 成员不存在于会议中
		return &meeting.Result{
			Code: code.MeetingUserNotInMeetingCode,
		}, nil
	}
	// 成员状态存在，从哈希中删除
	_, err = l.svcCtx.Redis.HdelCtx(l.ctx, memberStatusKey, strconv.FormatUint(in.UserId, 10))
	if err != nil {
		l.Logger.Errorf("Failed to remove member status from RedisHash: %v", err)
		return &meeting.Result{
			Code: code.ErrRedisOpCode,
		}, nil
	}

	return &meeting.Result{
		Code: code.SUCCESSCode,
	}, nil
}
