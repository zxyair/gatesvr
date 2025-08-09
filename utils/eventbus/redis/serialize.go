package redis

import (
	"encoding/json"
	"gatesvr/core/value"
	"gatesvr/utils/eventbus"
	"gatesvr/utils/xconv"
	"gatesvr/utils/xtime"
	"gatesvr/utils/xuuid"
)

type data struct {
	ID        string `json:"id"`        // 事件ID
	Topic     string `json:"topic"`     // 事件主题
	Payload   string `json:"payload"`   // 事件载荷
	Timestamp int64  `json:"timestamp"` // 事件时间
}

// 序列化
func serialize(topic string, payload interface{}) ([]byte, error) {
	return json.Marshal(&data{
		ID:        xuuid.UUID(),
		Topic:     topic,
		Payload:   xconv.String(payload),
		Timestamp: xtime.Now().UnixNano(),
	})
}

// 反序列化
func deserialize(v []byte) (*eventbus.Event, error) {
	d := &data{}

	err := json.Unmarshal(v, d)
	if err != nil {
		return nil, err
	}

	return &eventbus.Event{
		ID:        d.ID,
		Topic:     d.Topic,
		Payload:   value.NewValue(d.Payload),
		Timestamp: xtime.UnixNano(d.Timestamp),
	}, nil
}
