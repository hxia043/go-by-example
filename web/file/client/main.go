package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func uploadFile(uploadUrl, filename string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath := path.Join(currentDir, "../local/", filename)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(uploadUrl, contentType, bodyBuffer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
	return nil
}

func main() {
	uploadUrl := "http://localhost:9091/upload"
	filename := "local.html"
	err := uploadFile(uploadUrl, filename)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}
