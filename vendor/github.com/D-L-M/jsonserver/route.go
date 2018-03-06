package jsonserver

import (
	"strings"
)

// Route structs define executable HTTP routes
type Route struct {
	Path       string
	Action     RouteAction
	Middleware []Middleware
}

// MatchesPath checks whether the route's path matches a given path and returns any wildcard values
func (route *Route) MatchesPath(path string) (bool, RouteParams) {

	pathFragments := strings.Split(padPathWithSlashes(path), "/")
	routePathFragments := strings.Split(padPathWithSlashes(route.Path), "/")
	wildcardValues := RouteParams{}

	if len(pathFragments) == len(routePathFragments) {

		for i, routePathFragment := range routePathFragments {

			isWildcard := strings.HasPrefix(routePathFragment, "{") && strings.HasSuffix(routePathFragment, "}")

			if isWildcard == false && pathFragments[i] != routePathFragment {
				return false, RouteParams{}
			}

			if isWildcard {
				wildcardKey := routePathFragment[1 : len(routePathFragment)-1]
				wildcardValues[wildcardKey] = pathFragments[i]
			}

		}

		return true, wildcardValues

	}

	return false, RouteParams{}

}

// padPathWithSlashes ensures that a path has leading and trailing slashes
func padPathWithSlashes(path string) string {

	if strings.HasPrefix(path, "/") == false {
		path = "/" + path
	}

	if strings.HasSuffix(path, "/") == false {
		path += "/"
	}

	return path

}
