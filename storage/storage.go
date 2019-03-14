package storage

import (
	"bytes"
	"fmt"
	"io"
	"kart/utils"
	"os"
	"path/filepath"
	// "syscall"
	"github.com/spf13/viper"
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
		DirPath:   viper.GetString("FilePath"),
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
			// fmt.Println("find, block id = ", section.BlockID, "free size = ", section.Size)
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

func (st *Storage) Read(blockID int64, offset int64, size int64) []byte {
	content := make([]byte, size)
	block := st.FindBlockByID(int(blockID))
	_, err := block.FileHandler.ReadAt(content, offset)
	if err != nil {
		panic(err)
	}
	return content
}

// Write 将内容写到文件中
func (st *Storage) Write(content []byte) (int, int, int) {
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
	return block.ID, section.Offset - length, length
}

// AddFile 添加文件
func (st *Storage) AddFile(r io.Reader, fileName string) string {
	buf := bytes.NewBuffer([]byte{})
	_, err := buf.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	contentBytes := buf.Bytes()
	fileID := utils.ContentMd5(contentBytes)
	idx := st.Indexes.FindIndex(fileID)
	if idx != nil {
		fmt.Println("File already exists.")
		return string(idx.FileID[:])
	}
	blockID, offset, size := st.Write(contentBytes)
	index := NewIndex(fileID, blockID, offset, size)
	st.Indexes.AddIndex(index)
	return fileID
}

// FindByFileID 查找文件
func (st *Storage) FindByFileID(fileID string) *Index {
	index := st.Indexes.FindIndex(fileID)
	if index == nil {
		return nil
	}
	return index
}
