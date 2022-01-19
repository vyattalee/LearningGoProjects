proto:
	protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb

clean:
	rm pb/*

server:
	go run cmd/server/server.go -port 8080

client:
	go run cmd/client/client.go  -address 0.0.0.0:8080

memory.proto:
	protoc --proto_path=proto proto/memory_message.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb

test:
	go test -v -cover -covermode=atomic -coverprofile=cover.out ./...

.PHONY: proto test clean server client