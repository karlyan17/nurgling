//nurgling-logging.go
package logging

import(
	"fmt"
	"time"
	"os"
	"io"
)

// functions
func TimeStamp() string {
	return time.Now().Format(time.RFC3339)
}
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// structsures
type Log struct {
	Log_path string
	Err_path string
}

// methods

func (log Log) LogWrite(message string, err ...error) {
	time_message := "[" + time.Now().Format(time.RFC3339) + "] "  + message + "\n"
	WriteFile(log.Log_path, []byte(time_message), os.ModeAppend)
	fmt.Printf(time_message)
	if (len(err) != 0 && err[0] != nil) {
		time_error := "[" +time.Now().Format(time.RFC3339) + "] " + "ERROR: " + err[0].Error() + "\n"
		WriteFile(log.Err_path, []byte(time_error), os.ModeAppend)
		WriteFile(log.Log_path, []byte(time_error), os.ModeAppend)
		fmt.Fprint(os.Stderr, time_error)
	}
}
