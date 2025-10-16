package server

import (
	"GoMeeting/pkg/structs/message"
	"context"
	"errors"
	"fmt"
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
	fmt.Println("Consuming message")
	logx.Info("Consuming message")
	// 解析消息
	msg, err := message.ParseMessage(data)
	if err != nil || msg == nil {
		c.WsServer.Logger.Errorf("Failed to unmarshal message: %v", err)
		return err
	}
	fmt.Println("msg: %v", msg)
	// 根据消息类型进行处理
	switch msg.MessageType {
	case message.Ping_Message:
		return c.WsServer.routes[msg.Method](c.WsServer, msg)
	case message.Chat_Message:
		return c.WsServer.routes[msg.Method](c.WsServer, msg)
	case message.Notification_Message:
		return c.WsServer.routes[msg.Method](c.WsServer, msg)
	case message.WebRTC_Message:
		return c.WsServer.routes[msg.Method](c.WsServer, msg)
	default:
		c.WsServer.Logger.Infof("Unknown message type: %v", msg.MessageType)
		return errors.New("Unknown message type")
	}
}

// Start 启动消费者
func (c *WsConsumer) Start() {
	fmt.Println("Starting WebSocket consumer...")
	logx.Info("WebSocket consumer started")
	// 消费逻辑通常由外部框架调用 Consume 方法处理
	// 启动独立的 goroutine 来持续消费消息
	for i := 0; i < c.WsServer.svc.Config.KafkaConsumerConf.Consumers; i++ {
		go c.runConsumer(i)
	}
}

func (c *WsConsumer) runConsumer(id int) {
	consumer := c.WsServer.svc.CreateKafkaConsumer()
	defer consumer.Close()

	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			if err == context.Canceled {
				break
			}
			continue
		}

		if err := c.Consume(msg.Value); err != nil {
			logx.Errorf("Failed to consume message: %v", err)
		}

	}
}

// Stop 停止消费者
func (c *WsConsumer) Stop() {
	logx.Info("WebSocket consumer stopped")
}
