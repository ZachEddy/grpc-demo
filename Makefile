protobuf:
	protoc --go_out=plugins=grpc,import_path=pkg/pb:. demo-server/pkg/pb/demo.proto
