package user

import (
	code "GoMeeting/pkg/result"
	"GoMeeting/rpcs/meeting/rpc/meeting"
	"GoMeeting/rpcs/user/rpc/user"
	"context"
	"runtime/debug"

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

func (l *SignupLogic) Signup(req *types.SignUpReq) (resp *types.Result, err error) {
	//接收HTTP请求
	//参数格式校验
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Captcha == "" {
		return types.NewErrorResultWithCode(code.ParamErrorCode), nil
	}
	//调用RPC服务
	//统一错误处理和响应格式转换
	result, err := l.svcCtx.UserRpc.SignUp(l.ctx, &user.SignUpReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Sex:      req.Sex,
		Avatar:   req.Avatar,
		Captcha:  req.Captcha,
	})
	//系统错误
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("UserRpc.SignUp error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	//业务错误
	if result.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result), nil
	}

	//创建对应的会议记录与会议成员记录
	_, err = l.svcCtx.MeetingRpc.CreateMeeting(l.ctx, &meeting.CreateMeetingReq{
		UserId:    result.Id,
		Username:  req.Username,
		MeetingId: result.MeetingId,
	})
	//系统错误
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("MeetingRpc.CreateMeeting error: %v, stack: %s", err, debug.Stack())
		return types.NewSystemErrorResult(), nil
	}
	//业务错误
	if result.Code != code.SUCCESSCode {
		return types.NewErrorRpcResult(result), nil
	}

	return types.NewSuccessResult(), err
}
