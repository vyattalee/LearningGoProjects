package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/LearningGoProjects/ResourceMonitor/service"
	"github.com/LearningGoProjects/ResourceMonitor/utils"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

func loadTLSCredentials(onlyserversidetls bool) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile(utils.ClientCACertFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(utils.ServerCertFile, utils.ServerKeyFile)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		//ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientAuth: tls.NoClientCert,
		ClientCAs:  certPool,
	}

	return credentials.NewTLS(config), nil
}

func runGRPCServer(processorsServer pb.ProcessorsServiceServer, memoryServer pb.MemoryServiceServer, enableTLS bool, listener net.Listener) error {

	jwtManager := service.NewJWTManager(utils.SecretKey, utils.TokenDuration)

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
		grpc.ConnectionTimeout(30 * time.Second),
		grpc.MaxConcurrentStreams(10),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute, //这个连接最大的空闲时间，超过就释放，解决proxy等到网络问题（不通知grpc的client和server）
		}),
	}

	if enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			return fmt.Errorf("cannot load TLS credentials: %w", err)
		}

		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	grpcServer := grpc.NewServer(serverOptions...)

	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: ", err)
	}

	authServer := service.NewAuthServer(userStore, jwtManager)
	resourceMonitorServer := service.NewResourceMonitorServer()

	// Register the server
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	pb.RegisterResourceMonitorServiceServer(grpcServer, resourceMonitorServer)
	pb.RegisterProcessorsServiceServer(grpcServer, processorsServer)
	//pb.RegisterMemoryServiceServer(grpcServer, memoryServer)
	reflection.Register(grpcServer)

	// Start sending data to subscribers
	go resourceMonitorServer.StartService()

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

func accessibleRoles() map[string][]string {
	const resourceMonitorServicePath = "/LearningGoProjects.ResourceMonitor.ResourceMonitorService/"

	return map[string][]string{
		resourceMonitorServicePath + "Subscribe":   {"admin", "user"},
		resourceMonitorServicePath + "Unsubscribe": {"admin", "user"},
	}
}
