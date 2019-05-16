package handlers

import (
	"Kart/database"
	"Kart/storage"
	"Kart/utils"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	// "golang.org/x/image/bmp"
	"Kart/global"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
)

// AddFileHandler 上传文件
func AddFileHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.FormValue("bucket")
	fileObj, fileHeader, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	//userID := r.Header.Get("userID")
	userID := "4220857c1585416391054f447f875e48"
	fmt.Println("userID = ", userID)
	fmt.Println("bucket = ", bucketName)
	object := storage.CreateObject(fileObj, bucketName)
	uf := storage.NewUserFileModel(object.ID.GetValue(), userID, fileHeader.Filename, int64(object.Size.GetInt()))
	global.DBConn.Insert("UserFileModel", uf)
	utils.JSONResponse(database.ModelToMap(object), w)

}

// ListFileHandler 获取文件列表
func ListFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		utils.JSONResponse(nil, w)
		return
	}
	userID := r.Header.Get("userID")
	vars := mux.Vars(r)
	bucketID := vars["bucketID"]
	// 先要检查这个 bucket 是不是属于该用户
	//isEnable := global.StoreHandler.CheckBucketPermission(userID, bucketID)
	//if !isEnable {
	//	w.WriteHeader(http.StatusOK)
	//	fmt.Fprintf(w, "You are not permitted to access this bucket.")
	//	return
	//}
	objectList := storage.ListObject(bucketID)
	var result []interface{}
	for _, object := range objectList {
		uf := storage.GetUserFile(userID, object.ID.GetValue())
		if uf != nil {
			result = append(result, uf)
		}
	}
	utils.JSONResponse(result, w)
}

// GetFileHandler 获取文件
func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID := vars["objectID"]
	params := r.URL.Query()
	width, _ := strconv.Atoi(params.Get("width"))
	height, _ := strconv.Atoi(params.Get("height"))
	//index := global.StoreHandler.FindByFileID(fileID)
	object := storage.GetObject(objectID)
	content := storage.GetObjectContent(object)
	//content := global.StoreHandler.Read(index.BlockID, index.Offset, index.Size)
	tp := http.DetectContentType(content)
	fmt.Println("tp = ", tp)
	if width > 0 || height > 0 {
		if tp == "image/jpeg" {
			img, err := jpeg.Decode(bytes.NewReader(content))
			if err != nil {
				panic(err)
			}
			newImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
			jpeg.Encode(w, newImg, nil)
		} else if tp == "image/png" {
			img, err := png.Decode(bytes.NewReader(content))
			if err != nil {
				panic(err)
			}
			newImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
			png.Encode(w, newImg)
		} else if tp == "image/gif" {
			img, err := gif.Decode(bytes.NewReader(content))
			if err != nil {
				panic(err)
			}
			newImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
			gif.Encode(w, newImg, nil)
		} else {
			w.Write(content)
		}
	} else {
		w.Write(content)
	}
}
