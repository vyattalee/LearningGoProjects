proto:
	protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb

clean:
	rm pb/*

memory.proto:
	protoc --proto_path=proto proto/memory_message.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb


.PHONY: proto test clean