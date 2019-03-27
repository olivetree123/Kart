package handlers

import (
	"encoding/json"
	"fmt"
	// "golang.org/x/image/bmp"
	"github.com/google/uuid"
	"kart/global"
	"kart/utils"
	"net/http"
)

// AddUserHandler 创建用户
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	params := utils.JSONParam(r)
	fmt.Printf("%+v\n", params)
	user := global.StoreHandler.AddUser(params["email"].(string), params["passwd"].(string))
	w.WriteHeader(http.StatusOK)
	if user != nil {
		w.Write(user.ToBytes())
	} else {
		fmt.Fprintf(w, "Failed to Create User")
	}
}

// LoginHandler 用户登录
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	params := utils.JSONParam(r)
	fmt.Printf("%+v\n", params)
	user := global.StoreHandler.VerifyUser(params["email"].(string), params["passwd"].(string))
	w.WriteHeader(http.StatusOK)
	if user == nil {
		fmt.Fprintf(w, "Failed to login.")
		return
	}
	token := uuid.Must(uuid.NewRandom()).String()
	userMap := user.ToMap()
	userMap["token"] = token
	rt, err := json.Marshal(userMap)
	if err != nil {
		panic(err)
	}
	global.SetToken(token, user.ToObject())
	w.Write(rt)
}
