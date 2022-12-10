package main

import (
	"fmt"
	"net"
	"time"
)

func handleClient(conn net.Conn, result chan<- error) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	defer conn.Close()
	request := make([]byte, 256)

	for {
		if _, err := conn.Read(request); err != nil {
			// conn.Write([]byte(err.Error()))
			fmt.Println("Error: ", err)
			result <- err
			break
		}

		fmt.Println(string(request))

		if _, err := conn.Write([]byte("hello, client")); err != nil {
			fmt.Println("Warning: ", err)
			result <- err
			break
		}

		request = make([]byte, 256)
	}
}

func handleResult(result chan error) {
	for err := range result {
		if err != nil {
			fmt.Println(err)
			//os.Exit(1)
		}
	}
}

func main() {
	addr := "default-route-openshift-image-registry.apps.hztt-ecp-rcp-oe19-0302.ocp.hz.nsn-rdnet.net:6445"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		fmt.Println(err)
	}

	listener, err := net.ListenTCP("tcp4", tcpAddr)
	result := make(chan error)
	go handleResult(result)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go handleClient(conn, result)
	}
}
