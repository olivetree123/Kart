package storage

import (
	"Kart/database"
	"Kart/global"
	"Kart/utils"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"
)

type VolumeModel struct {
	ID       database.UUIDField
	DirPath  database.StringField `kart:"length=20"`
	MaxSize  database.IntegerField
	FreeSize database.IntegerField
	CreateAt database.IntegerField
}

// FreeSectionModel  组成了 freelist
type FreeSectionModel struct {
	ID       database.UUIDField
	VolumeID database.UUIDField
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
	VolumeID database.UUIDField
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
	ID         database.UUIDField
	ObjectID   database.UUIDField
	UserID     database.UUIDField
	ObjectName database.StringField `kart:"length=50"`
	Size       database.IntegerField
	CreateAt   database.IntegerField
}

func NewVolumeModel(dirPath string, maxSize int64, freeSize int64) VolumeModel {
	lenMap := database.GetFieldLenMap(VolumeModel{})
	return VolumeModel{
		ID:       database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		DirPath:  database.NewStringField("DirPath", dirPath, lenMap["DirPath"]),
		MaxSize:  database.NewIntegerField("MaxSize", maxSize),
		FreeSize: database.NewIntegerField("FreeSize", freeSize),
		CreateAt: database.NewIntegerField("CreateAt", time.Now().Unix()),
	}
}

func NewFreeSectionModel(volumeID string, offset int64, size int64) FreeSectionModel {
	return FreeSectionModel{
		ID:       database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		VolumeID: database.NewUUIDField("VolumeID", volumeID),
		Offset:   database.NewIntegerField("Offset", offset),
		Size:     database.NewIntegerField("Size", size),
		CreateAt: database.NewIntegerField("CreateAt", time.Now().Unix()),
	}
}

func NewObjectModel(volumeID string, bucketID string, offset int64, size int64) ObjectModel {
	return ObjectModel{
		ID:       database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		VolumeID: database.NewUUIDField("VolumeID", volumeID),
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

func NewUserFileModel(objectID string, userID string, fileName string, size int64) UserFileModel {
	lenMap := database.GetFieldLenMap(UserFileModel{})
	return UserFileModel{
		ID:         database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		ObjectID:   database.NewUUIDField("ObjectID", objectID),
		UserID:     database.NewUUIDField("UserID", userID),
		ObjectName: database.NewStringField("ObjectName", fileName, lenMap["ObjectName"]),
		Size:       database.NewIntegerField("Size", size),
		CreateAt:   database.NewIntegerField("CreateAt", time.Now().Unix()),
	}
}

func BucketModelFromMap(data map[string]string) BucketModel {
	lenMap := database.GetFieldLenMap(BucketModel{})
	createAt, _ := strconv.Atoi(data["CreatedAt"])
	public := true
	if data["Public"] == "0" {
		public = false
	}
	return BucketModel{
		ID:       database.NewUUIDField("ID", data["ID"]),
		UserID:   database.NewUUIDField("UserID", data["UserID"]),
		Name:     database.NewStringField("Name", data["Name"], lenMap["Name"]),
		Public:   database.NewBooleanField("Public", public),
		CreateAt: database.NewIntegerField("CreateAt", int64(createAt)),
	}
}

func FreeSectionModelFromMap(data map[string]string) FreeSectionModel {
	size, _ := strconv.Atoi(data["Size"])
	offset, _ := strconv.Atoi(data["Offset"])
	createAt, _ := strconv.Atoi(data["CreatedAt"])
	return FreeSectionModel{
		ID:       database.NewUUIDField("ID", data["ID"]),
		VolumeID: database.NewUUIDField("VolumeID", data["VolumeID"]),
		Offset:   database.NewIntegerField("Offset", int64(offset)),
		Size:     database.NewIntegerField("Size", int64(size)),
		CreateAt: database.NewIntegerField("CreateAt", int64(createAt)),
	}
}

func VolumeModelFromMap(data map[string]string) VolumeModel {
	lenMap := database.GetFieldLenMap(VolumeModel{})
	maxSize, _ := strconv.Atoi(data["MaxSize"])
	freeSize, _ := strconv.Atoi(data["FreeSize"])
	createAt, _ := strconv.Atoi(data["CreatedAt"])
	return VolumeModel{
		ID:       database.NewUUIDField("ID", data["ID"]),
		DirPath:  database.NewStringField("DirPath", data["DirPath"], lenMap["DirPath"]),
		MaxSize:  database.NewIntegerField("MaxSize", int64(maxSize)),
		FreeSize: database.NewIntegerField("FreeSize", int64(freeSize)),
		CreateAt: database.NewIntegerField("CreateAt", int64(createAt)),
	}
}

func ObjectModelFromMap(data map[string]string) ObjectModel {
	offset, _ := strconv.Atoi(data["Offset"])
	size, _ := strconv.Atoi(data["Size"])
	createAt, _ := strconv.Atoi(data["CreatedAt"])
	return ObjectModel{
		ID:       database.NewUUIDField("ID", data["ID"]),
		VolumeID: database.NewUUIDField("VolumeID", data["VolumeID"]),
		BucketID: database.NewUUIDField("BucketID", data["BucketID"]),
		Offset:   database.NewIntegerField("Offset", int64(offset)),
		Size:     database.NewIntegerField("Size", int64(size)),
		CreateAt: database.NewIntegerField("CreateAt", int64(createAt)),
	}
}

func UserFileModelFromMap(data map[string]string) UserFileModel {
	size, _ := strconv.Atoi(data["Size"])
	createAt, _ := strconv.Atoi(data["CreatedAt"])
	lenMap := database.GetFieldLenMap(UserFileModel{})
	return UserFileModel{
		ID:         database.NewUUIDField("ID", utils.UUIDToString(utils.NewUUID())),
		ObjectID:   database.NewUUIDField("ObjectID", data["ObjectID"]),
		UserID:     database.NewUUIDField("UserID", data["UserID"]),
		ObjectName: database.NewStringField("ObjectName", data["ObjectName"], lenMap["ObjectName"]),
		Size:       database.NewIntegerField("Size", int64(size)),
		CreateAt:   database.NewIntegerField("CreateAt", int64(createAt)),
	}
}

func CreateObject(f multipart.File, bucketName string) ObjectModel {
	// 1. 查找 bucket
	bucketCond := fmt.Sprintf("Name=%s", bucketName)
	bucketMap := global.DBConn.SelectOne("BucketModel", bucketCond)
	bucket := BucketModelFromMap(bucketMap)
	// 2. 找到剩余空间足够的 freeSection
	size, _ := f.Seek(0, io.SeekEnd)
	condition := fmt.Sprintf("Size>%d", size)
	sectionMap := global.DBConn.SelectOne("FreeSectionModel", condition)
	if sectionMap == nil {
		panic("剩余空间不足")
	}
	section := FreeSectionModelFromMap(sectionMap)

	cond := fmt.Sprintf("VolumeID=%s", section.VolumeID.GetValue())
	volumeMap := global.DBConn.SelectOne("VolumeModel", cond)
	volume := VolumeModelFromMap(volumeMap)
	filePath := path.Join(volume.DirPath.GetValue(), "binData.db")
	binFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	offset := section.Offset.GetInt()
	var content []byte
	_, err = f.Read(content)
	if err != nil {
		panic(err)
	}
	_, err = binFile.WriteAt(content, int64(offset))
	if err != nil {
		panic(err)
	}
	// 更新 SectionModel 表
	condition = fmt.Sprintf("ID=%s", section.ID.GetValue())
	data := make(map[string]string)
	beforeSize := section.Size.GetInt()
	afterSize := strconv.Itoa(beforeSize - int(size))
	afterOffset := strconv.Itoa(offset + int(size))
	data["Size"] = afterSize
	data["Offset"] = afterOffset
	global.DBConn.Update("SectionModel", condition, data)
	object := NewObjectModel(section.VolumeID.GetValue(), bucket.ID.GetValue(), int64(offset), size)
	global.DBConn.Insert("ObjectModel", object)
	return object
}

func GetObject(objectID string) ObjectModel {
	condition := fmt.Sprintf("ID=%s", objectID)
	objectMap := global.DBConn.SelectOne("ObjectModel", condition)
	object := ObjectModelFromMap(objectMap)
	return object
}

func ListObject(bucketID string) []ObjectModel {
	condition := fmt.Sprintf("BucketID=%s", bucketID)
	objectMapList := global.DBConn.Select("ObjectModel", condition)
	var objectList []ObjectModel
	for _, objectMap := range objectMapList {
		objectList = append(objectList, ObjectModelFromMap(objectMap))
	}
	return objectList
}

func GetObjectContent(object ObjectModel) []byte {
	volumeID := object.VolumeID.GetValue()
	cond := fmt.Sprintf("VolumeID=%s", volumeID)
	volumeMap := global.DBConn.SelectOne("VolumeModel", cond)
	volume := VolumeModelFromMap(volumeMap)
	filePath := path.Join(volume.DirPath.GetValue(), "binData.db")
	binFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	content := make([]byte, object.Size.GetInt())
	_, err = binFile.WriteAt(content, int64(object.Offset.GetInt()))
	if err != nil {
		panic(err)
	}
	return content
}

func GetUserFile(userID string, objectID string) *UserFileModel {
	condition := fmt.Sprintf("UserID=%s and ObjectID=%s", userID, objectID)
	ufMap := global.DBConn.SelectOne("UserFileModel", condition)
	if ufMap == nil {
		return nil
	}
	uf := UserFileModelFromMap(ufMap)
	return &uf
}
