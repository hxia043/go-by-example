package main

import (
	"fmt"
	"net/http"
	"time"
)

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, world")
}

func main() {
	t1 := time.Now()
	t2 := time.Now().Add(time.Second * 50)

	fmt.Println(t2.Sub(t1).Seconds())

	ut1 := time.Now().Unix()
	time.Sleep(5 * time.Second)
	ut2 := time.Now().Unix()
	fmt.Println(ut1, ut2)

	http.HandleFunc("/test", test)
	http.ListenAndServe(":9091", nil)
}

func sayHello() {
	fmt.Println("Hello, world")
	time.AfterFunc(10*time.Second, func() { sayHello() })
}

func init() {
	go sayHello()
}
