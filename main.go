package main

import (
	"fmt"
	"kart/storage"
	"os"
)

func main() {
	st := storage.NewStorage()
	for i := 1; i < 4; i++ {
		fpath := fmt.Sprintf("/Users/gao/Downloads/cover/%d.png", i)
		fmt.Println("file path = ", fpath)
		f, err := os.Open(fpath)
		if err != nil {
			panic(err)
		}
		st.AddFile(f, "1.png")
	}
}
