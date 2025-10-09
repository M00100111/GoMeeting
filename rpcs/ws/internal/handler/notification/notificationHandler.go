package notification

import (
	"GoMeeting/rpcs/ws/internal/message"
	"GoMeeting/rpcs/ws/internal/server"
	"encoding/json"
)

func NotificationHandler() server.WsHandlerFunc {
	return func(s *server.WsServer, msg *message.Message) error {
		var notiMsg message.NotificationData
		if err := json.Unmarshal(msg.Data, &notiMsg); err != nil {
			s.Logger.Errorf("Failed to unmarshal message.Data to message.NotificationMessage: %v", err)
			return err
		}
		conn := s.GetWsConnByUid(notiMsg.ReceiverId)
		err := s.SendMessage(&message.Message{
			MessageType: message.Notification_Message,
			Method:      message.Notification_Method,
		}, notiMsg, conn)
		if err != nil {
			s.Logger.Errorf("Failed to send notification message: %v", err)
			return err
		}
		return nil
	}
}
