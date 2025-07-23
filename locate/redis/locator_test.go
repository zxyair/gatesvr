package redis_test

import (
	"context"
	"fmt"
	"gatesvr/cluster"
	"gatesvr/etc"
	"gatesvr/locate/redis"
	"gatesvr/log"
	"gatesvr/utils/xuuid"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var locator = redis.NewLocator(
	redis.WithAddrs("127.0.0.1:6379"),
	redis.WithPassword("123456"),
)

const (
	defaultAddrsKey      = "etc.locate.redis.addrs"
	defaultDBKey         = "etc.locate.redis.db"
	defaultMaxRetriesKey = "etc.locate.redis.maxRetries"
	defaultPrefixKey     = "etc.locate.redis.prefix"
	defaultUsernameKey   = "etc.locate.redis.username"
	defaultPasswordKey   = "etc.locate.redis.password"
)

func printWorkingDir() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("当前工作目录:", dir)

	absPath, _ := filepath.Abs("./etc")
	fmt.Println("配置文件绝对路径:", absPath)
}
func TestNewLocator(t *testing.T) {
	printWorkingDir()
	log.Info(etc.Get(defaultAddrsKey).Strings())
	log.Info(etc.Get(defaultDBKey))
	log.Info(etc.Get(defaultMaxRetriesKey))
	log.Info(etc.Get(defaultPrefixKey))
	log.Info(etc.Get(defaultUsernameKey))
	log.Info(etc.Get(defaultPasswordKey))

	//if locator == nil {
	//	t.Fatal("locator is nil")
	//}
}
func TestLocator_BindGate(t *testing.T) {
	ctx := context.Background()
	uid := int64(1)
	gid := xuuid.UUID()

	if err := locator.BindGate(ctx, uid, gid); err != nil {
		t.Fatal(err)
	}
}

func TestLocator_BindNode(t *testing.T) {
	ctx := context.Background()
	uid := int64(1)
	nid := xuuid.UUID()
	name := "node1"

	if err := locator.BindNode(ctx, uid, name, nid); err != nil {
		t.Fatal(err)
	}
}

func TestLocator_UnbindGate(t *testing.T) {
	ctx := context.Background()
	uid := int64(1)
	gid := xuuid.UUID()

	if err := locator.BindGate(ctx, uid, gid); err != nil {
		t.Fatal(err)
	}

	if err := locator.UnbindGate(ctx, uid, gid); err != nil {
		t.Fatal(err)
	}
}

func TestLocator_UnbindNode(t *testing.T) {
	ctx := context.Background()
	uid := int64(1)
	nid1 := xuuid.UUID()
	nid2 := xuuid.UUID()
	name1 := "node1"
	name2 := "node2"

	if err := locator.BindNode(ctx, uid, name1, nid1); err != nil {
		t.Fatal(err)
	}

	if err := locator.BindNode(ctx, uid, name2, nid2); err != nil {
		t.Fatal(err)
	}

	if err := locator.UnbindNode(ctx, uid, name2, nid2); err != nil {
		t.Fatal(err)
	}
}

func TestLocator_Watch(t *testing.T) {
	watcher1, err := locator.Watch(context.Background(), cluster.Gate.String(), cluster.Node.String())
	if err != nil {
		t.Fatal(err)
	}

	watcher2, err := locator.Watch(context.Background(), cluster.Gate.String())
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			events, err := watcher1.Next()
			if err != nil {
				t.Errorf("goroutine 1: %v", err)
				return
			}

			fmt.Println("goroutine 1: new event entity")

			for _, event := range events {
				t.Logf("goroutine 1: %+v", event)
			}
		}
	}()

	go func() {
		for {
			events, err := watcher2.Next()
			if err != nil {
				t.Errorf("goroutine 2: %v", err)
				return
			}

			fmt.Println("goroutine 2: new event entity")

			for _, event := range events {
				t.Logf("goroutine 2: %+v", event)
			}
		}
	}()

	time.Sleep(60 * time.Second)
}
