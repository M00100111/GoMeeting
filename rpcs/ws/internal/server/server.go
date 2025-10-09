package server

import (
	"GoMeeting/pkg/ctxdata"
	"GoMeeting/rpcs/ws/internal/message"
	"GoMeeting/rpcs/ws/internal/svc"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
	"strconv"
	"sync"
)

const (
	defaultConcurrency = 10
)

// 定义全局WS服务器结构体
type WsServer struct {
	svc *svc.ServiceContext //调用其他服务

	Upgrader websocket.Upgrader // 用于升级为websocket服务
	// 用户id-连接映射,实现统一管理
	// 为应对高并发,修改时使用读写锁
	sync.RWMutex
	connToUser map[*WsConn]string
	userToConn map[string]*WsConn // 用于获取要转发对象的连接

	routes map[message.MessageMethod]WsHandlerFunc // 路由与WS处理函数映射

	*threading.TaskRunner // 异步任务管理
	logx.Logger           // 日志记录
}

func (s *WsServer) GetRedisClient() *redis.Redis {
	return s.svc.Redis
}

// 待实现选项模式
func NewWsServer(ctx *svc.ServiceContext) *WsServer {
	wsServer := &WsServer{
		svc:        ctx,
		connToUser: make(map[*WsConn]string),
		userToConn: make(map[string]*WsConn),
		routes:     make(map[message.MessageMethod]WsHandlerFunc),
		TaskRunner: threading.NewTaskRunner(defaultConcurrency),
		Logger:     logx.WithContext(nil),
	}
	// 简单配置，允许所有来源（仅用于开发环境）
	wsServer.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源
		},
	}
	return wsServer
}

// WsServer服务启动
func (s *WsServer) Start() {
	http.HandleFunc(s.svc.Config.Pattern, s.BuildWsConn)
	fmt.Println("websocket server start")
	s.Info(http.ListenAndServe(s.svc.Config.ListenOn, nil))
}

// WsServer服务关闭
func (s *WsServer) Stop() {

}

// 服务启动时注册方法与处理函数到全局路由
func (s *WsServer) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

// net/http的全局处理函数, 处理所有http连接以建立ws连接
func (s *WsServer) BuildWsConn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("检测到建立ws连接请求")
	//鉴权并获取用户名,只有已登录的用户才能建立ws连接

	//// 从Cookie获取token
	//cookie, err := r.Cookie("token")
	//if err != nil {
	//	http.Error(w, "缺少token", http.StatusUnauthorized)
	//	fmt.Println("缺失token")
	//	return
	//}
	//authHeader := cookie.Value

	// 修改为从查询参数获取：
	authHeader := r.URL.Query().Get("token")

	fmt.Println("token: %v", authHeader)
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
	conn := NewWsConn(s, w, r, strconv.FormatUint(uid, 10))
	if conn == nil {
		return
	}
	fmt.Println("WsConn:", conn)
	// 注册ws连接到全局连接管理映射
	s.addWsConn(conn)

	// 启用协程以接收当前连接发来的信息
	go s.readMessage(conn)
}

// 注册ws连接到全局连接管理映射
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
func (s *WsServer) GetWsConnByUid(uid string) *WsConn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.userToConn[uid]
}

func (s *WsServer) GetUidByWsConn(conn *WsConn) string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.connToUser[conn]
}

func (s *WsServer) GetWsConnsByUids(uids []string) []*WsConn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	var conns []*WsConn
	for _, uid := range uids {
		conns = append(conns, s.userToConn[uid])
	}
	return conns
}

func (s *WsServer) GetUidsByWsConns(conns []*WsConn) []string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	var uids []string
	for _, conn := range conns {
		uids = append(uids, s.connToUser[conn])
	}
	return uids
}

// 注销WsConn连接
func (s *WsServer) CloseWsServer(conn *WsConn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	//注销连接管理
	delete(s.connToUser, conn)
	delete(s.userToConn, conn.Uid)
	//关闭websocket连接
	conn.CloseWsConn()
}

// Conn负责接收二进制消息
// Message负责解析消息到结构体
// readMessage将消息转发到对应的处理函数
// 启动协程处理请求以处理同一连接对象的多个请求
func (s *WsServer) readMessage(conn *WsConn) {
	//处理消息,通过协程和chan等待处理好的消息，根据消息调用对应的处理函数
	go s.handleMessage(conn)

	//接收并解析和初步处理消息
	for {
		//读取消息以获取二进制消息
		_, data, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err %v", err)
			// 关闭连接对象
			s.CloseWsServer(conn)
			return
		}
		fmt.Println("接收到数据")

		//发布任务
		//if err = s.svc.KafkaWsPusher.Push(context.Background(), string(data)); err != nil {
		//	s.Logger.Errorf("kafka push err %v", err)
		//	s.CloseWsServer(conn)
		//	return
		//}

		// 直接处理
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

// 消息处理,由直接调用处理函数改为发布任务
func (s *WsServer) handleMessage(conn *WsConn) {
	for {
		select {
		//连接已关闭
		case <-conn.Done:
			return
		case msg := <-conn.Message:
			fmt.Println("处理消息")
			//根据消息类型进行处理
			switch msg.MessageType {
			case message.Ping_Message:
				fmt.Println("处理ping消息")
				handler, _ := s.routes[msg.Method]
				handler(s, msg)
			case message.Chat_Message:
				//根据请求的方法执行处理函数
				if handler, ok := s.routes[msg.Method]; ok {
					handler(s, msg)
				} else { // 方法不存在
					jsonData, err := json.Marshal(message.ErrorData{
						Msg: "方法不存在",
					})
					if err != nil {
						s.Errorf("json marshal err %v", err)
						return
					}
					s.SendMessage(&message.Message{
						MessageType: message.Chat_Message,
						Method:      message.Pong_Method,
						Data:        jsonData,
					}, conn)
				}
			case message.WebRTC_Message:
				fmt.Printf("接收到WebRTC请求,Type: %d,Method: %s", msg.MessageType, msg.Method)
				var result map[string]interface{}
				if err := json.Unmarshal(msg.Data, &result); err != nil {
					s.Logger.Errorf("Failed to unmarshal message.Data to map[string]interface{}: %v", err)
					return
				}
				result["senderId"] = s.GetUidByWsConn(conn)
				jsonResult, err := json.Marshal(result)
				if err != nil {
					s.Logger.Errorf("Failed to marshal result to JSON: %v", err)
					return
				}
				msg.Data = jsonResult
				var receiverConn *WsConn
				if result["receiverId"] == nil {
					s.Logger.Errorf("receiverId is nil")
					return
				}
				receiverConn = s.GetWsConnByUid(result["receiverId"].(string))
				if receiverConn != nil {
					s.SendMessage(msg, receiverConn)
				} else {
					s.Logger.Errorf("找不到接收者连接: %s", result["receiverId"])
				}
			}
		}
	}
}

// 消息推送
// Conn负责接收二进制消息
// Message负责解析消息到结构体
func (s *WsServer) SendMessage(msg *message.Message, data interface{}, conns ...*WsConn) error {
	result, err := message.BuildMessage(msg, data)
	if err != nil {
		s.Errorf("json marshal err %v", err)
		return err
	}
	//群发消息可单发
	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, result); err != nil {
			return err
		}
	}
	return nil
}
