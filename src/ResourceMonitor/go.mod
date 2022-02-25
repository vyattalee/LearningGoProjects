module github.com/LearningGoProjects/ResourceMonitor

go 1.15

//replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

require (
	github.com/Allenxuxu/ratelimit v0.0.0-20210131080358-1c878c80259b
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1
	github.com/go-log/log v0.1.0
	github.com/go-resty/resty/v2 v2.7.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.1.3
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hashicorp/consul/api v1.12.0
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/miekg/dns v1.1.41
	github.com/mitchellh/hashstructure v1.1.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.5.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.etcd.io/etcd v0.0.0-20200402134248-51bdeb39e698
	//go.etcd.io/etcd v2.3.8+incompatible
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	golang.org/x/sys v0.0.0-20220224120231-95c6836cb0e7 // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.27.1
	sigs.k8s.io/yaml v1.2.0 // indirect

)

//replace github.com/golang/protobuf => github.com/golang/protobuf v1.4.3
