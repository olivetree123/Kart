package storage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"kart/config"
	"kart/utils"
	"os"
	"path/filepath"
	"unsafe"
)

type UserAndFileDB struct {
	Filepath       string
	FileHandler    *os.File
	UserAndFileMap map[string]*UserAndFile
}

func (db *UserAndFileDB) Add(uf *UserAndFile) {
	err := binary.Write(db.FileHandler, binary.LittleEndian, uf)
	if err != nil {
		panic(err)
	}
	db.UserAndFileMap[utils.SliceToString(uf.ID[:])] = uf
}

func (db *UserAndFileDB) FindByFileID(fileID string) *UserAndFile {
	for _, val := range db.UserAndFileMap {
		if utils.UUIDToString(val.FileID) == fileID {
			return val
		}
	}
	return nil
}

func (db *UserAndFileDB) ListByUser(userID string) []*UserAndFile {
	var result []*UserAndFile
	for _, val := range db.UserAndFileMap {
		if utils.UUIDToString(val.UserID) == userID {
			result = append(result, val)
		}
	}
	return result
}

// LoadBucket 加载已有的索引文件
func (db *UserAndFileDB) LoadUserAndFile() {
	size, err := db.FileHandler.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err)
	}
	if size == 0 {
		return
	}
	length := unsafe.Sizeof(UserAndFile{})
	fmt.Println("length = ", length)
	var offset int64
	for size > 0 {
		indexBytes := make([]byte, length)
		_, err := db.FileHandler.ReadAt(indexBytes, offset)
		if err != nil {
			panic(err)
		}
		buf := bytes.NewReader(indexBytes)
		var uf UserAndFile
		err = binary.Read(buf, binary.LittleEndian, &uf)
		if err != nil {
			panic(err)
		}
		size -= int64(length)
		offset += int64(length)
		db.UserAndFileMap[utils.SliceToString(uf.ID[:])] = &uf
	}
}

// NewBucketDB 创建索引树
func NewUserAndFileDB() *UserAndFileDB {
	bucketFilePath := filepath.Join(config.Config.GetString("FilePath"), config.Config.GetString("UserAndFileName"))
	f, err := os.OpenFile(bucketFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	tree := &UserAndFileDB{
		Filepath:       bucketFilePath,
		FileHandler:    f,
		UserAndFileMap: make(map[string]*UserAndFile),
	}
	tree.LoadUserAndFile()
	return tree
}
