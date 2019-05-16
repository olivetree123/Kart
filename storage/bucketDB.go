package storage

import (
	"Kart/config"
	"Kart/utils"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"unsafe"
)

// BucketDB 文件夹存储
type BucketDB struct {
	Filepath    string
	FileHandler *os.File
	BucketMap   map[string]*Bucket
}

// AddBucket 添加索引，并将索引写入文件
func (tree *BucketDB) AddBucket(userID string, name string, public bool) *Bucket {
	var userIDBytes [32]byte
	copy(userIDBytes[:], userID)
	bucket := NewBucket(userIDBytes, name, public)
	err := binary.Write(tree.FileHandler, binary.LittleEndian, bucket)
	if err != nil {
		panic(err)
	}
	tree.BucketMap[utils.SliceToString(bucket.ID[:])] = bucket
	return bucket
}

// ListBucket Bucket列表
func (tree *BucketDB) ListBucket(userID string) []*Bucket {
	var rs []*Bucket
	for _, bucket := range tree.BucketMap {
		if utils.SliceToString(bucket.UserID[:]) == userID {
			rs = append(rs, bucket)
		}
	}
	return rs
}

// FindBucketByName 查找 bucket
func (tree *BucketDB) FindBucketByName(name string) *Bucket {
	for _, bucket := range tree.BucketMap {
		if utils.SliceToString(bucket.Name[:]) == name {
			return bucket
		}
	}
	return nil
}

// CheckPermission 检查用户是否有访问该 bucket 的权限
func (tree *BucketDB) CheckPermission(userID string, bucketID string) bool {
	if bucket, found := tree.BucketMap[bucketID]; found {
		if utils.SliceToString(bucket.UserID[:]) == userID {
			return true
		}
	}
	return false
}

// LoadBucket 加载已有的索引文件
func (tree *BucketDB) LoadBucket() {
	size, err := tree.FileHandler.Seek(0, io.SeekEnd)
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

// NewBucketDB 创建索引树
func NewBucketDB() *BucketDB {
	bucketFilePath := filepath.Join(config.Config.GetString("FilePath"), config.Config.GetString("BucketFileName"))
	f, err := os.OpenFile(bucketFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	tree := &BucketDB{
		Filepath:    bucketFilePath,
		FileHandler: f,
		BucketMap:   make(map[string]*Bucket),
	}
	tree.LoadBucket()
	return tree
}
