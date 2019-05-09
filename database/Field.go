package database

import (
	"fmt"
	"strconv"
)

type Field interface {
	GetName() string
	GetType() string
	GetLength() int
	Bytes() []byte
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
		Type:   "string",
		Length: length,
	}
}

func NewUUIDField(name string, value string) UUIDField {
	fmt.Println("value = ", value)
	if len(value) != 32 {
		panic("value error.")
	}
	return UUIDField{
		Name:   name,
		Value:  value,
		Type:   "uuid",
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
		Type:   "bool",
		Length: 1,
	}
}

func NewIntegerField(name string, value int) IntegerField {
	return IntegerField{
		Name:   name,
		Value:  strconv.Itoa(value),
		Type:   "int",
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
	return field.Type
}

func (field UUIDField) GetType() string {
	return field.Type
}

func (field BooleanField) GetType() string {
	return field.Type
}

func (field IntegerField) GetType() string {
	return field.Type
}

func (field StringField) GetLength() int {
	return field.Length
}

func (field UUIDField) GetLength() int {
	return field.Length
}

func (field BooleanField) GetLength() int {
	return field.Length
}

func (field IntegerField) GetLength() int {
	return field.Length
}
