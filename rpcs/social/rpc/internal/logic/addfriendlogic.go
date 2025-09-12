package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/pkg/rnum"
	"GoMeeting/rpcs/social/models"
	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type AddFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendLogic {
	return &AddFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddFriendLogic) AddFriend(in *social.AddFriendReq) (*social.AddFriendResp, error) {
	reqIdStr := rnum.GenerateNumber(20)
	reqId, _ := strconv.ParseUint(reqIdStr, 10, 64)
	_, err := l.svcCtx.FriendRequestsModel.Insert(l.ctx, &models.FriendRequests{
		ReqId:       reqId,
		UserIndex:   in.UserIndex,
		FriendIndex: in.FriendIndex,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
	})
	// 可以使用 result 获取插入的ID等信息
	// lastInsertId, _ := result.LastInsertId()
	//受影响的行数
	//result.RowsAffected()

	//发布WS通知服务(待实现)

	if err != nil {
		l.Logger.Errorf("Failed to insert friend request in Mysql: %v", err)
		return &social.AddFriendResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	return &social.AddFriendResp{
		Code: code.SUCCESSCode,
	}, nil
}
