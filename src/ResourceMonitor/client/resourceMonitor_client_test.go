package client_test

import (
	"github.com/LearningGoProjects/ResourceMonitor/client"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/LearningGoProjects/ResourceMonitor/service"
	"github.com/LearningGoProjects/ResourceMonitor/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"testing"
	"time"
)

//func TestNewResourceMonitorClient(t *testing.T) {
//	t.Parallel()
//
//	var ExpectedId int32 = 1
//	//serverAddress := startTestLaptopServer()
//	cc1, err := grpc.Dial("0.0.0.0", grpc.WithInsecure())
//	if err != nil {
//		log.Fatal("cannot dial server: ", err)
//	}
//
//	resourceMonitorClient := client.NewResourceMonitorClient(ExpectedId, cc1)
//
//	require.Equal(t, &client.ResourceMonitorClient{
//		service: pb.NewResourceMonitorServiceClient(cc1),
//		id:ExpectedId,
//		conn:cc1},
//		resourceMonitorClient)
//
//}

func TestResourceMonitorClient_Start(t *testing.T) {

	userStore := service.NewInMemoryUserStore()

	serverAddress, resouceMonitorServer := startTestResourceMonitorServer(t, userStore)
	transportOption := grpc.WithInsecure()

	var ExpectedId uint32 = 1

	clientX, err := client.MKResourceMonitorInterceptorClient(ExpectedId, serverAddress, transportOption)
	if err != nil {
		log.Fatal(err)
	}
	require.NoError(t, err)

	go clientX.Start("processor", "memory", "storage")

	require.NotNil(t, resouceMonitorServer.GetSubscribers())
	other, err := userStore.Find("admin1")
	require.NotNil(t, other)
}

func startTestResourceMonitorServer(t *testing.T, userStore service.UserStore) (string, *service.ResourceMonitorServer) {
	serverOptions := []grpc.ServerOption{
		grpc.ConnectionTimeout(30 * time.Second),
		grpc.MaxConcurrentStreams(10),
	}

	grpcServer := grpc.NewServer(serverOptions...)

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
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener.Addr().String(), resourceMonitorServer

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
