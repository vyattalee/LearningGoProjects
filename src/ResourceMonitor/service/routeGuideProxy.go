package service

import (
	"context"
	"encoding/json"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	//"github.com/shirou/gopsutil/cpu"

	"io"
	"log"
	"time"
)

type routeGuideServer struct {
	pb.UnimplementedRouteGuideServer
}

func NewServer() *routeGuideServer {
	return &routeGuideServer{}
}

func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	log.Println("[GetFeature]", point.Latitude)
	return &pb.Feature{Location: point}, nil
}

func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	log.Printf("[ListFeatures] %v", rect)

	for i := 0; i < 10; i++ {
		if err := stream.Send(&pb.Feature{
			Name: "feature",
			Location: &pb.Point{
				Latitude:  int32(i),
				Longitude: int32(i),
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}

		log.Printf("[RecordRoute] %v", point)
	}
}

func (s *routeGuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
	for {
		//in, err := stream.Recv()
		//if err == io.EOF {
		//	return nil
		//}
		//if err != nil {
		//	return err
		//}

		//log.Printf("[RouteChat] %v", in)

		//var byteData []byte
		storage, _ := GetStorageInfo()
		//byteData, err = json.Marshal(storage)
		//if err != nil {
		//	log.Printf("GetStorageInfo() error:", err)
		//}

		cpu, _ := CollectCPUGPUResource()
		j4cpu, _ := json.Marshal(cpu)

		mem, _ := GetMemoryInfo()

		if err := stream.Send(&pb.RouteNote{
			//Location: in.Location,
			Message: storage.String() + string(j4cpu) + mem.String(),
		}); err != nil {
			return err
		}
	}
}

//func (s *routeGuideServer) mustEmbedUnimplementedRouteGuideServer() {
//
//}
