package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"kart/global"
	"net/http"
)

// AddFileHandler 上传文件
func AddFileHandler(w http.ResponseWriter, r *http.Request) {
	fileObj, fileHeader, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	fmt.Println(fileHeader.Filename)
	fileID := global.StoreHandler.AddFile(fileObj, fileHeader.Filename)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success to Add File, fileID = %s.", fileID)
}

// GetFileHandler 获取文件
func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileID"]
	index := global.StoreHandler.FindByFileID(fileID)
	if index == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Failed to Get File By ID = %s.", fileID)
		return
	}
	w.WriteHeader(http.StatusOK)
	content := global.StoreHandler.Read(index.BlockID, index.Offset, index.Size)
	w.Write(content)
}
