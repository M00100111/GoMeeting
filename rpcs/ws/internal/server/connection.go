package server

import (
	"GoMeeting/pkg/structs/message"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

const MaxConnectionIdle = time.Minute * 10

// 对websocket连接对象的封装
type WsConn struct {
	//全局WsServer指针,便于调用全局WsServer的方法
	*WsServer
	Uid             string // 建立连接的用户id
	*websocket.Conn        //websocket连接

	Done chan struct{} //用于通知关闭Ws连接的通道

	lastActiveTimeMutex sync.Mutex
	lastActiveTime      time.Time //最后一次活跃时间

	//用于通信的管道
	Message chan *message.Message
}

// 将当前客户端的http连接升级为websocket连接
func NewWsConn(s *WsServer, w http.ResponseWriter, r *http.Request, uid string) *WsConn {
	//调用Ws的websocket的Upgrade方法以升级http连接为websocket连接
	c, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("websocket升级失败: %v", err)
		return nil
	}
	//新建WsConn连接对象
	conn := &WsConn{
		WsServer:       s,
		Conn:           c,
		Uid:            uid,
		Done:           make(chan struct{}),
		Message:        make(chan *message.Message, 1),
		lastActiveTime: time.Now(),
	}
	//闲置连接回收
	go conn.keepalive()
	return conn
}

// 读取连接对象发来的二进制消息
func (c *WsConn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()

	c.lastActiveTimeMutex.Lock()
	defer c.lastActiveTimeMutex.Unlock()

	c.lastActiveTime = time.Time{}
	return
}

// 向连接对象推送二进制消息
func (c *WsConn) WriteMessage(messageType int, data []byte) error {
	c.lastActiveTimeMutex.Lock()
	defer c.lastActiveTimeMutex.Unlock()

	err := c.Conn.WriteMessage(messageType, data)
	c.lastActiveTime = time.Now()
	return err
}

// UpdateLastActiveTime 更新最后活跃时间
func (c *WsConn) UpdateLastActiveTime() {
	c.lastActiveTimeMutex.Lock()
	defer c.lastActiveTimeMutex.Unlock()
	c.lastActiveTime = time.Now()
}

// 自动释放空闲的连接, 如果连接闲置时间过长,则关闭连接
func (c *WsConn) keepalive() {
	idleTimer := time.NewTimer(MaxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()
	//循环
	for {
		select {
		case <-idleTimer.C:
			//计时器到期,判断连接是否闲置
			c.lastActiveTimeMutex.Lock()
			lastTime := c.lastActiveTime
			idleDuration := time.Since(lastTime)
			if idleDuration > MaxConnectionIdle {
				//连接闲置时间过长,关闭连接
				c.lastActiveTimeMutex.Unlock()
				c.CloseWsConn()
				return
			} else {
				remainingTime := MaxConnectionIdle - idleDuration
				if remainingTime <= 0 {
					remainingTime = MaxConnectionIdle
				}
				idleTimer.Reset(remainingTime)
			}
			c.lastActiveTimeMutex.Unlock()
		case <-c.Done:
			//连接已关闭,停止
			return
		}
	}
}

// 关闭Ws连接
func (c *WsConn) CloseWsConn() error {
	fmt.Println("关闭Ws连接")
	//停止与客户端的心跳
	//避免关闭已关闭的通道
	select {
	case <-c.Done: // 已关闭会立即返回
	default:
		close(c.Done) // 关闭通道
	}
	return c.Conn.Close() // 关闭websocket连接

	//在会人员自动退出会议
}
