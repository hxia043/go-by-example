package main

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

func main() {
	timestamp := strconv.Itoa(time.Now().Nanosecond())
	hashMd5 := md5.New()
	hashMd5.Write([]byte(timestamp))

	fmt.Printf("%x", hashMd5.Sum(nil))
}
