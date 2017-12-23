package data

import (
	"flag"
	"strings"
)

// GetOptions returns options from the application's input flags
func GetOptions() (port int, peers []string) {

	port = *flag.Int("port", 9999, "Port on which to listen for requests")
	peersString := *flag.String("peers", "", "Comma-delimited list of peers serving the same database")

	flag.Parse()

	peers = strings.Split(peersString, ",")

	return

}
