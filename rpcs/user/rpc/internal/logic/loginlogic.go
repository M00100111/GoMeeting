package logic

import (
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/pkg/md5"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"strconv"
	"time"

	"GoMeeting/rpcs/user/rpc/internal/svc"
	"GoMeeting/rpcs/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	//邮箱登录
	if in.Email == "" {
		return &user.LoginResp{
			Msg: "邮箱不能为空",
		}, nil
	}
	//查询邮箱是否存在
	u, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return &user.LoginResp{
			Msg: "查询用户失败",
		}, err
	}
	if err == sql.ErrNoRows {
		logx.Errorf("用户不存在")
		return &user.LoginResp{
			Msg: "用户不存在",
		}, err
	}
	//验证密码
	if in.Password != "" && u.Password != md5.Encrypt(in.Password) {
		logx.Errorf("密码错误")
		return &user.LoginResp{
			Msg: "密码错误",
		}, err
	}
	//验证验证码
	if in.Captcha != "" {
		key := ctxdata.CAPTCHA_KEY_PREFIX + in.Email
		// 验证验证码是否正确
		value, err := l.svcCtx.Redis.Get(key)
		//Redis Get 操作的行为总结
		//键存在且有值：返回实际值和 nil 错误
		//键不存在：返回空字符串 "" 和 nil 错误
		//网络或其他错误：返回 "" 和具体的错误信息
		if err != nil {
			logx.Errorf("获取失败: %v", err)
			return &user.LoginResp{
				Msg: "获取验证码失败",
			}, err
		}
		if value == "" {
			logx.Errorf(",验证码已过期或不存在")
			return &user.LoginResp{
				Msg: "验证码已过期或不存在",
			}, err
		}
		if value != in.Captcha {
			logx.Errorf("验证码错误")
			return &user.LoginResp{
				Msg: "验证码错误",
			}, err
		}
		defer func() {
			_, err := l.svcCtx.Redis.Del(key)
			if err != nil {
				logx.Errorf("删除验证码失败: %v", err)
			}
		}()
	}
	//生成token并返回
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, strconv.FormatUint(u.UserId, 10))
	if err != nil {
		return nil, errors.Wrapf(err, "ctxdata get jwt token err %v")
	}
	return &user.LoginResp{
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
		UserId: strconv.FormatUint(u.UserId, 10),
		Msg:    "登录成功",
	}, nil
}
