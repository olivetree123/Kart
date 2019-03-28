package storage

import (
	"bytes"
	"encoding/binary"
	"kart/config"
	"os"
	"path/filepath"
	"unsafe"
)

// BucketTree 文件夹🌲
type BucketTree struct {
	Filepath    string
	FileHandler *os.File
	BucketMap   map[string]*Bucket
}

// AddBucket 添加索引，并将索引写入文件
func (tree *BucketTree) AddBucket(userID string, name string, public bool) *Bucket {
	var userIDBytes [32]byte
	copy(userIDBytes[:], userID)
	bucket := NewBucket(userIDBytes, name, public)
	binary.Write(tree.FileHandler, binary.LittleEndian, bucket)
	tree.BucketMap[string(bucket.ID[:])] = bucket
	return bucket
}

// FindBucketByName 查找 bucket
func (tree *BucketTree) FindBucketByName(name string) *Bucket {
	if bucket, found := tree.BucketMap[name]; found {
		return bucket
	}
	return nil
}

// LoadBucket 加载已有的索引文件
func (tree *BucketTree) LoadBucket() {
	size, err := tree.FileHandler.Seek(0, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	if size == 0 {
		return
	}
	length := unsafe.Sizeof(Bucket{})
	var offset int64
	for size > 0 {
		indexBytes := make([]byte, length)
		_, err := tree.FileHandler.ReadAt(indexBytes, offset)
		if err != nil {
			panic(err)
		}
		buf := bytes.NewReader(indexBytes)
		var bucket Bucket
		err = binary.Read(buf, binary.LittleEndian, &bucket)
		if err != nil {
			panic(err)
		}
		size -= int64(length)
		offset += int64(length)
		tree.BucketMap[string(bucket.ID[:])] = &bucket
	}
}

// NewBucketTree 创建索引树
func NewBucketTree() *BucketTree {
	bucketFilePath := filepath.Join(config.Config.GetString("FilePath"), config.Config.GetString("BucketFileName"))
	f, err := os.OpenFile(bucketFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	tree := &BucketTree{
		Filepath:    bucketFilePath,
		FileHandler: f,
		BucketMap:   make(map[string]*Bucket),
	}
	tree.LoadBucket()
	return tree
}
