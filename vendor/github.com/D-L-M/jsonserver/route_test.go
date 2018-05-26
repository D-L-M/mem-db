package jsonserver

import (
	"net/http"
	"net/url"
	"testing"
)

// TestNoRouteMatch tests route path not matching against a URL
func TestNoRouteMatch(t *testing.T) {

	action := func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	}

	route := Route{Path: "/foo", Action: action, Middleware: []Middleware{}}
	matches, params := route.MatchesPath("/bar")

	if matches != false {
		t.Errorf("Erroneous route match (pattern %v should not cover URL %v)", "/foo", "/bar")
	}

	if len(params) != 0 {
		t.Errorf("Param mismatch (expected: %v, actual: %v)", "[]", params)
	}

}

// TestCloseNoRouteMatch tests route path not matching against a URL that almost matches
func TestCloseNoRouteMatch(t *testing.T) {

	action := func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	}

	route := Route{Path: "/foo/bar", Action: action, Middleware: []Middleware{}}
	matches, params := route.MatchesPath("/foo")

	if matches != false {
		t.Errorf("Erroneous route match (pattern %v should not cover URL %v)", "/foo/bar", "/foo")
	}

	if len(params) != 0 {
		t.Errorf("Param mismatch (expected: %v, actual: %v)", "[]", params)
	}

}

// TestMatchesBaseURL tests route path matching against the base URL
func TestMatchesBaseURL(t *testing.T) {

	action := func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	}

	route := Route{Path: "/", Action: action, Middleware: []Middleware{}}
	matches, params := route.MatchesPath("/")

	if matches != true {
		t.Errorf("Route mismatch (pattern %v should cover URL %v)", "/", "/")
	}

	if len(params) != 0 {
		t.Errorf("Param mismatch (expected: %v, actual: %v)", "[]", params)
	}

}

// TestMatchesStaticURL tests route path matching against a full static URL
func TestMatchesStaticURL(t *testing.T) {

	action := func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	}

	route := Route{Path: "/shop/products", Action: action, Middleware: []Middleware{}}
	matches, params := route.MatchesPath("/shop/products")

	if matches != true {
		t.Errorf("Route mismatch (pattern %v should cover URL %v)", "/shop/products", "/shop/products")
	}

	if len(params) != 0 {
		t.Errorf("Param mismatch (expected: %v, actual: %v)", "[]", params)
	}

}

// TestMatchesDynamicURL tests route path matching against a full dynamic URL
func TestMatchesDynamicURL(t *testing.T) {

	action := func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	}

	route := Route{Path: "/shop/products/{id}", Action: action, Middleware: []Middleware{}}
	matches, params := route.MatchesPath("/shop/products/123")

	if matches != true {
		t.Errorf("Route mismatch (pattern %v should cover URL %v)", "/shop/products/{id}", "/shop/products/123")
	}

	if len(params) != 1 || params["id"] != "123" {
		t.Errorf("Param mismatch (expected: %v, actual: %v)", "map[id:123]", params)
	}

}

// TestMatchesMultipleDynamicURL tests route path matching against a full multiple dynamic URL
func TestMatchesMultipleDynamicURL(t *testing.T) {

	action := func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	}

	route := Route{Path: "/shop/{category}/products/{id}", Action: action, Middleware: []Middleware{}}
	matches, params := route.MatchesPath("/shop/kitchen/products/123")

	if matches != true {
		t.Errorf("Route mismatch (pattern %v should cover URL %v)", "/shop/{category}/products/{id}", "/shop/kitchen/products/123")
	}

	if len(params) != 2 || params["id"] != "123" || params["category"] != "kitchen" {
		t.Errorf("Param mismatch (expected: %v, actual: %v)", "map[category:kitchen id:123]", params)
	}

}

// TestMatchesMultipleDynamicURLWithFinalWildcard tests route path matching against a full multiple dynamic URL with a final wildcard
func TestMatchesMultipleDynamicURLWithFinalWildcard(t *testing.T) {

	action := func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	}

	route := Route{Path: "/shop/{category}/products/{id}/:", Action: action, Middleware: []Middleware{}}
	matches, params := route.MatchesPath("/shop/kitchen/products/123/foo/bar")

	if matches != true {
		t.Errorf("Route mismatch (pattern %v should cover URL %v)", "/shop/{category}/products/{id}/:", "/shop/kitchen/products/123/foo/bar")
	}

	if len(params) != 3 || params["id"] != "123" || params["category"] != "kitchen" || params["{catchAll}"] != "foo/bar" {
		t.Errorf("Param mismatch (expected: %v, actual: %v)", "map[category:kitchen id:123 {catchAll}:foo/bar]", params)
	}

}

// TestMatchesURLWithOnlyFinalWildcard tests route path matching against a URL with only a final wildcard
func TestMatchesURLWithOnlyFinalWildcard(t *testing.T) {

	action := func(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams RouteParams) {
	}

	route := Route{Path: "/:", Action: action, Middleware: []Middleware{}}
	matches, params := route.MatchesPath("/foo/bar")

	if matches != true {
		t.Errorf("Route mismatch (pattern %v should cover URL %v)", "/:", "/foo/bar")
	}

	if len(params) != 1 || params["{catchAll}"] != "foo/bar" {
		t.Errorf("Param mismatch (expected: %v, actual: %v)", "map[{catchAll}:foo/bar]", params)
	}

}

// TestNormalisePath tests normalisation of URL paths
func TestNormalisePath(t *testing.T) {

	paths := map[string]string{
		"":          "",
		"/":         "",
		"foo":       "foo",
		"/foo":      "foo",
		"foo/":      "foo",
		"/foo/":     "foo",
		"foo/bar":   "foo/bar",
		"/foo/bar/": "foo/bar",
	}

	for original, expected := range paths {

		actual := normalisePath(original)

		if actual != expected {
			t.Errorf("Path normalisation failure (expected: %v, actual: %v)", expected, actual)
		}

	}

}
