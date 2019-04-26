package main

import (
	"fmt"
	// "kart/storage"
	// "os"
	"github.com/gorilla/mux"
	// "github.com/spf13/viper"
	"kart/config"
	"kart/handlers"
	"kart/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	// st := storage.NewStorage()
	// for i := 1; i < 4; i++ {
	// 	fpath := fmt.Sprintf("/Users/gao/Downloads/cover/%d.png", i)
	// 	fmt.Println("file path = ", fpath)
	// 	f, err := os.Open(fpath)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	etag := st.AddFile(f, "1.png")
	// 	fmt.Println("etag = ", etag)
	// }
	// index := st.FindByFileID("1a6117d59aa13dd42c64d23e34ba4dcd")
	// if index == nil {
	// 	fmt.Println("index is nil")
	// } else {
	// 	fmt.Println("index found, ", index.BlockID, index.Offset, index.Size)
	// }
	router := mux.NewRouter()
	// API 创建用户
	router.HandleFunc("/kart/user", handlers.AddUserHandler).Methods("post")
	// API 登陆
	router.HandleFunc("/kart/login", handlers.LoginHandler).Methods("post")
	// API 创建 Bucket
	// router.HandleFunc("/kart/bucket", handlers.AddBucketHandler).Methods("post")
	router.Handle(
		"/kart/bucket",
		middleware.Auth(http.HandlerFunc(handlers.AddBucketHandler)),
	).Methods("post")
	// API 获取 Bucket 列表
	router.Handle(
		"/kart/buckets",
		middleware.Auth(http.HandlerFunc(handlers.ListBucketHandler)),
	).Methods("get")
	//router.HandleFunc("/kart/buckets", handlers.ListBucketHandler).Methods("get")
	// API 上传文件
	router.Handle(
		"/kart/file",
		middleware.Auth(http.HandlerFunc(handlers.AddFileHandler)),
	).Methods("post")
	// API 获取文件列表
	router.Handle(
		"/kart/files/{bucketID}",
		middleware.Auth(http.HandlerFunc(handlers.ListFileHandler)),
	).Methods("get")
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
