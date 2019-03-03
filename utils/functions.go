package utils

import (
	"crypto/md5"
	"fmt"
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
