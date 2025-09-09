package logic

import (
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/pkg/md5"
	code "GoMeeting/pkg/result"
	"GoMeeting/pkg/rnum"
	"GoMeeting/rpcs/user/models"
	"GoMeeting/rpcs/user/rpc/internal/svc"
	"GoMeeting/rpcs/user/rpc/user"
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"runtime/debug"
	"strconv"
)

type SignUpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSignUpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignUpLogic {
	return &SignUpLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SignUpLogic) SignUp(in *user.SignUpReq) (*user.SignUpResp, error) {
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
		return &user.SignUpResp{
			Code: code.ErrRedisOpCode,
		}, nil
	}
	if value == "" {
		logx.Errorf("验证码已过期或不存在")
		return &user.SignUpResp{
			Code: code.CaptchaExpireCode,
		}, nil
	}
	if value != in.Captcha {
		logx.Errorf("验证码错误")
		return &user.SignUpResp{
			Code: code.CaptchaErrorCode,
		}, nil
	}
	// 使用 defer 确保在函数结束时删除验证码
	//键存在且删除成功：
	//第一个返回值：1（删除的键数量）
	//第二个返回值：nil（无错误）
	//键不存在：
	//第一个返回值：0（删除的键数量）
	//第二个返回值：nil（无错误）
	//网络或其他错误：
	//第一个返回值：0
	//第二个返回值：具体的错误信息
	defer func() {
		_, err := l.svcCtx.Redis.Del(key)
		if err != nil {
			logx.Errorf("删除验证码失败: %v", err)
		}
	}()

	//2. 分开查询username、email是否已存在
	//找到记录 → err == nil
	//未找到记录 → err == sql.ErrNoRows（正常情况，不是错误）
	//数据库错误 → err 为其他值
	_, err = l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	if err == nil {
		logx.Errorf("用户名已存在")
		return &user.SignUpResp{
			Code: code.UserExistCode,
		}, nil
	}
	if err != nil && err != sql.ErrNoRows {
		// 记录错误日志
		l.Logger.Errorf("查询用户名失败: %v, stack: %s", err, debug.Stack())
		return &user.SignUpResp{
			Code: code.ErrDbOpCode,
		}, nil
	}
	_, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err == nil {
		logx.Errorf("邮箱已被注册")
		return &user.SignUpResp{
			Code: code.EmailExistCode,
		}, nil
	}
	if err != nil && err != sql.ErrNoRows {
		// 记录错误日志
		l.Logger.Errorf("查询邮箱失败: %v, stack: %s", err, debug.Stack())
		return &user.SignUpResp{
			Code: code.ErrDbOpCode,
		}, nil
	}

	//字段初始化
	userId, _ := strconv.Atoi(rnum.GenerateNumber(12))
	//3. password使用MD5加密
	password := md5.Encrypt(in.Password)
	//4. 字段初始化
	//未显式设置的字段会使用Go语言的零值，同时数据库层面可能也有默认值设置
	_, err = l.svcCtx.UserModel.Insert(l.ctx, &models.User{
		UserId:   uint64(userId),
		Username: in.Username,
		Password: password,
		Email:    in.Email,
		Sex:      uint64(in.Sex),
	})
	if err != nil && err != sql.ErrNoRows {
		// 记录错误日志
		l.Logger.Errorf("新增用户信息失败: %v, stack: %s", err, debug.Stack())
		return &user.SignUpResp{
			Code: code.ErrDbOpCode,
		}, nil
	}

	u, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, uint64(userId))
	if err != nil {
		// 记录错误日志
		l.Logger.Errorf("查询用户主键失败: %v, stack: %s", err, debug.Stack())
		return &user.SignUpResp{
			Code: code.ErrDbOpCode,
		}, nil
	}

	return &user.SignUpResp{
		Code:      code.SUCCESSCode,
		Id:        u.Id,
		MeetingId: u.UserId,
	}, nil
}
