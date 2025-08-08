package compress

type Compressor interface {
	// Name 名称
	Name() string
	// Compress 压缩
	Compress(data []byte) ([]byte, error)
	// Decompress 解压缩
	Decompress(data []byte) ([]byte, error)
}
