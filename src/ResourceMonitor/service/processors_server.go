package service

import (
	"context"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func NewCPU() *pb.CPU {

	//CpuId       int32    `protobuf:"varint,1,opt,name=cpuId,proto3" json:"cpuId,omitempty"`
	//VendorId    string   `protobuf:"bytes,2,opt,name=vendorId,proto3" json:"vendorId,omitempty"`
	//ModelName   string   `protobuf:"bytes,3,opt,name=modelName,proto3" json:"modelName,omitempty"`
	//Mhz         float64  `protobuf:"fixed64,4,opt,name=mhz,proto3" json:"mhz,omitempty"`
	//CacheSize   int32    `protobuf:"varint,5,opt,name=cacheSize,proto3" json:"cacheSize,omitempty"`
	//Flags       []string `protobuf:"bytes,6,rep,name=flags,proto3" json:"flags,omitempty"`
	//UsedPercent float64  `protobuf:"fixed64,7,opt,name=usedPercent,proto3" json:"usedPercent,omitempty"`

	cpu := &pb.CPU{
		CpuId:       0,
		VendorId:    "",
		ModelName:   "",
		Mhz:         0,
		CacheSize:   0,
		Flags:       []string{""},
		UsedPercent: 0.0,
	}

	return cpu
}

// NewGPU returns a new sample GPU
func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)

	minGhz := randomFloat64(1.0, 1.5)
	maxGhz := randomFloat64(minGhz, 2.0)
	memGB := randomInt(2, 6)

	gpu := &pb.GPU{
		Brand:  brand,
		Name:   name,
		MinGhz: minGhz,
		MaxGhz: maxGhz,
		Memory: &pb.Memory{
			Value: uint64(memGB),
			Unit:  pb.Memory_GIGABYTE,
		},
	}

	return gpu
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

	ci, err := collectResource()
	if err != nil {
		return nil, err
	}

	res := &pb.GetProcessorsResponse{
		Cpus: ci,
		//Cpus:     []*pb.CPU{NewCPU()},
		//Cpus: []*pb.CPU{ci},
	}
	return res, nil

}

func collectResource() ([]*pb.CPU, error) {
	cpuInfos, _ := GetCpuInfo() //总体信息

	if cpuInfos == nil {
		return nil, fmt.Errorf("no processors info")
	}

	var ci []*pb.CPU = make([]*pb.CPU, len(cpuInfos))

	for i, c := range cpuInfos {
		//log.Println("cpuInfo[", i, "]:", c)

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
	return ci, nil
}

//type ProcessorsServiceSubscribeProcesssorsInfoServer struct {
//	pb.ProcessorsService_SubscribeProcessorsInfoServer
//}
//
//// NewProcessorsServer returns a new ProcessorsServer
//func NewProcessorsServiceSubscribeProcesssorsInfoServer() *ProcessorsServiceSubscribeProcesssorsInfoServer {
//	return &ProcessorsServiceSubscribeProcesssorsInfoServer{}
//}
//SubscribeProcessorsInfo is stream RPC to get processors info
func (server *ProcessorsServer) SubscribeProcessorInfo(
	req *pb.GetProcessorsRequest,
	stream pb.ProcessorsService_SubscribeProcessorInfoServer,
) error { //*pb.GetProcessorsResponse,

	request := req.String()
	log.Printf("service subscribe local processor info %s", request)

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	waitResponse := make(chan error)

	go func() error {
		log.Println("processor server go routine")
		for {
			select {
			case <-ticker.C:
				// do stuff
				log.Println("begin to collect resource info")
				ci, err := collectResource()
				if err != nil {
					waitResponse <- fmt.Errorf("cannot collectResource err: %v", err)
					return err
				}
				log.Println(len(ci), "&&&&&&&&SubscribeProcessorInfo collect resource: ", ci)

				//cpuinfo := &pb.CPU{
				//	VendorId:"GenuineIntel" ,
				//	ModelName:"Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz",
				//	Mhz:2600,
				//	CacheSize:256,
				//
				//	UsedPercent:6.99999999999999}

				//cpuinfo
				res := &pb.GetProcessorsResponse{Cpu: ci[0]}

				log.Println("&pb.GetProcessorsResponse{Cpu: ci}", res)

				err = stream.Send(res)
				if err != nil {
					waitResponse <- fmt.Errorf("cannot send stream response: %v", err)
					return logError(status.Errorf(codes.Unknown, "cannot send stream response: %v", err, stream.SendMsg(nil)))
				}

			case <-quit:
				ticker.Stop()
				return logError(status.Errorf(codes.Aborted, "quit signal received and ticker stop "))
			}
		}
	}()

	err := <-waitResponse
	log.Println("err = <-waitResponse", err)

	return nil
}

//SubscribeProcessorsInfo is stream RPC to get processors info
func (server *ProcessorsServer) SubscribeProcessorsInfo(
	req *pb.GetProcessorsRequest,
	stream pb.ProcessorsService_SubscribeProcessorsInfoServer,
) error { //*pb.GetProcessorsResponse,

	request := req.String()
	log.Printf("service subscribe local processors info %s", request)

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	waitResponse := make(chan error)

	go func() error {
		log.Println("processors server go routine")
		for {
			select {
			case <-ticker.C:
				// do stuff
				log.Println("begin to collect resource info")
				ci, err := collectResource()
				if err != nil {
					waitResponse <- fmt.Errorf("cannot collectResource err: %v", err)
					return err
				}
				log.Println(len(ci), "&&&&&&&&SubscribeProcessorsInfo collect resource: ", ci)

				res := &pb.GetProcessorsResponse{Gpu: NewGPU()}
				//res := &pb.GetProcessorsResponse{Cpus: ci}

				log.Println("&pb.GetProcessorsResponse{Cpu: ci}", res)

				err = stream.Send(res)
				if err != nil {
					waitResponse <- fmt.Errorf("cannot send stream response: %v", err)
					return logError(status.Errorf(codes.Unknown, "cannot send stream response: %v", err, stream.SendMsg(nil)))
				}

			case <-quit:
				ticker.Stop()
				return logError(status.Errorf(codes.Aborted, "quit signal received and ticker stop "))
			}
		}
	}()

	err := <-waitResponse
	log.Println("err = <-waitResponse", err)

	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
