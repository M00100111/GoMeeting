package server

import (
	"GoMeeting/pkg/structs/message"
)

type Route struct {
	Method  message.MessageMethod
	Handler WsHandlerFunc
}

// 请求处理函数(全局服务对象Server、请求方连接对象、请求消息对象)
type WsHandlerFunc func(srv *WsServer, msg *message.Message) error
