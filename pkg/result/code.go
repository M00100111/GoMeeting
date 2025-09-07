package code

import "fmt"

var (
	codes = map[int32]Result{}     //标识注册过的错误码，防止重复注册
	msgs  = make(map[int32]string) //错误码对应的错误信息映射
)

func GetCodeResult(code int32) Result {
	r, ok := codes[code]
	if !ok {
		panic("未注册的错误码")
	}
	return r
}

func GetCodeMsg(code int32) string {
	return msgs[code]
}

type Result struct {
	code int32
	msg  string
}

func (r Result) GetCode() int32 {
	return r.code
}

func (r Result) GetMsg() string {
	return r.msg
}

func RegisterCode(code int32, msg string) Result {
	if code == 0 {
		panic("状态码不能为0")
	}
	//判断状态信息是否为空
	if msg == "" {
		panic("状态信息不能为空")
	}
	//判断状态码是否已注册
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("状态码%d已注册", code))
	}
	r := Result{code: code, msg: msg}
	codes[code] = r    //注册状态码
	msgs[code] = r.msg // 保存状态码及其信息
	return r
}

// 状态码
const (
	SUCCESSCode   = 200
	ERRORCode     = 500
	SYS_ERRORCode = 501
)

const (
	ParamErrorCode = 400 //参数错误
	ErrDbOpCode    = 1001
	ErrRedisOpCode = 1002
)

const (
	UserRegisterFailCode  = 10000
	EmailExistCode        = 10001
	UserExistCode         = 10002
	UserNotExistCode      = 10003
	UserLoginFailCode     = 10004
	CaptchaErrorCode      = 10005
	CaptchaExpireCode     = 10006
	CaptchaSendFailCode   = 10007
	UserPasswordErrorCode = 10008
)

const (
	TokenErrorCode         = 1003
	TokenExpireCode        = 1004
	TokenGenerateErrorCode = 1005
)

// 通用状态对象
var (
	Success  = RegisterCode(SUCCESSCode, "success")
	Error    = RegisterCode(ERRORCode, "unknown error")
	SysError = RegisterCode(SYS_ERRORCode, "system error")
)

var (
	ErrDbOp    = RegisterCode(ErrDbOpCode, "数据库操作异常")
	ErrRedisOp = RegisterCode(ErrRedisOpCode, "Redis操作异常")
)

var (
	ParamError        = RegisterCode(ParamErrorCode, "参数错误")
	UserRegisterFail  = RegisterCode(UserRegisterFailCode, "用户注册失败")
	EmailExist        = RegisterCode(EmailExistCode, "邮箱已存在")
	UserExist         = RegisterCode(UserExistCode, "用户已存在")
	UserNotExist      = RegisterCode(UserNotExistCode, "用户不存在")
	UserPasswordError = RegisterCode(UserPasswordErrorCode, "用户密码错误")
	UserLoginFail     = RegisterCode(UserLoginFailCode, "用户登录失败")

	UserCaptchaFail     = RegisterCode(CaptchaErrorCode, "验证码错误")
	UserCaptchaExpire   = RegisterCode(CaptchaExpireCode, "验证码已过期")
	UserCaptchaSendFail = RegisterCode(CaptchaSendFailCode, "验证码发送失败")
)
var (
	TokenError         = RegisterCode(TokenErrorCode, "token错误")
	TokenGenerateError = RegisterCode(TokenGenerateErrorCode, "token生成错误")
)
