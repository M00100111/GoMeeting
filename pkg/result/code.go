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
	ParamErrorCode = 400
)

const ( //参数错误
	ErrDbOpCode = 1001 + iota
	ErrRedisOpCode
	ErrRecordNotExistCode
	TokenErrorCode
	TokenExpireCode
	TokenGenerateErrorCode
)

const (
	UserRegisterFailCode = 10000 + iota
	EmailExistCode
	UserExistCode
	UserNotExistCode
	UserLoginFailCode
	CaptchaErrorCode
	CaptchaExpireCode
	CaptchaSendFailCode
	UserPasswordErrorCode
	UserNoPermissionCode
)

// 会议相关错误码
const (
	MeetingUserInOtherMeetingCode = 20000 + iota
	MeetingUserNotAllowedCode
	MeetingJoinCodeErrorCode
	MeetingNotExistCode
	MeetingAlreadyStartedCode
	MeetingAlreadyEndedCode
	MeetingStartFailCode
	MeetingNotStartedCode
	MeetingEndFailCode
)

// 通用状态对象
var (
	Success  = RegisterCode(SUCCESSCode, "success")
	Error    = RegisterCode(ERRORCode, "unknown error")
	SysError = RegisterCode(SYS_ERRORCode, "system error")
)

var (
	ParamError         = RegisterCode(ParamErrorCode, "参数错误")
	ErrDbOp            = RegisterCode(ErrDbOpCode, "数据库操作异常")
	ErrRedisOp         = RegisterCode(ErrRedisOpCode, "Redis操作异常")
	ErrRecordNotExist  = RegisterCode(ErrRecordNotExistCode, "记录不存在")
	TokenError         = RegisterCode(TokenErrorCode, "token错误")
	TokenGenerateError = RegisterCode(TokenGenerateErrorCode, "token生成错误")
)

var (
	UserRegisterFail  = RegisterCode(UserRegisterFailCode, "用户注册失败")
	EmailExist        = RegisterCode(EmailExistCode, "邮箱已存在")
	UserExist         = RegisterCode(UserExistCode, "用户已存在")
	UserNotExist      = RegisterCode(UserNotExistCode, "用户不存在")
	UserPasswordError = RegisterCode(UserPasswordErrorCode, "用户密码错误")
	UserLoginFail     = RegisterCode(UserLoginFailCode, "用户登录失败")

	UserNoPermission = RegisterCode(UserNoPermissionCode, "用户无此权限")

	UserCaptchaFail     = RegisterCode(CaptchaErrorCode, "验证码错误")
	UserCaptchaExpire   = RegisterCode(CaptchaExpireCode, "验证码已过期")
	UserCaptchaSendFail = RegisterCode(CaptchaSendFailCode, "验证码发送失败")
)

// 会议相关错误信息
var (
	MeetingUserInOtherMeeting = RegisterCode(MeetingUserInOtherMeetingCode, "用户已加入其他会议")
	MeetingUserNotAllowed     = RegisterCode(MeetingUserNotAllowedCode, "用户不被允许加入会议")
	MeetingJoinCodeError      = RegisterCode(MeetingJoinCodeErrorCode, "入会码错误")
	MeetingNotExist           = RegisterCode(MeetingNotExistCode, "会议不存在")
	MeetingAlreadyStarted     = RegisterCode(MeetingAlreadyStartedCode, "会议已开始")
	MeetingNotStarted         = RegisterCode(MeetingNotStartedCode, "会议未开始")
	MeetingAlreadyEnded       = RegisterCode(MeetingAlreadyEndedCode, "会议已结束")
	MeetingStartFail          = RegisterCode(MeetingStartFailCode, "会议开启失败")
	MeetingEndFail            = RegisterCode(MeetingEndFailCode, "结束会议失败")
)
