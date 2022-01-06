package service

import (
	"context"
	"errors"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

// ProcessorsServer is the server that provides services
type ProcessorsServer struct {
	pb.UnimplementedProcessorsServiceServer
}

// NewProcessorsServer returns a new ProcessorsServer
func NewProcessorsServer() *ProcessorsServer {
	return &ProcessorsServer{}
}

// CreateLaptop is a unary RPC to create a new laptop
func (server *ProcessorsServer) GetProcessorsInfo(
	ctx context.Context,
	req *pb.GetProcessorsRequest,
) (*pb.GetProcessorsResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		// check if it's a valid UUID
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		laptop.Id = id.String()
	}

	// some heavy processing
	// time.Sleep(6 * time.Second)

	if err := contextError(ctx); err != nil {
		return nil, err
	}

	// save the laptop to store
	err := server.laptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}

		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("saved laptop with id: %s", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}
