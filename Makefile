protobuf:
	protoc --go_out=plugins=grpc,import_path=pkg/pb:. grpc-demo-server/pkg/pb/demo.proto