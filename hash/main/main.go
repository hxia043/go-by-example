package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"hash/fnv"
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

	key := "3"
	var bucketMask uint32 = 1023
	h := fnv.New32a()
	h.Write([]byte(key))
	fmt.Println(h.Sum32())
	fmt.Println(h.Sum32() & bucketMask)
}
