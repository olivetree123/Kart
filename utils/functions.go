package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"reflect"
	"strconv"
	"strings"
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

func StringToSlice2(content string, length int) interface{} {
	fn := make([]byte, length)
	copy(fn[:], content)
	return fn
}

func StringToUUID(content string) [32]byte {
	var fn [32]byte
	copy(fn[:], content)
	return fn
}

func SliceToUUID(sl []byte) [32]byte {
	if len(sl) != 32 {
		panic("Invalid slice to uuid.")
	}
	var fn [32]byte
	copy(fn[:], sl)
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
	copy(idBytes[:], strings.Replace(uid.String(), "-", "", -1))
	return idBytes
}

func StructPtr2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.Elem().NumField(); i++ {
		fmt.Println(t.Elem().Field(i).Name)
		data[t.Elem().Field(i).Name] = v.Elem().Field(i).Interface()
	}
	return data
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func interface2String(inter interface{}) {
	switch inter.(type) {
	case string:
		fmt.Println("string", inter.(string))
		break
	case int:
		fmt.Println("int", inter.(int))
		break
	case float64:
		fmt.Println("float64", inter.(float64))
		break
	}
}

func GetLenFromTag(tag string) int {
	length := 0
	tagList := strings.Split(tag, ",")
	for _, t := range tagList {
		kv := strings.Split(t, "=")
		key := kv[0]
		value := kv[1]
		if key == "length" {
			length, _ = strconv.Atoi(value)
			break
		}
	}
	return length
}
