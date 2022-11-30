package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
)

func main() {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	id := base64.URLEncoding.EncodeToString(b)
	fmt.Println(id)

	urlId := url.QueryEscape(id)
	fmt.Println(urlId)

	id, _ = url.QueryUnescape(urlId)
	fmt.Println(id)
}
