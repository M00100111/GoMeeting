package logic

import (
	code "GoMeeting/pkg/result"
	"context"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendRequestListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendRequestListLogic {
	return &GetFriendRequestListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendRequestListLogic) GetFriendRequestList(in *social.GetFriendRequestListReq) (*social.GetFriendRequestListResp, error) {
	//获取待处理的好友请求
	list, err := l.svcCtx.FriendRequestsModel.FindRowsByFriendIndexAndHandleResult(l.ctx, in.UserIndex, 0)
	if err != nil {
		l.Logger.Errorf("Failed to find unhandled friend requests in Mysql: %v", err)
		return &social.GetFriendRequestListResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	var friendRequestslist []*social.FriendRequest
	for _, result := range list {
		friendRequestslist = append(friendRequestslist, &social.FriendRequest{
			FriendIndex: result.UserIndex,
			ReqMsg:      result.ReqMsg.String,
			ReqTime:     result.ReqTime.Unix(),
		})
	}
	return &social.GetFriendRequestListResp{
		Code:              code.SUCCESSCode,
		FriendRequestList: friendRequestslist,
	}, nil
}
