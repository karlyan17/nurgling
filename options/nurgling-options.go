package options

import(
	"fmt"
	"flag"
	"io/ioutil"
	"strings"
	"regexp"
	"os"
)

type options struct {
	Addr_listen string
	Port_listen string
	Ssl_port_listen string
	Workdir string
	Message_log_dir string
	Error_log_dir string
	Ssl_cert string
	Ssl_key string
	Cgi_path string
	Cgi_alias string
	Server_admin string
	Server_name string
}

func parseConfig(path *string) options {
	var opts options
	//read file
	config_file_bytes, err := ioutil.ReadFile(*path)
	if err != nil {
		fmt.Fprint(os.Stderr,"ERROR:" + err.Error() + "\n")
	} else {
		fmt.Println("config file read successfully; parsing...")
	}
	config_file := string(config_file_bytes)

	//delete comments
	comment_regexp, err := regexp.Compile("#.*\n")
	config_file = comment_regexp.ReplaceAllString(config_file, "\n")

	//clean whitespaces
	config_file = strings.Replace(config_file, " ", "", -1)
	config_file = strings.Replace(config_file, "\t", "", -1)
	config_lines := strings.Split(config_file, "\n")
	for i,line := range(config_lines) {
		if line == "" {
			continue
		}
		key_value := strings.Split(line, "=")
		switch key_value[0] {
			case "Addr_listen":
				opts.Addr_listen = key_value[1]
			case "Port_listen":
				opts.Port_listen = key_value[1]
			case "Ssl_port_listen":
				opts.Ssl_port_listen = key_value[1]
			case "Workdir":
				opts.Workdir = key_value[1]
			case "Message_log_dir":
				opts.Message_log_dir = key_value[1]
			case "Error_log_dir":
				opts.Error_log_dir = key_value[1]
			case "Ssl_cert":
				opts.Ssl_cert = key_value[1]
			case "Ssl_key":
				opts.Ssl_key = key_value[1]
			case "Cgi_path":
				opts.Cgi_path = key_value[1]
			case "Cgi_alias":
				opts.Cgi_alias = key_value[1]
			case "Server_admin":
				opts.Server_admin = key_value[1]
			case "Server_name":
				opts.Server_name = key_value[1]
			default:
				fmt.Fprintln(os.Stderr,"ERROR: error parsing line", i, ":", line, "\n")
		}
	}
	fmt.Println(opts)
	return opts
}

func Get() options {
	var config_file = flag.String("f", "nurgling.conf", "nurgling configuration file")
	flag.Parse()
	parsed_options := options{
		Addr_listen: "0.0.0.0",
		Port_listen: "7777",
		Workdir: "/home/nurgling",
		Message_log_dir: "/home/nurgling",
		Error_log_dir: "/home/nurgling",
	}
	parsed_options = parseConfig(config_file)
	return parsed_options
}
