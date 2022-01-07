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

func GetCpuInfo() ([]cpu.InfoStat, error) {
	cpuInfos, err := cpu.Info()
	if err != nil {
		fmt.Printf("get cpu info failed, err:%v", err)
		return nil, err
	}

	return cpuInfos, nil

}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

// GetProcessorsInfo is a unary RPC to get processors info
func (server *ProcessorsServer) GetProcessorsInfo(
	ctx context.Context,
	req *pb.GetProcessorsRequest,
) (*pb.GetProcessorsResponse, error) {
	request := req.String()
	log.Printf("service get local processors info %s", request)

	cpuInfos, _ := GetCpuInfo() //总体信息
	//output：	[{"cpu":0,cores":4,"modelName":"Intel(R) Core(TM) i5-2520M CPU @ 2.50GHz","mhz":2501,。。。]

	if cpuInfos == nil {
		return nil, fmt.Errorf("no processors info")
	}

	var ci []*pb.CPU
	for i, c := range cpuInfos {
		if i < len(cpuInfos) {

			ci[i] = &pb.CPU{
				CpuId:       c.CPU,
				VendorId:    c.VendorID,
				ModelName:   c.ModelName,
				Mhz:         c.Mhz,
				CacheSize:   c.CacheSize,
				Flags:       c.Flags,
				UsedPercent: GetCpuPercent(),
			}
		}
	}

	res := &pb.GetProcessorsResponse{Cpu: ci}
	return res, nil

}

type ProcessorsServiceSubscribeProcesssorsInfoServer struct {
	pb.ProcessorsService_SubscribeProcesssorsInfoServer
}

//SubscribeProcessorsInfo is stream RPC to get processors info
func (server *ProcessorsServiceSubscribeProcesssorsInfoServer) SubscribeProcesssorsInfo(
	ctx context.Context,
	req *pb.GetProcessorsRequest,
) (*pb.GetProcessorsResponse, error) {

	request := req.String()
	log.Printf("service subscribe local processors info %s", request)

	res := &pb.GetProcessorsResponse{}
	return res, nil
}
