package handlers

import (
	"fmt"
	"kart/database"
	"kart/storage"
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
	//user := global.StoreHandler.AddUser(params["email"].(string), params["password"].(string))
	user := storage.NewUserModel("", params["email"].(string), params["password"].(string), "")
	global.DBConn.Insert("UserModel", user)
	utils.JSONResponse(database.ModelToMap(user), w)
}

// LoginHandler 用户登录
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "OPTIONS" {
		utils.JSONResponse(nil, w)
		return
	}
	params := utils.JSONParam(r)
	fmt.Printf("%+v\n", params)
	//user := global.StoreHandler.VerifyUser(params["email"].(string), params["password"].(string))
	condition := fmt.Sprintf("Email=%s and PassWord=%s", params["email"].(string), params["password"].(string))
	rs := global.DBConn.Select("UserModel", condition)
	if len(rs) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Failed to login.")
		return
	}
	user := rs[0]
	token := uuid.Must(uuid.NewRandom()).String()
	//userMap := user.ToMap()
	//userMap["token"] = token
	//userInfo := user.ToObject()
	//userInfo.Token = token
	//rt, err := json.Marshal(userMap)
	//if err != nil {
	//	panic(err)
	//}
	global.SetToken(token, user)
	user["token"] = token
	//w.Write(rt)
	utils.JSONResponse(user, w)
}
