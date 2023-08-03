package main

import (
	"io"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"k8s.io/apiserver/pkg/server"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
)

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "world\n")
}

func main() {
	/*
		ws := new(restful.WebService)
		ws.Route(ws.GET("/hello").To(hello))
		restful.Add(ws)

		log.Fatal(http.ListenAndServe(":8083", nil))
	*/
	handler := server.NewAPIServerHandler(
		"test-server",
		legacyscheme.Codecs,
		func(apiHandler http.Handler) http.Handler {
			return apiHandler
		},
		nil)

	testApisV1 := new(restful.WebService).Path("/apis/test/v1")
	{
		testApisV1.Route(testApisV1.GET("hello").To(
			func(req *restful.Request, resp *restful.Response) {
				resp.WriteAsJson(map[string]interface{}{"k": "v"})
			},
		)).Doc("hello endpoint")
	}

	handler.GoRestfulContainer.Add(testApisV1)

	panic(http.ListenAndServe(":8080", handler))
}
