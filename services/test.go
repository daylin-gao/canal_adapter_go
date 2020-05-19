package main

import "github.com/gao111/canal-adapter-go/models"

func main() {
	sync := models.NewSync("test" , "t1")
	sync.UpdateSync()
}