package logic

import (
	code "GoMeeting/pkg/result"
	"context"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupListLogic {
	return &GetGroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupListLogic) GetGroupList(in *social.GetGroupListReq) (*social.GetGroupListResp, error) {
	list, err := l.svcCtx.GroupsModel.FindRowsByUserIndex(l.ctx, in.UserIndex)

	if err != nil {
		l.Logger.Errorf("GetGroupList GroupsModel.FindRowsByUserIndex error: %v", err)
		return &social.GetGroupListResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	var groupList []*social.Group
	for _, group := range list {
		groupList = append(groupList, &social.Group{
			GroupId:   group.GroupId,
			GroupName: group.GroupName,
		})
	}
	return &social.GetGroupListResp{
		Code:      code.SUCCESSCode,
		GroupList: groupList,
	}, nil
}
