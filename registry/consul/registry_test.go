package consul_test

import (
	"context"
	"fmt"
	"gatesvr/cluster"
	"gatesvr/registry"
	"gatesvr/registry/consul"
	"gatesvr/utils/xnet"

	"net"
	"testing"
	"time"
)

const (
	port        = 3553
	serviceName = "nodestart"
)

var reg = consul.NewRegistry()

func server(t *testing.T) {
	ls, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatal(err)
	}

	go func(ls net.Listener) {
		for {
			conn, err := ls.Accept()
			if err != nil {
				t.Error(err)
				return
			}
			var buff []byte
			if _, err = conn.Read(buff); err != nil {
				t.Error(err)
			}
		}
	}(ls)
}

func TestRegistry_Register(t *testing.T) {
	server(t)

	host, err := xnet.ExternalIP()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	ins := &registry.ServiceInstance{
		ID:       "test-1",
		Name:     serviceName,
		Kind:     cluster.Node.String(),
		Alias:    "mahjong",
		State:    cluster.Work.String(),
		Endpoint: fmt.Sprintf("grpc://%s:%d", host, port),
	}

	for i := 0; i < 100; i++ {
		ins.Routes = append(ins.Routes, registry.Route{
			ID:       int32(i),
			Stateful: true,
			Internal: true,
		})
	}

	if err = reg.Register(ctx, ins); err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	ins.State = cluster.Busy.String()
	if err = reg.Register(ctx, ins); err != nil {
		t.Fatal(err)
	}

	time.Sleep(30 * time.Second)
}

func TestRegistry_Services(t *testing.T) {
	services, err := reg.Services(context.Background(), serviceName)
	if err != nil {
		t.Fatal(err)
	}

	for _, service := range services {
		t.Logf("%+v", service)
	}
}

func TestRegistry_Watch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	watcher1, err := reg.Watch(ctx, serviceName)
	if err != nil {
		t.Fatal(err)
	}

	watcher2, err := reg.Watch(context.Background(), serviceName)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		time.Sleep(5 * time.Second)
		watcher1.Stop()
		time.Sleep(5 * time.Second)
		watcher2.Stop()
	}()

	go func() {
		for {
			services, err := watcher1.Next()
			if err != nil {
				t.Errorf("goroutine 1: %v", err)
				return
			}

			fmt.Println("goroutine 1: new event entity")

			for _, service := range services {
				t.Logf("goroutine 1: %+v", service)
			}
		}
	}()

	go func() {
		for {
			services, err := watcher2.Next()
			if err != nil {
				t.Errorf("goroutine 2: %v", err)
				return
			}

			fmt.Println("goroutine 2: new event entity")

			for _, service := range services {
				t.Logf("goroutine 2: %+v", service)
			}
		}
	}()

	time.Sleep(60 * time.Second)
}
