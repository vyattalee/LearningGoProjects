package main

import (
	"github.com/LearningGoProjects/ResourceMonitor/rest"
	restSelector "github.com/LearningGoProjects/ResourceMonitor/rest/client/selector"
	"github.com/LearningGoProjects/ResourceMonitor/rpc"
	rpcSelector "github.com/LearningGoProjects/ResourceMonitor/rpc/client/selector"
)

func NewRPCClient(name string, s rpcSelector.Selector, opt ...rpc.ClientOption) (*rpc.Client, error) {
	return rpc.NewClient(name, s, opt...)
}

func NewRestClient(name string, s restSelector.Selector, opt ...rest.ClientOption) (*rest.Client, error) {
	return rest.NewClient(name, s, opt...)
}
