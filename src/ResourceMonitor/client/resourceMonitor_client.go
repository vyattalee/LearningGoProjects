package client

import (
	"context"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"google.golang.org/grpc"
	"io"
	"log"
)

// ResourceMonitorClient is a client to call processors service RPCs
type ResourceMonitorClient struct {
	service pb.ResourceMonitorServiceClient
	id      int32 // id is the client ID used for subscribing
}

// NewResourceMonitorClient returns a new processors client
func NewResourceMonitorClient(id int32, cc *grpc.ClientConn) *ResourceMonitorClient {
	service := pb.NewResourceMonitorServiceClient(cc)
	return &ResourceMonitorClient{service: service, id: id}
}

// ProcessorsClient calls SubscribeProcessorInfo RPC
func (resourceMonitorClient *ResourceMonitorClient) Subscribe() error {
	req := &pb.Request{}

	ctx := context.Background()

	log.Println("resourceMonitorClient.service.Subscribe(ctx, req)")

	stream, err := resourceMonitorClient.service.Subscribe(ctx, req)
	if err != nil {
		return fmt.Errorf("cannot subscrible processor: %v", err)
	}

	waitResponse := make(chan error)
	// go routine to receive responses
	go func() {
		log.Println("processors client go routine")
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("no more responses")
				waitResponse <- nil
				return
			}
			if err != nil {
				stream.CloseSend()
				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
				return
			}

			// unmarshal to simulate coming off the wire

			// unmarshal the timestamp
			//var t2 timestamp.Timestamp
			//if err := ptypes.UnmarshalAny(res.ResourceData, &t2); err != nil {
			//	log.Fatalf("Could not unmarshal timestamp from anything field: %s", err)
			//}
			log.Println("received resourceMonitorClient.service.Subscribe response: ", res.ResourceData, stream.RecvMsg(nil))
		}
	}()

	err = <-waitResponse
	log.Println("err = <-waitResponse", err)
	return err
}
