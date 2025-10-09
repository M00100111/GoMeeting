package main

import (
	"GoMeeting/rpcs/ws/internal/handler"
	"GoMeeting/rpcs/ws/internal/server"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"

	"GoMeeting/rpcs/ws/internal/config"
	"GoMeeting/rpcs/ws/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/ws.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := server.NewWsServer(ctx)
	defer s.Stop()

	// 服务启动时注册路由与处理函数
	handler.RegisterHandlers(s)

	//统一管理消费者
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	for _, consumer := range server.Consumers(s) {
		serviceGroup.Add(consumer)
	}
	// 在单独的goroutine中启动服务组，避免阻塞
	go serviceGroup.Start()

	fmt.Printf("Starting ws server at %s...\n", c.ListenOn)
	s.Start()
}
