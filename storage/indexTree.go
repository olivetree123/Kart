package storage

import (
	"bytes"
	"encoding/binary"
	"os"
	"unsafe"
)

// IndexTree 索引树🌲
type IndexTree struct {
	Filepath    string
	FileHandler *os.File
	IndexList   map[string]*Index
}

// AddIndex 添加索引，并将索引写入文件
func (tree *IndexTree) AddIndex(index *Index) {
	binary.Write(tree.FileHandler, binary.LittleEndian, index)
	tree.IndexList[string(index.FileID[:])] = index
}

// FindIndex 查找索引
func (tree *IndexTree) FindIndex(fileID string) *Index {
	if index, found := tree.IndexList[fileID]; found {
		return index
	}
	return nil
}

// LoadIndex 加载已有的索引文件
func (tree *IndexTree) LoadIndex() {
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
		tree.IndexList[string(index.FileID[:])] = &index
	}
}

// NewIndexTree 创建索引树
func NewIndexTree() *IndexTree {
	indexFilePath := "/data/kart/kart.idx"
	f, err := os.OpenFile(indexFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	tree := &IndexTree{
		Filepath:    indexFilePath,
		FileHandler: f,
		IndexList:   make(map[string]*Index),
	}
	tree.LoadIndex()
	return tree
}
