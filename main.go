package main

import (
	"fmt"
	"kart/global"
	//"kart/utils"
	// "kart/storage"
	// "os"
	"github.com/gorilla/mux"
	// "github.com/spf13/viper"
	"kart/config"
	"kart/handlers"
	"kart/middleware"
	"kart/storage"
	"log"
	"net/http"
	"time"
)

func main() {
	//conn := database.NewConnection("kart")
	//fmt.Println("begin to create table")
	//conn.CreateTable(database.BucketModel{})
	//data := database.NewBucketModel(
	//	utils.UUIDToString(utils.NewUUID()),
	//	"5a81e3552e524e319e196676a91193b9",
	//	"gaojian",
	//	true,
	//	20,
	//)
	//fmt.Println("data = ", data)
	//conn.Insert("BucketModel", data)
	//result := conn.Select("BucketModel", "Age=20")
	//fmt.Println("result = ", result)
	global.DBConn.CreateTable(storage.BucketModel{})
	global.DBConn.CreateTable(storage.BlockModel{})
	global.DBConn.CreateTable(storage.ObjectModel{})
	global.DBConn.CreateTable(storage.SectionModel{})
	global.DBConn.CreateTable(storage.UserModel{})
	global.DBConn.CreateTable(storage.UserFileModel{})

	router := mux.NewRouter()
	// API 创建用户
	router.HandleFunc("/kart/user", handlers.AddUserHandler).Methods("post")
	// API 登陆
	router.HandleFunc("/kart/login", handlers.LoginHandler).Methods("post", "options")
	// API 创建 Bucket
	router.HandleFunc("/kart/bucket", handlers.AddBucketHandler).Methods("post")
	//router.Handle(
	//	"/kart/bucket",
	//	middleware.Auth(http.HandlerFunc(handlers.AddBucketHandler)),
	//).Methods("post")
	// API 获取 Bucket 列表
	//router.Handle(
	//	"/kart/buckets",
	//	middleware.Auth(http.HandlerFunc(handlers.ListBucketHandler)),
	//).Methods("get", "options")
	router.HandleFunc("/kart/buckets", handlers.ListBucketHandler).Methods("get")
	// API 上传文件
	router.Handle(
		"/kart/file",
		middleware.Auth(http.HandlerFunc(handlers.AddFileHandler)),
	).Methods("post")
	// API 获取文件列表
	router.Handle(
		"/kart/files/{bucketID}",
		middleware.Auth(http.HandlerFunc(handlers.ListFileHandler)),
	).Methods("get", "options")
	//router.HandleFunc("/kart/files", handlers.ListFileHandler).Methods("get")
	// API 获取文件
	router.HandleFunc("/kart/file/{fileID}", handlers.GetFileHandler).Methods("get")
	// host := "0.0.0.0"
	// port := 8000
	host := config.Config.GetString("Host")
	port := config.Config.GetInt("Port")
	fmt.Println("Listen ", host, " on ", port)
	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%d", host, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
