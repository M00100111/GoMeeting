package types

import (
	"GoMeeting/pkg/result"
	"fmt"
	"reflect"
)

// API通用结果类
type Result struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func newResult(code int32, msg string, data interface{}) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// 返回通用成功结果
func NewSuccessResult() *Result {
	return newResult(code.Success.GetCode(), code.Success.GetMsg(), nil)
}
func NewSuccessDataResult(data interface{}) *Result {
	return newResult(code.Success.GetCode(), code.Success.GetMsg(), data)
}

// 返回自定义成功信息结果
func NewSuccessMessageResult(msg string) *Result {
	return newResult(code.Success.GetCode(), msg, nil)
}

// 返回自定义格式成功信息结果
func NewSuccessMessageResultf(format string, a ...interface{}) *Result {
	return NewSuccessMessageResult(fmt.Sprintf(format, a...))
}

// 返回自定义错误信息结果
// 用于不常出现的错误,没必要专门注册错误码
func NewErrorResult(msg string) *Result {
	return newResult(code.ERRORCode, msg, nil)
}

func NewErrorResultf(format string, a ...interface{}) *Result {
	return NewErrorResult(fmt.Sprintf(format, a...))
}

// 根据传入错误码类型返回对应的结果
func NewErrorResultWithCode(c int32) *Result {
	result := code.GetCodeResult(c)
	return newResult(result.GetCode(), result.GetMsg(), nil)
}

// 根据传入错误码类型返回对应的结果
func NewErrorResultWithCodef(c int32, element interface{}) *Result {
	result := code.GetCodeResult(c)
	return newResult(result.GetCode(), result.GetMsg()+fmt.Sprintf(":%v", element), nil)
}

// 根据rpc的错误返回结果，返回对应的通用结果
func NewErrorRpcResult(ptr interface{}) *Result {
	//获取结构体指针对应的值对象
	valueOfPtr := reflect.ValueOf(ptr)
	//检查是否为指针类型
	if valueOfPtr.Kind() != reflect.Ptr {
		panic("ptr must be a pointer")
	}
	//获取结构体指针所指向的结构体对应的值对象
	pValueOfPtr := valueOfPtr.Elem()
	// 检查是否为结构体类型
	if pValueOfPtr.Kind() != reflect.Struct {
		panic("指针必须指向结构体")
	}

	//获取结构体属性所对应的字段的值对象
	ValueOfCode := pValueOfPtr.FieldByName("Code")
	// 安全地获取 Code 字段值
	// 如果字段不存在或类型不正确直接panic
	if !ValueOfCode.IsValid() {
		panic("field 'Code' not found or invalid")
	}
	if ValueOfCode.Kind() != reflect.Int32 {
		panic("field 'Code' is not int32 type")
	}
	c := int32(ValueOfCode.Int())

	//获取结构体属性所对应的字段的值对象
	ValueOfMsg := pValueOfPtr.FieldByName("Msg")
	// 如果字段不存在或类型不正确直接panic
	if !ValueOfMsg.IsValid() {
		panic("field 'Msg' not found or invalid")
	}
	if ValueOfMsg.Kind() != reflect.String {
		panic("field 'Msg' is not string type")
	}
	m := ValueOfMsg.String()

	if c == code.ERRORCode {
		if m == "" {
			panic("msg is empty")
		}
		// 通用错误代码
		return NewErrorResult(m)
	} else {
		if m != "" {
			return NewErrorResultWithCodef(c, m)
		}
		// 注册的错误代码
		return NewErrorResultWithCode(c)
	}
}

// 返回系统错误结果
func NewSystemErrorResult() *Result {
	return newResult(code.SysError.GetCode(), code.SysError.GetMsg(), nil)
}
