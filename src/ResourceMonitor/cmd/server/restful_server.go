package main

import (
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"log"
	"net"
	"net/http"
)

func runRESTServer(processorsServer pb.ProcessorsServiceServer, memoryServer pb.MemoryServiceServer, enableTLS bool, listener net.Listener, grpcEndpoint string) error {
	mux := runtime.NewServeMux()
	//dialOptions := []grpc.DialOption{grpc.WithInsecure()}

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// in-process handler
	// err := pb.RegisterAuthServiceHandlerServer(ctx, mux, authServer)

	log.Printf("Start REST server at %s, TLS = %t", listener.Addr().String(), enableTLS)
	//if enableTLS {
	//	return http.ServeTLS(listener, mux, serverCertFile, serverKeyFile)
	//}
	return http.Serve(listener, mux)
}
