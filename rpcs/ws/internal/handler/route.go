package handler

import (
	"GoMeeting/rpcs/ws/internal/server"
	"GoMeeting/rpcs/ws/internal/svc"
)

func RegisterHandlers(s *server.WsServer, svc *svc.ServiceContext) {
	s.AddRoutes([]server.Route{
		{
			Method:  "ping",
			Handler: PingHandler(svc),
		},
	})
}
