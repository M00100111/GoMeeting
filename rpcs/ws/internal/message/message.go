package message

import "encoding/json"

type MessageType uint8

// 聊天消息
const (
	Err_Message MessageType = iota

	Data_Message
)

type Message struct {
	SenderId    string `json:"senderId"` //发送者id
	Method      string `json:"method"`   //请求方法
	MessageType `json:"messageType"`
	Data        interface{} `json:"data"` // 数据
}

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
func BuildMessage(msg *Message) ([]byte, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return data, nil
}
