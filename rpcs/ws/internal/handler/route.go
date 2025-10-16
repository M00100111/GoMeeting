package handler

import (
	"GoMeeting/pkg/structs/message"
	"GoMeeting/rpcs/ws/internal/handler/chat"
	"GoMeeting/rpcs/ws/internal/handler/notification"
	"GoMeeting/rpcs/ws/internal/handler/ping"
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
		{
			Method:  message.Meeting_Member_Join_Notice_Method,
			Handler: notification.MeetingMemberJoinNoticeHandler(),
		},
	})
}
