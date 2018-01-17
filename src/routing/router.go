package routing

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"../auth"
	"../types"
)

var routes = map[string][]types.Route{}

// Register stores a closure to execute against a method and path
func Register(method string, path string, rootUserOnly bool, route func(request *http.Request, response *http.ResponseWriter, body *[]byte, id string, routeParams url.Values)) {

	methods := strings.Split(method, "|")

	for _, method := range methods {

		routes[method] = append(routes[method], types.Route{Path: path, Route: route, RootUserOnly: rootUserOnly})

	}

}

// Dispatch will search for and execute a route
func Dispatch(request *http.Request, response *http.ResponseWriter, method string, path string, params string, id string, body *[]byte) (bool, error) {

	if methodRoutes, ok := routes[method]; ok {

		for _, route := range methodRoutes {

			routeParams, _ := url.ParseQuery(params)

			if route.Path == path || route.Path == "/*" {

				username, _, _ := auth.GetCredentials(request)
				isRootUser := username == "root"

				if route.RootUserOnly == false || isRootUser {

					route.Route(request, response, body, id, routeParams)

					return true, nil

				}

				return false, errors.New("Route is restricted to root user only")

			}

		}

	}

	return false, nil

}

// GetFirstParamValue extracts the first value for a key from a route's URL parameters
func GetFirstParamValue(params url.Values, key string, fallback string) string {

	if _, ok := params[key]; ok {

		values := params[key]

		if len(values) >= 1 {
			return values[0]
		}

	}

	return fallback

}
