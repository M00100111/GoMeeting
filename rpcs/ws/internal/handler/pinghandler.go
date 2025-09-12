package handler

import (
	"GoMeeting/rpcs/ws/internal/message"
	"GoMeeting/rpcs/ws/internal/server"
	"GoMeeting/rpcs/ws/internal/svc"
)

func PingHandler(svcCtx *svc.ServiceContext) server.WsHandlerFunc {
	return func(s *server.WsServer, conn *server.WsConn, msg *message.Message) {
		err := s.SendMessage(&message.Message{
			MessageType: 1,
			Method:      "pong",
			SenderId:    "0",
			Data:        "pong",
		}, conn)
		s.Info("err ", err)
	}
}
