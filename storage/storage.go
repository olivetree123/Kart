package storage

import (
	"bytes"
	"fmt"
	"io"
	"kart/utils"
	"os"
	"path/filepath"
)

// Storage 存储
type Storage struct {
	DirPath   string
	BlockNum  int
	BlockList []*Block
	FreeList  []*Section
	Indexes   *IndexTree
}

// NewStorage 创建存储
func NewStorage() *Storage {
	st := &Storage{
		DirPath:   "/data/kart",
		BlockNum:  4,
		BlockList: nil,
		FreeList:  nil,
		Indexes:   NewIndexTree(),
	}
	st.Init()
	return st
}

// Init 初始化
func (st *Storage) Init() {
	maxSize := 100 * 1024 * 1024
	for i := 0; i < st.BlockNum; i++ {
		fileName := fmt.Sprintf("%d.db", i)
		fpath := filepath.Join(st.DirPath, fileName)
		st.AddBlock(i, fpath, maxSize)
	}
}

// AddBlock 添加 Block
func (st *Storage) AddBlock(id int, filePath string, maxSize int) {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	initOffset, err := f.Seek(0, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	block := &Block{
		ID:          id,
		FilePath:    filePath,
		MaxSize:     maxSize,
		FreeSize:    maxSize,
		FileHandler: f,
	}
	section := &Section{
		BlockID: id,
		Offset:  int(initOffset),
		Size:    maxSize - int(initOffset),
	}
	st.BlockList = append(st.BlockList, block)
	st.FreeList = append(st.FreeList, section)
}

// FindSection 从 FreeList 中查找大小合适的 Section
func (st *Storage) FindSection(length int) (int, *Section) {
	for i, section := range st.FreeList {
		if section.Size >= length {
			fmt.Println("find, block id = ", section.BlockID, "free size = ", section.Size)
			return i, section
		}
	}
	return -1, nil
}

// RemoveSection 从 FreeList 中删除 Section
func (st *Storage) RemoveSection(i int) {
	st.FreeList = append(st.FreeList[:i], st.FreeList[i+1:]...)
}

// FindBlockByID xxx
func (st *Storage) FindBlockByID(id int) *Block {
	for _, block := range st.BlockList {
		if block.ID == id {
			return block
		}
	}
	return nil
}

// Write 将内容写到文件中
func (st *Storage) Write(content []byte) (string, int, int, int) {
	length := len(content)
	i, section := st.FindSection(length)
	if i < 0 {
		panic("存储空间不足")
	}
	block := st.FindBlockByID(section.BlockID)
	_, err := block.FileHandler.WriteAt(content, int64(section.Offset))
	if err != nil {
		panic(err)
	}
	section.Offset += length
	section.Size -= length
	if section.Size == 0 {
		st.RemoveSection(i)
	}
	return utils.ContentMd5(content), block.ID, section.Offset - length, length
}

// AddFile 添加文件
func (st *Storage) AddFile(r io.Reader, fileName string) {
	buf := bytes.NewBuffer([]byte{})
	n, err := buf.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	fmt.Println("read ", n, "bytes")
	fileID, blockID, offset, size := st.Write(buf.Bytes())
	index := NewIndex(fileID, blockID, offset, size)
	st.Indexes.AddIndex(index)
}
