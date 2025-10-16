package notification

import (
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/pkg/structs/message"
	"GoMeeting/rpcs/ws/internal/server"
	"context"
	"encoding/json"
	"fmt"
)

func MeetingMemberJoinNoticeHandler() server.WsHandlerFunc {
	return func(s *server.WsServer, msg *message.Message) error {
		fmt.Println("调用notification.NewMemberJoinHandler")
		var MemberJoinMsg message.MeetingMemberJoinNoticeData
		if err := json.Unmarshal(msg.Data, &MemberJoinMsg); err != nil {
			s.Logger.Errorf("Failed to unmarshal message.Data to message.MeetingMemberJoinNoticeData: %v", err)
			return err
		}
		fmt.Println("MemberJoinMsg:", MemberJoinMsg)
		fmt.Println("MemberJoinMsg.UserId:", MemberJoinMsg.UserId)
		// 从Redis中获取会议中所有成员ID
		memberSetKey := fmt.Sprintf(ctxdata.MeetingMemberPrefix, MemberJoinMsg.MeetingId)
		members, err := s.GetRedisClient().SmembersCtx(context.Background(), memberSetKey)
		if err != nil {
			s.Logger.Errorf("Failed to get members from Redis Set: %v", err)
			return err
		}
		fmt.Println("members:", members)

		var uids []string
		for _, memberId := range members {
			// 替换第33行代码为以下内容：
			if memberId == fmt.Sprintf("%d", MemberJoinMsg.UserId) {
				continue
			}
			// 假设 memberId 是字符串形式的用户ID，直接添加进 uids
			uids = append(uids, memberId)
		}
		fmt.Println("uids:", uids)
		conns := s.GetWsConnsByUids(uids)
		err = s.SendMessage(&message.Message{
			MessageType: message.WebRTC_Message,
			Method:      message.Meeting_Member_Join_Notice_Method,
		}, MemberJoinMsg, conns...)
		if err != nil {
			s.Logger.Errorf("Failed to send member join message: %v", err)
			return err
		}
		return nil
	}
}
