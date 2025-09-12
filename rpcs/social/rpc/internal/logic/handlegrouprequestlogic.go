package logic

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"
	"context"
	"database/sql"

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
	request, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, in.ReqId)
	if err != nil {
		l.Logger.Errorf("Failed to find group request: %v", err)
		return &social.HandleGroupRequestResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	request.HandleResult = in.HandleResult
	request.HandlerIndex = in.HandlerIndex
	if in.HandleMsg != "" {
		request.HandleMsg = sql.NullString{
			String: in.HandleMsg,
			Valid:  true,
		}
	}
	err = l.svcCtx.GroupRequestsModel.Update(l.ctx, request)

	//不接受直接返回
	if in.HandleResult == HandleResultUnAccept {
		return &social.HandleGroupRequestResp{
			Code: code.SUCCESSCode,
		}, nil
	}

	//创建成员关系
	err = l.svcCtx.CallCreateGroupMember(l.ctx, request.GroupIndex, request.UserIndex)
	if err != nil {
		l.Logger.Errorf("HandleGroupRequest CallCreateGroupMember error: %v", err)
		return &social.HandleGroupRequestResp{
			Code: code.ErrDbOpCode,
		}, err
	}

	return &social.HandleGroupRequestResp{
		Code: code.SUCCESSCode,
	}, nil
}
