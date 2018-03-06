package jsonserver

import (
	"net/http"
	"net/url"
)

// JSON represents JSON documents in map form
type JSON map[string]interface{}

// RouteParams is an alias for a map to hold route wildcard parameters, where both keys and values will be strings
type RouteParams map[string]string

// RouteAction is a function signature for actions carried out when a route is matched
type RouteAction func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams)

// Middleware is a function signature for HTTP middleware that can be assigned routes
type Middleware func(request *http.Request, body *[]byte, queryParams url.Values, routeParams RouteParams) bool
