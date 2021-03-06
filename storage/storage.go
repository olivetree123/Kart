package storage

import (
	"Kart/utils"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	// "syscall"
	// "github.com/spf13/viper"
	"Kart/config"
)

// Storage 存储
type Storage struct {
	DirPath string
	// BlockNum     int
	BlockList    []*Block
	FreeList     []*Section
	Indexes      *IndexDB
	Buckets      *BucketDB
	Users        *UserDB
	UserAndFiles *UserAndFileDB
	BlockMaxSize int
}

// NewStorage 创建存储
func NewStorage() *Storage {
	st := &Storage{
		DirPath: config.Config.GetString("FilePath"),
		// BlockNum:     4,
		BlockList:    nil,
		FreeList:     nil,
		Users:        NewUserDB(),
		Indexes:      NewIndexDB(),
		Buckets:      NewBucketDB(),
		UserAndFiles: NewUserAndFileDB(),
		BlockMaxSize: config.Config.GetInt("BlockMaxSize") * 1024 * 1024,
	}
	st.Init()
	return st
}

// Init 初始化
func (st *Storage) Init() {
	// maxSize := 100 * 1024 * 1024
	// maxSize := config.Config.GetInt("BlockMaxSize") * 1024 * 1024
	// for i := 0; i < st.BlockNum; i++ {
	// 	fileName := fmt.Sprintf("%d.db", i)
	// 	fpath := filepath.Join(st.DirPath, fileName)
	// 	st.AddBlock(i, fpath, maxSize)
	// }
	st.AddBlock()
}

// AddBlock 添加 Block
func (st *Storage) AddBlock() {
	blockID := len(st.BlockList)
	fileName := fmt.Sprintf("%d.db", blockID)
	filePath := filepath.Join(st.DirPath, fileName)
	// maxSize := config.Config.GetInt("BlockMaxSize") * 1024 * 1024
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	initOffset, err := f.Seek(0, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	block := &Block{
		ID:          blockID,
		FilePath:    filePath,
		MaxSize:     st.BlockMaxSize,
		FreeSize:    st.BlockMaxSize,
		FileHandler: f,
	}
	section := &Section{
		BlockID: blockID,
		Offset:  int(initOffset),
		Size:    st.BlockMaxSize - int(initOffset),
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
	fmt.Println("offset = ", offset)
	_, err := block.FileHandler.ReadAt(content, offset)
	if err != nil {
		panic(err)
	}
	return content
}

// Write 将内容写到文件中
func (st *Storage) Write(content []byte) (int, int, int) {
	length := len(content)
	if length > st.BlockMaxSize {
		panic("单个文件体积过大，无法存储")
	}
	i, section := st.FindSection(length)
	if i < 0 {
		st.AddBlock()
		i, section = st.FindSection(length)
		if i < 0 {
			panic("未知错误")
		}
	}
	fmt.Println("i = ", i)
	fmt.Println("blockID = ", section.BlockID)
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
func (st *Storage) AddFile(userID string, r io.Reader, fileName string, bucketName string) string {
	buf := bytes.NewBuffer([]byte{})
	_, err := buf.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	fmt.Println("bucketName = ", bucketName)
	bucket := st.Buckets.FindBucketByName(bucketName)
	if bucket == nil {
		return ""
	}
	fmt.Println("bucket = ", bucket)
	contentBytes := buf.Bytes()
	fileID := utils.ContentMd5(contentBytes)
	idx := st.Indexes.FindIndex(fileID)
	if idx != nil {
		fmt.Println("File already exists.")
		return utils.SliceToString(idx.ID[:])
	}
	blockID, offset, size := st.Write(contentBytes)
	fmt.Println(fileID)
	fmt.Println(blockID, offset, size)
	fmt.Println(string(bucket.ID[:]))
	index := NewIndex(fileID, string(bucket.ID[:]), blockID, offset, size)
	st.Indexes.AddIndex(index)
	uf := NewUserAndFile(userID, fileID, fileName, index.Size)
	st.UserAndFiles.Add(uf)
	return fileID
}

// AddBucket 添加 Bucket
func (st *Storage) AddBucket(userID string, name string, public bool) *Bucket {
	bucket := st.Buckets.AddBucket(userID, name, public)
	return bucket
}

// ListBucket Bucket列表
func (st *Storage) ListBucket(userID string) []*Bucket {
	buckets := st.Buckets.ListBucket(userID)
	return buckets
}

// ListByBucket 文件列表
func (st *Storage) ListByBucket(bucketID string) []*Index {
	files := st.Indexes.ListByBucket(bucketID)
	return files
}

// FindByFileID 查找文件
func (st *Storage) FindByFileID(fileID string) *Index {
	index := st.Indexes.FindIndex(fileID)
	if index == nil {
		return nil
	}
	return index
}

// AddUser 添加用户
func (st *Storage) AddUser(email string, password string) *User {
	user := st.Users.AddUser(email, password)
	return user
}

// VerifyUser 验证用户
func (st *Storage) VerifyUser(email string, password string) *User {
	user := st.Users.VerifyUser(email, password)
	return user
}

func (st *Storage) GetUserFileInfo(fileID string) *UserAndFile {
	uf := st.UserAndFiles.FindByFileID(fileID)
	return uf
}

func (st *Storage) CheckBucketPermission(userID string, bucketID string) bool {
	return st.Buckets.CheckPermission(userID, bucketID)
}
