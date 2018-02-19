package output

import (
	"fmt"
	"time"
)

// Log outputs a message to the console
func Log(message string) {

	fmt.Println(time.Now().Format(time.RFC3339) + ": " + message)

}
