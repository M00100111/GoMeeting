package logic

import (
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/pkg/md5"
	"GoMeeting/pkg/rnum"
	"GoMeeting/rpcs/user/models"
	"GoMeeting/rpcs/user/rpc/internal/svc"
	"GoMeeting/rpcs/user/rpc/user"
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
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
		logx.Errorf("获取失败: %v", err)
		return &user.SignUpResp{
			Msg: "获取验证码失败",
		}, err
	}
	if value == "" {
		logx.Errorf(",验证码已过期或不存在")
		return &user.SignUpResp{
			Msg: "验证码已过期或不存在",
		}, err
	}
	if value != in.Captcha {
		logx.Errorf("验证码错误")
		return &user.SignUpResp{
			Msg: "验证码错误",
		}, err
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
			Msg: "用户名已存在",
		}, err
	}
	if err != nil && err != sql.ErrNoRows {
		logx.Errorf("查询用户名失败: %v", err)
		return &user.SignUpResp{
			Msg: "查询用户名失败",
		}, err
	}
	_, err = l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err == nil {
		logx.Errorf("邮箱已注册")
		return &user.SignUpResp{
			Msg: "邮箱已注册",
		}, err
	}
	if err != nil && err != sql.ErrNoRows {
		logx.Errorf("查询邮箱失败: %v", err)
		return &user.SignUpResp{
			Msg: "查询邮箱失败",
		}, err
	}

	//字段初始化
	meetId := rnum.GenerateNumber(12)
	userId, _ := strconv.Atoi(meetId)
	//3. password使用MD5加密
	password := md5.Encrypt(in.Password)
	//4. 字段初始化
	//未显式设置的字段会使用Go语言的零值，同时数据库层面可能也有默认值设置
	l.svcCtx.UserModel.Insert(l.ctx, &models.User{
		UserId:    uint64(userId),
		Username:  in.Username,
		Password:  password,
		Email:     in.Email,
		MeetingId: meetId,
		Sex:       uint64(in.Sex),
	})

	return &user.SignUpResp{
		Msg: "注册成功",
	}, nil
}
