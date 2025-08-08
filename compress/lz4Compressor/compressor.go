package lz4Compressor

import (
	"github.com/pierrec/lz4/v4"
)

type LZ4Compressor struct{}

func (l *LZ4Compressor) Name() string {
	return "lz4"
}

// 创建压缩器
func NewCompressor() *LZ4Compressor {
	return &LZ4Compressor{}
}

func (l *LZ4Compressor) Compress(data []byte) ([]byte, error) {
	buf := make([]byte, lz4.CompressBlockBound(len(data)))
	n, err := lz4.CompressBlock(data, buf, nil)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (l *LZ4Compressor) Decompress(data []byte) ([]byte, error) {
	buf := make([]byte, len(data)*10) // 假设解压后的数据不超过原始数据的10倍
	n, err := lz4.UncompressBlock(data, buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}
