//nurgling-logging.go
package logging

import(
	"fmt"
	"time"
	"io/ioutil"
	"os"
)

// functions
func TimeStamp() string {
	return time.Now().Format(time.RFC3339)
}

// structsures
type Log struct {
	Log_path string
	Err_path string
}

// methods

func (log Log) LogWrite(message string, err ...error) {
	time_message := "[" + time.Now().Format(time.RFC3339) + "] "  + message + "\n"
	go ioutil.WriteFile(log.Log_path, []byte(time_message), os.ModeAppend)
	fmt.Printf(time_message)
	if (len(err) != 0 && err[0] != nil) {
		time_error := "[" +time.Now().Format(time.RFC3339) + "] " + "ERROR: " + err[0].Error() + "\n"
		go ioutil.WriteFile(log.Err_path, []byte(time_error), os.ModeAppend)
		go ioutil.WriteFile(log.Log_path, []byte(time_error), os.ModeAppend)
		fmt.Fprint(os.Stderr, time_error)
	}
}

func (log Log) LogCrit(message string, err ...error) {
	if (len(err) != 0 && err[0] != nil) {
		time_error := "[" +time.Now().Format(time.RFC3339) + "] " + err[0].Error() + "\n"
		ioutil.WriteFile(log.Err_path, []byte(time_error), os.ModeAppend)
		fmt.Fprintln(os.Stderr, time_error)
	} else{
		time_message := "[" + time.Now().Format(time.RFC3339) + "] "  + message + "\n"
		ioutil.WriteFile(log.Log_path, []byte(time_message), os.ModeAppend)
		fmt.Printf(time_message)
	}
}
