package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/LearningGoProjects/ResourceMonitor/client"
	"github.com/LearningGoProjects/ResourceMonitor/conf"
	"github.com/LearningGoProjects/ResourceMonitor/log"
	"github.com/LearningGoProjects/ResourceMonitor/pb"
	"github.com/LearningGoProjects/ResourceMonitor/registry/consul"
	"github.com/LearningGoProjects/ResourceMonitor/rpc"
	"github.com/LearningGoProjects/ResourceMonitor/rpc/client/balancer"
	"github.com/LearningGoProjects/ResourceMonitor/rpc/client/selector"
	"github.com/LearningGoProjects/ResourceMonitor/rpc/client/selector/registry"
	"github.com/LearningGoProjects/ResourceMonitor/utils"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"sync"
	"time"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed server's certificate
	pemServerCA, err := ioutil.ReadFile(utils.ClientCACertFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	//Load client's certificate and private key
	clientCert, err := tls.LoadX509KeyPair(utils.ClientCertFile, utils.ClientKeyFile)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {

	if utils.DirectlySubscribing == true {
		SubscribeToServerDirectly()
	} else {
		SubscribeToServerByServcieDiscovery(true)
	}

}

func SubscribeToServerDirectly() {
	var serverAddress string
	var enableTLS bool

	config, err := conf.LoadConfig("./conf/client.yaml")
	//config, err := conf.LoadEnvConfig("./conf/")
	log.Infof("client LoadConfig:", config)

	if err != nil {
		log.Infof("cannot load config:", err)
		log.Infof("Use runtime parameters Now!")
		ss := flag.String("address", "", "the server address")
		tls := flag.Bool("tls", false, "enable SSL/TLS")
		flag.Parse()
		serverAddress = *ss
		enableTLS = *tls
	} else {
		log.Infof("Use YAML config file Now!")
		//serverAddress = config.GetServerAddress()
		//enableTLS = config.GetTLS()
		serverAddress = viper.GetString("server.address") + ":" + viper.GetString("server.port")
		enableTLS = viper.GetBool("tls")
		//serverAddress = config.GetServerAddress()
		//enableTLS = config.GetTLS()
	}

	log.Infof("connecting server %s, TLS = %t", serverAddress, enableTLS)

	transportOption := grpc.WithInsecure()

	if enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}

		transportOption = grpc.WithTransportCredentials(tlsCredentials)
	}

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		clientX, err := client.MKResourceMonitorInterceptorClient(int32(i), serverAddress, transportOption)
		//clientX, err := clientX.MKResourceMonitorClient(int32(i), *serverAddress, transportOption)
		if err != nil {
			log.Fatal(err)
		}
		// Dispatch clientX goroutine
		services := []string{"processor", "memory", "storage"}
		//log.Infof("%%%%%%%%%%%%%%%%", strings.Join(services[:((i-1)%len(services)+1)], ","))

		//go clientX.Start(strings.Join(services[:((i-1)%len(services)+1)],","))
		go clientX.Start(services[:((i-1)%len(services) + 1)]...)

		time.Sleep(time.Second * 2)

	}
	//time.Sleep(time.Second * 30)
	//
	//for i := 3; i <= 5; i++ {
	//	wg.Add(1)
	//	cl, err := client.MKResourceMonitorInterceptorClient(int32(i), *serverAddress, transportOption)
	//
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	cl.Stop()
	//}

	// The wait group purpose is to avoid exiting, the clients do not exit
	wg.Wait()

	//cc1, err := grpc.Dial(*serverAddress, transportOption)
	//if err != nil {
	//	log.Fatal("cannot dial server: ", err)
	//}

	//resourceMonitorClient := client.NewResourceMonitorClient(int32(i), cc1)
	//resourceMonitorClient.Subscribe_OLD()

	//processorsClient := client.NewProcessorsClient(cc1)
	//log.Infof("processorsClient call GetProcessorsInfo RPC")
	//processorsClient.GetProcessorsInfo()
	//
	//log.Infof("processorsClient call SubscribeProcessorInfo RPC")
	//processorsClient.SubscribeProcessorInfo()
	//
	////SubscribeProcessorsInfo stream RPC has some issues:rpc error: code = Internal desc = grpc: failed to unmarshal the received message failed to unmarshal, message is <nil>, want proto.Message
	////log.Infof("processorsClient call SubscribeProcessorsInfo RPC")
	////processorsClient.SubscribeProcessorsInfo()
	//
	//memoryClient := client.NewMemoryClient(cc1)
	//log.Infof("memoryClient call GetMemoryInfo RPC")
	//memoryClient.GetMemoryInfo()

	log.Infof("Resource Monitor Finished! ")
}

func SubscribeToServerByServcieDiscovery(enableTLS bool) {
	rg, err := consul.NewRegistry()
	//rg, err := mdns.NewRegistry()
	//rg, err := etcd.NewRegistry()
	if err != nil {
		panic(err)
	}

	s, err := registry.NewSelector(rg, selector.BalancerName(balancer.RoundRobin))
	/*selector.BalancerName(balancer.RoundRobin) */
	if err != nil {
		panic(err)
	}

	log.Infof("registry.NewSelector ", rg.String(), s.Address("ResourceMonitor.CPU"), err)

	transportOption := grpc.WithInsecure()
	//if enableTLS {
	//	tlsCredentials, err := loadTLSCredentials()
	//	if err != nil {
	//		log.Fatal("cannot load TLS credentials: ", err)
	//	}
	//
	//	transportOption = grpc.WithTransportCredentials(tlsCredentials)
	//}

	//authClient := client.NewAuthClient(conn, username, password)
	//interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	//if err != nil {
	//	log.Fatal("cannot create auth interceptor: ", err)
	//}
	//
	//transportOption = append(transportOption,
	//	grpc.WithUnaryInterceptor(interceptor.Unary()),
	//	grpc.WithStreamInterceptor(interceptor.Stream()))

	client, err := NewRPCClient("ResourceMonitor", s,
		rpc.GrpcDialOption(
			transportOption,
		),
	)
	if err != nil {
		panic(err)
	}
	c := pb.NewRouteGuideClient(client.Conn())

	for i := 0; i < 20; i++ {
		resp, err := c.GetFeature(context.Background(), &pb.Point{
			Latitude:  int32(i),
			Longitude: int32(i),
		})
		if err != nil {
			log.Error(err)
			time.Sleep(time.Second)
			continue
		}

		fmt.Println(i, resp.Name, resp.Location.Latitude, resp.Location.Longitude)
	}

	stream, err := c.RouteChat(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		if err := stream.Send(&pb.RouteNote{
			Location: &pb.Point{
				Latitude:  1,
				Longitude: 1,
			},
			Message: "hello @ " + time.Now().String(),
		}); err != nil {
			panic(err)
		}

		in, err := stream.Recv()
		if err == io.EOF {
			panic(err)

		}
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second * 5)
		log.Infof("[RouteChat] %v received at client @: %s", in, time.Now().String())
	}
}
