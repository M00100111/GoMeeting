package logic

import (
	"context"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupRequestListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupRequestListLogic {
	return &GetGroupRequestListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupRequestListLogic) GetGroupRequestList(in *social.GetGroupRequestListReq) (*social.GetGroupRequestListResp, error) {
	// todo: add your logic here and delete this line

	return &social.GetGroupRequestListResp{}, nil
}
