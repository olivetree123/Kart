package storage

import (
	"kart/utils"
	"time"
)

// UserAndFile 用户与文件的关联。不同的用户上传同一份文件时，文件只会保存一份，但是关联关系需要多份。
type UserAndFile struct {
	ID         [32]byte
	FileID     [32]byte // 32 byte
	UserID     [32]byte
	FileName   [64]byte
	UploadTime int64
	UpdateTime int64
}

type UserAndFileObject struct {
	ID         string
	FileID     string
	UserID     string
	FileName   string
	UploadTime int64
	UpdateTime int64
}

// NewUserAndFile 当用户上传文件时，需要将用户和文件进行关联。
func NewUserAndFile(userID string, fileID string, fileName string) *UserAndFile {
	now := time.Now().Unix()
	var fn [64]byte
	copy(fn[:], fileName)
	obj := &UserAndFile{
		ID:         utils.NewUUID(),
		FileID:     utils.StringToUUID(fileID),
		UserID:     utils.StringToUUID(userID),
		FileName:   fn,
		UploadTime: now,
		UpdateTime: now,
	}
	return obj
}

func (uf *UserAndFile) ToObject() *UserAndFileObject {
	obj := &UserAndFileObject{
		ID:         utils.SliceToString(uf.ID[:]),
		FileID:     utils.SliceToString(uf.FileID[:]),
		UserID:     utils.SliceToString(uf.UserID[:]),
		FileName:   utils.SliceToString(uf.FileName[:]),
		UploadTime: uf.UploadTime,
		UpdateTime: uf.UpdateTime,
	}
	return obj
}
