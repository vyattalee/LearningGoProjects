proto:
	protoc --go_out=./  --go-grpc_out=./ --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative protos/longlived.proto
#

clean:
	rm protos/longlived.pb.go protos/longlived_grpc.pb.go

server:
	go run server/server.go
# 	./server

client:
	go run client/client.go
# 	./client


.PHONY: proto test clean server client