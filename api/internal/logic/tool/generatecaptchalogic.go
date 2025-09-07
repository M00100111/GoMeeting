package tool

import (
	"GoMeeting/pkg/captcha"
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/pkg/email"
	code "GoMeeting/pkg/result"
	"context"
	"runtime/debug"
	"time"

	"GoMeeting/api/internal/svc"
	"GoMeeting/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const CAPTCHA_EXPIRE_TIME = time.Minute * 5

type GenerateCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请生成验证码发送到邮箱并存入Redis
func NewGenerateCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateCaptchaLogic {
	return &GenerateCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateCaptchaLogic) GenerateCaptcha(req *types.CaptchaReq) (resp *types.Result, err error) {
	value := captcha.GenerateCaptcha()
	key := ctxdata.CAPTCHA_KEY_PREFIX + req.Email

	err = email.SendEmail(req.Email, value)
	//业务错误
	if err != nil {
		l.Logger.Errorf("发送验证码到邮箱 %v 失败: %v, stack: %s", req.Email, err, debug.Stack())
		return types.NewErrorResultWithCodef(code.CaptchaSendFailCode, req.Email), nil
	}

	// 设置单个键值对并设置过期时间（5分钟）
	err = l.svcCtx.Redis.Setex(key, value, int(CAPTCHA_EXPIRE_TIME.Seconds()))
	if err != nil {
		l.Logger.Errorf("设置邮箱 %v 的验证码到redis中失败: %v, stack: %s", req.Email, err, debug.Stack())
		return types.NewErrorResultWithCodef(code.ErrRedisOpCode, req.Email), nil
	}

	return types.NewSuccessMessageResult("已发送验证码到邮箱" + req.Email), nil
}
