package handler

import (
	"GoMeeting/rpcs/ws/internal/handler/chat"
	"GoMeeting/rpcs/ws/internal/handler/notification"
	"GoMeeting/rpcs/ws/internal/handler/ping"
	"GoMeeting/rpcs/ws/internal/message"
	"GoMeeting/rpcs/ws/internal/server"
)

func RegisterHandlers(s *server.WsServer) {
	s.AddRoutes([]server.Route{
		{
			//Method:  "Ping",
			Method:  message.Ping_Method,
			Handler: ping.PingHandler(),
		},
		{
			//Method:  "Chat",
			Method:  message.Chat_Method,
			Handler: chat.ChatHandler(),
		},
		{
			//Method:  "Notification",
			Method:  message.Notification_Method,
			Handler: notification.NotificationHandler(),
		},
		//{
		//	Method:  message.WebRTC_Method,
		//	Handler: webRTC.WebRTCHandler(),
		//},
	})
}
