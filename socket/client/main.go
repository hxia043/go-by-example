package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	s := "10.69.56.19"
	ip := net.ParseIP(s)

	fmt.Println(ip.String())

	addr := "default-route-openshift-image-registry.apps.hztt-ecp-rcp-oe19-0302.ocp.hz.nsn-rdnet.net:6445"
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(tcpAddr.IP, tcpAddr.Port, tcpAddr.Zone)

	tcpConn, err := net.DialTCP("tcp4", nil, tcpAddr)
	defer tcpConn.Close()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(tcpConn)

	for {
		if _, err = tcpConn.Write([]byte("hello, server")); err != nil {
			fmt.Println("Error1: ", err)
			break
		}

		reply := make([]byte, 256)
		_, err := tcpConn.Read(reply)
		if err != nil {
			fmt.Println("Error2: ", err)
			break
		}
		fmt.Println(string(reply))

		time.Sleep(30 * time.Second)
	}

	reply, err := io.ReadAll(tcpConn)
	if err != nil {
		fmt.Println("Error3: ", err)
	}

	fmt.Println(string(reply))
}
