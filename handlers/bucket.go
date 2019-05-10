package handlers

import (
	"fmt"
	"kart/database"
	"kart/global"
	"kart/storage"
	"kart/utils"
	"net/http"
)

// AddBucketHandler 添加文件夹
func AddBucketHandler(w http.ResponseWriter, r *http.Request) {
	//userID := r.Header.Get("userID")
	userID := "4220857c1585416391054f447f875e48"
	fmt.Println("userID = ", userID)
	params := utils.JSONParam(r)
	fmt.Printf("%+v\n", params)
	name := params["name"].(string)
	// public := params["public"].(bool)
	//bucket := global.StoreHandler.AddBucket(userID, name, true)
	//utils.JSONResponse(bucket.ToObject(), w)
	bucket := storage.NewBucketModel(userID, name, true)
	global.DBConn.Insert("BucketModel", bucket)
	utils.JSONResponse(database.ModelToMap(bucket), w)
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
	//userID := r.Header.Get("userID")
	userID := "4220857c1585416391054f447f875e48"
	fmt.Println("userID = ", userID)
	//buckets := global.StoreHandler.ListBucket(userID)
	buckets := global.DBConn.Select("BucketModel", "UserID=4220857c1585416391054f447f875e48")
	//var rs []interface{}
	//for _, bucket := range buckets {
	//	rs = append(rs, bucket.ToObject())
	//}
	//fmt.Println("rs = ", rs)
	//utils.JSONResponse(rs, w)
	utils.JSONResponse(buckets, w)
}
