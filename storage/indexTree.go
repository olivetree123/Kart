package storage

import (
	"encoding/binary"
	"os"
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

// NewIndexTree 创建索引树
func NewIndexTree() *IndexTree {
	indexFilePath := "/data/kart/kart.idx"
	f, err := os.OpenFile(indexFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	return &IndexTree{
		Filepath:    indexFilePath,
		FileHandler: f,
		IndexList:   nil,
	}
}
