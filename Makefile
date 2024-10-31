proto:
	protoc -I. -I./api/googleapis --go_out ./pkg/ --go-grpc_out ./pkg/ --grpc-gateway_out ./pkg/ api/liveness.proto