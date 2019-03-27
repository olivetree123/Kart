package storage

import (
	"bytes"
	"encoding/binary"
	"kart/config"
	"os"
	"path/filepath"
	"unsafe"
)

// UserTree 文件夹🌲
type UserTree struct {
	Filepath    string
	FileHandler *os.File
	UserMap     map[string]*User
}

// AddUser 添加用户
func (tree *UserTree) AddUser(email string, password string) *User {
	if _, found := tree.UserMap[email]; found {
		return nil
	}
	user := NewUser(email, password)
	binary.Write(tree.FileHandler, binary.LittleEndian, user)
	tree.UserMap[email] = user
	return user
}

// VerifyUser 验证用户
func (tree *UserTree) VerifyUser(email string, password string) *User {
	if user, found := tree.UserMap[email]; found {
		return user
	}
	return nil
}

// LoadUser 加载已有的索引文件
func (tree *UserTree) LoadUser() {
	size, err := tree.FileHandler.Seek(0, os.SEEK_END)
	if err != nil {
		panic(err)
	}
	if size == 0 {
		return
	}
	length := unsafe.Sizeof(User{})
	var offset int64
	for size > 0 {
		indexBytes := make([]byte, length)
		_, err := tree.FileHandler.ReadAt(indexBytes, offset)
		if err != nil {
			panic(err)
		}
		buf := bytes.NewReader(indexBytes)
		var user User
		err = binary.Read(buf, binary.LittleEndian, &user)
		if err != nil {
			panic(err)
		}
		size -= int64(length)
		offset += int64(length)
		tree.UserMap[string(bytes.TrimRight(user.Email[:], "\x00"))] = &user
	}
}

// NewUserTree 创建索引树
func NewUserTree() *UserTree {
	userFileName := filepath.Join(config.Config.GetString("FilePath"), config.Config.GetString("UserFileName"))
	f, err := os.OpenFile(userFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	tree := &UserTree{
		Filepath:    userFileName,
		FileHandler: f,
		UserMap:     make(map[string]*User),
	}
	tree.LoadUser()
	return tree
}
