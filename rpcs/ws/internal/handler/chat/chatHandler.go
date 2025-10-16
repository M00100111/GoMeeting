package chat

import (
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/pkg/structs/message"
	"GoMeeting/rpcs/ws/internal/server"
	"context"
	"encoding/json"
	"fmt"
)

func ChatHandler() server.WsHandlerFunc {
	return func(s *server.WsServer, msg *message.Message) error {
		var chatMsg message.ChatData
		if err := json.Unmarshal(msg.Data, &chatMsg); err != nil {
			s.Logger.Errorf("Failed to unmarshal message.Data to message.ChatMessage: %v", err)
			return err
		}

		var conns []*server.WsConn
		switch chatMsg.ChatType {
		case message.SingleChat:
			conns = append(conns, s.GetWsConnByUid(chatMsg.ReceiverId))
		case message.GroupChat:
			// 从Redis中获取群组成员集合
			// Redis中所有数据都以字符串形式存储，即使存入整数，也会被转换为字符串
			members, err := s.GetRedisClient().SmembersCtx(context.Background(), fmt.Sprintf(ctxdata.GroupMemberPrefix, chatMsg.ReceiverId))
			if err != nil {
				s.Logger.Errorf("Failed to get group membersId from redis: %v", err)
				return err
			}

			// 根据成员ID获取对应的WebSocket连接
			for _, memberId := range members {
				if conn := s.GetWsConnByUid(memberId); conn != nil {
					conns = append(conns, conn)
				}
			}
		}

		err := s.SendMessage(msg, chatMsg, conns...)
		if err != nil {
			s.Logger.Errorf("Failed to send chat message: %v", err)
			return err
		}
		return nil
	}
}
