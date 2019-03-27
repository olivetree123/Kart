package storage

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
)

// Bucket 文件夹
type Bucket struct {
	ID     [32]byte // 32 byte
	Name   [32]byte // max_length = 32 byte
	UserID [32]byte
	Public bool // is public or not
}

// BucketObject 文件夹
type BucketObject struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"userID"`
	Public bool   `json:"public"` // is public or not
}

// NewBucket 创建 Bucket
func NewBucket(userID [32]byte, name string, public bool) *Bucket {
	uid := uuid.Must(uuid.NewRandom())
	fmt.Println("bucket uid = ", uid)
	var nameBytes, idBytes [32]byte
	copy(nameBytes[:], name)
	copy(idBytes[:], uid.String())
	bucket := &Bucket{
		ID:     idBytes,
		Name:   nameBytes,
		UserID: userID,
		Public: public,
	}
	return bucket
}

// ToObject 转换为 BucketObject
func (bucket *Bucket) ToObject() *BucketObject {
	obj := &BucketObject{
		ID:     string(bytes.TrimRight(bucket.ID[:], "\x00")),
		Name:   string(bytes.TrimRight(bucket.Name[:], "\x00")),
		UserID: string(bytes.TrimRight(bucket.UserID[:], "\x00")),
		Public: bucket.Public,
	}
	return obj
}
