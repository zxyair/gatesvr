package redis

import (
	"gatesvr/log"
	"gatesvr/utils/eventbus"
	"gatesvr/utils/task"
	"reflect"
	"sync"
)

type consumer struct {
	rw       sync.RWMutex
	handlers map[uintptr][]eventbus.EventHandler
}

// 添加处理器
func (c *consumer) addHandler(handler eventbus.EventHandler) int {
	pointer := reflect.ValueOf(handler).Pointer()

	c.rw.Lock()
	defer c.rw.Unlock()

	if _, ok := c.handlers[pointer]; !ok {
		c.handlers[pointer] = make([]eventbus.EventHandler, 0, 1)
	}

	c.handlers[pointer] = append(c.handlers[pointer], handler)

	return len(c.handlers[pointer])
}

// 移除处理器
func (c *consumer) remHandler(handler eventbus.EventHandler) int {
	pointer := reflect.ValueOf(handler).Pointer()

	c.rw.Lock()
	defer c.rw.Unlock()

	delete(c.handlers, pointer)

	return len(c.handlers)
}

// 分发数据
func (c *consumer) dispatch(data []byte) {
	event, err := deserialize(data)
	if err != nil {
		log.Error("invalid event data")
		return
	}

	c.rw.RLock()
	defer c.rw.RUnlock()

	for _, handlers := range c.handlers {
		for i := range handlers {
			handler := handlers[i]
			task.AddTask(func() { handler(event) })
		}
	}
}
