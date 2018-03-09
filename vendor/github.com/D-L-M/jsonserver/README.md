# jsonserver

[![Build Status](https://travis-ci.org/D-L-M/jsonserver.svg?branch=master)](https://travis-ci.org/D-L-M/jsonserver) [![Coverage Status](https://coveralls.io/repos/github/D-L-M/jsonserver/badge.svg?branch=master)](https://coveralls.io/github/D-L-M/jsonserver?branch=master) [![GoDoc](https://godoc.org/github.com/D-L-M/jsonserver?status.svg)](https://godoc.org/github.com/D-L-M/jsonserver)

jsonserver is a simple Golang HTTP server and routing component that can be used to create a JSON API.

An example of basic usage with a couple of routes is:

```go
package main

import (
    "net/http"
    "net/url"

    "github.com/D-L-M/jsonserver"
)

func main() {

    middleware := []jsonserver.Middleware{authenticationMiddleware}

    jsonserver.RegisterRoute("GET", "/", middleware, index)
    jsonserver.RegisterRoute("GET", "/products/{id}", middleware, products)

    jsonserver.Start(9999) // Port number

    select{}

}

// Middleware to ensure that the user is logged in
func authenticationMiddleware(request *http.Request, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) (bool, int) {

    if /* some authentication logic */ {
        return true, 0
    }

    return false, 401 // Failure and HTTP code to send back

}

// Index route
func index(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

    responseBody := jsonserver.JSON{"categories": "/categories", "basket": "/shopping-basket", "logout": "/log-out"}

    jsonserver.WriteResponse(response, &responseBody, http.StatusOK)

}

// Product route
func products(request *http.Request, response http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

    product := GetProduct(routeParams["id"])
    responseBody := jsonserver.JSON{"id": product.ID, "name": product.Name, "price": product.Price}

    jsonserver.WriteResponse(response, &responseBody, http.StatusOK)

}
```

## Middleware

Middleware (if assigned) can block execution of a route if it returns `false`, and also returns the HTTP status code that will be returned to the client.

Middleware slices are executed in the order that they are specified, so it would make sense, for example, to list generic login middleware prior to permission-checking middleware — the first one to fail will halt execution of the route and any other middleware in the slice will not be run.

## HTTP Methods

The HTTP method on which a route will listen is provided as the first argument to `jsonserver.RegisterRoute()`. To register a route against multiple HTTP methods you can provide them in the following format: `GET|OPTIONS|DELETE`.

## Route Parameters

The values of named `{wildcard}` fragments in routes are provided in the `routeParams` map, where the wildcard names (excluding curly braces) form the keys.

If a route path ends with `/:` all URL fragments at (and following) that point are collected into a route parameter named `{catchAll}` (with curly braces).

## Query Parameters

Query string parameters from a URL are made available as a `url.Values` object in the `queryParams` argument. To get a single value, the key can be passed to `queryParams.Get()`.