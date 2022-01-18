package client

import (
	"context"
	"encoding/json"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/LearningGoProjects/ResourceMonitor/utils"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"time"
)

// ResourceMonitorClient is a client to call processors service RPCs
type ResourceMonitorClient struct {
	service      pb.ResourceMonitorServiceClient
	conn         *grpc.ClientConn // conn is the client gRPC connection
	id           int32            // id is the client ID used for subscribing
	sub_services *utils.BitMap
}

// NewResourceMonitorClient returns a new processors client
func NewResourceMonitorClient(id int32, cc *grpc.ClientConn) *ResourceMonitorClient {
	service := pb.NewResourceMonitorServiceClient(cc)
	return &ResourceMonitorClient{service: service, id: id, sub_services: utils.NewBitMap(8)}
}

func (resourceMonitorClient *ResourceMonitorClient) Start(sub_services ...string) {
	var err error
	// stream is the client side of the RPC stream
	var stream pb.ResourceMonitorService_SubscribeClient
	for {
		if stream == nil {
			if stream, err = resourceMonitorClient.subscribe(sub_services...); err != nil {
				log.Printf("Failed to subscribe: %v", err)
				resourceMonitorClient.sleep()
				// Retry on failure
				continue
			}
		}
		response, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive message: %v", err)
			// Clearing the stream will force the client to resubscribe on next iteration
			stream = nil
			resourceMonitorClient.sleep()
			// Retry on failure
			continue
		}

		//map[string]interface{} is good for everything
		var data map[string]interface{}

		err = json.Unmarshal(response.AnyResourceData.GetValue(), &data)
		if err != nil {
			log.Print("Error while unmarshaling the response.AnyResourceData.GetValue()")
		}

		//TypeUrl can be used for select which kind of anyResourceData type.
		unmarshal := &pb.Response{}
		switch response.AnyResourceData.GetTypeUrl() {
		case "anyResourceData_cpu":
			ptypes.UnmarshalAny(response.AnyResourceData, unmarshal)
		case "anyResourceData_memory":
			ptypes.UnmarshalAny(response.AnyResourceData, unmarshal)

		}

		switch response.Resource.(type) {
		case *pb.Response_Cpu:
			log.Printf("Client ID %d got response: %q", resourceMonitorClient.id, response.ResourceData,
				"\r\nresource-->CPU info:", response.GetCpu(), "\r\n #########Any type Data:", data) //response.GetAnyResourceData()
		case *pb.Response_Memory:
			log.Printf("Client ID %d got response: %q", resourceMonitorClient.id, response.ResourceData,
				"\r\nresource-->Memory info:", response.GetMemory(), "\r\n #########Any type Data:", data)
		}

	}
}

// subscribe subscribes to messages from the gRPC server
func (resourceMonitorClient *ResourceMonitorClient) subscribe(sub_services ...string) (pb.ResourceMonitorService_SubscribeClient, error) {
	log.Printf("Subscribing client ID: %d", resourceMonitorClient.id, "	Subscribe Service:", sub_services)
	for _, sub_service := range sub_services {
		log.Printf("^^^^^^^^^^^^^^", sub_service)
		if sub_service == "processor" {
			resourceMonitorClient.sub_services.Set(int(pb.ServiceType_ProcessorService))
		} else if sub_service == "memory" {
			resourceMonitorClient.sub_services.Set(int(pb.ServiceType_MemoryService))
		} else if sub_service == "storage" {
			resourceMonitorClient.sub_services.Set(int(pb.ServiceType_StorageService))
		}
	}
	log.Println("**************", resourceMonitorClient.sub_services.String())
	return resourceMonitorClient.service.Subscribe(context.Background(),
		&pb.Request{Id: resourceMonitorClient.id, Filter: &pb.Filter{SubService: resourceMonitorClient.sub_services.Bytes()}}) //, Filter: &pb.Filter{ServiceType: &pb.ServiceType.ProcessorService}
}

// sleep is used to give the server time to unsubscribe the client and reset the stream
func (resourceMonitorClient *ResourceMonitorClient) sleep() {
	time.Sleep(time.Second * 5)
}

func mkConnection(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(target, opts...)
}

// MKResourceMonitorClient creates a new client instance
func MKResourceMonitorClient(id int32, target string, opts ...grpc.DialOption) (*ResourceMonitorClient, error) {
	conn, err := mkConnection(target, opts...)
	if err != nil {
		return nil, err
	}
	return &ResourceMonitorClient{
		service:      pb.NewResourceMonitorServiceClient(conn),
		conn:         conn,
		id:           id,
		sub_services: utils.NewBitMap(8),
	}, nil
}

const (
	username        = "admin1"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const resourceMonitorServicePath = "/LearningGoProjects.ResourceMonitor.ResourceMonitorService/"

	return map[string]bool{
		//resourceMonitorServicePath + "Subscribe":   true,
		resourceMonitorServicePath + "Unsubscribe": true,
	}
}

// MKResourceMonitorClient creates a new client instance
func MKResourceMonitorInterceptorClient(id int32, target string, opts ...grpc.DialOption) (*ResourceMonitorClient, error) {
	conn, err := mkConnection(target, opts...)
	if err != nil {
		return nil, err
	}

	authClient := NewAuthClient(conn, username, password)
	interceptor, err := NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	transportOption := grpc.WithInsecure()

	conn2, err := grpc.Dial(
		target,
		transportOption,
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	return &ResourceMonitorClient{
		service:      pb.NewResourceMonitorServiceClient(conn2),
		conn:         conn2,
		id:           id,
		sub_services: utils.NewBitMap(8),
	}, nil
}

// close is not used but is here as an example of how to close the gRPC client connection
func (resourceMonitorClient *ResourceMonitorClient) close() {
	if err := resourceMonitorClient.conn.Close(); err != nil {
		log.Fatal(err)
	}
}
