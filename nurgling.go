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
	"strings"
)


//variables
var addr_listen string
var port_listen string
var err error
var nurgling_workdir string

// functions
func parseHTTP(message_raw string) ([]string,map[string]string,string) {
	// variables that will be returned
	var http_request []string
	var http_header_fields map[string]string
	var message_body string

	// variables that are garbage
	var message_raw_split []string
	var http_request_lines []string
	var http_header_field_split []string

	// split http header from message body and save body
	message_raw_split = strings.Split(message_raw, "\r\n\r\n")
	message_body = message_raw_split[1]

	// split up the header in lines
	http_request_lines = strings.Split(message_raw_split[0], "\r\n")

	// go through the lines and 
	// 1) isolate the http request line
	// 2) make a map of the request header key/value pairs
	http_header_fields = make(map[string]string)
	for i,line := range(http_request_lines) {
		if i == 0 {
			// this is the actual meat, the http request line
			http_request = strings.Split(line, " ")
		} else {
			http_header_field_split = strings.Split(line, ": ")
			http_header_fields[http_header_field_split[0]] = http_header_field_split[1]
		}
	}

	return http_request,http_header_fields,message_body
}

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
			fmt.Printf("started connection to %v\n", connect.RemoteAddr())
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
			parseHTTP(string(message))
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
