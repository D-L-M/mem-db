package data

import (
	"flag"
	"os/user"
	"strconv"
	"strings"
)

// Options will be cached once they have been initially retrieved
var optionsCached = false
var cachedPort = 9999
var cachedHostname = ""
var cachedPeers = []string{}
var cachedBaseDirectory = ""

// GetOptions returns options from the application's input flags
func GetOptions() (port int, hostname string, peers []string, baseDirectory string) {

	if optionsCached {
		return cachedPort, cachedHostname, cachedPeers, cachedBaseDirectory
	}

	flag.IntVar(&port, "port", 9999, "Port on which to listen for requests")
	flag.StringVar(&hostname, "hostname", "", "Publicly accessible hostname of the instance")
	flag.StringVar(&baseDirectory, "base-directory", "", "Base directory in which to store files")

	peersString := flag.String("peers", "", "Comma-delimited list of peers serving the same database")

	flag.Parse()

	peers = strings.Split(*peersString, ",")

	if hostname == "" {
		hostname = "http://127.0.0.1:" + strconv.Itoa(port)
	}

	if baseDirectory == "" {

		user, err := user.Current()

		if err == nil {
			baseDirectory = user.HomeDir
		}

	}

	// Cache options
	cachedPort = port
	cachedHostname = hostname
	cachedPeers = peers
	cachedBaseDirectory = baseDirectory
	optionsCached = true

	return

}
