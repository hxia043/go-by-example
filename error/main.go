package main

import (
	"fmt"
	"net/http"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error
type testHandler func(name string) error

type fileErr struct {
	msg string
}

func (fe *fileErr) Error() string {
	return fe.msg
}

func testHello(name string) error {
	_, err := fmt.Println(name)
	return err
}

func sayHello(w http.ResponseWriter, r *http.Request) error {
	_, err := w.Write([]byte(fmt.Sprintln("hello, world")))
	if err != nil {
		return err
	}

	return nil
}

func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := ah(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func testError() error {
	var fErr *fileErr
	var flag bool = true
	if flag {
		return fErr
	}

	return fErr
}

func main() {
	if err := testError(); err != nil {
		fmt.Println(err)
	}

	if err := testHandler(testHello)("hxia"); err != nil {
		fmt.Println(err)
	}

	http.Handle("/hello", appHandler(sayHello))
	http.ListenAndServe(":9091", nil)
}
