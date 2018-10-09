package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":7777")
	if err != nil {
		fmt.Printf("error connecting to socket:\n")
		fmt.Print(err)
	}
	for {
		connect, err := listen.Accept()
		if err != nil {
			fmt.Printf("connection error:\n")
			fmt.Print(err)
		}
		connect.Write([]byte("HTTP/1.1 200 OK\n\nSUCC"))
		connect.Close()
	}
}
