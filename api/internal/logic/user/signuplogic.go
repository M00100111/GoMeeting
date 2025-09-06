package user

import (
	"GoMeeting/rpcs/meeting/rpc/meeting"
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

func (l *SignupLogic) Signup(req *types.SignUpReq) (resp *types.Result, err error) {
	result, err := l.svcCtx.UserRpc.SignUp(l.ctx, &user.SignUpReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Sex:      req.Sex,
		Avatar:   req.Avatar,
		Captcha:  req.Captcha,
	})
	if err != nil {
		return &types.Result{
			Msg: result.Msg,
		}, err
	}

	//创建对应的会议记录与会议成员记录
	_, err = l.svcCtx.MeetingRpc.CreateMeeting(l.ctx, &meeting.CreateMeetingReq{
		UserId:   result.Id,
		Username: req.Username,
	})
	if err != nil {
		return &types.Result{
			Msg: "创建会议信息和用户成员信息失败",
		}, err
	}

	resp = &types.Result{
		Msg: result.Msg,
	}
	return resp, err
}
