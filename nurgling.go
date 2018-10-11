// # TODO
// 	- variables for options
//	- config file for options
//	- command line arguments
//	- parse HTTP requests (only deliver on GET, otherwise respond with 405)
//	- deliver custom files from workdir

package main

import (
	"fmt"
	"net"
	"io/ioutil"
)


//variables
var addr_listen string
var port_listen string
var err error
var nurgling_workdir string


func main() {
	// default options
	addr_listen = "0.0.0.0"
	port_listen = "7777"
	nurgling_workdir = "/home/karlyan/go/src/nurgling"
	//

	index, err := ioutil.ReadFile(nurgling_workdir + "/" + "index.html")
	if err != nil {
		fmt.Printf("error reading index.html:\n")
		fmt.Print(err)
	} else {
		fmt.Printf("index.html read\n")
	}
	//start the Listener for tcp on port 7777
	listen, err := net.Listen("tcp", addr_listen + ":" + port_listen)
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
		nbytes, err = connect.Write([]byte("HTTP/1.1 200 OK\n\n"))
		nbytes, err = connect.Write(index)
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
