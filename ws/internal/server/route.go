package server

import (
	"GoMeeting/ws/internal/message"
)

type Route struct {
	Method  string
	Handler WsHandlerFunc
}

// 请求处理函数(全局服务对象Server、请求方连接对象、请求消息对象)
type WsHandlerFunc func(srv *WsServer, conn *WsConn, msg *message.Message)
