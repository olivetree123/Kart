package main

import (
	"fmt"
	// "kart/storage"
	// "os"
	"github.com/gorilla/mux"
	"kart/handlers"
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
	// API 上传文件
	router.HandleFunc("/kart/file", handlers.AddFileHandler).Methods("post")
	// API 获取文件
	router.HandleFunc("/kart/file/{fileID}", handlers.GetFileHandler).Methods("get")
	host := "0.0.0.0"
	port := 8000
	fmt.Println("Listen ", host, " on ", port)
	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%d", host, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
