package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/LearningGoProjects/ResourceMonitor/registry/consul"
	"github.com/LearningGoProjects/ResourceMonitor/rpc/server/middleware/ratelimit"
	"github.com/LearningGoProjects/ResourceMonitor/service"
	"github.com/LearningGoProjects/ResourceMonitor/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/LearningGoProjects/ResourceMonitor/registry"
	"github.com/LearningGoProjects/ResourceMonitor/rest"
	restSelector "github.com/LearningGoProjects/ResourceMonitor/rest/client/selector"
	"github.com/LearningGoProjects/ResourceMonitor/rpc"
	rpcSelector "github.com/LearningGoProjects/ResourceMonitor/rpc/client/selector"

	"github.com/Allenxuxu/ratelimit/tokenbucket"
)

func NewRPCServer(rg registry.Registry, opt ...rpc.ServerOption) *rpc.Server {
	return rpc.NewServer(rg, opt...)
}

func NewRPCClient(name string, s rpcSelector.Selector, opt ...rpc.ClientOption) (*rpc.Client, error) {
	return rpc.NewClient(name, s, opt...)
}

func NewRestServer(rg registry.Registry, handler http.Handler, opts ...rest.ServerOption) *rest.Server {
	return rest.NewSever(rg, handler, opts...)
}

func NewRestClient(name string, s restSelector.Selector, opt ...rest.ClientOption) (*rest.Client, error) {
	return rest.NewClient(name, s, opt...)
}

func interceptorFun(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	st := time.Now()
	resp, err = handler(ctx, req)

	p, _ := peer.FromContext(ctx)
	log.Printf("method: %s time: %v, peer : %s\n", info.FullMethod, time.Since(st), p.Addr)
	return resp, err
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
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
		ClientAuth:   tls.RequireAndVerifyClientCert, //it is used for mutual TLS, client & server need to provide the certificate to each other
		//ClientAuth: tls.NoClientCert,		//it is used for only server side TLS
		ClientCAs: certPool,
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

	//rg, err := etcd.NewRegistry()
	rg, err := consul.NewRegistry()
	if err != nil {
		panic(err)
	}

	rpcServer := NewRPCServer(rg,
		rpc.Name("ResourceMonitor.CPU"),
		rpc.Version("v1.0.0"),
		rpc.Metadata(map[string]string{
			"server":           "rpc",
			"resource_monitor": "1",
		}),
		rpc.MetricsAddress(":9091"),
		rpc.UnaryServerInterceptor(
			interceptorFun,
			ratelimit.UnaryServerInterceptor(tokenbucket.New(10, 10)),
		),
		rpc.StreamServerInterceptor(
			ratelimit.StreamServerInterceptor(tokenbucket.New(10, 10)),
		),

		rpc.GrpcOpts(serverOptions),
	)

	userStore := service.NewInMemoryUserStore()
	err = seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed users: ", err)
	}

	authServer := service.NewAuthServer(userStore, jwtManager)
	resourceMonitorServer := service.NewResourceMonitorServer()

	// Register the server
	pb.RegisterAuthServiceServer(rpcServer.GrpcServer(), authServer)
	pb.RegisterResourceMonitorServiceServer(rpcServer.GrpcServer(), resourceMonitorServer)
	pb.RegisterProcessorsServiceServer(rpcServer.GrpcServer(), processorsServer)
	//pb.RegisterMemoryServiceServer(rpcServer, memoryServer)
	reflection.Register(rpcServer.GrpcServer())

	// Start sending data to subscribers
	go resourceMonitorServer.StartService()

	log.Printf("Start GRPC server at %s, TLS = %t", listener.Addr().String(), enableTLS)
	return rpcServer.GrpcServer().Serve(listener)
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
	const resourceMonitorServicePath = "/LearningGoProjects.ResourceMonitor.ResourceMonitorService.Server/"

	return map[string][]string{
		resourceMonitorServicePath + "Subscribe":   {"admin", "user"},
		resourceMonitorServicePath + "Unsubscribe": {"admin", "user"},
	}
}
