# GoMeeting 视频会议系统

GoMeeting 是一个基于 Go 语言开发的实时视频会议系统，采用微服务架构，使用 go-zero 框架，集成了 gRPC 和 WebSocket 实现服务间通信与实时消息传输。

前后端交互效果呈现见前端项目的 README 文档：[GoMeetingClient](https://github.com/M00100111/GoMeetingClient)

## 技术栈

- 后端语言：Go 1.23.1
- 微服务框架：go-zero
- 数据库：MySQL、Redis
- 实时通信：WebSocket、WebRTC
- 消息队列：Kafka（计划中）
- 服务发现：etcd
- 容器化：Docker
- API 文档：自动生成的 Swagger 文档

## 系统架构

系统采用微服务架构，主要包括以下几个服务：

1. **API 网关服务** - 统一入口，处理 HTTP 请求
2. **用户服务 (user)** - 用户注册、登录、信息管理
3. **会议服务 (meeting)** - 会议创建、管理、控制
4. **社交服务 (social)** - 好友、群组功能
5. **WebSocket 服务 (ws)** - 实时通信处理

## 功能特性

### 用户管理
- 邮箱验证码注册
- 用户登录认证
- JWT Token 鉴权

### 会议功能
- 创建会议
- 加入/离开会议
- 开始/结束会议
- 会议室管理（公开/私密）

### 社交功能
- 好友管理
- 群组聊天
- 好友/群组申请处理

### 实时通信
- WebSocket 连接管理
- WebRTC 信令传输
- 实时消息推送

## 项目结构

```
GoMeeting/
├── api/              # API网关服务
├── rpcs/             # RPC服务
│   ├── user/         # 用户服务
│   ├── meeting/      # 会议服务
│   ├── social/       # 社交服务
│   └── ws/           # WebSocket服务
├── pkg/              # 公共包
├── deploy/           # 部署配置
└── sql/              # 数据库脚本
```

## 数据库设计

系统使用三个主要数据库：

1. **user** - 用户信息表
2. **meeting** - 会议信息和成员表
3. **social** - 社交关系表（好友、群组等）

详细表结构请查看 [sql](./sql) 目录下的 SQL 文件。

## 快速开始

### 环境要求

- Go 1.23.1
- Docker & Docker Compose
- MySQL 5.7
- Redis
- etcd
- Kafka (可选)

### 安装步骤

1. 克隆项目代码：
```bash
git clone <项目地址>
cd GoMeeting
```

2. 启动依赖服务：
```bash
cd deploy
docker-compose up -d
```

3. 初始化数据库：
```bash
# 执行 sql 目录下的 SQL 文件初始化数据库
```

4. 启动各服务：
```bash
# 分别启动各个微服务
```

### 配置说明

各服务的配置文件位于对应服务的 `etc/` 目录下，请根据实际环境修改相关配置。

## API 文档

系统使用 go-zero 自动生成 API 文档，具体接口请参考各服务的 handler 文件。

主要接口包括：

### 用户相关
- `POST /user/signup` - 用户注册
- `POST /user/login` - 用户登录
- `GET /user/pinguser` - 用户服务健康检查

### 工具相关
- `POST /tool/captcha` - 获取验证码

### 会议相关
- `POST /meeting/startmeeting` - 开始会议
- `POST /meeting/joinmeeting` - 加入会议
- `POST /meeting/leavemeeting` - 离开会议
- `POST /meeting/endmeeting` - 结束会议
- `GET /meeting/getmeetinginfo` - 获取会议信息

## 部署说明

项目支持 Docker 容器化部署，使用 docker-compose 管理所有依赖服务。

构建和部署步骤：

```bash
# 构建各服务镜像
# 编写 Dockerfile...

# 使用 docker-compose 启动所有服务
cd deploy
docker-compose up -d
```

## 开发指南

### 代码生成

本项目使用 go-zero 提供的代码生成工具，可以快速生成 API 和 RPC 代码：

```bash
# 生成 API 代码
goctl api go -api *.api -dir .

# 生成 RPC 代码
goctl rpc protoc *.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

### 代码规范

- 遵循 Go 语言官方编码规范
- 使用 go-zero 框架推荐的最佳实践
- 统一错误处理和响应格式

## 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进项目。

## 许可证

[MIT License](LICENSE)