package database

import (
	"kart/utils"
	"reflect"
)

type Model interface {
	Bytes() []byte
	Len() int
	Fields() []Field
	GetID() [32]byte
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
}

func (model BucketModel) Bytes() []byte {
	var r []byte
	r = append(r, model.ID.Bytes()...)
	r = append(r, model.UserID.Bytes()...)
	r = append(r, model.Name.Bytes()...)
	r = append(r, model.Public.Bytes()...)
	return r
}

func (model BucketModel) Len() int {
	length := model.ID.Length +
		model.UserID.Length +
		model.Name.Length +
		model.Public.Length
	return length
}

func (model BucketModel) GetID() [32]byte {
	return utils.StringToUUID(model.ID.Value)
}

func (model BucketModel) Fields() []Field {
	r := []Field{model.ID, model.UserID, model.Name, model.Public}
	return r
}

func NewBucketModel(id string, userID string, name string, public bool) *BucketModel {
	t := reflect.TypeOf(BucketModel{})
	lenMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type.Name() == "StringField" {
			length := utils.GetLenFromTag(t.Field(i).Tag.Get("orm"))
			if length <= 0 {
				panic("Invalid length.")
			}
			lenMap[t.Field(i).Name] = length
		}
	}
	return &BucketModel{
		ID:     NewUUIDField("ID", id),
		UserID: NewUUIDField("UserID", userID),
		Name:   NewStringField("Name", name, lenMap["Name"]),
		Public: NewBooleanField("Public", public),
	}
}

func ModelToMap(model Model) map[string]string {
	data := make(map[string]string)
	for _, field := range model.Fields() {
		data[field.GetName()] = field.GetValue()
	}
	return data
}
