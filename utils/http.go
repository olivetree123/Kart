package utils

import (
	"encoding/json"
	//"fmt"
	"net/http"
)

func PermitCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization")
}

// JSONResponse 返回 json 对象
func JSONResponse(data interface{}, w http.ResponseWriter) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	PermitCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonData)
	if err != nil {
		panic(err)
	}
}

// JSONParam 获取 json 参数
func JSONParam(r *http.Request) map[string]interface{} {
	params := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)
	return params
}
