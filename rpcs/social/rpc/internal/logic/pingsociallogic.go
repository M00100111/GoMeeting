package logic

import (
	code "GoMeeting/pkg/result"
	"context"

	"GoMeeting/rpcs/social/rpc/internal/svc"
	"GoMeeting/rpcs/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingSocialLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingSocialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingSocialLogic {
	return &PingSocialLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingSocialLogic) PingSocial(in *social.PingReq) (*social.Result, error) {
	return &social.Result{
		Code: code.SUCCESSCode,
		Msg:  in.Msg + "!!!",
	}, nil
}
