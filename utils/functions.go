package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
)

// ContentMd5 计算内容的 md5 值
func ContentMd5(content []byte) string {
	h := md5.New()
	_, err := h.Write(content)
	if err != nil {
		panic(err)
	}
	value := fmt.Sprintf("%x", h.Sum(nil))
	return value
}

func StringToSlice(content string, length int) []byte {
	fn := make([]byte, length)
	copy(fn[:], content)
	return fn
}

func StringToUUID(content string) [32]byte {
	var fn [32]byte
	copy(fn[:], content)
	return fn
}

func UUIDToString(uid [32]byte) string {
	return string(uid[:])
}

func SliceToString(sl []byte) string {
	return string(bytes.TrimRight(sl, "\x00"))
}

func NewUUID() [32]byte {
	uid := uuid.Must(uuid.NewRandom())
	var idBytes [32]byte
	copy(idBytes[:], uid.String())
	return idBytes
}
