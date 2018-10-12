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
	"strconv"
)


// variables
var addr_listen string
var port_listen string
var err error
var nurgling_workdir string

// structures
type httpRequest struct {
	request []string
	header map[string]string
	message string
}

// functions
func parseHTTP(message_raw string) httpRequest {
	// variables that will be returned
	var parsed_request httpRequest
	var http_request []string
	var http_header_fields map[string]string
	var message_body string

	// variables that are garbage
	var message_raw_split []string
	var http_request_lines []string
	var http_header_field_split []string

	// split http header from message body and save body
	message_raw_split = strings.Split(message_raw, "\r\n\r\n")
	if len(message_raw_split) > 1 {
		message_body = message_raw_split[1]
	}

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
	parsed_request = httpRequest{
		request: http_request,
		header: http_header_fields,
		message: message_body,
	}
	return parsed_request
}

func handleHTTP(http_request httpRequest) []byte{
	// variables
	var request_method string = http_request.request[0]
	var request_resource string = http_request.request[1]
	var response_head []byte
	var response_body []byte

	response_head = append(response_head, []byte("Server: nurgling/0.1\r\n")...)

	switch request_method {
	case "GET":
		// GET request
		if rune(request_resource[len(request_resource) - 1]) == '/' {
			response_body, err = ioutil.ReadFile(nurgling_workdir + request_resource + "index.html")
			if err != nil {
				go fmt.Println("error reading " + nurgling_workdir + request_resource + "index.html :")
				go fmt.Println(err)
				response_head = append([]byte("HTTP/1.1 404 Not Found\r\n"), response_head...)
				response_body = []byte("404 stop trying\r\n")
			} else {
				go fmt.Println(nurgling_workdir + request_resource + "index.html read SUCCsesfully")
				response_head = append([]byte("HTTP/1.1 200 OK\r\n"), response_head...)
			}
		} else {
			response_body, err = ioutil.ReadFile(nurgling_workdir + request_resource)
			if err != nil {
				go fmt.Println("error reading " + nurgling_workdir + request_resource + " :")
				go fmt.Println(err)
				response_head = append([]byte("HTTP/1.1 404 Not Found\r\n"), response_head...)
				response_body = []byte("404 stop trying\r\n")
			} else {
				go fmt.Println(nurgling_workdir + request_resource + " read SUCCesfully")
				response_head = append([]byte("HTTP/1.1 200 OK\r\n"), response_head...)
			}
		}
	case "HEAD":
		// HEAD request
	case "POST":
		// POST request
	case "PUT":
		// PUT request
	case "DELETE":
		// DELETE request
	case "TRACE":
		// TRACE request
	case "OPTIONS":
		// TRACE request
	case "CONNECT":
		// CONNECT request
	case "PATCH":
		//PATCH request
	}
	response_head = append(response_head, []byte("Content-Length: " + strconv.Itoa(len(response_body)) + "\r\n\r\n")...)
	fmt.Println(string(response_head))
	return append(response_head, response_body...)
}
func main() {
	// variables
	var http_response []byte
	var http_request_parsed httpRequest

	// default options
	addr_listen = "0.0.0.0"
	port_listen = "7777"
	nurgling_workdir = "/home/karlyan/go/src/nurgling"
	//

	//index, err := ioutil.ReadFile(nurgling_workdir + "/" + "index.html")
	if err != nil {
		go fmt.Printf("error reading index.html:\n")
		go fmt.Print(err)
	} else {
		go fmt.Printf("index.html read\n")
	}
	//start the Listener for tcp on port 7777
	listen, err := net.Listen("tcp", addr_listen + ":" + port_listen)
	if err != nil {
		go fmt.Printf("error connecting to socket:\n")
		go fmt.Print(err)
	} else {
		go fmt.Printf("Opened socket on port 7777 to listen on\n")
	}
	//begin infinite serving loop
	for {
		//wait for connection
		connect, err := listen.Accept()
		//////////////////////////////FORK HERE///////////////////////////////
		if err != nil {
			go fmt.Printf("connection error:\n")
			go fmt.Print(err)
		} else {
			go fmt.Printf("started connection to %v\n", connect.RemoteAddr())
		}
		//read incomming request
		message := make([]byte, 1024)
		nbytes, err :=connect.Read(message)
		if err != nil {
			go fmt.Printf("reading error (%v bytes read):\n", nbytes)
			go fmt.Print(err)
		} else {
			go fmt.Printf("read %v bytes:\n", nbytes)
			go fmt.Print(string(message))
		}

		//respond with message
		http_request_parsed = parseHTTP(string(message))
		http_response = handleHTTP(http_request_parsed)
		nbytes, err = connect.Write(http_response)
		if err != nil {
			go fmt.Printf("response error (%v bytes were written):\n", nbytes)
			go fmt.Print(err)
		} else {
			go fmt.Printf("%v bytes written SUCCessfully\n", nbytes)
		}
		//close connection
		connect.Close()
	}
}
