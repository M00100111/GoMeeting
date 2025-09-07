package user

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/user/rpc/user"
	"context"
	"runtime/debug"
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

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.Result, err error) {
	//参数验证
	if req.Email == "" || (req.Password == "" && req.Captcha == "") {
		return types.NewErrorResultWithCode(code.ParamErrorCode), nil
	}

	result, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Email:    req.Email,
		Captcha:  req.Captcha,
		Password: req.Password,
	})
	//系统错误
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("UserRpc.Login error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	//业务错误
	if result.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result), nil
	}
	uid, _ := strconv.Atoi(result.UserId)

	return types.NewSuccessDataResult(&types.LoginResp{
		Msg:    result.Msg,
		UserId: int64(uid),
		Token:  result.Token,
		Expire: result.Expire,
	}), nil
}
