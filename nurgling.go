package main

import (
	"fmt"
	"net"
)

func main() {
	//start the Listener for tcp on port 7777
	listen, err := net.Listen("tcp", ":7777")
	if err != nil {
		fmt.Printf("error connecting to socket:\n")
		fmt.Print(err)
	} else {
		fmt.Printf("Opened socket on port 7777 to listen on\n")
	}
	//begin infinite serving loop
	for {
		//wait for connection
		connect, err := listen.Accept()
		if err != nil {
			fmt.Printf("connection error:\n")
			fmt.Print(err)
		} else {
			fmt.Printf("start listening\n")
		}
		//read incomming request
		message := make([]byte, 1024)
		nbytes, err :=connect.Read(message)
		if err != nil {
			fmt.Printf("reading error (%v bytes read):\n", nbytes)
			fmt.Print(err)
		} else {
			fmt.Printf("read %v bytes:\n", nbytes)
			fmt.Print(string(message))
		}
		//respond with message
		nbytes, err = connect.Write([]byte("HTTP/1.1 200 OK\n\nSUCC"))
		if err != nil {
			fmt.Printf("response error (%v bytes were written):\n", nbytes)
			fmt.Print(err)
		} else {
			fmt.Printf("%v bytes written SUCCessfully\n", nbytes)
		}
		//close connection
		connect.Close()
	}
}
