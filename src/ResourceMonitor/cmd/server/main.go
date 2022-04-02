package main

import (
	"flag"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/conf"
	"github.com/LearningGoProjects/ResourceMonitor/service"
	"github.com/LearningGoProjects/ResourceMonitor/utils"
	"github.com/spf13/viper"
	"log"
	"net"
)

func main() {

	if utils.DirectlySubscribing == true {
		ServiceProviderBySubscribeDirectly()
	} else {
		ServiceProviderByServiceDiscovery()
	}

}

func ServiceProviderBySubscribeDirectly() {
	var serverAddress, serverType string
	var endPoint string
	var enableTLS bool

	config, err := conf.LoadConfig("./conf/server.yaml")
	//config, err := conf.LoadEnvConfig("./conf/")
	log.Println("client LoadConfig:", config)

	if err != nil {
		log.Println("cannot load config:", err)
		log.Println("Use runtime parameters Now!")
		server_port := flag.Int("port", 0, "the server port")
		tls := flag.Bool("tls", false, "enable SSL/TLS")
		server_type := flag.String("type", "grpc", "type of server (grpc/rest)")
		end_point := flag.String("endpoint", "", "gRPC endpoint")
		flag.Parse()

		serverAddress = fmt.Sprintf("0.0.0.0:%d", *server_port)
		enableTLS = *tls
		serverType = *server_type
		endPoint = *end_point

	} else {
		log.Println("Use YAML config file Now!")
		//serverAddress = config.GetServerAddress()
		//enableTLS = config.GetTLS()
		serverAddress = viper.GetString("server.address") + ":" + viper.GetString("server.port")
		enableTLS = viper.GetBool("tls")
		serverType = viper.GetString("type")
		endPoint = viper.GetString("endpoint")
		//serverAddress = config.GetServerAddress()
		//enableTLS = config.GetTLS()
	}

	//each server's service must be register to grpcserver
	processors_server := service.NewProcessorsServer()
	memory_server := service.NewMemoryServer()

	//listener listen
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	if serverType == "grpc" {
		err = runGRPCServer(processors_server, memory_server, enableTLS, listener)
	} else {
		err = runRESTServer(processors_server, memory_server, enableTLS, listener, endPoint)
	}

	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}

func ServiceProviderByServiceDiscovery() {
	err := runRemodeledGRPCServer(true)
	if err != nil {
		log.Fatal("cannot start ServiceProviderByServiceDiscovery server: ", err)
	}
}
