package logic

import (
	"context"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleGroupRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleGroupRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleGroupRequestLogic {
	return &HandleGroupRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HandleGroupRequestLogic) HandleGroupRequest(in *social.HandleGroupRequestReq) (*social.HandleGroupRequestResp, error) {
	// todo: add your logic here and delete this line

	return &social.HandleGroupRequestResp{}, nil
}
