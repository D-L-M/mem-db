package data

import (
	"flag"
	"strings"
)

// GetOptions returns options from the application's input flags
func GetOptions() (port int, peers []string) {

	portInt := flag.Int("port", 9999, "Port on which to listen for requests")
	peersString := flag.String("peers", "", "Comma-delimited list of peers serving the same database")

	flag.Parse()

	port = *portInt
	peers = strings.Split(*peersString, ",")

	return

}
