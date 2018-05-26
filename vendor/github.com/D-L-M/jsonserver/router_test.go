package jsonserver

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Set up some routes
func testRouteSetUp() {

	RegisterRoute("GET", "/", []Middleware{}, func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
		response.Write([]byte("GET /"))
	})

	RegisterRoute("GET|PUT", "/foo", []Middleware{}, func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
		response.Write([]byte("GET|PUT /foo"))
	})

	RegisterRoute("GET", "/foo/{bar}", []Middleware{}, func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
		response.Write([]byte("GET /foo/{bar} " + routeParams["bar"]))
	})

	RegisterRoute("GET", "/foo/{bar}/:", []Middleware{}, func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
		response.Write([]byte("GET /foo/{bar}/: " + routeParams["bar"] + " " + routeParams["{catchAll}"]))
	})

	RegisterRoute("GET", "/all", []Middleware{}, func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {

		requestURL := (*request).URL.String()
		bodyString := string(*body)
		queryParam := queryParams.Get("foo")

		response.Write([]byte(requestURL + " " + bodyString + " " + queryParam))

	})

	allowMiddleware := func(request *http.Request, body *[]byte, queryParams url.Values, routeParams RouteParams) (bool, int) {
		return true, 0
	}

	denyMiddleware := func(request *http.Request, body *[]byte, queryParams url.Values, routeParams RouteParams) (bool, int) {
		return false, 401
	}

	RegisterRoute("GET", "/middleware_allow", []Middleware{allowMiddleware}, func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
		response.Write([]byte("middleware_allow"))
	})

	RegisterRoute("GET", "/middleware_deny", []Middleware{allowMiddleware, denyMiddleware}, func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
		response.Write([]byte("middleware_deny"))
	})

}

// TestRegisterRoute tests registering a route with the router
func TestRegisterRoute(t *testing.T) {

	RegisterRoute("GET", "/foo", []Middleware{}, func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	})

	routesLock.RLock()
	route := routes["GET"][0]
	routesLock.RUnlock()

	if route.Path != "/foo" {
		t.Errorf("Route path mismatch (expected: %v, actual: %v)", "/foo", route.Path)
	}

	if len(route.Middleware) != 0 {
		t.Errorf("Route middleware count mismatch (expected: %v, actual: %v)", 0, len(route.Middleware))
	}

	if route.Action == nil {
		t.Errorf("Route action missing")
	}

	testRouteTearDown()

}

// TestRegisterRouteToMultipleMethods tests registering a route with the router against multiple HTTP methods
func TestRegisterRouteToMultipleMethods(t *testing.T) {

	RegisterRoute("GET|PUT", "/bar", []Middleware{func(request *http.Request, body *[]byte, queryParams url.Values, routeParams RouteParams) (bool, int) {
		return false, 401
	}}, func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	})

	routesLock.RLock()
	getRoute := routes["GET"][0]
	putRoute := routes["PUT"][0]
	routesLock.RUnlock()

	if getRoute.Path != "/bar" {
		t.Errorf("Route path mismatch (expected: %v, actual: %v)", "/foo", getRoute.Path)
	}

	if len(getRoute.Middleware) != 1 {
		t.Errorf("Route middleware count mismatch (expected: %v, actual: %v)", 1, len(getRoute.Middleware))
	}

	if getRoute.Action == nil {
		t.Errorf("Route action missing")
	}

	if putRoute.Path != "/bar" {
		t.Errorf("Route path mismatch (expected: %v, actual: %v)", "/foo", putRoute.Path)
	}

	if len(putRoute.Middleware) != 1 {
		t.Errorf("Route middleware count mismatch (expected: %v, actual: %v)", 1, len(putRoute.Middleware))
	}

	if putRoute.Action == nil {
		t.Errorf("Route action missing")
	}

	testRouteTearDown()

}

// TestDispatchUnmatchedRoute tests dispatching a route that doesn't match
func TestDispatchUnmatchedRoute(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("GET", "http://localhost:9999/no_match", nil)
	response := httptest.NewRecorder()

	success, code, err := dispatch(request, response, "GET", "/no_match", "", &[]byte{})

	if response.Body.String() != "" {
		t.Errorf("Route erroneously executed")
	}

	if success != false {
		t.Errorf("Route erroneously executed")
	}

	if code != 0 {
		t.Errorf("Middleware incorrectly returned HTTP code")
	}

	if err != nil {
		t.Errorf("Middleware incorrectly blocked route execution")
	}

	testRouteTearDown()

}

// TestDispatchRouteForMethodWithNoRoutes tests dispatching a route for a method that has no routes
func TestDispatchRouteForMethodWithNoRoutes(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("OPTIONS", "http://localhost:9999/", nil)
	response := httptest.NewRecorder()

	success, code, err := dispatch(request, response, "OPTIONS", "/", "", &[]byte{})

	if response.Body.String() != "" {
		t.Errorf("Route erroneously executed")
	}

	if success != false {
		t.Errorf("Route erroneously executed")
	}

	if code != 0 {
		t.Errorf("Middleware incorrectly returned HTTP code")
	}

	if err != nil {
		t.Errorf("Middleware incorrectly blocked route execution")
	}

	testRouteTearDown()

}

// TestDispatchBasicRoute tests dispatching a basic route
func TestDispatchBasicRoute(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("GET", "http://localhost:9999/", nil)
	response := httptest.NewRecorder()

	success, code, err := dispatch(request, response, "GET", "/", "", &[]byte{})

	if response.Body.String() != "GET /" {
		t.Errorf("Correct route did not execute")
	}

	if success != true {
		t.Errorf("Correct route did not execute")
	}

	if code != 0 {
		t.Errorf("Middleware incorrectly returned HTTP code")
	}

	if err != nil {
		t.Errorf("Middleware incorrectly blocked route execution")
	}

	testRouteTearDown()

}

// TestDispatchMultiMethodGetRoute tests dispatching a GET route registed against multiple HTTP methods
func TestDispatchMultiMethodGetRoute(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("GET", "http://localhost:9999/foo", nil)
	response := httptest.NewRecorder()

	success, code, err := dispatch(request, response, "GET", "/foo", "", &[]byte{})

	if response.Body.String() != "GET|PUT /foo" {
		t.Errorf("Correct route did not execute")
	}

	if success != true {
		t.Errorf("Correct route did not execute")
	}

	if code != 0 {
		t.Errorf("Middleware incorrectly returned HTTP code")
	}

	if err != nil {
		t.Errorf("Middleware incorrectly blocked route execution")
	}

	testRouteTearDown()

}

// TestDispatchMultiMethodPutRoute tests dispatching a PUT route registed against multiple HTTP methods
func TestDispatchMultiMethodPutRoute(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("PUT", "http://localhost:9999/foo", nil)
	response := httptest.NewRecorder()

	success, code, err := dispatch(request, response, "PUT", "/foo", "", &[]byte{})

	if response.Body.String() != "GET|PUT /foo" {
		t.Errorf("Correct route did not execute")
	}

	if success != true {
		t.Errorf("Correct route did not execute")
	}

	if code != 0 {
		t.Errorf("Middleware incorrectly returned HTTP code")
	}

	if err != nil {
		t.Errorf("Middleware incorrectly blocked route execution")
	}

	testRouteTearDown()

}

// TestDispatchRouteWithWildcard tests dispatching a route with a regular wildcard
func TestDispatchRouteWithWildcard(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("GET", "http://localhost:9999/foo/baz", nil)
	response := httptest.NewRecorder()

	success, code, err := dispatch(request, response, "GET", "/foo/baz", "", &[]byte{})

	if response.Body.String() != "GET /foo/{bar} baz" {
		t.Errorf("Correct route did not execute")
	}

	if success != true {
		t.Errorf("Correct route did not execute")
	}

	if code != 0 {
		t.Errorf("Middleware incorrectly returned HTTP code")
	}

	if err != nil {
		t.Errorf("Middleware incorrectly blocked route execution")
	}

	testRouteTearDown()

}

// TestDispatchRouteWithWildcardAndFinalWildcard tests dispatching a route with a regular and final wildcard
func TestDispatchRouteWithWildcardAndFinalWildcard(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("GET", "http://localhost:9999/foo/bar/baz/foo", nil)
	response := httptest.NewRecorder()

	success, code, err := dispatch(request, response, "GET", "/foo/bar/baz/foo", "", &[]byte{})

	if response.Body.String() != "GET /foo/{bar}/: bar baz/foo" {
		t.Errorf("Correct route did not execute")
	}

	if success != true {
		t.Errorf("Correct route did not execute")
	}

	if code != 0 {
		t.Errorf("Middleware incorrectly returned HTTP code")
	}

	if err != nil {
		t.Errorf("Middleware incorrectly blocked route execution")
	}

	testRouteTearDown()

}

// TestDispatchRouteParams tests the parameters passed to a dispatched route
func TestDispatchRouteParams(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("GET", "http://localhost:9999/all?foo=bar", nil)
	response := httptest.NewRecorder()
	body := []byte("Body Text")

	success, code, err := dispatch(request, response, "GET", "/all", "foo=bar", &body)

	if response.Body.String() != "http://localhost:9999/all?foo=bar Body Text bar" {
		t.Errorf("Correct route did not execute")
	}

	if success != true {
		t.Errorf("Correct route did not execute")
	}

	if code != 0 {
		t.Errorf("Middleware incorrectly returned HTTP code")
	}

	if err != nil {
		t.Errorf("Middleware incorrectly blocked route execution")
	}

	testRouteTearDown()

}

// TestMiddlewarePermitsRoute tests dispatching a route with middleware works when middleware allows access
func TestMiddlewarePermitsRoute(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("GET", "http://localhost:9999/middleware_allow", nil)
	response := httptest.NewRecorder()

	success, code, err := dispatch(request, response, "GET", "/middleware_allow", "", &[]byte{})

	if response.Body.String() != "middleware_allow" {
		t.Errorf("Correct route did not execute")
	}

	if success != true {
		t.Errorf("Correct route did not execute")
	}

	if code != 0 {
		t.Errorf("Middleware incorrectly returned HTTP code")
	}

	if err != nil {
		t.Errorf("Middleware incorrectly blocked route execution")
	}

	testRouteTearDown()

}

// TestMiddlewareDeniesRoute tests dispatching a route with middleware doesn't work when middleware denies access
func TestMiddlewareDeniesRoute(t *testing.T) {

	testRouteSetUp()

	request := httptest.NewRequest("GET", "http://localhost:9999/middleware_deny", nil)
	response := httptest.NewRecorder()

	success, code, err := dispatch(request, response, "GET", "/middleware_deny", "", &[]byte{})

	if response.Body.String() != "" {
		t.Errorf("Correct route did not execute")
	}

	if success != false {
		t.Errorf("Middleware did not deny access")
	}

	if code != 401 {
		t.Errorf("Middleware did not return new HTTP code")
	}

	if err == nil {
		t.Errorf("Middleware did not report that it had denied access")
	}

	testRouteTearDown()

}

// Reset the routes
func testRouteTearDown() {

	routesLock.Lock()
	routes = map[string][]Route{}
	routesLock.Unlock()

	http.DefaultServeMux = new(http.ServeMux)

}
