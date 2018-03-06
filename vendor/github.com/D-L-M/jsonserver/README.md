# jsonserver

jsonserver is a simple Golang TCP server and routing component that can be used to create a simple JSON API.

Simple usage is:

```go
package main

import (
    "net/http"
    "net/url"

    "github.com/D-L-M/jsonserver"
)

func main() {

    middleware := []jsonserver.Middleware{} // Optional slice of Middleware functions

    jsonserver.RegisterRoute("GET", "/products/{id}", middleware, func(request *http.Request, response *http.ResponseWriter, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) {

        jsonserver.WriteResponse(response, jsonserver.JSON{"foo": "bar", "query_params": queryParams, "route_params": routeParams}, http.StatusOK)

    })

    jsonserver.Start(9999)

    select{}

}
```

A route can listen on multiple HTTP methods by pipe-delimiting them, e.g. `GET|POST`.

Middleware functions have the signature `func(request *http.Request, body *[]byte, queryParams url.Values, routeParams jsonserver.RouteParams) bool` and will prevent the route from loading if they return `false`.