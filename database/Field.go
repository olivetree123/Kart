package database

import (
	"strconv"
)

type Field interface {
	GetName() string
	GetValue() string
	Bytes() []byte
	// GetType 和 GetLength 比较特殊，返回的是结构的属性，其他函数返回的都是对象的属性
	GetType() string
	GetLength() int
}

type StringField struct {
	Name   string
	Value  string
	Type   string
	Length int
}

type UUIDField struct {
	Name   string
	Value  string
	Type   string
	Length int
}

type BooleanField struct {
	Name   string
	Value  string
	Type   string
	Length int
}

type IntegerField struct {
	Name   string
	Value  string
	Type   string
	Length int
}

type DateTimeField struct {
	Name  string
	Value string
	Type  string
}

func NewStringField(name string, value string, length int) StringField {
	if len(value) > length {
		panic("length error.")
	}
	return StringField{
		Name:   name,
		Value:  value,
		Length: length,
	}
}

func NewUUIDField(name string, value string) UUIDField {
	if len(value) != 32 {
		panic("value error.")
	}
	return UUIDField{
		Name:   name,
		Value:  value,
		Length: 32,
	}
}

func NewBooleanField(name string, value bool) BooleanField {
	v := "1"
	if value == false {
		v = "0"
	}
	return BooleanField{
		Name:   name,
		Value:  v,
		Length: 1,
	}
}

func NewIntegerField(name string, value int64) IntegerField {
	return IntegerField{
		Name:   name,
		Value:  strconv.FormatInt(value, 10),
		Length: 20,
	}
}

func (field StringField) Bytes() []byte {
	data := make([]byte, field.Length)
	copy(data[:], field.Value)
	return data
}

func (field UUIDField) Bytes() []byte {
	data := make([]byte, field.Length)
	copy(data[:], field.Value)
	return data
}

func (field BooleanField) Bytes() []byte {
	data := make([]byte, field.Length)
	copy(data[:], field.Value)
	return data
}

func (field IntegerField) Bytes() []byte {
	data := make([]byte, field.Length)
	copy(data[:], field.Value)
	return data
}

func (field StringField) GetName() string {
	return field.Name
}

func (field UUIDField) GetName() string {
	return field.Name
}

func (field BooleanField) GetName() string {
	return field.Name
}

func (field IntegerField) GetName() string {
	return field.Name
}

func (field StringField) GetType() string {
	return "string"
}

func (field UUIDField) GetType() string {
	return "uuid"
}

func (field BooleanField) GetType() string {
	return "bool"
}

func (field IntegerField) GetType() string {
	return "int"
}

func (field StringField) GetLength() int {
	// StringField 默认为 255，用户可在 tag 中自定义
	return 255
}

func (field UUIDField) GetLength() int {
	return 32
}

func (field BooleanField) GetLength() int {
	return 1
}

func (field IntegerField) GetLength() int {
	return 20
}

func (field StringField) GetValue() string {
	return field.Value
}

func (field UUIDField) GetValue() string {
	return field.Value
}

func (field BooleanField) GetValue() string {
	return field.Value
}

func (field IntegerField) GetValue() string {
	return field.Value
}

func (field IntegerField) GetInt() int {
	r, err := strconv.Atoi(field.GetValue())
	if err != nil {
		panic(err)
	}
	return r
}
