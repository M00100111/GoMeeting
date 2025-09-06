package types

import "fmt"

// 通用结果类
const (
	SUCCESS = 200
	ERROR   = 500
)

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

func NewSuccessResult(data interface{}) *Result {
	return newResult(SUCCESS, "success", data)
}

func NewSuccessMessageResult(msg string) *Result {
	return newResult(SUCCESS, msg, nil)
}

func NewSuccessMessageResultf(format string, a ...interface{}) *Result {
	return NewSuccessMessageResult(fmt.Sprintf(format, a...))
}

func NewErrorResult(msg string) *Result {
	return newResult(ERROR, msg, nil)
}

func NewErrorResultf(format string, a ...interface{}) *Result {
	return NewErrorResult(fmt.Sprintf(format, a...))
}
