package handler

import (
	"GoMeeting/ws/internal/server"
	"GoMeeting/ws/internal/svc"
)

func RegisterHandlers(s *server.WsServer, svc *svc.ServiceContext) {
	s.AddRoutes([]server.Route{
		{
			Method:  "ping",
			Handler: PingHandler(svc),
		},
	})
}
