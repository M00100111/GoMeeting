package main

import (
	"GoMeeting/ws/internal/handler"
	"GoMeeting/ws/internal/server"
	"flag"
	"fmt"

	"GoMeeting/ws/internal/config"
	"GoMeeting/ws/internal/svc"

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
	handler.RegisterHandlers(s, ctx)

	fmt.Printf("Starting ws server at %s...\n", c.ListenOn)
	s.Start()
}
