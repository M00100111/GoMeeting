package user

import (
	"GoMeeting/rpcs/user/rpc/user"
	"context"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignupLogic) Signup(req *types.SignUpReq) (resp *types.SignUpResp, err error) {
	result, err := l.svcCtx.User.SignUp(l.ctx, &user.SignUpReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Sex:      req.Sex,
		Avatar:   req.Avatar,
	})
	resp = &types.SignUpResp{
		Msg: result.Msg,
	}
	return resp, err
}
