package benchtest

import (
	"fmt"
	"gatesvr/cluster"
	"gatesvr/cluster/client"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/clientstart/pojo"
	"gatesvr/functiontest/send-and-receive-message-compress-encrypto/route"
	"gatesvr/log"
	"gatesvr/utils/xrand"
	"gatesvr/utils/xtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	wg        *sync.WaitGroup
	startTime int64
	totalSent int64
	totalRecv int64
	message   string
)

func TestBenchMark(proxy *client.Proxy) {
	samples := []struct {
		c    int // 并发数
		n    int // 请求数
		size int // 数据包大小
	}{
		{
			c:    50,
			n:    1000000,
			size: 1024,
		},
		{
			c:    100,
			n:    1000000,
			size: 1024,
		},
		{
			c:    200,
			n:    1000000,
			size: 1024,
		},
		{
			c:    300,
			n:    1000000,
			size: 1024,
		},
		{
			c:    400,
			n:    1000000,
			size: 1024,
		},
		{
			c:    500,
			n:    1000000,
			size: 1024,
		},
		{
			c:    1000,
			n:    1000000,
			size: 1024,
		},
		{
			c:    1000,
			n:    1000000,
			size: 2 * 1024,
		},
	}

	for _, sample := range samples {
		doPressureTest(proxy, sample.c, sample.n, sample.size)
	}
}
func doPressureTest(proxy *client.Proxy, c, n, size int) {
	wg = &sync.WaitGroup{}
	message = xrand.Letters(size)

	atomic.StoreInt64(&totalSent, 0)
	atomic.StoreInt64(&totalRecv, 0)

	wg.Add(n)

	chSeq := make(chan int32, n)

	// 创建连接
	for i := 0; i < c; i++ {
		conn, err := proxy.Dial()
		if err != nil {
			log.Errorf("gate connect failed: %v", err)
			return
		}

		go func(conn *client.Conn) {
			defer func() {
				_ = conn.Close()
			}()

			for {
				select {
				case seq, ok := <-chSeq:
					if !ok {
						return
					}

					err := conn.Push(&cluster.Message{
						Route: route.Greet,
						Seq:   seq,
						Data:  &pojo.GreetReq{Message: message},
					})
					if err != nil {
						log.Errorf("push message failed: %v", err)
						return
					}

					atomic.AddInt64(&totalSent, 1)
				}
			}
		}(conn)
	}

	startTime = xtime.Now().UnixNano()

	for i := 1; i <= n; i++ {
		chSeq <- int32(i)
	}

	wg.Wait()

	close(chSeq)

	totalTime := float64(time.Now().UnixNano()-startTime) / float64(time.Second)

	fmt.Printf("server               : %s\n", proxy.Client().Protocol())
	fmt.Printf("concurrency          : %d\n", c)
	fmt.Printf("latency              : %fs\n", totalTime)
	fmt.Printf("data size            : %s\n", convBytes(size))
	fmt.Printf("sent requests        : %d\n", totalSent)
	fmt.Printf("received requests    : %d\n", totalRecv)
	fmt.Printf("throughput (TPS)     : %d\n", int64(float64(totalRecv)/totalTime))
	fmt.Printf("--------------------------------\n")
}
func convBytes(bytes int) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes < KB:
		return fmt.Sprintf("%.2fB", float64(bytes))
	case bytes < MB:
		return fmt.Sprintf("%.2fKB", float64(bytes)/KB)
	case bytes < GB:
		return fmt.Sprintf("%.2fMB", float64(bytes)/MB)
	case bytes < TB:
		return fmt.Sprintf("%.2fGB", float64(bytes)/GB)
	default:
		return fmt.Sprintf("%.2fTB", float64(bytes)/TB)
	}
}
