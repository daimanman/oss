package main

import (
	"log"
	"oss/handler"
)

func main() {
	log.Printf("************************ Start a OSS Server ******************\n")
	handler.StartOssServer()
}
