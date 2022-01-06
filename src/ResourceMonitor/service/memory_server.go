package service

import (
	"context"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/shirou/gopsutil/mem"
	"log"
)

// MemoryServer is the server that provides services
type MemoryServer struct {
	pb.UnimplementedMemoryServiceServer
}

// NewMemoryServer returns a new MemoryServer
func NewMemoryServer() *MemoryServer {
	return &MemoryServer{}
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func ConvertMemory(total uint64) (value uint64, unit pb.Memory_Unit) {
	if total < 1024 {
		//return strconv.FormatInt(total, 10) + "B"
		return (uint64(total) / uint64(1)), pb.Memory_BYTE
	} else if total < (1024 * 1024) {
		return (uint64(total) / uint64(1024)), pb.Memory_KILOBYTE
	} else if total < (1024 * 1024 * 1024) {
		return (uint64(total) / uint64(1024*1024)), pb.Memory_MEGABYTE
	} else if total < (1024 * 1024 * 1024 * 1024) {
		return (uint64(total) / uint64(1024*1024*1024)), pb.Memory_GIGABYTE
	} else if total < (1024 * 1024 * 1024 * 1024 * 1024) {
		return (uint64(total) / uint64(1024*1024*1024*1024)), pb.Memory_TERABYTE
	} else { //if total < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return (uint64(total) / uint64(1024*1024*1024*1024*1024)), pb.Memory_TERABYTE
	}
}

// CreateLaptop is a unary RPC to create a new laptop
func (server *MemoryServer) GetMemoryInfo(
	ctx context.Context,
	req *pb.GetMemoryRequest,
) (*pb.GetMemoryResponse, error) {
	request := req.String()
	log.Printf("service get local Memory info %s", request)

	info, _ := mem.VirtualMemory()
	fmt.Println(info)
	info2, _ := mem.SwapMemory()
	fmt.Println(info2)

	val, unit := ConvertMemory(info.Total)

	mm := &pb.Memory{
		Value: val,
		Unit:  unit,
	}

	res := &pb.GetMemoryResponse{Mem: mm}
	return res, nil
}
