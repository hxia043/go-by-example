package main

import (
	"fmt"
	"net/http"
	"net/rpc"
)

type Person struct {
	Name string
}

type Mock int

func (p *Mock) SayHello(args *Person, reply *string) error {
	*reply = fmt.Sprintf("Hello %s", args.Name)
	return nil
}

func main() {
	p := new(Mock)
	rpc.Register(p)
	rpc.HandleHTTP()

	err := http.ListenAndServe("default-route-openshift-image-registry.apps.hztt-ecp-rcp-oe19-0302.ocp.hz.nsn-rdnet.net:6445", nil)
	if err != nil {
		fmt.Println(err)
	}
}
