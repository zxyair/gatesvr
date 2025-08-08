package lz4Compressor

import (
	"bytes"
	"fmt"
	"gatesvr/compress/lz4/pb"
	"gatesvr/encoding/json"
	"gatesvr/encoding/proto"
	"testing"
)

func TestLZ4Compressor_Name(t *testing.T) {
	compressor := &LZ4Compressor{}
	name := compressor.Name()
	if name != "lz4" {
		t.Errorf("expected name 'lz4', got '%s'", name)
	}
}

type greetReq struct {
	Message string
}

func TestLZ4Compressor_JSONCompression(t *testing.T) {
	compressor := &LZ4Compressor{}
	// Generate larger and more repetitive data
	data := greetReq{
		Message: "This is a much longer test message to demonstrate compression effects. " +
			"It contains repeated patterns like 'abcabcabc' and longer strings to show better compression ratios. " +
			"Repeated data compresses better because of the algorithm's ability to find and eliminate redundancy.",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	fmt.Printf("JSON data: len=%d\n", len(jsonData))

	compressed, err := compressor.Compress(jsonData)
	if err != nil {
		t.Fatalf("compression failed: %v", err)
	}
	fmt.Printf("JSON compressed data: len=%d, compression ratio=%.2f\n", len(compressed), float64(len(compressed))/float64(len(jsonData)))

	decompressed, err := compressor.Decompress(compressed)
	fmt.Printf("JSON decompressed data: len=%d\n", len(decompressed))
	if err != nil {
		t.Fatalf("decompression failed: %v", err)
	}
	if !bytes.Equal(decompressed, jsonData) {
		t.Fatal("decompressed JSON data does not match original")
	}
}

func TestLZ4Compressor_ProtoCompression(t *testing.T) {
	compressor := &LZ4Compressor{}
	protoData := &pb.HelloArgs{
		Name: "This is a much longer test message to demonstrate compression effects. " +
			"It contains repeated patterns like 'abcabcabc' and longer strings to show better compression ratios. " +
			"Repeated data compresses better because of the algorithm's ability to find and eliminate redundancy.",
	}
	protoBytes, err := proto.Marshal(protoData)
	if err != nil {
		t.Fatalf("Proto marshal failed: %v", err)
	}
	fmt.Printf("Proto data: len=%d\n", len(protoBytes))

	compressed, err := compressor.Compress(protoBytes)
	if err != nil {
		t.Fatalf("compression failed: %v", err)
	}
	fmt.Printf("Proto compressed data: len=%d, compression ratio=%.2f\n", len(compressed), float64(len(compressed))/float64(len(protoBytes)))

	decompressed, err := compressor.Decompress(compressed)
	fmt.Printf("Proto decompressed data: len=%d\n", len(decompressed))
	if err != nil {
		t.Fatalf("decompression failed: %v", err)
	}
	if !bytes.Equal(decompressed, protoBytes) {
		t.Fatal("decompressed Proto data does not match original")
	}
}

func BenchmarkLZ4Compressor_ProtoCompress(b *testing.B) {
	compressor := &LZ4Compressor{}
	protoData := &pb.HelloArgs{
		Name: "This is a much longer test message to demonstrate compression effects. " +
			"It contains repeated patterns like 'abcabcabc' and longer strings to show better compression ratios. " +
			"Repeated data compresses better because of the algorithm's ability to find and eliminate redundancy.",
	}
	data, _ := proto.Marshal(protoData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressor.Compress(data)
	}
}

func BenchmarkLZ4Compressor_ProtoDecompress(b *testing.B) {
	compressor := &LZ4Compressor{}
	protoData := &pb.HelloArgs{
		Name: "This is a much longer test message to demonstrate compression effects. " +
			"It contains repeated patterns like 'abcabcabc' and longer strings to show better compression ratios. " +
			"Repeated data compresses better because of the algorithm's ability to find and eliminate redundancy.",
	}
	data, _ := proto.Marshal(protoData)
	compressed, _ := compressor.Compress(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressor.Decompress(compressed)
	}
}

func BenchmarkLZ4Compressor_JSONCompress(b *testing.B) {
	compressor := &LZ4Compressor{}
	data := greetReq{
		Message: "This is a much longer test message to demonstrate compression effects. " +
			"It contains repeated patterns like 'abcabcabc' and longer strings to show better compression ratios. " +
			"Repeated data compresses better because of the algorithm's ability to find and eliminate redundancy.",
	}
	jsonData, _ := json.Marshal(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressor.Compress(jsonData)
	}
}

func BenchmarkLZ4Compressor_JSONDecompress(b *testing.B) {
	compressor := &LZ4Compressor{}
	data := greetReq{
		Message: "This is a much longer test message to demonstrate compression effects. " +
			"It contains repeated patterns like 'abcabcabc' and longer strings to show better compression ratios. " +
			"Repeated data compresses better because of the algorithm's ability to find and eliminate redundancy.",
	}
	jsonData, _ := json.Marshal(data)
	compressed, _ := compressor.Compress(jsonData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = compressor.Decompress(compressed)
	}
}
