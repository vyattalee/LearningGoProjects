package service

import (
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"
	"log"
	"sync"
	"time"
)

// ProcessorsServer is the server that provides services
type ResourceMonitorServer struct {
	pb.UnimplementedResourceMonitorServiceServer
	subscribers sync.Map // subscribers is a concurrent map that holds mapping from a client ID to it's subscriber
}

// NewProcessorsServer returns a new ProcessorsServer
func NewResourceMonitorServer() *ResourceMonitorServer {
	return &ResourceMonitorServer{}
}

type sub struct {
	stream   pb.ResourceMonitorService_SubscribeServer // stream is the server side of the RPC stream
	finished chan<- bool                               // finished is used to signal closure of a client subscribing goroutine
}

func (server *ResourceMonitorServer) Subscribe(
	req *pb.Request,
	stream pb.ResourceMonitorService_SubscribeServer,
) error {

	//switch req.GetServiceType() {
	//case pb.ServiceType_ProcessorService.Enum():
	//	//
	//	log.Println("Subscribe processor service")
	//
	//case pb.ServiceType_MemoryService.Enum():
	//	log.Println("Subscribe memory service")
	//
	//}

	// Handle subscribe request
	log.Printf("Received subscribe request from ID: %d", req.Id)

	fin := make(chan bool)
	// Save the subscriber stream according to the given client ID
	server.subscribers.Store(req.Id, sub{stream: stream, finished: fin})

	ctx := stream.Context()
	// Keep this scope alive because once this scope exits - the stream is closed
	for {
		select {
		case <-fin:
			log.Printf("Closing stream for client ID: %d", req.Id)
			return nil
		case <-ctx.Done():
			log.Printf("Client ID %d has disconnected", req.Id)
			return nil
		}
	}

}

func (server *ResourceMonitorServer) StartService() {
	log.Println("Starting resource monitor background service")
	for {
		time.Sleep(time.Second)

		// A list of clients to unsubscribe in case of error
		var unsubscribe []int32

		// Iterate over all subscribers and send data to each client
		server.subscribers.Range(func(k, v interface{}) bool {
			id, ok := k.(int32)
			if !ok {
				log.Printf("Failed to cast subscriber key: %T", k)
				return false
			}
			sub, ok := v.(sub)
			if !ok {
				log.Printf("Failed to cast subscriber value: %T", v)
				return false
			}
			// Send data over the gRPC stream to the client
			//s1 := fmt.Sprintf("data mock for: %d", id)
			t1 := &timestamp.Timestamp{}
			serialized, err := proto.Marshal(t1)
			if err != nil {
				log.Fatal("could not serialize string")
			}

			anydata := &any.Any{
				TypeUrl: "example.com/resourcemonitor/" + proto.MessageName(t1),
				Value:   serialized,
			}

			if err := sub.stream.Send(&pb.Response{ResourceData: anydata}); err != nil {
				log.Printf("Failed to send data to client: %v", err)
				select {
				case sub.finished <- true:
					log.Printf("Unsubscribed client: %d", id)
				default:
					// Default case is to avoid blocking in case client has already unsubscribed
				}
				// In case of error the client would re-subscribe so close the subscriber stream
				unsubscribe = append(unsubscribe, id)
			}
			return true
		})

		// Unsubscribe erroneous client streams
		for _, id := range unsubscribe {
			server.subscribers.Delete(id)
		}
	}
}
