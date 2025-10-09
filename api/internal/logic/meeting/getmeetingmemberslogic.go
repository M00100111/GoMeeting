package meeting

import (
	"GoMeeting/pkg/ctxdata"
	code "GoMeeting/pkg/result"
	"context"
	"fmt"
	"strconv"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMeetingMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMeetingMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeetingMembersLogic {
	return &GetMeetingMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeetingMembersLogic) GetMeetingMembers(req *types.GetMeetingMembersReq) (resp *types.Result, err error) {
	// 根据成员ID从 Redis Hash 中获取会议ID
	meetingId, err := l.svcCtx.Redis.HgetCtx(l.ctx, ctxdata.OnMeetingUserPrefix, strconv.FormatUint(req.UserId, 10))
	if err != nil {
		l.Logger.Errorf("Failed to search member ID in RedisHash: %v", err)
		return types.NewErrorResultWithCode(code.ErrRedisOpCode), nil
	}
	if meetingId == "" {
		return types.NewErrorResultWithCode(code.UserNotInMeetingCode), nil
	}

	// 从Redis Set中获取当前会议的所有成员
	memberSetKey := fmt.Sprintf(ctxdata.MeetingMemberPrefix, meetingId)
	members, err := l.svcCtx.Redis.SmembersCtx(l.ctx, memberSetKey)
	if err != nil {
		l.Logger.Errorf("Failed to get members from RedisSet: %v", err)
		return types.NewErrorResultWithCode(code.ErrRedisOpCode), nil
	}
	var result []uint64
	for _, memberStr := range members {
		memberId, _ := strconv.ParseUint(memberStr, 10, 64)
		result = append(result, memberId)
	}
	return types.NewSuccessDataResult(result), nil
}
