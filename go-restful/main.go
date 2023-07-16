package main

import (
	"io"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world\n")
}

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/hello").To(hello))
	restful.Add(ws)

	log.Fatal(http.ListenAndServe(":8083", nil))
}
