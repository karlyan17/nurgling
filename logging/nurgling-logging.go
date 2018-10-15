//nurgling-logging.go
package logging

import(
	"fmt"
	"time"
	"io/ioutil"
	"os"
)

// functions
func TimeStamp() {
	fmt.Println(time.Now().Format(time.RFC822))
}

// structsures
type Log struct {
	Log_path string
	Err_path string
}

// methods

func (log Log) LogWrite(message string, err ...error) {
	if (len(err) != 0 && err[0] != nil) {
		time_error := "[" +time.Now().Format(time.RFC822) + "] " + "ERROR: " + err[0].Error() + "\n"
		ioutil.WriteFile(log.Err_path, []byte(time_error), os.ModeAppend)
		ioutil.WriteFile(log.Log_path, []byte(time_error), os.ModeAppend)
		fmt.Fprintln(os.Stderr, time_error)
	}
	time_message := "[" + time.Now().Format(time.RFC822) + "] "  + message + "\n"
	ioutil.WriteFile(log.Log_path, []byte(time_message), os.ModeAppend)
	fmt.Printf(time_message)
}

func (log Log) LogCrit(message string, err ...error) {
	if (len(err) != 0 && err[0] != nil) {
		time_error := "[" +time.Now().Format(time.RFC822) + "] " + err[0].Error() + "\n"
		ioutil.WriteFile(log.Err_path, []byte(time_error), os.ModeAppend)
		fmt.Fprintln(os.Stderr, time_error)
	} else{
		time_message := "[" + time.Now().Format(time.RFC822) + "] "  + message + "\n"
		ioutil.WriteFile(log.Log_path, []byte(time_message), os.ModeAppend)
		fmt.Printf(time_message)
	}
}
