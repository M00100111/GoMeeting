package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/social/models"
	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupMemberLogic {
	return &CreateGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupMemberLogic) CreateGroupMember(in *social.CreateGroupMemberReq) (*social.CreateGroupMemberResp, error) {
	_, err := l.svcCtx.GroupMembersModel.Insert(l.ctx, &models.GroupMembers{
		GroupIndex: in.GroupIndex,
		UserIndex:  in.UserIndex,
	})
	if err != nil {
		l.Logger.Errorf("CreateGroupMember GroupMembersModel.Insert error: %v", err)
		return &social.CreateGroupMemberResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	return &social.CreateGroupMemberResp{
		Code: code.SUCCESSCode,
	}, nil
}
