package service

import (
	"encoding/json"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/LearningGoProjects/ResourceMonitor/utils"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/shirou/gopsutil/mem"
	"log"
	"sync"
	"time"
)

// ProcessorsServer is the server that provides services
type ResourceMonitorServer struct {
	pb.UnimplementedResourceMonitorServiceServer
	subscribers sync.Map // subscribers is a concurrent map that holds mapping from a client ID to it's subscriber
	ticker      *time.Ticker
}

// NewProcessorsServer returns a new ProcessorsServer
func NewResourceMonitorServer() *ResourceMonitorServer {
	return &ResourceMonitorServer{ticker: time.NewTicker(5 * time.Second)}
}

type sub struct {
	stream       pb.ResourceMonitorService_SubscribeServer // stream is the server side of the RPC stream
	finished     chan<- bool                               // finished is used to signal closure of a client subscribing goroutine
	sub_services *utils.BitMap                             //utils.BitMap
}

func (server *ResourceMonitorServer) Subscribe(
	req *pb.Request,
	stream pb.ResourceMonitorService_SubscribeServer,
) error {

	// Handle subscribe request
	log.Printf("Received subscribe request from ID: %d", req.Id)

	fin := make(chan bool)
	// Save the subscriber stream according to the given client ID
	server.subscribers.Store(req.Id, sub{stream: stream, finished: fin, sub_services: utils.NewBitMapFromBytes(req.Filter.SubService)})

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
		//time.Sleep(time.Second)

		quit := make(chan struct{})
		//waitResponse := make(chan error)

		select {
		case <-server.ticker.C:
			// do stuff
			log.Println("begin to do ticker stuff")
			server.doTickerJobs(quit)

		case <-quit:
			server.ticker.Stop()
			break

		}

	}
}

func (server *ResourceMonitorServer) doTickerJobs(quit chan struct{}) {

	// A list of clients to unsubscribe in case of error
	var unsubscribe []int32
	//subsciber length, if subscriber_len is zero, quit chan<- struct{}{} to stop time.Ticker
	//var subscriber_len int = 0

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

		//if client id is odd, then send the cpu info to client
		//if client id is even, then send the memory info to client
		var err error
		var cpu []*pb.CPU
		//var memory *pb.Memory
		//var response *pb.Response
		var byteData []byte

		log.Println("$$$$$$$$$$$$$client:", id, "	subscribe services:", sub.sub_services.String())
		if sub.sub_services.Bit(1<<pb.ServiceType_ProcessorService-1) == utils.IsSet {
			cpu, err = CollectResource()
			if err != nil {

				return false
			}
			resource := &pb.Response_Cpu{Cpu: cpu[0]}

			byteData, err = json.Marshal(resource)
			if err != nil {
				log.Printf("Could not convert data input to bytes")
			}
			resourceData := &any.Any{
				TypeUrl: "anyResourceData_cpu",
				Value:   byteData,
			}
			err = sub.stream.Send(&pb.Response{
				ResourceData:    fmt.Sprintf("data mock for: %d", id),
				Resource:        resource,
				AnyResourceData: resourceData,
			})

		} else if sub.sub_services.Bit(1<<pb.ServiceType_MemoryService-1) == utils.IsSet {
			info, _ := mem.VirtualMemory()
			val, unit := ConvertMemory(info.Total)
			resource := &pb.Response_Memory{
				Memory: &pb.Memory{Value: val, Unit: unit},
			}

			byteData, err = json.Marshal(resource)
			if err != nil {
				log.Printf("Could not convert data input to bytes")
			}
			resourceData := &any.Any{
				TypeUrl: "anyResourceData_memory",
				Value:   byteData,
			}
			err = sub.stream.Send(&pb.Response{
				ResourceData:    fmt.Sprintf("data mock for: %d", id),
				Resource:        resource,
				AnyResourceData: resourceData,
			})
		}
		//else {
		//	log.Printf("not support resource service type: ", sub.sub_services.Xor(utils.NewBitMapFromString("11")))
		//}
		//whether Can it be integrated into the following single function
		//if err = sub.stream.Send(&pb.Response{
		//	ResourceData: fmt.Sprintf("data mock for: %d", id),
		//	Resource: func()*pb.{
		//		return "anonymous stringy\n"
		//	};},
		//}); err != nil {
		if err != nil {
			//if err := sub.stream.Send(&pb.Response{ResourceData: anydata}); err != nil {
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
