package xfile

import (
	"gatesvr/utils/xpath"
	"io/fs"
	"os"
	"path/filepath"
)

// WriteFile 写文件
func WriteFile(file string, data []byte) error {
	path := filepath.Dir(file)

	if !xpath.IsDir(path) {
		err := os.MkdirAll(path, fs.ModePerm)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(file, data, fs.ModePerm)
}
