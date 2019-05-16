package global

import (
	"Kart/database"
	//"Kart/storage"
	"time"
)

// StoreHandler xxx
//var StoreHandler = storage.NewStorage()
var DBConn = database.NewConnection("kart")

// DBConn 数据库连接
//var DBConn = database.NewConnection("kart")

// TokenMap token 存储
var TokenMap = make(map[string]*Cache)

// var Cache = make(map[string]interface{})

// Cache 缓存
type Cache struct {
	Key        string
	Value      interface{}
	CreateTime int64
	Duration   int
}

// NewCache 创建 cache
func NewCache(key string, value interface{}, duration int) *Cache {
	cache := &Cache{
		Key:        key,
		Value:      value,
		CreateTime: time.Now().Unix(),
		Duration:   duration,
	}
	return cache
}

var cache = &Cache{}

// SetToken 设置 token
func SetToken(key string, value map[string]string) {
	cache := NewCache(key, value, 3600)
	TokenMap[key] = cache
}

// GetToken 获取 token
//func GetToken(key string) *storage.UserObject {
//	if r, found := TokenMap[key]; found {
//		now := time.Now().Unix()
//		if r.Duration == 0 || r.CreateTime+int64(r.Duration) > now {
//			return r.Value.(*storage.UserObject)
//		}
//	}
//	return nil
//}
