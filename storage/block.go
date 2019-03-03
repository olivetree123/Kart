package storage

import (
	"os"
)

// Block 存储块
type Block struct {
	ID          int
	FilePath    string
	MaxSize     int
	FreeSize    int
	FileHandler *os.File
}
