package main

import (
	"fmt"
	"log"

	"github.com/artinZareie/iori-sync/internal/filesystem"
)

func main() {
	files, err := filesystem.WalkAsList("./")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(files[0].Info.IsDir())
}
