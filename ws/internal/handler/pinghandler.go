package handler

import (
	"GoMeeting/ws/internal/message"
	"GoMeeting/ws/internal/server"
	"GoMeeting/ws/internal/svc"
)

func PingHandler(svcCtx *svc.ServiceContext) server.WsHandlerFunc {
	return func(s *server.WsServer, conn *server.WsConn, msg *message.Message) {
		err := s.SendMessage(&message.Message{
			MsgId:       "",
			MessageType: 1,
			Method:      "pong",
			SenderId:    "0",
			Data:        "pong",
		}, conn)
		s.Info("err ", err)
	}
}
