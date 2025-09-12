package logic

import (
	code "GoMeeting/pkg/result"
	"context"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListLogic) GetFriendList(in *social.GetFriendListReq) (*social.GetFriendListResp, error) {
	list, err := l.svcCtx.FriendsModel.FindRowsByUserIndex(l.ctx, in.UserIndex)
	if err != nil {
		l.Logger.Errorf("FriendsModel.FindRowsByUserIndex error: %v", err)
		return &social.GetFriendListResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	var friendList []*social.Friend
	for _, friend := range list {
		friendList = append(friendList, &social.Friend{
			FriendIndex: friend.FriendIndex,
			Comment:     friend.Comment.String,
		})
	}
	return &social.GetFriendListResp{
		Code:       code.SUCCESSCode,
		FriendList: friendList,
	}, nil
}
