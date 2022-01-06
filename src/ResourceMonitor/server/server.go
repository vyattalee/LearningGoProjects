package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"github.com/LearningGoProjects/ResourceMonitor/service"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	serverType := flag.String("type", "grpc", "type of server (grpc/rest)")
	endPoint := flag.String("endpoint", "", "gRPC endpoint")
	flag.Parse()


	processors_server = server.

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	if *serverType == "grpc" {
		err = runGRPCServer(laptopServer, *enableTLS, listener)
	} else {
		err = runRESTServer(laptopServer, *enableTLS, listener, *endPoint)
	}

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}