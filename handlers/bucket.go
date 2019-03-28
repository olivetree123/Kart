package handlers

import (
	"fmt"
	"kart/global"
	"kart/utils"
	"net/http"
)

// AddBucketHandler 上传文件
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
