//nurgling.go


package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"mime"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"strconv"
	"nurgling/logging"
	"nurgling/options"
	"regexp"
	"time"
	"bytes"
	"syscall"
)
//#include <stdio.h>
//#include <sys/types.h>
//#include <unistd.h>
import "C"


// variables
var addr_listen string = "0.0.0.0"
var port_listen string = "7777"
var ssl_port_listen string = "8888"
var nurgling_workdir string = "/home/nurgling/www"
var error_log string
var error_log_dir string = "/home/nurgling/log"
var message_log string
var message_log_dir string = "/home/nurgling/log"
var log logging.Log
var ssl_cert string = "/home/nurgling/cert/server.crt"
var ssl_key string = "/home/nurgling/key/server.key"
var server_name string = "www.example.com"
var server_software string = "nurgling/0.1"
var cgi_path string = "/home/nurgling/www/cgi"
var cgi_alias string = "/cgi"
var server_admin string = "nobody@example.com"
var default_page string = "index.html"
var default_cgi string = "spr"
var www_user string = "nurgling"

// structures
type httpRequest struct {
	request []string
	header map[string]string
	message []byte
}

// functions
func parseHTTP(message_raw []byte) httpRequest {
	// variables that will be returned
	var parsed_request httpRequest
	var http_request []string
	var http_header_fields map[string]string
	var message_body []byte

	// variables that are garbage
	var message_raw_split []string
	var http_request_lines []string
	var http_header_field_split []string


	// split http header from message body 
	message_raw_split = strings.Split(string(message_raw), "\r\n\r\n")
	// split up the header in lines
	http_request_lines = strings.Split(message_raw_split[0], "\r\n")
	// save body
	if len(message_raw_split) > 1 {
		message_body = []byte(message_raw_split[1])
	}


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
			if len(http_header_field_split) > 1 {
				http_header_fields[http_header_field_split[0]] = http_header_field_split[1]
			}
		}
	}
	parsed_request = httpRequest{
		request: http_request,
		header: http_header_fields,
		message: message_body,
	}
	return parsed_request
}

func handleHTTP(http_request httpRequest, connect net.Conn, is_https bool) []byte{
	// variables
	var request_method string
	var request_resource string
	var response_head []byte
	var response_body []byte
	var err error

        response_head = append(response_head, []byte("Server: " + server_software + "\r\n")...)

	if len(http_request.request) > 1 {
		request_method = http_request.request[0]
		request_resource = http_request.request[1]
	} else {
		go log.LogWrite("bad request received")
		response_head = append([]byte("HTTP/1.1 400 Bad Request\r\n"), response_head...)
		response_body = []byte("400 bad request\r\n")
		return append(response_head, response_body...)
	}
	if match,_ := regexp.MatchString("^" + cgi_alias, request_resource); match {
		go log.LogWrite("CGI requested")
		var query_string string
		var path_info string
		var script_name string
		var script_path string

		// starting with request_resource = "/cgi/script.sh/additional/path?query=1"
		query_split := strings.SplitN(request_resource, "?", 2)
		// now query_split = {"/cgi/script.sh/additional/path","query=1"}
		if len(query_split) == 2 {
			query_string = query_split[1]
			//if a query exists query_string = "query=1"
		}
		cgi_path_stripped := strings.SplitN(query_split[0], cgi_alias, 2)[1]
		// from query_split[0] = "/cgi/script.sh/additional/path"
		// get cgi_path_stripped = "/script.sh/additional/path"
		script_name_split := strings.SplitN(cgi_path_stripped, "/", 3)
		// now script_name_split = {"", "script.sh", "additional/path" }
		if len(script_name_split) > 1 {
			script_name = cgi_alias + "/" + script_name_split[1]
		} else {
			script_name = cgi_alias
		}
		// now script_name = "/cgi/script.sh"
		if len(script_name_split) > 1 {
			script_path = cgi_path + "/" + script_name_split[1]
		} else {
			script_path = cgi_path
		}
		if len(script_name_split) == 3 {
			path_info = "/" + script_name_split[2]
			// if an additional path exists path_info = "/additional/path"
		}
		remote_conn := strings.Split(fmt.Sprint(connect.RemoteAddr()), ":")
		cgi_env := []string{
				"GATEWAY_INTERFACE=CGI/1.1",
				"DOCUMENT_ROOT=" + nurgling_workdir,
				"HTTP_COOKIE=",
				"HTTP_HOST=" + http_request.header["Host"],
				"HTTP_REFERER=" + http_request.header["Referer"],
				"HTTP_USER_AGENT=" + http_request.header["User-Agent"],
				"HTTP_CONNECTION=" + http_request.header["Connection"],
				//"PATH=",
				"QUERY_STRING=" + query_string,
				"REMOTE_ADDR=" + remote_conn[0],
				"REMOTE_HOST=" + remote_conn[0],
				"REMOTE_PORT=" + remote_conn[1],
				"REQUEST_METHOD=" + request_method,
				"REQUEST_URI=" + request_resource,
				"SCRIPT_FILENAME=" + script_path,
				"SCRIPT_NAME=" + script_name,
				"PATH_INFO=" + path_info,
				"SERVER_ADMIN=" + server_admin,
				"SERVER_NAME=" + server_name,
				"SERVER_SOFTWARE=" + server_software,
				}
		if is_https {
			cgi_env = append(cgi_env, []string{
							"HTTPS=on",
							"SERVER_PORT=" + ssl_port_listen,
							}...)
		} else{
			cgi_env = append(cgi_env, []string{
							"HTTPS=",
							"SERVER_PORT=" + port_listen,
							}...)
		}
		fmt.Println("CGI ENV:", cgi_env)
		script_path_info,_ := os.Stat(script_path)
		var cmd *exec.Cmd
		if script_path_info != nil && script_path_info.Mode().IsDir() {
			cmd = exec.Command(script_path + "/" + default_cgi, string(http_request.message))
		} else {
			cmd = exec.Command(script_path, string(http_request.message))
		}
		cmd.Env = cgi_env
		out,err := cmd.Output()
		//fmt.Println(string(out))
		lenOut := fmt.Sprintf("Content-Type: text/html; charset=utf-8\r\nContent-Length: %v\r\n\r\n", len(out))
		log.LogWrite("CGI request processed",err)
		response_head = append([]byte("HTTP/1.1 200 OK\r\n"), response_head...)
		response_head = append(response_head, []byte(lenOut)...)
		return append(response_head, out...)
	}
	switch request_method {
	case "GET":
		// GET request
		//if rune(request_resource[len(request_resource) - 1]) == '/' {
		//	request_resource = request_resource + "index.html"
		//}
		resource_file_info,_ := os.Stat(nurgling_workdir + request_resource)
		if resource_file_info != nil && resource_file_info.Mode().IsDir() {
			request_resource = request_resource + default_page
		}
		response_body, err = ioutil.ReadFile(nurgling_workdir + request_resource)
		if err != nil {
			go log.LogWrite(fmt.Sprint("error reading " + nurgling_workdir + request_resource + " :"), err)
			response_head = append([]byte("HTTP/1.1 404 Not Found\r\n"), response_head...)
			response_body = []byte("404 stop trying\r\n")
		} else {
			go log.LogWrite(fmt.Sprint(nurgling_workdir + request_resource + " read SUCCesfully"), err)
			request_resource_sep := strings.Split(request_resource, ".")
			request_resource_extension := mime.TypeByExtension("." + request_resource_sep[len(request_resource_sep) - 1 ])
			response_head = append(response_head, []byte("Content-Type: " + request_resource_extension + "\r\n")...)
			response_head = append([]byte("HTTP/1.1 200 OK\r\n"), response_head...)
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
	log.LogWrite(string(response_head))
	return append(response_head, response_body...)
}

func serveConnection(connect net.Conn, is_https bool) {
/*	if syscall.Getuid() == 0 {
		log.LogWrite("Running as root. Dropping priveleges to user: " + www_user)
		user, err := user.Lookup(www_user)
		if err != nil {
			log.LogWrite("user " + www_user + " not found, exiting:", err)
			os.Exit(1)
		}
		uid,err := strconv.Atoi(user.Uid)
		if err != nil {
			log.LogWrite("error in UID:", err)
			os.Exit(1)
		}
		gid,err := strconv.Atoi(user.Gid)
		if err != nil {
			log.LogWrite("error in GID:", err)
			os.Exit(1)
		}
		_, err = C.setgid(C.__gid_t(gid))
		if err != nil {
			log.LogWrite("Unable to set GID due to error:", err)
		}
		_, err = C.setuid(C.__uid_t(uid))
		if err != nil {
			log.LogWrite("Unable to set UID due to error:", err)
		}
	}
*/
	for{
		//variables
		var http_response []byte
		var http_request_parsed httpRequest
		var message []byte
		var err error

		//read incomming request
		buff := make([]byte, 1)
		//var message []byte
		for {
			_, err = connect.Read(buff)
			message = append(message, buff...)
			if err != nil {
				go log.LogWrite("error while reading http request", err)
				break
			}
			if l := len(message); l >= 4 && bytes.Equal(message[(l-4):(l)], []byte("\r\n\r\n")) {
				break
			}
		}
		//respond with message
		if err == nil {
		go log.LogWrite(fmt.Sprintf("%v: read %v bytes:\n" + string(message), connect.RemoteAddr(), len(message)), err)
			http_request_parsed = parseHTTP(message)
			if cl := http_request_parsed.header["Content-Length"]; cl != "" {
				m_len,_ := strconv.Atoi(cl)
				http_request_parsed.message = make([]byte, m_len)
				nbytes, err := connect.Read(http_request_parsed.message)
				go log.LogWrite(fmt.Sprintf("%v: read %v bytes%v\n", connect.RemoteAddr(), nbytes, http_request_parsed.message), err)
			}
			http_response = handleHTTP(http_request_parsed, connect, is_https)
			nbytes, err := connect.Write(http_response)
			connect.SetDeadline(time.Now().Add(5 * time.Minute))
			go log.LogWrite(fmt.Sprintf("%v: %v bytes written SUCCessfully", connect.RemoteAddr(), nbytes), err)
		} else {
			log.LogWrite(fmt.Sprintf("Connection to %v closed  with " + fmt.Sprint(err), connect.RemoteAddr()))
			break
		}
	}
}

func listenHTTP(listen net.Listener, done chan int) {
	log.LogWrite("HTTP listener started")
	for {
		connect, err := listen.Accept()
		connect.SetDeadline(time.Now().Add(5 * time.Minute))
		go log.LogWrite(fmt.Sprintf("started HTTP connection to %v", connect.RemoteAddr()), err)
		go serveConnection(connect, false)
	}
	done <- 0
}

func listenHTTPS(listen net.Listener, done chan int) {
	log.LogWrite("HTTPS listener started")
	for {
		connect, err := listen.Accept()
		connect.SetDeadline(time.Now().Add(5 * time.Minute))
		go log.LogWrite(fmt.Sprintf("started HTTPS connection to %v", connect.RemoteAddr()), err)
		go serveConnection(connect, true)
	}
	done <- 0
}

func main() {
	// setup
	opts := options.Get()
	addr_listen = opts.Addr_listen
	port_listen = opts.Port_listen
	ssl_port_listen = opts.Ssl_port_listen
	nurgling_workdir = opts.Workdir
	message_log_dir = opts.Message_log_dir
	error_log_dir = opts.Error_log_dir
	ssl_cert = opts.Ssl_cert
	ssl_key = opts.Ssl_key
	cgi_path = opts.Cgi_path
	cgi_alias = opts.Cgi_alias
	server_admin = opts.Server_admin
	server_name = opts.Server_name
	default_page = opts.Default_page
	default_cgi = opts.Default_cgi
	www_user = opts.Www_user

	fmt.Println("[" + logging.TimeStamp() + "]", "options parsed successfully")

	message_log = message_log_dir + "/" + logging.TimeStamp() + "_message.log"
	error_log = error_log_dir + "/" + logging.TimeStamp() + "_error.log"
	log = logging.Log {
		Log_path: message_log,
		Err_path: error_log,
	}
	f, err := os.OpenFile(message_log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	f, err = os.OpenFile(error_log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	f.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	log = logging.Log {
		Log_path: message_log,
		Err_path: error_log,
	}
	log.LogWrite("Logger started")

	//start the Listener for tcp on port 7777
	listen, err := net.Listen("tcp", addr_listen + ":" + port_listen)
	log.LogWrite("Opened socket on port " + port_listen + " to listen on", err)

	// start TLS Listener 
	ssl_cert_key,err := tls.LoadX509KeyPair(ssl_cert, ssl_key)
	log.LogWrite("Loading certificate and key", err)
	tls_config := &tls.Config{Certificates: []tls.Certificate{ssl_cert_key}}
	listen_tls,err := tls.Listen("tcp", addr_listen + ":" + ssl_port_listen, tls_config)
	log.LogWrite("Opened socket on port " + ssl_port_listen + " to listen on", err)
	if syscall.Getuid() == 0 {
		log.LogWrite("Running as root. Dropping priveleges to user: " + www_user)
		user, err := user.Lookup(www_user)
		if err != nil {
			log.LogWrite("user " + www_user + " not found, exiting:", err)
			os.Exit(1)
		}
		uid,err := strconv.Atoi(user.Uid)
		if err != nil {
			log.LogWrite("error in UID:", err)
			os.Exit(1)
		}
		gid,err := strconv.Atoi(user.Gid)
		if err != nil {
			log.LogWrite("error in GID:", err)
			os.Exit(1)
		}
		err = os.Chown(message_log, uid, gid)
		if err != nil {
			log.LogWrite("error giving message_log to " + www_user + ":", err)
			os.Exit(1)
		}
		err = os.Chown(error_log, uid, gid)
		if err != nil {
			log.LogWrite("error giving error_log to " + www_user + ":", err)
			os.Exit(1)
		}
		_, err = C.setgid(C.__gid_t(gid))
		if err != nil {
			log.LogWrite("Unable to set GID due to error:", err)
		}
		_, err = C.setuid(C.__uid_t(uid))
		if err != nil {
			log.LogWrite("Unable to set UID due to error:", err)
		}
	}
	//begin infinite serving loop
	chan_http := make(chan int)
	chan_https := make(chan int)
	go listenHTTP(listen, chan_http)
	go listenHTTPS(listen_tls, chan_https)
	<-chan_http
	<-chan_https
	for {
		//wait for connection
		//connect, err := listen.Accept()
		//go log.LogWrite(fmt.Sprintf("started connection to %v", connect.RemoteAddr()), err)
		//go serveConnection(connect)
	}
}
