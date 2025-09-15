package handler

import (
	"GoMeeting/rpcs/ws/internal/handler/chat"
	"GoMeeting/rpcs/ws/internal/handler/notification"
	"GoMeeting/rpcs/ws/internal/handler/ping"
	"GoMeeting/rpcs/ws/internal/server"
)

func RegisterHandlers(s *server.WsServer) {
	s.AddRoutes([]server.Route{
		{
			Method:  "ping",
			Handler: ping.PingHandler(),
		},
		{
			Method:  "chat",
			Handler: chat.ChatHandler(),
		},
		{
			Method:  "notification",
			Handler: notification.NotificationHandler(),
		},
	})
}
