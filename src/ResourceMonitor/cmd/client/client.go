package main

import (
	"flag"
	"github.com/LearningGoProjects/ResourceMonitor/client"
	"google.golang.org/grpc"
	"log"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	enableTLS := flag.Bool("tls", false, "enable SSL/TLS")
	flag.Parse()
	log.Printf("connecting server %s, TLS = %t", *serverAddress, *enableTLS)

	transportOption := grpc.WithInsecure()

	//if *enableTLS {
	//	tlsCredentials, err := loadTLSCredentials()
	//	if err != nil {
	//		log.Fatal("cannot load TLS credentials: ", err)
	//	}
	//
	//	transportOption = grpc.WithTransportCredentials(tlsCredentials)
	//}

	cc1, err := grpc.Dial(*serverAddress, transportOption)
	if err != nil {
		log.Fatal("cannot connecting server: ", err)
	}

	processorsClient := client.NewProcessorsClient(cc1)
	log.Println("processorsClient call GetProcessorsInfo RPC")
	processorsClient.GetProcessorsInfo()

	log.Println("processorsClient call SubscribeProcessorInfo RPC")
	processorsClient.SubscribeProcessorInfo()

	//SubscribeProcessorsInfo stream RPC has some issues:rpc error: code = Internal desc = grpc: failed to unmarshal the received message failed to unmarshal, message is <nil>, want proto.Message
	//log.Println("processorsClient call SubscribeProcessorsInfo RPC")
	//processorsClient.SubscribeProcessorsInfo()

	memoryClient := client.NewMemoryClient(cc1)
	log.Println("memoryClient call GetMemoryInfo RPC")
	memoryClient.GetMemoryInfo()

	log.Println("Resource Monitor Finished! ")

}
