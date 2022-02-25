prepare:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

proto:
	protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb

clean:
	rm pb/*

server:
	go run cmd/server/main.go cmd/server/grpc_server.go cmd/server/restful_server.go -port 8080

server-tls:
	go run cmd/server/main.go cmd/server/grpc_server.go cmd/server/restful_server.go -port 8080 -tls

client:
	go run cmd/client/main.go  -address 0.0.0.0:8080

client-tls:
	go run cmd/client/main.go  -address 0.0.0.0:8080 -tls

memory.proto:
	protoc --proto_path=proto proto/memory_message.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb

test:
	go test -v -cover -covermode=atomic -coverprofile=cover.out ./...

cert:
	cd cert; ./gen.sh; cd ..


.PHONY: proto test clean server client cert