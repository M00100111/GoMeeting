package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinGroupLogic) JoinGroup(in *social.JoinGroupReq) (*social.JoinGroupResp, error) {
	group, err := l.svcCtx.GroupsModel.FindOneByGroupId(l.ctx, in.GroupId)
	if err != nil {
		l.Logger.Errorf("GroupsModel.FindOneByGroupId error: %v", err)
		return &social.JoinGroupResp{
			Code: code.ErrDbOpCode,
		}, err
	}
	//开放直接加入群组
	if group.JoinStatus == 0 {
		err = l.svcCtx.CallCreateGroupMember(l.ctx, group.Id, in.UserIndex)
		if err != nil {
			l.Logger.Errorf("JoinGroup CallCreateGroupMember error: %v", err)
			return &social.JoinGroupResp{
				Code: code.ErrDbOpCode,
			}, err
		}
	} else { // 提交申请
		err = l.svcCtx.CallCreateGroupMemberRequest(l.ctx, group.Id, in.UserIndex, in.ReqMsg)
		if err != nil {
			l.Logger.Errorf("JoinGroup CallCreateGroupMemberRequest error: %v", err)
			return &social.JoinGroupResp{
				Code: code.ErrDbOpCode,
			}, err
		}
	}

	return &social.JoinGroupResp{
		Code: code.SUCCESSCode,
	}, nil
}
