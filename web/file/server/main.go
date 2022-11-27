package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(8 << 20) // 8Mib

	file, header, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println("Error", err.Error())
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "%v", header.Header)
	f, err := os.OpenFile("../test/upload.html", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	respFile, err := os.Open("../test/upload.html")
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	body, err := io.ReadAll(respFile)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	if _, err = w.Write(body); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
}

func main() {
	http.HandleFunc("/upload", Upload)
	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}
