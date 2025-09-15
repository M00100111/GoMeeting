package server

import (
	"GoMeeting/rpcs/ws/internal/message"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

// 统一管理所有消费者
func Consumers(s *WsServer) []service.Service {
	return []service.Service{
		NewWsConsumer(s),
	}
}

// Ws消费者
type WsConsumer struct {
	WsServer *WsServer // 保存 WebSocket 服务器引用
}

func NewWsConsumer(s *WsServer) *WsConsumer {
	return &WsConsumer{
		WsServer: s,
	}
}

func (c *WsConsumer) Consume(data []byte) error {
	// 解析消息
	msg, err := message.ParseMessage(data)
	if err != nil || msg == nil {
		c.WsServer.Logger.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	// 根据消息类型进行处理
	switch msg.MessageType {
	case message.Ping_Message:
		return c.WsServer.routes[msg.Method](c.WsServer, msg)
	case message.Chat_Message:
		return c.WsServer.routes[msg.Method](c.WsServer, msg)
	case message.Notification_Message:
		return c.WsServer.routes[msg.Method](c.WsServer, msg)
	default:
		c.WsServer.Logger.Infof("Unknown message type: %v", msg.MessageType)
		return errors.New("Unknown message type")
	}
}

// Start 启动消费者
func (c *WsConsumer) Start() {
	logx.Info("WebSocket consumer started")
	// 消费逻辑通常由外部框架调用 Consume 方法处理
}

// Stop 停止消费者
func (c *WsConsumer) Stop() {
	logx.Info("WebSocket consumer stopped")
}
