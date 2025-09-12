package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/pkg/rnum"
	"GoMeeting/rpcs/social/models"
	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"
	"context"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupLogic) CreateGroup(in *social.CreateGroupReq) (*social.CreateGroupResp, error) {
	groupIdStr := rnum.GenerateNumber(12)
	groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
	_, err := l.svcCtx.GroupsModel.Insert(l.ctx, &models.Groups{
		GroupId:    groupId,
		GroupName:  in.GroupName,
		UserIndex:  in.UserIndex,
		JoinStatus: in.JoinStatus,
	})
	if err != nil {
		l.Logger.Errorf("GroupsModel.Insert error: %v", err)
		return &social.CreateGroupResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	return &social.CreateGroupResp{
		Code: code.SUCCESSCode,
	}, nil
}
