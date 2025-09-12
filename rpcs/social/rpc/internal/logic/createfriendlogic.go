package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/social/models"
	"context"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFriendLogic {
	return &CreateFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateFriendLogic) CreateFriend(in *social.CreateFriendReq) (*social.CreateFriendResp, error) {
	friend := &models.Friends{
		UserIndex:   in.UserIndex,
		FriendIndex: in.FriendIndex,
	}
	_, err := l.svcCtx.FriendsModel.Insert(l.ctx, friend)
	if err != nil {
		l.Logger.Errorf("insert friends error: %v", err)
		return &social.CreateFriendResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	return &social.CreateFriendResp{
		Code: code.SUCCESSCode,
	}, nil
}
