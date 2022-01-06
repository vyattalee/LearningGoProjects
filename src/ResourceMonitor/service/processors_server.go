package service

import (
	"context"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	//"github.com/shirou/gopsutil"
	"github.com/shirou/gopsutil/cpu"
	"time"

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

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

// CreateLaptop is a unary RPC to create a new laptop
func (server *ProcessorsServer) GetProcessorsInfo(
	ctx context.Context,
	req *pb.GetProcessorsRequest,
) (*pb.GetProcessorsResponse, error) {
	request := req.String()
	log.Printf("service get local processors info %s", request)

	info, _ := cpu.Info() //总体信息
	fmt.Println(info)
	//output：	[{"cpu":0,cores":4,"modelName":"Intel(R) Core(TM) i5-2520M CPU @ 2.50GHz","mhz":2501,。。。]

	ci := &pb.CPU{
		Name:          info[0].ModelName,
		Brand:         info[0].VendorID,
		NumberCores:   uint32(info[0].Cores),
		NumberThreads: uint32(info[0].Cores),
		MaxGhz:        info[0].Mhz,
		MinGhz:        info[0].Mhz,
	}

	res := &pb.GetProcessorsResponse{Cpu: ci}
	return res, nil
}
