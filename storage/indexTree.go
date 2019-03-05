package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"unsafe"
)

// IndexTree 索引树🌲
type IndexTree struct {
	Filepath    string
	FileHandler *os.File
	IndexList   []*Index
}

// AddIndex xxx
func (tree *IndexTree) AddIndex(index *Index) {
	binary.Write(tree.FileHandler, binary.LittleEndian, index)
	tree.IndexList = append(tree.IndexList, index)
}

// FindIndex 查找索引
func (tree *IndexTree) FindIndex(fileID string) *Index {
	var fID [32]byte
	copy(fID[:], fileID)
	for _, index := range tree.IndexList {
		if index.FileID == fID {
			return index
		}
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
	fmt.Println("index file size = ", size)
	length := unsafe.Sizeof(Index{})
	fmt.Println("index struct size = ", length)
	var offset int64
	for size > 0 {
		indexBytes := make([]byte, length)
		n, err := tree.FileHandler.ReadAt(indexBytes, offset)
		if err != nil {
			panic(err)
		}
		fmt.Println("read index ", n, " bytes")
		buf := bytes.NewReader(indexBytes)
		var index Index
		err = binary.Read(buf, binary.LittleEndian, &index)
		if err != nil {
			panic(err)
		}
		fmt.Println("get index, fileID = ", string(index.FileID[:]))
		size -= int64(length)
		offset += int64(length)
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
		IndexList:   nil,
	}
	tree.LoadIndex()
	return tree
}
