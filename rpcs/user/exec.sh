#存放常用执行命令

//构建或更新UserRPC服务
goctl rpc protoc ./rpcs/user/rpc/user.proto --go_out=./rpcs/user/rpc/ --go-grpc_out=./rpcs/user/rpc/ --zrpc_out=./rpcs/user/rpc/

//构建或更新usersModel
goctl model mysql ddl --src="./sql/user/user.sql" --dir "./rpcs/user/models" -c

//构建或更新UserAPI服务
goctl api go -api ./apps/user/api/user.api -dir ./apps/user/api -style gozero