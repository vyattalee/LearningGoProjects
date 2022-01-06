package client

import (
	"context"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// CreateLaptop calls create laptop RPC
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
	gpuinfo := res.GetGpu()

	log.Printf("processors info  \r\n -->cpu: %s,\r\n -->gpu: %s", cpuinfo, gpuinfo)
}
