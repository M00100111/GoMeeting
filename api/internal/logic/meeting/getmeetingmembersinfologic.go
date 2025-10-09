package meeting

import (
	"GoMeeting/pkg/ctxdata"
	code "GoMeeting/pkg/result"
	"GoMeeting/pkg/structs"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMeetingMembersInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMeetingMembersInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeetingMembersInfoLogic {
	return &GetMeetingMembersInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeetingMembersInfoLogic) GetMeetingMembersInfo(req *types.GetMeetingMembersInfoReq) (resp *types.Result, err error) {
	// 根据成员ID从 Redis Hash 中获取会议ID
	meetingId, err := l.svcCtx.Redis.HgetCtx(l.ctx, ctxdata.OnMeetingUserPrefix, strconv.FormatUint(req.UserId, 10))
	if err != nil {
		l.Logger.Errorf("Failed to search member ID in RedisHash: %v", err)
		return types.NewErrorResultWithCode(code.ErrRedisOpCode), nil
	}
	if meetingId == "" {
		return types.NewErrorResultWithCode(code.UserNotInMeetingCode), nil
	}

	// 获取会议中所有成员的详细信息
	fmt.Println("会议ID:", meetingId)
	memberStatusKey := fmt.Sprintf(ctxdata.MeetingMemberDetailPrefix, meetingId)
	// 从 Redis Hash 中获取所有成员信息
	allMembersData, err := l.svcCtx.Redis.HgetallCtx(l.ctx, memberStatusKey)
	if err != nil {
		l.Logger.Errorf("Failed to get all member status from RedisHash: %v", err)
		return types.NewErrorResultWithCode(code.ErrRedisOpCode), nil
	}

	fmt.Println("所有成员信息:", allMembersData)

	// 处理所有成员信息
	var members []structs.MemberStatus
	for userId, statusJson := range allMembersData {
		var memberStatus structs.MemberStatus
		err := json.Unmarshal([]byte(statusJson), &memberStatus)
		if err != nil {
			l.Logger.Errorf("Failed to unmarshal member status for user %s: %v", userId, err)
			continue
		}
		members = append(members, memberStatus)
	}
	fmt.Println("会议成员信息:", members)
	return types.NewSuccessDataResult(members), nil
}
