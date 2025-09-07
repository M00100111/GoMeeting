package server

import (
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/rpcs/ws/internal/message"
	"GoMeeting/rpcs/ws/internal/svc"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

// 定义全局WS服务器结构体
type WsServer struct {
	svc *svc.ServiceContext

	Upgrader websocket.Upgrader // 用于升级为websocket服务
	// 用户id-连接映射,实现统一管理
	// 为应对高并发,修改时使用读写锁
	sync.RWMutex
	connToUser map[*WsConn]string
	userToConn map[string]*WsConn // 用于获取要转发对象的连接

	routes map[string]WsHandlerFunc // 路由与WS处理函数映射

	logx.Logger // 日志记录
}

func NewWsServer(ctx *svc.ServiceContext) *WsServer {
	return &WsServer{
		svc:        ctx,
		connToUser: make(map[*WsConn]string),
		userToConn: make(map[string]*WsConn),
		routes:     make(map[string]WsHandlerFunc),
		Logger:     logx.WithContext(nil),
	}
}

// 添加路由
func (s *WsServer) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

// net/http的全局处理函数, 处理所有http连接以建立ws连接
func (s *WsServer) BuildWsConn(w http.ResponseWriter, r *http.Request) {
	//鉴权并获取用户名,只有已登录的用户才能建立ws连接
	// 从Authorization获取token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "缺少token", http.StatusUnauthorized)
		return
	}
	// 验证token
	if !ctxdata.ValidateToken(authHeader, s.svc.Config.Jwt.AccessSecret) {
		http.Error(w, "token无效", http.StatusUnauthorized)
		return
	}
	uid, err := ctxdata.GetUidFromToken(authHeader, s.svc.Config.Jwt.AccessSecret)
	if err != nil {
		http.Error(w, "无法解析token", http.StatusUnauthorized)
		return
	}

	//将当前客户端的http连接升级为websocket连接
	conn := NewWsConn(s, w, r, uid)
	if conn == nil {
		return
	}
	//连接统一管理
	s.addWsConn(conn)

	// 启用协程以同时处理多个连接的请求
	go s.handlerWsConn(conn)
}

// 用于建立ws连接并管理ws连接
func (s *WsServer) addWsConn(conn *WsConn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	//重复登录时关闭之前未断开的连接
	if c := s.userToConn[conn.Uid]; c != nil {
		c.CloseWsConn()
	}
	s.userToConn[conn.Uid] = conn
	s.connToUser[conn] = conn.Uid
}

// 全局连接映射管理的查操作
func (s *WsServer) getWsConnByUid(uid string) *WsConn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.userToConn[uid]
}

func (s *WsServer) getUidByWsConn(conn *WsConn) string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.connToUser[conn]
}

func (s *WsServer) getWsConnsByUids(uids []string) []*WsConn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	var conns []*WsConn
	for _, uid := range uids {
		conns = append(conns, s.userToConn[uid])
	}
	return conns
}

func (s *WsServer) getUidsByWsConns(conns []*WsConn) []string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	var uids []string
	for _, conn := range conns {
		uids = append(uids, s.connToUser[conn])
	}
	return uids
}

func (s *WsServer) CloseWsServer(conn *WsConn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	//注销连接管理
	delete(s.connToUser, conn)
	delete(s.userToConn, conn.Uid)
	//关闭websocket连接
	conn.CloseWsConn()
}

func (s *WsServer) Start() {
	http.HandleFunc(s.svc.Config.Pattern, s.BuildWsConn)
	fmt.Println("websocket server start")
	s.Info(http.ListenAndServe(s.svc.Config.ListenOn, nil))
}

func (s *WsServer) Stop() {

}

// 启动协程为当前连接对象处理任务请求
func (s *WsServer) handlerWsConn(conn *WsConn) {
	//处理消息,通过协程和chan等待处理好的消息，根据消息调用对应的处理函数
	go s.handlerMessage(conn)

	//接收并解析和初步处理消息
	for {
		//获取二进制消息
		_, data, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			// 关闭连接对象
			s.CloseWsServer(conn)
			return
		}

		//解析消息
		msg, err := message.ParseMessage(data)
		if err != nil {
			s.Errorf("json unmarshal err %v, msg %v", err, string(data))
			// 关闭连接对象
			s.CloseWsServer(conn)
			return
		}

		//处理消息
		conn.Message <- msg
	}
}

// 消息处理
func (s *WsServer) handlerMessage(conn *WsConn) {
	for {
		select {
		//连接已关闭
		case <-conn.Done:
			return
		case msg := <-conn.Message:
			//根据消息类型进行处理
			switch msg.MessageType {
			case message.Data_Message:
				//根据请求的方法执行处理函数
				if handler, ok := s.routes[msg.Method]; ok {
					handler(s, conn, msg)
				} else {
					s.SendMessage(&message.Message{MessageType: message.Data_Message, Data: fmt.Sprintf("method %v not found", msg.Method)}, conn)
				}
			}
		}
	}
}

func (s *WsServer) SendMessage(msg *message.Message, conns ...*WsConn) error {
	data, err := message.BuildMessage(msg)
	if err != nil {
		s.Errorf("json marshal err %v", err)
		return err
	}
	//群发消息可单发
	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}
	return nil
}
