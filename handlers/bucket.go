package handlers

import (
	"fmt"
	// "golang.org/x/image/bmp"
	"kart/global"
	"kart/utils"
	"net/http"
)

// AddBucketHandler 上传文件
func AddBucketHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	fmt.Println("token = ", token)
	params := utils.JSONParam(r)
	fmt.Printf("%+v\n", params)
	name := params["name"].(string)
	userID := params["userID"].(string)
	public := params["public"].(bool)
	fmt.Println("public = ", public)
	bucket := global.StoreHandler.AddBucket(userID, name, true)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, "Success to Add File, fileID = %s.", fileID)
	utils.JSONResponse(bucket.ToObject(), w)
}
