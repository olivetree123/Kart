package handlers

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
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
	token := r.Header.Get("Authorization")
	fmt.Println("token = ", token)
	bucket := r.Form.Get("bucket")
	fileObj, fileHeader, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	// _, err = jpeg.Decode(fileObj)
	// if err != nil {
	// 	fmt.Println("EEEEEEEEE")
	// 	panic(err)
	// }
	fmt.Println(fileHeader.Filename)
	fileID := global.StoreHandler.AddFile(fileObj, fileHeader.Filename, bucket)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success to Add File, fileID = %s.", fileID)
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
