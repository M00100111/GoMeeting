//构建或更新MeetingRPC服务
goctl rpc protoc ./rpcs/social/rpc/social.proto --go_out=./rpcs/social/rpc/ --go-grpc_out=./rpcs/social/rpc/ --zrpc_out=./rpcs/social/rpc/

//构建或更新socialModel
goctl model mysql ddl --src="./sql/social.sql" --dir "./rpcs/social/models" -c



