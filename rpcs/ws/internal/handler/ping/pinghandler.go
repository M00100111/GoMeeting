package ping

import (
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/rpcs/ws/internal/message"
	"GoMeeting/rpcs/ws/internal/server"
	"encoding/json"
	"fmt"
)

func PingHandler() server.WsHandlerFunc {
	return func(s *server.WsServer, msg *message.Message) error {
		fmt.Println("start ping handler")
		var pingMsg message.PingData
		if err := json.Unmarshal(msg.Data, &pingMsg); err != nil {
			s.Logger.Errorf("Failed to unmarshal message.Data to message.PingMessage: %v", err)
			return err
		}
		fmt.Println("pingMsg:", pingMsg.Msg)
		conn := s.GetWsConnByUid(pingMsg.SenderId)
		result := &message.PingData{
			SenderId: ctxdata.Root,
			Msg:      pingMsg.Msg + "!!!",
		}

		err := s.SendMessage(&message.Message{
			MessageType: message.Pong_Message,
			Method:      message.Pong_Method,
		}, result, conn)
		if err != nil {
			s.Logger.Errorf("Failed to send ping message: %v", err)
			return err
		}
		return nil
	}
}
