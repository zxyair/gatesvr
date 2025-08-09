package proto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义测试用的protobuf消息
type testMessage struct {
	Value string
}

func (m *testMessage) ProtoMessage() {}

func TestCodec_Name(t *testing.T) {
	c := codec{}
	assert.Equal(t, "proto", c.Name())
}

func TestMarshalUnmarshal(t *testing.T) {
	// 准备测试数据
	msg := &testMessage{Value: "test value"}

	// 编码
	data, err := Marshal(msg)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// 解码
	var decoded testMessage
	err = Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "test value", decoded.Value)
}

func TestMarshal_Error(t *testing.T) {
	// 测试非protobuf类型
	_, err := Marshal("invalid type")
	assert.Error(t, err)
	assert.EqualError(t, err, "can't marshal a value that not implements proto.Buffer interface")
}

func TestUnmarshal_Error(t *testing.T) {
	// 测试非protobuf类型
	var target string
	err := Unmarshal([]byte("test data"), &target)
	assert.Error(t, err)
	assert.EqualError(t, err, "can't unmarshal to a value that not implements proto.Buffer")

	// 测试无效数据
	var msg testMessage
	err = Unmarshal([]byte("invalid data"), &msg)
	assert.Error(t, err)
}
