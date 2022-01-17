package main

import (
	"flag"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/LearningGoProjects/ResourceMonitor/service"
	"github.com/LearningGoProjects/ResourceMonitor/utils"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

func runGRPCServer(processorsServer pb.ProcessorsServiceServer, memoryServer pb.MemoryServiceServer, enableTLS bool, listener net.Listener) error {

	serverOptions := []grpc.ServerOption{
		//grpc.UnaryInterceptor(interceptor.Unary()),
		//grpc.StreamInterceptor(interceptor.Stream()),
		grpc.ConnectionTimeout(30 * time.Second),
		grpc.MaxConcurrentStreams(10),
		//grpc.KeepaliveParams(keepalive.ServerParameters{
		//	MaxConnectionIdle: 5 * time.Minute, //这个连接最大的空闲时间，超过就释放，解决proxy等到网络问题（不通知grpc的client和server）
		//}),
	}

	//if enableTLS {
	//	tlsCredentials, err := loadTLSCredentials()
	//	if err != nil {
	//		return fmt.Errorf("cannot load TLS credentials: %w", err)
	//	}
	//
	//	serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	//}

	grpcServer := grpc.NewServer(serverOptions...)

	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: ", err)
	}

	jwtManager := service.NewJWTManager(utils.SecretKey, utils.TokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	resourceMonitorServer := service.NewResourceMonitorServer()

	// Register the server
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterResourceMonitorServiceServer(grpcServer, resourceMonitorServer)

	// Start sending data to subscribers
	go resourceMonitorServer.StartService()

	pb.RegisterProcessorsServiceServer(grpcServer, processorsServer)
	//pb.RegisterMemoryServiceServer(grpcServer, memoryServer)

	log.Printf("Start GRPC server at %s, TLS = %t", listener.Addr().String(), enableTLS)
	return grpcServer.Serve(listener)
}

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

func main() {
	port := flag.Int("port", 0, "the server port")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	serverType := flag.String("type", "grpc", "type of server (grpc/rest)")
	endPoint := flag.String("endpoint", "", "gRPC endpoint")
	flag.Parse()

	//each server's service must be register to grpcserver
	processors_server := service.NewProcessorsServer()
	memory_server := service.NewMemoryServer()

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	if *serverType == "grpc" {
		err = runGRPCServer(processors_server, memory_server, *enableTLS, listener)
	} else {
		err = runRESTServer(processors_server, memory_server, *enableTLS, listener, *endPoint)
	}

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

}

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}
