package main

import (
	"fmt"
	"kart/storage"
	// "os"
)

func main() {
	st := storage.NewStorage()
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
	index := st.FindByFileID("1a6117d59aa13dd42c64d23e34ba4dcd")
	if index == nil {
		fmt.Println("index is nil")
	} else {
		fmt.Println("index found, ", index.BlockID, index.Offset, index.Size)
	}
}
