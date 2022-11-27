package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		timestamp := strconv.Itoa(time.Now().Nanosecond())
		hashWr := md5.New()
		hashWr.Write([]byte(timestamp))
		token := fmt.Sprintf("%x", hashWr.Sum(nil))

		t, _ := template.ParseFiles("login.go.tpl")
		log.Println(t.Execute(w, token))
	} else {
		r.ParseForm()

		if token := r.Form.Get("token"); token != "" {
			fmt.Println("Info: form token is:", token)
		} else {
			fmt.Println("Error: please get the form first!!")
			return
		}

		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])

		age, err := strconv.Atoi(r.Form.Get("age"))
		if err != nil {
			fmt.Println("err:", err.Error())
		}
		fmt.Println("age:", age)
	}
}

func main() {
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
