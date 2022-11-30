package main

import (
	"fmt"
	"net/http"
)

func readCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("username")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cookie.Name, cookie.Value)
	fmt.Fprintln(w, cookie.Name+"="+cookie.Value)
}

func main() {
	http.HandleFunc("/cookie", readCookie)
	http.ListenAndServe(":9091", nil)
}
