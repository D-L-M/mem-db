package jsonserver

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

var routes = map[string][]Route{}

// RegisterRoute stores a closure to execute against a method and path
func RegisterRoute(method string, path string, middleware []Middleware, action RouteAction) {

	methods := strings.Split(method, "|")

	for _, method := range methods {
		routes[method] = append(routes[method], Route{Path: path, Action: action, Middleware: middleware})
	}

}

// Dispatch will search for and execute a route
func dispatch(request *http.Request, response *http.ResponseWriter, method string, path string, params string, body *[]byte) (bool, error) {

	if methodRoutes, ok := routes[method]; ok {

		for _, route := range methodRoutes {

			routeMatches, routeParams := route.MatchesPath(path)

			// TODO: Implement a check here that works with (and extracts) wildcards
			if routeMatches {

				queryParams, _ := url.ParseQuery(params)

				for _, middleware := range route.Middleware {

					// Execute all middleware and halt execution if one of them
					// returns FALSE
					if middleware(request, body, queryParams, routeParams) == false {
						return false, errors.New("Access denied to route")
					}

				}

				route.Action(request, response, body, queryParams, routeParams)

				return true, nil

			}

		}

	}

	return false, nil

}
