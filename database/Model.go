package database

import (
	"errors"
	"kart/utils"
	"reflect"
)

type Model interface {
}

type BaseModel struct {
	ID        UUIDField
	CreatedAt DateTimeField
	DeletedAt DateTimeField
}

type BucketModel struct {
	ID     UUIDField
	UserID UUIDField
	Name   StringField `orm:"length=10"`
	Public BooleanField
	Age    IntegerField
}

func GetModelID(model interface{}) (string, error) {
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)
	for i := 0; i < t.NumField(); i++ {
		if v.Field(i).Interface().(Field).GetName() == "ID" {
			return v.Field(i).Interface().(Field).GetValue(), nil
		}
	}
	return "", errors.New("ID not found")
}

func GetModelFields(model interface{}) []Field {
	var r []Field
	t := reflect.TypeOf(model)
	v := reflect.ValueOf(model)
	for i := 0; i < t.NumField(); i++ {
		r = append(r, v.Field(i).Interface().(Field))
	}
	//r := []Field{model.ID, model.UserID, model.Name, model.Public}
	return r
}

func GetFieldLenMap(model interface{}) map[string]int {
	t := reflect.TypeOf(model)
	lenMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Name() == "StringField" {
			length := utils.GetLenFromTag(t.Field(i).Tag.Get("kart"))
			if length <= 0 {
				panic("Invalid length.")
			}
			lenMap[t.Field(i).Name] = length
		}
	}
	return lenMap
}

func ModelToMap(model interface{}) map[string]string {
	data := make(map[string]string)
	for _, field := range GetModelFields(model) {
		data[field.GetName()] = field.GetValue()
	}
	return data
}

func NewBucketModel(id string, userID string, name string, public bool, age int) BucketModel {
	lenMap := GetFieldLenMap(BucketModel{})
	return BucketModel{
		ID:     NewUUIDField("ID", id),
		UserID: NewUUIDField("UserID", userID),
		Name:   NewStringField("Name", name, lenMap["Name"]),
		Public: NewBooleanField("Public", public),
		Age:    NewIntegerField("Age", int64(age)),
	}
}
