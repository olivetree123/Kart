package storage

import (
	"kart/database"
	"kart/utils"
	"time"
)

type BlockModel struct {
	ID       database.UUIDField
	MaxSize  database.IntegerField
	FreeSize database.IntegerField
	CreateAt database.IntegerField
}

type SectionModel struct {
	ID       database.UUIDField
	BlockID  database.IntegerField
	Offset   database.IntegerField
	Size     database.IntegerField
	CreateAt database.IntegerField
}

type BucketModel struct {
	ID       database.UUIDField
	UserID   database.UUIDField
	Name     database.StringField `kart:"length=10"`
	Public   database.BooleanField
	CreateAt database.IntegerField
}

type ObjectModel struct {
	ID       database.UUIDField
	BlockID  database.IntegerField
	BucketID database.UUIDField
	Offset   database.IntegerField
	Size     database.IntegerField
	CreateAt database.IntegerField
}

type UserModel struct {
	ID       database.UUIDField
	NickName database.StringField `kart:"length=20"`
	Email    database.StringField `kart:"length=20"`
	PassWord database.StringField `kart:"length=40"`
	Avatar   database.StringField `kart:"length=100"`
	CreateAt database.IntegerField
}

type UserFileModel struct {
	ID       database.UUIDField
	FileID   database.UUIDField
	UserID   database.UUIDField
	FileName database.StringField `kart:"length=50"`
	Size     database.IntegerField
	CreateAt database.IntegerField
}

func NewBlockModel(maxSize int64, freeSize int64) BlockModel {
	return BlockModel{
		ID:       database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		MaxSize:  database.NewIntegerField("MaxSize", maxSize),
		FreeSize: database.NewIntegerField("FreeSize", freeSize),
		CreateAt: database.NewIntegerField("CreateAt", time.Now().Unix()),
	}
}

func NewSectionModel(blockID int64, offset int64, size int64) SectionModel {
	return SectionModel{
		ID:       database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		BlockID:  database.NewIntegerField("BlockID", blockID),
		Offset:   database.NewIntegerField("Offset", offset),
		Size:     database.NewIntegerField("Size", size),
		CreateAt: database.NewIntegerField("CreateAt", time.Now().Unix()),
	}
}

func NewObjectModel(blockID int64, bucketID string, offset int64, size int64) ObjectModel {
	return ObjectModel{
		ID:       database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		BlockID:  database.NewIntegerField("BlockID", blockID),
		BucketID: database.NewUUIDField("BucketID", bucketID),
		Offset:   database.NewIntegerField("Offset", offset),
		Size:     database.NewIntegerField("Size", size),
		CreateAt: database.NewIntegerField("CreateAt", time.Now().Unix()),
	}
}

func NewBucketModel(userID string, name string, public bool) BucketModel {
	lenMap := database.GetFieldLenMap(BucketModel{})
	return BucketModel{
		ID:       database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		UserID:   database.NewUUIDField("UserID", userID),
		Name:     database.NewStringField("Name", name, lenMap["Name"]),
		Public:   database.NewBooleanField("Public", public),
		CreateAt: database.NewIntegerField("CreateAt", time.Now().Unix()),
	}
}

func NewUserModel(nickName string, email string, password string, avatar string) UserModel {
	lenMap := database.GetFieldLenMap(UserModel{})
	return UserModel{
		ID:       database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		NickName: database.NewStringField("NickName", nickName, lenMap["NickName"]),
		Email:    database.NewStringField("Email", email, lenMap["Email"]),
		PassWord: database.NewStringField("PassWord", password, lenMap["PassWord"]),
		Avatar:   database.NewStringField("Avatar", avatar, lenMap["Avatar"]),
		CreateAt: database.NewIntegerField("CreateAt", time.Now().Unix()),
	}
}

func NewUserFileModel(fileID string, userID string, fileName string, size int64) UserFileModel {
	lenMap := database.GetFieldLenMap(UserFileModel{})
	return UserFileModel{
		ID:       database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		FileID:   database.NewUUIDField("FileID", fileID),
		UserID:   database.NewUUIDField("UserID", userID),
		FileName: database.NewStringField("FileName", fileName, lenMap["FileName"]),
		Size:     database.NewIntegerField("Size", size),
		CreateAt: database.NewIntegerField("CreateAt", time.Now().Unix()),
	}
}
