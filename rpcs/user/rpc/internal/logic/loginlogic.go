package logic

import (
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/pkg/md5"
	code "GoMeeting/pkg/result"
	"context"
	"database/sql"
	"runtime/debug"
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
	//查询邮箱是否存在
	u, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("查询用户失败: %v, stack: %s", err, debug.Stack())
		return &user.LoginResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	if err == sql.ErrNoRows {
		l.Logger.Errorf("用户不存在:%v", in.Email)
		return &user.LoginResp{
			Code: code.UserNotExistCode,
		}, nil
	}
	//验证密码
	if in.Password != "" && u.Password != md5.Encrypt(in.Password) {
		l.Logger.Errorf("密码错误:%v", in.Password)
		return &user.LoginResp{
			Code: code.UserPasswordErrorCode,
		}, nil
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
			// 记录错误日志
			l.Logger.Errorf("获取Redis中的验证码失败: %v, stack: %s", err, debug.Stack())
			return &user.LoginResp{
				Code: code.ErrRedisOpCode,
			}, nil
		}
		if value == "" {
			logx.Errorf("验证码已过期或不存在")
			return &user.LoginResp{
				Code: code.CaptchaExpireCode,
			}, nil
		}
		if value != in.Captcha {
			logx.Errorf("验证码错误")
			return &user.LoginResp{
				Code: code.CaptchaErrorCode,
			}, nil
		}
		defer func() {
			_, err := l.svcCtx.Redis.Del(key)
			if err != nil {
				l.Logger.Errorf("删除验证码失败: %v", err)
			}
		}()
	}

	//修改用户上次登录时间
	u.LastLoginTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	err = l.svcCtx.UserModel.Update(l.ctx, u)
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("更新用户登录时间失败: %v, stack: %s", err, debug.Stack())
		return &user.LoginResp{
			Code: code.ErrDbOpCode,
		}, nil
	}

	//生成token并返回
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, strconv.FormatUint(u.UserId, 10))
	if err != nil {
		l.Logger.Errorf("生成token失败: %v, stack: %s", err, debug.Stack())
		return &user.LoginResp{
			Code: code.TokenGenerateErrorCode,
		}, nil
	}
	return &user.LoginResp{
		Code:   code.SUCCESSCode,
		Token:  token,
		Expire: now + l.svcCtx.Config.Jwt.AccessExpire,
		UserId: strconv.FormatUint(u.UserId, 10),
	}, nil
}
