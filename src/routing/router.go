package routing

import (
	"net/http"
	"strings"

	"../types"
)

var routes = map[string][]types.Route{}

// Register stores a closure to execute against a method and path
func Register(method string, path string, route func(response *http.ResponseWriter, body *[]byte, id string)) {

	methods := strings.Split(method, "|")

	for _, method := range methods {

		routes[method] = append(routes[method], types.Route{Path: path, Route: route})

	}

}

// Dispatch will search for and execute a route
func Dispatch(response *http.ResponseWriter, method string, path string, id string, body *[]byte) bool {

	if methodRoutes, ok := routes[method]; ok {

		for _, route := range methodRoutes {

			if route.Path == path || route.Path == "/*" {

				route.Route(response, body, id)

				return true

			}

		}

	}

	return false

}
