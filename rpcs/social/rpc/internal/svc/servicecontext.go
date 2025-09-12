package svc

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/social/models"
	"GoMeeting/rpcs/social/rpc/internal/config"
	"GoMeeting/rpcs/social/rpc/internal/logic"
	"GoMeeting/rpcs/social/rpc/social"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	models.FriendsModel
	models.FriendRequestsModel
	models.GroupsModel
	models.GroupMembersModel
	models.GroupRequestsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:              c,
		FriendsModel:        models.NewFriendsModel(sqlConn, c.Cache),
		FriendRequestsModel: models.NewFriendRequestsModel(sqlConn, c.Cache),
		GroupsModel:         models.NewGroupsModel(sqlConn, c.Cache),
		GroupMembersModel:   models.NewGroupMembersModel(sqlConn, c.Cache),
		GroupRequestsModel:  models.NewGroupRequestsModel(sqlConn, c.Cache),
	}
}

// 调用同一个服务内的 CreateGroupMember 方法
func (s *ServiceContext) CallCreateGroupMember(ctx context.Context, groupIndex, userIndex uint64) error {
	// 调用同一个服务内的 createFriend 方法
	creatReq := &social.CreateGroupMemberReq{
		GroupIndex: groupIndex,
		UserIndex:  userIndex,
	}
	// 创建 CreateFriendLogic 实例并调用
	createLogic := logic.NewCreateGroupMemberLogic(ctx, s)
	createResult, err := createLogic.CreateGroupMember(creatReq)
	if err != nil {
		return err
	}
	// 检查创建好友记录的结果
	if createResult.Code != code.SUCCESSCode {
		return errors.New("调用RPC本地CreateGroupMember失败")
	}
	return nil
}

// 调用同一个服务内的 CreateGroupMemberRequest 方法
func (s *ServiceContext) CallCreateGroupMemberRequest(ctx context.Context, groupIndex, userIndex uint64, reqMsg string) error {
	// 调用同一个服务内的 createFriend 方法
	creatReq := &social.CreateGroupMemberRequestReq{
		GroupIndex: groupIndex,
		UserIndex:  userIndex,
		ReqMsg:     reqMsg,
	}
	// 创建 CreateGroupMemberRequestLogic 实例并调用
	createLogic := logic.NewCreateGroupMemberRequestLogic(ctx, s)
	createResult, err := createLogic.CreateGroupMemberRequest(creatReq)
	if err != nil {
		return err
	}
	// 检查创建好友记录的结果
	if createResult.Code != code.SUCCESSCode {
		return errors.New("调用RPC本地CreateGroupMemberRequest失败")
	}
	return nil
}
