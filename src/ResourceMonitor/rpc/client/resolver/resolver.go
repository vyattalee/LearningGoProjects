package resolver

import (
	"errors"
	"sync/atomic"

	"github.com/LearningGoProjects/ResourceMonitor/registry"

	"github.com/LearningGoProjects/ResourceMonitor/log"
	"github.com/LearningGoProjects/ResourceMonitor/rpc/client/selector"
	"google.golang.org/grpc/resolver"
)

// "stark-registry:///{service}"
type builder struct {
	scheme   string
	selector *atomic.Value
}

func NewBuilder(scheme string, selector *atomic.Value) resolver.Builder {
	return &builder{
		scheme:   scheme,
		selector: selector,
	}
}

func (d *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	s, ok := d.selector.Load().(selector.Selector)
	if !ok {
		return nil, errors.New("grpc resolver selector is nil")
	}

	sr := &resourceMonitorResolver{
		selector: s,
		cc:       cc,
		service:  target.Endpoint,
	}

	if err := sr.run(); err != nil {
		return nil, err
	}

	return sr, nil
}

func (d *builder) Scheme() string {
	return d.scheme
}

type resourceMonitorResolver struct {
	selector selector.Selector
	cc       resolver.ClientConn
	watcher  registry.Watcher
	service  string
}

func (r *resourceMonitorResolver) run() (err error) {
	if err = r.updateState(); err != nil {
		return err
	}

	r.watcher, err = r.selector.Watch(r.service)
	if err != nil {
		return err
	}

	// for static selector
	if r.watcher == nil {
		return nil
	}

	go func() {
		for {
			_, err := r.watcher.Next()
			if err != nil {
				// watcher close
				return
			}

			if err := r.updateState(); err != nil {
				log.Errorf("resource monitor resolver update state error: %v", err)
			}
		}
	}()

	return nil
}

func (r *resourceMonitorResolver) updateState() error {
	services, err := r.selector.GetService(r.service)
	if err != nil {
		return err
	}
	var status resolver.State
	for _, s := range services {
		for _, node := range s.Nodes {
			status.Addresses = append(status.Addresses, resolver.Address{Addr: node.Address})
		}
	}

	r.cc.UpdateState(status)

	return nil
}

func (r *resourceMonitorResolver) Close() {
	if r.watcher != nil {
		r.watcher.Stop()
	}
	if r.selector != nil {
		_ = r.selector.Close()
	}
}

func (r *resourceMonitorResolver) ResolveNow(options resolver.ResolveNowOptions) {}
