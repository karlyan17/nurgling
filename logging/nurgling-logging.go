//nurgling-logging.go
package logging

import(
	"fmt"
	"time"
	//"io/ioutils"
)

// functions
func TimeStamp() {
	fmt.Println(time.Now().Format(time.RFC822))
}

// structsures
type Log struct {
	log_path string

}

// methods

