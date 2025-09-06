//构建或更新MeetingRPC服务
goctl rpc protoc ./rpcs/meeting/rpc/meeting.proto --go_out=./rpcs/meeting/rpc/ --go-grpc_out=./rpcs/meeting/rpc/ --zrpc_out=./rpcs/meeting/rpc/

//构建或更新meetingsModel
goctl model mysql ddl --src="./sql/meeting.sql" --dir "./rpcs/meeting/models" -c


//构建或更新MeetingRPC服务
goctl rpc protoc ./rpcs/meeting/rpc/meeting.proto --go_out=./rpcs/meeting/rpc/ --go-grpc_out=./rpcs/meeting/rpc/ --zrpc_out=./rpcs/meeting/rpc/ --proto_path=$(go env GOMODCACHE)/google.golang.org/protobuf@v1.36.8


