package logic

import (
	code "GoMeeting/pkg/result"
	"context"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(in *social.GetGroupMemberListReq) (*social.GetGroupMemberListResp, error) {
	list, err := l.svcCtx.GroupMembersModel.FindRowsByGroupIndex(l.ctx, in.GroupIndex)

	if err != nil {
		l.Logger.Errorf("GetGroupMemberList GroupMembersModel.FindRowsByGroupIndex error: %v", err)
		return &social.GetGroupMemberListResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	var memberList []*social.GroupMember
	for _, member := range list {
		memberList = append(memberList, &social.GroupMember{
			UserIndex:  member.UserIndex,
			UserType:   member.UserType,
			UserStatus: member.UserStatus,
		})
	}
	return &social.GetGroupMemberListResp{
		Code:            code.SUCCESSCode,
		GroupMemberList: memberList,
	}, nil
}
