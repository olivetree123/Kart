package storage

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
)

// User 用户
type User struct {
	ID       [32]byte
	NickName [32]byte
	Email    [40]byte
	Passwd   [40]byte
	Avatar   [100]byte
}

// UserObject 用户对象
type UserObject struct {
	ID       string `json:"id"`
	NickName string `json:"nickName"`
	Email    string `json:"email"`
	// Passwd   string `json:"password"`
	Avatar string `json:"avatar"`
	Token  string `json:"token"`
}

// NewUserObject 创建 UserObject
func NewUserObject(user *User) *UserObject {
	u := &UserObject{}
	u.ID = string(bytes.Trim(user.ID[:], "\x00"))
	u.NickName = string(bytes.Trim(user.NickName[:], "\x00"))
	u.Email = string(bytes.Trim(user.Email[:], "\x00"))
	u.Avatar = string(bytes.Trim(user.Avatar[:], "\x00"))
	return u
}

// ToObject 转换为 UserObject
func (user *User) ToObject() *UserObject {
	return NewUserObject(user)
}

// ToBytes 转换成 map 格式
func (user *User) ToBytes() []byte {
	u := NewUserObject(user)
	r, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return r
}

// ToMap 转换成 map 格式
func (user *User) ToMap() map[string]string {
	u := NewUserObject(user)
	r, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	m := make(map[string]string)
	decoder := json.NewDecoder(bytes.NewBuffer(r))
	decoder.Decode(&m)
	return m
}

// NewUser 创建用户
func NewUser(email string, password string) *User {
	uid := uuid.Must(uuid.NewRandom())
	var avatar [100]byte
	var userID, nickName [32]byte
	var emailBytes, passwdBytes [40]byte
	copy(userID[:], uid.String())
	copy(emailBytes[:], email)
	copy(passwdBytes[:], password)
	user := &User{
		ID:       userID,
		NickName: nickName,
		Email:    emailBytes,
		Passwd:   passwdBytes,
		Avatar:   avatar,
	}
	return user
}
