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

// MemoryClient is a client to call Memory service RPCs
type MemoryClient struct {
	service pb.MemoryServiceClient
}

// NewMemoryClient returns a new Memory client
func NewMemoryClient(cc *grpc.ClientConn) *MemoryClient {
	service := pb.NewMemoryServiceClient(cc)
	return &MemoryClient{service}
}

// CreateLaptop calls create laptop RPC
func (memoryClient *MemoryClient) GetMemoryInfo() {
	req := &pb.GetMemoryRequest{}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := memoryClient.service.GetMemoryInfo(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			// not a big deal
			log.Print("Memory already exists")
		} else {
			log.Fatal("cannot get Memory info: ", err)
		}
		return
	}

	meminfo := res.GetMem()

	log.Printf("Memory info  \r\n -->memory: %s", meminfo)

}
