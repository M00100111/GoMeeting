package logic

import (
	code "GoMeeting/pkg/result"
	"context"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupRequestListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupRequestListLogic {
	return &GetGroupRequestListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupRequestListLogic) GetGroupRequestList(in *social.GetGroupRequestListReq) (*social.GetGroupRequestListResp, error) {
	list, err := l.svcCtx.GroupRequestsModel.FindRowsByGroupIndex(l.ctx, in.GroupIndex)

	if err != nil {
		l.Logger.Errorf("GetGroupRequestList GroupRequestsModel.FindRowsByGroupIndex error: %v", err)
		return &social.GetGroupRequestListResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	var requestList []*social.GroupRequest
	for _, request := range list {
		requestList = append(requestList, &social.GroupRequest{
			ReqId:        request.ReqId,
			UserIndex:    request.UserIndex,
			ReqMsg:       request.ReqMsg.String,
			ReqTime:      request.ReqTime.Unix(),
			HandlerIndex: request.HandlerIndex,
			HandleResult: request.HandleResult,
			HandleMsg:    request.HandleMsg.String,
			HandleTime:   request.HandleTime.Time.Unix(),
		})
	}
	return &social.GetGroupRequestListResp{
		Code:             code.SUCCESSCode,
		GroupRequestList: requestList,
	}, nil
}
