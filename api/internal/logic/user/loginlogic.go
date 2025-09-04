package user

import (
	"GoMeeting/rpcs/user/rpc/user"
	"context"
	"strconv"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	result, err := l.svcCtx.Login(l.ctx, &user.LoginReq{
		Email:    req.Email,
		Captcha:  req.Captcha,
		Password: req.Password,
	})
	uid, _ := strconv.Atoi(result.UserId)
	resp = &types.LoginResp{
		Msg:    result.Msg,
		UserId: int64(uid),
		Token:  result.Token,
		Expire: result.Expire,
	}
	return resp, err
}
