package main

import (
	"fmt"

	"github.com/bitly/go-simplejson"
)

func main() {
	js, _ := simplejson.NewJson([]byte(`{
		"test": {
			"array": [1, "2", 3],
			"int": 10,
			"float": 5.150,
			"bignum": 9223372036854775807,
			"string": "simplejson",
			"bool": true
		}
	}`))

	arr, _ := js.Get("test").Get("array").Array()
	fmt.Println(arr)

	i, _ := js.Get("test").Get("int").Int()
	fmt.Println(i)

	ms := js.Get("test").Get("string").MustString()
	fmt.Println(ms)
}
