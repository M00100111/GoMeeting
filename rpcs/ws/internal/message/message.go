package message

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

type MessageType uint8

// 聊天消息
const (
	Ping_Message MessageType = iota
	Pong_Message
	Notification_Message
	Chat_Message
)

type Message struct {
	MessageType MessageType     `json:"message_type"`
	Method      string          `json:"method"` //请求方法
	Data        json.RawMessage `json:"data"`   // 数据
}

type PingMessage struct {
	SenderId string `json:"sender_id"` //发送者id
	Msg      string `json:"msg"`
}

type ChatMessage struct {
	SenderId   string   `json:"sender_id"` //发送者id
	ChatType   ChatType `json:"chat_type"`
	ReceiverId string   `json:"receiver_id"`
	MsgType    MsgType  `json:"msg_type"`
	Msg        string   `json:"msg"`
	SendTime   int64    `json:"send_time"`
}

type ChatType uint8

const (
	SingleChat ChatType = iota
	GroupChat
)

type MsgType uint8

const (
	Text_Msg MsgType = iota
	Image_Msg
	Video_Msg
	File_Msg
)

type NotificationMessage struct {
	ReceiverId string `json:"receiver_id"`
	Msg        string `json:"msg"`
}

type MessageMethod string

const (
	Ping_Method                        MessageMethod = "Ping"
	Pong_Method                                      = "Pong"
	Chat_Method                                      = "Chat"
	Notification_Method                              = "Notification"
	Meeting_Start_Notice_Method                      = "Meeting_Start_Notification"
	Meeting_End_Notice_Method                        = "Meeting_End_Notification"
	Meeting_Member_Join_Notice_Method                = "Meeting_Member_Join"
	Meeting_Member_Leave_Notice_Method               = "Meeting_Member_Leave"
	Meeting_Message_Notice_Method                    = "Meeting_Message_Leave"
	Group_Member_Join_Notice_Method                  = "Meeting_Member_Join"
	Group_Member_Leave_Notice_Method                 = "Meeting_Member_Leave"
	Group_Message_Notice_Method                      = "Meeting_Member_Notice"
)

// 解析conn接收到的二进制消息
func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// 将消息转为二进制消息
func BuildMessage(msg *Message, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logx.Errorf("Failed to marshal Data to json: %v", err)
		return nil, err
	}
	msg.Data = jsonData
	result, err := json.Marshal(msg)
	if err != nil {
		logx.Errorf("Failed to marshal message to json: %v", err)
		return nil, err
	}
	return result, nil
}
