package output

import (
	"fmt"
	"time"

	"github.com/D-L-M/mem-db/src/data"
)

// Log outputs a message to the console
func Log(message string) {

	_, _, _, _, logMode := data.GetOptions()

	if logMode != "silent" {
		fmt.Println(time.Now().Format(time.RFC3339) + ": " + message)
	}

}
