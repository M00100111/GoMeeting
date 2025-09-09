package logic

import (
	code "GoMeeting/pkg/result"
	"context"

	"GoMeeting/rpcs/user/rpc/internal/svc"
	"GoMeeting/rpcs/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIndexByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIndexByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIndexByUserIdLogic {
	return &GetIndexByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetIndexByUserIdLogic) GetIndexByUserId(in *user.GetIndexByUserIdReq) (*user.GetIndexByUserIdResp, error) {
	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil {
		return &user.GetIndexByUserIdResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	return &user.GetIndexByUserIdResp{
		Code:  code.SUCCESSCode,
		Index: u.Id,
	}, nil
}
