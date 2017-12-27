package data

import (
	"flag"
	"strconv"
	"strings"
)

// GetOptions returns options from the application's input flags
func GetOptions() (port int, hostname string, peers []string) {

	flag.IntVar(&port, "port", 9999, "Port on which to listen for requests")
	flag.StringVar(&hostname, "hostname", "", "Publicly accessible hostname of the instance")

	peersString := flag.String("peers", "", "Comma-delimited list of peers serving the same database")

	flag.Parse()

	peers = strings.Split(*peersString, ",")

	if hostname == "" {
		hostname = "http://127.0.0.1:" + strconv.Itoa(port)
	}

	return

}
