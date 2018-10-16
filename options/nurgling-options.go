package options

import(
	"fmt"
	"flag"
	"io/ioutil"
	"strings"
	"regexp"
)

type options struct {
	Addr_listen string
	Port_listen string
	Workdir string
	Message_log_dir string
	Error_log_dir string
}

func parseConfig(path *string) options {
	var opts options
	//read file
	config_file_bytes, err := ioutil.ReadFile(*path)
	fmt.Println(err)
	config_file := string(config_file_bytes)

	//delete comments
	comment_regexp, err := regexp.Compile("#.*\n")
	fmt.Println(err)
	config_file = comment_regexp.ReplaceAllString(config_file, "\n")

	//clean whitespaces
	config_file = strings.Replace(config_file, " ", "", -1)
	config_file = strings.Replace(config_file, "\t", "", -1)

	config_lines := strings.Split(config_file, "\n")
	for _,line := range(config_lines) {
		key_value := strings.Split(line, "=")
		switch key_value[0] {
			case "Addr_listen":
				opts.Addr_listen = key_value[1]
			case "Port_listen":
				opts.Port_listen = key_value[1]
			case "Workdir":
				opts.Workdir = key_value[1]
			case "Message_log_dir":
				opts.Message_log_dir = key_value[1]
			case "Error_log_dir":
				opts.Error_log_dir = key_value[1]
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
