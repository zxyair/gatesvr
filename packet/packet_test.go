package packet_test

import (
	"bytes"
	"gatesvr/packet"
	"gatesvr/utils/xrand"
	"testing"
)

var packer = packet.NewPacker(
	packet.WithHeartbeatTime(true),
)

func TestDefaultPacker_PackMessage(t *testing.T) {
	data, err := packer.PackMessage(&packet.Message{
		Seq:    1,
		Route:  1,
		Buffer: []byte("hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(data)

	message, err := packer.UnpackMessage(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("seq: %d", message.Seq)
	t.Logf("route: %d", message.Route)
	t.Logf("buffer: %s", string(message.Buffer))
}

func TestPackHeartbeat(t *testing.T) {
	data, err := packer.PackHeartbeat()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(data)

	isHeartbeat, err := packer.CheckHeartbeat(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(isHeartbeat)
}

func TestMagicCheck(t *testing.T) {
	// 构造一个包含正确magic值的消息
	msg := &packet.Message{
		Seq:    1,
		Route:  1,
		Buffer: []byte("test magic"),
	}

	data, err := packer.PackMessage(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Packed message with magic: %x", data[len(data)-2:])

	// 验证解包时是否正确识别magic值
	unpackedMsg, err := packer.UnpackMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Unpacked message: Seq=%d, Route=%d, Buffer=%s", unpackedMsg.Seq, unpackedMsg.Route, string(unpackedMsg.Buffer))

	if unpackedMsg.Seq != msg.Seq || unpackedMsg.Route != msg.Route || string(unpackedMsg.Buffer) != string(msg.Buffer) {
		t.Fatal("unpacked message does not match original")
	}

	// 构造一个不包含正确magic值的消息（模拟错误情况）
	invalidData := make([]byte, len(data))
	copy(invalidData, data)
	invalidData[len(invalidData)-2] = 0x00 // 修改magic值
	invalidData[len(invalidData)-1] = 0x00
	t.Logf("Modified message with invalid magic: %x", invalidData[len(invalidData)-2:])

	// 验证解包时是否返回错误
	_, err = packer.UnpackMessage(invalidData)
	if err == nil {
		t.Fatal("expected error for invalid magic, got nil")
	}
	t.Logf("Expected error for invalid magic: %v", err)
}

func TestCriticalMessagePackUnpack(t *testing.T) {
	msg := &packet.Message{
		Seq:        1,
		Route:      1,
		Buffer:     []byte("critical data"),
		IsCritical: true,
	}

	data, err := packer.PackMessage(msg)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Packed byte array: %v", data)

	unpackedMsg, err := packer.UnpackMessage(data)
	if err != nil {
		t.Fatal(err)
	}

	if unpackedMsg.Seq != msg.Seq || unpackedMsg.Route != msg.Route || string(unpackedMsg.Buffer) != string(msg.Buffer) || unpackedMsg.IsCritical != msg.IsCritical {
		t.Fatal("unpacked message does not match original")
	}

	t.Logf("Original: Seq=%d, Route=%d, Buffer=%s, IsCritical=%v", msg.Seq, msg.Route, string(msg.Buffer), msg.IsCritical)

	// Control group: Simulate unpack failure by corrupting the data
	corruptedData := make([]byte, len(data))
	copy(corruptedData, data)

	// Control group for normal (non-critical) message
	normalMsg := &packet.Message{
		Seq:        2,
		Route:      2,
		Buffer:     []byte("normal data"),
		IsCritical: false,
	}

	normalData, err := packer.PackMessage(normalMsg)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Packed normal byte array: %v", normalData)

	unpackedNormalMsg, err := packer.UnpackMessage(normalData)
	if err != nil {
		t.Fatal(err)
	}

	if unpackedNormalMsg.Seq != normalMsg.Seq || unpackedNormalMsg.Route != normalMsg.Route || string(unpackedNormalMsg.Buffer) != string(normalMsg.Buffer) || unpackedNormalMsg.IsCritical != normalMsg.IsCritical {
		t.Fatal("unpacked normal message does not match original")
	}

	t.Logf("Original normal: Seq=%d, Route=%d, Buffer=%s, IsCritical=%v", normalMsg.Seq, normalMsg.Route, string(normalMsg.Buffer), normalMsg.IsCritical)
}

func BenchmarkDefaultPacker_PackMessage(b *testing.B) {
	buffer := []byte(xrand.Letters(1024))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := packet.PackMessage(&packet.Message{
			Seq:    1,
			Route:  1,
			Buffer: buffer,
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnpack(b *testing.B) {
	buf, err := packet.PackMessage(&packet.Message{
		Seq:    1,
		Route:  1,
		Buffer: []byte("hello world"),
	})
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, err := packet.UnpackMessage(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDefaultPacker_ReadMessage(b *testing.B) {
	buf, err := packer.PackMessage(&packet.Message{
		Seq:    1,
		Route:  1,
		Buffer: []byte(xrand.Letters(1024)),
	})
	if err != nil {
		b.Fatal(err)
	}

	reader := bytes.NewReader(buf)

	b.ResetTimer()
	b.SetBytes(int64(len(buf)))

	for i := 0; i < b.N; i++ {
		if _, err = packer.ReadMessage(reader); err != nil {
			b.Fatal(err)
		}

		reader.Reset(buf)
	}
}

func BenchmarkDefaultPacker_UnpackMessage(b *testing.B) {
	buf, err := packer.PackMessage(&packet.Message{
		Seq:    1,
		Route:  1,
		Buffer: []byte(xrand.Letters(1024)),
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.SetBytes(int64(len(buf)))

	for i := 0; i < b.N; i++ {
		_, err := packer.UnpackMessage(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}
