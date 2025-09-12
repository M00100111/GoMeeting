package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/pkg/rnum"
	"GoMeeting/rpcs/social/models"
	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"
	"context"
	"database/sql"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupMemberRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupMemberRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupMemberRequestLogic {
	return &CreateGroupMemberRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupMemberRequestLogic) CreateGroupMemberRequest(in *social.CreateGroupMemberRequestReq) (*social.CreateGroupMemberRequestResp, error) {
	reqIdStr := rnum.GenerateNumber(20)
	reqId, _ := strconv.ParseUint(reqIdStr, 10, 64)
	request := &models.GroupRequests{
		ReqId:      reqId,
		UserIndex:  in.UserIndex,
		GroupIndex: in.GroupIndex,
	}
	if in.ReqMsg != "" {
		request.ReqMsg = sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		}
	}
	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, request)
	if err != nil {
		l.Logger.Errorf("CreateGroupMemberRequest GroupRequestsModel.Insert err: %v", err)
		return &social.CreateGroupMemberRequestResp{
			Code: code.ErrDbOpCode,
		}, nil
	}

	return &social.CreateGroupMemberRequestResp{
		Code: code.SUCCESSCode,
	}, nil
}
