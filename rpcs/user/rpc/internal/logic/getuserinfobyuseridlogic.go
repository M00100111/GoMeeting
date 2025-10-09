package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/user/models"
	"context"

	"GoMeeting/rpcs/user/rpc/internal/svc"
	"GoMeeting/rpcs/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoByUserIdLogic {
	return &GetUserInfoByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoByUserIdLogic) GetUserInfoByUserId(in *user.GetUserInfoByUserIdReq) (*user.GetUserInfoByUserIdResp, error) {
	userInfo, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserId)
	if err != nil && err != models.ErrNotFound {
		return &user.GetUserInfoByUserIdResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	if err == models.ErrNotFound {
		return &user.GetUserInfoByUserIdResp{
			Code: code.UserNotExistCode,
		}, nil
	}
	return &user.GetUserInfoByUserIdResp{
		Code: code.SUCCESSCode,
		UserInfo: &user.UserInfo{
			Index:    userInfo.Id,
			UserId:   userInfo.UserId,
			Email:    userInfo.Email,
			Username: userInfo.Username,
			Sex:      userInfo.Sex,
			Status:   userInfo.Status,
		},
	}, nil
}
