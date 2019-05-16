package storage

import (
	"Kart/config"
	"Kart/utils"
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"unsafe"
)

// IndexDB 索引存储
type IndexDB struct {
	Filepath    string
	FileHandler *os.File
	IndexMap    map[string]*Index
}

// AddIndex 添加索引，并将索引写入文件
func (tree *IndexDB) AddIndex(index *Index) {
	err := binary.Write(tree.FileHandler, binary.LittleEndian, index)
	if err != nil {
		panic(err)
	}
	tree.IndexMap[utils.SliceToString(index.ID[:])] = index
}

// FindIndex 查找索引
func (tree *IndexDB) FindIndex(fileID string) *Index {
	if index, found := tree.IndexMap[fileID]; found {
		return index
	}
	return nil
}

// ListByBucket 根据 bucketID 获取文件列表
func (tree *IndexDB) ListByBucket(bucketID string) []*Index {
	var rs []*Index
	for _, value := range tree.IndexMap {
		if utils.SliceToString(value.BucketID[:]) == bucketID {
			rs = append(rs, value)
		}
	}
	return rs
}

// LoadIndex 加载已有的索引文件
func (tree *IndexDB) LoadIndex() {
	size, err := tree.FileHandler.Seek(0, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	if size == 0 {
		return
	}
	length := unsafe.Sizeof(Index{})
	var offset int64
	for size > 0 {
		indexBytes := make([]byte, length)
		_, err := tree.FileHandler.ReadAt(indexBytes, offset)
		if err != nil {
			panic(err)
		}
		buf := bytes.NewReader(indexBytes)
		var index Index
		err = binary.Read(buf, binary.LittleEndian, &index)
		if err != nil {
			panic(err)
		}
		size -= int64(length)
		offset += int64(length)
		tree.IndexMap[utils.SliceToString(index.ID[:])] = &index
	}
}

// NewIndexDB 创建索引树
func NewIndexDB() *IndexDB {
	// indexFilePath := "/data/kart/kart.idx"
	indexFilePath := filepath.Join(config.Config.GetString("FilePath"), config.Config.GetString("IndexFileName"))
	f, err := os.OpenFile(indexFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	tree := &IndexDB{
		Filepath:    indexFilePath,
		FileHandler: f,
		IndexMap:    make(map[string]*Index),
	}
	tree.LoadIndex()
	return tree
}
