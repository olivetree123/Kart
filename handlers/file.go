package handlers

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"kart/utils"
	// "golang.org/x/image/bmp"
	"image/gif"
	"image/jpeg"
	"image/png"
	"kart/global"
	"net/http"
	"strconv"
)

// AddFileHandler 上传文件
func AddFileHandler(w http.ResponseWriter, r *http.Request) {
	bucket := r.FormValue("bucket")
	fileObj, fileHeader, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	// _, err = jpeg.Decode(fileObj)
	// if err != nil {
	// 	fmt.Println("EEEEEEEEE")
	// 	panic(err)
	// }
	userID := r.Header.Get("userID")
	fmt.Println("userID = ", userID)
	fmt.Println("bucket = ", bucket)
	fmt.Println(fileHeader.Filename)
	fileID := global.StoreHandler.AddFile(userID, fileObj, fileHeader.Filename, bucket)
	w.WriteHeader(http.StatusOK)
	if fileID == "" {
		fmt.Fprintf(w, "Failed to Add File.")
		return
	}
	fmt.Fprintf(w, "Success to Add File, fileID = %s.", fileID)
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
	isEnable := global.StoreHandler.CheckBucketPermission(userID, bucketID)
	if !isEnable {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "You are not permitted to access this bucket.")
		return
	}
	files := global.StoreHandler.ListByBucket(bucketID)
	var result []interface{}
	for _, f := range files {
		fmt.Println("file id = ", utils.SliceToString(f.ID[:]))
		obj := global.StoreHandler.GetUserFileInfo(utils.SliceToString(f.ID[:]))
		if obj != nil {
			result = append(result, obj.ToObject())
		}
	}
	utils.JSONResponse(result, w)
}

// GetFileHandler 获取文件
func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileID"]
	params := r.URL.Query()
	width, _ := strconv.Atoi(params.Get("width"))
	height, _ := strconv.Atoi(params.Get("height"))
	index := global.StoreHandler.FindByFileID(fileID)
	if index == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Failed to Get File By ID = %s.", fileID)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(index.BlockID, index.Offset, index.Size)
	content := global.StoreHandler.Read(index.BlockID, index.Offset, index.Size)
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
