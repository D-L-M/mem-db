package jsonserver

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

var routes = map[string][]Route{}
var routesLock = sync.RWMutex{}

// RegisterRoute stores a closure to execute against a method and path
func RegisterRoute(method string, path string, middleware []Middleware, action RouteAction) {

	methods := strings.Split(strings.ToUpper(method), "|")

	for _, method := range methods {
		routesLock.Lock()
		routes[method] = append(routes[method], Route{Path: path, Action: action, Middleware: middleware})
		routesLock.Unlock()
	}

}

// Dispatch will search for and execute a route
func dispatch(request http.Request, response http.ResponseWriter, method string, path string, params string, body *[]byte) (bool, int, error) {

	routesLock.RLock()

	if methodRoutes, ok := routes[strings.ToUpper(method)]; ok {

		routesLock.RUnlock()

		for _, route := range methodRoutes {

			routeMatches, routeParams := route.MatchesPath(path)

			if routeMatches {

				queryParams, _ := url.ParseQuery(params)

				for _, middleware := range route.Middleware {

					// Execute all middleware and halt execution if one of them
					// returns FALSE
					middlewareDecision, middlewareResponseCode := middleware(&request, body, queryParams, routeParams)

					if middlewareDecision == false {
						return false, middlewareResponseCode, errors.New("Access denied to route")
					}

				}

				route.Action(&request, &response, body, queryParams, routeParams)

				return true, 0, nil

			}

		}

	} else {
		routesLock.RUnlock()
	}

	return false, 0, nil

}
