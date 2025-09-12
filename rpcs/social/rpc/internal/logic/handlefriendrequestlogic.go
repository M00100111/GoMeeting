package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/social/models"
	"context"
	"database/sql"
	"errors"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleFriendRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleFriendRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleFriendRequestLogic {
	return &HandleFriendRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HandleFriendRequestLogic) HandleFriendRequest(in *social.HandleFriendRequestReq) (*social.HandleFriendRequestResp, error) {
	result, err := l.svcCtx.FriendRequestsModel.FindOneByReqId(l.ctx, in.ReqId)
	if err != nil && err != models.ErrNotFound {
		l.Logger.Errorf("Failed to find friend request in Mysql: %v", err)
		return &social.HandleFriendRequestResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	if err == models.ErrNotFound {
		return &social.HandleFriendRequestResp{
			Code: code.SocialFriendRequestNotExistCode,
		}, nil
	}
	result.HandleResult = in.HandleResult
	if in.HandleMsg != "" {
		result.HandleMsg = sql.NullString{
			String: in.HandleMsg,
			Valid:  true,
		}
	}
	err = l.svcCtx.FriendRequestsModel.Update(l.ctx, result)
	if err != nil {
		l.Logger.Errorf("Failed to update friend request in Mysql: %v", err)
		return &social.HandleFriendRequestResp{
			Code: code.ErrDbOpCode,
		}, nil
	}

	//不同意直接返回
	if in.HandleResult == HandleResultUnAccept {
		return &social.HandleFriendRequestResp{
			Code: code.SUCCESSCode,
		}, nil
	}

	//发布添加好友任务(待实现)
	//添加好友关系记录,冗余存储
	err = createFriend(l, result.UserIndex, result.FriendIndex)
	if err != nil {
		l.Logger.Errorf("Failed to create friend in Mysql: %v", err)
		return &social.HandleFriendRequestResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	err = createFriend(l, result.FriendIndex, result.UserIndex)
	if err != nil {
		l.Logger.Errorf("Failed to create friend in Mysql: %v", err)
		return &social.HandleFriendRequestResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	//发布WS会话任务通知(待实现)

	return &social.HandleFriendRequestResp{
		Code: code.SUCCESSCode,
	}, nil
}

// 调用同一个服务内的 CreateFriend 方法
func createFriend(l *HandleFriendRequestLogic, userIndex, friendIndex uint64) error {
	// 调用同一个服务内的 createFriend 方法
	creatReq := &social.CreateFriendReq{
		UserIndex:   userIndex,
		FriendIndex: friendIndex,
	}
	// 创建 CreateFriendLogic 实例并调用
	createLogic := NewCreateFriendLogic(l.ctx, l.svcCtx)
	createResult, err := createLogic.CreateFriend(creatReq)
	if err != nil {
		l.Logger.Errorf("调用RPC本地CreateFriend 失败: %v", err)
		return err
	}
	// 检查创建好友记录的结果
	if createResult.Code != code.SUCCESSCode {
		l.Logger.Errorf("调用RPC本地CreateFriend 失败: %v", err)
		return errors.New("调用RPC本地CreateFriend创建好友关系失败")
	}
	return nil
}
