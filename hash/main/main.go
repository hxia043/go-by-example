package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/scrypt"
)

func main() {
	timestamp := strconv.Itoa(time.Now().Nanosecond())
	hashMd5 := md5.New()
	hashMd5.Write([]byte(timestamp))

	fmt.Printf("%x", hashMd5.Sum(nil))

	// refer: https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/09.5.md
	salt := []byte("#$%n*~+x")
	dk, err := scrypt.Key([]byte("system123"), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(base64.StdEncoding.EncodeToString(dk))
}
