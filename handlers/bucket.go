package handlers

import (
	"fmt"
	"kart/global"
	"kart/utils"
	"net/http"
)

// AddBucketHandler 添加文件夹
func AddBucketHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userID")
	fmt.Println("userID = ", userID)
	params := utils.JSONParam(r)
	fmt.Printf("%+v\n", params)
	name := params["name"].(string)
	// public := params["public"].(bool)
	bucket := global.StoreHandler.AddBucket(userID, name, true)
	utils.JSONResponse(bucket.ToObject(), w)
}

// ListBucketHandler 文件夹列表
func ListBucketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		//w.WriteHeader(http.StatusOK)
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//w.Write(nil)
		utils.JSONResponse(nil, w)
		return
	}
	userID := r.Header.Get("userID")
	fmt.Println("userID = ", userID)
	buckets := global.StoreHandler.ListBucket(userID)
	var rs []interface{}
	for _, bucket := range buckets {
		rs = append(rs, bucket.ToObject())
	}
	fmt.Println("rs = ", rs)
	utils.JSONResponse(rs, w)
}
