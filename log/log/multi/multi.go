package main

import (
	"io"
	"log"
	"os"
)

func main() {
	file, _ := os.OpenFile("demo.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.Default().Println("demo for log")
}
