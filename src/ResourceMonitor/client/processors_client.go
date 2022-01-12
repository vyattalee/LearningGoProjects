package client

import (
	"context"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

// ProcessorsClient is a client to call processors service RPCs
type ProcessorsClient struct {
	service pb.ProcessorsServiceClient
}

// NewProcessorsClient returns a new processors client
func NewProcessorsClient(cc *grpc.ClientConn) *ProcessorsClient {
	service := pb.NewProcessorsServiceClient(cc)
	return &ProcessorsClient{service}
}

// ProcessorsClient calls GetProcessorsInfo RPC
func (processorsClient *ProcessorsClient) GetProcessorsInfo() {
	req := &pb.GetProcessorsRequest{}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := processorsClient.service.GetProcessorsInfo(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			// not a big deal
			log.Print("processors already exists")
		} else {
			log.Fatal("cannot get processors info: ", err)
		}
		return
	}

	cpuinfo := res.GetCpu()
	//cpuinfo := res.GetCpus()
	gpuinfo := res.GetGpu()

	log.Printf("processors info  \r\n -->cpu: %s,\r\n -->gpu: %s", cpuinfo, gpuinfo)
}

// ProcessorsClient calls SubscribeProcessorInfo RPC
func (processorsClient *ProcessorsClient) SubscribeProcessorInfo() error {

	req := &pb.GetProcessorsRequest{}

	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	log.Println("processorsClient.service.SubscribeProcessorInfo(ctx, req)")

	stream, err := processorsClient.service.SubscribeProcessorInfo(ctx, req)
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
				continue
			}

			log.Println("received SubscribeProcessorInfo response: ", res.GetCpu(), stream.RecvMsg(nil))
		}
	}()

	err = <-waitResponse
	log.Println("err = <-waitResponse", err)
	return err
}

// ProcessorsClient calls SubscribeProcessorsInfo RPC
//func (processorsClient *ProcessorsClient) SubscribeProcessorsInfo() error {
//
//	req := &pb.GetProcessorsRequest{}
//
//	ctx := context.Background()
//	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	//defer cancel()
//
//	log.Println("processorsClient.service.SubscribeProcessorsInfo(ctx, req)")
//
//	stream, err := processorsClient.service.SubscribeProcessorsInfo(ctx, req)
//	if err != nil {
//		return fmt.Errorf("cannot subscrible processors: %v", err)
//	}
//
//	waitResponse := make(chan error)
//	// go routine to receive responses
//	go func() {
//		log.Println("processors client go routine")
//		for {
//			res, err := stream.Recv()
//			if err == io.EOF {
//				log.Print("no more responses")
//				waitResponse <- nil
//				return
//			}
//			if err != nil {
//				stream.CloseSend()
//				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
//				return
//			}
//
//			log.Println("received SubscribeProcessorsInfo response cpus: ", res.GetCpus(), stream.RecvMsg(nil))
//			log.Println("received SubscribeProcessorsInfo response gpu: ", res.GetGpu(), stream.RecvMsg(nil))
//		}
//	}()
//
//	err = <-waitResponse
//	log.Println("err = <-waitResponse", err)
//	return err
//}
