package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse 返回 json 对象
func JSONResponse(data interface{}, w http.ResponseWriter) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// JSONParam 获取 json 参数
func JSONParam(r *http.Request) map[string]interface{} {
	params := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)
	return params
}
