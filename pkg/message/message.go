package message

type MessageType uint8

// 推送类型
const (
	Ping_Message MessageType = iota
	Notice_Message
	Chat_Message
)
